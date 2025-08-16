package engine

import (
	"fmt"
	"os"
	"sync"
	"time"

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

// AnalyzeDirectory analyzes all supported files in a directory
func (e *Engine) AnalyzeDirectory(rootPath string) (*ProjectAnalysis, error) {
	return e.AnalyzeDirectoryWithProgress(rootPath, nil)
}

// AnalyzeDirectoryWithProgress analyzes all supported files in a directory with progress reporting
func (e *Engine) AnalyzeDirectoryWithProgress(rootPath string, progressCallback func(current, total int, filePath string)) (*ProjectAnalysis, error) {
	// Create file walker
	walker := NewFileWalker(e.parserRegistry, e.config)
	
	// Start worker pool
	e.workerPool.Start()
	defer e.workerPool.Stop()
	
	// Walk directory to find files
	walkResultChan, err := walker.Walk(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}
	
	// Collect all files first to know total count
	var walkResults []WalkResult
	for result := range walkResultChan {
		if result.Error != nil {
			// Log error but continue
			fmt.Printf("Warning: %v\n", result.Error)
			continue
		}
		walkResults = append(walkResults, result)
	}
	
	totalFiles := len(walkResults)
	if totalFiles == 0 {
		return &ProjectAnalysis{
			RootPath:     rootPath,
			TotalFiles:   0,
			TotalLines:   0,
			Languages:    make(map[string]LanguageStats),
			FileResults:  []*parser.AnalysisResult{},
		}, nil
	}
	
	// Submit jobs to worker pool
	for _, walkResult := range walkResults {
		content, err := e.readFileContent(walkResult.FilePath)
		if err != nil {
			fmt.Printf("Warning: failed to read file %s: %v\n", walkResult.FilePath, err)
			continue
		}
		
		job := AnalysisJob{
			FilePath: walkResult.FilePath,
			Content:  content,
			Parser:   walkResult.Parser,
		}
		
		if err := e.workerPool.SubmitJob(job); err != nil {
			return nil, fmt.Errorf("failed to submit job for %s: %w", walkResult.FilePath, err)
		}
	}
	
	// Collect results
	var results []*parser.AnalysisResult
	var errors []error
	processedCount := 0
	
	resultChan := e.workerPool.GetResultChannel()
	
	for processedCount < totalFiles {
		select {
		case jobResult := <-resultChan:
			processedCount++
			
			if progressCallback != nil {
				filePath := ""
				if jobResult.Result != nil {
					filePath = jobResult.Result.FilePath
				}
				progressCallback(processedCount, totalFiles, filePath)
			}
			
			if jobResult.Error != nil {
				errors = append(errors, jobResult.Error)
				continue
			}
			
			if jobResult.Result != nil {
				results = append(results, jobResult.Result)
			}
		}
	}
	
	// Aggregate results into project analysis
	analysis := e.aggregateResults(rootPath, results)
	
	// Report any errors that occurred
	if len(errors) > 0 {
		fmt.Printf("Analysis completed with %d errors\n", len(errors))
		for _, err := range errors {
			fmt.Printf("  Error: %v\n", err)
		}
	}
	
	return analysis, nil
}

// AnalyzeFile analyzes a single file
func (e *Engine) AnalyzeFile(filePath string) (*parser.AnalysisResult, error) {
	// Get parser for file
	parser, err := e.parserRegistry.GetParser(filePath)
	if err != nil {
		return nil, fmt.Errorf("no parser available for file %s: %w", filePath, err)
	}
	
	// Read file content
	content, err := e.readFileContent(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	
	// Parse file
	result, err := parser.Parse(filePath, content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}
	
	return result, nil
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

// readFileContent reads the content of a file with size limits
func (e *Engine) readFileContent(filePath string) ([]byte, error) {
	// Check file size if limit is set
	if e.config.MaxFileSize > 0 {
		info, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}
		
		if info.Size() > e.config.MaxFileSize {
			return nil, fmt.Errorf("file %s exceeds size limit (%d bytes)", filePath, e.config.MaxFileSize)
		}
	}
	
	return os.ReadFile(filePath)
}

// aggregateResults combines individual file analysis results into a project analysis
func (e *Engine) aggregateResults(rootPath string, results []*parser.AnalysisResult) *ProjectAnalysis {
	analysis := &ProjectAnalysis{
		RootPath:     rootPath,
		TotalFiles:   len(results),
		TotalLines:   0,
		Languages:    make(map[string]LanguageStats),
		FileResults:  results,
		GeneratedAt:  time.Now(),
	}
	
	// Aggregate statistics
	for _, result := range results {
		analysis.TotalLines += result.LineCount
		
		// Update language statistics
		langStats, exists := analysis.Languages[result.Language]
		if !exists {
			langStats = LanguageStats{}
		}
		
		langStats.FileCount++
		langStats.LineCount += result.LineCount
		langStats.FunctionCount += len(result.Functions)
		langStats.ClassCount += len(result.Classes)
		langStats.Complexity += result.Complexity
		
		analysis.Languages[result.Language] = langStats
	}
	
	return analysis
}

// GetSupportedExtensions returns all supported file extensions
func (e *Engine) GetSupportedExtensions() []string {
	return e.parserRegistry.GetSupportedExtensions()
}

// IsFileSupported checks if a file is supported for analysis
func (e *Engine) IsFileSupported(filePath string) bool {
	return e.parserRegistry.IsSupported(filePath)
}

// GetFileWalkerStats returns statistics about files in a directory
func (e *Engine) GetFileWalkerStats(rootPath string) (*WalkStats, error) {
	walker := NewFileWalker(e.parserRegistry, e.config)
	return walker.GetStats(rootPath)
}