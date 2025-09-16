package userinfo

import (
	"context"
	"errors"
	"fmt"
	"gihub.com/mefourr/tgdevob/worker/internal/storage/user/posgresql"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserInfo interface {
	Create(ctx context.Context, user *posgresql.User) error
	FindAll(ctx context.Context) (u []posgresql.User, err error)
	FindOne(ctx context.Context, id string) (posgresql.User, error)
	Update(ctx context.Context, user posgresql.User) error
	Delete(ctx context.Context, id string) error
}

type db struct {
	client posgresql.Client
}

func (d db) Create(ctx context.Context, user *posgresql.User) error {
	// TODO: do fields validation on up level
	q := `
		INSERT INTO 
		    public.tg_users (tg_uid, login, first_name, last_name)
		VALUES($1, $2, $3, $4)
		RETURNING id`

	if err := d.client.QueryRow(
		ctx, q, user.TgUserId, user.UserName, user.FirstName, user.LastName,
	).Scan(&user.ID); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			return errors.New(fmt.Sprintf("SQL Error: Message: %s, Detail: %s, Where: %s, SQL State: %s, Code: %s", pgError.Message, pgError.Detail, pgError.Where, pgError.SQLState(), pgError.Code))
		}
		return err
	}
	return nil
}

func (d db) FindAll(ctx context.Context) (u []posgresql.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (d db) FindOne(ctx context.Context, id string) (posgresql.User, error) {
	//TODO implement me
	panic("implement me")
}

func (d db) Update(ctx context.Context, user posgresql.User) error {
	//TODO implement me
	panic("implement me")
}

func (d db) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func New(client posgresql.Client) UserInfo {
	return &db{client: client}
}
