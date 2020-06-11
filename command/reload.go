package command

import (
	"github.com/PulseDevelopmentGroup/0x626f74/multiplexer"
)

// Reload is a command
// TODO: Make this a better description
type Reload struct {
	Command  string
	HelpText string

	Mux *multiplexer.Mux
}

// Init is called by the multiplexer before the bot starts to initialize any
// variables the command needs.
func (c Reload) Init(m *multiplexer.Mux) {
	// Nothing to init
}

// Handle is called by the multiplexer whenever a user triggers the command.
func (c Reload) Handle(ctx *multiplexer.Context) {
	err := cfg.Update()
	if err != nil {
		cmdErr(ctx, err, "Unable to update config")
	}

	/* Re-init simple commands */
	c.Mux.ClearSimple()
	RegisterSimple(c.Mux)

	/* Reload permissions */
	c.Mux.SetPermissions(cfg.Permissions)

	ctx.ChannelSend("Done")
}

// HandleHelp is called by whatever help command is in place when a user enters
// "!help [command name]". If the help command is not being handled, return
// false.
func (c Reload) HandleHelp(ctx *multiplexer.Context) bool {
	ctx.ChannelSend("Reload the bot's permissions")
	return true
}

// Settings is called by the multiplexer on startup to process any settings
// associated with that command.
func (c Reload) Settings() *multiplexer.CommandSettings {
	return &multiplexer.CommandSettings{
		Command:  c.Command,
		HelpText: c.HelpText,
	}
}
