package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Run(bot *Bot, args string, message *discordgo.Message) error
}

type SubCommand map[string]Command

func (c SubCommand) Run(bot *Bot, args string, message *discordgo.Message) error {
	spaceIndex := strings.IndexRune(args, ' ')
	var commandName, subArgs string
	if spaceIndex == -1 {
		commandName = subArgs
		args = ""
	} else {
		commandName = args[:spaceIndex]
		args = args[spaceIndex+1:]
	}

	if command := c[commandName]; command != nil {
		command.Run(bot, args, message)
	}

	return nil
}
