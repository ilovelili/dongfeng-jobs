package jobs

import (
	"fmt"
	"os"

	gocsv "github.com/gocarina/gocsv"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/controllers"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/micro/cli"
)

// IngredientNutritionUpload ingredient nutrition file upload
func IngredientNutritionUpload(ctx *cli.Context) int {
	operationname := "upload ingredient nutrition csv file"

	if len(ctx.String("ingredient_nutrition_file_path")) == 0 {
		errorlog(fmt.Sprintf("Error: %s", "must assign ingredient_nutrition_file_path"), operationname)
		return 1
	}

	filepath := ctx.String("ingredient_nutrition_file_path")
	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		errorlog(fmt.Sprintf("Error: failed to open %s", filepath), operationname)
		return 1
	}
	defer file.Close()

	ingredientnutritions := []*models.IngredientNutrition{}
	if err := gocsv.UnmarshalFile(file, &ingredientnutritions); err != nil {
		errorlog(fmt.Sprintf("Error: invalid ingredient nutrition file %s", filepath), operationname)
		return 1
	}

	nutritioncontroller := controllers.NewNutritionController()
	err = nutritioncontroller.SaveIngredientNutrition(ingredientnutritions)
	if err != nil {
		errorlog(fmt.Sprintf("Error on saving ingredient nutritions: %s", err), operationname)
		return 1
	}

	systemlog("job ended", operationname)
	return 0
}

// RecipeNutritionUpload recipe nutrition file upload
func RecipeNutritionUpload(ctx *cli.Context) int {
	operationname := "upload recipe nutrition csv file"

	if len(ctx.String("recipe_nutrition_file_path")) == 0 {
		errorlog(fmt.Sprintf("Error: %s", "must assign recipe_nutrition_file_path"), operationname)
		return 1
	}

	filepath := ctx.String("recipe_nutrition_file_path")
	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		errorlog(fmt.Sprintf("Error: failed to open %s", filepath), operationname)
		return 1
	}
	defer file.Close()

	recipenutritions := []*models.RecipeNutrition{}
	if err := gocsv.UnmarshalFile(file, &recipenutritions); err != nil {
		errorlog(fmt.Sprintf("Error: invalid recipe nutrition file %s", filepath), operationname)
		return 1
	}

	nutritioncontroller := controllers.NewNutritionController()
	err = nutritioncontroller.SaveRecipeNutrition(recipenutritions)
	if err != nil {
		errorlog(fmt.Sprintf("Error on saving recipe nutritions: %s", err), operationname)
		return 1
	}

	systemlog("job ended", operationname)
	return 0
}
