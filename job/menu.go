package job

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ilovelili/dongfeng-jobs/core/model"
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
	operation := "upload menu csv file"

	if len(ctx.String("menu_file_path")) == 0 {
		errorLog(fmt.Sprintf("Error: %s", "must assign menu_file_path"), operation)
		return 1
	}

	filepath := ctx.String("menu_file_path")
	file, err := os.Open(filepath)
	if err != nil {
		errorLog(fmt.Sprintf("Error: failed to open %s", filepath), operation)
		return 1
	}
	defer file.Close()

	menus := make([]*model.Menu, 5)
	// init with empty menu
	for i := range menus {
		menus[i] = new(model.Menu)
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
			errorLog(fmt.Sprintf("Error: %s", err), operation)
			return 1
		}

		// the csv file has 5 columns
		if len(line) != 5 {
			errorLog(fmt.Sprintf("Error: invalid file %s", filepath), operation)
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

	_menus := []*model.Menu{}
	for _, menu := range menus {
		_submenus, err := parseRecipes(menu)
		if err != nil {
			errorLog(err.Error(), operation)
			return 1
		}

		_menus = append(_menus, _submenus...)
	}

	if len(_menus) == 0 {
		return 0
	}

	if err := menuRepo.SaveAll(_menus); err != nil {
		errorLog(fmt.Sprintf("Error on saving menu: %s", err), operation)
		return 1
	}

	systemLog("job ended", operation)
	return 0
}

// parseRecipes parse recipes
func parseRecipes(m *model.Menu) ([]*model.Menu, error) {
	result := []*model.Menu{}
	recipes := parse(m.JuniorBreakfastRaw)

	for _, recipeName := range recipes {
		menu := &model.Menu{
			Date:             m.Date,
			BreakfastOrLunch: model.Breakfast,
			JuniorOrSenior:   model.Junior,
		}

		recipe, err := recipeRepo.FirstOrCreate(&model.Recipe{Name: recipeName})
		if err != nil {
			return result, err
		}
		menu.RecipeID = recipe.ID
		result = append(result, menu)
	}

	recipes = parse(m.JuniorLunchRaw)
	for _, recipeName := range recipes {
		menu := &model.Menu{
			Date:             m.Date,
			BreakfastOrLunch: model.Lunch,
			JuniorOrSenior:   model.Junior,
		}

		recipe, err := recipeRepo.FirstOrCreate(&model.Recipe{Name: recipeName})
		if err != nil {
			return result, err
		}
		menu.RecipeID = recipe.ID
		result = append(result, menu)
	}

	recipes = parse(m.JuniorSnackRaw)
	for _, recipeName := range recipes {
		menu := &model.Menu{
			Date:             m.Date,
			BreakfastOrLunch: model.Snack,
			JuniorOrSenior:   model.Junior,
		}

		recipe, err := recipeRepo.FirstOrCreate(&model.Recipe{Name: recipeName})
		if err != nil {
			return result, err
		}
		menu.RecipeID = recipe.ID
		result = append(result, menu)
	}

	recipes = parse(m.SeniorBreakfastRaw)
	for _, recipeName := range recipes {
		menu := &model.Menu{
			Date:             m.Date,
			BreakfastOrLunch: model.Breakfast,
			JuniorOrSenior:   model.Senior,
		}

		recipe, err := recipeRepo.FirstOrCreate(&model.Recipe{Name: recipeName})
		if err != nil {
			return result, err
		}
		menu.RecipeID = recipe.ID
		result = append(result, menu)
	}

	recipes = parse(m.SeniorLunchRaw)
	for _, recipeName := range recipes {
		menu := &model.Menu{
			Date:             m.Date,
			BreakfastOrLunch: model.Lunch,
			JuniorOrSenior:   model.Senior,
		}

		recipe, err := recipeRepo.FirstOrCreate(&model.Recipe{Name: recipeName})
		if err != nil {
			return result, err
		}
		menu.RecipeID = recipe.ID
		result = append(result, menu)
	}

	recipes = parse(m.SeniorSnackRaw)
	for _, recipeName := range recipes {
		menu := &model.Menu{
			Date:             m.Date,
			BreakfastOrLunch: model.Snack,
			JuniorOrSenior:   model.Senior,
		}

		recipe, err := recipeRepo.FirstOrCreate(&model.Recipe{Name: recipeName})
		if err != nil {
			return result, err
		}
		menu.RecipeID = recipe.ID
		result = append(result, menu)
	}

	return result, nil
}

func parse(raw string) []string {
	return strings.Split(raw, "|")
}

func formatDate(rawdate string) string {
	return fmt.Sprintf("%s-%s", time.Now().Format("2006"), rawdate)
}
