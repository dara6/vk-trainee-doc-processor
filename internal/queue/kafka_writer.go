package queue

import (
	"log"

	"vk/pkg/model"
	"vk/pkg/proto"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	gproto "google.golang.org/protobuf/proto"
)

type KafkaQueueWriter struct {
	topic    string
	producer *kafka.Producer
}

func NewKafkaQueueWriter(topic string, producer *kafka.Producer) *KafkaQueueWriter {
	return &KafkaQueueWriter{topic: topic, producer: producer}
}

func (q *KafkaQueueWriter) WriteDoc(doc model.Document) error {
	deliveryChan := make(chan kafka.Event)

	protoDoc := proto.TDocument{
		Url:            doc.Url,
		PubDate:        doc.PubDate,
		Text:           doc.Text,
		FetchTime:      doc.FetchTime,
		FirstFetchTime: doc.FirstFetchTime,
	}

	buf, _ := gproto.Marshal(&protoDoc)

	err := q.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &q.topic, Partition: kafka.PartitionAny},
		Value:          buf,
	}, deliveryChan)
	if err != nil {
		log.Fatalf("Failed to produce message: %v\n", err)
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Fatalf("Failed to deliver message: %v\n", m.TopicPartition.Error)
	} else {
		log.Printf("Produced message to %v[%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
	return nil
}
