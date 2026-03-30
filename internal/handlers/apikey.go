package handlers

import (
	"time"

	"fiber/internal/repository/dbgen"
	"fiber/internal/services"
	"fiber/middleware"
	"fiber/pkg/errorx"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CreateAPIKeyRequest struct {
	Name      string    `json:"name"`
	Scopes    []string  `json:"scopes"`
	ExpiresAt *string   `json:"expires_at"`
}

type APIKeyResponse struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Prefix     string   `json:"prefix"`
	Scopes     []string `json:"scopes"`
	Status     string   `json:"status"`
	LastUsedAt *string  `json:"last_used_at,omitempty"`
	ExpiresAt  *string  `json:"expires_at,omitempty"`
	CreatedAt  string   `json:"created_at"`
}

type CreatedAPIKeyResponse struct {
	APIKeyResponse
	Key string `json:"key"`
}

func (h *Handler) CreateAPIKey(c fiber.Ctx) error {
	userID, err := middleware.UserID(c)
	if err != nil {
		return err
	}

	req := new(CreateAPIKeyRequest)
	if err := c.Bind().Body(req); err != nil {
		return errorx.ErrParamsInvalid
	}

	var expiresAt *time.Time
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		parsed, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			return errorx.New(400, "expires_at must be RFC3339")
		}
		expiresAt = &parsed
	}

	result, err := h.S.APIKey.Create(c.Context(), userID, services.CreateAPIKeyInput{
		Name:      req.Name,
		Scopes:    req.Scopes,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(CreatedAPIKeyResponse{
		APIKeyResponse: toAPIKeyResponse(*result.Record),
		Key:            result.Key,
	})
}

func (h *Handler) ListAPIKeys(c fiber.Ctx) error {
	userID, err := middleware.UserID(c)
	if err != nil {
		return err
	}

	items, err := h.S.APIKey.List(c.Context(), userID)
	if err != nil {
		return err
	}

	response := make([]APIKeyResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toAPIKeyResponse(item))
	}
	return c.JSON(response)
}

func (h *Handler) RevokeAPIKey(c fiber.Ctx) error {
	userID, err := middleware.UserID(c)
	if err != nil {
		return err
	}

	apiKeyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errorx.ErrParamsInvalid
	}

	if err := h.S.APIKey.Revoke(c.Context(), userID, apiKeyID); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) DeleteAPIKey(c fiber.Ctx) error {
	return h.RevokeAPIKey(c)
}

func toAPIKeyResponse(item dbgen.ApiKey) APIKeyResponse {
	var lastUsedAt *string
	if item.LastUsedAt.Valid {
		value := item.LastUsedAt.Time.Format(time.RFC3339)
		lastUsedAt = &value
	}
	var expiresAt *string
	if item.ExpiresAt.Valid {
		value := item.ExpiresAt.Time.Format(time.RFC3339)
		expiresAt = &value
	}

	return APIKeyResponse{
		ID:         item.ID.String(),
		Name:       item.Name,
		Prefix:     item.KeyPrefix,
		Scopes:     item.Scopes,
		Status:     item.Status,
		LastUsedAt: lastUsedAt,
		ExpiresAt:  expiresAt,
		CreatedAt:  item.CreatedAt.Format(time.RFC3339),
	}
}
