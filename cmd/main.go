package main

import (
	"musicCloud-bot/bot"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go handleSignal()
	bot.Start()
}

func handleSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-c
	bot.B.Stop()
	os.Exit(0)
}
