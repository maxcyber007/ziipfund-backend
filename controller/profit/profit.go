package profit

import (
	"net/http"
	"time"
	"ziipfund/jwt-api/orm"

	"github.com/gin-gonic/gin"
)

var hmacSampleSecret []byte
var currentTime = time.Now()

// Biding from Profit JSON
type ProfitBody struct {
	//Id     string `json:"id" binding:"required"`
	Profit string `json:"profit" binding:"required"`
}

func PostProfit(c *gin.Context) {

	var json ProfitBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Create Profit
	tbl_profit := orm.Tbl_profit{Profit: json.Profit, Created_at: currentTime.Format("2006-01-02 15:04:05")}
	orm.Db.Create(&tbl_profit)
	if tbl_profit.Id > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Profit Created", "Profit : ": json.Profit})

	} else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Profit Failed"})
	}

}

func GetProfitAll(c *gin.Context) {
	var profits []orm.Tbl_profit
	orm.Db.Order("id desc").Find(&profits)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Profits Read Success", "profits": profits})
}

func GetProfitById(c *gin.Context) {
	id := c.Param("id")
	var profits []orm.Tbl_profit
	orm.Db.First(&profits, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Profits By Id Read Success", "profits": profits})
}

func GetProfitToday(c *gin.Context) {
	var profits []orm.Tbl_profit
	orm.Db.Last(&profits)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Profits Read Success", "profits": profits})
}

func DeleteProfit(c *gin.Context) {
	id := c.Param("id")
	var profits []orm.Tbl_profit
	orm.Db.Delete(&profits, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Delete Profit Success"})
}
