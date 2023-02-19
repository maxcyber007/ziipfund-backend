package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"ziipfund/jwt-api/orm"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecret []byte

// Biding from Register JSON
type RegisterBody struct {
	//Id     string  `json:"id" binding:"required"`
	Ref_code       string `json:"ref_code" binding:"required"`
	Child_ref_code string `json:"child_ref_code" binding:"required"`
	Firstname      string `json:"firstname" binding:"required"`
	Lastname       string `json:"lastname" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {

	var json RegisterBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Check User Exists
	var userExist orm.Tbl_member
	orm.Db.Where("ref_code = ?", json.Ref_code).First(&userExist)
	if userExist.Id > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Does Not Exists"})
		return
	}

	//Create User
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	tbl_member := orm.Tbl_member{Ref_code: json.Ref_code, Child_ref_code: json.Child_ref_code, Firstname: json.Firstname, Lastname: json.Lastname, Email: json.Email, Password: string(encryptedPassword)}
	orm.Db.Create(&tbl_member)
	if tbl_member.Id > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Created", "userID": tbl_member.Id, "Child_ref_code": json.Child_ref_code})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Failed", "Child_ref_code": json.Child_ref_code})
	}

}

// Biding from Login JSON
type LoginBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login Function
func Login(c *gin.Context) {

	var json LoginBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Check User Exists
	var userExist orm.Tbl_member
	orm.Db.Where("email = ?", json.Email).First(&userExist)
	if userExist.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Does Not Exists"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(json.Password))
	if err == nil {

		hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"Ref_code":  userExist.Ref_code,
			"userId":    userExist.Id,
			"Firstname": userExist.Firstname,
			"Lastname":  userExist.Lastname,
			"exp":       time.Now().Add(time.Minute * 24).Unix(),
		})
		tokenString, err := token.SignedString(hmacSampleSecret)
		fmt.Println(tokenString, err)

		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Login Success", "token": tokenString})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Login Failed"})
	}
}

func Auth(c *gin.Context) {
	hmacSampleSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
	header := c.Request.Header.Get("Authorization")
	tokenString := strings.Replace(header, "Bearer ", "", 1)

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method used in the token is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSampleSecret, nil
	})

	// Check if the token is valid
	if err == nil && token.Valid {
		log.Println("Token is Success")
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("Ref_code", claims["Ref_code"])
			c.Set("userId", claims["userId"])
			c.Set("Firstname", claims["Firstname"])
			c.Set("Lastname", claims["Lastname"])
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Token Success", "Ref_code": claims["Ref_code"], "userId": claims["userId"], "Firstname": claims["Firstname"], "Lastname": claims["Lastname"]})
		}
	} else {
		log.Println("Token is invalid:", err)
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Token Failed"})
	}

}

// func Auth(c *gin.Context){

//   hmacSampleSecret := []byte(os.Getenv("JWT_SECRET_KEY"))

//   header := c.Request.Header.Get("Authorization")
//   tokenString := strings.Replace(header, "Bearer ", "", 1)

//   token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return hmacSampleSecret, nil
// 	})

//   if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//     c.Set("userId", claims["userId"])
//     c.Set("Firstname", claims["Firstname"])
//     c.Set("Lastname", claims["Lastname"])
//     c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Token Success", "userId": claims["userId"], "Firstname": claims["Firstname"], "Lastname": claims["Lastname"]})
// } else {
//   c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "forbidden", "message": err.Error()})
// }
// c.Next()

// }
