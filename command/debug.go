package command

import (
	"fmt"
	"strings"

	"github.com/PulseDevelopmentGroup/0x626f74/multiplexer"
	"github.com/bwmarrin/discordgo"
)

// Debug is a command
// TODO: Make this a better description
type Debug struct {
	Command  string
	HelpText string
}

// Init is called by the multiplexer before the bot starts to initialize any
// variables the command needs.
func (c Debug) Init(m *multiplexer.Mux) {
	// Nothing to init
}

// Handle is called by the multiplexer whenever a user triggers the command.
func (c Debug) Handle(ctx *multiplexer.Context) {
	if len(ctx.Arguments) == 0 {
		c.HandleHelp(ctx)
		return
	}

	switch strings.ToLower(ctx.Arguments[0]) {
	case "config":

		var sb strings.Builder
		sb.WriteString(
			fmt.Sprintf("`Simple Commands: %+v`\n\n", commandConfig.SimpleCommands))
		sb.WriteString(
			fmt.Sprintf("`Permissions: %+v`", commandConfig.Permissions))
		ctx.ChannelSend(sb.String())

	case "args":
		ctx.ChannelSend(fmt.Sprintf("%+v", ctx.Arguments))
	case "stats":
		c.stats(ctx)
	default:
		ctx.ChannelSend("Debug")
	}
}

// HandleHelp is called by whatever help command is in place when a user enters
// "!help [command name]". If the help command is not being handled, return
// false.
func (c Debug) HandleHelp(ctx *multiplexer.Context) bool {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		"`%s%s config`: Returns the contents of the JSON config file.\n",
		ctx.Prefix, c.Command,
	))
	sb.WriteString(fmt.Sprintf(
		"`%s%s args`: Returns the supplied arguments.",
		ctx.Prefix, c.Command,
	))
	ctx.ChannelSend(sb.String())
	return true
}

// Settings is called by the multiplexer on startup to process any settings
// associated with that command.
func (c Debug) Settings() *multiplexer.CommandSettings {
	return &multiplexer.CommandSettings{
		Command:  c.Command,
		HelpText: c.HelpText,
		Permissions: &multiplexer.CommandPermissions{
			RoleIDs: commandConfig.Permissions.RoleIDs[c.Command],
		},
	}
}

func (c Debug) stats(ctx *multiplexer.Context) {
	if len(ctx.Arguments) > 1 {
		switch strings.ToLower(ctx.Arguments[1]) {
		case "cpu":
			ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID,
				c.generateStatsEmbed(
					"CPU Stats", "now-15m", "now", 1000, 500, 50,
				),
			)
		case "memory":
			ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID,
				c.generateStatsEmbed(
					"Memory Stats", "now-15m", "now", 1000, 500, 51,
				),
			)
		case "network":
			ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID,
				c.generateStatsEmbed(
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

func (c Debug) generateStatsEmbed(title, from, to string, width, height, panelID int) *discordgo.MessageEmbed {
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
