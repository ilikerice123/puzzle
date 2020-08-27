package game

import (
	"log"
	"os"
	"time"

	"github.com/ilikerice123/puzzle/fs"
)

// PuzzlePoolBase represents a PuzzlePool interface
type PuzzlePoolBase interface {
	AddPuzzle(LivePuzzleBase)

	GetPuzzle(id string) LivePuzzleBase

	Prune()
}

// PuzzlePool represents the pool of interactable puzzles
type PuzzlePool struct {
	puzzles map[string]LivePuzzleBase
}

// GlobalPuzzlePool represents all the puzzles
var GlobalPuzzlePool PuzzlePoolBase

// InitPuzzlePool assigns value to globalUserPool
func InitPuzzlePool() {
	GlobalPuzzlePool = NewPuzzlePool()
}

// NewPuzzlePool creates a new puzzle pool
func NewPuzzlePool() *PuzzlePool {
	p := &PuzzlePool{puzzles: make(map[string]LivePuzzleBase)}
	scheduler := time.NewTicker(12 * time.Hour)
	go func() {
		for range scheduler.C {
			p.Prune()
		}
	}()
	return p
}

// AddPuzzle adds a puzzle to the pool
func (p *PuzzlePool) AddPuzzle(puzzle LivePuzzleBase) {
	p.puzzles[puzzle.ID()] = puzzle
}

// GetPuzzle gets a puzzle from the pool
func (p *PuzzlePool) GetPuzzle(id string) LivePuzzleBase {
	return p.puzzles[id]
}

// Prune removes all puzzles that are complete, including their images
// also removes all directories that doesn't have a puzzle associated to it
func (p *PuzzlePool) Prune() {
	for id, puzzle := range p.puzzles {
		if puzzle.Complete() {
			// remove directory for puzzle
			dir := "images/" + id
			if fs.DirExists(dir) {
				os.RemoveAll(dir)
			}
			// remove active puzzles from user
			for userID := range puzzle.Results() {
				if user := GlobalUserPool.GetUser(userID); user != nil {
					delete(user.PieceCount, id)
				}
			}
			delete(p.puzzles, id)
		}
	}
	imageFolder, err := os.Open("images")
	if err != nil {
		log.Printf("Error opening images/: %s", err.Error())
	}
	folders, err := imageFolder.Readdirnames(0)
	if err != nil {
		log.Printf("Error reading images/*: %s", err.Error())
	}
	for _, name := range folders {
		if _, exists := p.puzzles[name]; !exists {
			os.RemoveAll("images/" + name)
		}
	}
}
