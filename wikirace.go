package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type articleInfo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type wikiResult struct {
	BatchComplete string `json:"batchComplete"`
	Query         struct {
		Random []articleInfo `json:"random"`
	} `json:"query"`
}

func handleWikirace(ctx *context) {
	/* TODO: Maybe float these erros up to the handler? */
	resp, err := http.Get("https://en.wikipedia.org/w/api.php?action=query&format=json&list=random&rnnamespace=0&rnlimit=2")
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":   err,
			"command": ctx.Command,
		}).Error("Unable to get random wikipedia page")

		ctx.channelSend("Hmm... I seem to have run into an issue... Try again later?")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":   err,
			"command": ctx.Command,
		}).Error("Unable to read page")

		ctx.channelSend("Hmm... I seem to have run into an issue... Try again later?")
		return
	}

	var search wikiResult
	err = json.Unmarshal(body, &search)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":   err,
			"command": ctx.Command,
		}).Error("Unable to unmarshal page")

		ctx.channelSend("Hmm... I seem to have run into an issue... Try again later?")
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

	ctx.channelSend(msgBuilder.String())
}
