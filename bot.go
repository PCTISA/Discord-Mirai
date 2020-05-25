// This codebase has really turned into a disaster. It should really totally be
// reworked... Eh, maybe someday

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/PulseDevelopmentGroup/0x626f74/command"
	"github.com/PulseDevelopmentGroup/0x626f74/config"
	"github.com/PulseDevelopmentGroup/0x626f74/log"

	"github.com/CS-5/disgomux"

	"github.com/bwmarrin/discordgo"
	goenv "github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
)

type environment struct {
	Token   string `env:"BOT_TOKEN"`
	Debug   bool   `env:"DEBUG" envDefault:"false"`
	DataDir string `env:"DATA_DIR" envDefault:"data/"`
	Fuzzy   bool   `env:"USE_FUZZY" envDefault:"false"`
}

var (
	env  = environment{}
	cfg  *config.BotConfig
	logs *log.Logs

	prefix = "!"
)

func init() {
	/* Parse enviorment variables */
	if err := goenv.Parse(&env); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	/* Parse config */
	var err error
	cfg, err = config.Get(env.DataDir + "config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	/* Define logging setup */
	logs = log.New(env.Debug)
}

func main() {
	/* Initialize DiscordGo */
	logs.Primary.Info("Starting Bot...")
	dg, err := discordgo.New("Bot " + env.Token)
	if err != nil {
		logs.Primary.WithError(err).Error("Problem starting bot")
	}
	logs.Primary.Info("Bot started")

	/* Initialize Mux */
	dMux, err := disgomux.New(prefix)
	if err != nil {
		logs.Primary.WithError(err).Fatalf("Unable to create multixplexer")
	}

	dMux.UseMiddleware(logs.MuxMiddleware)

	/* Setup Errors */
	dMux.SetErrors(disgomux.ErrorTexts{
		CommandNotFound: "Command not found.",
		NoPermissions:   "You do not have permissions to execute that command.",
	})

	/* === Register all the things === */

	command.InitGlobals(cfg, logs)

	dMux.Register(
		command.Debug{
			Command:  "debug",
			HelpText: "Debuging info for bot-wranglers",
		},
		command.Wiki{
			Command:  "wikirace",
			HelpText: "Start a wikirace",
		},
		command.Gatekeeper{
			Command:  "role",
			HelpText: "Manage your access to roles, and their related channels",
		},
		command.Help{
			Command:  "help",
			HelpText: "Displays help information regarding the bot's commands",
		},
		command.Inspire{
			Command:  "inspire",
			HelpText: "Get an inspirational quote from inspirobot.me (use at your own risk)",
		},
		command.JPEG{
			Command:  "jpeg",
			HelpText: "More JPEG for the last image. 'nuff said",
		},
	)

	dMux.Options(&disgomux.Options{
		IgnoreDMs:        true,
		IgnoreBots:       true,
		IgnoreNonDefault: true,
		IgnoreEmpty:      true,
	})

	dMux.Initialize()

	if env.Fuzzy {
		dMux.InitializeFuzzy()
	}

	/* Register commands from the config file */
	for k := range cfg.SimpleCommands {
		k := k
		dMux.RegisterSimple(disgomux.SimpleCommand{
			Command:  k,
			Content:  cfg.SimpleCommands[k],
			HelpText: "This is a simple command",
		})
	}

	/* === End Register === */

	/* Handle commands and start DiscordGo */
	dg.AddHandler(dMux.Handle)

	err = dg.Open()
	if err != nil {
		logs.Primary.WithError(err).Error(
			"Problem opening websocket connection.",
		)
		return
	}

	idle := 0
	dg.UpdateStatusComplex(discordgo.UpdateStatusData{
		IdleSince: &idle,
		Game: &discordgo.Game{
			Name: "you",
			Type: discordgo.GameTypeWatching,
			Assets: discordgo.Assets{
				LargeImageID: "watching",
				LargeText:    "Watching...",
			},
		},
		Status: "online",
	})

	defer dg.Close()

	/* Wait for interrupt */
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
