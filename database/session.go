package database

import (
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/models"
)

var sessions = models.Sessions{
	Sessions: make(map[string]uint),
}

func CreateNewSession(sessionID string, userID uint) error {
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

func GetIDFromSession(sessionID string) (uint, error) {
	sessions.Lock()
	id, ok := sessions.Sessions[sessionID]
	sessions.Unlock()
	if !ok {
		return 0, ErrSessionNotFound
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
