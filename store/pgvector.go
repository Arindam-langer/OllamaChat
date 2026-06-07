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
		embedding vector(768),
		source_file text,
		file_hash text
	);`

	_, err := s.pool.Exec(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to create documents table: %w", err)
	}

	// Migrate existing tables that lack the new columns
	migrate := `
	ALTER TABLE documents ADD COLUMN IF NOT EXISTS source_file text;
	ALTER TABLE documents ADD COLUMN IF NOT EXISTS file_hash text;`
	_, err = s.pool.Exec(ctx, migrate)
	if err != nil {
		return fmt.Errorf("failed to migrate documents table: %w", err)
	}

	return nil
}

// GetFileHashes returns a map of source_file -> file_hash for all tracked documents
func (s *VectorStore) GetFileHashes(ctx context.Context) (map[string]string, error) {
	rows, err := s.pool.Query(ctx,
		"SELECT DISTINCT source_file, file_hash FROM documents WHERE source_file IS NOT NULL",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query file hashes: %w", err)
	}
	defer rows.Close()

	hashes := make(map[string]string)
	for rows.Next() {
		var file, hash string
		if err := rows.Scan(&file, &hash); err != nil {
			return nil, fmt.Errorf("failed to scan file hash row: %w", err)
		}
		hashes[file] = hash
	}
	return hashes, nil
}

// DeleteBySourceFile removes all embeddings associated with a specific source file
func (s *VectorStore) DeleteBySourceFile(ctx context.Context, sourceFile string) error {
	_, err := s.pool.Exec(ctx, "DELETE FROM documents WHERE source_file = $1", sourceFile)
	if err != nil {
		return fmt.Errorf("failed to delete embeddings for %s: %w", sourceFile, err)
	}
	log.Printf("Deleted old embeddings for: %s", sourceFile)
	return nil
}

// CleanupOrphans removes embeddings for files that no longer exist on disk
func (s *VectorStore) CleanupOrphans(ctx context.Context, activeFiles map[string]bool) (int, error) {
	storedHashes, err := s.GetFileHashes(ctx)
	if err != nil {
		return 0, err
	}

	removed := 0
	for storedFile := range storedHashes {
		if !activeFiles[storedFile] {
			if err := s.DeleteBySourceFile(ctx, storedFile); err != nil {
				return removed, err
			}
			removed++
		}
	}
	return removed, nil
}

func (s *VectorStore) InsertEmbeddings(ctx context.Context, embeddings []extract.Embedding, sourceFile, fileHash string) error {
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
			"INSERT INTO documents (content, embedding, source_file, file_hash) VALUES ($1, $2, $3, $4)",
			emb.Text, pgvector.NewVector(emb.Vector), sourceFile, fileHash,
		)
		if err != nil {
			return fmt.Errorf("failed to insert embedding: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully inserted %d documents for %s into vector store", len(embeddings), sourceFile)
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

func (s *VectorStore) Flush(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, "TRUNCATE TABLE documents RESTART IDENTITY;")
	if err != nil {
		return fmt.Errorf("failed to truncate documents table: %w", err)
	}
	return nil
}

