package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"net/http"
	"strings"

	"github.com/CS-5/disgomux"
	"github.com/bwmarrin/discordgo"
	"github.com/disintegration/imaging"
)

type cJPEG struct {
	Command  string
	HelpText string
}

var (
	imgSaturation float64 = 100
	imgBlur       float64 = 3
	imgQuality    int     = 1
)

func (i cJPEG) Init(m *disgomux.Mux) {
	// Nothing to init
}

func (i cJPEG) Handle(ctx *disgomux.Context) {
	var message *discordgo.Message

	if len(ctx.Arguments) == 0 {
		messages, err := ctx.Session.ChannelMessages(
			ctx.Message.ChannelID, 2, ctx.Message.ID, "", "",
		)
		if err != nil {
			cmdIssue(ctx, err, "There was a problem getting the lastest messages")
			return
		}

		messageLen := len(messages)
		for i := range messages {
			message = messages[messageLen-1-i]
		}
	} else {
		var err error
		message, err = ctx.Session.ChannelMessage(
			ctx.Message.ChannelID, ctx.Arguments[0],
		)
		if err != nil {
			ctx.ChannelSendf(
				"No message with ID `%s` found in this channel!",
				ctx.Arguments[0],
			)
			return
		}
	}

	if len(message.Attachments) == 0 || message.Attachments[0] == nil {
		ctx.ChannelSend("That message doesn't have any attachments!")
		return
	}

	attachment := message.Attachments[0]
	if strings.HasSuffix(attachment.ProxyURL, ".png") ||
		strings.HasSuffix(attachment.ProxyURL, ".jpg") ||
		strings.HasSuffix(attachment.ProxyURL, ".jpeg") {
		req, err := http.Get(attachment.ProxyURL)
		if err != nil {
			cmdIssue(ctx, err, "There was a problem getting the attachment")
			return
		}
		defer req.Body.Close()

		imgIn, _, err := image.Decode(req.Body)
		if err != nil {
			cmdIssue(ctx, err, "There was a problem decoding the image")
			return
		}

		/* Tweak these values to adjust JPEGness */
		img1 := imaging.AdjustSaturation(imgIn, imgSaturation)
		imgOut := imaging.Blur(img1, imgBlur)

		var buf bytes.Buffer // Buffer to return image
		err = jpeg.Encode(&buf, imgOut, &jpeg.Options{
			Quality: imgQuality,
		})
		if err != nil {
			cmdIssue(ctx, err, "There was a problem endoding the image")
			return
		}

		ctx.Session.ChannelFileSend(
			ctx.Message.ChannelID,
			attachment.Filename,
			&buf,
		)
		return
	}
	ctx.ChannelSend("No valid image to JPEGify (must be .jpg or .png)!")
}

func (i cJPEG) HandleHelp(ctx *disgomux.Context) bool {
	ctx.ChannelSend("`!jpeg` to JPEGify the image that was just sent.\n`!jpeg [message ID]` to JPEGify a specific image in this channel.")
	return true
}

func (i cJPEG) Settings() *disgomux.CommandSettings {
	return &disgomux.CommandSettings{
		Command:  i.Command,
		HelpText: i.HelpText,
	}
}

func (i cJPEG) Permissions() *disgomux.CommandPermissions {
	return &disgomux.CommandPermissions{}
}
