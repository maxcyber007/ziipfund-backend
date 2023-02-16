package user

import (
	"net/http"
	  "github.com/gin-gonic/gin"
	  "ziipfund/jwt-api/orm"
  )

  var hmacSampleSecret []byte

func ReadAll(c *gin.Context){
var users []orm.Tbl_member
orm.Db.Find(&users)
c.JSON(http.StatusOK, gin.H{"staus": "ok", "message": "User Read Success", "users": users})
}

func Profile(c *gin.Context){

	userId := c.MustGet("userId").(float64)
	var users orm.Tbl_member
	orm.Db.First(&users, userId)
	c.JSON(http.StatusOK, gin.H{"staus": "ok", "message": "User Read Success", "users": users})
}