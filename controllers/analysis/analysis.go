package analysis

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"io/ioutil"
	"fmt"
	"hzl.im/gin-platform/models"
	"hzl.im/gin-platform/services"
)

func ShowData(ctx *gin.Context) {
	currency := ctx.Param("currency")
	currency = currency[:3] + "_" + currency[3:]

	currencyData := []models.CurrencyData{}

	rows := services.DB.Table(currency).Order("id asc").Find(&currencyData)
	if rows.Error != nil {
		fmt.Println(rows.Error)
		ctx.JSON(http.StatusOK, "")
		return
	}

	var data []interface{}

	for _, v := range currencyData {
		obj := make(map[string]interface{})
		obj["date"] = v.Time.Format("2006-01-02")
		obj["Highest"] = v.High
		obj["Lowest"] = v.Low
		data = append(data, obj)

	}

	ctx.HTML(http.StatusOK, "analysis.tmpl", gin.H{
		"title": "Analysis Page",
		"data": data,
		"currency": currency,
	})

}

func GetData(ctx *gin.Context) {
	url := "http://apilayer.net/api/live?access_key=2d003e9906243170959cbab33604b029&format=1"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}