package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ilovelili/dongfeng-jobs/services/server/jobs"
	logger "github.com/ilovelili/dongfeng-logger"
	"github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	_ "github.com/micro/go-plugins/registry/kubernetes"
)

func commands() []cli.Command {
	return []cli.Command{
		// test fail
		cli.Command{
			Name:  "test_fail",
			Usage: "test cmd",
			Action: func(c *cli.Context) {
				run(c, jobs.HeIsDeadJim)
			},
		},
		// test
		cli.Command{
			Name:  "test",
			Usage: "test cmd",
			Action: func(c *cli.Context) {
				run(c, func() int {
					fmt.Println("Good to go")
					return 0
				})
			},
		},
	}
}

func main() {
	app := cmd.App()
	app.Commands = append(app.Commands, commands()...)
	cmd.Init()
}

func run(c *cli.Context, fn func() int) {
	operationname := c.Command.FullName()
	fmt.Println("job starts: ", operationname)
	start := time.Now()

	// fire
	returnCode := fn()
	systemlog := &logger.Log{
		Category: "CRONJOB:",
		Content:  fmt.Sprintf("Batch [%s] elapsed time: %v\n", operationname, time.Since(start).Seconds()),
		Time:     time.Now(),
	}

	systemlog.SystemLog(logger.CronJobs)
	os.Exit(returnCode)
}
