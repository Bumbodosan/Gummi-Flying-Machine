package bot

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type EncodeCommand struct {
	encode func(string) string
	name   string
}

func EncodeSubCommand() SubCommand {
	return SubCommand(map[string]Command{
		"base64": EncodeCommand{
			name: "base64",
			encode: func(text string) string {
				return base64.StdEncoding.EncodeToString([]byte(text))
			},
		},
		"hex": EncodeCommand{
			name: "hex",
			encode: func(text string) string {
				return hex.EncodeToString([]byte(text))
			},
		},
		"url-path": EncodeCommand{
			name: "url path",
			encode: func(text string) string {
				return url.PathEscape(text)
			},
		},
		"url-query": EncodeCommand{
			name: "url query",
			encode: func(text string) string {
				return url.QueryEscape(text)
			},
		},
		"utf8": EncodeCommand{
			name: "utf8",
			encode: func(text string) string {
				return text
			},
		},
	})
}

func (c EncodeCommand) Run(
	bot *Bot,
	args string,
	message *discordgo.Message,
) error {
	if args == "" {
		sentMessage, err := bot.Session.ChannelMessageSend(
			message.ChannelID,
			"Got nothing to encode",
		)
		if err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
		bot.Session.ChannelMessagesBulkDelete(message.ChannelID, []string{
			message.ID,
			sentMessage.ID,
		})
		return nil
	}

	encoded := c.encode(args)

	text := "```\n" +
		strings.ReplaceAll(args, "```", "`​``") +
		"```" +
		c.name +
		" encoded: ```\n" +
		encoded +
		"```"

	_, err := bot.Session.ChannelMessageSend(message.ChannelID, text)

	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	bot.Session.ChannelMessageDelete(
		message.ChannelID,
		message.ID,
	)

	return nil
}
