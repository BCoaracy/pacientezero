package repository

import (
	"context"

	"github.com/BCoaracy/pacientezero/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsuarioRepository interface {
	Create(ctx context.Context, u *model.Usuario) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Usuario, error)
	GetByEmail(ctx context.Context, email string) (*model.Usuario, error)
	List(ctx context.Context) ([]*model.Usuario, error)
	Update(ctx context.Context, u *model.Usuario) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type pgUsuarioRepository struct {
	pool *pgxpool.Pool
}

func NewUsuarioRepository(pool *pgxpool.Pool) UsuarioRepository {
	return &pgUsuarioRepository{pool: pool}
}

func (r *pgUsuarioRepository) Create(ctx context.Context, u *model.Usuario) error {
	q := `
		INSERT INTO usuarios (nome, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email_verified, criado_em, atualizado_em`
	err := r.pool.QueryRow(ctx, q, u.Nome, u.Email, u.PasswordHash, u.Role).
		Scan(&u.ID, &u.EmailVerified, &u.CriadoEm, &u.AtualizadoEm)
	return mapErr(err)
}

func (r *pgUsuarioRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Usuario, error) {
	q := `SELECT id, nome, email, password_hash, role, email_verified, criado_em, atualizado_em
	      FROM usuarios WHERE id = $1`
	u := &model.Usuario{}
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.Nome, &u.Email, &u.PasswordHash,
		&u.Role, &u.EmailVerified, &u.CriadoEm, &u.AtualizadoEm,
	)
	if err != nil {
		return nil, mapErr(err)
	}
	return u, nil
}

func (r *pgUsuarioRepository) GetByEmail(ctx context.Context, email string) (*model.Usuario, error) {
	q := `SELECT id, nome, email, password_hash, role, email_verified, criado_em, atualizado_em
	      FROM usuarios WHERE email = $1`
	u := &model.Usuario{}
	err := r.pool.QueryRow(ctx, q, email).Scan(
		&u.ID, &u.Nome, &u.Email, &u.PasswordHash,
		&u.Role, &u.EmailVerified, &u.CriadoEm, &u.AtualizadoEm,
	)
	if err != nil {
		return nil, mapErr(err)
	}
	return u, nil
}

func (r *pgUsuarioRepository) List(ctx context.Context) ([]*model.Usuario, error) {
	q := `SELECT id, nome, email, password_hash, role, email_verified, criado_em, atualizado_em
	      FROM usuarios ORDER BY criado_em DESC`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.Usuario
	for rows.Next() {
		u := &model.Usuario{}
		if err := rows.Scan(
			&u.ID, &u.Nome, &u.Email, &u.PasswordHash,
			&u.Role, &u.EmailVerified, &u.CriadoEm, &u.AtualizadoEm,
		); err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, rows.Err()
}

func (r *pgUsuarioRepository) Update(ctx context.Context, u *model.Usuario) error {
	q := `
		UPDATE usuarios
		SET nome = $1, email = $2, password_hash = $3, role = $4, atualizado_em = NOW()
		WHERE id = $5
		RETURNING atualizado_em`
	err := r.pool.QueryRow(ctx, q, u.Nome, u.Email, u.PasswordHash, u.Role, u.ID).
		Scan(&u.AtualizadoEm)
	return mapErr(err)
}

func (r *pgUsuarioRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM usuarios WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}