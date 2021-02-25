package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Demo struct {
}

func NewDemo() Demo {
	return Demo{}
}


func (d Demo) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "PONG")
}
