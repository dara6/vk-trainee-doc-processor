package repository_test

import (
	"sync"
	"testing"
	"vk/pkg/model"
	"vk/pkg/repository"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepository(t *testing.T) {
	repo := repository.NewInMemoryRepository()

	doc := &model.Document{
		Url:            "http://example.com",
		PubDate:        123456789,
		FetchTime:      12346789,
		Text:           "12346789",
		FirstFetchTime: 12346789,
	}

	t.Run("SaveDocument", func(t *testing.T) {
		err := repo.SaveDocument(doc)
		assert.NoError(t, err, "expected no error saving document")

		savedDoc, err := repo.GetDocument(doc.Url)
		assert.NoError(t, err, "expected no error getting document")
		assert.Equal(t, doc, savedDoc, "expected to get the saved document")
	})

	t.Run("GetDocument_NotFound", func(t *testing.T) {
		_, err := repo.GetDocument("http://notfound.com")
		assert.Error(t, err, "expected error for not found document")
		assert.Equal(t, "document not found", err.Error(), "expected 'document not found' error")
	})

	t.Run("LockDocument", func(t *testing.T) {
		err := repo.LockDocument(doc.Url)
		assert.NoError(t, err, "expected no error locking document")

		// emulate concurrent lock
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			err := repo.LockDocument(doc.Url)
			assert.NoError(t, err, "expected no error locking document again")
			repo.UnlockDocument(doc.Url)
		}()

		repo.UnlockDocument(doc.Url)
		wg.Wait()
	})

	t.Run("UnlockDocument_NotFound", func(t *testing.T) {
		err := repo.UnlockDocument("http://notfound.com")
		assert.Error(t, err, "expected error unlocking not found document")
		assert.Equal(t, "document not found", err.Error(), "expected 'document not found' error")
	})
}
