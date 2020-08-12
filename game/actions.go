package game

type action int

// actions that could be performed
const (
	SWAP action = iota
	HOLD
)

// Request representing a request to move something
type Request struct {
	UserID string
	PieceX int
	PieceY int
}

// Update representing a state change of the puzzle
// delta represents change in number of correct pieces
// - if update is a SWAP, piece1ID and piece2ID are populated
//   * swap is implicitly a RELEASE state change if piece1ID == piece2
// - if update is a HOLD, piece1ID and userID are populated
type Update struct {
	Action   action
	UserID   string
	Piece1ID int
	Piece2ID int
	Delta    int
}
