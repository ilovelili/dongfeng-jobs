package job

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/ilovelili/dongfeng-jobs/core/model"
	"github.com/micro/cli"
)

// RecipeNutritionUpload recipe nutrition file upload
func RecipeNutritionUpload(ctx *cli.Context) int {
	operation := "upload recipe nutrition csv file"

	if len(ctx.String("recipe_nutrition_file_dir")) == 0 {
		errorLog(fmt.Sprintf("Error: %s", "must assign recipe_nutrition_file_dir"), operation)
		return 1
	}

	filepath := ctx.String("recipe_nutrition_file_dir")
	files, err := findNutritionFiles(filepath)
	if err != nil {
		errorLog(fmt.Sprintf("Error: failed to open %s", filepath), operation)
		return 1
	}

	if len(files) == 0 {
		return 0
	}

	for _, filepath := range files {
		file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
		if err != nil {
			errorLog(fmt.Sprintf("Error: failed to open %s", filepath), operation)
			return 1
		}
		defer file.Close()

		nutritions := []*model.RecipeNutrition{}
		if err := gocsv.UnmarshalFile(file, &nutritions); err != nil {
			errorLog(fmt.Sprintf("Error: invalid recipe nutrition file %s", filepath), operation)
			return 1
		}

		if len(nutritions) == 0 {
			continue
		}

		// distinct
		distinctedNutritions := []*model.RecipeNutrition{}
		keys := make(map[string]bool)
		for _, nutrition := range nutritions {
			if _, ok := keys[nutrition.Recipe]; ok {
				// already added, ignores
				continue
			} else {
				keys[nutrition.Recipe] = true
				distinctedNutritions = append(distinctedNutritions, nutrition)
			}
		}

		for _, nutrition := range distinctedNutritions {
			nutrition, err = nutritionRepo.FirstOrCreate(nutrition)
			if err != nil {
				errorLog(fmt.Sprintf("Error on saving recipe nutrition: %s", err), operation)
				return 1
			}

			recipe := &model.Recipe{Name: nutrition.Recipe}
			recipe, err = recipeRepo.Find(recipe)
			if err == nil {
				recipe.RecipeNutritionID = &nutrition.ID
				err = recipeRepo.Save(recipe)
				if err != nil {
					errorLog(fmt.Sprintf("Error on saving recipe: %s", err), operation)
					return 1
				}
			}
		}

	}

	systemLog("job ended", operation)
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
