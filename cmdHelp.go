package main

import (
	"fmt"

	"github.com/CS-5/disgomux"
	"github.com/bwmarrin/discordgo"
)

type cHelp struct {
	Command  string
	HelpText string
}

var (
	helpHandlers = make(map[string]func(ctx *disgomux.Context) bool)
	helpFields   []*discordgo.MessageEmbedField
)

func (h cHelp) Init(m *disgomux.Mux) {
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

	cLog.WithField("command", h.Command).Infof(
		"Loaded help handlers and messages for %d commands", i,
	)
}

func (h cHelp) Handle(ctx *disgomux.Context) {
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

func (h cHelp) HandleHelp(ctx *disgomux.Context) bool {
	ctx.ChannelSend("Are you sure _you_ don't need help?")
	return true
}

func (h cHelp) Settings() *disgomux.CommandSettings {
	return &disgomux.CommandSettings{
		Command:  h.Command,
		HelpText: h.HelpText,
	}
}

func (h cHelp) Permissions() *disgomux.CommandPermissions {
	return &disgomux.CommandPermissions{}
}
