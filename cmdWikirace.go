package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/CS-5/disgomux"
	"github.com/bwmarrin/discordgo"
)

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
	resp, err := http.Get("https://en.wikipedia.org/w/api.php?action=query&format=json&list=random&rnnamespace=0&rnlimit=2")
	if err != nil {
		cmdIssue(ctx, err, "Unable to get random wikipedia page")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		cmdIssue(ctx, err, "Unable to read page")
		return
	}

	var search wikiResult
	err = json.Unmarshal(body, &search)
	if err != nil {
		cmdIssue(ctx, err, "Unable to unmarshal page")
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

func (w cWiki) HandleHelp(ctx *disgomux.Context) bool {
	var sb strings.Builder
	sb.WriteString(
		"Use `!wikirace` to start a new race! The rules are simple:\n",
	)
	sb.WriteString("1. Only blue links _within_ the article are allowed\n")
	sb.WriteString("2. You cannot use the back button or the search function\n")
	sb.WriteString("3. Whoever gets to end article in the fewest clicks wins\n")

	ctx.ChannelSend(sb.String())
	return true
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
