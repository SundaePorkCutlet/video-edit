package controller

import (
	"net/http"
	"time"

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

func ResponseSignedUrl(ginCtx *gin.Context, s3ObjectKey string, signedUrl string) {
	ginCtx.JSON(http.StatusOK, gin.H{
		"code":      Success,
		"timestamp": time.Now(),
		"objectKey": s3ObjectKey,
		"url":       signedUrl,
	})
}

func LoginResponse(ginCtx *gin.Context, code int, date time.Time, token string) {
	role := ginCtx.GetInt("role")
	ginCtx.JSON(code, gin.H{
		"code":  code,
		"date":  date,
		"token": token,
		"role":  role,
	})

}
