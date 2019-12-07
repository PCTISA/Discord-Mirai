package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/CS-5/disgomux"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

/* === Start Debug Command === */

type cDebug struct {
	Command  string
	HelpText string
}

func (d cDebug) Init(m *disgomux.Mux) {
	// Nothing to init
}

func (d cDebug) Handle(ctx *disgomux.Context) {
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

func (d cDebug) HandleHelp(ctx *disgomux.Context) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		"`%s%s config`: Returns the contents of the JSON config file.\n",
		ctx.Prefix, ctx.Command,
	))
	sb.WriteString(fmt.Sprintf(
		"`%s%s args`: Returns the supplied arguments.",
		ctx.Prefix, ctx.Command,
	))
	ctx.ChannelSend(sb.String())
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

/* === End Debug Command === */

/* === Start WikiRace Command === */

type (
	cWiki struct {
		Command  string
		HelpText string
	}

	articleInfo struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}

	wikiResult struct {
		BatchComplete string `json:"batchComplete"`
		Query         struct {
			Random []articleInfo `json:"random"`
		} `json:"query"`
	}
)

const issueText = "Hmm... I seem to have run into an issue... Try again later?"

func (w cWiki) Init(m *disgomux.Mux) {
	// Nothing to init
}

func (w cWiki) Handle(ctx *disgomux.Context) {
	/* TODO: Maybe float these erros up to the handler? */
	resp, err := http.Get("https://en.wikipedia.org/w/api.php?action=query&format=json&list=random&rnnamespace=0&rnlimit=2")
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":   err,
			"command": ctx.Command,
		}).Error("Unable to get random wikipedia page")

		ctx.ChannelSend(issueText)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":   err,
			"command": ctx.Command,
		}).Error("Unable to read page")

		ctx.ChannelSend(issueText)
		return
	}

	var search wikiResult
	err = json.Unmarshal(body, &search)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":   err,
			"command": ctx.Command,
		}).Error("Unable to unmarshal page")

		ctx.ChannelSend(issueText)
		return
	}

	articles := search.Query.Random[:2]

	ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID,
		&discordgo.MessageEmbed{
			Title:       "Wikipedia Race",
			Author:      &discordgo.MessageEmbedAuthor{},
			Color:       0x0080ff,
			Description: "Start at the start and use only blue links in the article to get to the end page!",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name: "Start:vertical_traffic_light:",
					Value: fmt.Sprintf(
						"[%s](%s%d)",
						articles[0].Title,
						"https://en.wikipedia.org/?curid=",
						articles[0].ID,
					),
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name: "End :checkered_flag:",
					Value: fmt.Sprintf(
						"[%s](%s%d)",
						articles[1].Title,
						"https://en.wikipedia.org/?curid=",
						articles[1].ID,
					),
					Inline: false,
				},
			},
		})
}

func (w cWiki) HandleHelp(ctx *disgomux.Context) {
	// TODO: Finish this
}

func (w cWiki) Settings() *disgomux.CommandSettings {
	return &disgomux.CommandSettings{
		Command:  w.Command,
		HelpText: w.HelpText,
	}
}

func (w cWiki) Permissions() *disgomux.CommandPermissions {
	return &disgomux.CommandPermissions{
		RoleIDs: config.permissions[w.Command],
	}
}

/* === End WikiRace Command */

/* === Start Gatekeeper Command === */

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
		log.WithField("error", err).Errorf(
			"Problem getting roles for guild `%v`", guildID,
		)
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
		ctx.ChannelSend(fmt.Sprintf(backendError, err))
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

func (g cGate) HandleHelp(ctx *disgomux.Context) {
	//TODO: Finish this
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

/* === End Gatekeeper Command === */
