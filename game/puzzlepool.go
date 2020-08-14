package game

// PuzzlePoolBase represents a PuzzlePool interface
type PuzzlePoolBase interface {
	AddPuzzle(LivePuzzleBase)

	GetPuzzle(id string) LivePuzzleBase
}

// PuzzlePool represents the pool of interactable puzzles
type PuzzlePool struct {
	puzzles map[string]LivePuzzleBase
}

// GlobalPuzzlePool represents all the puzzles
var GlobalPuzzlePool PuzzlePoolBase

// InitPuzzlePool assigns value to globalUserPool
func InitPuzzlePool() {
	GlobalPuzzlePool = &PuzzlePool{puzzles: make(map[string]LivePuzzleBase)}
}

// AddPuzzle adds a puzzle to the pool
func (p *PuzzlePool) AddPuzzle(puzzle LivePuzzleBase) {
	p.puzzles[puzzle.ID()] = puzzle
}

// GetPuzzle gets a puzzle from the pool
func (p *PuzzlePool) GetPuzzle(id string) LivePuzzleBase {
	return p.puzzles[id]
}
