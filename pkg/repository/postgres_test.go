package repository_test

import (
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"vk/internal/config"
	"vk/pkg/model"
	"vk/pkg/repository"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var repo *repository.PostgresRepository

func TestMain(m *testing.M) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dsn := "user=" + cfg.PostgresUser + " password=" + cfg.PostgresPassword +
		" dbname=" + cfg.PostgresDB + " sslmode=disable" +
		" host=" + cfg.PostgresHost + " port=" + cfg.PostgresPort

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec("TRUNCATE TABLE documents")

	repo = repository.NewPostgresRepository(db)

	code := m.Run()

	db.MustExec("TRUNCATE TABLE documents")
	os.Exit(code)
}

func TestPostgresRepository(t *testing.T) {
	doc := &model.Document{
		Url:            "http://example.com",
		PubDate:        uint64(time.Now().Unix()),
		FetchTime:      uint64(time.Now().Unix()),
		Text:           "example content",
		FirstFetchTime: uint64(time.Now().Unix()),
	}

	t.Run("SaveDocument", func(t *testing.T) {
		err := repo.SaveDocument(doc)
		assert.NoError(t, err, "expected no error saving document")

		savedDoc, err := repo.GetDocument(doc.Url)
		assert.NoError(t, err, "expected no error getting document")
		assert.Equal(t, doc, savedDoc)
	})

	t.Run("GetDocument_NotFound", func(t *testing.T) {
		_, err := repo.GetDocument("http://notfound.com")
		assert.Error(t, err, "expected error for not found document")
		assert.Equal(t, repository.ErrDocumentNotFound, err, "expected 'document not found' error")
	})

	t.Run("LockDocument", func(t *testing.T) {
		gorutinesCount := 3
		sleepTime := time.Second * 1
		expectedWorkDuration := sleepTime * time.Duration(gorutinesCount)

		wg := sync.WaitGroup{}
		wg.Add(gorutinesCount)

		startTime := time.Now()

		idx := 0
		for idx < gorutinesCount {
			go func() {
				defer wg.Done()
				err := repo.LockDocument(doc.Url)
				assert.NoError(t, err, "expected no error locking document again")
				<-time.After(sleepTime)
				repo.UnlockDocument(doc.Url)
			}()
			idx += 1
		}

		wg.Wait()

		endTime := time.Now()
		workDuration := endTime.Sub(startTime)

		assert.GreaterOrEqual(
			t, workDuration, expectedWorkDuration,
			"expected total work duration to be at least %v, got %v",
			expectedWorkDuration, workDuration,
		)
	})

	t.Run("UnlockDocument_NotFound", func(t *testing.T) {
		err := repo.UnlockDocument("http://notfound.com")
		// Since we use PostgreSQL advisory locks, unlocking a non-existing document does not raise an error.
		assert.NoError(t, err, "expected no error unlocking not found document")
	})
}
