package main

import (
	"context"
	"log"
	"os"

	tgClient "pass-keeper-bot/clients/telegram"
	"pass-keeper-bot/consumer/event-consumer"
	"pass-keeper-bot/events/telegram"
	"pass-keeper-bot/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "storage.db"
	batchSize         = 100
)

func main() {
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token, ok := os.LookupEnv("BOT-TOKEN")

	if !ok {
		log.Panic("token is not specified")
	}

	return token
}
