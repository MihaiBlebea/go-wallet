package bot

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/MihaiBlebea/go-wallet/domain"
	"github.com/mileusna/crontab"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Service interface {
	Start()
	SendMessageToUser(message string) error
	ToAmount(amount int, currency string) string
}

type service struct {
	bot        *tb.Bot
	userChatID int
	cron       *crontab.Crontab
	wallet     domain.Wallet
}

// https://golangrepo.com/repo/tucnak-telebot-go-bot-building#sendable

func New(wallet domain.Wallet) (Service, error) {
	bot, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 30 * time.Second},
	})
	if err != nil {
		return &service{}, err
	}

	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	intChatID, err := strconv.Atoi(chatID)
	if err != nil {
		return &service{}, err
	}

	serv := &service{
		bot:        bot,
		userChatID: intChatID,
		cron:       crontab.New(),
		wallet:     wallet,
	}

	return serv, nil
}

func (s *service) Start() {
	s.bot.Handle("/hello", func(m *tb.Message) {
		// res, err := s.ytube.Search("Harald Baldr")
		// if err != nil {
		// 	s.bot.Send(m.Sender, err.Error())
		// }

		s.bot.Send(m.Sender, "Hello back to you")
		// for _, v := range res.Items {
		// 	fmt.Printf("%+v", v)
		// 	title := v.Snippet.Title
		// 	thumbnail := &tb.Photo{File: tb.FromURL(v.Snippet.Thumbnails.High.URL)}
		// 	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", v.Snippet.ID.VideoID)

		// 	s.bot.Send(m.Sender, thumbnail)
		// 	s.bot.Send(m.Sender, fmt.Sprintf("<a href='%s'>%s</a>", url, title), &tb.SendOptions{
		// 		ParseMode:             tb.ModeHTML,
		// 		DisableWebPagePreview: true,
		// 	})
		// 	// s.bot.Send(m.Sender, url, &tb.SendOptions{DisableWebPagePreview: true})
		// 	time.Sleep(time.Second * 2)
		// }
	})

	s.bot.Handle(OnAdd, func(m *tb.Message) {
		p, err := s.parseReceiptCommand(m.Text)
		if err != nil {
			s.bot.Send(m.Sender, err.Error())
			return
		}

		shortID, remainingAmount, err := s.wallet.StoreReceipt(
			p.Description,
			p.Amount,
			p.Currency,
			p.Quantity,
			p.Unit,
		)
		if err != nil {
			s.bot.Send(m.Sender, err.Error())
			return
		}

		if remainingAmount > 0 {
			s.bot.Send(
				m.Sender,
				fmt.Sprintf(
					"%s - I added it, remaining amount %s",
					shortID,
					s.ToAmount(remainingAmount, "GBP"),
				),
			)
		} else {
			s.bot.Send(m.Sender, "No more remaining amount, marking transaction as completed.")
		}
	})

	s.bot.Handle(OnDelete, func(m *tb.Message) {
		shortID, err := s.parseDeleteCommand(m.Text)
		if err != nil {
			s.bot.Send(m.Sender, err.Error())
			return
		}

		err = s.wallet.DeleteReceipt(shortID)
		if err != nil {
			s.bot.Send(m.Sender, err.Error())
			return
		}

		s.bot.Send(m.Sender, fmt.Sprintf("%s - I deleted it", shortID))
	})

	s.bot.Handle(OnConfirm, func(m *tb.Message) {
		err := s.wallet.MarkTransactionCompleted()
		if err != nil {
			s.bot.Send(m.Sender, err.Error())
			return
		}

		s.bot.Send(m.Sender, "Confirmed. Thanks!")
	})

	s.bot.Handle(OnCancel, func(m *tb.Message) {
		err := s.wallet.MarkTransactionCompleted()
		if err != nil {
			s.bot.Send(m.Sender, err.Error())
			return
		}

		s.bot.Send(m.Sender, "Oki, no worries, I cancelled that.")
	})

	s.bot.Handle(tb.OnText, func(m *tb.Message) {

		s.bot.Send(m.Sender, fmt.Sprintf("%d", m.Sender.ID))
		// chatId := m.Chat.ID

		// s.bot.Send(m.Sender, resp)
	})

	// TODO: just for testing, remove
	s.bot.Handle(OnTest, func(m *tb.Message) {
		var (
			// Universal markup builders.
			menu     = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
			selector = &tb.ReplyMarkup{}

			// Reply buttons.
			btnHelp     = menu.Text("ℹ Help")
			btnSettings = menu.Text("⚙ Settings")
			btnUrl      = menu.URL("Visit", "https://google.com")

			// Inline buttons.
			//
			// Pressing it will cause the client to
			// send the bot a callback.
			//
			// Make sure Unique stays unique as per button kind,
			// as it has to be for callback routing to work.
			//
			btnPrev = selector.Data("⬅", "prev")
			btnNext = selector.Data("➡", "next")
		)

		menu = &tb.ReplyMarkup{ResizeReplyKeyboard: true}

		menu.Reply(
			menu.Row(btnHelp),
			menu.Row(btnSettings),
		)
		selector.Inline(
			selector.Row(btnPrev, btnNext, btnUrl),
		)

		s.bot.Send(m.Sender, "Hello!", selector)
	})

	// TODO: just for testing, remove
	s.bot.Handle(tb.OnQuery, func(q *tb.Query) {
		urls := []string{
			"https://images.app.goo.gl/fyD4H4g1iKuW27Ee8",
			"https://images.app.goo.gl/fyD4H4g1iKuW27Ee8",
		}

		results := make(tb.Results, len(urls)) // []tb.Result
		for i, url := range urls {
			result := &tb.PhotoResult{
				URL: url,

				// required for photos
				ThumbURL: url,
			}

			results[i] = result
			// needed to set a unique string ID for each result
			results[i].SetResultID(strconv.Itoa(i))
		}

		err := s.bot.Answer(q, &tb.QueryResponse{
			Results:   results,
			CacheTime: 60, // a minute
		})

		if err != nil {
			log.Println(err)
		}
	})

	// Send messages at intervals o time
	s.cron.MustAddJob("* * * * *", func() {
		s.bot.Send(&tb.User{ID: s.userChatID}, "hey there")
	})

	s.bot.Start()
}

func (s *service) SendMessageToUser(message string) error {
	_, err := s.bot.Send(&tb.User{ID: s.userChatID}, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) ToAmount(amount int, currency string) string {
	if amount < 0 {
		amount = -amount
	}

	if currency == "GBP" {
		currency = "£"
	}
	floatAmount := float64(amount) / 100

	return fmt.Sprintf("%s%.2f", currency, floatAmount)
}
