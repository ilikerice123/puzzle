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
	LifetimePieces int            `json:"lifetime_pieces"`
}

// InitUserPool assigns value to globalUserPool
func InitUserPool() {
	GlobalUserPool = &UserPool{users: make(map[string]*User)}
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
