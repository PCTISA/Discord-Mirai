package main

import (
	"errors"
	"fmt"
	"strings"
)

const (
	unknownCommand = "Unknown command. Usage: `!role [give|take] [Role Name]`"
	backendError   = "Backend error: `%q`"
)

func idToName(ctx *context, id *string) (string, error) {
	role, err := ctx.Session.State.Role(ctx.Message.GuildID, *id)
	if err != nil {
		return "", err
	}
	return role.Name, nil
}

func nameToID(ctx *context, name *string) (string, error) {
	roles, err := ctx.Session.GuildRoles(ctx.Message.GuildID)
	if err != nil {
		return "", errors.New("Unable to fetch roles for this guild")
	}

	for _, r := range roles {
		if strings.ToLower(r.Name) == *name {
			return r.ID, nil
		}
	}
	return "", fmt.Errorf("Role name %s does not exist", *name)
}

/* TODO: Maybe move this to util if we need such functionality elsewhere? */
func arrayContains(array []string, value string) bool {
	for _, e := range array {
		if e == value {
			return true
		}
	}
	return false
}

func handleGatekeeper(ctx *context) {
	/* Get role IDs and names */
	roleIDs := config.requestableRoles
	roleNames := []string{}
	for _, id := range roleIDs {
		name, err := idToName(ctx, &id)
		if err != nil {
			log.WithField("error", err).Errorf(
				"Problem converting ID: %q to name", id,
			)
			ctx.channelSend("There was a problem executing the command")
			return
		}
		roleNames = append(roleNames, strings.ToLower(name))
	}

	/* If there are no arguments (Give/Take). Provide the user with options */
	if len(ctx.Arguments) == 0 {

		var msg strings.Builder
		msg.WriteString("Available roles are: ")

		for _, n := range roleNames {
			msg.WriteString(fmt.Sprintf("\n- `%s`", n))
		}

		ctx.channelSend(msg.String())
		return
	}

	/* If there was an argument, was it give, take, or something invalid? */
	var give bool
	switch strings.ToLower(ctx.Arguments[0]) {
	case "give":
		give = true
	case "take":
		give = false
	default:
		ctx.channelSend(unknownCommand)
		return
	}

	/* If there was just one argument, inform the user */
	if len(ctx.Arguments) < 2 {
		ctx.channelSend(unknownCommand)
		return
	}

	/* Get the user ID and the user object */
	userID := ctx.Message.Author.ID
	user, err := ctx.Session.User(userID)
	if err != nil {
		ctx.channelSend(fmt.Sprintf(backendError, err))
	}

	/* Check to see if the requested role is valid */
	req := strings.ToLower(ctx.Arguments[1])

	if !arrayContains(
		roleNames, req,
	) {
		ctx.channelSend(
			fmt.Sprintf("Unable to give/take role %s, %s", req, user.Mention()),
		)
		return
	}

	/* Get the guild and role IDs */
	guildID := ctx.Message.GuildID
	roleID, err := nameToID(ctx, &req)
	if err != nil {
		ctx.channelSend(fmt.Sprintf(backendError, err))
		return
	}

	/* Give a role */
	if give {
		ctx.Session.GuildMemberRoleAdd(guildID, userID, roleID)
		ctx.channelSend(
			fmt.Sprintf("You have been given role %s, %s", req, user.Mention()),
		)

		return
	}

	/* Take a role */
	ctx.Session.GuildMemberRoleRemove(guildID, userID, roleID)
	ctx.channelSend(fmt.Sprintf("Taking role %s away, %s", req, user.Mention()))
}
