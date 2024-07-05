package telegram

import (
	"context"
	"fmt"
	"gopkg.in/telebot.v3"
)

func (t *Telegram) editDisplayName(c telebot.Context) error {
	c.Delete()

	return t.editDisplayNamePrompt(c, "سلام اسمت چیه؟")
}

func (t *Telegram) editDisplayNamePrompt(c telebot.Context, promptText string) error {
	account := GetAccount(c)
	msg, err := t.Input(c, InputConfig{
		Prompt: promptText,
		Confirm: &Confirm{
			ConfirmText: func(msg *telebot.Message) string {
				return fmt.Sprintf("We call you %s from now on", msg.Text)
			},
		},
		Validator: Validator{
			Validator: func(msg *telebot.Message) bool {
				l := len([]rune(msg.Text))

				return l >= 3 && l <= 20
			},
			OnInvalid: func(msg *telebot.Message) string {
				return "name chars should be between 3 and 20"
			},
		},
	})
	if err != nil {
		return err
	}

	displayName := msg.Text
	// TODO: Validation
	account.DisplayName = displayName
	if err := t.App.Account.Update(context.Background(), account); err != nil {
		return err
	}

	c.Set("account", account)
	c.Reply(fmt.Sprint("اسم شما ", displayName, " است"))

	return nil
}
