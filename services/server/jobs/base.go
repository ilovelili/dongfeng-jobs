package jobs

import (
	"fmt"
	"time"

	"github.com/ilovelili/dongfeng-jobs/services/utils"
	logger "github.com/ilovelili/dongfeng-logger"
)

var (
	config = utils.GetConfig()
)

// sysmtemlog system log
func systemlog(msg, operation string) {
	systemlog := &logger.Log{
		Category: fmt.Sprintf("CRONJOB: %s\n", operation),
		Content:  msg,
		Time:     time.Now(),
	}
	systemlog.SystemLog(logger.CronJobs)
}

// errorlog error log
func errorlog(msg, operation string) {
	errorlog := &logger.Log{
		Category: fmt.Sprintf("CRONJOB Error: %s\n", operation),
		Content:  msg,
		Time:     time.Now(),
	}
	errorlog.ErrorLog(logger.CronJobs)
}
