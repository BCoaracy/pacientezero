package service

import (
	"context"
	"testing"
	"time"

	"github.com/BCoaracy/pacientezero/internal/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// mockUsuarioRepo implementa repository.UsuarioRepository para testes
type mockUsuarioRepo struct {
	usuarios map[uuid.UUID]*model.Usuario
	byEmail  map[string]*model.Usuario
}

func newMockUsuarioRepo() *mockUsuarioRepo {
	return &mockUsuarioRepo{
		usuarios: make(map[uuid.UUID]*model.Usuario),
		byEmail:  make(map[string]*model.Usuario),
	}
}

func (m *mockUsuarioRepo) Create(_ context.Context, u *model.Usuario) error {
	if _, exists := m.byEmail[u.Email]; exists {
		return model.ErrConflict
	}
	u.ID = uuid.New()
	u.CriadoEm = time.Now()
	u.AtualizadoEm = time.Now()
	m.usuarios[u.ID] = u
	m.byEmail[u.Email] = u
	return nil
}

func (m *mockUsuarioRepo) GetByID(_ context.Context, id uuid.UUID) (*model.Usuario, error) {
	u, ok := m.usuarios[id]
	if !ok {
		return nil, model.ErrNotFound
	}
	return u, nil
}

func (m *mockUsuarioRepo) GetByEmail(_ context.Context, email string) (*model.Usuario, error) {
	u, ok := m.byEmail[email]
	if !ok {
		return nil, model.ErrNotFound
	}
	return u, nil
}

func (m *mockUsuarioRepo) List(_ context.Context) ([]*model.Usuario, error) {
	result := make([]*model.Usuario, 0, len(m.usuarios))
	for _, u := range m.usuarios {
		result = append(result, u)
	}
	return result, nil
}

func (m *mockUsuarioRepo) Update(_ context.Context, u *model.Usuario) error {
	if _, ok := m.usuarios[u.ID]; !ok {
		return model.ErrNotFound
	}
	u.AtualizadoEm = time.Now()
	m.usuarios[u.ID] = u
	m.byEmail[u.Email] = u
	return nil
}

func (m *mockUsuarioRepo) Delete(_ context.Context, id uuid.UUID) error {
	u, ok := m.usuarios[id]
	if !ok {
		return model.ErrNotFound
	}
	delete(m.byEmail, u.Email)
	delete(m.usuarios, id)
	return nil
}

var testJWTSecret = []byte("secret-para-testes")

func TestUsuarioService_Create(t *testing.T) {
	svc := NewUsuarioService(newMockUsuarioRepo(), testJWTSecret)

	t.Run("cria usuário com sucesso", func(t *testing.T) {
		resp, err := svc.Create(context.Background(), model.CreateUsuarioRequest{
			Nome:  "Ana Silva",
			Email: "ANA@EXAMPLE.COM",
			Senha: "senha123",
		})
		if err != nil {
			t.Fatalf("esperava nil, got %v", err)
		}
		if resp.Email != "ana@example.com" {
			t.Errorf("email deveria ser lowercase: %s", resp.Email)
		}
		if resp.Role != "usuario" {
			t.Errorf("role padrão deveria ser 'usuario': %s", resp.Role)
		}
	})

	t.Run("rejeita campos obrigatórios vazios", func(t *testing.T) {
		_, err := svc.Create(context.Background(), model.CreateUsuarioRequest{Nome: "X"})
		if err != model.ErrBadRequest {
			t.Errorf("esperava ErrBadRequest, got %v", err)
		}
	})
}

func TestUsuarioService_GetByID(t *testing.T) {
	svc := NewUsuarioService(newMockUsuarioRepo(), testJWTSecret)

	resp, _ := svc.Create(context.Background(), model.CreateUsuarioRequest{
		Nome: "Bob", Email: "bob@example.com", Senha: "abc123",
	})

	t.Run("retorna usuário existente", func(t *testing.T) {
		got, err := svc.GetByID(context.Background(), resp.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != resp.ID {
			t.Errorf("IDs divergem")
		}
	})

	t.Run("retorna ErrNotFound para ID inexistente", func(t *testing.T) {
		_, err := svc.GetByID(context.Background(), uuid.New())
		if err != model.ErrNotFound {
			t.Errorf("esperava ErrNotFound, got %v", err)
		}
	})
}

func TestUsuarioService_Update(t *testing.T) {
	svc := NewUsuarioService(newMockUsuarioRepo(), testJWTSecret)

	resp, _ := svc.Create(context.Background(), model.CreateUsuarioRequest{
		Nome: "Carlos", Email: "carlos@example.com", Senha: "senha123",
	})

	novoNome := "Carlos Atualizado"
	got, err := svc.Update(context.Background(), resp.ID, model.UpdateUsuarioRequest{
		Nome: &novoNome,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Nome != novoNome {
		t.Errorf("nome não atualizado: %s", got.Nome)
	}
}

func TestUsuarioService_Delete(t *testing.T) {
	svc := NewUsuarioService(newMockUsuarioRepo(), testJWTSecret)

	resp, _ := svc.Create(context.Background(), model.CreateUsuarioRequest{
		Nome: "Deletar", Email: "del@example.com", Senha: "senha123",
	})

	if err := svc.Delete(context.Background(), resp.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := svc.Delete(context.Background(), resp.ID); err != model.ErrNotFound {
		t.Errorf("esperava ErrNotFound na segunda deleção, got %v", err)
	}
}

func TestUsuarioService_Login(t *testing.T) {
	svc := NewUsuarioService(newMockUsuarioRepo(), testJWTSecret)

	svc.Create(context.Background(), model.CreateUsuarioRequest{
		Nome: "Login User", Email: "login@example.com", Senha: "correta123",
	})

	t.Run("login com credenciais corretas retorna token", func(t *testing.T) {
		resp, err := svc.Login(context.Background(), model.LoginRequest{
			Email: "login@example.com",
			Senha: "correta123",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Token == "" {
			t.Error("token não deveria ser vazio")
		}
	})

	t.Run("login com senha errada retorna ErrUnauthorized", func(t *testing.T) {
		_, err := svc.Login(context.Background(), model.LoginRequest{
			Email: "login@example.com",
			Senha: "errada",
		})
		if err != model.ErrUnauthorized {
			t.Errorf("esperava ErrUnauthorized, got %v", err)
		}
	})

	t.Run("login com email inexistente retorna ErrUnauthorized", func(t *testing.T) {
		_, err := svc.Login(context.Background(), model.LoginRequest{
			Email: "naoexiste@example.com",
			Senha: "qualquer",
		})
		if err != model.ErrUnauthorized {
			t.Errorf("esperava ErrUnauthorized, got %v", err)
		}
	})
}

func TestUsuarioService_PasswordNotExposedInResponse(t *testing.T) {
	repo := newMockUsuarioRepo()
	svc := NewUsuarioService(repo, testJWTSecret)

	svc.Create(context.Background(), model.CreateUsuarioRequest{
		Nome: "Teste", Email: "teste@example.com", Senha: "senha_secreta",
	})

	// UsuarioResponse não tem campo PasswordHash — verificamos via repo diretamente
	u, _ := repo.GetByEmail(context.Background(), "teste@example.com")
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("senha_secreta")); err != nil {
		t.Error("hash bcrypt inválido no repositório")
	}
}
