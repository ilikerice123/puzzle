package puzzle

// WebRequest Represents a request from the web that should be served
type WebRequest struct {
	Request
	url      string
	puzzleID string
}

// Pool represents a moving, interactable puzzle
type Pool struct {
	requests   chan WebRequest
	puzzles    map[string]Base
	numWorkers int
}

// Start starts the puzzle pool
func (p Pool) Start() {
	for i := 0; i < p.numWorkers; i++ {
		go p.worker(i + 1)
	}
}

func (p Pool) worker(id int) {
	for req := range p.requests {
		puzzle := p.puzzles[req.puzzleID]
		puzzle.Do(req.Request)
	}
}

func (p Pool) addPuzzle(file string, xSize int, ySize int) error {
	puzzle, err := NewPuzzle(file, xSize, ySize)
	if err != nil {
		return err
	}
	p.puzzles[puzzle.id] = puzzle
	return nil
}
