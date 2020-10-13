package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bumbodosan/Gummi-Flying-Machine/bot"

	"github.com/joho/godotenv"
)

func main() {
	b := &bot.Bot{}

	godotenv.Load();

	flag.StringVar(&b.Token, "token", "", "Discord Bot Token")
	flag.StringVar(&b.Prefix, "prefix", "", "Command prefix")
	flag.Parse()
	if b.Token == "" {
		var found bool
		b.Token, found = os.LookupEnv("TOKEN")
		if !found {
			fmt.Println("Missing token")
			os.Exit(-1)
		}
	}
	if b.Prefix == "" {
		var found bool
		b.Prefix, found = os.LookupEnv("PREFIX")
		if !found {
			fmt.Println("Missing prefix")
			os.Exit(-1)
		}
	}

	if err := b.Start(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println("Gummi Flying Machine running")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	fmt.Println(<-ch, "received, quitting...")

	if err := b.Stop(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
