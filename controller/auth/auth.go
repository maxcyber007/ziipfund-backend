package auth

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
	"ziipfund/jwt-api/orm"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecret []byte

/*
type info struct {
	Name string
}

func (i info) sendMail() {

	t := template.New("template.html")

	var err error
	t, err = t.ParseFiles("template.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, i); err != nil {
		log.Println(err)
	}

	result := tpl.String()
	m := gomail.NewMessage()
	m.SetHeader("From", "ziipfund@gmail.com")
	m.SetHeader("To", "rungsiyanon@gmail.com")
	m.SetHeader("Subject", "[ZiipFund] Completing Your Registration")
	m.SetBody("text/html", result)
	m.Attach("template.html") // attach whatever you want

	d := gomail.NewDialer("smtp.gmail.com", 587, "ziipfund@gmail.com", "doadmcllvyvbodrn")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
*/

// Biding from Register JSON
type RegisterBody struct {
	//Id     string  `json:"id" binding:"required"`
	Ref_code       string `json:"ref_code" binding:"required"`
	Child_ref_code string `json:"child_ref_code" binding:"required"`
	Firstname      string `json:"firstname" binding:"required"`
	Lastname       string `json:"lastname" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Phone          string `json:"phone" binding:"required"`
}

type FormData struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func Register(c *gin.Context) {

	var json RegisterBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Check User Exists
	var userExist orm.Tbl_user
	orm.Db.Where("ref_code = ?", json.Ref_code).First(&userExist)
	if userExist.Id > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Does Not Exists"})
		return
	}

	//Create User
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	tbl_user := orm.Tbl_user{Ref_code: json.Ref_code, Child_ref_code: json.Child_ref_code, Firstname: json.Firstname, Lastname: json.Lastname, Email: json.Email, Password: string(encryptedPassword), Phone: json.Phone}
	orm.Db.Create(&tbl_user)
	if tbl_user.Id > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Created", "userID": tbl_user.Id, "Child_ref_code": json.Child_ref_code})
		/*
			//Sent Mail
			//"gopkg.in/gomail.v2"
			mail := gomail.NewMessage()
			mail.SetHeader("From", "ziipfund@gmail.com")
			mail.SetHeader("To", json.Email)
			mail.SetHeader("Subject", "[ZiipFund] Completing Your Registration")
			// mail.SetBody("text/plain", "Username : "+json.Email+" "+"Password:"+json.Password)
			mail.SetBody("text/html",
				"<head><title>"+
					"<link rel=\"preconnect\" href=\"https://fonts.googleapis.com\">"+
					"<link rel=\"preconnect\" href=\"https://fonts.gstatic.com\" crossorigin>"+
					"<link href=\"https://fonts.googleapis.com/css2?family=Prompt:wght@200&display=swap\" rel=\"stylesheet\"></title></head>"+
					"<div align=\"center\">"+
					"<table width=\"100%\">"+
					"<tr bgcolor=\"#000000\">"+
					"<td align=\"left\"><img src=\"https://ziipfund.com/_next/image?url=%2Fassets%2Fimg%2Ffavicon.png&w=96&q=75\" width=\"50\" heist=\"50\"></td>"+
					"<tr>"+
					"<td style=\"font-family: 'Prompt', sans-serif;\"><b>Dear</b> "+json.Firstname+"  "+json.Lastname+"</td>"+
					"<tr>"+
					"<td>Thank you for completing your registration with ziipfund company limited.</td>"+
					"<tr>"+
					"<td>This email serves as a confirmation that your account is activated and that you are officially a part of the ziipfund company limited family. Enjoy!.</td>"+
					"<tr>"+
					"<td><tr><tr></td>"+
					"<tr>"+
					"<td><b>Username :</b> "+json.Email+"</td>"+
					"<tr>"+
					"<td><b>Password :</b> "+json.Password+"</td>"+
					"<tr>"+
					"<td><hr></td>"+
					"<tr>"+
					"<td><b>Risk warning : </b>Investing in investment units is not a deposit. And there is a risk of investment, investors may receive a return of investment more or less than the initial investment. Investors should invest in mutual funds that are suitable for their investment objectives and the risks that may arise from such investments are acceptable to investors.</td>"+
					"<tr><tr><tr>"+
					"<td><b>Regards</b></td>"+
					"<tr>"+
					"<td><b>The ZiipFund Company Team Services.</b></td>"+
					"<tr>"+
					"<td align=\"center\"><b>© 2023 Copyright www.ziipfund.com All Rights Reserved.</b></td>"+
					"</table>"+
					"</div>")
			dialer := gomail.NewPlainDialer("smtp.gmail.com", 587, "ziipfund@gmail.com", "doadmcllvyvbodrn")
			//dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
			if err := dialer.DialAndSend(mail); err != nil {
				fmt.Print(err)
				panic(err)
			}
		*/
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
	var userExist orm.Tbl_user
	orm.Db.Where("email = ?", json.Email).First(&userExist)
	if userExist.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Does Not Exists"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(json.Password))
	if err == nil {

		hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"Ref_code":   userExist.Ref_code,
			"userId":     userExist.Id,
			"Firstname":  userExist.Firstname,
			"Lastname":   userExist.Lastname,
			"Avatar":     userExist.Avatar,
			"Privileges": userExist.Privileges,
			"exp":        time.Now().Add(time.Minute * 24).Unix(),
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
			c.Set("Level", claims["Level"])
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Token Success", "Ref_code": claims["Ref_code"], "userId": claims["userId"], "Firstname": claims["Firstname"], "Lastname": claims["Lastname"], "Avatar": claims["Avatar"], "Privileges": claims["Privileges"]})
		}
	} else {
		log.Println("Token is invalid:", err)
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Token Failed"})
	}

}

func GetTeamWork(c *gin.Context) {
	//id := c.Param("id")
	var users []orm.Tbl_user
	//orm.Db.Find(&users, "child_ref_code = ?", id)
	orm.Db.Find(&users)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Teamwork Read Success", "teams": users})
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

func Upload(c *gin.Context) {
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

}
