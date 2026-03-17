// Package worker предоставляет pool для обработки пакетов.
package worker

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/internal/protocol"
)

// Task represents a packet processing task.
type Task struct {
	Packet *protocol.Packet
}

// Processor processes a packet.
type Processor interface {
	Process(pkt *protocol.Packet) error
}

// Pool is a worker pool for packet processing.
type Pool struct {
	mu        sync.RWMutex
	workers   int
	tasks     chan Task
	processor Processor
	wg        sync.WaitGroup
	started   atomic.Bool
	stopped   atomic.Bool

	// Stats
	processed atomic.Uint64
	errors    atomic.Uint64

	// Logger
	log *logger.Logger
}

// New creates a new worker pool.
func New(workers int, processor Processor) *Pool {
	log := logger.Default().With("worker")
	return &Pool{
		workers:   workers,
		processor: processor,
		log:       log,
	}
}

// Start starts the worker pool.
func (p *Pool) Start(ctx context.Context) error {
	if !p.started.CompareAndSwap(false, true) {
		return fmt.Errorf("already started")
	}

	p.tasks = make(chan Task, p.tasksBufferSize())

	p.log.Info("Starting worker pool with %d workers", p.workers)

	p.wg.Add(p.workers)
	for i := 0; i < p.workers; i++ {
		go p.worker(ctx, i)
	}

	return nil
}

// tasksBufferSize returns the buffer size for the task queue.
func (p *Pool) tasksBufferSize() int {
	return p.workers * 10
}

// Submit submits a task to the pool.
func (p *Pool) Submit(pkt *protocol.Packet) error {
	if p.stopped.Load() {
		return fmt.Errorf("pool stopped")
	}

	select {
	case p.tasks <- Task{Packet: pkt}:
		return nil
	default:
		p.log.Warn("Task queue full, dropping packet 0x%02X", pkt.ID)
		return fmt.Errorf("task queue full")
	}
}

// worker processes tasks from the queue.
func (p *Pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()
	p.log.Debug("Worker %d started", id)

	for {
		select {
		case <-ctx.Done():
			p.log.Debug("Worker %d stopped by context", id)
			return

		case task, ok := <-p.tasks:
			if !ok {
				p.log.Debug("Worker %d stopped (channel closed)", id)
				return
			}

			p.processTask(task, id)
		}
	}
}

// processTask processes a single task.
func (p *Pool) processTask(task Task, workerID int) {
	defer func() {
		if r := recover(); r != nil {
			p.log.Error("Worker %d panic: %v", workerID, r)
			p.errors.Add(1)
		}
	}()

	err := p.processor.Process(task.Packet)
	if err != nil {
		p.log.Error("Worker %d failed to process packet 0x%02X: %v",
			workerID, task.Packet.ID, err)
		p.errors.Add(1)
		return
	}

	p.processed.Add(1)
}

// Stop stops the worker pool gracefully.
func (p *Pool) Stop() {
	if !p.stopped.CompareAndSwap(false, true) {
		return
	}

	p.log.Info("Stopping worker pool...")
	close(p.tasks)
	p.wg.Wait()
	p.log.Info("Worker pool stopped")
}

// Stats returns worker pool statistics.
func (p *Pool) Stats() (processed, errors uint64) {
	return p.processed.Load(), p.errors.Load()
}
