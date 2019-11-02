package main

import (
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func handleRequest(ctx *context) {
	if len(ctx.Arguments) == 0 {
		var roles = []string{
			"tester",
			"anotherTester",
		}

		var msg strings.Builder

		msg.WriteString("Available roles are: ")

		for i := 0; i < len(roles); i++ {
			msg.WriteString("\n`" + strconv.Itoa(i+1) + ": " + roles[i] + "`")
		}

		ctx.channelSend(msg.String())
		return
	}

	if ctx.Message.ChannelID == channelMap["BotTesting"] {
		userID := ctx.Message.Author.ID
		guildID := ctx.Message.GuildID
		roleID := "640291240834760726"

		// Hardcoded Siege role id
		ctx.Session.GuildMemberRoleAdd(guildID, userID, roleID)
		role, err := ctx.Session.State.Role(guildID, roleID)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error":   err,
				"command": ctx.Command,
			}).Error("Something went wrong when reading role")
		}

		var msg strings.Builder

		msg.WriteString("You have been granted the role of ")
		msg.WriteString(role.Name)
		msg.WriteString(". I trust that you can wield your newfound powers wisely")

		ctx.channelSend(msg.String())
	}
}

func handleTake(ctx *context) {
	if ctx.Message.ChannelID == channelMap["BotTesting"] {
		userID := ctx.Message.Author.ID
		guildID := ctx.Message.GuildID
		roleID := "640291240834760726"

		ctx.Session.GuildMemberRoleRemove(guildID, userID, roleID)
		role, err := ctx.Session.State.Role(guildID, roleID)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error":   err,
				"command": ctx.Command,
			}).Error("Something went wrong when reading role")
		}

		var msg strings.Builder

		// User @mention here?
		msg.WriteString("Alas, all good things must come to an end. You no longer have the role ")
		msg.WriteString(role.Name)

		ctx.channelSend(msg.String())
	}
}
