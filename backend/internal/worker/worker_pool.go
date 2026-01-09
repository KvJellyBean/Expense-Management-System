package worker

import (
	"context"
	"expense-management-system/internal/usecase"
	"expense-management-system/pkg/logger"
	"sync"
)

type WorkerPool struct {
	paymentChan    chan usecase.PaymentJob
	paymentService *PaymentService
	workerCount    int
	maxRetries     int
	wg             sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
}

func NewWorkerPool(paymentChan chan usecase.PaymentJob, paymentService *PaymentService, workerCount, maxRetries int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		paymentChan:    paymentChan,
		paymentService: paymentService,
		workerCount:    workerCount,
		maxRetries:     maxRetries,
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (wp *WorkerPool) Start() {
	logger.InfoLogger.Printf("Starting worker pool with %d workers", wp.workerCount)

	for i := 1; i <= wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	logger.InfoLogger.Printf("Worker %d started", id)

	for {
		select {
		case <-wp.ctx.Done():
			logger.InfoLogger.Printf("Worker %d shutting down", id)
			return
		case job, ok := <-wp.paymentChan:
			if !ok {
				logger.InfoLogger.Printf("Worker %d: channel closed, shutting down", id)
				return
			}

			logger.InfoLogger.Printf("Worker %d processing payment for expense %d", id, job.ExpenseID)

			if err := wp.paymentService.ProcessPaymentWithRetry(wp.ctx, job, wp.maxRetries); err != nil {
				logger.ErrorLogger.Printf("Worker %d failed to process payment for expense %d: %v", id, job.ExpenseID, err)
			}
		}
	}
}

func (wp *WorkerPool) Stop() {
	logger.InfoLogger.Println("Stopping worker pool...")
	wp.cancel()
	close(wp.paymentChan)
	wp.wg.Wait()
	logger.InfoLogger.Println("Worker pool stopped")
}
