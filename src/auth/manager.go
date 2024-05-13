package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthManager struct {
	loginManager    LoginDatabaseManager
	sessionsManager SessionManager
	//logger          Logger
}

func (manager *AuthManager) setupSession(ctx *gin.Context) string {
	var err *HttpError
	sessionID, noCookie := ctx.Cookie("session_id")
	if noCookie != nil {
		sessionID, err = manager.sessionsManager.newSession()
		//manager.logger.Log(fmt.Sprintf("Created new session: %s", sessionID))
	} else {
		sessionID, err = manager.sessionsManager.RenewSession(sessionID)
	}
	if err != nil {
		ctx.String(err.Code(), err.Error())
		//manager.logger.Log(fmt.Sprintf("Failed to setup session: %s", err.Error()))
		return ""
	}
	ctx.SetCookie("session_id", sessionID, 3600, "/", "", false, true)
	return sessionID
}

func (manager *AuthManager) InitializeSession(ctx *gin.Context) {
	sessionID := manager.setupSession(ctx)
	if sessionID == "" {
		return
	}

	ctx.Status(http.StatusOK)
}

func (manager *AuthManager) LogIn(ctx *gin.Context) {
	sessionID := manager.setupSession(ctx)
	if sessionID == "" {
		return
	}

	login := ctx.Query("login")
	password := ctx.Query("password")

	id, userType, err := manager.loginManager.loginAccount(login, password)
	if err != nil {
		//manager.logger.Log(fmt.Sprintf("Failed to log in user %s: %s", login, err.Error()))
		ctx.String(err.Code(), err.Error())
		return
	}
	sessionID, err = manager.sessionsManager.SetID(sessionID, id, userType)
	if err != nil {
		//manager.logger.Log(fmt.Sprintf("Failed to set user id: %s", err.Error()))
		ctx.String(err.Code(), err.Error())
	} else {
		//manager.logger.Log(fmt.Sprintf("Successfully logged in user %s for sesion %s", login, sessionID))
		ctx.Status(http.StatusOK)
	}
}

func (manager *AuthManager) SingUp(ctx *gin.Context) {
	sessionID := manager.setupSession(ctx)
	if sessionID == "" {
		return
	}

	login := ctx.Query("login")
	password := ctx.Query("password")
	userType := UserType(ctx.Query("user_type"))

	if !manager.loginManager.CheckUserType(userType) {
		ctx.String(http.StatusBadRequest, "Invalid user type: %s", userType)
		//manager.logger.Log(fmt.Sprintf("Received invalid user type: %s", userType))
		return
	}

	id, err := manager.loginManager.CreateAccount(login, password, userType)
	if err != nil {
		//manager.logger.Log(fmt.Sprintf("Failed to create an account: %s", err.Error()))
		ctx.String(err.Code(), err.Error())
		return
	}
	sessionID, err = manager.sessionsManager.SetID(sessionID, id, userType)
	if err != nil {
		//manager.logger.Log(fmt.Sprintf("Failed to set user id: %s", err.Error()))
		ctx.String(err.Code(), err.Error())
	} else {
		//manager.logger.Log(fmt.Sprintf("Successfully created user %s for sesion %s", login, sessionID))
		ctx.Status(http.StatusOK)
	}
}

func (manager *AuthManager) GetID(ctx *gin.Context) {
	sessionID := manager.setupSession(ctx)
	if sessionID == "" {
		return
	}

	id, err := manager.sessionsManager.GetID(sessionID)
	if err != nil {
		//manager.logger.Log(fmt.Sprintf("Failed to get user id: %s", err.Error()))
		ctx.String(err.Code(), err.Error())
	} else {
		//manager.logger.Log(fmt.Sprintf("Successfully retrieved user id %s for session %s", id, sessionID))
		ctx.String(http.StatusOK, id)
	}
}

func (manager *AuthManager) Close() {
	check(manager.loginManager.Close())
	check(manager.sessionsManager.Close())
}
