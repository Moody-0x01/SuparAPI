package routes;

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Moody0101-X/Go_Api/models"
	"github.com/Moody0101-X/Go_Api/database"
	"github.com/Moody0101-X/Go_Api/crypto"
)

func isEmpty(s string) bool { return len(s) == 0 }

func GetFieldFromContext(c *gin.Context, field string) string { return c.Query(field) }

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusOK, models.MakeServerResponse(100, "Not implemented!"))	
}
