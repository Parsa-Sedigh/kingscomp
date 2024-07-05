package telegram

import (
	"errors"
	"fmt"
	"github.com/Parsa-Sedigh/kingscomp/internal/service"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
	"time"
)

var (
	ErrInputTimeout = errors.New("input timeout")
	ErrInvalidInput = errors.New(" ")
)

type Telegram struct {
	App        *service.App
	bot        *telebot.Bot
	Teleprompt *TelePrompt
}

type Confirm struct {
	ConfirmText func(msg *telebot.Message) string
}

type Validator struct {
	Validator func(msg *telebot.Message) bool
	OnInvalid func(msg *telebot.Message) string
}

type InputConfig struct {
	Prompt         any
	OnTimeout      any
	HasConfirm     bool
	PromptKeyboard [][]string
	Confirm        *Confirm
	Validator      Validator
}

func NewTelegram(app *service.App, apiKey string) (*Telegram, error) {
	t := &Telegram{
		App:        app,
		Teleprompt: NewTelePrompt(),
	}

	pref := telebot.Settings{
		Token:   apiKey,
		Poller:  &telebot.LongPoller{Timeout: 60 * time.Second},
		OnError: t.onError,
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		logrus.WithError(err).Error("couldn't connect to telegram servers")

		return nil, err
	}

	t.bot = bot

	t.setupHandlers()

	return t, nil
}

func (t *Telegram) Start() {
	t.bot.Start()
}

func (t *Telegram) Input(c telebot.Context, config InputConfig) (*telebot.Message, error) {
getInput:
	// makes a prompt to the client and ask for the data
	if config.Prompt != nil {
		if config.PromptKeyboard != nil {
			c.Send(config.Prompt, generateKeyboard(config.PromptKeyboard))
		} else {
			c.Send(config.Prompt)
		}
	}

	// waits for the client until the response is fetched
	response, isTimeout := t.Teleprompt.AsMessage(c.Sender().ID, DefaultInputTimeout)
	if isTimeout {
		if config.OnTimeout != nil {
			c.Reply(config.OnTimeout)
		} else {
			c.Reply(DefaultTimeoutText)
		}

		return nil, ErrInputTimeout
	}

	// validate the response
	if config.Validator.Validator != nil && !config.Validator.Validator(response) {
		c.Reply(config.Validator.OnInvalid(response))

		goto getInput
	}

	// client has to confirm
	if config.Confirm.ConfirmText != nil {
		confirmText := config.Confirm.ConfirmText(response)

		confirmMessage, err := t.Input(c, InputConfig{
			Prompt:         confirmText,
			PromptKeyboard: [][]string{{TxtDecline}, {TxtConfirm}},
			Validator:      choiceValidator(TxtDecline, TxtConfirm),
		})
		if err != nil {
			return nil, err
		}

		// on confirm we need to do nothing, but on decline, do this:
		if confirmMessage.Text == TxtDecline {
			goto getInput
		}
	}

	return response, nil
}

func (t *Telegram) onError(err error, c telebot.Context) {
	if errors.Is(err, ErrInputTimeout) {
		return
	}

	errID := uuid.New().String()

	logrus.WithError(err).WithField("tracing_id", errID).Errorln("unhandled error")

	c.Reply(fmt.Sprintf("در پردازش اطلاعات مشکلی پیش آمد.\n کد بررسی: %s", errID))
}

func generateKeyboard(rows [][]string) *telebot.ReplyMarkup {
	mu := &telebot.ReplyMarkup{
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	mu.Reply(lo.Map(rows, func(row []string, _ int) telebot.Row {
		return mu.Row(lo.Map(row, func(btn string, _ int) telebot.Btn {
			return mu.Text(btn)
		})...)
	})...)

	return mu
}
