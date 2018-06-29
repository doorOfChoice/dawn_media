package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Page404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "common/404", h(gin.H{}, c))
}
