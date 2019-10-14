package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type (
	multiplexer struct {
		prefix, errorText string
		commands          map[string]func(ctx *context)
		logger            *logrus.Logger
		debug             bool
	}

	context struct {
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
		prefix:    prefix,
		errorText: errorText,
		commands:  make(map[string]func(ctx *context)),
		logger:    logger,
		debug:     debug,
	}, nil
}

// Handle handles commands. Called directly from dg.AddHandler()
func (m *multiplexer) handle(
	ses *discordgo.Session,
	msg *discordgo.MessageCreate,
) {
	if msg.Author.ID == ses.State.User.ID {
		return
	}

	/* TODO: The way arguments are handled here by splitting up slices is
	probably really inefficent. */
	args := strings.Split(msg.Message.Content, " ")
	if args[0][:1] != m.prefix {
		return
	}

	if m.debug {
		ch, _ := ses.Channel(msg.ChannelID)
		gu, _ := ses.Guild(msg.GuildID)
		m.logger.WithFields(logrus.Fields{
			"messageGuild":   gu.Name,
			"messageChannel": ch.Name,
			"messageAuthor":  msg.Author.Username,
			"messageContent": msg.Content,
		}).Info("Message Recieved")
	}

	handler, ok := m.commands[args[0][1:]]
	if !ok {
		ses.ChannelMessageSend(msg.ChannelID, m.errorText)
		return
	}

	go handler(&context{
		Arguments: args[1:],
		Session:   ses,
		Message:   msg,
	})
}

// Register adds a command to the bot
func (m *multiplexer) register(command string, handler func(ctx *context)) {
	m.commands[command] = handler
}
