package engine

import (
	"fmt"
	"sync"

	"github.com/tito-sala/codebasereaderv2/internal/parser"
)

// Engine orchestrates the codebase analysis process
type Engine struct {
	parserRegistry *parser.ParserRegistry
	config         *Config
	workerPool     *WorkerPool
}

// NewEngine creates a new analysis engine with the given configuration
func NewEngine(config *Config) *Engine {
	if config == nil {
		config = DefaultConfig()
	}

	return &Engine{
		parserRegistry: parser.NewParserRegistry(),
		config:         config,
		workerPool:     NewWorkerPool(config.MaxWorkers),
	}
}

// GetParserRegistry returns the parser registry for registering new parsers
func (e *Engine) GetParserRegistry() *parser.ParserRegistry {
	return e.parserRegistry
}

// GetConfig returns the current configuration
func (e *Engine) GetConfig() *Config {
	return e.config
}

// UpdateConfig updates the engine configuration
func (e *Engine) UpdateConfig(config *Config) {
	if config != nil {
		e.config = config
		// Update worker pool if max workers changed
		if e.workerPool.maxWorkers != config.MaxWorkers {
			e.workerPool.Stop()
			e.workerPool = NewWorkerPool(config.MaxWorkers)
		}
	}
}

// WorkerPool manages concurrent file processing
type WorkerPool struct {
	maxWorkers  int
	jobQueue    chan AnalysisJob
	resultQueue chan AnalysisJobResult
	workers     []*worker
	wg          sync.WaitGroup
	stopChan    chan struct{}
	running     bool
	mutex       sync.RWMutex
}

// NewWorkerPool creates a new worker pool with the specified number of workers
func NewWorkerPool(maxWorkers int) *WorkerPool {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}

	return &WorkerPool{
		maxWorkers:  maxWorkers,
		jobQueue:    make(chan AnalysisJob, maxWorkers*2),
		resultQueue: make(chan AnalysisJobResult, maxWorkers*2),
		workers:     make([]*worker, 0, maxWorkers),
		stopChan:    make(chan struct{}),
		running:     false,
	}
}

// Start initializes and starts all workers in the pool
func (wp *WorkerPool) Start() {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	if wp.running {
		return
	}

	wp.running = true
	wp.workers = make([]*worker, wp.maxWorkers)

	for i := 0; i < wp.maxWorkers; i++ {
		wp.workers[i] = &worker{
			id:          i,
			jobQueue:    wp.jobQueue,
			resultQueue: wp.resultQueue,
			stopChan:    wp.stopChan,
		}
		wp.wg.Add(1)
		go wp.workers[i].start(&wp.wg)
	}
}

// Stop gracefully shuts down all workers
func (wp *WorkerPool) Stop() {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	if !wp.running {
		return
	}

	close(wp.stopChan)
	wp.wg.Wait()
	wp.running = false
}

// SubmitJob adds a job to the worker pool queue
func (wp *WorkerPool) SubmitJob(job AnalysisJob) error {
	wp.mutex.RLock()
	defer wp.mutex.RUnlock()

	if !wp.running {
		return fmt.Errorf("worker pool is not running")
	}

	select {
	case wp.jobQueue <- job:
		return nil
	default:
		return fmt.Errorf("job queue is full")
	}
}

// GetResultChannel returns the channel for receiving job results
func (wp *WorkerPool) GetResultChannel() <-chan AnalysisJobResult {
	return wp.resultQueue
}

// worker represents a single worker in the pool
type worker struct {
	id          int
	jobQueue    <-chan AnalysisJob
	resultQueue chan<- AnalysisJobResult
	stopChan    <-chan struct{}
}

// start begins the worker's job processing loop
func (w *worker) start(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case job := <-w.jobQueue:
			result, err := job.Parser.Parse(job.FilePath, job.Content)
			w.resultQueue <- AnalysisJobResult{
				Result: result,
				Error:  err,
			}
		case <-w.stopChan:
			return
		}
	}
}