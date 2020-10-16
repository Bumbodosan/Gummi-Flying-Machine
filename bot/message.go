package bot

import (
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Sendable interface {
	Send(*Bot, *discordgo.Message) error
}

type Message struct {
	Content         string                            `json:"content,omitempty"`
	Embed           *discordgo.MessageEmbed           `json:"embed,omitempty"`
	TTS             bool                              `json:"tts"`
	Files           []File                            `json:"-"`
	AllowedMentions *discordgo.MessageAllowedMentions `json:"allowed_mentions,omitempty"`

	Timeout time.Duration // zero for don't remove
	// maybe different timeout for message being replied to
}

type File struct {
	Name   string
	Reader io.Reader
}

type ErrorMessage struct {
	Content string

	// Add more things from discordgo.MessageSend here when needed
}

func (msg Message) Send(bot *Bot, replyingTo *discordgo.Message) error {
	var msgsSent []*discordgo.Message
	if msg.Embed != nil {
		ms, err := bot.Session.ChannelMessageSendEmbed(replyingTo.ChannelID, msg.Embed)
		if err != nil {
			return err
		}
		msgsSent = append(msgsSent, ms)
	}

	if msg.Files != nil {
		for _, file := range msg.Files {
			ms, err := bot.Session.ChannelMessageSendComplex(
				replyingTo.ChannelID,
				&discordgo.MessageSend{
					Files: []*discordgo.File{
						{
							Name:   file.Name,
							Reader: file.Reader,
						},
					},
				},
			)
			if err != nil {
				return err
			}
			msgsSent = append(msgsSent, ms)
		}
	}

	if msg.Timeout != 0 {
		time.Sleep(msg.Timeout)

		for _, msgSent := range msgsSent {
			err := bot.Session.ChannelMessagesBulkDelete(
				replyingTo.ChannelID,
				[]string{replyingTo.ID, msgSent.ID},
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (msg ErrorMessage) Send(bot *Bot, replyingTo *discordgo.Message) error {
	msgSent, err := bot.Session.ChannelMessageSend(replyingTo.ChannelID, msg.Content)

	if err != nil {
		return err
	}

	time.Sleep(time.Second * 20)

	err = bot.Session.ChannelMessagesBulkDelete(
		replyingTo.ChannelID,
		[]string{replyingTo.ID, msgSent.ID},
	)

	if err != nil {
		return err
	}

	return nil
}
