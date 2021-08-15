package cmd

import (
	"os"

	"github.com/MihaiBlebea/go-wallet/bot"
	"github.com/MihaiBlebea/go-wallet/domain"
	"github.com/MihaiBlebea/go-wallet/domain/account"
	"github.com/MihaiBlebea/go-wallet/domain/receipt"
	"github.com/MihaiBlebea/go-wallet/domain/transaction"

	server "github.com/MihaiBlebea/go-wallet/server"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use: "start",
	RunE: func(cmd *cobra.Command, args []string) error {

		l := logrus.New()

		l.SetFormatter(&logrus.JSONFormatter{})
		l.SetOutput(os.Stdout)
		l.SetLevel(logrus.InfoLevel)

		conn, err := gorm.Open(sqlite.Open(os.Getenv("SQLITE_PATH")), &gorm.Config{})
		if err != nil {
			return err
		}

		conn.AutoMigrate(
			&account.Account{},
			&receipt.Receipt{},
			&transaction.Transaction{},
		)

		wallet := domain.New(conn)

		telebot, err := bot.New(wallet)
		if err != nil {
			return err
		}

		go telebot.Start()

		server.NewServer(wallet, telebot, l)

		return nil
	},
}
