package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bumbodosan/Gummi-Flying-Machine/bot"

	"github.com/joho/godotenv"
)

func main() {
	b := &bot.Bot{}

	godotenv.Load()

	if os.Getenv("TOKEN") == "" {
		fmt.Println("Missing token.")
		os.Exit(-1)
	}
	b.Token = os.Getenv("TOKEN")
	b.Prefix = os.Getenv("PREFIX")
	if b.Prefix == "" {
		b.Prefix = "!"
	}

	if err := b.Start(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println("Gummi Flying Machine is in the skies! 🛩️")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	fmt.Println(<-ch, "Exiting...")

	if err := b.Stop(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
