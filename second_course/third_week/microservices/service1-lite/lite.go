package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {

	// создаем сессию
	sessId, err := AuthCreateSession(
		&Session{
			Login:     "rvasily",
			Useragent: "chrome",
		})
	fmt.Println("sessId", sessId, err)

	// проеряем сессию
	sess := AuthCheckSession(
		&SessionID{
			ID: sessId.ID,
		})
	fmt.Println("sess", sess)

	// удаляем сессию
	AuthSessionDelete(
		&SessionID{
			ID: sessId.ID,
		})

	// проверяем еще раз
	sess = AuthCheckSession(
		&SessionID{
			ID: sessId.ID,
		})
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

var (
	sessions = map[SessionID]*Session{}
	mu       = &sync.RWMutex{}
)

func AuthCreateSession(in *Session) (*SessionID, error) {
	mu.Lock()
	id := SessionID{RandStringRunes(sessKeyLen)}
	mu.Unlock()
	sessions[id] = in
	return &id, nil
}

func AuthCheckSession(in *SessionID) *Session {
	mu.RLock()
	defer mu.RUnlock()
	if sess, ok := sessions[*in]; ok {
		return sess
	}
	return nil
}

func AuthSessionDelete(in *SessionID) {
	mu.Lock()
	defer mu.Unlock()
	delete(sessions, *in)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
