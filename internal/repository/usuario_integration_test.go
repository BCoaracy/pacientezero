//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/BCoaracy/pacientezero/internal/model"
	"github.com/BCoaracy/pacientezero/internal/repository"
)

func TestUsuarioRepository_CRUD(t *testing.T) {
	truncateTables(t)
	repo := repository.NewUsuarioRepository(testPool)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		u := &model.Usuario{
			Nome:         "Test User",
			Email:        "test@example.com",
			PasswordHash: "$2a$12$hashdummy",
			Role:         "usuario",
		}
		if err := repo.Create(ctx, u); err != nil {
			t.Fatalf("Create: %v", err)
		}
		if u.ID.String() == "" {
			t.Fatal("ID nao preenchido apos Create")
		}
	})

	t.Run("GetByEmail", func(t *testing.T) {
		u, err := repo.GetByEmail(ctx, "test@example.com")
		if err != nil {
			t.Fatalf("GetByEmail: %v", err)
		}
		if u.Nome != "Test User" {
			t.Errorf("nome esperado 'Test User', got '%s'", u.Nome)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		u, _ := repo.GetByEmail(ctx, "test@example.com")
		got, err := repo.GetByID(ctx, u.ID)
		if err != nil {
			t.Fatalf("GetByID: %v", err)
		}
		if got.Email != u.Email {
			t.Errorf("email esperado '%s', got '%s'", u.Email, got.Email)
		}
	})

	t.Run("List", func(t *testing.T) {
		usuarios, err := repo.List(ctx)
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(usuarios) == 0 {
			t.Fatal("List retornou vazio")
		}
	})

	t.Run("Update", func(t *testing.T) {
		u, _ := repo.GetByEmail(ctx, "test@example.com")
		u.Nome = "Test Atualizado"
		if err := repo.Update(ctx, u); err != nil {
			t.Fatalf("Update: %v", err)
		}
		got, _ := repo.GetByID(ctx, u.ID)
		if got.Nome != "Test Atualizado" {
			t.Errorf("nome esperado 'Test Atualizado', got '%s'", got.Nome)
		}
	})

	t.Run("ConflictEmail", func(t *testing.T) {
		u2 := &model.Usuario{
			Nome: "Outro", Email: "test@example.com",
			PasswordHash: "$2a$12$hashdummy2", Role: "usuario",
		}
		err := repo.Create(ctx, u2)
		if err != model.ErrConflict {
			t.Errorf("esperado ErrConflict, got %v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		u, _ := repo.GetByEmail(ctx, "test@example.com")
		if err := repo.Delete(ctx, u.ID); err != nil {
			t.Fatalf("Delete: %v", err)
		}
		_, err := repo.GetByID(ctx, u.ID)
		if err != model.ErrNotFound {
			t.Errorf("esperado ErrNotFound apos Delete, got %v", err)
		}
	})
}