package store

import (
	"context"
	"fmt"
	"log"
	"time"

	extract "github.com/Arindam-langer/OllamaChat/preprocessing"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
)

type VectorStore struct {
	pool *pgxpool.Pool
}

func Connect(ctx context.Context, databaseURL string) (*VectorStore, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// This requires the 'vector' extension to ALREADY exist in the database!
		return pgxvec.RegisterTypes(ctx, conn)
	}

	var pool *pgxpool.Pool
	// Retry connection a few times in case docker is still starting
	for i := 0; i < 5; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, config)
		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				return &VectorStore{pool: pool}, nil
			}
			pool.Close()
		}
		log.Printf("Waiting for database connection... (%d/5)\n", i+1)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database: %w", err)
}

func (s *VectorStore) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

func (s *VectorStore) InitDB(ctx context.Context) error {
	// Using 768 dimensions because nomic-embed-text generates 768-dimensional embeddings
	schema := `
	CREATE TABLE IF NOT EXISTS documents (
		id bigserial PRIMARY KEY,
		content text NOT NULL,
		embedding vector(768)
	);`

	_, err := s.pool.Exec(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to create documents table: %w", err)
	}

	return nil
}

func (s *VectorStore) InsertEmbeddings(ctx context.Context, embeddings []extract.Embedding) error {
	if len(embeddings) == 0 {
		return nil
	}

	// Begin a transaction for bulk insert
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, emb := range embeddings {
		_, err := tx.Exec(ctx,
			"INSERT INTO documents (content, embedding) VALUES ($1, $2)",
			emb.Text, pgvector.NewVector(emb.Vector),
		)
		if err != nil {
			return fmt.Errorf("failed to insert embedding: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully inserted %d documents into vector store", len(embeddings))
	return nil
}

func (s *VectorStore) SearchSimilar(ctx context.Context, queryVector []float32, limit int) ([]string, error) {
	rows, err := s.pool.Query(ctx,
		"SELECT content FROM documents ORDER BY embedding <=> $1 LIMIT $2",
		pgvector.NewVector(queryVector), limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search similar documents: %w", err)
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, content)
	}

	return results, nil
}
