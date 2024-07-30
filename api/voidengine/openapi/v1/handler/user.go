package handler

import (
	"BlackHole/internal/voidengine/controller"
	"BlackHole/pkg/env"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListUser(c *gin.Context, e *env.Env) {
	user := controller.NewUser()
	res := user.ListUser(c, e)
	c.JSON(http.StatusOK, res)
}

func AddUer(c *gin.Context, e *env.Env) {
	user := controller.NewUser()
	res := user.AddUser(c, e)
	c.JSON(http.StatusOK, res)
}

func ModifyUer(c *gin.Context, e *env.Env) {
	user := controller.NewUser()
	res := user.ModifyUser(c, e)
	c.JSON(http.StatusOK, res)
}

func DeleteUer(c *gin.Context, e *env.Env) {
	user := controller.NewUser()
	res := user.DeleteUser(c, e)
	c.JSON(http.StatusOK, res)
}
