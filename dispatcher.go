package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type (
	multiplexer struct {
		Prefix, ErrorText string
		Commands          map[string]func(ctx *context)
		HelpText 		  map[string]string
		Logger            *logrus.Logger
		Debug             bool
	}

	context struct {
		Command   string
		Arguments []string
		Session   *discordgo.Session
		Message   *discordgo.MessageCreate
	}
)

// NewMux creates a new multiplexer.
func newMux(
	prefix, errorText string,
	logger *logrus.Logger,
	debug bool,
) (multiplexer, error) {
	if len(prefix) > 1 {
		return multiplexer{},
			fmt.Errorf("prefix %q longer than max length of 1", prefix)
	}

	/* TODO: Make errorText optional..
	Would have to check for it's existance elsewhere */
	if len(errorText) == 0 {
		return multiplexer{}, fmt.Errorf("error text %q nonexistant", errorText)
	}

	if logger == nil {
		return multiplexer{},
			/* Technically not a use for Errorf,
			but better than importing errors */
			fmt.Errorf("logger invalid")
	}

	return multiplexer{
		Prefix:    prefix,
		ErrorText: errorText,
		Commands:  make(map[string]func(ctx *context)),
		HelpText:  make(map[string]string),
		Logger:    logger,
		Debug:     debug,
	}, nil
}

// Handle handles commands. Called directly from dg.AddHandler()
func (m *multiplexer) handle(
	session *discordgo.Session,
	message *discordgo.MessageCreate,
) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	/* TODO: The way arguments are handled here by splitting up slices is
	probably really inefficent. */
	args := strings.Split(message.Content, " ")
	if args[0][:1] != m.Prefix {
		return
	}

	if m.Debug {
		ch, _ := session.Channel(message.ChannelID)
		gu, _ := session.Guild(message.GuildID)
		m.Logger.WithFields(logrus.Fields{
			"messageGuild":   gu.Name,
			"messageChannel": ch.Name,
			"messageAuthor":  message.Author.Username,
			"messageContent": message.Content,
		}).Info("Message Recieved")
	}

	handler, ok := m.Commands[args[0][1:]]
	if !ok {
		session.ChannelMessageSend(message.ChannelID, m.ErrorText)
		return
	}

	go handler(&context{
		Command:   args[0][1:],
		Arguments: args[1:],
		Session:   session,
		Message:   message,
	})
}

// Register adds a command to the bot
func (m *multiplexer) register(
	command, helpText string, 
	handler func(ctx *context), 
) error {
	if len(command) == 0 {
		return fmt.Errorf("Command '%v' too short", command)
	}

	if len(helpText) == 0 {
		return fmt.Errorf("Help text '%v' too short", helpText)
	}

	m.Commands[command] = handler
	m.HelpText[command] = helpText
	return nil
}

// HandleHelp adds a !help command with auto-generated output. Must be called 
// after all register commands
func (m *multiplexer) handleHelp(description string) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%v\n", description))

	for k, v := range m.HelpText {
		b.WriteString(fmt.Sprintf("`!%v`: %v\n", k, v))
	}

	m.register("help", "Lists all commands and their functions.", 
		func(ctx *context) {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, b.String())
		},
	)
}
