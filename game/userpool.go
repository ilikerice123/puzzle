package game

import (
	"time"

	"github.com/ilikerice123/puzzle/store"
)

// GlobalUserPool represents all the players
var GlobalUserPool UserPoolBase

// UserPoolBase is the interface for a pool of users
type UserPoolBase interface {
	AddUser(*store.User)

	GetUser(string) *store.User

	Prune()
}

// UserPool implements UserPoolBase
type UserPool struct {
	users map[string]*store.User
}

// InitUserPool assigns value to globalUserPool
func InitUserPool() {
	GlobalUserPool = NewUserPool()
}

// NewUserPool creates a new user pool
func NewUserPool() *UserPool {
	p := &UserPool{users: make(map[string]*store.User)}
	scheduler := time.NewTicker(12 * time.Hour)
	go func() {
		for range scheduler.C {
			p.Prune()
		}
	}()
	return p
}

// AddUser adds a user to the pool
func (p *UserPool) AddUser(u *store.User) {
	p.users[u.ID] = u
	return
}

// GetUser gets a user from the pool
func (p *UserPool) GetUser(id string) *store.User {
	return p.users[id]
}

// AuthUser authenticates a user from the pool
func (p *UserPool) AuthUser(name string, password string) {

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
