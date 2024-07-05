package service

import (
	"context"
	"github.com/Parsa-Sedigh/kingscomp/internal/entity"
	"github.com/Parsa-Sedigh/kingscomp/internal/repository"
	"github.com/Parsa-Sedigh/kingscomp/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAccountService_CreateOrUpdateWithUserNotExists(t *testing.T) {
	accRepo := &mocks.AccountRepository{}
	accSvc := NewAccountService(accRepo)

	// prepare the mock
	accRepo.
		On("Get", mock.Anything, entity.NewID("account", 12)).
		Return(entity.Account{}, repository.ErrNotFound).
		Once()

	accRepo.
		On("Save", mock.Anything, mock.MatchedBy(func(acc entity.Account) bool {
			return acc.FirstName == "Parsa"
		})).
		Return(nil).
		Once()

	/* By looking at how we prepared the mocks in prev lines, after calling CreateOrUpdate(), we will get "Reza" as the saved account and
	not "Parsa" for ID: 12. Therefore, the CreateOrUpdate() method will try to update the curr saved record. So it will call the
	Save() method. So the `created` var should be false, because we did an update op not create.*/
	newAcc, created, err := accSvc.CreateOrUpdate(context.Background(), entity.Account{
		ID:        12,
		FirstName: "Parsa",
	})
	assert.NoError(t, err)
	assert.Equal(t, true, created)
	assert.Equal(t, "Parsa", newAcc.FirstName)

	accRepo.AssertExpectations(t)
}

func TestAccountService_CreateOrUpdateWithUserHasNotChanged(t *testing.T) {
	accRepo := &mocks.AccountRepository{}
	accSvc := NewAccountService(accRepo)

	accRepo.
		On("Get", mock.Anything, entity.NewID("account", 12)).
		Return(entity.Account{ID: 12, FirstName: "Parsa"}, nil).
		Once()

	newAcc, created, err := accSvc.CreateOrUpdate(context.Background(), entity.Account{
		ID:        12,
		FirstName: "Parsa",
	})
	assert.NoError(t, err)
	assert.Equal(t, false, created)
	assert.Equal(t, "Parsa", newAcc.FirstName)

	accRepo.AssertExpectations(t)
}
