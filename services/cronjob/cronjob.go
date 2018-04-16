package cronjob

import (
	"github.com/robfig/cron"
	"hzl.im/gin-platform/controllers"
	"fmt"
)

func InitCronJob()  {
	cj := cron.New()
	cj.AddFunc("0 30 * * * *", func1)//Every hour on the half hour
	//cj.AddFunc("@hourly", func1)//Every hour
	//cj.AddFunc("@every 1h30m", func1)//Every hour thirty

	cj.Start()
}

func func1() {
	fmt.Println("cronjob func1 start")
	controllers.GetValueFunc()
}