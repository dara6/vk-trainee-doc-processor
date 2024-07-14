package repository

import (
	"sync"
	"vk/pkg/model"
)

type InMemoryRepository struct {
	data           map[string]*model.Document
	dataMutex      sync.RWMutex
	conditionMutex sync.Mutex
	condition      map[string]*sync.Mutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data:      make(map[string]*model.Document),
		condition: make(map[string]*sync.Mutex),
	}
}

func (repo *InMemoryRepository) GetDocument(url string) (*model.Document, error) {
	repo.dataMutex.RLock()
	defer repo.dataMutex.RUnlock()

	doc, exists := repo.data[url]
	if !exists {
		return nil, ErrDocumentNotFound
	}
	return doc, nil
}

func (repo *InMemoryRepository) SaveDocument(doc *model.Document) error {
	repo.dataMutex.Lock()
	defer repo.dataMutex.Unlock()

	repo.data[doc.Url] = doc
	return nil
}

func (repo *InMemoryRepository) LockDocument(url string) error {
	repo.conditionMutex.Lock()
	defer repo.conditionMutex.Unlock()

	mutex, exists := repo.condition[url]
	if !exists {
		mutex = &sync.Mutex{}
		repo.condition[url] = mutex
	}

	mutex.Lock()
	return nil
}

func (repo *InMemoryRepository) UnlockDocument(url string) error {
	repo.conditionMutex.Lock()
	defer repo.conditionMutex.Unlock()

	mutex, exists := repo.condition[url]
	if !exists {
		return ErrDocumentNotFound
	}

	mutex.Unlock()
	return nil
}
