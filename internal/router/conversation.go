package router

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/middleware"
	"app/internal/service"

	"github.com/gofiber/fiber/v3"
)

func RegisterConversationRoutes(r fiber.Router, c *handler.ConversationHandler, m *handler.MessageHandler, cfg *config.Config, userService service.UserService) {
	conversations := r.Group("/conversations", middleware.JWT(userService, cfg))

	conversations.Get("/", c.GetByUserID)
	conversations.Post("/", c.Create)
	conversations.Get("/:id", c.GetByID)
	conversations.Delete("/:id", c.SoftDelete)
	conversations.Get("/:id/messages", m.GetByConversationID)
	conversations.Post("/:id/messages", m.Create)
	conversations.Get("/:id/members", c.GetMembers)
	conversations.Post("/:id/members", c.AddMembers)
	conversations.Put("/:id/name", c.UpdateName)
	conversations.Put("/:id/avatar", c.UpdateAvatarURL)
	conversations.Delete("/:id/messages", c.ClearMessages)
	conversations.Delete("/:id/members/:uid", c.RemoveMember)
}
