package main
import (
  "fmt"
  AuthController "ziipfund/jwt-api/controller/auth"
  UserController "ziipfund/jwt-api/controller/user"
  "ziipfund/jwt-api/orm"
  "ziipfund/jwt-api/middleware"
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/cors"
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

  authorized := r.Group("/users", middleware.JwtAuthen())
  authorized.GET("/readall", UserController.ReadAll)
  authorized.GET("/profile", UserController.Profile)

  r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}