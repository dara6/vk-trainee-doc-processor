package queue

import "vk/pkg/model"

type QueueReader interface {
	ReadDoc(doc []byte) (*model.Document, error)
}
