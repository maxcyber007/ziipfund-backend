package orm
import (
	"os"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
  )
  
var Db *gorm.DB
var err error

func InitDB() {
	dsn := os.Getenv("MYSQL_DNS")
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//Migrate the schema
	//db.AutoMigrate(&Tbl_member{})


}