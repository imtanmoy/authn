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

func (r *repoMock) Save(ctx context.Context, org *models.Organization) error {
	args := r.Called(ctx, org)
	return args.Error(0)
}

var _ Repository = (*repoMock)(nil)

func Test_Save(t *testing.T) {
	repo := new(repoMock)
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
