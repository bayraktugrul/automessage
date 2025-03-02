package integration

import (
	"automsg/pkg/persistence"
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbName     = "testdb"
	dbUser     = "testuser"
	dbPassword = "testpassword"
)

func setupPostgresContainer(t *testing.T) (*postgres.PostgresContainer, *sql.DB) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)
	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)

	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port.Port(), dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	err = initializeSchema(db)
	assert.NoError(t, err)

	return pgContainer, db
}

func initializeSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE recipients (
		id SERIAL PRIMARY KEY,
		phone_number VARCHAR(15) NOT NULL UNIQUE,
		name VARCHAR(100),
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE messages (
		id SERIAL PRIMARY KEY,
		content TEXT NOT NULL CHECK (LENGTH(content) <= 160),
		recipient_id INTEGER NOT NULL REFERENCES recipients(id),
		is_sent BOOLEAN DEFAULT FALSE,
		sent_at TIMESTAMP,
		message_id VARCHAR(50),
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	CREATE INDEX idx_messages_is_sent ON messages(is_sent);
	CREATE INDEX idx_messages_created_at ON messages(created_at);

	INSERT INTO recipients (phone_number, name, created_at, updated_at)
	VALUES 
		('+905551234567', 'Test Recipient 1', NOW(), NOW()),
		('+905551234568', 'Test Recipient 2', NOW(), NOW());

	INSERT INTO messages (content, recipient_id, is_sent, created_at, updated_at)
	VALUES 
		('test message 1', 1, false, NOW(), NOW()),
		('test message 2', 1, false, NOW(), NOW()),
		('test message 3', 2, true, NOW(), NOW()),
		('test message 4', 2, false, NOW(), NOW());
	`

	_, err := db.Exec(schema)
	return err
}

func Test_PostgresMessageRepository_should_get_unsent_processing_messages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	//given
	ctx := context.Background()
	pgContainer, db := setupPostgresContainer(t)
	defer func() {
		db.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	repo := persistence.NewPostgresMessageRepository(db)

	//when
	messages, err := repo.GetUnsentProcessingMessages(ctx, 10)

	//then
	assert.NoError(t, err)
	assert.Len(t, messages, 3, "Expected 3 unsent messages")

	assert.Equal(t, messages[0].Content, "test message 1")
	assert.Equal(t, messages[0].PhoneNumber, "+905551234567")
	assert.Equal(t, messages[1].Content, "test message 2")
	assert.Equal(t, messages[1].PhoneNumber, "+905551234567")
	assert.Equal(t, messages[2].Content, "test message 4")
	assert.Equal(t, messages[2].PhoneNumber, "+905551234568")
}

func Test_PostgresMessageRepository_should_get_sent_messages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	//given
	ctx := context.Background()
	pgContainer, db := setupPostgresContainer(t)
	defer func() {
		db.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	repo := persistence.NewPostgresMessageRepository(db)

	//when
	messages, totalCount, err := repo.GetSentMessages(ctx, 1, 10)

	//then
	assert.NoError(t, err)
	assert.Equal(t, 1, totalCount)
	assert.Len(t, messages, 1)
}

func Test_PostgresMessageRepository_should_mark_message_as_sent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	//given
	ctx := context.Background()
	pgContainer, db := setupPostgresContainer(t)
	defer func() {
		db.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	repo := persistence.NewPostgresMessageRepository(db)

	//when
	messages, err := repo.GetUnsentProcessingMessages(ctx, 1)

	assert.NoError(t, err)
	messageID := messages[0].Id

	tx, err := repo.BeginTx(ctx)
	assert.NoError(t, err)

	locked, err := repo.LockMessageForProcessing(ctx, tx, messageID)
	assert.NoError(t, err)
	assert.True(t, locked)

	testMessageID := "test-message-id-123"
	err = repo.MarkMessageAsSentTx(ctx, tx, messageID, testMessageID)
	assert.NoError(t, err)

	err = tx.Commit()
	assert.NoError(t, err)

	var isSent bool
	var storedMessageID string
	err = db.QueryRowContext(ctx, "SELECT is_sent, message_id FROM messages WHERE id = $1", messageID).
		Scan(&isSent, &storedMessageID)

	//then
	assert.NoError(t, err)
	assert.True(t, isSent)
	assert.Equal(t, testMessageID, storedMessageID)
}
