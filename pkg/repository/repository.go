package repository

import (
	"vk/pkg/model"

	"errors"
)

var ErrDocumentNotFound = errors.New("document not found")

type Repository interface {
	GetDocument(url string) (*model.Document, error)
	SaveDocument(doc *model.Document) error
	LockDocument(url string) error
	UnlockDocument(url string) error
}
