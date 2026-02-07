package handlers

import (
	"context"

	"fiber/internal/repository/dbgen"
	"fiber/internal/services"
	"fiber/pkg/errorx"

	"github.com/gofiber/fiber/v3"
)

type CheckEmailRequest struct {
	Email string `json:"email"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username,omitempty"` // Optional for login
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *dbgen.User `json:"user"`
}

func (h *Handler) CheckEmail(c fiber.Ctx) error {
	req := new(CheckEmailRequest)
	if err := c.Bind().Body(req); err != nil {
		return errorx.ErrParamsInvalid
	}

	if req.Email == "" {
		return errorx.New(fiber.StatusBadRequest, "Email is required")
	}

	exists, err := h.S.Auth.CheckEmailExists(context.Background(), req.Email)
	if err != nil {
		return err // Middleware will handle it
	}

	return c.JSON(fiber.Map{"exists": exists})
}

func (h *Handler) Register(c fiber.Ctx) error {
	req := new(AuthRequest)
	if err := c.Bind().Body(req); err != nil {
		return errorx.ErrParamsInvalid
	}

	if req.Email == "" || req.Password == "" {
		return errorx.New(fiber.StatusBadRequest, "Email and password are required")
	}

	output, err := h.S.Auth.Register(context.Background(), services.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(AuthResponse{
		Token: output.Token,
		User:  output.User,
	})
}

func (h *Handler) Login(c fiber.Ctx) error {
	req := new(AuthRequest)
	if err := c.Bind().Body(req); err != nil {
		return errorx.ErrParamsInvalid
	}

	if req.Email == "" || req.Password == "" {
		return errorx.New(fiber.StatusBadRequest, "Email and password are required")
	}

	output, err := h.S.Auth.Login(context.Background(), req.Email, req.Password)
	if err != nil {
		// Map implementation errors to HTTP errors if needed, or let middleware handle specific error types
		// For login, generic error usually means unauthorized
		return errorx.New(fiber.StatusUnauthorized, err.Error())
	}

	return c.JSON(AuthResponse{
		Token: output.Token,
		User:  output.User,
	})
}
