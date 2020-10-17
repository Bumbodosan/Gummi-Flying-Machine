package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Run(bot *Bot, args string, message *discordgo.Message) Sendable
}

type CommandGroup map[string]Command

func (c CommandGroup) Run(
	bot *Bot,
	args string,
	message *discordgo.Message,
) Sendable {
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

	embed := discordgo.MessageEmbed{
		Description: fmt.Sprintf(
			"Must be one of: %s\n",
			oneOf,
		),
	}

	if args == "" {
		embed.Title = "Missing subcommand"
	} else {
		embed.Title = fmt.Sprintf("Unknown subcommand: `%s`", commandName)
	}
	return ErrorMessage{Embed: &embed}
}
