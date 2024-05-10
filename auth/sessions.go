package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/hazelcast/hazelcast-go-client"
	"time"
)

type SessionManager struct {
	hzCTX    context.Context
	hzClient *hazelcast.Client
	hzMap    *hazelcast.Map
}

type session struct {
	ExpirationTime time.Time `json:"expiration_time,omitempty"`
	IsShopAccount  bool      `json:"is_shop_account,omitempty"`
	AccountID      int64     `json:"account_id,omitempty"`
}

const UnLogged = -1

func (sm *SessionManager) Initialize(clusterName, mapName string) (err error) {
	sm.hzCTX = context.Background()
	cfg := hazelcast.Config{}
	cfg.Cluster.Name = clusterName
	sm.hzClient, err = hazelcast.StartNewClientWithConfig(sm.hzCTX, cfg)
	if err != nil {
		return err
	}
	sm.hzMap, err = sm.hzClient.GetMap(sm.hzCTX, mapName)
	if err != nil {
		return err
	}
	return nil
}
func (sm *SessionManager) Close() error {
	return sm.hzClient.Shutdown(sm.hzCTX)
}

func (sm *SessionManager) newSession() (string, error) {
	sessionUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	sessionID := sessionUUID.String()
	expirationTime := time.Now().Add(5 * time.Minute)
	value, err := json.Marshal(session{ExpirationTime: expirationTime, AccountID: UnLogged})
	if err != nil {
		return "", err
	}
	err = sm.hzMap.Set(sm.hzCTX, sessionID, value)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (sm *SessionManager) getSession(sessionID string) (*session, error) {
	value, err := sm.hzMap.Get(sm.hzCTX, sessionID)
	if err != nil || value == nil {
		return nil, err
	}
	var sessionObj session
	err = json.Unmarshal(value.([]byte), &sessionObj)
	if err != nil {
		return nil, err
	}
	if time.Now().After(sessionObj.ExpirationTime) {
		return nil, nil
	}
	return &sessionObj, nil
}

func (sm *SessionManager) RenewSession(sessionID string) (string, error) {
	sessionObj, err := sm.getSession(sessionID)
	if err != nil {
		return "", err
	}
	if sessionObj == nil {
		return sm.newSession()
	}
	sessionObj.ExpirationTime = time.Now().Add(5 * time.Minute)
	value, err := json.Marshal(sessionObj)
	if err != nil {
		return "", err
	}
	err = sm.hzMap.Set(sm.hzCTX, sessionID, value)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (sm *SessionManager) GetUserID(sessionID string) (int64, error) {
	sessionObj, err := sm.getSession(sessionID)
	if err != nil || sessionObj == nil {
		return UnLogged, err
	}
	if !sessionObj.IsShopAccount {
		return sessionObj.AccountID, nil
	}
	return UnLogged, nil
}
func (sm *SessionManager) GetShopID(sessionID string) (int64, error) {
	sessionObj, err := sm.getSession(sessionID)
	if err != nil || sessionObj == nil {
		return UnLogged, err
	}
	if sessionObj.IsShopAccount {
		return sessionObj.AccountID, nil
	}
	return UnLogged, nil
}

func (sm *SessionManager) SetUserID(sessionID string, userID int64) (string, error) {
	sessionObj, err := sm.getSession(sessionID)
	if err != nil {
		return "", err
	}
	if sessionObj == nil {
		sessionID, err = sm.newSession()
		if err != nil {
			return "", err
		}
		sessionObj, err = sm.getSession(sessionID)
		if err != nil {
			return "", err
		}
		if sessionObj == nil {
			return "", errors.New("failed to create new session and log in")
		}
	}
	sessionObj.IsShopAccount = false
	sessionObj.AccountID = userID
	value, err := json.Marshal(sessionObj)
	if err != nil {
		return "", err
	}
	err = sm.hzMap.Set(sm.hzCTX, sessionID, value)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}
func (sm *SessionManager) SetShopID(sessionID string, shopID int64) (string, error) {
	sessionObj, err := sm.getSession(sessionID)
	if err != nil {
		return "", err
	}
	if sessionObj == nil {
		sessionID, err = sm.newSession()
		if err != nil {
			return "", err
		}
		sessionObj, err = sm.getSession(sessionID)
		if err != nil {
			return "", err
		}
		if sessionObj == nil {
			return "", errors.New("failed to create new session and log in")
		}
	}
	sessionObj.IsShopAccount = true
	sessionObj.AccountID = shopID
	value, err := json.Marshal(sessionObj)
	if err != nil {
		return "", err
	}
	err = sm.hzMap.Set(sm.hzCTX, sessionID, value)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}
