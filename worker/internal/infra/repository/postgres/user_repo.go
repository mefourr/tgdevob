package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mefourr/tgdevob/worker/internal/domain"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"log/slog"
)

type Repository struct {
	client Client
}

func New(client Client) *Repository {
	return &Repository{client: client}
}

func formatQuery(_ string) string {
	panic("implement me")
}

func (d *Repository) FindByID(ctx context.Context, id string) (domain.PgUser, error) {

	q := `
		SELECT 
		    u.id, u.tg_uid, u.login, u.first_name, u.last_name
		FROM public.tg_users u
		WHERE u.tg_uid = $1
	`
	slog.DebugContext(ctx, "SQL Query", ":", q, "id:", id)

	var u domain.PgUser
	if err := d.client.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.TgUserId, &u.UserName, &u.FirstName, &u.LastName,
	); err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "SQL Query", ":", err.Error())
		return domain.PgUser{}, err
	}

	slog.DebugContext(ctx, "Retrieved entity", "user", u)
	return u, nil
}

func (d *Repository) Create(ctx context.Context, user *domain.PgUser) error {
	// TODO: do fields validation on up level
	q := `
		INSERT INTO 
		    public.tg_users (tg_uid, login, first_name, last_name)
		VALUES($1, $2, $3, $4)
		RETURNING id
	`
	slog.DebugContext(ctx, "SQL", "query", q)

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

func (d *Repository) FindAll(ctx context.Context) (u []domain.PgUser, err error) {
	q := `
		SELECT 
		    u.id, u.tg_uid, u.login, u.first_name, u.last_name
		FROM public.tg_users u
	`
	slog.DebugContext(ctx, "SQL query", ":", q)

	query, err := d.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	users := make([]domain.PgUser, 0)
	for query.Next() {
		var u domain.PgUser
		if err = query.Scan(&u.ID, &u.TgUserId, &u.UserName, &u.FirstName, &u.LastName); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = query.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (d *Repository) Update(_ context.Context, _ domain.PgUser) error {
	_ = `
		UPDATE public.tg_users
		SET tg_uid = $1, login = $2, first_name = $3, last_name = $4
		WHERE id = $5
	`

	//TODO implement me
	panic("implement me")
}

func (d *Repository) Delete(_ context.Context, _ string) error {
	_ = `
		DELETE FROM public.tg_users WHERE id = $1
	` // TODO: or add a flag isDeleted

	//TODO implement me
	panic("implement me")
}
