package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BCoaracy/pacientezero/internal/handler"
	"github.com/BCoaracy/pacientezero/internal/model"
	"github.com/google/uuid"
)

type mockUsuarioSvc struct {
	createFn  func(ctx context.Context, req model.CreateUsuarioRequest) (*model.UsuarioResponse, error)
	getByIDFn func(ctx context.Context, id uuid.UUID) (*model.UsuarioResponse, error)
	listFn    func(ctx context.Context) ([]*model.UsuarioResponse, error)
	updateFn  func(ctx context.Context, id uuid.UUID, req model.UpdateUsuarioRequest) (*model.UsuarioResponse, error)
	deleteFn  func(ctx context.Context, id uuid.UUID) error
	loginFn   func(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error)
}

func (m *mockUsuarioSvc) Create(ctx context.Context, req model.CreateUsuarioRequest) (*model.UsuarioResponse, error) {
	return m.createFn(ctx, req)
}
func (m *mockUsuarioSvc) GetByID(ctx context.Context, id uuid.UUID) (*model.UsuarioResponse, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockUsuarioSvc) List(ctx context.Context) ([]*model.UsuarioResponse, error) {
	return m.listFn(ctx)
}
func (m *mockUsuarioSvc) Update(ctx context.Context, id uuid.UUID, req model.UpdateUsuarioRequest) (*model.UsuarioResponse, error) {
	return m.updateFn(ctx, id, req)
}
func (m *mockUsuarioSvc) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}
func (m *mockUsuarioSvc) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	return m.loginFn(ctx, req)
}

func fixedUsuario() *model.UsuarioResponse {
	return &model.UsuarioResponse{
		ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Nome:     "Ana",
		Email:    "ana@example.com",
		Role:     "usuario",
		CriadoEm: time.Time{},
	}
}

func TestUsuarioHandler_Create(t *testing.T) {
	svc := &mockUsuarioSvc{
		createFn: func(_ context.Context, req model.CreateUsuarioRequest) (*model.UsuarioResponse, error) {
			return fixedUsuario(), nil
		},
	}
	h := handler.NewUsuarioHandler(svc)

	body, _ := json.Marshal(model.CreateUsuarioRequest{Nome: "Ana", Email: "ana@example.com", Senha: "secret"})
	req := httptest.NewRequest(http.MethodPost, "/usuarios", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("esperado 201, got %d", w.Code)
	}
	var resp model.UsuarioResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Email != "ana@example.com" {
		t.Errorf("email inesperado: %s", resp.Email)
	}
}

func TestUsuarioHandler_Create_BadRequest(t *testing.T) {
	h := handler.NewUsuarioHandler(&mockUsuarioSvc{})

	req := httptest.NewRequest(http.MethodPost, "/usuarios", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, got %d", w.Code)
	}
}

func TestUsuarioHandler_Login(t *testing.T) {
	svc := &mockUsuarioSvc{
		loginFn: func(_ context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
			return &model.LoginResponse{Token: "tok", User: *fixedUsuario()}, nil
		},
	}
	h := handler.NewUsuarioHandler(svc)

	body, _ := json.Marshal(model.LoginRequest{Email: "ana@example.com", Senha: "secret"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperado 200, got %d", w.Code)
	}
	var resp model.LoginResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Token != "tok" {
		t.Errorf("token inesperado: %s", resp.Token)
	}
}

func TestUsuarioHandler_Login_Unauthorized(t *testing.T) {
	svc := &mockUsuarioSvc{
		loginFn: func(_ context.Context, _ model.LoginRequest) (*model.LoginResponse, error) {
			return nil, model.ErrUnauthorized
		},
	}
	h := handler.NewUsuarioHandler(svc)

	body, _ := json.Marshal(model.LoginRequest{Email: "x@x.com", Senha: "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("esperado 401, got %d", w.Code)
	}
}

func TestUsuarioHandler_GetByID_NotFound(t *testing.T) {
	svc := &mockUsuarioSvc{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.UsuarioResponse, error) {
			return nil, model.ErrNotFound
		},
	}
	h := handler.NewUsuarioHandler(svc)

	id := uuid.New()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /usuarios/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/usuarios/"+id.String(), nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("esperado 404, got %d", w.Code)
	}
}

func TestUsuarioHandler_List(t *testing.T) {
	svc := &mockUsuarioSvc{
		listFn: func(_ context.Context) ([]*model.UsuarioResponse, error) {
			return []*model.UsuarioResponse{fixedUsuario()}, nil
		},
	}
	h := handler.NewUsuarioHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/usuarios", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperado 200, got %d", w.Code)
	}
	var resp []*model.UsuarioResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp) != 1 {
		t.Errorf("esperado 1 usuário, got %d", len(resp))
	}
}

func TestUsuarioHandler_Delete(t *testing.T) {
	svc := &mockUsuarioSvc{
		deleteFn: func(_ context.Context, _ uuid.UUID) error { return nil },
	}
	h := handler.NewUsuarioHandler(svc)

	id := uuid.New()
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /usuarios/{id}", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/usuarios/"+id.String(), nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("esperado 204, got %d", w.Code)
	}
}
