package auth

import (
	"apigo/internal/modules/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *postgres.Pgdb
}

func NewRepository(db *postgres.Pgdb) *Repository {
	return &Repository{db: db}
}

func (r Repository) CodeInsert(ctx context.Context, data *CodeInsertData) (string, error) {
	op := "Auth.Repository.CodeInsert"
	qry := `INSERT INTO users_codes_original (code, phone) VALUES (@code, @phone) RETURNING id;`

	var ref string
	if err := r.db.QueryRow(
		ctx,
		qry,
		pgx.NamedArgs{
			"code":  data.Code,
			"phone": data.Phone,
		},
	).Scan(&ref); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return ref, nil
}

func (r Repository) CodeDelete(ctx context.Context, ref string) (int64, error) {
	op := "Auth.Repository.CodeDelete"
	qry := `DELETE FROM users_codes_original WHERE id = $1`

	res, err := r.db.Exec(ctx, qry, ref)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return res.RowsAffected(), nil
}

func (r Repository) CodeSelect(ctx context.Context, ref string) (*Code, error) {
	const op = "Auth.Repository.CodeSelect"
	qry := `SELECT id, code, phone, date_created, date_expired FROM users_codes_original WHERE id = $1`

	rows, err := r.db.Query(ctx, qry, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	raw, err := pgx.CollectOneRow[_CodeRaw](rows, pgx.RowToStructByNameLax[_CodeRaw])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, WrapCodeNotFound(err))
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return raw.ToModel(), nil
}

func (r Repository) UserRefByPhone(ctx context.Context, phone string) (string, error) {
	op := "Auth.Repository.UserRefByPhone"

	qry := `
	SELECT id
	FROM users
	WHERE phone = $1
	LIMIT 1
	`

	var userRef string
	if err := r.db.QueryRow(ctx, qry, phone).Scan(&userRef); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, WrapUserNotFound(err))
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userRef, nil
}

func (r Repository) SessionInsert(ctx context.Context, data *SessionInsertData) (string, error) {
	op := "Auth.Repository.SessionInsert"
	qry := `INSERT INTO auth_sessions (uid, token_hash) VALUES (@uid, @token_hash)`

	var ref string
	if err := r.db.QueryRow(
		ctx,
		qry,
		pgx.NamedArgs{
			"uid":        data.UserRef,
			"token_hash": data.TokenHash,
		},
	).Scan(&ref); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return ref, nil
}

func (r Repository) SessionSelect(ctx context.Context, ref string) (*Session, error) {
	const op = "Auth.Repository.SessionSelect"

	qry := `SELECT id, uid, token_hash, date_expired, date_created, date_revoked FROM auth_sessions WHERE id = $1`

	rows, err := r.db.Query(ctx, qry, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[Session])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, WrapSessionNotFound(err))
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return new(result), nil
}

func (r Repository) SessionSelectByToken(ctx context.Context, token string) (*Session, error) {
	const op = "Auth.Repository.SessionSelectByToken"
	qry := `SELECT id, uid, token_hash, date_expired, date_created, date_revoked FROM auth_sessions WHERE token_hash = $1`

	rows, err := r.db.Query(ctx, qry, token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[Session])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = WrapSessionNotFound(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return new(result), nil
}
