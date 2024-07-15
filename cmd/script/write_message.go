package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"vk/internal/config"
	"vk/internal/queue"
	"vk/pkg/model"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

func main() {
	doc, err := ParseDocumentFromFlags()
	if err != nil {
		log.Fatalf("Error parse doc: %v", err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", cfg.KafkaBrokerHost, cfg.KafkaBrokerPort),
	})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	q := queue.NewKafkaQueueWriter(cfg.KafkaInTopic, producer)
	q.WriteDoc(*doc)
}

func ParseDocumentFromFlags() (*model.Document, error) {
	urlFlag := flag.String("url", "", "Document URL")
	pubDateFlag := flag.String("pubDate", "", "Publication date")
	fetchTimeFlag := flag.String("fetchTime", "", "Fetch time")
	textFlag := flag.String("text", "", "Text content")
	firstFetchTimeFlag := flag.String("firstFetchTime", "", "First fetch time")

	flag.Parse()

	pubDateInt, err := strconv.ParseUint(*pubDateFlag, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid PubDate value: %v", err)
	}

	fetchTimeInt, err := strconv.ParseUint(*fetchTimeFlag, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid FetchTime value: %v", err)
	}

	firstFetchTimeInt, err := strconv.ParseUint(*firstFetchTimeFlag, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid FirstFetchTime value: %v", err)
	}

	doc := &model.Document{
		Url:            *urlFlag,
		PubDate:        pubDateInt,
		FetchTime:      fetchTimeInt,
		Text:           *textFlag,
		FirstFetchTime: firstFetchTimeInt,
	}

	return doc, nil
}
