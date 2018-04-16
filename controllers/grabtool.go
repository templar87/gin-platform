package controllers

import (
	"hzl.im/gin-platform/models"
	"hzl.im/gin-platform/services"
	"net/http"
	"io/ioutil"
	"time"
	"encoding/json"
	"fmt"
	"strings"
)

func GetValueFunc() {

	url := "http://apilayer.net/api/live?access_key=2d003e9906243170959cbab33604b029&format=1"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("cache-control", "no-cache")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)


	today := time.Now().UTC()

	var data map[string]interface{}
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	newData := data["quotes"].(map[string]interface{})

	currencyArr := []string{"USDHKD","USDEUR","USDGBP","USDJPY","USDAUD","USDNZD","USDCAD"}

	for _, currency := range currencyArr {
		currency_ := strings.ToLower(currency[:3] + "_" + currency[3:])
		todayValue, err := getValue(today.Format("2006-01-02"), currency_)
		fmt.Println(todayValue)
		fmt.Println(newData[currency])

		if err != nil {
			todayValue = models.CurrencyData{
				High: newData[currency].(float64),
				Low: newData[currency].(float64),
				Time: today,
				Name: currency_,
			}
			if err := services.DB.Create(&todayValue).Error; err != nil {
				fmt.Println(err.Error())
			}

		} else {
			flag1 :=  newData[currency].(float64) > todayValue.High
			flag2 :=  newData[currency].(float64) < todayValue.Low
			if flag1 {
				todayValue.High =  newData[currency].(float64)
			}
			if flag2 {
				todayValue.Low =  newData[currency].(float64)
			}
			if err := services.DB.Save(&todayValue).Error; err != nil {
				fmt.Println(err.Error())
			}

		}


	}





}

func getValue(today,currency string) (models.CurrencyData,error) {

	value := models.CurrencyData{Name:currency}
	err := services.DB.Where("time = ?", today).First(&value)
	if err.Error != nil {

		return value, err.Error
	}
	return value, nil
}