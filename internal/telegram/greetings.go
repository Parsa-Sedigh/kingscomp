package telegram

import (
	"fmt"
	"gopkg.in/telebot.v3"
)

func (t *Telegram) start(c telebot.Context) error {
	isJustCreated := c.Get("is_just_created").(bool)
	if !isJustCreated {
		return t.myInfo(c)
	}

	if err := t.editDisplayNamePrompt(c, "what do you want to be called?"); err != nil {
		return err
	}

	return t.myInfo(c)
}

func (t *Telegram) myInfo(c telebot.Context) error {
	account := GetAccount(c)
	selector := &telebot.ReplyMarkup{}
	selector.Inline(selector.Row(btnEditDisplayName))

	return c.Send(fmt.Sprintf("King %s, what we can do for you?", account.DisplayName), selector)
}
