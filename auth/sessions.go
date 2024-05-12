package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/hazelcast/hazelcast-go-client"
	"net/http"
	"strconv"
	"time"
)

type SessionManager struct {
	hzCTX    context.Context
	hzClient *hazelcast.Client
	hzMap    *hazelcast.Map
}

type session struct {
	ExpirationTime time.Time `json:"expiration_time,omitempty"`
	IsShopAccount  UserType  `json:"user_type,omitempty"`
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

func (sm *SessionManager) newSession() (string, *HttpError) {
	sessionUUID, err := uuid.NewUUID()
	if err != nil {
		return "", NewHttpError(err, "Unable to create UUID", http.StatusInternalServerError)
	}
	sessionID := sessionUUID.String()
	expirationTime := time.Now().Add(5 * time.Minute)
	value, err := json.Marshal(session{ExpirationTime: expirationTime, AccountID: UnLogged})
	if err != nil {
		return "", NewHttpError(err, "Unable to marshal session object into JSON", http.StatusInternalServerError)
	}
	err = sm.hzMap.Set(sm.hzCTX, sessionID, value)
	if err != nil {
		return "", NewHttpError(err, "Unable to save session into Hazelcast map", http.StatusInternalServerError)
	}
	return sessionID, nil
}

func (sm *SessionManager) getSession(sessionID string) (*session, *HttpError) {
	value, err := sm.hzMap.Get(sm.hzCTX, sessionID)
	if err != nil {
		return nil, NewHttpError(err, "Unable to retrieve value from Hazelcast map", http.StatusServiceUnavailable)
	} else if value == nil {
		return nil, nil
	}
	var sessionObj session
	err = json.Unmarshal(value.([]byte), &sessionObj)
	if err != nil {
		return nil, NewHttpError(err, "Unable to parse session object", http.StatusInternalServerError)
	}
	if time.Now().After(sessionObj.ExpirationTime) {
		return nil, nil
	}
	return &sessionObj, nil
}

func (sm *SessionManager) RenewSession(sessionID string) (string, *HttpError) {
	sessionObj, httpErr := sm.getSession(sessionID)
	if httpErr != nil {
		return "", httpErr
	}
	if sessionObj == nil {
		return sm.newSession()
	}
	sessionObj.ExpirationTime = time.Now().Add(5 * time.Minute)
	value, err := json.Marshal(sessionObj)
	if err != nil {
		return "", NewHttpError(err, "Unable to marshal session object into JSON", http.StatusInternalServerError)
	}
	err = sm.hzMap.Set(sm.hzCTX, sessionID, value)
	if err != nil {
		return "", NewHttpError(err, "Unable to save session into Hazelcast map", http.StatusInternalServerError)
	}
	return sessionID, nil
}

func (sm *SessionManager) GetID(sessionID string) (string, *HttpError) {
	sessionObj, err := sm.getSession(sessionID)
	if err != nil {
		return "", err
	}
	if sessionObj == nil || sessionObj.AccountID == UnLogged {
		return "", NewHttpError(errors.New(""), "Not authorised", http.StatusUnauthorized)
	}
	return string(sessionObj.IsShopAccount) + ":" + strconv.FormatInt(sessionObj.AccountID, 10), nil
}

func (sm *SessionManager) SetID(sessionID string, id int64, userType UserType) (string, *HttpError) {
	sessionObj, httpErr := sm.getSession(sessionID)
	if httpErr != nil {
		return "", httpErr
	}
	if sessionObj == nil {
		sessionID, httpErr = sm.newSession()
		if httpErr != nil {
			return "", httpErr
		}
		sessionObj, httpErr = sm.getSession(sessionID)
		if httpErr != nil {
			return "", httpErr
		}
		if sessionObj == nil {
			return "", NewHttpError(errors.New(""), "failed to create new session and log in", http.StatusInternalServerError)
		}
	}
	sessionObj.IsShopAccount = userType
	sessionObj.AccountID = id
	value, err := json.Marshal(sessionObj)
	if err != nil {
		return "", NewHttpError(err, "Unable to marshal session object into JSON", http.StatusInternalServerError)
	}
	err = sm.hzMap.Set(sm.hzCTX, sessionID, value)
	if err != nil {
		return "", NewHttpError(err, "Unable to save session into Hazelcast map", http.StatusInternalServerError)
	}
	return sessionID, nil
}
