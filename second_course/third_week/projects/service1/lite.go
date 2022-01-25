package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {

	var sessManager SessionManagerI

	sessManager = NewSessManager()

	// создаем сессию
	sessID, err := sessManager.Create(
		&Session{
			Login:     "rvasily",
			Useragent: "chrome",
		},
	)
	fmt.Println("sessId", sessID, err)

	// проверяем сессию
	sess := sessManager.Check(
		&SessionID{
			ID: sessID.ID,
		})

	fmt.Println("sess", sess)

	// удаляем сессию
	sessManager.Delete(
		&SessionID{
			ID: sessID.ID,
		})

	// проверяем еще раз
	sess = sessManager.Check(
		&SessionID{
			ID: sessID.ID,
		},
	)
	fmt.Println("sess", sess)
}

type Session struct {
	Login     string
	Useragent string
}

type SessionID struct {
	ID string
}

const sessKeyLen = 10

type SessionManager struct {
	mu      sync.Mutex
	session map[SessionID]*Session
}

type SessionManagerI interface {
	Create(*Session) (*SessionID, error)
	Check(*SessionID) *Session
	Delete(*SessionID)
}

func NewSessManager() *SessionManager {
	return &SessionManager{
		mu:      sync.Mutex{},
		session: map[SessionID]*Session{},
	}
}

func (sm *SessionManager) Create(in *Session) (*SessionID, error) {
	sm.mu.Lock()
	id := SessionID{RandStringRunes(sessKeyLen)}
	sm.mu.Unlock()
	sm.session[id] = in
	return &id, nil
}

func (sm *SessionManager) Check(in *SessionID) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sess, ok := sm.session[*in]; ok {
		return sess
	}
	return nil
}

func (sm *SessionManager) Delete(in *SessionID) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.session, *in)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
