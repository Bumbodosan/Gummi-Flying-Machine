package bot

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type EncodeCommand struct {
	encode func(string) string
	name   string
}

type DecodeCommand struct {
	decode func(string) ([]byte, error)
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
		"url-path-segment": EncodeCommand{
			name: "url path segment",
			encode: func(text string) string {
				return url.PathEscape(text)
			},
		},
		"url-query-component": EncodeCommand{
			name: "url query component",
			encode: func(text string) string {
				return url.QueryEscape(text)
			},
		},
	})
}

func DecodeCommands() CommandGroup {
	return map[string]Command{
		"base64": DecodeCommand{
			name: "base64",
			decode: func(text string) ([]byte, error) {
				return base64.StdEncoding.DecodeString(text)
			},
		},
		"hex": DecodeCommand{
			name: "hex",
			decode: func(text string) ([]byte, error) {
				return hex.DecodeString(text)
			},
		},
		"url-query": DecodeCommand{
			name: "url query",
			decode: func(text string) ([]byte, error) {
				values, err := url.ParseQuery(text)
				if err != nil {
					return nil, err
				}
				var output []byte

				for key, value := range values {
					if len(value) == 0 {
						output = append(output, []byte(key)...)
					} else if len(value) == 1 {
						output = append(output, []byte(key+": "+value[0])...)
					} else {
						for _, val := range value {
							output = append(output, []byte(key+": "+val)...)
						}
					}
				}

				return output, nil
			},
		},
	}
}

// Run : Runs the command,
func (c EncodeCommand) Run(
	bot *Bot,
	args string,
	message *discordgo.Message,
) Sendable {
	if args == "" {
		return ErrorMessage{Title: "Got nothing to encode"}
	}

	encoded := c.encode(args)

	return Message{
		Content: encoded,
		Embed: &discordgo.MessageEmbed{
			Title: "Encoded to `" + c.name + "`",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Input",
					Value: "`" + args + "`",
				},
				{
					Name:  "Output",
					Value: "`" + encoded + "`",
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				IconURL: message.Author.AvatarURL("1024"),
				Text:    "Requested by " + message.Author.Username,
			},
		},
	}
}

func (c DecodeCommand) Run(
	bot *Bot,
	args string,
	message *discordgo.Message,
) Sendable {
	if args == "" {
		return ErrorMessage{Title: "Got nothing to decode"}
	}

	decoded, err := c.decode(args)

	if err != nil {
		return ErrorMessage{Title: "Decoding error", Description: err.Error()}
	}

	decodedString := string(decoded)

	isPrintable := true
	for _, r := range decodedString {
		if r < 0x20 && r != '\n' && r != '\t' && r != '\r' || r == 0x7f {
			isPrintable = false
			break
		}
	}

	msg := Message{
		Content: decodedString,
		Timeout: -1,
	}

	if isPrintable {
		msg.Embed = &discordgo.MessageEmbed{
			Title: "Decoded from `" + c.name + "`",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Input",
					Value: "```" + args + "```",
				},
				{
					Name:  "Output",
					Value: "```" + decodedString + "```",
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				IconURL: message.Author.AvatarURL("1024"),
				Text:    "Requested by " + message.Author.Username,
			},
		}
	} else {
		msg.Files = []File{
			{
				Name:   "decoded",
				Reader: strings.NewReader(decodedString),
			},
		}
	}
	return msg
}
