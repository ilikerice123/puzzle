package puzzle

type action int

// actions that could be performed
const (
	SWAP action = iota
	HOLD
	RELEASE
)

// Request representing a request to move something
type Request struct {
	action action
	userID string
	pieceX int
	pieceY int
}

// Update representing a state change of the puzzle
// if update is a SWAP, piece1ID and piece2ID are populated
// if update is a HOLD, piece1ID and userID are populated
// if update is a RELEASE, piece1ID is populated
type Update struct {
	action   action
	userID   string
	piece1ID int
	piece2ID int
}
