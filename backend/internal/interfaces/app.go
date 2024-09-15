package interfaces

import "github.com/gin-gonic/gin"

type App interface {
	Init(r *gin.Engine)
}
