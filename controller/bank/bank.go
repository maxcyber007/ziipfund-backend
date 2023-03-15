package bank

import (
	_ "encoding/json"
	_ "fmt"
	"mime/multipart"
	"net/http"
	"time"
	"ziipfund/jwt-api/orm"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var hmacSampleSecret []byte
var currentTime = time.Now()

// Biding from Register JSON
type BankBody struct {
	User_id   string `json:"user_id" binding:"required"`
	Bank_file string `json:"bank_file" binding:"required"`
}

type FormData struct {
	User_id string                `form:"user_id" binding:"required"`
	File    *multipart.FileHeader `form:"file" binding:"required"`
}

func GetBy(c *gin.Context) {
	//id := c.Param("id")
	userId := c.MustGet("userId").(float64)
	var banks []orm.Tbl_bank
	orm.Db.First(&banks, "user_id = ?", userId)
	//orm.Db.Find(&banks)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Bank Read Success", "banks": banks})
}

func Upload(c *gin.Context) {

	var formData FormData

	//Bind the form data to the struct
	if err := c.ShouldBindWith(&formData, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the uploaded file to disk
	if err := c.SaveUploadedFile(formData.File, "uploads/bank/"+formData.User_id+"_"+formData.File.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return the path to the uploaded file
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Created", "path": "uploads/bank/" + formData.File.Filename})
}

func Post(c *gin.Context) {

	var json BankBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentTime := time.Now()
	tbl_bank := orm.Tbl_bank{User_id: json.User_id, Bank_file: json.User_id + "_" + json.Bank_file, Created_at: currentTime.Format("2006.01.02 15:04:05")}
	orm.Db.Create(&tbl_bank)
	if tbl_bank.Id > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Created", "tbl_bank": tbl_bank})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Failed", "tbl_bank": tbl_bank})
	}
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	var banks []orm.Tbl_bank
	orm.Db.Delete(&banks, "user_id = ?", id)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Bank Delete Success", "banks": banks})
}
