package service

import (
	"context"
	"strconv"
	"time"

	"app/internal/model/response"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const ticketTTL = 30 * time.Second

type TicketService interface {
	Create(ctx context.Context, userID int) (*response.TicketResponse, error)
	ValidateAndConsume(ticket string) (*int, error)
}

type ticketService struct {
	storage fiber.Storage
}

func NewTicketService(storage fiber.Storage) TicketService {
	return &ticketService{storage: storage}
}

func (t *ticketService) Create(ctx context.Context, userID int) (*response.TicketResponse, error) {
	ticket := uuid.New().String()
	err := t.storage.Set(ticket, []byte(strconv.Itoa(userID)), ticketTTL)

	if err != nil {
		return nil, err
	}

	return &response.TicketResponse{
		Ticket:    ticket,
		ExpiresAt: time.Now().Add(ticketTTL),
	}, nil
}

func (t *ticketService) ValidateAndConsume(ticket string) (*int, error) {
	val, err := t.storage.Get(ticket)

	if err != nil {
		return nil, err
	}

	// ticket has expired
	if val == nil {
		return nil, &Error{
			Code:    ErrCodeForbidden,
			Message: "ticket is expired",
		}
	}

	userID, err := strconv.Atoi(string(val))

	if err != nil {
		return nil, err
	}

	t.storage.Delete(ticket)

	return &userID, nil
}
