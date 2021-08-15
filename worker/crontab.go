package worker

import (
	"os"
	"strconv"

	"github.com/MihaiBlebea/go-wallet/domain"
	"github.com/mileusna/crontab"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Service interface {
}

type service struct {
	ctab   *crontab.Crontab
	wallet domain.Wallet
	bot    *tb.Bot
}

func New(wallet domain.Wallet, bot *tb.Bot) Service {
	return &service{crontab.New(), wallet, bot}
}

func (s *service) Run() error {
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")
	chatID, err := strconv.Atoi(telegramChatID)
	if err != nil {
		return err
	}

	err = s.ctab.AddJob("* * * * *", func() {
		s.bot.Send(&tb.User{ID: chatID}, "hey there")
	})
	if err != nil {
		return err
	}

	return nil
}
