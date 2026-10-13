package story

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/ws"
	"context"
	"errors"
	"testing"
)

type fakeUsersRepository struct {
	getByIDUser *domain.User
	getByIDErr  error

	getByUsernameUser *domain.User
	getByUsernameErr  error

	getByIDCalls       int
	getByUsernameCalls int
}

func (f *fakeUsersRepository) GetByID(
	ctx context.Context,
	userID int64,
) (*domain.User, error) {
	f.getByIDCalls++
	return f.getByIDUser, f.getByIDErr
}

func (f *fakeUsersRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*domain.User, error) {
	f.getByUsernameCalls++
	return f.getByUsernameUser, f.getByUsernameErr
}

func (f *fakeUsersRepository) Create(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	return domain.User{}, nil
}

func (f *fakeUsersRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {
	return nil, nil
}

func (f *fakeUsersRepository) SearchByUsernamePrefix(
	ctx context.Context,
	prefix string,
	limit int,
) ([]domain.User, error) {
	return nil, nil
}

func (f *fakeUsersRepository) UpdateLastSeen(
	ctx context.Context,
	userID int64,
) error {
	return nil
}

type fakeDirectChatsRepository struct {
	createChat *domain.DirectChat
	createErr  error

	gotUser1ID    int64
	gotUser2ID    int64
	gotUser1Title string
	gotUser2Title string
	createCalls   int
}

func (f *fakeDirectChatsRepository) Create(
	ctx context.Context,
	user1ID int64,
	user2ID int64,
	user1Title string,
	user2Title string,
) (*domain.DirectChat, error) {

	f.gotUser1ID = user1ID
	f.gotUser2ID = user2ID
	f.gotUser1Title = user1Title
	f.gotUser2Title = user2Title

	f.createCalls++

	return f.createChat, f.createErr
}

func (f *fakeDirectChatsRepository) GetByUser(
	ctx context.Context,
	userID int64,
) ([]domain.DirectChat, error) {
	return nil, nil
}

func (f *fakeDirectChatsRepository) GetByUsers(
	ctx context.Context,
	user1ID int64,
	user2ID int64,
) (*domain.DirectChat, error) {
	return nil, nil
}

func (f *fakeDirectChatsRepository) GetByID(
	ctx context.Context,
	chatID int64,
) (*domain.DirectChat, error) {
	return nil, nil
}

func TestCreateDirectChatService_CreateSuccess(t *testing.T) {
	var currentUserID int64 = 10 // arrange
	currentUsername := "Alice"
	var otherUserID int64 = 20
	otherUsername := "Bob"
	userRepo := &fakeUsersRepository{getByIDUser: &domain.User{
		ID:       currentUserID,
		Username: currentUsername,
	},
		getByUsernameUser: &domain.User{
			ID:       otherUserID,
			Username: otherUsername,
		},
	}
	chatRepo := &fakeDirectChatsRepository{
		createChat: &domain.DirectChat{
			ID: 100,
		},
	}
	hub := ws.NewHub()

	service := NewCreateDirectChatService(chatRepo, userRepo, hub)

	chat, err := service.Create( // act
		context.Background(),
		CreateDirectChatInput{
			CurrentUserID: currentUserID,
			OtherUsername: otherUsername,
		},
	)
	if err != nil { // assert
		t.Fatalf("unexpected error: %v", err)
	}

	if chat == nil {
		t.Fatal("expected chat, got nil")
	}

	if chat.ID != 100 {
		t.Errorf("expected chat ID 100, got %d", chat.ID)
	}

	if chatRepo.gotUser1ID != currentUserID {
		t.Errorf("expected user1ID %d, got %d", currentUserID, chatRepo.gotUser1ID)
	}

	if chatRepo.gotUser2ID != otherUserID {
		t.Errorf("expected user1ID %d, got %d", otherUserID, chatRepo.gotUser2ID)
	}

	if chatRepo.gotUser1Title != "chat with Bob" {
		t.Errorf("expected user1Title %q, got %q",
			"chat with Bob",
			chatRepo.gotUser1Title)
	}

	if chatRepo.gotUser2Title != "chat with Alice" {
		t.Errorf("expected user2Title %q, got %q", "chat with Alice", chatRepo.gotUser2Title)
	}
}

func TestCreateDirectChatService_NormalizesUserOrder(t *testing.T) {
	var currentUserID int64 = 20 // arrange
	currentUsername := "Alice"
	var otherUserID int64 = 10
	otherUsername := "Bob"
	userRepo := &fakeUsersRepository{getByIDUser: &domain.User{
		ID:       currentUserID,
		Username: currentUsername,
	},
		getByUsernameUser: &domain.User{
			ID:       otherUserID,
			Username: otherUsername,
		},
	}
	chatRepo := &fakeDirectChatsRepository{
		createChat: &domain.DirectChat{
			ID: 100,
		},
	}
	hub := ws.NewHub()

	service := NewCreateDirectChatService(chatRepo, userRepo, hub)

	chat, err := service.Create( // act
		context.Background(),
		CreateDirectChatInput{
			CurrentUserID: currentUserID,
			OtherUsername: otherUsername,
		},
	)
	if err != nil { // assert
		t.Fatalf("unexpected error: %v", err)
	}

	if chat == nil {
		t.Fatal("expected chat, got nil")
	}

	if chat.ID != 100 {
		t.Errorf("expected chat ID 100, got %d", chat.ID)
	}

	if chatRepo.gotUser1ID != otherUserID {
		t.Errorf("expected user1ID %d, got %d", otherUserID, chatRepo.gotUser1ID)
	}

	if chatRepo.gotUser2ID != currentUserID {
		t.Errorf("expected user1ID %d, got %d", currentUserID, chatRepo.gotUser2ID)
	}

	if chatRepo.gotUser1Title != "chat with Alice" {
		t.Errorf("expected user1Title %q, got %q",
			"chat with Alice",
			chatRepo.gotUser1Title)
	}

	if chatRepo.gotUser2Title != "chat with Bob" {
		t.Errorf("expected user2Title %q, got %q", "chat with Bob", chatRepo.gotUser2Title)
	}
}

func TestCreateDirectChatService_CantCreateChatWithSelf(t *testing.T) {
	currentUserID := int64(10) // arrange
	otherUserID := int64(10)

	userRepo := &fakeUsersRepository{getByIDUser: &domain.User{
		ID: currentUserID,
	},
		getByUsernameUser: &domain.User{
			ID: otherUserID,
		},
	}
	chatRepo := &fakeDirectChatsRepository{}
	hub := ws.NewHub()

	service := NewCreateDirectChatService(chatRepo, userRepo, hub)

	chat, err := service.Create( // act
		context.Background(),
		CreateDirectChatInput{
			CurrentUserID: currentUserID,
			OtherUsername: "Alice",
		},
	)

	if !errors.Is(err, domain.ErrCantCreateChat) { // assert
		t.Fatalf("expeted error can't create chat, got %v", err)
	}
	if chat != nil {
		t.Fatalf("expeted not created chat, got %v", chat)
	}
	if chatRepo.createCalls != 0 {
		t.Fatalf("expected create not to be called, got %d calls", chatRepo.createCalls)
	}
}

func TestCreateDirectChatService_OtherUserNotFound(t *testing.T) {
	currentUserID := int64(10)

	userRepo := &fakeUsersRepository{
		getByIDUser: &domain.User{
			ID: currentUserID,
		},
		getByUsernameErr: domain.ErrUserNotFound,
	}

	chatRepo := &fakeDirectChatsRepository{}
	hub := ws.NewHub()

	service := NewCreateDirectChatService(
		chatRepo,
		userRepo,
		hub,
	)

	chat, err := service.Create(
		context.Background(),
		CreateDirectChatInput{
			CurrentUserID: currentUserID,
			OtherUsername: "Bob",
		},
	)

	if !errors.Is(err, domain.ErrUserNotFound) { // assert
		t.Fatalf("expeted error user not found, got %v", err)
	}
	if chat != nil {
		t.Fatalf("expected not created chat, got %v", chat)
	}
	if chatRepo.createCalls != 0 {
		t.Fatalf("expected create not to be called, got %d calls", chatRepo.createCalls)
	}
}

var errCreateChat = errors.New("create chat failed")

func TestCreateDirectChatService_CreateRepositoryError(t *testing.T) {

	currentUserID := int64(10)
	otherUserID := int64(20)

	chatRepo := &fakeDirectChatsRepository{
		createErr: errCreateChat,
	}

	userRepo := &fakeUsersRepository{getByIDUser: &domain.User{
		ID:       currentUserID,
		Username: "Alice",
	},
		getByUsernameUser: &domain.User{
			ID:       otherUserID,
			Username: "Bob",
		}}

	hub := ws.NewHub()
	service := NewCreateDirectChatService(chatRepo, userRepo, hub)

	chat, err := service.Create(context.Background(), CreateDirectChatInput{
		CurrentUserID: currentUserID,
		OtherUsername: "Bob"})

	if !errors.Is(err, errCreateChat) {
		t.Fatalf("expeted error create chat error, got %v", err)
	}

	if chat != nil {
		t.Fatalf("expected not created chat, got %v", chat)
	}

	if chatRepo.createCalls != 1 {
		t.Fatalf("expected Create to be called once, got %d calls", chatRepo.createCalls)
	}
}

var errGetCurrentUser = errors.New("get current user failed")

func TestCreateDirectChatService_CurrentUserRepositoryError(t *testing.T) {

}
