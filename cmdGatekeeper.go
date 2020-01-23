package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/CS-5/disgomux"
)

type cGate struct {
	Command  string
	HelpText string
}

const (
	unknownCommand = "Unknown command. Usage: `!role [give|take] [Role Name]`"
	backendError   = "Backend error: `%q`"
)

func (g cGate) Init(m *disgomux.Mux) {
	// Nothing to init
}

func (g cGate) Handle(ctx *disgomux.Context) {
	guildID := ctx.Message.GuildID
	roles, err := ctx.Session.GuildRoles(guildID)
	if err != nil {
		cmdIssue(ctx, err, "There was a problem getting the roles of the guild")
		return
	}

	requestableRoles := make(map[string]string)
	printNames := []string{}
	for _, r := range roles {
		if strings.HasPrefix(r.Name, ":") {
			requestableRoles[strings.ToLower(r.Name[1:])] = r.ID
			printNames = append(printNames, r.Name[1:])
		}
	}

	/* If there are no arguments (Give/Take). Provide the user with options */
	if len(ctx.Arguments) == 0 {
		var msg strings.Builder
		msg.WriteString("Available roles are: ")

		for _, n := range printNames {
			msg.WriteString(fmt.Sprintf("\n- `%s`", n))
		}

		ctx.ChannelSend(msg.String())
		return
	}

	/* If there was an argument, was it give, take, or something invalid? */
	var give bool
	switch strings.ToLower(ctx.Arguments[0]) {
	case "give", "g":
		give = true
	case "take", "t":
		give = false
	default:
		ctx.ChannelSend(unknownCommand)
		return
	}

	/* If there was just one argument, inform the user */
	if len(ctx.Arguments) < 2 {
		ctx.ChannelSend(unknownCommand)
		return
	}

	/* Get the user ID and the user object */
	userID := ctx.Message.Author.ID
	member, err := ctx.Session.GuildMember(ctx.Message.GuildID, userID)
	if err != nil {
		cmdIssue(ctx, err, "There was a problem getting the user id")
	}

	/* Check to see if the requested role is valid */
	req := strings.ToLower(ctx.Arguments[1])

	roleID, ok := requestableRoles[req]
	if !ok {
		ctx.ChannelSend(
			fmt.Sprintf("Unable to find role `%s`", req),
		)
		return
	}

	hasRole := false
	if arrayContains(member.Roles, roleID, false) {
		hasRole = true
	}

	/* Give a role */
	if give {
		if hasRole {
			ctx.ChannelSend(fmt.Sprintf(
				"You appear to already have that role, %s", member.Mention(),
			))
			return
		}
		ctx.Session.GuildMemberRoleAdd(guildID, userID, roleID)
		ctx.ChannelSend(
			fmt.Sprintf(
				"You have been given role `%s`, %s", req, member.Mention(),
			),
		)
		return
	}

	/* Take a role */
	if !hasRole {
		ctx.ChannelSend(fmt.Sprintf(
			"You don't have that role... How do you expect me to take it, %s?",
			member.Mention(),
		))
		return
	}
	ctx.Session.GuildMemberRoleRemove(guildID, userID, roleID)
	ctx.ChannelSend(fmt.Sprintf(
		"Taking role `%s` away, %s", req, member.Mention(),
	))
}

func (g cGate) HandleHelp(ctx *disgomux.Context) bool {
	//TODO: Finish this
	return false
}

func (g cGate) Settings() *disgomux.CommandSettings {
	return &disgomux.CommandSettings{
		Command:  g.Command,
		HelpText: g.HelpText,
	}
}

func (g cGate) Permissions() *disgomux.CommandPermissions {
	return &disgomux.CommandPermissions{
		RoleIDs: config.permissions[g.Command],
	}
}

func idToName(ctx *disgomux.Context, id *string) (string, error) {
	role, err := ctx.Session.State.Role(ctx.Message.GuildID, *id)
	if err != nil {
		return "", err
	}
	return role.Name, nil
}

func nameToID(ctx *disgomux.Context, name *string) (string, error) {
	roles, err := ctx.Session.GuildRoles(ctx.Message.GuildID)
	if err != nil {
		return "", errors.New("Unable to fetch roles for this guild")
	}

	for _, r := range roles {
		if strings.ToLower(r.Name) == *name {
			return r.ID, nil
		}
	}
	return "", fmt.Errorf("Role name '%s' does not exist", *name)
}
