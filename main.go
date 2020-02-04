package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ilovelili/dongfeng-jobs/job"
	"github.com/ilovelili/logger"
	"github.com/micro/cli"
	"github.com/micro/go-micro/config/cmd"
)

const (
	appName = "dongfeng-jobs"
)

func commands() []cli.Command {
	return []cli.Command{
		// test
		cli.Command{
			Name:  "test",
			Usage: "test cmd",
			Action: func(c *cli.Context) {
				run(c, func(*cli.Context) int {
					fmt.Println("Good to go")
					return 0
				})
			},
		},
		// test error
		cli.Command{
			Name:  "test_error",
			Usage: "test cmd to return error",
			Action: func(c *cli.Context) {
				run(c, job.HeIsDeadJim)
			},
		},
		// ebook convert
		cli.Command{
			Name:  "ebook_convert",
			Usage: "use chrome headless to convert ebook html to pdf",
			Action: func(c *cli.Context) {
				run(c, job.ConvertEbookToPDF)
			},
		},
		// ebook merge
		cli.Command{
			Name:  "ebook_merge",
			Usage: "use pdftk to merge ebook pdfs into one file",
			Action: func(c *cli.Context) {
				run(c, job.MergeEbook)
			},
		},
		// menu csv file upload
		cli.Command{
			Name:  "menu_upload",
			Usage: "menu csv file upload",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "menu_file_path",
					Usage: "Menu file path",
				}},
			Action: func(c *cli.Context) {
				run(c, job.MenuUpload)
			},
		},
		// menu csv file upload
		cli.Command{
			Name:  "recipe_upload",
			Usage: "recipe excel file upload",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "recipe_file_dir",
					Usage: "Recipe file directory",
				}},
			Action: func(c *cli.Context) {
				run(c, job.RecipeUpload)
			},
		},
		// recipe nutrition csv file upload
		cli.Command{
			Name:  "recipe_nutrition_upload",
			Usage: "recipe nutrition csv file upload",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "recipe_nutrition_file_dir",
					Usage: "recipe nutrition file directory",
				}},
			Action: func(c *cli.Context) {
				run(c, job.RecipeNutritionUpload)
			},
		},
	}
}

func main() {
	app := cmd.App()
	app.Commands = append(app.Commands, commands()...)
	cmd.Init()
}

func run(c *cli.Context, fn func(*cli.Context) int) {
	operation := c.Command.FullName()
	start := time.Now()

	// fire
	returnCode := fn(c)

	logger.Type("application").WithFields(logger.Fields{
		"operation":    operation,
		"elapsed time": time.Since(start).Seconds(),
	}).Infoln("")

	os.Exit(returnCode)
}
