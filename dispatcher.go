package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type (
	multiplexer struct {
		prefix, errorText string
		commands          map[string]func(ctx *context)
	}

	context struct {
		Arguments []string
		Session   *discordgo.Session
		Message   *discordgo.MessageCreate
	}
)

func newMux(p, et string) (*multiplexer, error) {
	if len(p) > 1 {
		return &multiplexer{},
			fmt.Errorf("prefix %q longer than max length of 1", p)
	}

	return &multiplexer{
		prefix:    p,
		errorText: et,
		commands:  make(map[string]func(ctx *context)),
	}, nil
}

func (m multiplexer) handle(ses *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == ses.State.User.ID {
		return
	}

	args := strings.Split(msg.Message.Content, " ")

	if args[0][:1] != m.prefix {
		return
	}

	handler, ok := m.commands[args[0][1:]]
	if !ok {
		ses.ChannelMessageSend(msg.ChannelID, m.errorText)
		return
	}

	handler(&context{
		Arguments: args,
		Session:   ses,
		Message:   msg,
	})
}

// Register adds a command to the bot
func (m multiplexer) register(cmd string, handler func(ctx *context)) {
	m.commands[cmd] = handler
}
