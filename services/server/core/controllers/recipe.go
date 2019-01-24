package controllers

import (
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/repositories"
)

// RecipeController recipe controller
type RecipeController struct {
	ingredientcontroller *IngredientController
	repository           *repositories.RecipeRepository
}

// NewRecipeController new controller
func NewRecipeController() *RecipeController {
	return &RecipeController{
		ingredientcontroller: NewIngredientController(),
		repository:           repositories.NewRecipeRepository(),
	}
}

// Get get recipe
func (c *RecipeController) Get(name string) ([]*models.Recipe, error) {
	return c.repository.Select(name)
}

// Save save recipe
func (c *RecipeController) Save(recipes []*models.Recipe) error {
	for _, recipe := range recipes {
		ingredient, err := c.SelectIngredient(recipe.IngredientName)
		if err != nil {
			return err
		}
		recipe.Ingredient = ingredient.ID
	}

	return c.repository.Upsert(recipes)
}

// IngredientExists ingredient exists or not
func (c *RecipeController) IngredientExists(ingredient string) (bool, error) {
	return c.ingredientcontroller.Exists(ingredient)
}

// SaveIngredient save ingredient
func (c *RecipeController) SaveIngredient(ingredient string) error {
	return c.ingredientcontroller.Save(&models.Ingredient{
		Name: ingredient,
	})
}

// SelectIngredient select ingredient
func (c *RecipeController) SelectIngredient(ingredient string) (*models.Ingredient, error) {
	return c.ingredientcontroller.SelectByName(ingredient)
}
