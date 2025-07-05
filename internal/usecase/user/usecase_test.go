package user_test

import (
	"context"
	"errors"
	"testing"

	"merchshop/internal/entity"
	"merchshop/internal/usecase/user"

	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	CreateUserFunc    func(ctx context.Context, username string, password string) (*entity.User, error)
	GetByUsernameFunc func(ctx context.Context, username string) (*entity.User, error)
	GetByIDFunc       func(ctx context.Context, id int) (*entity.User, error)
}

func (m *mockUserRepo) CreateUser(ctx context.Context, username string, password string) (*entity.User, error) {
	return m.CreateUserFunc(ctx, username, password)
}

func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return m.GetByUsernameFunc(ctx, username)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int) (*entity.User, error) {
	return m.GetByIDFunc(ctx, id)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := &mockUserRepo{
		CreateUserFunc: func(ctx context.Context, username, password string) (*entity.User, error) {
			return &entity.User{ID: 1, Username: username, Password: password}, nil
		},
	}

	uc := user.NewUseCase(mockRepo)
	user, err := uc.Register(context.Background(), "testuser", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "password123", user.Password)
}

func TestRegister_Failure(t *testing.T) {
	mockRepo := &mockUserRepo{
		CreateUserFunc: func(ctx context.Context, username, password string) (*entity.User, error) {
			return nil, errors.New("username already taken")
		},
	}

	uc := user.NewUseCase(mockRepo)
	user, err := uc.Register(context.Background(), "testuser", "password123")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "failed to create user")
}

func TestGetByUsername_Success(t *testing.T) {
	mockRepo := &mockUserRepo{
		GetByUsernameFunc: func(ctx context.Context, username string) (*entity.User, error) {
			return &entity.User{ID: 1, Username: username}, nil
		},
	}

	uc := user.NewUseCase(mockRepo)
	user, err := uc.GetByUsername(context.Background(), "testuser")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
}

func TestGetByUsername_Failure(t *testing.T) {
	mockRepo := &mockUserRepo{
		GetByUsernameFunc: func(ctx context.Context, username string) (*entity.User, error) {
			return nil, errors.New("user not found")
		},
	}

	uc := user.NewUseCase(mockRepo)
	user, err := uc.GetByUsername(context.Background(), "nonexistentuser")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "failed to get user by username")
}

func TestGetByID_Success(t *testing.T) {
	mockRepo := &mockUserRepo{
		GetByIDFunc: func(ctx context.Context, id int) (*entity.User, error) {
			return &entity.User{ID: id, Username: "testuser"}, nil
		},
	}

	uc := user.NewUseCase(mockRepo)
	user, err := uc.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "testuser", user.Username)
}

func TestGetByID_Failure(t *testing.T) {
	mockRepo := &mockUserRepo{
		GetByIDFunc: func(ctx context.Context, id int) (*entity.User, error) {
			return nil, errors.New("user not found")
		},
	}

	uc := user.NewUseCase(mockRepo)
	user, err := uc.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "failed to get user by id")
}
