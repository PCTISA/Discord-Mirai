package main

import (
	"github.com/CS-5/disgomux"
	"github.com/sirupsen/logrus"
)

type muxLog struct {
	logger      *logrus.Logger
	logMessages bool
}

func (ml muxLog) Init(mux *disgomux.Mux) {

}

func (ml muxLog) MessageRecieved(ctx *disgomux.Context) {
	if ml.logMessages {
		ch, _ := ctx.Session.Channel(ctx.Message.ChannelID)
		gu, _ := ctx.Session.Guild(ctx.Message.GuildID)

		ml.logger.WithFields(logrus.Fields{
			"messageGuild":   gu.Name,
			"messageChannel": ch.Name,
			"messageAuthor":  ctx.Message.Author.Username,
			"messageContent": ctx.Message.Content,
		}).Info("Message Recieved")
	}
}

func (ml muxLog) Info(ctx *disgomux.Context, message string) {
	ml.logger.Info(message)
}

func (ml muxLog) Warn(ctx *disgomux.Context, message string) {
	ml.logger.Warn(message)
}

func (ml muxLog) Error(ctx *disgomux.Context, message string) {
	ml.logger.Error(message)
}

func (ml muxLog) Done(mux *disgomux.Mux) {

}
