// Package store keeps the games being played keyed by session ID
package store

import (
	"sync"
	"time"
	"uuid"

	"github.com/Joshwantt/gordle/cmd/gordle/session"
)

func New() *Store {
	return &Store{
		sessions: make(map[uuid.UUID]*session.Session),
	}
}

type Store struct {
	mutex    sync.Mutex
	sessions map[uuid.UUID]*session.Session
}

func (store *Store) AddSession(newSession session.Session) uuid.UUID {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	sessionID := uuid.New()

	newSession.LastInteraction = time.Now()

	store.sessions[sessionID] = &newSession

	return sessionID
}

func (store *Store) DeleteSession(sessionID uuid.UUID) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	delete(store.sessions, sessionID)
}

func (store *Store) GetSession(sessionID uuid.UUID) (session.Session, bool) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	storedSession, exists := store.sessions[sessionID]

	if !exists {
		return session.Session{}, false
	}

	return *storedSession, true
}

func (store *Store) UpdateSession(sessionID uuid.UUID, change func(storedSession *session.Session)) bool {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	storedSession, exists := store.sessions[sessionID]

	if !exists {
		return false
	}

	change(storedSession)
	storedSession.LastInteraction = time.Now()
	return true
}

func (store *Store) DeleteInactive(cutoff time.Time) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	for sessionID, storedSession := range store.sessions {
		if storedSession.LastInteraction.Before(cutoff) {
			delete(store.sessions, sessionID)
		}
	}
}

func (store *Store) DeleteInactiveTicker(interval time.Duration, maxInactivity time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		store.DeleteInactive(time.Now().Add(-maxInactivity))
	}
}
