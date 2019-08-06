package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ilovelili/dongfeng-jobs/services/server/jobs"
	logger "github.com/ilovelili/dongfeng-logger"
	"github.com/micro/cli"
	"github.com/micro/go-micro/cmd"	
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
				run(c, jobs.HeIsDeadJim)
			},
		},
		// ebook convert
		cli.Command{
			Name:  "ebook_convert",
			Usage: "use chrome headless to convert ebook html to pdf",
			Action: func(c *cli.Context) {
				run(c, jobs.ConvertEbookToPDF)
			},
		},
		// ebook merge
		cli.Command{
			Name:  "ebook_merge",
			Usage: "use pdftk to merge ebook pdfs into one file",
			Action: func(c *cli.Context) {
				run(c, jobs.MergeEbook)
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
				run(c, jobs.MenuUpload)
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
				run(c, jobs.RecipeUpload)
			},
		},
		// ingredient nutrition csv file upload
		cli.Command{
			Name:  "ingredient_nutrition_upload",
			Usage: "ingredient nutrition csv file upload",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ingredient_nutrition_file_dir",
					Usage: "Ingredient nutrition file path",
				}},
			Action: func(c *cli.Context) {
				run(c, jobs.IngredientNutritionUpload)
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
				run(c, jobs.RecipeNutritionUpload)
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
	operationname := c.Command.FullName()
	fmt.Println("job starts: ", operationname)
	start := time.Now()

	// fire
	returnCode := fn(c)
	log := &logger.Log{
		Category: "CRONJOB:",
		Content:  fmt.Sprintf("Batch [%s] elapsed time: %v\n", operationname, time.Since(start).Seconds()),
		Time:     time.Now(),
	}
	log.SystemLog(logger.CronJobs)
	os.Exit(returnCode)
}
