package main

import (
	"fmt"
	AuthController "ziipfund/jwt-api/controller/auth"
	BankController "ziipfund/jwt-api/controller/bank"
	CommissionController "ziipfund/jwt-api/controller/commission"
	DocumentController "ziipfund/jwt-api/controller/doc"
	ProfitController "ziipfund/jwt-api/controller/profit"
	UserController "ziipfund/jwt-api/controller/user"
	"ziipfund/jwt-api/middleware"
	"ziipfund/jwt-api/orm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	orm.InitDB()

	r := gin.Default()
	//r.Use(cors.Default())
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("authorization")
	r.Use(cors.New(config))

	r.POST("/register", AuthController.Register)
	r.POST("/login", AuthController.Login)
	r.POST("/auth", AuthController.Auth)
	r.POST("/upload", AuthController.Upload)
	r.Static("/uploads", "./uploads")
	r.GET("/teamwork/:id", AuthController.GetTeamWork)

	authorized := r.Group("/users", middleware.JwtAuthen())
	authorized.GET("/readall", UserController.ReadAll)
	authorized.GET("/profile", UserController.Profile)
	authorized.PUT("/updateprofile", UserController.UpdateProfile)
	authorized.PUT("/updateavatar", UserController.UpdateAvatar)
	authorized.PUT("/changepassword", UserController.ChangePassword)

	authorized.GET("/getdeposits", UserController.GetDeposits)
	authorized.GET("/getdeposit", UserController.GetDeposit)
	authorized.POST("/postdeposit", UserController.PostDeposit)
	// authorized.GET("/teamwork/:id", UserController.GetTeamWork)

	authorized.POST("/postprofit", ProfitController.PostProfit)
	authorized.GET("/profit/:id", ProfitController.GetProfitById)
	authorized.GET("/profitall", ProfitController.GetProfitAll)
	authorized.GET("/profittoday", ProfitController.GetProfitToday)
	authorized.DELETE("/deleteprofit/:id", ProfitController.DeleteProfit)

	//Bank Document Upload
	authorized.GET("/bank/getby", BankController.GetBy)
	authorized.POST("/bank/post", BankController.Post)
	authorized.POST("/bank/upload", BankController.Upload)
	authorized.DELETE("/bank/delete/:id", BankController.Delete)

	//Photo ID Document Upload
	authorized.GET("/document/getby", DocumentController.GetBy)
	authorized.POST("/document/post", DocumentController.Post)
	authorized.POST("/document/upload", DocumentController.Upload)
	authorized.DELETE("/document/delete/:id", DocumentController.Delete)

	//Commission
	authorized.GET("/commission/getall", CommissionController.GetAll)

	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
