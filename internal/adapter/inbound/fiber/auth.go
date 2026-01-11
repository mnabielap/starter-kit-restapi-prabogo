package fiber_inbound_adapter

import (
	"github.com/gofiber/fiber/v2"
	"prabogo/internal/domain"
	"prabogo/internal/model"
	inbound_port "prabogo/internal/port/inbound"
)

type authAdapter struct {
	domain domain.Domain
}

func NewAuthAdapter(domain domain.Domain) inbound_port.AuthHttpPort {
	return &authAdapter{domain: domain}
}

func (h *authAdapter) Register(a any) error {
	c := a.(*fiber.Ctx)
	var req model.UserInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Success: false, Error: "Invalid body"})
	}

	user, tokens, err := h.domain.Auth().Register(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Success: false, Error: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Success: true,
		Data: fiber.Map{
			"user":   user,
			"tokens": tokens,
		},
	})
}

func (h *authAdapter) Login(a any) error {
	c := a.(*fiber.Ctx)
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Success: false, Error: "Invalid body"})
	}

	user, tokens, err := h.domain.Auth().Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{Success: false, Error: err.Error()})
	}

	return c.JSON(model.Response{
		Success: true,
		Data: fiber.Map{
			"user":   user,
			"tokens": tokens,
		},
	})
}

func (h *authAdapter) RefreshToken(a any) error {
	c := a.(*fiber.Ctx)
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Success: false, Error: "Invalid body"})
	}

	tokens, err := h.domain.Auth().RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{Success: false, Error: err.Error()})
	}

	return c.JSON(model.Response{Success: true, Data: tokens})
}

func (h *authAdapter) Logout(a any) error {
	c := a.(*fiber.Ctx)
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Success: false, Error: "Invalid body"})
	}

	if err := h.domain.Auth().Logout(c.Context(), req.RefreshToken); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{Success: false, Error: err.Error()})
	}

	return c.JSON(model.Response{Success: true})
}

func (h *authAdapter) ForgotPassword(a any) error {
	c := a.(*fiber.Ctx)
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Success: false, Error: "Invalid body"})
	}

	h.domain.Auth().ForgotPassword(c.Context(), req.Email)
	return c.JSON(model.Response{Success: true, Message: "If email exists, reset link sent"})
}

func (h *authAdapter) ResetPassword(a any) error {
	c := a.(*fiber.Ctx)
	token := c.Query("token")
	var req struct {
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Success: false, Error: "Invalid body"})
	}

	if err := h.domain.Auth().ResetPassword(c.Context(), token, req.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Success: false, Error: err.Error()})
	}

	return c.JSON(model.Response{Success: true})
}