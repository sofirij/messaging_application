package service

import (
	"context"
	"errors"

	"app/internal/model/request"
	"app/internal/model/response"
	"app/internal/model/service"
	"app/internal/repository"
	"app/internal/model/db"
	"app/internal/config"

	"github.com/jackc/pgx/v5"
)

const (
	direct = "direct"
	group  = "group"
)

type ConversationService interface {
	Create(ctx context.Context, userID int, req request.ConversationCreateRequest) (*response.ConversationResponse, error)
	ClearMessages(ctx context.Context, userID, conversationID int) error
	SoftDelete(ctx context.Context, userID, conversationID int) error
	AddMember(ctx context.Context, userID, conversationID int, req request.ConversationAddMemberRequest) error
	RemoveMember(ctx context.Context, userID, conversationID, memberID int) error
	GetByUserID(ctx context.Context, userID int) ([]response.ConversationResponse, error)
	GetByID(ctx context.Context, userID, conversationID int) (*response.ConversationResponse, error)
	UpdateName(ctx context.Context, userID, conversationID int, req request.ConversationRenameRequest) error
	UpdateAvatarURL(ctx context.Context, userID, conversationID int, req request.ConversationAvatarRequest) error
}

type conversationService struct {
	conversationRepo repository.ConversationRepository
	userRepo repository.UserRepository
}

func (c *conversationService) Create(ctx context.Context, userID int, req request.ConversationCreateRequest) (*response.ConversationResponse, error) {
	// invalid conversation type
	if req.Type != direct && req.Type != group {
		return nil, &service.Error{}
	}

	// should not add more than one member to a conversation
	if len(req.UserIDs) > 1 && req.Type == direct {
		return nil, &service.Error{}
	}

	// if direct conversation, ensure to not create duplicate conversation
	if req.Type == direct {
		memberID := req.UserIDs[0]
		userConversations, err := c.conversationRepo.GetByUserID(ctx, userID)

		if err != nil {
			return nil, err
		}

		conversationMap := make(map[int]bool)

		for _, conversation := range userConversations {
			conversationMap[conversation.ID] = true
		}

		memberConversations, err := c.conversationRepo.GetByUserID(ctx, memberID)

		if err != nil {
			return nil, err
		}

		for _, conversation := range memberConversations {
			if conversationMap[conversation.ID] {
				return nil, &service.Error{}
			}
		}
	}

	req.UserIDs = append(req.UserIDs, userID)

	conversation, err := c.conversationRepo.Create(ctx, userID, req.Type, req.Name, req.UserIDs)

	if err != nil {
		return nil, err
	}

	resp := response.ConversationResponse{
		ID:            conversation.ID,
		Name:          conversation.Name,
		AvatarURL:     conversation.AvatarURL,
		Type:          conversation.Type,
		LastMessageAt: conversation.LastMessageAt,
		CreatedAt:     conversation.CreatedAt,
		CreatedBy:     conversation.CreatedBy,
	}

	return &resp, nil
}

func (c *conversationService) ClearMessages(ctx context.Context, userID, conversationID int) error {
	err := userInConversation(ctx, c.conversationRepo, conversationID, userID)

	if err != nil {
		return err
	}

	err = c.conversationRepo.ClearMessages(ctx, conversationID, userID)

	if err != nil {
		return err
	}

	return nil
}

func (c *conversationService) SoftDelete(ctx context.Context, userID, conversationID int) error {
	err := c.ClearMessages(ctx, userID, conversationID)

	if err != nil {
		return err
	}

	err = c.conversationRepo.SoftDeleteMember(ctx, conversationID, userID)

	if err != nil {
		return err
	}

	return nil
}

func (c *conversationService) AddMember(ctx context.Context, userID, conversationID int, req request.ConversationAddMemberRequest) error {
	err := userInConversation(ctx, c.conversationRepo, conversationID, userID)

	if err != nil {
		return err
	}

	return c.conversationRepo.AddMembers(ctx, conversationID, req.UserIDs)
}

func (c *conversationService) RemoveMember(ctx context.Context, userID, conversationID, memberID int) error {
	err := userInConversation(ctx, c.conversationRepo, conversationID, userID)

	if err != nil {
		return err
	}

	return c.conversationRepo.RemoveMember(ctx, conversationID, userID)
}

func (c *conversationService) GetByUserID(ctx context.Context, userID int) ([]response.ConversationResponse, error) {
	conversations, err := c.conversationRepo.GetByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	resp := make([]response.ConversationResponse, len(conversations))

	for i, conversation := range conversations {
		resp[i] = response.ConversationResponse{
			ID: conversation.ID,
			Name: conversation.Name,
			AvatarURL: conversation.AvatarURL,
			Type: conversation.Type,
			LastMessageAt: conversation.LastMessageAt,
			CreatedAt: conversation.CreatedAt,
			CreatedBy: conversation.CreatedBy,
		}
	}

	return resp, nil
}

func (c *conversationService) GetByID(ctx context.Context, userID, conversationID int) (*response.ConversationResponse, error) {
	err := userInConversation(ctx, c.conversationRepo, conversationID, userID)

	if err != nil {
		return nil, err
	}

	conversation, err := c.conversationRepo.GetByID(ctx, conversationID)

	if err != nil {
		return nil, err
	}

	members, err := c.conversationRepo.GetMembers(ctx, conversationID)

	userIDs := make([]int, len(members))

	for i, member := range members {
		userIDs[i] = member.UserID
	}

	users, err := c.userRepo.GetByIDs(ctx, userIDs)

	if err != nil {
		return nil, err
	}
	
	userMap := make(map[int]db.User, len(members))

	for _, user := range users {
		userMap[user.ID] = user
	}

	memberResp := make([]response.MemberResponse, len(members))

	for i, member := range members {
		memberResp[i] = response.MemberResponse{
			ID: member.UserID,
			Username: userMap[member.UserID].Username,
			AvatarURL: userMap[member.UserID].AvatarURL,
			JoinedAt: member.JoinedAt,
		}
	}

	resp := response.ConversationResponse{
		ID: conversation.ID,
		Name: conversation.Name,
		AvatarURL: conversation.AvatarURL,
		Type: conversation.Type,
		LastMessageAt: conversation.LastMessageAt,
		CreatedAt: conversation.CreatedAt,
		CreatedBy: conversation.CreatedBy,
		Members: memberResp,
	}

	return &resp, nil
}

func (c *conversationService) UpdateName(ctx context.Context, userID, conversationID int, req request.ConversationRenameRequest) error {
	err := userInConversation(ctx, c.conversationRepo, conversationID, userID)

	if err != nil {
		return err
	}

	_, err = c.conversationRepo.UpdateName(ctx, conversationID, req.Name)

	return err
}

func (c *conversationService) UpdateAvatarURL(ctx context.Context, userID, conversationID int, req request.ConversationAvatarRequest) error {
	err := userInConversation(ctx, c.conversationRepo, conversationID, userID)

	if err != nil {
		return err
	}

	_, err = c.conversationRepo.UpdateAvatarURL(ctx, conversationID, req.AvatarURL)

	return err
}

func NewConversationRepository(conversationRepo repository.ConversationRepository, userRepo repository.UserRepository, cfg *config.Config) ConversationService {
	return &conversationService{
		conversationRepo: conversationRepo,
		userRepo: userRepo,
	}
}


func userInConversation(ctx context.Context, repo repository.ConversationRepository, conversationID, userID int) error {
	_, err := repo.GetMember(ctx, conversationID, userID)

	// sender doesn't belong in the conversation
	if errors.Is(err, pgx.ErrNoRows) {
		return &service.Error{}
	}

	return nil
}
