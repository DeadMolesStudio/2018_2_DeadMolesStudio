package sessions

import (
	"errors"

	"github.com/gomodule/redigo/redis"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/logger"
)

var sm *sessionManager

var (
	ErrKeyNotFound = errors.New("key not found")
)

type sessionManager struct {
	redisConn redis.Conn
}

func (sm *sessionManager) Close() {
	sm.redisConn.Close()
}

func ConnectSessionDB(address, database string) *sessionManager {
	var err error
	sm = &sessionManager{}
	sm.redisConn, err = redis.DialURL("redis://" + address + "/" + database)
	if err != nil {
		logger.Panic(err)
	}

	logger.Infof("Successfully connected to %v, database %v", address, database)

	return sm
}

func Create(sID string, uID uint) (bool, error) {
	res, err := sm.redisConn.Do("SET", sID, uID, "NX", "EX", 30*24*60*60)
	if err != nil {
		return false, err
	}
	if res != "OK" {
		logger.Infow("collision, session not created",
			"sID", sID,
			"uID", uID,
		)
		return false, nil
	}

	logger.Infow("session created",
		"sID", sID,
		"uID", uID,
	)

	return true, nil
}

func Get(sID string) (uint, error) {
	res, err := redis.Uint64(sm.redisConn.Do("GET", sID))
	if err != nil {
		if err == redis.ErrNil {
			return 0, ErrKeyNotFound
		}
		return 0, err
	}

	return uint(res), nil
}

func Delete(sID string) error {
	_, err := redis.Int(sm.redisConn.Do("DEL", sID))
	if err != nil {
		return err
	}

	return nil
}
