package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

	var msgBuilder strings.Builder

	msgBuilder.WriteString("Race starts at ")
	msgBuilder.WriteString(articles[0].Title)
	// <> characters prevent embed
	msgBuilder.WriteString(" (<https://en.wikipedia.org/?curid=" + strconv.Itoa(articles[0].ID) + ">)")
	msgBuilder.WriteString(" and goes to ")
	msgBuilder.WriteString(articles[1].Title + " (<https://en.wikipedia.org/?curid=" + strconv.Itoa(articles[1].ID) + ">)")
	msgBuilder.WriteString(".")

	s.ChannelMessageSend(m.ChannelID, msgBuilder.String())
}
