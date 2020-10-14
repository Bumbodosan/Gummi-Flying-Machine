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
	if args == "" {
		return ErrorMessage{Content: "Missing subcommand"}
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

	return ErrorMessage{Content: fmt.Sprintf(
		"Unknown subcommand '%s'. Must be one of: %s\n",
		commandName,
		oneOf,
	)}
}
