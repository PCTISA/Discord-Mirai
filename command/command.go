package command

import (
	"github.com/PulseDevelopmentGroup/0x626f74/config"
	"github.com/PulseDevelopmentGroup/0x626f74/log"
)

var (
	commandConfig *config.BotConfig
	commandLogs   *log.Logs
)

// InitGlobals is used to set global variables used by all commands. Must be
// called before commands are initialized or you'll have problems.
func InitGlobals(config *config.BotConfig, logs *log.Logs) {
	commandConfig = config
	commandLogs = logs
}
