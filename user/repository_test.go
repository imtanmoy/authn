package user

import (
	"context"
	"github.com/imtanmoy/authn/internal/authx"
	"github.com/imtanmoy/authn/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type userRepoMock struct {
	mock.Mock
}

func (o *userRepoMock) FindAll(ctx context.Context) ([]*models.User, error) {
	panic("implement me")
}

func (o *userRepoMock) Save(ctx context.Context, u *models.User) error {
	panic("implement me")
}

func (o *userRepoMock) ExistsByEmail(ctx context.Context, email string) bool {
	panic("implement me")
}

func (o *userRepoMock) Delete(ctx context.Context, u *models.User) error {
	panic("implement me")
}

func (o *userRepoMock) Update(ctx context.Context, u *models.User) error {
	panic("implement me")
}

func (o *userRepoMock) FindByID(ctx context.Context, id int) (*models.User, error) {
	args := o.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (o *userRepoMock) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	panic("implement me")
}

func (o *userRepoMock) GetByEmail(ctx context.Context, identity string) (authx.AuthUser, error) {
	panic("implement me")
}

func (o *userRepoMock) ExistsByID(ctx context.Context, id int) bool {
	args := o.Called(ctx, id)
	return args.Bool(0)
}

var _ Repository = (*userRepoMock)(nil)

func (o *userRepoMock) Find(ctx context.Context, id int) (*models.User, error) {
	args := o.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func Test_UserExist(t *testing.T) {
	repo := new(userRepoMock)
	c := context.Background()
	repo.On("ExistsByID", mock.Anything, mock.AnythingOfType("int")).Return(true).Once()
	repo.ExistsByID(c, 1)
	repo.AssertExpectations(t)
}

func Test_FindByID(t *testing.T) {
	repo := new(userRepoMock)
	c := context.Background()
	repo.On("Find", mock.Anything, mock.AnythingOfType("int")).Return(&models.User{
		ID:        1,
		Name:      "test",
		Email:     "test@test.com",
		Password:  "password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Now(),
	}, nil).Once()
	u, err := repo.Find(c, 1)
	assert.Nil(t, err)
	assert.Equal(t, 1, u.ID)
	repo.AssertExpectations(t)
}
