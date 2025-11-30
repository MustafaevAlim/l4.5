package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"l4.5/internal/api"
	"l4.5/internal/api/controllers"
	"l4.5/internal/app"
	"l4.5/internal/config"
	"l4.5/internal/repository"
	"l4.5/scripts"
)

// Загрузка переменных окружения
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	go scripts.WriteInKafka()
	conf := config.New()
	db := repository.Init()
	cache := repository.NewLRUCache(conf.Cache.CacheSize, db)
	kafka := repository.NewKafkaReader([]string{conf.Kafka.KafkaBrokers}, conf.Kafka.KafkaTopic, conf.Kafka.KafkaConsumerGroup)
	ctrl := controllers.Controller{DB: db, Cache: cache}
	mux := api.RouteController(&ctrl)

	a := app.NewApp(db, cache, kafka, mux)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := a.Run(ctx)
	if err != nil {
		return
	}

}
