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

// RecipeUpload recipe excel file upload
func RecipeUpload(ctx *cli.Context) int {
	operation := "upload recipe csv files"

	if len(ctx.String("recipe_file_dir")) == 0 {
		errorLog(fmt.Sprintf("Error: %s", "must assign recipe_file_dir"), operation)
		return 1
	}

	filepath := ctx.String("recipe_file_dir")
	recipefiles, err := findRecipeFiles(filepath)
	if err != nil {
		errorLog(fmt.Sprintf("Error: failed to open %s", filepath), operation)
		return 1
	}

	// no file
	if len(recipefiles) == 0 {
		return 0
	}

	// exec sequentially
	if err := updateRecipeAndIngredient(recipefiles); err != nil {
		errorLog(fmt.Sprintf("UpdateRecipeAndIngredient error: %s", err), operation)
		return 1
	}

	systemLog("job ended", operation)
	return 0
}

func findRecipeFiles(dir string) (recipefiles []string, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.Index(info.Name(), ".csv.0") > -1 {
			recipefiles = append(recipefiles, path)
		}
		return nil
	})

	return
}

func updateRecipeAndIngredient(recipefiles []string) error {
	recipes := []*model.Recipe{}
	for _, file := range recipefiles {
		_recipes := []*model.Recipe{}
		csv, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to open %s", file)
		}
		defer csv.Close()

		if err := gocsv.Unmarshal(csv, &_recipes); err != nil {
			return fmt.Errorf("failed to unmarshal %s", file)
		}

		recipes = append(recipes, _recipes...)
	}

	if len(recipes) == 0 {
		return nil
	}

	// merge recipes with same name
	recipeMap := make(map[string]*model.Recipe)
	for _, recipe := range recipes {
		ingredient := &model.Ingredient{
			Ingredient: recipe.CSVIngredient,
			Alias:      recipe.CSVIngredient,
		}

		_ingredient, err := ingredientRepo.FirstOrCreate(ingredient)
		if err != nil {
			return err
		}

		if _recipe, ok := recipeMap[recipe.Name]; ok {
			ingredientFound := false
			for _, ingredient := range _recipe.Ingredients {
				if ingredient.Ingredient == _ingredient.Ingredient {
					ingredientFound = true
					break
				}
			}
			if !ingredientFound {
				_recipe.Ingredients = append(_recipe.Ingredients, _ingredient)
			}
		} else {
			recipe.Ingredients = append(recipe.Ingredients, _ingredient)
			recipeMap[recipe.Name] = recipe
		}
	}

	for _, recipe := range recipeMap {
		if err := recipeRepo.Save(recipe); err != nil {
			return err
		}
	}

	return nil
}
