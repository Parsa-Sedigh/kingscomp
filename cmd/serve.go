package cmd

import (
	"github.com/Parsa-Sedigh/kingscomp/internal/repository"
	"github.com/Parsa-Sedigh/kingscomp/internal/repository/redis"
	"github.com/Parsa-Sedigh/kingscomp/internal/service"
	"github.com/Parsa-Sedigh/kingscomp/internal/telegram"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the telegram bot",
	Run:   serve,
}

func serve(cmd *cobra.Command, args []string) {
	_ = godotenv.Load()

	// set up repositories
	redisClient, err := redis.NewRedisClient(os.Getenv("REDIS_URL"))
	if err != nil {
		logrus.WithError(err).Fatalln("couldn't connect to the redis server")
	}

	accountRepository := repository.NewAccountRedisRepository(redisClient)

	// set up app
	app := service.NewApp(service.NewAccountService(accountRepository))

	tg, err := telegram.NewTelegram(app, os.Getenv("BOT_API"))
	if err != nil {
		logrus.WithError(err).Fatalln("couldn't connect to the telegram server")
	}

	tg.Start()
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
