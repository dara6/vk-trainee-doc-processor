package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"vk/internal/config"
	"vk/internal/queue"
	"vk/pkg/repository"
	processor "vk/pkg/service"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// configure
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", cfg.KafkaBrokerHost, cfg.KafkaBrokerPort),
	})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	qw := queue.NewKafkaQueueWriter(cfg.KafkaOutTopic, producer)

	// consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", cfg.KafkaBrokerHost, cfg.KafkaBrokerPort),
		"group.id":          "consumer-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	consumer.Subscribe(cfg.KafkaInTopic, nil)

	qr := queue.NewKafkaQueueReader()

	// repo
	dsn := "user=" + cfg.PostgresUser + " password=" + cfg.PostgresPassword +
		" dbname=" + cfg.PostgresDB + " sslmode=disable" +
		" host=" + cfg.PostgresHost + " port=" + cfg.PostgresPort

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	repo := repository.NewPostgresRepository(db)

	// processor
	p := processor.NewProcessor(repo)

	// logic
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true

	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			msg, err := consumer.ReadMessage(-1)
			if err == nil {
				_ = processMessage(msg.Value, qr, qw, p)
			} else {
				log.Fatalf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}
}

func processMessage(msg []byte, qr *queue.KafkaQueueReader, qw *queue.KafkaQueueWriter, p processor.Processor) error {
	doc, err := qr.ReadDoc(msg)
	if err != nil {
		log.Fatalf("Can't read doc: %v: %s", err, msg)
		return err
	}

	newDoc, err := p.Process(doc)
	if err != nil {
		log.Fatalf("Can't process doc: %v: %s", err, msg)
		return err
	}

	err = qw.WriteDoc(*newDoc)
	if err != nil {
		log.Fatalf("Can't write doc: %v: %s", err, msg)
		return err
	}

	return nil
}
