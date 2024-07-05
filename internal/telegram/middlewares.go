package telegram

import (
	"context"
	"github.com/Parsa-Sedigh/kingscomp/internal/entity"
	"gopkg.in/telebot.v3"
	"time"
)

func (t *Telegram) registerMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		acc := entity.Account{
			ID:        c.Sender().ID,
			FirstName: c.Sender().Username,
			LastName:  c.Sender().LastName,
			UserName:  c.Sender().Username,
			JoinedAt:  time.Now(),
		}

		account, created, err := t.App.Account.CreateOrUpdate(context.Background(), acc)
		if err != nil {
			return err
		}

		c.Set("account", account)
		c.Set("is_just_created", created)

		return next(c)
	}
}
