package game

type action int

// actions that could be performed
const (
	SWAP action = iota
	HOLD
	JOIN
	LEAVE
)

// Request representing a request to move something
type Request struct {
	Action action
	UserID string
	PieceX int
	PieceY int
}

// Update representing a state change of the puzzle
// delta represents change in number of correct pieces
// - if Action is a SWAP, piece1ID and piece2ID are populated
//   * swap is implicitly a RELEASE state change if piece1ID == piece2
// - if Action is a HOLD, piece1ID and userID are populated
// - if Action is a JOIN or LEAVE, only userID is populated
type Update struct {
	ID       int
	Action   action
	UserID   string
	Piece1ID int
	Piece2ID int
	Delta    int
}
