package bot

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"

	"github.com/bwmarrin/discordgo"
)

type EncodeCommand struct {
	encode func(string) string
	name   string
}

func EncodeCommands() CommandGroup {
	return CommandGroup(map[string]Command{
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
	})
}

func (c EncodeCommand) Run(
	bot *Bot,
	args string,
	message *discordgo.Message,
) Sendable {
	if args == "" {
		return ErrorMessage{Content: "Got nothing to encode"}
	}

	encoded := c.encode(args)

	text := "```\n" + encoded + "```"

	return Message{Content: text}
}
