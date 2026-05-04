package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Paciente struct {
	ID             uuid.UUID       `db:"id"`
	Nome           string          `db:"nome"`
	DataNascimento time.Time       `db:"data_nascimento"`
	AlturaCm       *int            `db:"altura_cm"`
	PesoKg         *float64        `db:"peso_kg"`
	Sexo           *string         `db:"sexo"`
	Anamnese       json.RawMessage `db:"anamnese"`
	CriadoEm      time.Time       `db:"criado_em"`
	AtualizadoEm  time.Time       `db:"atualizado_em"`
}

type CreatePacienteRequest struct {
	Nome           string          `json:"nome"`
	DataNascimento string          `json:"data_nascimento"`
	AlturaCm       *int            `json:"altura_cm,omitempty"`
	PesoKg         *float64        `json:"peso_kg,omitempty"`
	Sexo           *string         `json:"sexo,omitempty"`
	Anamnese       json.RawMessage `json:"anamnese,omitempty"`
}

type UpdatePacienteRequest struct {
	Nome           *string         `json:"nome,omitempty"`
	DataNascimento *string         `json:"data_nascimento,omitempty"`
	AlturaCm       *int            `json:"altura_cm,omitempty"`
	PesoKg         *float64        `json:"peso_kg,omitempty"`
	Sexo           *string         `json:"sexo,omitempty"`
	Anamnese       json.RawMessage `json:"anamnese,omitempty"`
}

type PacienteResponse struct {
	ID             uuid.UUID       `json:"id"`
	Nome           string          `json:"nome"`
	DataNascimento string          `json:"data_nascimento"`
	AlturaCm       *int            `json:"altura_cm,omitempty"`
	PesoKg         *float64        `json:"peso_kg,omitempty"`
	Sexo           *string         `json:"sexo,omitempty"`
	Anamnese       json.RawMessage `json:"anamnese,omitempty"`
	CriadoEm      time.Time       `json:"criado_em"`
}