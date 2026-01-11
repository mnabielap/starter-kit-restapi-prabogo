package fiber_inbound_adapter

import (
	"context"
	"github.com/gofiber/fiber/v2"
	inbound_port "prabogo/internal/port/inbound"
)

func InitRoute(ctx context.Context, app *fiber.App, port inbound_port.HttpPort) {
	// --- AUTH ROUTES ---
	auth := app.Group("/v1/auth")
	auth.Post("/register", func(c *fiber.Ctx) error { return port.Auth().Register(c) })
	auth.Post("/login", func(c *fiber.Ctx) error { return port.Auth().Login(c) })
	auth.Post("/refresh-tokens", func(c *fiber.Ctx) error { return port.Auth().RefreshToken(c) })
	auth.Post("/logout", func(c *fiber.Ctx) error { return port.Auth().Logout(c) })
	auth.Post("/forgot-password", func(c *fiber.Ctx) error { return port.Auth().ForgotPassword(c) })
	auth.Post("/reset-password", func(c *fiber.Ctx) error { return port.Auth().ResetPassword(c) })

	// --- USER ROUTES ---
	users := app.Group("/v1/users")
	
	// Middleware for Users
	authMiddleware := func(c *fiber.Ctx) error { return port.Middleware().Auth(c) }
	adminMiddleware := func(c *fiber.Ctx) error { return port.Middleware().RequireAdmin(c) }
	adminOrSelfMiddleware := func(c *fiber.Ctx) error { return port.Middleware().RequireAdminOrSelf(c) }

	users.Use(authMiddleware)

	// Create: Admin Only
	users.Post("/", adminMiddleware, func(c *fiber.Ctx) error { return port.User().Create(c) })
	
	// Get List: Admin Only
	users.Get("/", adminMiddleware, func(c *fiber.Ctx) error { return port.User().GetList(c) })
	
	// Get One: Admin OR Self
	users.Get("/:id", adminOrSelfMiddleware, func(c *fiber.Ctx) error { return port.User().GetOne(c) })
	
	// Update: Admin Only
	users.Patch("/:id", adminMiddleware, func(c *fiber.Ctx) error { return port.User().Update(c) })
	
	// Delete: Admin Only
	users.Delete("/:id", adminMiddleware, func(c *fiber.Ctx) error { return port.User().Delete(c) })
}