package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Run(bot *Bot, args string, message *discordgo.Message) error
}

type SubCommand map[string]Command

func (c SubCommand) Run(bot *Bot, args string, message *discordgo.Message) error {
	if args == "" {
		sentMessage, err := bot.Session.ChannelMessageSend(
			message.ChannelID,
			"Missing subcommand.",
		)
		if err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
		bot.Session.ChannelMessagesBulkDelete(
			message.ChannelID,
			[]string{
				message.ID,
				sentMessage.ID,
			},
		)
	}

	spaceIndex := strings.IndexRune(args, ' ')
	var commandName, subArgs string
	if spaceIndex == -1 {
		commandName = args
		subArgs = ""
	} else {
		commandName = args[:spaceIndex]
		subArgs = args[spaceIndex+1:]
	}

	if command := c[commandName]; command != nil {
		return command.Run(bot, subArgs, message)
	}

	var oneOf string
	for name := range c {
		oneOf += "\n- " + name
	}

	sentMessage, err := bot.Session.ChannelMessageSend(
		message.ChannelID,
		fmt.Sprintf(
			"Unknown subcommand '%s'. Must be one of: %s\n",
			commandName,
			oneOf,
		),
	)

	if err != nil {
		return err
	}

	time.Sleep(1 * time.Minute)

	bot.Session.ChannelMessagesBulkDelete(
		message.ChannelID,
		[]string{
			message.ID,
			sentMessage.ID,
		},
	)

	return nil
}
