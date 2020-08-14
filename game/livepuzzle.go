package game

import "sync"

// LivePuzzleBase represents a threadsafe puzzle object
type LivePuzzleBase interface {
	Start()

	AddRequest(*Request)

	AddCallback(func(*Update))

	ID() string
}

// LivePuzzle implements the LivePuzzleBase interface
type LivePuzzle struct {
	Puzzle PuzzleBase

	requests     chan *Request
	updates      chan *Update
	callbacks    []func(*Update)
	callbackLock sync.Locker
}

// NewLivePuzzle creates new live puzzle
func NewLivePuzzle(
	id string,
	file string,
	ySize int,
	xSize int,
	users UserPoolBase) *LivePuzzle {
	updates := make(chan *Update)
	p := NewPuzzle(id, file, ySize, xSize, updates, users)
	if p == nil {
		return nil
	}

	return &LivePuzzle{
		Puzzle:       p,
		requests:     make(chan *Request),
		updates:      updates,
		callbacks:    make([]func(*Update), 0),
		callbackLock: &sync.Mutex{}}
}

// ID returns the id of the puzzle
func (p *LivePuzzle) ID() string {
	return p.Puzzle.GetID()
}

// AddRequest adds a request to the LivePuzzle
func (p *LivePuzzle) AddRequest(r *Request) {
	p.requests <- r
}

// AddCallback registers a function callback
func (p *LivePuzzle) AddCallback(f func(*Update)) {
	p.callbackLock.Lock()
	p.callbacks = append(p.callbacks, f)
	p.callbackLock.Unlock()
}

// Start starts the puzzle
func (p *LivePuzzle) Start() {
	// goroutine to process requests
	go func() {
		for req := range p.requests {
			p.Puzzle.Do(*req)
		}
	}()
	// goroutine to send updates
	go func() {
		for update := range p.updates {
			p.callbackLock.Lock()
			for _, f := range p.callbacks {
				f(update)
			}
			p.callbackLock.Unlock()
		}
	}()
}
