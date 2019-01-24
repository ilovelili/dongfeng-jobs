package jobs

import (
	"fmt"

	excelize "github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/controllers"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	sharedlib "github.com/ilovelili/dongfeng-shared-lib"
	"github.com/micro/cli"
)

// RecipeUpload recipe excel file upload
func RecipeUpload(ctx *cli.Context) int {
	operationname := "upload recipe excel file"

	if len(ctx.String("recipe_file_path")) == 0 {
		errorlog(fmt.Sprintf("Error: %s", "must assign recipe_file_path"), operationname)
		return 1
	}

	filepath := ctx.String("recipe_file_path")
	excel, err := excelize.OpenFile(filepath)
	if err != nil {
		errorlog(fmt.Sprintf("Error: failed to open %s", filepath), operationname)
		return 1
	}

	recipeingredientmap := make(map[string][]string)
	for _, sheet := range excel.WorkBook.Sheets.Sheet {
		// get the "原料" sheet
		if sheet.Name != "原料" {
			continue
		}

		rows := excel.GetRows(sheet.Name)
		recipecolindex, ingredientcolindex := 0, 0
		for rindex, row := range rows {
			if rindex == 0 {
				for cindex, col := range row {
					if col == "菜品名称" {
						recipecolindex = cindex
					} else if col == "原料名称" {
						ingredientcolindex = cindex
					}
				}
			} else {
				if recipecolindex == 0 && ingredientcolindex == 0 {
					errorlog(fmt.Sprintf("Error: invalid recipe file %s", filepath), operationname)
					return 1
				}

				recipe, ingredient := row[recipecolindex], row[ingredientcolindex]
				if v, ok := recipeingredientmap[recipe]; !ok {
					recipeingredientmap[recipe] = []string{ingredient}
				} else {
					if !sharedlib.ContainString(v, ingredient) {
						recipeingredientmap[recipe] = append(v, ingredient)
					}
				}
			}
		}
	}

	recipes := make([]*models.Recipe, 0)
	recipecontroller := controllers.NewRecipeController()
	for recipe, ingredients := range recipeingredientmap {
		for _, ingredient := range ingredients {
			exist, err := recipecontroller.IngredientExists(ingredient)
			if err != nil {
				errorlog(fmt.Sprintf("Error: failed to get ingredient %s", err), operationname)
				return 1
			}

			if !exist {
				if err = recipecontroller.SaveIngredient(ingredient); err != nil {
					errorlog(fmt.Sprintf("Error: failed to save ingredient %s", err), operationname)
					return 1
				}
			}

			recipes = append(recipes, &models.Recipe{
				Name:           recipe,
				IngredientName: ingredient,
				CreatedBy:      "AgentSmith",
			})
		}
	}

	err = recipecontroller.Save(recipes)
	if err != nil {
		errorlog(fmt.Sprintf("Error: failed to save recipe %s", err), operationname)
		return 1
	}

	systemlog("job ended", operationname)
	return 0
}
