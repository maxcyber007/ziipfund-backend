package user

import (
	"mime/multipart"
	"net/http"
	"time"
	"ziipfund/jwt-api/orm"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecret []byte

// Biding from Register JSON
type UpdateBody struct {
	Id        string `json:"id" binding:"required"`
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Address   string `json:"address" binding:"required"`
	City      string `json:"city" binding:"required"`
	Province  string `json:"province" binding:"required"`
	Zipcode   string `json:"zipcode" binding:"required"`
	Country   string `json:"country" binding:"required"`
	//Avatar    string `json:"avatar" binding:"required"`
}

type FormData struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type ChangePasswordBody struct {
	Id       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Biding from Register JSON
type DepositBody struct {
	User_id string `json:"user_id" binding:"required"`
	Money   string `json:"money" binding:"required"`
}

func ReadAll(c *gin.Context) {
	var users []orm.Tbl_user
	orm.Db.Find(&users)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "users": users})
}

func Profile(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var users orm.Tbl_user
	orm.Db.First(&users, userId)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "users": users})
}

func GetDeposits(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var deposits []orm.Tbl_deposit
	orm.Db.Find(&deposits, "user_id = ?", userId)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Orders Read Success", "deposits": deposits})
}

func GetDeposit(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var deposits []orm.Tbl_deposit
	orm.Db.Find(&deposits, "user_id = ? AND status = ?", userId, "Active")
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Orders Read Success", "deposits": deposits})
}

func PostDeposit(c *gin.Context) {
	var json DepositBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentTime := time.Now()
	tbl_deposit := orm.Tbl_deposit{User_id: json.User_id, Money: json.Money, Status: "Active", Created_at: currentTime.Format("2006.01.02 15:04:05")}
	orm.Db.Create(&tbl_deposit)
	if tbl_deposit.Id > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Created", "tbl_fund": tbl_deposit})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Failed", "tbl_fund": tbl_deposit})
	}

}

func UpdateProfile(c *gin.Context) {

	var users UpdateBody

	var UpdateUsers orm.Tbl_user
	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Update Profile
	orm.Db.First(&UpdateUsers, users.Id)
	UpdateUsers.Firstname = users.Firstname
	UpdateUsers.Lastname = users.Lastname
	UpdateUsers.Email = users.Email
	UpdateUsers.Phone = users.Phone
	UpdateUsers.Address = users.Address
	UpdateUsers.City = users.City
	UpdateUsers.Province = users.Province
	UpdateUsers.Zipcode = users.Zipcode
	UpdateUsers.Country = users.Country
	//UpdateUsers.Avatar = users.Avatar
	orm.Db.Save(&UpdateUsers)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Updated"})

}

func UpdateAvatar(c *gin.Context) {

	var formData FormData
	// Bind the form data to the struct
	if err := c.ShouldBindWith(&formData, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Save the uploaded file to disk
	if err := c.SaveUploadedFile(formData.File, "uploads/avatar/"+formData.File.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return the path to the uploaded file
	c.JSON(http.StatusOK, gin.H{"path": "uploads/avatar/" + formData.File.Filename})

	//Update avatar
	var users UpdateBody
	var UpdateUsers orm.Tbl_user

	//Update Profile
	orm.Db.First(&UpdateUsers, users.Id)

	UpdateUsers.Avatar = formData.File.Filename
	orm.Db.Save(&UpdateUsers)

}

func ChangePassword(c *gin.Context) {

	var users ChangePasswordBody

	var UpdateUsers orm.Tbl_user
	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Update Profile
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(users.Password), 10)
	orm.Db.First(&UpdateUsers, users.Id)
	UpdateUsers.Password = string(encryptedPassword)
	orm.Db.Save(&UpdateUsers)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Password Updated"})

}
