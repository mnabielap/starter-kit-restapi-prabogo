package fiber_inbound_adapter

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"prabogo/internal/domain"
	"prabogo/internal/model"
	"prabogo/utils/jwt"
)

type MiddlewareAdapter interface {
	Auth(a any) error
	RequireAdmin(a any) error
	RequireAdminOrSelf(a any) error
	InternalAuth(a any) error
	ClientAuth(a any) error
}

type middlewareAdapter struct {
	domain domain.Domain
}

func NewMiddlewareAdapter(domain domain.Domain) MiddlewareAdapter {
	return &middlewareAdapter{domain: domain}
}

// --- New Middleware ---

func (m *middlewareAdapter) Auth(a any) error {
	c := a.(*fiber.Ctx)
	authHeader := c.Get("Authorization")
	
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{Success: false, Error: "Missing Authorization header"})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{Success: false, Error: "Invalid Authorization header format"})
	}

	tokenString := parts[1]
	claims, err := jwt.ValidateLocalToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{Success: false, Error: "Invalid or expired token"})
	}

	if claims.Type != "access" {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{Success: false, Error: "Invalid token type"})
	}

	c.Locals("userID", claims.Sub)
	
	return c.Next()
}

func (m *middlewareAdapter) RequireAdmin(a any) error {
	c := a.(*fiber.Ctx)
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{Success: false, Error: "Unauthorized"})
	}

	user, err := m.domain.User().GetByID(c.Context(), userID)
	if err != nil || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{Success: false, Error: "User not found"})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(model.Response{Success: false, Error: "Forbidden: Admins only"})
	}

	return c.Next()
}

func (m *middlewareAdapter) RequireAdminOrSelf(a any) error {
	c := a.(*fiber.Ctx)
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{Success: false, Error: "Unauthorized"})
	}

	targetID := c.Params("id")

	if targetID == userID {
		return c.Next()
	}

	user, err := m.domain.User().GetByID(c.Context(), userID)
	if err != nil || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{Success: false, Error: "User not found"})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(model.Response{Success: false, Error: "Forbidden: Access denied"})
	}

	return c.Next()
}

// --- Legacy Middleware (Stubs to satisfy interface) ---

func (m *middlewareAdapter) InternalAuth(a any) error {
	// Implement original logic if needed, or redirect to Auth
	return m.Auth(a)
}

func (m *middlewareAdapter) ClientAuth(a any) error {
	// Implement original logic if needed, or redirect to Auth
	return m.Auth(a)
}