package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type articleInfo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type queryResult struct {
	Random []articleInfo `json:"random"`
}

type wikiResult struct {
	BatchComplete string      `json:"batchComplete"`
	Query         queryResult `json:"query"`
}

func initWikiRace(ctx *context) {
	m := ctx.Message
	s := ctx.Session

	resp, err := http.Get("https://en.wikipedia.org/w/api.php?action=query&format=json&list=random&rnnamespace=0&rnlimit=2")
	if err != nil {
		log.Error(err)
		s.ChannelMessageSend(m.ChannelID, "Hmm... I seem to have run into an issue... Try again later?")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		s.ChannelMessageSend(m.ChannelID, "Hmm... I seem to have run into an issue... Try again later?")
		return
	}

	var search wikiResult
	err = json.Unmarshal(body, &search)
	if err != nil {
		log.Error(err)
		s.ChannelMessageSend(m.ChannelID, "Hmm... I seem to have run into an issue... Try again later?")
		return
	}

	articles := search.Query.Random[:2]

	ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, &discordgo.MessageEmbed{
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
