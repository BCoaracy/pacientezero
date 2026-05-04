package service

import (
	"context"
	"testing"
	"time"

	"github.com/BCoaracy/pacientezero/internal/model"
	"github.com/google/uuid"
)

// mockPacienteRepo implementa repository.PacienteRepository para testes
type mockPacienteRepo struct {
	pacientes map[uuid.UUID]*model.Paciente
}

func newMockPacienteRepo() *mockPacienteRepo {
	return &mockPacienteRepo{pacientes: make(map[uuid.UUID]*model.Paciente)}
}

func (m *mockPacienteRepo) Create(_ context.Context, p *model.Paciente) error {
	p.ID = uuid.New()
	p.CriadoEm = time.Now()
	p.AtualizadoEm = time.Now()
	m.pacientes[p.ID] = p
	return nil
}

func (m *mockPacienteRepo) GetByID(_ context.Context, id uuid.UUID) (*model.Paciente, error) {
	p, ok := m.pacientes[id]
	if !ok {
		return nil, model.ErrNotFound
	}
	return p, nil
}

func (m *mockPacienteRepo) List(_ context.Context) ([]*model.Paciente, error) {
	result := make([]*model.Paciente, 0, len(m.pacientes))
	for _, p := range m.pacientes {
		result = append(result, p)
	}
	return result, nil
}

func (m *mockPacienteRepo) Update(_ context.Context, p *model.Paciente) error {
	if _, ok := m.pacientes[p.ID]; !ok {
		return model.ErrNotFound
	}
	p.AtualizadoEm = time.Now()
	m.pacientes[p.ID] = p
	return nil
}

func (m *mockPacienteRepo) Delete(_ context.Context, id uuid.UUID) error {
	if _, ok := m.pacientes[id]; !ok {
		return model.ErrNotFound
	}
	delete(m.pacientes, id)
	return nil
}

func TestPacienteService_Create(t *testing.T) {
	svc := NewPacienteService(newMockPacienteRepo())

	t.Run("cria paciente com sucesso", func(t *testing.T) {
		altura := 175
		peso := 70.5
		sexo := "M"
		resp, err := svc.Create(context.Background(), model.CreatePacienteRequest{
			Nome:           "João Souza",
			DataNascimento: "1990-05-15",
			AlturaCm:       &altura,
			PesoKg:         &peso,
			Sexo:           &sexo,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Nome != "João Souza" {
			t.Errorf("nome incorreto: %s", resp.Nome)
		}
		if resp.DataNascimento != "1990-05-15" {
			t.Errorf("data incorreta: %s", resp.DataNascimento)
		}
	})

	t.Run("rejeita nome vazio", func(t *testing.T) {
		_, err := svc.Create(context.Background(), model.CreatePacienteRequest{
			DataNascimento: "1990-01-01",
		})
		if err != model.ErrBadRequest {
			t.Errorf("esperava ErrBadRequest, got %v", err)
		}
	})

	t.Run("rejeita data de nascimento inválida", func(t *testing.T) {
		_, err := svc.Create(context.Background(), model.CreatePacienteRequest{
			Nome:           "Teste",
			DataNascimento: "não-é-data",
		})
		if err != model.ErrBadRequest {
			t.Errorf("esperava ErrBadRequest, got %v", err)
		}
	})
}

func TestPacienteService_GetByID(t *testing.T) {
	svc := NewPacienteService(newMockPacienteRepo())

	resp, _ := svc.Create(context.Background(), model.CreatePacienteRequest{
		Nome: "Maria", DataNascimento: "1985-03-20",
	})

	t.Run("retorna paciente existente", func(t *testing.T) {
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

func TestPacienteService_Update(t *testing.T) {
	svc := NewPacienteService(newMockPacienteRepo())

	resp, _ := svc.Create(context.Background(), model.CreatePacienteRequest{
		Nome: "Pedro", DataNascimento: "2000-07-10",
	})

	novoNome := "Pedro Atualizado"
	novaData := "2000-07-11"
	got, err := svc.Update(context.Background(), resp.ID, model.UpdatePacienteRequest{
		Nome:           &novoNome,
		DataNascimento: &novaData,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Nome != novoNome {
		t.Errorf("nome não atualizado: %s", got.Nome)
	}
	if got.DataNascimento != novaData {
		t.Errorf("data não atualizada: %s", got.DataNascimento)
	}
}

func TestPacienteService_Update_DataInvalida(t *testing.T) {
	svc := NewPacienteService(newMockPacienteRepo())

	resp, _ := svc.Create(context.Background(), model.CreatePacienteRequest{
		Nome: "Teste", DataNascimento: "1995-01-01",
	})

	dataInvalida := "abc"
	_, err := svc.Update(context.Background(), resp.ID, model.UpdatePacienteRequest{
		DataNascimento: &dataInvalida,
	})
	if err != model.ErrBadRequest {
		t.Errorf("esperava ErrBadRequest, got %v", err)
	}
}

func TestPacienteService_Delete(t *testing.T) {
	svc := NewPacienteService(newMockPacienteRepo())

	resp, _ := svc.Create(context.Background(), model.CreatePacienteRequest{
		Nome: "Deletar", DataNascimento: "1980-12-25",
	})

	if err := svc.Delete(context.Background(), resp.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := svc.Delete(context.Background(), resp.ID); err != model.ErrNotFound {
		t.Errorf("esperava ErrNotFound na segunda deleção, got %v", err)
	}
}

func TestPacienteService_List(t *testing.T) {
	svc := NewPacienteService(newMockPacienteRepo())

	for i := 0; i < 3; i++ {
		svc.Create(context.Background(), model.CreatePacienteRequest{
			Nome:           "Paciente",
			DataNascimento: "1990-01-01",
		})
	}

	result, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("esperava 3 pacientes, got %d", len(result))
	}
}
