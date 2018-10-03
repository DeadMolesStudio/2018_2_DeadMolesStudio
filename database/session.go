package database

import (
	"fmt"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/models"
)

var sessions = make(models.Session)

func CreateNewSession(sessionID string, userID int) error {
	sessions[sessionID] = userID

	return nil
}

func DeleteSession(sessionID string) error {
	delete(sessions, sessionID)

	return nil
}

func GetIDFromSession(sessionID string) (int, error) {
	id, ok := sessions[sessionID]
	if !ok {
		return -1, fmt.Errorf("no session in database")
	}

	return id, nil
}

func CheckExistenceOfSession(sessionID string) (bool, error) {
	_, ok := sessions[sessionID]
	if !ok {
		return false, nil
	}

	return true, nil
}
