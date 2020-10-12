package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token            string
	Prefix           string
	PrefixNeedsSpace bool
)

func init() {
	flag.StringVar(&Token, "t", "", "Discord Bot Token")
	flag.StringVar(&Prefix, "p", "gfm", "Prefix for messages")
	flag.Parse()
	if regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(Prefix) {
		Prefix += " "
	}
}

func main() {
	discord, err := discordgo.New("Bot " + Token)

	if err != nil {
		fmt.Println("Error starting bot:", err)
		return
	}

	discord.AddHandler(messageCreated)

	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = discord.Open()
	if err != nil {
		fmt.Println("Error starting bot:", err)
		return
	}

	fmt.Println("Bot is running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func messageCreated(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}
	if !strings.HasPrefix(message.Content, Prefix) {
		return
	}
	handleMessage(session, message.Message)
}

func handleMessage(session *discordgo.Session, message *discordgo.Message) {
	s := strings.ToLower(string(message.Content[len(Prefix):]))

	var returningMessage string

	if s == "ping" {
		returningMessage = "Pong!"
	} else if s == "pong" {
		returningMessage = "Ping!"
	} else {
		return
	}

	sentMessage, err := session.ChannelMessageSend(message.ChannelID, returningMessage)
	if err != nil {
		return
	}

	go deleteMessage(3, session, message)
	go deleteMessage(5, session, sentMessage)
}

func deleteMessage(seconds int, session *discordgo.Session, message *discordgo.Message) {
	<-time.NewTimer(time.Duration(seconds) * time.Second).C

	err := session.ChannelMessageDelete(message.ChannelID, message.ID)
	if err != nil {
		fmt.Println(err)
	}
}
