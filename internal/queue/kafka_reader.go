package queue

import (
	"vk/pkg/model"
	"vk/pkg/proto"

	gproto "google.golang.org/protobuf/proto"
)

type KafkaQueueReader struct{}

func NewKafkaQueueReader() *KafkaQueueReader {
	return &KafkaQueueReader{}
}

func (q *KafkaQueueReader) ReadDoc(doc []byte) (*model.Document, error) {
	parsedDoc := &proto.TDocument{}

	err := gproto.Unmarshal(doc, parsedDoc)
	if err != nil {
		return nil, err
	}

	return &model.Document{
		Url:            parsedDoc.Url,
		PubDate:        parsedDoc.PubDate,
		Text:           parsedDoc.Text,
		FetchTime:      parsedDoc.FetchTime,
		FirstFetchTime: parsedDoc.FirstFetchTime,
	}, nil
}
