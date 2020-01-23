package main

import (
	"fmt"
	"strings"

	"github.com/CS-5/disgomux"
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

	switch ctx.Arguments[0] {
	case "config":
		var sb strings.Builder
		sb.WriteString(
			fmt.Sprintf("`Simple Commands: %+v`\n\n", config.simpleCommands))
		sb.WriteString(
			fmt.Sprintf("`Permissions: %+v`", config.permissions))
		ctx.ChannelSend(sb.String())
	case "args":
		ctx.ChannelSend(fmt.Sprintf("%+v", ctx.Arguments))
	default:
		ctx.ChannelSend("Debug")
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
