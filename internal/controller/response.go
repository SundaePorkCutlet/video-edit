package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	Success = 0
)

func ResponseSuccess(ginCtx *gin.Context) {
	ginCtx.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

func ResponseFailure(ginCtx *gin.Context, errorCode int) {
	ginCtx.JSON(http.StatusOK, gin.H{
		"code": errorCode,
	})
}

func ResponseWithResult(ginCtx *gin.Context, code int) {
	ginCtx.JSON(http.StatusOK, gin.H{
		"code": code,
	})
}

func ResponseWithData(ginCtx *gin.Context, data interface{}) {
	ginCtx.JSON(http.StatusOK, data)
}
