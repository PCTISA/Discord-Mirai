package command

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
)

// Init is called by the multiplexer before the bot starts to initialize any
// variables the command needs.
func (c JPEG) Init(m *disgomux.Mux) {
	// Nothing to init
}

// Handle is called by the multiplexer whenever a user triggers the command.
func (c JPEG) Handle(ctx *disgomux.Context) {
	var message *discordgo.Message

	if len(ctx.Arguments) == 0 {
		messages, err := ctx.Session.ChannelMessages(
			ctx.Message.ChannelID, 2, ctx.Message.ID, "", "",
		)
		if err != nil {
			commandLogs.CmdErr(ctx, err, "There was a problem getting the lastest messages")
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
			commandLogs.CmdErr(ctx, err, "There was a problem getting the attachment")
			return
		}
		defer req.Body.Close()

		imgIn, _, err := image.Decode(req.Body)
		if err != nil {
			commandLogs.CmdErr(ctx, err, "There was a problem decoding the image")
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
			commandLogs.CmdErr(ctx, err, "There was a problem endoding the image")
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

// HandleHelp is called by whatever help command is in place when a user enters
// "!help [command name]". If the help command is not being handled, return
// false.
func (c JPEG) HandleHelp(ctx *disgomux.Context) bool {
	ctx.ChannelSend("`!jpeg` to JPEGify the image that was just sent.\n`!jpeg [message ID]` to JPEGify a specific image in this channel.")
	return true
}

// Settings is called by the multiplexer on startup to process any settings
// associated with that command.
func (c JPEG) Settings() *disgomux.CommandSettings {
	return &disgomux.CommandSettings{
		Command:  c.Command,
		HelpText: c.HelpText,
	}
}

// Permissions is called by the multiplexer on startup to collect the list of
// permissions required to run the given command.
func (c JPEG) Permissions() *disgomux.CommandPermissions {
	return &disgomux.CommandPermissions{}
}
