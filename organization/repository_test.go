package organization

import (
	"context"
	"github.com/imtanmoy/authn/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type repoMock struct {
	mock.Mock
}

var repo *repoMock

func init() {
	repo = new(repoMock)
}

func (r *repoMock) FindByID(ctx context.Context, id int) (*models.Organization, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(*models.Organization), args.Error(1)
}

func (r *repoMock) Save(ctx context.Context, org *models.Organization) error {
	args := r.Called(ctx, org)
	return args.Error(0)
}

var _ Repository = (*repoMock)(nil)

func Test_Save(t *testing.T) {
	c := context.Background()
	org := &models.Organization{
		Name:    "Test Orgs",
		OwnerID: 1,
	}
	repo.On("Save", mock.Anything, org).Return(nil).Once()
	err := repo.Save(c, org)
	assert.Nil(t, err)
	repo.AssertExpectations(t)
}

func Test_FindByID(t *testing.T) {
	c := context.Background()
	org := &models.Organization{
		Name:    "Test Orgs",
		OwnerID: 1,
	}
	repo.On("FindByID", mock.Anything, mock.AnythingOfType("int")).Return(org, nil).Once()
	o, err := repo.FindByID(c, 1)
	assert.Nil(t, err)
	assert.NotNil(t, o)
	repo.AssertExpectations(t)

}
