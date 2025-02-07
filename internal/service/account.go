package service

import (
	"context"
	"errors"
	"github.com/Parsa-Sedigh/kingscomp/internal/entity"
	"github.com/Parsa-Sedigh/kingscomp/internal/repository"
	"time"
)

const (
	DefaultState = "home"
)

type AccountService struct {
	accounts repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) *AccountService {
	return &AccountService{
		accounts: accountRepo,
	}
}

// CreateOrUpdate creates a new user in the data store or updates the existing user
func (a *AccountService) CreateOrUpdate(ctx context.Context, account entity.Account) (entity.Account, bool, error) {
	savedAccount, err := a.accounts.Get(ctx, account.EntityID())

	// user exists
	if err == nil {
		// fields have changed, we need to update the savedAccount and then update it
		if savedAccount.UserName != account.UserName || savedAccount.FirstName != account.FirstName {
			savedAccount.UserName = account.UserName
			savedAccount.FirstName = account.FirstName

			return savedAccount, false, a.accounts.Save(ctx, savedAccount)
		}

		return savedAccount, false, nil
	}

	// user does not exist in the database
	if errors.Is(err, repository.ErrNotFound) {
		account.JoinedAt = time.Now()
		account.State = DefaultState

		return account, true, a.accounts.Save(ctx, account)
	}

	return entity.Account{}, false, err
}

func (a *AccountService) Update(ctx context.Context, account entity.Account) error {
	return a.accounts.Save(ctx, account)
}
