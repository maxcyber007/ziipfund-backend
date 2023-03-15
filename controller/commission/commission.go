package commission

import (
	_ "encoding/json"
	_ "fmt"
	"net/http"
	"time"
	"ziipfund/jwt-api/orm"

	"github.com/gin-gonic/gin"
)

var hmacSampleSecret []byte
var currentTime = time.Now()

// Biding from Register JSON
type CommissionBody struct {
	User_id   string `json:"user_id" binding:"required"`
	Bank_file string `json:"bank_file" binding:"required"`
}

func GetAll(c *gin.Context) {
	//id := c.Param("id")
	var users []orm.Tbl_user
	//orm.Db.Find(&users, "child_ref_code = ?", id)
	orm.Db.Find(&users)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "users": users})
}
