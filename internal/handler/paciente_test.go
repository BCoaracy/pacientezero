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

type mockPacienteSvc struct {
	createFn  func(ctx context.Context, req model.CreatePacienteRequest) (*model.PacienteResponse, error)
	getByIDFn func(ctx context.Context, id uuid.UUID) (*model.PacienteResponse, error)
	listFn    func(ctx context.Context) ([]*model.PacienteResponse, error)
	updateFn  func(ctx context.Context, id uuid.UUID, req model.UpdatePacienteRequest) (*model.PacienteResponse, error)
	deleteFn  func(ctx context.Context, id uuid.UUID) error
}

func (m *mockPacienteSvc) Create(ctx context.Context, req model.CreatePacienteRequest) (*model.PacienteResponse, error) {
	return m.createFn(ctx, req)
}
func (m *mockPacienteSvc) GetByID(ctx context.Context, id uuid.UUID) (*model.PacienteResponse, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockPacienteSvc) List(ctx context.Context) ([]*model.PacienteResponse, error) {
	return m.listFn(ctx)
}
func (m *mockPacienteSvc) Update(ctx context.Context, id uuid.UUID, req model.UpdatePacienteRequest) (*model.PacienteResponse, error) {
	return m.updateFn(ctx, id, req)
}
func (m *mockPacienteSvc) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}

func fixedPaciente() *model.PacienteResponse {
	return &model.PacienteResponse{
		ID:             uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		Nome:           "Carlos",
		DataNascimento: "1990-05-15",
		CriadoEm:      time.Time{},
	}
}

func TestPacienteHandler_Create(t *testing.T) {
	svc := &mockPacienteSvc{
		createFn: func(_ context.Context, _ model.CreatePacienteRequest) (*model.PacienteResponse, error) {
			return fixedPaciente(), nil
		},
	}
	h := handler.NewPacienteHandler(svc)

	body, _ := json.Marshal(model.CreatePacienteRequest{Nome: "Carlos", DataNascimento: "1990-05-15"})
	req := httptest.NewRequest(http.MethodPost, "/pacientes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("esperado 201, got %d", w.Code)
	}
	var resp model.PacienteResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Nome != "Carlos" {
		t.Errorf("nome inesperado: %s", resp.Nome)
	}
}

func TestPacienteHandler_Create_BadRequest(t *testing.T) {
	h := handler.NewPacienteHandler(&mockPacienteSvc{})

	req := httptest.NewRequest(http.MethodPost, "/pacientes", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, got %d", w.Code)
	}
}

func TestPacienteHandler_GetByID(t *testing.T) {
	svc := &mockPacienteSvc{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.PacienteResponse, error) {
			return fixedPaciente(), nil
		},
	}
	h := handler.NewPacienteHandler(svc)

	id := uuid.New()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /pacientes/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/pacientes/"+id.String(), nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperado 200, got %d", w.Code)
	}
}

func TestPacienteHandler_GetByID_NotFound(t *testing.T) {
	svc := &mockPacienteSvc{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.PacienteResponse, error) {
			return nil, model.ErrNotFound
		},
	}
	h := handler.NewPacienteHandler(svc)

	id := uuid.New()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /pacientes/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/pacientes/"+id.String(), nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("esperado 404, got %d", w.Code)
	}
}

func TestPacienteHandler_List(t *testing.T) {
	svc := &mockPacienteSvc{
		listFn: func(_ context.Context) ([]*model.PacienteResponse, error) {
			return []*model.PacienteResponse{fixedPaciente()}, nil
		},
	}
	h := handler.NewPacienteHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/pacientes", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperado 200, got %d", w.Code)
	}
	var resp []*model.PacienteResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp) != 1 {
		t.Errorf("esperado 1 paciente, got %d", len(resp))
	}
}

func TestPacienteHandler_Delete(t *testing.T) {
	svc := &mockPacienteSvc{
		deleteFn: func(_ context.Context, _ uuid.UUID) error { return nil },
	}
	h := handler.NewPacienteHandler(svc)

	id := uuid.New()
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /pacientes/{id}", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/pacientes/"+id.String(), nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("esperado 204, got %d", w.Code)
	}
}

func TestPacienteHandler_GetByID_InvalidUUID(t *testing.T) {
	h := handler.NewPacienteHandler(&mockPacienteSvc{})

	mux := http.NewServeMux()
	mux.HandleFunc("GET /pacientes/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/pacientes/not-a-uuid", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, got %d", w.Code)
	}
}
