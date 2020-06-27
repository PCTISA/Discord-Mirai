package command

import (
	"bytes"
	"image"
	"image/jpeg"
	"net/http"
	"regexp"
	"strings"

	"github.com/PCTISA/Discord-Mirai/multiplexer"
	"github.com/bwmarrin/discordgo"
	"github.com/disintegration/imaging"
)

// JPEG is a command
// TODO: Make this a better description
type JPEG struct {
	Command  string
	HelpText string
}

var (
	imgSaturation float64 = 100
	imgBlur       float64 = 3
	imgQuality    int     = 1

	reURL = regexp.MustCompile(`(http(s?):)([/|.|\w|\s|-])*\.(?:jpg|jpeg|png)`)
)

// Init is called by the multiplexer before the bot starts to initialize any
// variables the command needs.
func (c JPEG) Init(m *multiplexer.Mux) {
	// Nothing to init
}

// Handle is called by the multiplexer whenever a user triggers the command.
func (c JPEG) Handle(ctx *multiplexer.Context) {
	var message *discordgo.Message

	if len(ctx.Arguments) == 0 {
		messages, err := ctx.Session.ChannelMessages(
			ctx.Message.ChannelID, 2, ctx.Message.ID, "", "",
		)
		if err != nil {
			cmdErr(ctx, err, "There was a problem getting the lastest messages")
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

	var urls []string
	if len(message.Attachments) == 0 || message.Attachments[0] == nil {
		urls = reURL.FindAllString(message.Content, -1)
	} else {
		for _, attach := range message.Attachments {
			if strings.HasSuffix(attach.ProxyURL, ".png") ||
				strings.HasSuffix(attach.ProxyURL, ".jpg") ||
				strings.HasSuffix(attach.ProxyURL, ".jpeg") {
				urls = append(urls, attach.ProxyURL)
			}
		}
	}

	if len(urls) == 0 {
		ctx.ChannelSend("No valid image to JPEGify (must be .jpg or .png)!")
		return
	}

	ctx.Session.ChannelTyping(message.ChannelID)
	for _, url := range urls {
		req, err := http.Get(url)
		if err != nil {
			cmdErr(ctx, err, "There was a problem getting the attachment")
			return
		}
		defer req.Body.Close()

		imgIn, _, err := image.Decode(req.Body)
		if err != nil {
			cmdErr(ctx, err, "There was a problem decoding the image")
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
			cmdErr(ctx, err, "There was a problem endoding the image")
			return
		}

		ctx.Session.ChannelFileSend(
			ctx.Message.ChannelID,
			"compressed.jpeg",
			&buf,
		)
	}
}

// HandleHelp is called by whatever help command is in place when a user enters
// "!help [command name]". If the help command is not being handled, return
// false.
func (c JPEG) HandleHelp(ctx *multiplexer.Context) bool {
	ctx.ChannelSend(
		"`!jpeg` to JPEGify the image that was just sent.\n" +
			"`!jpeg [message ID]` to JPEGify a specific image in this channel.",
	)
	return true
}

// Settings is called by the multiplexer on startup to process any settings
// associated with that command.
func (c JPEG) Settings() *multiplexer.CommandSettings {
	return &multiplexer.CommandSettings{
		Command:  c.Command,
		HelpText: c.HelpText,
	}
}
