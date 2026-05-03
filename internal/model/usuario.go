package model

import (
	"time"

	"github.com/google/uuid"
)

type Usuario struct {
	ID            uuid.UUID `db:"id"`
	Nome          string    `db:"nome"`
	Email         string    `db:"email"`
	PasswordHash  string    `db:"password_hash"`
	Role          string    `db:"role"`
	EmailVerified bool      `db:"email_verified"`
	CriadoEm     time.Time `db:"criado_em"`
	AtualizadoEm time.Time `db:"atualizado_em"`
}

type CreateUsuarioRequest struct {
	Nome  string `json:"nome"`
	Email string `json:"email"`
	Senha string `json:"senha"`
}

type UpdateUsuarioRequest struct {
	Nome  *string `json:"nome,omitempty"`
	Email *string `json:"email,omitempty"`
	Senha *string `json:"senha,omitempty"`
}

type UsuarioResponse struct {
	ID            uuid.UUID `json:"id"`
	Nome          string    `json:"nome"`
	Email         string    `json:"email"`
	Role          string    `json:"role"`
	EmailVerified bool      `json:"email_verified"`
	CriadoEm     time.Time `json:"criado_em"`
}