package doc

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
type DocumentBody struct {
	User_id       string `json:"user_id" binding:"required"`
	Document_file string `json:"document_file" binding:"required"`
}

type FormData struct {
	User_id string                `form:"user_id" binding:"required"`
	File    *multipart.FileHeader `form:"file" binding:"required"`
}

func GetBy(c *gin.Context) {
	//id := c.Param("id")
	userId := c.MustGet("userId").(float64)
	var docs []orm.Tbl_document
	orm.Db.First(&docs, "user_id = ?", userId)
	//orm.Db.Find(&banks)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Bank Read Success", "docs": docs})
}

func Upload(c *gin.Context) {

	var formData FormData

	//Bind the form data to the struct
	if err := c.ShouldBindWith(&formData, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the uploaded file to disk
	if err := c.SaveUploadedFile(formData.File, "uploads/document/"+formData.User_id+"_"+formData.File.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return the path to the uploaded file
	//c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Created", "path": "uploads/document/" + formData.File.Filename})
}

func Post(c *gin.Context) {

	var json DocumentBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentTime := time.Now()
	tbl_document := orm.Tbl_document{User_id: json.User_id, Document_file: json.User_id + "_" + json.Document_file, Created_at: currentTime.Format("2006.01.02 15:04:05")}
	orm.Db.Create(&tbl_document)
	if tbl_document.Id > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Created", "tbl_document": tbl_document})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Failed", "tbl_document": tbl_document})
	}
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	var docs []orm.Tbl_document
	orm.Db.Delete(&docs, "user_id = ?", id)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Document Delete Success", "docs": docs})
}
