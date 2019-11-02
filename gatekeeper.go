package main

import (
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func findRoleByIndex(ctx *context) (string, error) {
	roleIndex, err := strconv.Atoi(ctx.Arguments[0])
	if err != nil {
		ctx.channelSend("Sorry, I don't understand which role you want :(")
		return "", err
	}

	roles, err := ctx.Session.GuildRoles(ctx.Message.GuildID)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":   err,
			"command": ctx.Command,
		}).Error("Something went wrong when reading role")
	}

	roleID := roles[roleIndex-1].ID

	return roleID, nil
}

func isValidChannel(ctx *context) bool {
	return ctx.Message.ChannelID == channelMap["BotTesting"] || ctx.Message.ChannelID == channelMap["BotSpam"]
}

func handleRequest(ctx *context) {
	if len(ctx.Arguments) == 0 {
		roles, err := ctx.Session.GuildRoles(ctx.Message.GuildID)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error":   err,
				"command": ctx.Command,
			}).Error("Something went wrong when reading role")
		}

		var msg strings.Builder

		msg.WriteString("Available roles are: ")

		for i := 0; i < len(roles); i++ {
			msg.WriteString("\n- `" + strconv.Itoa(i+1) + ": " + roles[i].Name + "`")
		}

		ctx.channelSend(msg.String())
		return
	}

	if isValidChannel(ctx) {
		userID := ctx.Message.Author.ID
		guildID := ctx.Message.GuildID

		roleID, err := findRoleByIndex(ctx)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error":   err,
				"command": ctx.Command,
			}).Error("Unable to find desired role. Was it deleted?")
		}

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
	if isValidChannel(ctx) {
		userID := ctx.Message.Author.ID
		guildID := ctx.Message.GuildID

		roleID, err := findRoleByIndex(ctx)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error":   err,
				"command": ctx.Command,
			}).Error("Unable to find desired role. Was it deleted?")
		}

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
