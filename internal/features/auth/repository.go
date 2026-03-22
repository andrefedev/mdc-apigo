package auth

import (
	"context"
	"errors"

	"apigo/internal/modules/postgres"
	"apigo/internal/platforms/apperr"

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
	qry := `INSERT INTO users_codes (code, phone) VALUES(@code, @phone) returning id;`

	// CodeData
	var ref string
	if err := r.db.QueryRow(
		ctx,
		qry,
		pgx.NamedArgs{
			"code":  data.Code,
			"phone": data.Phone,
		},
	).Scan(&ref); err != nil {
		return "", apperr.Internal(op, err)
	}

	return ref, nil
}

func (r Repository) CodeDelete(ctx context.Context, ref string) (int64, error) {
	op := "Auth.Repository.CodeDelete"
	query := "DELETE FROM users_codes WHERE id = $1"

	res, err := r.db.Exec(ctx, query, ref)
	if err != nil {
		return 0, apperr.Internal(op, err)
	}

	return res.RowsAffected(), nil
}

func (r Repository) CodeSelect(ctx context.Context, ref string) (*Code, error) {
	op := "Auth.Repository.CodeSelect"

	qry := `
	SELECT
	id, code, phone,
	date_created, date_expired
	FROM users_codes WHERE id = $1
	`

	rows, err := r.db.Query(ctx, qry, ref)
	if err != nil {
		return nil, apperr.Internal(op, err)
	}
	defer rows.Close()

	raw, err := pgx.CollectOneRow[_CodeRaw](rows, pgx.RowToStructByNameLax[_CodeRaw])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFound(op, err)
		}
		return nil, apperr.Internal(op, err)
	}

	return raw.ToModel(), nil
}

func (r Repository) IdentitySelectByIdToken(ctx context.Context, idToken string) (*Identity, error) {
	op := "Auth.Repository.IdentitySelectByIdToken"

	qry := `
	SELECT
	id, idk, is_staff, is_super, is_active
	FROM users WHERE idk = $1
	`

	rows, err := r.db.Query(ctx, qry, idToken)
	if err != nil {
		return nil, apperr.Internal(op, err)
	}
	defer rows.Close()

	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[Identity])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFound(op, err)
		}
		return nil, apperr.Internal(op, err)
	}

	return &res, nil
}

func (r Repository) SessionSelectByAccessToken(ctx context.Context, accessToken string) (*Session, error) {
	return nil, nil
}
