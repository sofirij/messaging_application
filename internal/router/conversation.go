package router

import (
	"app/internal/handler"
	"app/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterConversationRoutes(r fiber.Router, c *handler.ConversationHandler, m *handler.MessageHandler, jwtSecret string) {
	conversations := r.Group("/conversations", middleware.JWT(jwtSecret))

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
