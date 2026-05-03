//go:build integration

package repository_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/BCoaracy/pacientezero/internal/model"
	"github.com/BCoaracy/pacientezero/internal/repository"
)

func TestPacienteRepository_CRUD(t *testing.T) {
	truncateTables(t)
	repo := repository.NewPacienteRepository(testPool)
	ctx := context.Background()
	sexoM := "M"
	anamnese := json.RawMessage(`{"alergias":["dipirona"]}`)

	t.Run("Create", func(t *testing.T) {
		p := &model.Paciente{
			Nome:           "Paciente Teste",
			DataNascimento: time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
			Sexo:           &sexoM,
			Anamnese:       anamnese,
		}
		if err := repo.Create(ctx, p); err != nil {
			t.Fatalf("Create: %v", err)
		}
		if p.ID.String() == "" {
			t.Fatal("ID nao preenchido apos Create")
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		lista, _ := repo.List(ctx)
		got, err := repo.GetByID(ctx, lista[0].ID)
		if err != nil {
			t.Fatalf("GetByID: %v", err)
		}
		if got.Nome != "Paciente Teste" {
			t.Errorf("nome esperado 'Paciente Teste', got '%s'", got.Nome)
		}
	})

	t.Run("List", func(t *testing.T) {
		pacientes, err := repo.List(ctx)
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(pacientes) == 0 {
			t.Fatal("List retornou vazio")
		}
	})

	t.Run("Update", func(t *testing.T) {
		lista, _ := repo.List(ctx)
		p := lista[0]
		p.Nome = "Paciente Atualizado"
		if err := repo.Update(ctx, p); err != nil {
			t.Fatalf("Update: %v", err)
		}
		got, _ := repo.GetByID(ctx, p.ID)
		if got.Nome != "Paciente Atualizado" {
			t.Errorf("nome esperado 'Paciente Atualizado', got '%s'", got.Nome)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		lista, _ := repo.List(ctx)
		id := lista[0].ID
		if err := repo.Delete(ctx, id); err != nil {
			t.Fatalf("Delete: %v", err)
		}
		_, err := repo.GetByID(ctx, id)
		if err != model.ErrNotFound {
			t.Errorf("esperado ErrNotFound apos Delete, got %v", err)
		}
	})
}