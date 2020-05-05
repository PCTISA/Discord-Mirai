package main

import (
	"fmt"
	"strings"

	"github.com/CS-5/disgomux"
	"github.com/bwmarrin/discordgo"
)

type cDebug struct {
	Command  string
	HelpText string
}

func (d cDebug) Init(m *disgomux.Mux) {
	// Nothing to init
}

func (d cDebug) Handle(ctx *disgomux.Context) {
	if len(ctx.Arguments) == 0 {
		d.HandleHelp(ctx)
		return
	}

	switch strings.ToLower(ctx.Arguments[0]) {
	case "config":
		var sb strings.Builder
		sb.WriteString(
			fmt.Sprintf("`Simple Commands: %+v`\n\n", config.simpleCommands))
		sb.WriteString(
			fmt.Sprintf("`Permissions: %+v`", config.permissions))
		ctx.ChannelSend(sb.String())
	case "args":
		ctx.ChannelSend(fmt.Sprintf("%+v", ctx.Arguments))
	case "stats":
		d.stats(ctx)
	default:
		ctx.ChannelSend("Debug")
		cmdIssue(ctx, nil, "Debug")
	}
}

func (d cDebug) HandleHelp(ctx *disgomux.Context) bool {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		"`%s%s config`: Returns the contents of the JSON config file.\n",
		ctx.Prefix, d.Command,
	))
	sb.WriteString(fmt.Sprintf(
		"`%s%s args`: Returns the supplied arguments.",
		ctx.Prefix, d.Command,
	))
	ctx.ChannelSend(sb.String())
	return true
}

func (d cDebug) Settings() *disgomux.CommandSettings {
	return &disgomux.CommandSettings{
		Command:  d.Command,
		HelpText: d.HelpText,
	}
}

func (d cDebug) Permissions() *disgomux.CommandPermissions {
	return &disgomux.CommandPermissions{
		RoleIDs: config.permissions[d.Command],
	}
}

func (d cDebug) stats(ctx *disgomux.Context) {
	if len(ctx.Arguments) > 1 {
		switch strings.ToLower(ctx.Arguments[1]) {
		case "cpu":
			ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID,
				d.generateStatsEmbed(
					"CPU Stats", "now-15m", "now", 1000, 500, 50,
				),
			)
		case "memory":
			ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID,
				d.generateStatsEmbed(
					"Memory Stats", "now-15m", "now", 1000, 500, 51,
				),
			)
		case "network":
			ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID,
				d.generateStatsEmbed(
					"Network Stats", "now-15m", "now", 1000, 500, 52,
				),
			)
		default:
			ctx.ChannelSendf("Unknown option `%s`", ctx.Arguments[1])
		}

		return
	}

	ctx.ChannelSendf("Monitoring URL: https://crsn.link/fzer0")
}

func (d cDebug) generateStatsEmbed(title, from, to string, width, height, panelID int) *discordgo.MessageEmbed {
	fsBase := "https://status.carsonseese.com/d/lWIounQWk/vm-docker?orgId=2&var-Container=TestServerBot"
	base := "https://status.carsonseese.com/render/d-solo/lWIounQWk/vm-docker?orgId=2&var-Container=TestServerBot"

	return &discordgo.MessageEmbed{
		Title: title,
		Description: fmt.Sprintf(
			"[Fullscreen](%s&from=%s&to=%s&panelId=%d&fullscreen)",
			fsBase, from, to, panelID,
		),

		Footer: &discordgo.MessageEmbedFooter{
			Text: "This image will not live update",
		},
		Color: 0xf55142,
		Image: &discordgo.MessageEmbedImage{
			URL: fmt.Sprintf(
				"%s&from=%s&to=%s&panelId=%d&width=%d&height=%d",
				base, from, to, panelID, width, height,
			),
		},
	}
}
