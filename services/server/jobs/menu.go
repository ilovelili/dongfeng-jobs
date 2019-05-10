package jobs

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ilovelili/dongfeng-jobs/services/server/core/controllers"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/micro/cli"
)

// menuIndexer menu file index
type menuIndexer int64

const (
	header menuIndexer = iota + 1
	seniorbreakfast
	seniorlunch
	seniorsnack
	juniorbreakfast
	juniorlunch
	juniorsnack
)

// MenuUpload menu csv file upload
func MenuUpload(ctx *cli.Context) int {
	operationname := "upload menu csv file"

	if len(ctx.String("menu_file_path")) == 0 {
		errorlog(fmt.Sprintf("Error: %s", "must assign menu_file_path"), operationname)
		return 1
	}

	filepath := ctx.String("menu_file_path")
	file, err := os.Open(filepath)
	if err != nil {
		errorlog(fmt.Sprintf("Error: failed to open %s", filepath), operationname)
		return 1
	}
	defer file.Close()

	menus := make([]*models.Menu, 5)
	// init with empty menu
	for i := range menus {
		menus[i] = new(models.Menu)
	}

	reader := csv.NewReader(file)
	// options are available at:
	// http://golang.org/src/pkg/encoding/csv/reader.go?s=3213:3671#L94
	lineCount := 1
	for {
		line, err := reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			errorlog(fmt.Sprintf("Error: %s", err), operationname)
			return 1
		}

		// the csv file has 5 columns
		if len(line) != 5 {
			errorlog(fmt.Sprintf("Error: invalid file %s", filepath), operationname)
			return 1
		}

		switch lineCount {
		case int(header):
			for i := range menus {
				menus[i].Date = formatDate(line[i])
			}
		case int(seniorbreakfast):
			for i := range menus {
				menus[i].SeniorBreakfastRaw = line[i]
			}
		case int(seniorlunch):
			for i := range menus {
				menus[i].SeniorLunchRaw = line[i]
			}
		case int(seniorsnack):
			for i := range menus {
				menus[i].SeniorSnackRaw = line[i]
			}
		case int(juniorbreakfast):
			for i := range menus {
				menus[i].JuniorBreakfastRaw = line[i]
			}
		case int(juniorlunch):
			for i := range menus {
				menus[i].JuniorLunchRaw = line[i]
			}
		case int(juniorsnack):
			for i := range menus {
				menus[i].JuniorSnackRaw = line[i]
			}
		}

		lineCount++
	}

	_menus := []*models.Menu{}
	for _, menu := range menus {
		_menus = append(_menus, menu.ParseRecipes()...)
	}

	if len(_menus) == 0 {
		return 0
	}

	menucontroller := controllers.NewMenuController()
	err = menucontroller.Save(_menus)
	if err != nil {
		errorlog(fmt.Sprintf("Error on saving menu: %s", err), operationname)
		return 1
	}

	systemlog("job ended", operationname)
	return 0
}

func formatDate(rawdate string) string {
	return fmt.Sprintf("%s-%s", time.Now().Format("2006"), rawdate)
}
