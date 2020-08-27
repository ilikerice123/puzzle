package game

import (
	"time"

	"github.com/google/uuid"
)

// GlobalUserPool represents all the players
var GlobalUserPool UserPoolBase

// UserPoolBase is the interface for a pool of users
type UserPoolBase interface {
	AddUser(*User)

	GetUser(string) *User

	Prune()
}

// UserPool implements UserPoolBase
type UserPool struct {
	users map[string]*User
}

// User represents a user
type User struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Created        time.Time      `json:"created"`
	PieceCount     map[string]int `json:"-"`
	LifetimePieces int            `json:"lifetimePieces"`
}

// InitUserPool assigns value to globalUserPool
func InitUserPool() {
	GlobalUserPool = NewUserPool()
}

// NewUserPool creates a new user pool
func NewUserPool() *UserPool {
	p := &UserPool{users: make(map[string]*User)}
	scheduler := time.NewTicker(12 * time.Hour)
	go func() {
		for range scheduler.C {
			p.Prune()
		}
	}()
	return p
}

// NewUser creates a new user
func NewUser(name string) *User {
	return &User{
		ID:             uuid.New().String(),
		Name:           name,
		Created:        time.Now(),
		PieceCount:     make(map[string]int),
		LifetimePieces: 0}
}

// AddUser adds a user to the pool
func (p *UserPool) AddUser(u *User) {
	p.users[u.ID] = u
	return
}

// GetUser gets a user from the pool
func (p *UserPool) GetUser(id string) *User {
	return p.users[id]
}

// DeleteUser deletes a user from the pool
func (p *UserPool) DeleteUser(id string) {
	delete(p.users, id)
	return
}

// Prune removes all puzzles from pieceCount that no longer exist
func (p *UserPool) Prune() {
	for _, user := range p.users {
		for puzzleID := range user.PieceCount {
			if puzzle := GlobalPuzzlePool.GetPuzzle(puzzleID); puzzle == nil {
				delete(user.PieceCount, puzzleID)
			}
		}
	}
}
