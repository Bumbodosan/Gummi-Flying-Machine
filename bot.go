package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

var (
	Token string
	Prefix string
)

func init() {
	flag.StringVar(&Token, "t", "", "Discord Bot Token")
	flag.StringVar(&Prefix, "p", "!", "Prefix for messages")
	flag.Parse()
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
	if message.Author.ID == session.State.User.ID {return}
	if !strings.HasPrefix(message.Content, Prefix) {return}

	s := strings.ToLower(string(message.Content[len(Prefix):]))

	if s == "ping" {
		session.ChannelMessageSend(message.ChannelID, "Pong!")
	} else if s == "pong" {
		session.ChannelMessageSend(message.ChannelID, "Ping!")
	} else {
		fmt.Println(s)
	}
	
}