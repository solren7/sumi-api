package services

import (
	"context"
	"errors"
	"strings"

	"fiber/config"
	"fiber/internal/repository/dbgen"
	"fiber/pkg/utils"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	q   *dbgen.Queries
	cfg *config.Config
}

func NewAuthService(q *dbgen.Queries, cfg *config.Config) *AuthService {
	return &AuthService{q: q, cfg: cfg}
}

func (s *AuthService) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	_, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

type RegisterInput struct {
	Email    string
	Password string
	Username string
}

type AuthOutput struct {
	Token string
	User  *dbgen.User
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthOutput, error) {
	if input.Username == "" {
		parts := strings.Split(input.Email, "@")
		if len(parts) > 0 {
			input.Username = parts[0]
		} else {
			input.Username = input.Email
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.q.CreateUser(ctx, dbgen.CreateUserParams{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(utils.UUIDToString(user.ID), user.Email, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		Token: token,
		User:  &user,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthOutput, error) {
	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(utils.UUIDToString(user.ID), user.Email, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		Token: token,
		User:  &user,
	}, nil
}
