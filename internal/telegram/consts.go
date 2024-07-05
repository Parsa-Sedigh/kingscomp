package telegram

import (
	"github.com/Parsa-Sedigh/kingscomp/internal/entity"
	"gopkg.in/telebot.v3"
	"time"
)

var (
	DefaultInputTimeout = time.Minute * 5
	DefaultTimeoutText  = "timeout"

	TxtConfirm = "بله"
	TxtDecline = "خیر"

	selector           = &telebot.ReplyMarkup{}
	btnEditDisplayName = selector.Data("ویرایش نام نمایشی", "btnEditDisplayName")
	btnNext            = selector.Data("-", "next")
)

func GetAccount(c telebot.Context) entity.Account {
	return c.Get("account").(entity.Account)
}
