package bot

import (
	"fmt"
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
	Embed   *discordgo.MessageEmbed
	// Add more things from discordgo.MessageSend here when needed
}

func (msg Message) Send(bot *Bot, replyingTo *discordgo.Message) error {
	files := make([]*discordgo.File, len(msg.Files))
	for i, file := range msg.Files {
		files[i] = &discordgo.File{
			Name:   file.Name,
			Reader: file.Reader,
		}
	}

	// Only send plain context when we have no embed
	var content string
	if msg.Embed == nil {
		content = msg.Content
	}

	msgSent, err := bot.Session.ChannelMessageSendComplex(
		replyingTo.ChannelID,
		&discordgo.MessageSend{
			Content: content,
			Files:   files,
			Embed:   msg.Embed,
		},
	)
	if err != nil {
		return err
	}

	go func() {
		time.Sleep(5 * time.Second)
		err = bot.Session.ChannelMessageDelete(
			replyingTo.ChannelID,
			replyingTo.ID,
		)

		if err != nil {
			fmt.Println(err)
		}
	}()

	if msg.Timeout == 0 {
		msg.Timeout = time.Second * 10
	}
	if msg.Timeout > 0 {
		time.Sleep(msg.Timeout)
		err := bot.Session.ChannelMessagesBulkDelete(
			replyingTo.ChannelID,
			[]string{replyingTo.ID, msgSent.ID},
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (msg ErrorMessage) Send(bot *Bot, replyingTo *discordgo.Message) error {
	msg.Embed.Color = 15158332
	msgSent, err := bot.Session.ChannelMessageSendEmbed(replyingTo.ChannelID, msg.Embed)

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
