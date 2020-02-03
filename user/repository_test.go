package user

import (
	"context"
	"github.com/imtanmoy/authn/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type userRepoMock struct {
	mock.Mock
}

func (o *userRepoMock) Exists(ctx context.Context, id int) bool {
	args := o.Called(ctx, id)
	return args.Bool(0)
}

func (o *userRepoMock) Find(ctx context.Context, id int) (*models.User, error) {
	args := o.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func Test_UserExist(t *testing.T) {
	repo := new(userRepoMock)
	c := context.Background()
	repo.On("Exists", mock.Anything, mock.AnythingOfType("int")).Return(true).Once()
	repo.Exists(c, 1)
	repo.AssertExpectations(t)
}

func Test_UserFind(t *testing.T) {
	repo := new(userRepoMock)
	c := context.Background()
	repo.On("Find", mock.Anything, mock.AnythingOfType("int")).Return(&models.User{
		ID:             1,
		Name:           "",
		Designation:    "",
		Email:          "",
		Password:       "",
		Enabled:        false,
		OrganizationId: 0,
		CreatedBy:      0,
		UpdatedBy:      0,
		DeletedBy:      0,
		JoinedAt:       time.Time{},
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
		DeletedAt:      time.Time{},
	}, nil).Once()
	u, err := repo.Find(c, 1)
	assert.Nil(t, err)
	assert.Equal(t, 1, u.ID)
	repo.AssertExpectations(t)
}
