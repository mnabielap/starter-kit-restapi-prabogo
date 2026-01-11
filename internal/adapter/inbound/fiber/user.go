package fiber_inbound_adapter

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"prabogo/internal/domain"
	"prabogo/internal/model"
	inbound_port "prabogo/internal/port/inbound"
)

type userAdapter struct {
	domain domain.Domain
}

func NewUserAdapter(domain domain.Domain) inbound_port.UserHttpPort {
	return &userAdapter{domain: domain}
}

func (h *userAdapter) Create(a any) error {
	c := a.(*fiber.Ctx)
	var req model.UserInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Success: false, Error: "Invalid body"})
	}

	user, err := h.domain.User().Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Success: false, Error: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(model.Response{Success: true, Data: user})
}

func (h *userAdapter) GetList(a any) error {
	c := a.(*fiber.Ctx)
	
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "created_at:desc")
	
	filters := model.UserFilter{
		Search: c.Query("search"),
		Role:   c.Query("role"),
	}

	users, total, err := h.domain.User().GetAll(c.Context(), filters, page, limit, sortBy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{Success: false, Error: err.Error()})
	}

	return c.JSON(model.Response{
		Success: true,
		Data: fiber.Map{
			"results":       users,
			"page":          page,
			"limit":         limit,
			"total_results": total,
		},
	})
}

func (h *userAdapter) GetOne(a any) error {
	c := a.(*fiber.Ctx)
	id := c.Params("id")

	user, err := h.domain.User().GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{Success: false, Error: err.Error()})
	}
	return c.JSON(model.Response{Success: true, Data: user})
}

func (h *userAdapter) Update(a any) error {
	c := a.(*fiber.Ctx)
	id := c.Params("id")
	var req model.UserInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Success: false, Error: "Invalid body"})
	}

	user, err := h.domain.User().Update(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Success: false, Error: err.Error()})
	}
	return c.JSON(model.Response{Success: true, Data: user})
}

func (h *userAdapter) Delete(a any) error {
	c := a.(*fiber.Ctx)
	id := c.Params("id")

	if err := h.domain.User().Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{Success: false, Error: err.Error()})
	}
	return c.JSON(model.Response{Success: true})
}