package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CS-5/disgomux"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

type (
	statistics struct {
		guildMsgs, chanMsgs, userMsgs, cmdsMsgs *cache.Cache
		logger                                  *logrus.Entry
	}

	msgStats struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
)

func statsInit(exp time.Duration, log *logrus.Entry) *statistics {
	s := &statistics{
		guildMsgs: cache.New(exp, exp),
		chanMsgs:  cache.New(exp, exp),
		userMsgs:  cache.New(exp, exp),
		cmdsMsgs:  cache.New(exp, exp),
		logger:    log,
	}

	s.logger.Info("Initialized stats collector.")

	go s.startHTTP()

	return s
}

func (s *statistics) startHTTP() {
	http.HandleFunc("/stats/channels", func(w http.ResponseWriter, r *http.Request) {
		payload, err := json.Marshal(s.buildStats(s.chanMsgs.Items()))
		if err != nil {
			s.logger.WithError(err).Warn(
				"Problem marshaling channel message stats",
			)
		}
		w.Write(payload)
	})
	http.HandleFunc("/stats/users", func(w http.ResponseWriter, r *http.Request) {
		payload, err := json.Marshal(s.buildStats(s.userMsgs.Items()))
		if err != nil {
			s.logger.WithError(err).Warn(
				"Problem marshaling users message stats",
			)
		}
		w.Write(payload)
	})
	http.HandleFunc("/stats/guilds", func(w http.ResponseWriter, r *http.Request) {
		payload, err := json.Marshal(s.buildStats(s.guildMsgs.Items()))
		if err != nil {
			s.logger.WithError(err).Warn(
				"Problem marshaling guilds message stats",
			)
		}
		w.Write(payload)
	})
	http.HandleFunc("/stats/commands", func(w http.ResponseWriter, r *http.Request) {
		payload, err := json.Marshal(s.buildStats(s.cmdsMsgs.Items()))
		if err != nil {
			s.logger.WithError(err).Warn(
				"Problem marshaling commands message stats",
			)
		}
		w.Write(payload)
	})

	http.ListenAndServe(":8080", nil)
}

func (s *statistics) buildStats(items map[string]cache.Item) []msgStats {
	out := []msgStats{}
	for k, v := range items {
		out = append(out, msgStats{
			Name:  k,
			Count: v.Object.(int),
		})
	}
	return out
}

func (s *statistics) handle(session *discordgo.Session, message *discordgo.MessageCreate) {
	guild, err := session.Guild(message.GuildID)
	if err != nil {
		s.logger.WithError(err).Warn(
			"Problem getting guild",
		)
	}

	channel, err := session.Channel(message.ChannelID)
	if err != nil {
		s.logger.WithError(err).Warn(
			"Problem getting channel",
		)
	}

	s.guildMsgs.Add(guild.Name, 0, cache.DefaultExpiration)
	s.chanMsgs.Add(channel.Name, 0, cache.DefaultExpiration)
	s.userMsgs.Add(message.Author.Username, 0, cache.DefaultExpiration)

	s.guildMsgs.Increment(guild.Name, 1)
	s.chanMsgs.Increment(channel.Name, 1)
	s.userMsgs.Increment(message.Author.Username, 1)
}

func (s *statistics) middleware(ctx *disgomux.Context) {
	s.cmdsMsgs.Add(ctx.Command, 0, cache.DefaultExpiration)
	s.cmdsMsgs.Increment(ctx.Command, 1)
}
