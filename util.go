package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/CS-5/disgomux"
	"github.com/sirupsen/logrus"
)

func initFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)

		if err != nil {
			return &os.File{}, err
		}
		return file, err
	}
	return file, err
}

func initLogging(debug bool) *logrus.Logger {
	logrus.SetOutput(os.Stdout)
	log := logrus.New()

	logrus.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	if !debug {
		logrus.SetLevel(logrus.InfoLevel)
		log.SetFormatter(&logrus.JSONFormatter{})
	}
	return log
}

func arrayContains(array []string, value string, ignoreCase bool) bool {
	for _, e := range array {
		if ignoreCase {
			e = strings.ToLower(e)
		}

		if e == value {
			return true
		}
	}
	return false
}

func cmdIssue(ctx *disgomux.Context, e error, msg string) {
	cLog.WithError(e).WithField("command", ctx.Command).Error(msg)

	/* Send to error message to designated channel. Not error handling here since
	   errors thrown are effectivly meaningless */
	channel, _ := ctx.Session.Channel(ctx.Message.ChannelID)
	ctx.Session.ChannelMessageSend(config.errChan, fmt.Sprintf(
		"**Error:**\n```Channel: %s\nCommand: %s\nError: %v\nMessage: %s```",
		channel.Name, ctx.Command, e, msg,
	))

	if env.Debug {
		ctx.ChannelSendf("Error: Message: `%s`, Error Message: `%v`", msg, e)
		return
	}
	ctx.ChannelSend(msg)
}
