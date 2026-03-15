package worker

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/dukerupert/skalkaho/internal/repository"
)

// JobHandler processes a single import job.
type JobHandler func(ctx context.Context, job repository.MaterialImportJob) error

// Worker polls for pending import jobs and dispatches them to a handler.
type Worker struct {
	queries  *repository.Queries
	db       *sql.DB
	interval time.Duration
	handler  JobHandler
	logger   *slog.Logger
}

// New creates a new background worker.
func New(queries *repository.Queries, db *sql.DB, interval time.Duration, logger *slog.Logger) *Worker {
	return &Worker{
		queries:  queries,
		db:       db,
		interval: interval,
		logger:   logger,
	}
}

// SetHandler registers the job handler function.
func (w *Worker) SetHandler(h JobHandler) {
	w.handler = h
}

// Start begins the polling loop. It blocks until ctx is cancelled.
func (w *Worker) Start(ctx context.Context) {
	w.logger.Info("import worker started", "interval", w.interval)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("import worker shutting down")
			return
		case <-ticker.C:
			w.poll(ctx)
		}
	}
}

func (w *Worker) poll(ctx context.Context) {
	if w.handler == nil {
		return
	}

	job, err := w.queries.ClaimNextImportJob(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return // no pending jobs
		}
		w.logger.Error("claiming import job", "error", err)
		return
	}

	w.logger.Info("processing import job", "job_id", job.ID, "supplier_id", job.SupplierID, "filename", job.Filename)

	// Run handler with panic recovery
	func() {
		defer func() {
			if r := recover(); r != nil {
				errMsg := fmt.Sprintf("panic: %v", r)
				w.logger.Error("import job panicked", "job_id", job.ID, "panic", r)
				_, _ = w.queries.FailImportJob(ctx, repository.FailImportJobParams{
					ID:           job.ID,
					ErrorMessage: sql.NullString{String: errMsg, Valid: true},
				})
			}
		}()

		if err := w.handler(ctx, job); err != nil {
			w.logger.Error("import job failed", "job_id", job.ID, "error", err)
			_, _ = w.queries.FailImportJob(ctx, repository.FailImportJobParams{
				ID:           job.ID,
				ErrorMessage: sql.NullString{String: err.Error(), Valid: true},
			})
		}
	}()
}
