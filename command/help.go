package command

import (
	"fmt"

	"github.com/CS-5/disgomux"
	"github.com/bwmarrin/discordgo"
)

// Help is a command
// TODO: Make this a better description
type Help struct {
	Command  string
	HelpText string
}

var (
	helpHandlers = make(map[string]func(ctx *disgomux.Context) bool)
	helpFields   []*discordgo.MessageEmbedField
)

// Init is called by the multiplexer before the bot starts to initialize any
// variables the command needs.
func (c Help) Init(m *disgomux.Mux) {
	i := 0
	for k, v := range m.Commands {
		msg := v.Settings().HelpText

		/* If there is no description, omit command from help */
		if len(msg) == 0 {
			continue
		}

		helpHandlers[k] = v.HandleHelp
		helpFields = append(helpFields, &discordgo.MessageEmbedField{
			Name:   m.Prefix + k,
			Value:  msg,
			Inline: true,
		})
		i++
	}

	commandLogs.Command.WithField("command", c.Command).Infof(
		"Loaded help handlers and messages for %d commands", i,
	)
}

// Handle is called by the multiplexer whenever a user triggers the command.
func (c Help) Handle(ctx *disgomux.Context) {
	if len(ctx.Arguments) == 0 {
		ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID,
			&discordgo.MessageEmbed{
				Title:       ":regional_indicator_h::regional_indicator_e::regional_indicator_l::regional_indicator_p:",
				Author:      &discordgo.MessageEmbedAuthor{},
				Color:       0xfdd329,
				Description: "Available commands:",
				Fields:      helpFields,
			})
		return
	}

	command, ok := helpHandlers[ctx.Arguments[0]]
	if !ok {
		ctx.ChannelSend(fmt.Sprintf(
			"Unable to find help info for command `%s`", ctx.Arguments[0],
		))
		return
	}

	command(ctx)
}

// HandleHelp is called by whatever help command is in place when a user enters
// "!help [command name]". If the help command is not being handled, return
// false.
func (c Help) HandleHelp(ctx *disgomux.Context) bool {
	ctx.ChannelSend("Are you sure _you_ don't need help?")
	return true
}

// Settings is called by the multiplexer on startup to process any settings
// associated with that command.
func (c Help) Settings() *disgomux.CommandSettings {
	return &disgomux.CommandSettings{
		Command:  c.Command,
		HelpText: c.HelpText,
	}
}

// Permissions is called by the multiplexer on startup to collect the list of
// permissions required to run the given command.
func (c Help) Permissions() *disgomux.CommandPermissions {
	return &disgomux.CommandPermissions{}
}
