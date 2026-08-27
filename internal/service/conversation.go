package service

import (
	"context"
	"errors"

	"app/internal/model/db"
	"app/internal/model/request"
	"app/internal/model/response"
	"app/internal/repository"

	"github.com/jackc/pgx/v5"
)

type ConversationService interface {
	Create(ctx context.Context, userID int, req request.ConversationCreateRequest) (*response.ConversationResponse, error)
	ClearMessages(ctx context.Context, userID, conversationID int) error
	SoftDelete(ctx context.Context, userID, conversationID int) error
	GetMembers(ctx context.Context, userID, conversationID int) ([]response.UserResponse, error)
	AddMembers(ctx context.Context, userID, conversationID int, req request.ConversationAddMemberRequest) ([]response.UserResponse, error)
	RemoveMember(ctx context.Context, userID, conversationID, memberID int) error
	GetByUserID(ctx context.Context, userID int) ([]response.ConversationResponse, error)
	GetByID(ctx context.Context, userID, conversationID int) (*response.ConversationResponse, error)
	UpdateName(ctx context.Context, userID, conversationID int, req request.ConversationRenameRequest) error
	UpdateAvatarURL(ctx context.Context, userID, conversationID int, req request.ConversationAvatarRequest) error
}

type conversationService struct {
	conversationRepo repository.ConversationRepository
	userRepo         repository.UserRepository
	messageRepo      repository.MessageRepository
	hub              HubService
}

func NewConversationService(conversationRepo repository.ConversationRepository, userRepo repository.UserRepository, messageRepo repository.MessageRepository, hub HubService) ConversationService {
	return &conversationService{
		conversationRepo: conversationRepo,
		userRepo:         userRepo,
		hub:              hub,
		messageRepo:      messageRepo,
	}
}

func (c *conversationService) Create(ctx context.Context, userID int, req request.ConversationCreateRequest) (*response.ConversationResponse, error) {
	for i, val := range req.UserIDs {
		if val == userID {
			req.UserIDs = append(req.UserIDs[:i], req.UserIDs[i+1:]...)
		}
	}

	// invalid conversation type
	if req.Type != db.DirectConversation && req.Type != db.GroupConversation {
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

	// invalid group name
	if req.Type == db.GroupConversation {
		if req.Name == nil || len(*req.Name) < 3 {
			return nil, &Error{
				Code:    ErrCodeBadRequest,
				Message: "invalid group name",
			}
		}
	}

	// should not add more than one member to a conversation
	if len(req.UserIDs) > 1 && req.Type == db.DirectConversation {
		return nil, &Error{
			Code:    ErrCodeBadRequest,
			Message: "too many user ids",
		}
	}

	// if direct conversation, ensure to not create duplicate conversation
	if req.Type == db.DirectConversation {
		memberID := req.UserIDs[0]
		userConversations, err := c.conversationRepo.GetAllByUserID(ctx, userID)

		if err != nil {
			return nil, err
		}

		conversationMap := make(map[int]bool)

		for _, conversation := range userConversations {
			conversationMap[conversation.ID] = true
		}

		memberConversations, err := c.conversationRepo.GetAllByUserID(ctx, memberID)

		if err != nil {
			return nil, err
		}

		// resume and return the previously created conversation
		for _, conversation := range memberConversations {
			if conversationMap[conversation.ID] && conversation.Type == db.DirectConversation {
				member, err := c.conversationRepo.GetMember(ctx, conversation.ID, userID)
				if err != nil {
					return nil, err
				}

				// cannot resume a member that is not deleted so this is a conflict
				if member.DeletedAt == nil {
					return nil, &Error{
						Code: ErrCodeConflict,
						Message: "Duplicate conversation",
					}
				}

				_, err = c.conversationRepo.ResumeMember(ctx, conversation.ID, userID)

				if err != nil {
					return nil, err
				}

				setRecipientNameAndAvatar(ctx, userID, &conversation, c.conversationRepo, c.userRepo)

				lastMessageReadByUser, err := getLastMessageReadByUser(ctx, c.conversationRepo, userID, conversation.ID)

				if err != nil {
					return nil, err
				}

				lastMessageSentInConversation, err := getLastMessageSentInConversation(ctx, c.messageRepo, conversation.ID)

				if err != nil {
					return nil, err
				}

				lastMessageReadInConversation, err := getLastMessageReadInConversation(ctx, c.conversationRepo, userID, conversation.ID)

				if err != nil {
					return nil, err
				}

				resp := response.ConversationResponse{
					ID:                            conversation.ID,
					Name:                          *conversation.Name,
					AvatarURL:                     conversation.AvatarURL,
					Type:                          conversation.Type,
					CreatedAt:                     conversation.CreatedAt,
					CreatedBy:                     conversation.CreatedBy,
					LastMessageReadByUser:         lastMessageReadByUser,
					LastMessageReadInConversation: lastMessageReadInConversation,
					LastMessageSentInConversation: lastMessageSentInConversation,
				}

				return &resp, nil
			}
		}
	}

	req.UserIDs = append(req.UserIDs, userID)

	conversation, err := c.conversationRepo.Create(ctx, userID, req.Type, req.Name, req.UserIDs)

	if err != nil {
		return nil, err
	}

	// if direct conversation set the name to that of the recipient
	if conversation.Type == db.DirectConversation {
		setRecipientNameAndAvatar(ctx, userID, conversation, c.conversationRepo, c.userRepo)
	}

	lastMessageReadByUser, err := getLastMessageReadByUser(ctx, c.conversationRepo, userID, conversation.ID)

	if err != nil {
		return nil, err
	}

	lastMessageSentInConversation, err := getLastMessageSentInConversation(ctx, c.messageRepo, conversation.ID)

	if err != nil {
		return nil, err
	}

	lastMessageReadInConversation, err := getLastMessageReadInConversation(ctx, c.conversationRepo, userID, conversation.ID)

	if err != nil {
		return nil, err
	}

	resp := response.ConversationResponse{
		ID:                            conversation.ID,
		Name:                          *conversation.Name,
		AvatarURL:                     conversation.AvatarURL,
		Type:                          conversation.Type,
		CreatedAt:                     conversation.CreatedAt,
		CreatedBy:                     conversation.CreatedBy,
		LastMessageReadByUser:         lastMessageReadByUser,
		LastMessageReadInConversation: lastMessageReadInConversation,
		LastMessageSentInConversation: lastMessageSentInConversation,
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
	err := userInConversation(ctx, c.conversationRepo, conversationID, userID)

	if err != nil {
		return err
	}

	err = c.conversationRepo.SoftDeleteMember(ctx, conversationID, userID)

	if err != nil {
		return err
	}

	return nil
}

func (c *conversationService) AddMembers(ctx context.Context, userID, conversationID int, req request.ConversationAddMemberRequest) ([]response.UserResponse, error) {
	err := userInConversation(ctx, c.conversationRepo, conversationID, userID)

	if err != nil {
		return nil, err
	}

	conversation, err := c.conversationRepo.GetByID(ctx, conversationID)

	if err != nil {
		return nil, err
	}

	if conversation.Type == db.DirectConversation {
		return nil, &Error{
			Code:    ErrCodeBadRequest,
			Message: "cannot add to direct conversation",
		}
	}

	// ignore already existing members
	members, err := c.conversationRepo.GetMembers(ctx, conversationID)

	if err != nil {
		return nil, err
	}

	userSet := make(map[int]bool)
	for _, id := range req.UserIDs {
		userSet[id] = true
	}

	for _, member := range members {
		if userSet[member.UserID] {
			delete(userSet, member.UserID)
		}
	}

	k := 0
	toAdd := make([]int, len(userSet))
	for id := range userSet {
		toAdd[k] = id
		k += 1
	}

	if len(toAdd) == 0 {
		return nil, &Error{
			Code:    ErrCodeBadRequest,
			Message: "no new members to add to this conversation",
		}
	}

	err = c.conversationRepo.AddMembers(ctx, conversationID, toAdd)

	if err != nil {
		return nil, err
	}

	users, err := c.userRepo.GetByIDs(ctx, toAdd)

	if err != nil {
		return nil, err
	}

	resp := make([]response.UserResponse, len(users))

	for i, user := range users {
		resp[i] = response.UserResponse{
			ID:         user.ID,
			Username:   user.Username,
			AvatarURL:  user.AvatarURL,
			IsOnline:   c.hub.IsOnline(user.ID),
			LastSeenAt: user.LastSeenAt,
			CreatedAt:  user.CreatedAt,
		}
	}

	return resp, nil
}

func (c *conversationService) RemoveMember(ctx context.Context, userID, conversationID, memberID int) error {
	err := userInConversation(ctx, c.conversationRepo, conversationID, userID)

	if err != nil {
		return err
	}

	err = userInConversation(ctx, c.conversationRepo, conversationID, memberID)

	if err != nil {
		return err
	}

	conversation, err := c.conversationRepo.GetByID(ctx, conversationID)

	if err != nil {
		return err
	}

	if conversation.Type == db.DirectConversation {
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
		if conversation.Type == db.DirectConversation {
			setRecipientNameAndAvatar(ctx, userID, &conversation, c.conversationRepo, c.userRepo)
		}

		lastMessageReadByUser, err := getLastMessageReadByUser(ctx, c.conversationRepo, userID, conversation.ID)

		if err != nil {
			return nil, err
		}

		lastMessageSentInConversation, err := getLastMessageSentInConversation(ctx, c.messageRepo, conversation.ID)

		if err != nil {
			return nil, err
		}

		lastMessageReadInConversation, err := getLastMessageReadInConversation(ctx, c.conversationRepo, userID, conversation.ID)

		if err != nil {
			return nil, err
		}

		resp[i] = response.ConversationResponse{
			ID:                            conversation.ID,
			Name:                          *conversation.Name,
			AvatarURL:                     conversation.AvatarURL,
			Type:                          conversation.Type,
			CreatedAt:                     conversation.CreatedAt,
			CreatedBy:                     conversation.CreatedBy,
			LastMessageReadByUser:         lastMessageReadByUser,
			LastMessageReadInConversation: lastMessageReadInConversation,
			LastMessageSentInConversation: lastMessageSentInConversation,
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

	// if direct conversation set the name to that of the recipient
	if conversation.Type == db.DirectConversation {
		setRecipientNameAndAvatar(ctx, userID, conversation, c.conversationRepo, c.userRepo)
	}

	lastMessageReadByUser, err := getLastMessageReadByUser(ctx, c.conversationRepo, userID, conversation.ID)

	if err != nil {
		return nil, err
	}

	lastMessageSentInConversation, err := getLastMessageSentInConversation(ctx, c.messageRepo, conversation.ID)

	if err != nil {
		return nil, err
	}

	lastMessageReadInConversation, err := getLastMessageReadInConversation(ctx, c.conversationRepo, userID, conversation.ID)

	if err != nil {
		return nil, err
	}

	resp := response.ConversationResponse{
		ID:                            conversation.ID,
		Name:                          *conversation.Name,
		AvatarURL:                     conversation.AvatarURL,
		Type:                          conversation.Type,
		CreatedAt:                     conversation.CreatedAt,
		CreatedBy:                     conversation.CreatedBy,
		LastMessageReadByUser:         lastMessageReadByUser,
		LastMessageReadInConversation: lastMessageReadInConversation,
		LastMessageSentInConversation: lastMessageSentInConversation,
	}

	return &resp, nil
}

func (c *conversationService) UpdateName(ctx context.Context, userID, conversationID int, req request.ConversationRenameRequest) error {
	err := userInConversation(ctx, c.conversationRepo, conversationID, userID)

	if err != nil {
		return err
	}

	conversation, err := c.conversationRepo.GetByID(ctx, conversationID)

	if err != nil {
		return err
	}

	if conversation.Type == db.DirectConversation {
		return &Error{
			Code:    ErrCodeBadRequest,
			Message: "cannot update name of direct conversation",
		}
	}

	_, err = c.conversationRepo.UpdateName(ctx, conversationID, req.Name)

	return err
}

func (c *conversationService) UpdateAvatarURL(ctx context.Context, userID, conversationID int, req request.ConversationAvatarRequest) error {
	err := userInConversation(ctx, c.conversationRepo, conversationID, userID)

	if err != nil {
		return err
	}

	conversation, err := c.conversationRepo.GetByID(ctx, conversationID)

	if err != nil {
		return err
	}

	if conversation.Type == db.DirectConversation {
		return &Error{
			Code:    ErrCodeBadRequest,
			Message: "cannot update picture of direct conversation",
		}
	}

	_, err = c.conversationRepo.UpdateAvatarURL(ctx, conversationID, req.AvatarURL)

	return err
}

func (c *conversationService) GetMembers(ctx context.Context, userID int, conversationID int) ([]response.UserResponse, error) {
	err := userInConversation(ctx, c.conversationRepo, conversationID, userID)

	if err != nil {
		return nil, err
	}

	members, err := c.conversationRepo.GetMembers(ctx, conversationID)

	if err != nil {
		return nil, err
	}

	memberIDs := make([]int, len(members))
	for i, member := range members {
		memberIDs[i] = member.UserID
	}

	users, err := c.userRepo.GetByIDs(ctx, memberIDs)

	if err != nil {
		return nil, err
	}

	resp := make([]response.UserResponse, len(members))
	for i, user := range users {
		resp[i] = response.UserResponse{
			ID:         user.ID,
			Username:   user.Username,
			AvatarURL:  user.AvatarURL,
			IsOnline:   c.hub.IsOnline(user.ID),
			LastSeenAt: user.LastSeenAt,
			CreatedAt:  user.CreatedAt,
		}
	}

	return resp, nil
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

func setRecipientNameAndAvatar(ctx context.Context, userID int, conversation *db.Conversation, conversationRepo repository.ConversationRepository, userRepo repository.UserRepository) error {
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
			conversation.AvatarURL = user.AvatarURL
		}
	}

	return nil
}

func getLastMessageReadByUser(ctx context.Context, conversationRepo repository.ConversationRepository, userID, conversationID int) (*int, error) {
	member, err := conversationRepo.GetMember(ctx, conversationID, userID)

	if err != nil {
		return nil, err
	}

	return member.LastMessageRead, nil
}

func getLastMessageReadInConversation(ctx context.Context, conversationRepo repository.ConversationRepository, userID, conversationID int) (*int, error) {
	members, err := conversationRepo.GetMembers(ctx, conversationID)

	if err != nil {
		return nil, err
	}

	var lastMessageRead *int
	for _, member := range members {
		if lastMessageRead == nil {
			lastMessageRead = member.LastMessageRead
		} else if member.LastMessageRead != nil && *member.LastMessageRead > *lastMessageRead {
			lastMessageRead = member.LastMessageRead
		}
	}

	return lastMessageRead, nil
}

func getLastMessageSentInConversation(ctx context.Context, messageRepo repository.MessageRepository, conversationID int) (*response.MessageResponse, error) {
	message, err := messageRepo.GetLastInConversation(ctx, conversationID)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	resp, err := getMessageByID(ctx, messageRepo, message.ID)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
