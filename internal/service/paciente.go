package service

import (
	"context"
	"time"

	"github.com/BCoaracy/pacientezero/internal/model"
	"github.com/BCoaracy/pacientezero/internal/repository"
	"github.com/google/uuid"
)

type PacienteService interface {
	Create(ctx context.Context, req model.CreatePacienteRequest) (*model.PacienteResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.PacienteResponse, error)
	List(ctx context.Context) ([]*model.PacienteResponse, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdatePacienteRequest) (*model.PacienteResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type pacienteService struct {
	repo repository.PacienteRepository
}

func NewPacienteService(repo repository.PacienteRepository) PacienteService {
	return &pacienteService{repo: repo}
}

func (s *pacienteService) Create(ctx context.Context, req model.CreatePacienteRequest) (*model.PacienteResponse, error) {
	if req.Nome == "" || req.DataNascimento == "" {
		return nil, model.ErrBadRequest
	}

	dataNasc, err := time.Parse("2006-01-02", req.DataNascimento)
	if err != nil {
		return nil, model.ErrBadRequest
	}

	p := &model.Paciente{
		Nome:           req.Nome,
		DataNascimento: dataNasc,
		AlturaCm:       req.AlturaCm,
		PesoKg:         req.PesoKg,
		Sexo:           req.Sexo,
		Anamnese:       req.Anamnese,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return toPacienteResponse(p), nil
}

func (s *pacienteService) GetByID(ctx context.Context, id uuid.UUID) (*model.PacienteResponse, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toPacienteResponse(p), nil
}

func (s *pacienteService) List(ctx context.Context) ([]*model.PacienteResponse, error) {
	pacientes, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*model.PacienteResponse, len(pacientes))
	for i, p := range pacientes {
		result[i] = toPacienteResponse(p)
	}
	return result, nil
}

func (s *pacienteService) Update(ctx context.Context, id uuid.UUID, req model.UpdatePacienteRequest) (*model.PacienteResponse, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Nome != nil {
		p.Nome = *req.Nome
	}
	if req.DataNascimento != nil {
		dataNasc, err := time.Parse("2006-01-02", *req.DataNascimento)
		if err != nil {
			return nil, model.ErrBadRequest
		}
		p.DataNascimento = dataNasc
	}
	if req.AlturaCm != nil {
		p.AlturaCm = req.AlturaCm
	}
	if req.PesoKg != nil {
		p.PesoKg = req.PesoKg
	}
	if req.Sexo != nil {
		p.Sexo = req.Sexo
	}
	if req.Anamnese != nil {
		p.Anamnese = req.Anamnese
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}

	return toPacienteResponse(p), nil
}

func (s *pacienteService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func toPacienteResponse(p *model.Paciente) *model.PacienteResponse {
	return &model.PacienteResponse{
		ID:             p.ID,
		Nome:           p.Nome,
		DataNascimento: p.DataNascimento.Format("2006-01-02"),
		AlturaCm:       p.AlturaCm,
		PesoKg:         p.PesoKg,
		Sexo:           p.Sexo,
		Anamnese:       p.Anamnese,
		CriadoEm:      p.CriadoEm,
	}
}
