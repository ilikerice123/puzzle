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
	Action action `json:"action"`
	UserID string `json:"userID"`
	PieceX int    `json:"pieceX"`
	PieceY int    `json:"pieceY"`
}

// Update representing a state change of the puzzle
// delta represents change in number of correct pieces
// - if Action is a SWAP, piece1ID and piece2ID are populated
//   * swap is implicitly a RELEASE state change if piece1ID == piece2
// - if Action is a HOLD, piece1ID and userID are populated
// - if Action is a JOIN or LEAVE, only userID is populated
type Update struct {
	ID       int    `json:"id"`
	Action   action `json:"action"`
	UserID   string `json:"userID"`
	Piece1ID int    `json:"piece1ID"`
	Piece2ID int    `json:"piece2ID"`
	Delta    int    `json:"delta"`
}
