package command

import (
	"fmt"
	"strings"

	"github.com/PCTISA/Discord-Mirai/multiplexer"
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
			fmt.Sprintf("`Simple Commands: %+v`\n\n", cfg.SimpleCommands))
		sb.WriteString(
			fmt.Sprintf("`Permissions: %+v`", cfg.Permissions))
		ctx.ChannelSend(sb.String())

	case "args":
		ctx.ChannelSend(fmt.Sprintf("%+v", ctx.Arguments))
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
	}
}
