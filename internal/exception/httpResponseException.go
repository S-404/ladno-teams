package exception

import (
	"log"

	"github.com/gin-gonic/gin"
)

func HttpResponseException(c *gin.Context, apiError *ApiError) {
	log.Printf("HttpResponseException guid:%s, code:%s, status:%d, message:%s", apiError.Guid, apiError.Code, apiError.Status, apiError.Message)
	c.AbortWithStatusJSON(apiError.Status, gin.H{"error": apiError})
}
