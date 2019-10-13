package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	goenv "github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

type (
	enviorment struct {
		Token string `env:"BOT_TOKEN"`
		Debug bool   `env:"DEBUG" envDefault:"false"`
	}
)

var (
	env = enviorment{}
	log *logrus.Logger
)

func init() {
	/* Parse enviorment variables */
	if err := goenv.Parse(&env); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", env)

	/* Define logging setup */
	logrus.SetOutput(os.Stdout)
	log = logrus.New()

	logrus.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	if !env.Debug {
		logrus.SetLevel(logrus.InfoLevel)
		log.SetFormatter(&logrus.JSONFormatter{})
	}
}

func main() {
	log.Info("Starting Bot...")
	dg, err := discordgo.New("Bot " + env.Token)
	if err != nil {
		log.WithField("error", err).Error("Problem starting bot")
	}

	log.Info("Bot started")

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		/* Log every message recieved */
		ch, _ := s.Channel(m.ChannelID)
		log.WithFields(logrus.Fields{
			"author":  m.Author.Username,
			"channel": ch.Name,
			"message": m.Content,
		}).Info("Message recieved")
	})

	err = dg.Open()
	if err != nil {
		log.WithField("error", err).Error(
			"Problem opening websocket connection.",
		)
		return
	}
	defer dg.Close()

	/* Wait for interrupt */
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
