package jobs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/controllers"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/micro/cli"
)

// RecipeUpload recipe excel file upload
func RecipeUpload(ctx *cli.Context) int {
	operationname := "upload recipe csv files"

	if len(ctx.String("recipe_file_dir")) == 0 {
		errorlog(fmt.Sprintf("Error: %s", "must assign recipe_file_dir"), operationname)
		return 1
	}

	filepath := ctx.String("recipe_file_dir")
	recipefiles, ingredientfiles, err := findRecipeFiles(filepath)
	if err != nil {
		errorlog(fmt.Sprintf("Error: failed to open %s", filepath), operationname)
		return 1
	}

	// no file
	if len(recipefiles) == 0 {
		return 0
	}

	// exec sequentially
	if err := updateRecipeAndIngredient(recipefiles); err != nil {
		errorlog(fmt.Sprintf("UpdateRecipeAndIngredient error: %s", err), operationname)
		return 1
	}

	if err := updateIngredientMaterial(ingredientfiles); err != nil {
		errorlog(fmt.Sprintf("UpdateIngredientMaterial error: %s", err), operationname)
		return 1
	}

	systemlog("job ended", operationname)
	return 0
}

func findRecipeFiles(dir string) (recipefiles, ingredientfiles []string, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.Index(info.Name(), ".csv") > -1 {
			if strings.Index(info.Name(), ".csv.0") > -1 {
				// recipefiles has name xxx.csv.0 since it's the first sheet of excel
				recipefiles = append(recipefiles, path)
			} else {
				// otherwise has xxx.csv.1 or xxx.csv.2 ...
				ingredientfiles = append(ingredientfiles, path)
			}
		}
		return nil
	})
	return
}

func updateRecipeAndIngredient(recipefiles []string) error {
	recipes := []*models.Recipe{}
	for _, file := range recipefiles {
		_recipes := []*models.Recipe{}
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

	ingredientmap := make(map[string][]string)
	for _, recipe := range recipes {
		if ingredients, ok := ingredientmap[recipe.Name]; !ok {
			ingredientmap[recipe.Name] = []string{recipe.IngredientName}
		} else {
			ingredientmap[recipe.Name] = append(ingredients, recipe.IngredientName)
		}
	}

	_recipes := []*models.Recipe{}
	recipecontroller := controllers.NewRecipeController()
	for recipe, ingredients := range ingredientmap {
		for _, ingredient := range ingredients {
			exist, err := recipecontroller.IngredientExists(ingredient)
			if err != nil {
				return fmt.Errorf("failed to get ingredient %s", err)
			}

			if !exist {
				if err = recipecontroller.SaveIngredient(ingredient); err != nil {
					return fmt.Errorf("failed to save ingredient %s", err)
				}
			}

			_recipes = append(_recipes, &models.Recipe{
				Name:           recipe,
				IngredientName: ingredient,
				CreatedBy:      "AgentSmith",
			})
		}
	}

	// distinct by recipe_ingredient
	distinctedrecipes := []*models.Recipe{}
	keys := make(map[string]bool)
	for _, recipe := range _recipes {
		key := fmt.Sprintf("%s_%s", recipe.Name, recipe.IngredientName)
		if _, ok := keys[key]; ok {
			// already added, ignore
			continue
		} else {
			keys[key] = true
			distinctedrecipes = append(distinctedrecipes, recipe)
		}
	}

	err := recipecontroller.Save(distinctedrecipes)
	if err != nil {
		return fmt.Errorf("failed to save recipe %s", err)
	}

	return nil
}

func updateIngredientMaterial(ingredientfiles []string) error {
	ingredients := []*models.Ingredient{}
	for _, file := range ingredientfiles {
		_ingredients := []*models.Ingredient{}
		csv, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)

		fmt.Println("ingredient csv ", csv.Name())

		if err != nil {
			return fmt.Errorf("failed to open %s", file)
		}
		defer csv.Close()

		if err := gocsv.Unmarshal(csv, &_ingredients); err != nil {
			return fmt.Errorf("failed to unmarshal %s", file)
		}

		ingredients = append(ingredients, _ingredients...)
	}

	if len(ingredients) == 0 {
		return nil
	}

	ingredientcontroller := controllers.NewIngredientController()
	for _, ingredient := range ingredients {
		if err := ingredientcontroller.Save(ingredient); err != nil {
			return fmt.Errorf("failed to save ingredient %s", err)
		}
	}

	return nil
}
