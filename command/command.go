package command

import (
	"github.com/PCTISA/Discord-Mirai/config"
	"github.com/PCTISA/Discord-Mirai/log"
	"github.com/PCTISA/Discord-Mirai/multiplexer"
)

/*


TODO: The existance of this file and it's functionality is horrifying, but
      I really have no better solution? Will definitely have to revisit this.


*/

var (
	cfg  *config.BotConfig
	logs *log.Logs
)

// InitGlobals is used to set global variables used by all commands. Must be
// called before commands are initialized or you'll have problems.
// TODO: Using global variables may be bad practice here? -- Confirmed
func InitGlobals(config *config.BotConfig, log *log.Logs) {
	cfg = config
	logs = log
}

// RegisterSimple registers any simple commands in the config struct. Must be
// called after InitGlobals().
func RegisterSimple(mux *multiplexer.Mux) {
	for k := range cfg.SimpleCommands {
		k := k
		mux.RegisterSimple(multiplexer.SimpleCommand{
			Command:  k,
			Content:  cfg.SimpleCommands[k],
			HelpText: "This is a simple command",
		})
	}
}

// CmdErr is used for handling errors within commands which should be reported
// to the user. Takes a multiplexer context, error message, and user-readable
// message which are sent to the channel where the command was executed.
func cmdErr(ctx *multiplexer.Context, err error, msg string) {
	ctx.ChannelSendf(
		"%s Something weird happened. Might need to report this on Github\nError:```%s```", msg, err.Error(),
	)
	logs.Command.Error(err)
}
