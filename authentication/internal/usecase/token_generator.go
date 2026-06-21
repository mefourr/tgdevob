package usecase

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mefourr/tgdevob/authentication/internal/domain"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/iam/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/iamkey"
	"log/slog"
	"os"
	"time"
)

type Generator interface {
	GetOrCreateToken(ctx context.Context) (*domain.Token, error)
}

type IamTokenGenerator struct {
	token *domain.Token
}

func New() Generator {
	return &IamTokenGenerator{
		token: new(domain.Token),
	}
}

var ErrDurationIsTooLong = errors.New("generator limit (2000sec) exceeded")

func (i IamTokenGenerator) GetOrCreateToken(ctx context.Context) (*domain.Token, error) {
	if i.token.IamToken != "" &&
		time.Now().Add(30*time.Second).Before(i.token.ExpiresAt) {
		slog.DebugContext(ctx, "returning cached token", "expires_at", i.token.ExpiresAt)
		return i.token, nil
	}

	slog.InfoContext(ctx, "cached token missing or expiring soon, requesting new token")
	immToken, err := getIamToken(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get IAM token from Yandex Cloud", "err", err)
		return nil, err
	}

	// TODO: think of put key to db
	i.token.IamToken = immToken.GetIamToken()
	i.token.ExpiresAt = immToken.GetExpiresAt().AsTime()
	slog.InfoContext(ctx, "new token stored", "expires_at", i.token.ExpiresAt)
	return i.token, nil
}

var (
	keyID            = os.Getenv("ID")
	serviceAccountID = os.Getenv("SERVICE_ACCOUNT_ID")
	keyFile          = os.Getenv("AUTH_FILE")
)

// token exchange
func getIamToken(ctx context.Context) (*iam.CreateIamTokenResponse, error) {
	slog.DebugContext(ctx, "reading private key", "key_file", keyFile)
	authKey, err := readPrivateKey()
	if err != nil {
		slog.ErrorContext(ctx, "failed to read private key", "key_file", keyFile, "err", err)
		return nil, err
	}

	credentials, err := ycsdk.ServiceAccountKey(authKey)
	if err != nil {
		slog.ErrorContext(ctx, "failed to build service account credentials", "err", err)
		return nil, err
	}

	slog.DebugContext(ctx, "building Yandex Cloud SDK")
	sdk, err := ycsdk.Build(ctx, ycsdk.Config{
		Credentials: credentials,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to build Yandex Cloud SDK", "err", err)
		return nil, err
	}

	slog.DebugContext(ctx, "signing JWT token", "service_account_id", serviceAccountID, "key_id", keyID)
	token, err := signedToken()
	if err != nil {
		slog.ErrorContext(ctx, "failed to sign JWT token", "err", err)
		return nil, err
	}

	iamRequest := &iam.CreateIamTokenRequest{
		Identity: &iam.CreateIamTokenRequest_Jwt{Jwt: token},
	}

	slog.DebugContext(ctx, "requesting IAM token from Yandex Cloud")
	newKey, err := sdk.IAM().IamToken().Create(ctx, iamRequest)
	if err != nil {
		slog.ErrorContext(ctx, "Yandex Cloud IAM token request failed", "err", err)
		return nil, err
	}

	return newKey, nil
}

// JWT building.
func signedToken() (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    serviceAccountID,
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		NotBefore: jwt.NewNumericDate(time.Now().UTC()),
		Audience:  []string{"https://iam.api.cloud.yandex.net/iam/v1/tokens"},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	token.Header["kid"] = keyID

	privateKey, err := loadPrivateKey()
	if err != nil {
		return "", err
	}

	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	keyData, err := readPrivateKey()
	if err != nil {
		return nil, err
	}

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(keyData.PrivateKey))
	if err != nil {
		return nil, err
	}
	return rsaPrivateKey, nil
}

func readPrivateKey() (*iamkey.Key, error) {
	data, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	var keyData *iamkey.Key
	if err := json.Unmarshal(data, &keyData); err != nil {
		return nil, err
	}

	return keyData, nil
}
