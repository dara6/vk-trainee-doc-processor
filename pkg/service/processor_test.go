package processor

import (
	"errors"
	"testing"
	"time"

	"vk/pkg/model"
	"vk/pkg/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository acts as a mock implementation of the Repository interface.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetDocument(url string) (*model.Document, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Document), args.Error(1)
}

func (m *MockRepository) SaveDocument(doc *model.Document) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *MockRepository) LockDocument(url string) error {
	args := m.Called(url)
	return args.Error(0)
}

func (m *MockRepository) UnlockDocument(url string) error {
	args := m.Called(url)
	return args.Error(0)
}

func TestProcessor(t *testing.T) {
	mockRepo := new(MockRepository)
	processor := NewProcessor(mockRepo)

	doc := model.Document{
		Url:       "http://example.com",
		PubDate:   uint64(time.Now().Add(-time.Hour * 1).Unix()),
		FetchTime: uint64(time.Now().Unix()),
		Text:      "new content",
	}

	t.Run("Process_NewDocument", func(t *testing.T) {
		updatedDoc := &model.Document{
			Url:            doc.Url,
			PubDate:        doc.PubDate,
			FetchTime:      doc.FetchTime,
			Text:           doc.Text,
			FirstFetchTime: doc.FetchTime,
		}

		newDoc := doc

		mockRepo.On("LockDocument", doc.Url).Return(nil)
		mockRepo.On("UnlockDocument", doc.Url).Return(nil)
		mockRepo.On("GetDocument", doc.Url).Return(nil, repository.ErrDocumentNotFound)
		mockRepo.On("SaveDocument", updatedDoc).Return(nil)

		result, err := processor.Process(&newDoc)
		assert.NoError(t, err, "expected no error processing new document")
		assert.Equal(t, updatedDoc, result, "expected the result to be the same as the new document")

		mockRepo.AssertExpectations(t)
	})

	mockRepo = new(MockRepository)
	processor = NewProcessor(mockRepo)

	t.Run("Process_UpdateNewDocument", func(t *testing.T) {
		existingDoc := &model.Document{
			Url:            doc.Url,
			PubDate:        uint64(time.Now().Add(-time.Hour * 2).Unix()),
			FetchTime:      uint64(time.Now().Add(-time.Hour * 1).Unix()),
			Text:           "old content",
			FirstFetchTime: uint64(time.Now().Add(-time.Hour * 1).Unix()),
		}

		newDoc := &model.Document{
			Url:       doc.Url,
			PubDate:   uint64(time.Now().Add(-time.Hour * 2).Unix()),
			FetchTime: uint64(time.Now().Unix()),
			Text:      "new content",
		}

		updatedDoc := &model.Document{
			Url:            doc.Url,
			PubDate:        existingDoc.PubDate,
			FetchTime:      newDoc.FetchTime,
			Text:           newDoc.Text,
			FirstFetchTime: existingDoc.FetchTime,
		}

		mockRepo.On("LockDocument", doc.Url).Return(nil)
		mockRepo.On("UnlockDocument", doc.Url).Return(nil)
		mockRepo.On("GetDocument", doc.Url).Return(existingDoc, nil)
		mockRepo.On("SaveDocument", updatedDoc).Return(nil)

		result, err := processor.Process(newDoc)
		assert.NoError(t, err, "expected no error updating document")
		assert.Equal(t, updatedDoc, result, "expected the result to be the updated document")

		mockRepo.AssertExpectations(t)
	})

	mockRepo = new(MockRepository)
	processor = NewProcessor(mockRepo)

	t.Run("Process_UpdateOldDocument", func(t *testing.T) {
		existingDoc := &model.Document{
			Url:            doc.Url,
			PubDate:        uint64(time.Now().Unix()),
			FetchTime:      uint64(time.Now().Unix()),
			Text:           "old content",
			FirstFetchTime: uint64(time.Now().Unix()),
		}

		newDoc := &model.Document{
			Url:       doc.Url,
			PubDate:   uint64(time.Now().Add(-time.Hour * 2).Unix()),
			FetchTime: uint64(time.Now().Add(-time.Hour * 1).Unix()),
			Text:      "new content",
		}

		updatedDoc := &model.Document{
			Url:            doc.Url,
			PubDate:        newDoc.PubDate,
			FetchTime:      existingDoc.FetchTime,
			Text:           existingDoc.Text,
			FirstFetchTime: newDoc.FetchTime,
		}

		mockRepo.On("LockDocument", doc.Url).Return(nil)
		mockRepo.On("UnlockDocument", doc.Url).Return(nil)
		mockRepo.On("GetDocument", doc.Url).Return(existingDoc, nil)
		mockRepo.On("SaveDocument", updatedDoc).Return(nil)

		result, err := processor.Process(newDoc)
		assert.NoError(t, err, "expected no error updating document")
		assert.Equal(t, updatedDoc, result, "expected the result to be the updated document")

		mockRepo.AssertExpectations(t)
	})

	mockRepo = new(MockRepository)
	processor = NewProcessor(mockRepo)

	t.Run("Process_UpdateCopycatDocument", func(t *testing.T) {
		existingDoc := &model.Document{
			Url:            doc.Url,
			PubDate:        uint64(time.Now().Unix()),
			FetchTime:      uint64(time.Now().Unix()),
			Text:           "old content",
			FirstFetchTime: uint64(time.Now().Unix()),
		}

		newDoc := existingDoc

		updatedDoc := existingDoc

		mockRepo.On("LockDocument", doc.Url).Return(nil)
		mockRepo.On("UnlockDocument", doc.Url).Return(nil)
		mockRepo.On("GetDocument", doc.Url).Return(existingDoc, nil)
		mockRepo.On("SaveDocument", updatedDoc).Return(nil)

		result, err := processor.Process(newDoc)
		assert.NoError(t, err, "expected no error updating document")
		assert.Equal(t, updatedDoc, result, "expected the result to be the updated document")

		mockRepo.AssertExpectations(t)
	})

	mockRepo = new(MockRepository)
	processor = NewProcessor(mockRepo)

	t.Run("Process_LockFailure", func(t *testing.T) {
		newDoc := doc

		mockRepo.On("LockDocument", doc.Url).Return(errors.New("lock error"))

		result, err := processor.Process(&newDoc)
		assert.Error(t, err)
		assert.Nil(t, result, "expected nil result on lock failure")
		assert.Equal(t, "lock error", err.Error(), "expected lock error")

		mockRepo.AssertExpectations(t)
	})

	mockRepo = new(MockRepository)
	processor = NewProcessor(mockRepo)

	t.Run("Process_SaveFailure", func(t *testing.T) {
		updatedDoc := &model.Document{
			Url:            doc.Url,
			PubDate:        doc.PubDate,
			FetchTime:      doc.FetchTime,
			Text:           doc.Text,
			FirstFetchTime: doc.FetchTime,
		}

		newDoc := doc

		mockRepo.On("LockDocument", doc.Url).Return(nil)
		mockRepo.On("UnlockDocument", doc.Url).Return(nil)
		mockRepo.On("GetDocument", doc.Url).Return(nil, repository.ErrDocumentNotFound)
		mockRepo.On("SaveDocument", updatedDoc).Return(errors.New("save error"))

		result, err := processor.Process(&newDoc)
		assert.Error(t, err)
		assert.Nil(t, result, "expected nil result on save failure")
		assert.Equal(t, "save error", err.Error(), "expected save error")

		mockRepo.AssertExpectations(t)
	})

	mockRepo = new(MockRepository)
	processor = NewProcessor(mockRepo)

	t.Run("Process_GetFailure", func(t *testing.T) {
		newDoc := doc

		mockRepo.On("LockDocument", doc.Url).Return(nil)
		mockRepo.On("UnlockDocument", doc.Url).Return(nil)
		mockRepo.On("GetDocument", doc.Url).Return(nil, errors.New("get error"))

		result, err := processor.Process(&newDoc)
		assert.Error(t, err)
		assert.Nil(t, result, "expected nil result on get failure")
		assert.Equal(t, "get error", err.Error(), "expected get error")

		mockRepo.AssertExpectations(t)
	})
}
