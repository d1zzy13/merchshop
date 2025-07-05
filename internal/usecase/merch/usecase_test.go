package merch_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"merchshop/internal/entity"
	"merchshop/internal/usecase/merch"
)

// мок репозитория
type mockMerchRepo struct {
	ListFunc      func(ctx context.Context) ([]entity.Merchandise, error)
	GetByNameFunc func(ctx context.Context, name string) (*entity.Merchandise, error)
}

func (m *mockMerchRepo) List(ctx context.Context) ([]entity.Merchandise, error) {
	return m.ListFunc(ctx)
}

func (m *mockMerchRepo) GetByName(ctx context.Context, name string) (*entity.Merchandise, error) {
	return m.GetByNameFunc(ctx, name)
}

func TestUseCase_List(t *testing.T) {
	expected := []entity.Merchandise{
		{Name: "hoody", Price: 300},
		{Name: "cup", Price: 20},
	}

	mockRepo := &mockMerchRepo{
		ListFunc: func(ctx context.Context) ([]entity.Merchandise, error) {
			return expected, nil
		},
	}

	uc := merch.NewUseCase(mockRepo)

	result, err := uc.List(context.Background())
	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestUseCase_GetByName_Success(t *testing.T) {
	expected := &entity.Merchandise{Name: "hoody", Price: 300}

	mockRepo := &mockMerchRepo{
		GetByNameFunc: func(ctx context.Context, name string) (*entity.Merchandise, error) {
			if name == "hoody" {
				return expected, nil
			}
			return nil, errors.New("not found")
		},
	}

	uc := merch.NewUseCase(mockRepo)

	result, err := uc.GetByName(context.Background(), "hoody")
	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestUseCase_GetByName_EmptyName(t *testing.T) {
	mockRepo := &mockMerchRepo{}

	uc := merch.NewUseCase(mockRepo)

	result, err := uc.GetByName(context.Background(), "")
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "empty merchandise name")
}

func TestUseCase_GetByName_ErrorFromRepo(t *testing.T) {
	mockRepo := &mockMerchRepo{
		GetByNameFunc: func(ctx context.Context, name string) (*entity.Merchandise, error) {
			return nil, errors.New("db error")
		},
	}

	uc := merch.NewUseCase(mockRepo)

	result, err := uc.GetByName(context.Background(), "unknown")
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "failed to get merchandise by name")
}
