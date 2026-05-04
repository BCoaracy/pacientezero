package repository

import (
	"context"

	"github.com/BCoaracy/pacientezero/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PacienteRepository interface {
	Create(ctx context.Context, p *model.Paciente) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Paciente, error)
	List(ctx context.Context) ([]*model.Paciente, error)
	Update(ctx context.Context, p *model.Paciente) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type pgPacienteRepository struct {
	pool *pgxpool.Pool
}

func NewPacienteRepository(pool *pgxpool.Pool) PacienteRepository {
	return &pgPacienteRepository{pool: pool}
}

func (r *pgPacienteRepository) Create(ctx context.Context, p *model.Paciente) error {
	q := `
		INSERT INTO pacientes (nome, data_nascimento, altura_cm, peso_kg, sexo, anamnese)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, criado_em, atualizado_em`
	err := r.pool.QueryRow(ctx, q,
		p.Nome, p.DataNascimento, p.AlturaCm, p.PesoKg, p.Sexo, p.Anamnese,
	).Scan(&p.ID, &p.CriadoEm, &p.AtualizadoEm)
	return mapErr(err)
}

func (r *pgPacienteRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Paciente, error) {
	q := `SELECT id, nome, data_nascimento, altura_cm, peso_kg, sexo, anamnese, criado_em, atualizado_em
	      FROM pacientes WHERE id = $1`
	p := &model.Paciente{}
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&p.ID, &p.Nome, &p.DataNascimento, &p.AlturaCm, &p.PesoKg,
		&p.Sexo, &p.Anamnese, &p.CriadoEm, &p.AtualizadoEm,
	)
	if err != nil {
		return nil, mapErr(err)
	}
	return p, nil
}

func (r *pgPacienteRepository) List(ctx context.Context) ([]*model.Paciente, error) {
	q := `SELECT id, nome, data_nascimento, altura_cm, peso_kg, sexo, anamnese, criado_em, atualizado_em
	      FROM pacientes ORDER BY criado_em DESC`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.Paciente
	for rows.Next() {
		p := &model.Paciente{}
		if err := rows.Scan(
			&p.ID, &p.Nome, &p.DataNascimento, &p.AlturaCm, &p.PesoKg,
			&p.Sexo, &p.Anamnese, &p.CriadoEm, &p.AtualizadoEm,
		); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, rows.Err()
}

func (r *pgPacienteRepository) Update(ctx context.Context, p *model.Paciente) error {
	q := `
		UPDATE pacientes
		SET nome = $1, data_nascimento = $2, altura_cm = $3, peso_kg = $4,
		    sexo = $5, anamnese = $6, atualizado_em = NOW()
		WHERE id = $7
		RETURNING atualizado_em`
	err := r.pool.QueryRow(ctx, q,
		p.Nome, p.DataNascimento, p.AlturaCm, p.PesoKg, p.Sexo, p.Anamnese, p.ID,
	).Scan(&p.AtualizadoEm)
	return mapErr(err)
}

func (r *pgPacienteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM pacientes WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}