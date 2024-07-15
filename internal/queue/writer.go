package queue

import "vk/pkg/model"

type QueueWriter interface {
	WriteDoc(doc model.Document) error
}
