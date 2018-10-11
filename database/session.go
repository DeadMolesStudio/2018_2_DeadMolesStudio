package database

import (
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/models"
)

var sessions = models.Sessions{
	Sessions: make(map[string]int),
}

func CreateNewSession(sessionID string, userID int) error {
	sessions.Lock()
	sessions.Sessions[sessionID] = userID
	sessions.Unlock()

	return nil
}

func DeleteSession(sessionID string) error {
	sessions.Lock()
	delete(sessions.Sessions, sessionID)
	sessions.Unlock()

	return nil
}

func GetIDFromSession(sessionID string) (int, error) {
	sessions.Lock()
	id, ok := sessions.Sessions[sessionID]
	sessions.Unlock()
	if !ok {
		return -1, ErrSessionNotFound
	}

	return id, nil
}

func CheckExistenceOfSession(sessionID string) (bool, error) {
	sessions.Lock()
	_, ok := sessions.Sessions[sessionID]
	sessions.Unlock()
	if !ok {
		return false, nil
	}

	return true, nil
}
