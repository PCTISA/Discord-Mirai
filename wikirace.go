package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
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
	}

	defer resp.Body.Close()
	if err != nil {
		log.Error(err)
		s.ChannelMessageSend(m.ChannelID, "Hmm... I seem to have run into an issue... Try again later?")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		s.ChannelMessageSend(m.ChannelID, "Hmm... I seem to have run into an issue... Try again later?")
	}

	var search wikiResult
	err = json.Unmarshal(body, &search)
	if err != nil {
		log.Error(err)
		s.ChannelMessageSend(m.ChannelID, "Hmm... I seem to have run into an issue... Try again later?")
	}

	articles := search.Query.Random[:2]

	// @carson is there a better way of doing this?
	// <> characters prevent the embed
	msg := "Race starts at " + articles[0].Title + " (<https://en.wikipedia.org/?curid=" + strconv.Itoa(articles[0].ID) + ">)" + " and goes to " + articles[1].Title + " (<https://en.wikipedia.org/?curid=" + strconv.Itoa(articles[1].ID) + ">)" + "."

	s.ChannelMessageSend(m.ChannelID, msg)
}
