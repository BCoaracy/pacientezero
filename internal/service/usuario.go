package service

import (
	"context"
	"strings"
	"time"

	"github.com/BCoaracy/pacientezero/internal/model"
	"github.com/BCoaracy/pacientezero/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

type UsuarioService interface {
	Create(ctx context.Context, req model.CreateUsuarioRequest) (*model.UsuarioResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.UsuarioResponse, error)
	List(ctx context.Context) ([]*model.UsuarioResponse, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateUsuarioRequest) (*model.UsuarioResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error)
}

type usuarioService struct {
	repo      repository.UsuarioRepository
	jwtSecret []byte
}

func NewUsuarioService(repo repository.UsuarioRepository, jwtSecret []byte) UsuarioService {
	return &usuarioService{repo: repo, jwtSecret: jwtSecret}
}

func (s *usuarioService) Create(ctx context.Context, req model.CreateUsuarioRequest) (*model.UsuarioResponse, error) {
	if req.Nome == "" || req.Email == "" || req.Senha == "" {
		return nil, model.ErrBadRequest
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Senha), bcryptCost)
	if err != nil {
		return nil, err
	}

	u := &model.Usuario{
		Nome:         req.Nome,
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: string(hash),
		Role:         "usuario",
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	return toUsuarioResponse(u), nil
}

func (s *usuarioService) GetByID(ctx context.Context, id uuid.UUID) (*model.UsuarioResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toUsuarioResponse(u), nil
}

func (s *usuarioService) List(ctx context.Context) ([]*model.UsuarioResponse, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*model.UsuarioResponse, len(users))
	for i, u := range users {
		result[i] = toUsuarioResponse(u)
	}
	return result, nil
}

func (s *usuarioService) Update(ctx context.Context, id uuid.UUID, req model.UpdateUsuarioRequest) (*model.UsuarioResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Nome != nil {
		u.Nome = *req.Nome
	}
	if req.Email != nil {
		u.Email = strings.ToLower(strings.TrimSpace(*req.Email))
	}
	if req.Senha != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Senha), bcryptCost)
		if err != nil {
			return nil, err
		}
		u.PasswordHash = string(hash)
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	return toUsuarioResponse(u), nil
}

func (s *usuarioService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *usuarioService) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	if req.Email == "" || req.Senha == "" {
		return nil, model.ErrBadRequest
	}

	u, err := s.repo.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		return nil, model.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Senha)); err != nil {
		return nil, model.ErrUnauthorized
	}

	token, err := s.generateJWT(u)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User:  *toUsuarioResponse(u),
	}, nil
}

func (s *usuarioService) generateJWT(u *model.Usuario) (string, error) {
	claims := jwt.MapClaims{
		"sub":  u.ID.String(),
		"role": u.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func toUsuarioResponse(u *model.Usuario) *model.UsuarioResponse {
	return &model.UsuarioResponse{
		ID:            u.ID,
		Nome:          u.Nome,
		Email:         u.Email,
		Role:          u.Role,
		EmailVerified: u.EmailVerified,
		CriadoEm:     u.CriadoEm,
	}
}
