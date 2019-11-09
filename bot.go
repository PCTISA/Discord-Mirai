package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	goenv "github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

type (
	environment struct {
		Token   string `env:"BOT_TOKEN"`
		Debug   bool   `env:"DEBUG" envDefault:"false"`
		DataDir string `env:"DATA_DIR" envDefault:"data/"`
	}
)

var (
	env    = environment{}
	log    *logrus.Logger
	config *botConfig
)

func init() {
	/* Parse enviorment variables */
	if err := goenv.Parse(&env); err != nil {
		fmt.Printf("%+v\n", err)
	}

	/* Define logging setup */
	log = initLogging(env.Debug)

	/* Parse config */
	var err error
	config, err = getConfig(env.DataDir + "config.json")
	if err != nil {
		log.WithField("error", err).Error("Problem executing config command")
	}
}

func main() {
	/* Initialize DiscordGo */
	log.Info("Starting Bot...")
	dg, err := discordgo.New("Bot " + env.Token)
	if err != nil {
		log.WithField("error", err).Error("Problem starting bot")
	}
	log.Info("Bot started")

	/* Initialize mux */
	mux, err := newMux("!", "Unknown command D:", log, env.Debug)
	if err != nil {
		log.WithField("error", err).Fatalf("Unable to create multiplexer")
	}

	/* --- Register all the things --- */

	mux.register("test", "", func(ctx *context) {
		ctx.channelSend(fmt.Sprintf("%+v", ctx.Arguments))
	})

	mux.register("config", "", func(ctx *context) {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("`Requestable Roles: %+v`\n", config.requestableRoles))
		sb.WriteString(fmt.Sprintf("`Simple Commands: %+v`\n", config.simpleCommands))
		sb.WriteString(fmt.Sprintf("`Permissions: %+v`", config.permissions))

		ctx.channelSend(sb.String())
	})

	mux.register("wikirace", "Start a wikirace", handleWikirace)

	mux.handleHelp("Available commands:")

	/* --- End Register --- */

	/* Handle commands and start DiscordGo */
	dg.AddHandler(mux.handle)
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
