package service

import (
	"context"
	"errors"

	"app/internal/config"
	"app/internal/model/db"
	"app/internal/model/request"
	"app/internal/model/response"
	"app/internal/repository"

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
	userRepo         repository.UserRepository
}

func (c *conversationService) Create(ctx context.Context, userID int, req request.ConversationCreateRequest) (*response.ConversationResponse, error) {
	// invalid conversation type
	if req.Type != direct && req.Type != group {
		return nil, &Error{
			Code:    ErrCodeBadRequest,
			Message: "invalid conversation type",
		}
	}

	// missing user ids to add
	if len(req.UserIDs) == 0 {
		return nil, &Error{
			Code:    ErrCodeBadRequest,
			Message: "missing user ids",
		}
	}

	// missing group name
	if req.Name == nil && req.Type == group {
		return nil, &Error{
			Code:    ErrCodeBadRequest,
			Message: "missing group name",
		}
	}

	// should not add more than one member to a conversation
	if len(req.UserIDs) > 1 && req.Type == direct {
		return nil, &Error{
			Code:    ErrCodeBadRequest,
			Message: "too many user ids",
		}
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
				return nil, &Error{
					Code:    ErrCodeConflict,
					Message: "duplicate conversation",
				}
			}
		}
	}

	req.UserIDs = append(req.UserIDs, userID)

	conversation, err := c.conversationRepo.Create(ctx, userID, req.Type, req.Name, req.UserIDs)

	if err != nil {
		return nil, err
	}

	// if direct conversation set the name to that of the recipient
	if conversation.Type == direct {
		setRecipientName(ctx, userID, conversation, c.conversationRepo, c.userRepo)
	}

	resp := response.ConversationResponse{
		ID:              conversation.ID,
		Name:            *conversation.Name,
		AvatarURL:       conversation.AvatarURL,
		Type:            conversation.Type,
		LastMessageID:   conversation.LastMessageID,
		CreatedAt:       conversation.CreatedAt,
		CreatedBy:       conversation.CreatedBy,
		LastMessageRead: conversation.LastMessageRead,
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

	conversation, err := c.conversationRepo.GetByID(ctx, conversationID)

	if err != nil {
		return err
	}

	if conversation.Type == direct {
		return &Error{
			Code:    ErrCodeBadRequest,
			Message: "cannot add to direct conversation",
		}
	}

	return c.conversationRepo.AddMembers(ctx, conversationID, req.UserIDs)
}

func (c *conversationService) RemoveMember(ctx context.Context, userID, conversationID, memberID int) error {
	err := userInConversation(ctx, c.conversationRepo, conversationID, userID)

	if err != nil {
		return err
	}

	conversation, err := c.conversationRepo.GetByID(ctx, conversationID)

	if err != nil {
		return err
	}

	if conversation.Type == direct {
		return &Error{
			Code:    ErrCodeBadRequest,
			Message: "cannot remove member from direct conversation",
		}
	}

	return c.conversationRepo.RemoveMember(ctx, conversationID, memberID)
}

func (c *conversationService) GetByUserID(ctx context.Context, userID int) ([]response.ConversationResponse, error) {
	conversations, err := c.conversationRepo.GetByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	resp := make([]response.ConversationResponse, len(conversations))

	for i, conversation := range conversations {
		// if direct conversation set the name to that of the recipient
		if conversation.Type == direct {
			setRecipientName(ctx, userID, &conversation, c.conversationRepo, c.userRepo)
		}

		resp[i] = response.ConversationResponse{
			ID:              conversation.ID,
			Name:            *conversation.Name,
			AvatarURL:       conversation.AvatarURL,
			Type:            conversation.Type,
			LastMessageID:   conversation.LastMessageID,
			CreatedAt:       conversation.CreatedAt,
			CreatedBy:       conversation.CreatedBy,
			LastMessageRead: conversation.LastMessageRead,
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

	if err != nil {
		return nil, err
	}

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
			ID:        member.UserID,
			Username:  userMap[member.UserID].Username,
			AvatarURL: userMap[member.UserID].AvatarURL,
			JoinedAt:  member.JoinedAt,
		}
	}

	// if direct conversation set the name to that of the recipient
	if conversation.Type == direct {
		setRecipientName(ctx, userID, conversation, c.conversationRepo, c.userRepo)
	}

	resp := response.ConversationResponse{
		ID:              conversation.ID,
		Name:            *conversation.Name,
		AvatarURL:       conversation.AvatarURL,
		Type:            conversation.Type,
		LastMessageID:   conversation.LastMessageID,
		CreatedAt:       conversation.CreatedAt,
		CreatedBy:       conversation.CreatedBy,
		LastMessageRead: conversation.LastMessageRead,
		Members:         memberResp,
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

func NewConversationService(conversationRepo repository.ConversationRepository, userRepo repository.UserRepository, cfg *config.Config) ConversationService {
	return &conversationService{
		conversationRepo: conversationRepo,
		userRepo:         userRepo,
	}
}

func userInConversation(ctx context.Context, repo repository.ConversationRepository, conversationID, userID int) error {
	_, err := repo.GetMember(ctx, conversationID, userID)

	// sender doesn't belong in the conversation
	if errors.Is(err, pgx.ErrNoRows) {
		return &Error{
			Code:    ErrCodeForbidden,
			Message: "user not in conversation",
		}
	}

	return err
}

func setRecipientName(ctx context.Context, userID int, conversation *db.Conversation, conversationRepo repository.ConversationRepository, userRepo repository.UserRepository) error {
	members, err := conversationRepo.GetMembers(ctx, conversation.ID)

	if err != nil {
		return err
	}

	for _, member := range members {
		if member.UserID != userID {
			user, err := userRepo.GetByID(ctx, member.UserID)

			if err != nil {
				return err
			}

			name := user.Username
			conversation.Name = &name
		}
	}

	return nil
}