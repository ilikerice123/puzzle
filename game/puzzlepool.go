package game

// WebRequest Represents a request from the web that should be served
type WebRequest struct {
	Request
	url      string
	puzzleID string
}

// PoolBase represents a PuzzlePool interface
type PoolBase interface {
	Start()

	AddRequest(r WebRequest)

	AddPuzzle(file string, xSize int, ySize int) error
}

// Pool represents the pool of interactable puzzles
type Pool struct {
	requests   chan WebRequest
	puzzles    map[string]Base
	numWorkers int
}

// GlobalPuzzlePool represents all the puzzles
var GlobalPuzzlePool PoolBase

// InitPuzzlePool assigns value to globalUserPool
func InitPuzzlePool() {
	GlobalPuzzlePool = &Pool{}
}

// Start starts the puzzle pool
func (p *Pool) Start() {
	for i := 0; i < p.numWorkers; i++ {
		go p.worker(i + 1)
	}
}

func (p *Pool) worker(id int) {
	for req := range p.requests {
		puzzle := p.puzzles[req.puzzleID]
		puzzle.Do(req.Request)
	}
}

// AddPuzzle adds a puzzle to the pool
func (p *Pool) AddPuzzle(file string, xSize int, ySize int) error {
	puzzle, err := NewPuzzle(file, xSize, ySize)
	if err != nil {
		return err
	}
	p.puzzles[puzzle.id] = puzzle
	return nil
}

// AddRequest adds a web request to be processed by the workers
func (p *Pool) AddRequest(r WebRequest) {
	return
}
