package jobs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gocsv "github.com/gocarina/gocsv"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/controllers"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/micro/cli"
)

// IngredientNutritionUpload ingredient nutrition file upload
func IngredientNutritionUpload(ctx *cli.Context) int {
	operationname := "upload ingredient nutrition csv file"

	if len(ctx.String("ingredient_nutrition_file_dir")) == 0 {
		errorlog(fmt.Sprintf("Error: %s", "must assign ingredient_nutrition_file_dir"), operationname)
		return 1
	}

	filepath := ctx.String("ingredient_nutrition_file_dir")
	files, err := findNutritionFiles(filepath)
	if err != nil {
		errorlog(fmt.Sprintf("Error: failed to open %s", filepath), operationname)
		return 1
	}

	if len(files) == 0 {
		return 0
	}

	for _, filepath := range files {
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

		if len(ingredientnutritions) == 0 {
			continue
		}

		// distinct
		distinctedingredientnutritions := []*models.IngredientNutrition{}
		keys := make(map[string]bool)
		for _, ingredientnutrition := range ingredientnutritions {
			if _, ok := keys[ingredientnutrition.Ingredient]; ok {
				// already added, ignore
				continue
			} else {
				keys[ingredientnutrition.Ingredient] = true
				distinctedingredientnutritions = append(distinctedingredientnutritions, ingredientnutrition)
			}
		}

		nutritioncontroller := controllers.NewNutritionController()
		err = nutritioncontroller.SaveIngredientNutrition(distinctedingredientnutritions)
		if err != nil {
			errorlog(fmt.Sprintf("Error on saving ingredient nutritions: %s", err), operationname)
			return 1
		}
	}

	systemlog("job ended", operationname)
	return 0
}

// RecipeNutritionUpload recipe nutrition file upload
func RecipeNutritionUpload(ctx *cli.Context) int {
	operationname := "upload recipe nutrition csv file"

	if len(ctx.String("recipe_nutrition_file_dir")) == 0 {
		errorlog(fmt.Sprintf("Error: %s", "must assign recipe_nutrition_file_dir"), operationname)
		return 1
	}

	filepath := ctx.String("recipe_nutrition_file_dir")
	files, err := findNutritionFiles(filepath)
	if err != nil {
		errorlog(fmt.Sprintf("Error: failed to open %s", filepath), operationname)
		return 1
	}

	if len(files) == 0 {
		return 0
	}

	for _, filepath := range files {
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

		if len(recipenutritions) == 0 {
			continue
		}

		// distinct
		distinctedrecipenutritions := []*models.RecipeNutrition{}
		keys := make(map[string]bool)
		for _, recipenutrition := range recipenutritions {
			if _, ok := keys[recipenutrition.Recipe]; ok {
				// already added, ignore
				continue
			} else {
				keys[recipenutrition.Recipe] = true
				distinctedrecipenutritions = append(distinctedrecipenutritions, recipenutrition)
			}
		}

		nutritioncontroller := controllers.NewNutritionController()
		err = nutritioncontroller.SaveRecipeNutrition(distinctedrecipenutritions)
		if err != nil {
			errorlog(fmt.Sprintf("Error on saving recipe nutritions: %s", err), operationname)
			return 1
		}
	}

	systemlog("job ended", operationname)
	return 0
}

func findNutritionFiles(dir string) (files []string, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.Index(info.Name(), ".csv") > -1 {
			files = append(files, path)
		}
		return nil
	})
	return
}
