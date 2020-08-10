package puzzle

// WebRequest Represents a request from the web that should be served
type WebRequest struct {
	Request
	url      string
	puzzleID string
}

// Pool represents a moving, interactable puzzle
type Pool struct {
	requests chan WebRequest
	puzzles  map[string]Puzzle
}
