package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Bot struct
type Bot struct {
	Session *discordgo.Session
	// DB *gorm.DB

	Prefix string
	Token  string

	Commands map[string]Command
}

// Start the bot
func (bot *Bot) Start() error {
	if err := bot.initDiscord(); err != nil {
		return err
	}

	if err := bot.initCommands(); err != nil {
		return err
	}

	return nil
}

// Stop the bot
func (bot *Bot) Stop() error {
	return bot.Session.Close()
}

func (bot *Bot) initDiscord() error {
	var err error
	bot.Session, err = discordgo.New("Bot " + bot.Token)
	if err != nil {
		return err
	}

	bot.Session.AddHandler(bot.onMessage)
	bot.Session.Identify.Intents = discordgo.MakeIntent(
		discordgo.IntentsGuildMessages,
	)

	if err := bot.Session.Open(); err != nil {
		return err
	}

	if err := bot.Session.UpdateStatus(0, "in the skies!"); err != nil {
		return err
	}

	return nil
}

func (bot *Bot) onMessage(
	_ *discordgo.Session,
	message *discordgo.MessageCreate,
) {
	if message.Author.ID == bot.Session.State.User.ID ||
		!strings.HasPrefix(message.Content, bot.Prefix) {
		return
	}

	rest := strings.TrimSpace(message.Content[len(bot.Prefix):])

	spaceIndex := strings.IndexRune(rest, ' ')
	var commandName, args string
	if spaceIndex == -1 {
		commandName = rest
		args = ""
	} else {
		commandName = rest[:spaceIndex]
		args = rest[spaceIndex+1:]
	}

	command := bot.Commands[commandName]
	var reply Sendable
	if command == nil {
		reply = ErrorMessage{Title: "Unknown command"}
	} else {
		reply = command.Run(bot, args, message.Message)
	}

	err := reply.Send(bot, message.Message)

	if err != nil {
		fmt.Println(err)
	}
}

func (bot *Bot) initCommands() error {
	bot.Commands = map[string]Command{
		"ping":   PingCommand{},
		"encode": EncodeCommands(),
		"decode": DecodeCommands(),
	}

	return nil
}
