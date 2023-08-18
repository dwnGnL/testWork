package public

import (
	"net/http"

	"github.com/dwnGnL/testWork/internal/application"
	"github.com/dwnGnL/testWork/internal/config"
	"github.com/dwnGnL/testWork/internal/handler/models"
	"github.com/dwnGnL/testWork/lib/goerrors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func newHandler(cfg *config.Config) *Handler {
	return &Handler{}
}

func GenRouting(r *gin.RouterGroup, cfg *config.Config) {
	handler := newHandler(cfg)
	r.POST("/login", handler.Login)
}

func (h Handler) Login(c *gin.Context) {
	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	var req models.UserLogin
	err = c.ShouldBindJSON(&req)
	if err != nil {
		goerrors.Log().WithError(err).Warnln("req unmarshal err")
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	auth := app.GetAuth()
	token, err := auth.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		goerrors.Log().WithError(err).Warnln("login err")
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
