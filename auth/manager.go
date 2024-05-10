package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AuthManager struct {
	loginManager    LoginDatabaseManager
	sessionsManager SessionManager
}

func (manager *AuthManager) InitializeSession(ctx *gin.Context) {
	sessionID := ctx.Query("session_id")
	var err error
	if sessionID, err = manager.sessionsManager.RenewSession(sessionID); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, sessionID)
}

func (manager *AuthManager) LogIn(ctx *gin.Context) {
	login := ctx.Query("login")
	password := ctx.Query("password")
	sessionID := ctx.Query("session_id")
	isShop, err := strconv.ParseBool(ctx.Query("is_shop"))
	if err != nil {
		fmt.Printf("log in with invalid is_shop parameter: \"%s\" defaulting to false\n", ctx.Query("is_shop"))
	}
	if sessionID, err = manager.sessionsManager.RenewSession(sessionID); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	var id int64
	if !isShop {
		id, err = manager.loginManager.LoginUser(login, password)
	} else {
		id, err = manager.loginManager.LoginShop(login, password)
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	if !isShop {
		sessionID, err = manager.sessionsManager.SetUserID(sessionID, id)
	} else {
		sessionID, err = manager.sessionsManager.SetShopID(sessionID, id)
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		ctx.JSON(http.StatusOK, sessionID)
	}
}

func (manager *AuthManager) SingUp(ctx *gin.Context) {
	login := ctx.Query("login")
	password := ctx.Query("password")
	sessionID := ctx.Query("session_id")
	isShop, err := strconv.ParseBool(ctx.Query("is_shop"))
	if err != nil {
		fmt.Printf("log in with invalid is_shop parameter: \"%s\" defaulting to false\n", ctx.Query("is_shop"))
	}
	if sessionID, err = manager.sessionsManager.RenewSession(sessionID); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	var id int64
	if !isShop {
		id, err = manager.loginManager.CreateUser(login, password)
	} else {
		id, err = manager.loginManager.CreateShop(login, password)
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	if !isShop {
		sessionID, err = manager.sessionsManager.SetUserID(sessionID, id)
	} else {
		sessionID, err = manager.sessionsManager.SetShopID(sessionID, id)
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		ctx.JSON(http.StatusOK, sessionID)
	}
}

func (manager *AuthManager) UserID(ctx *gin.Context) {
	sessionID := ctx.Query("session_id")
	var err error
	if sessionID, err = manager.sessionsManager.RenewSession(sessionID); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	id, err := manager.sessionsManager.GetUserID(sessionID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		ctx.JSON(http.StatusOK, id)
	}
}

func (manager *AuthManager) ShopID(ctx *gin.Context) {
	sessionID := ctx.Query("session_id")
	var err error
	if sessionID, err = manager.sessionsManager.RenewSession(sessionID); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	id, err := manager.sessionsManager.GetShopID(sessionID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		ctx.JSON(http.StatusOK, id)
	}
}

func (manager *AuthManager) Close() {
	check(manager.loginManager.Close())
	check(manager.sessionsManager.Close())
}
