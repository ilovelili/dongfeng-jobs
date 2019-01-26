package controllers

import (
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
)

// NutritionController nutrition controller
type NutritionController struct {
	ingredientnutritioncontroller *IngredientNutritionController
	recipenutritioncontroller     *RecipeNutritionController
}

// NewNutritionController new controller
func NewNutritionController() *NutritionController {
	return &NutritionController{
		ingredientnutritioncontroller: NewIngredientNutritionController(),
		recipenutritioncontroller:     NewRecipeNutritionController(),
	}
}

// GetIngredientNutrition get ingredient nutrition
func (c *NutritionController) GetIngredientNutrition(ingredient string) (*models.IngredientNutrition, error) {
	return c.ingredientnutritioncontroller.Get(ingredient)
}

// SaveIngredientNutrition save ingredient nutrition
func (c *NutritionController) SaveIngredientNutrition(nutritions []*models.IngredientNutrition) error {
	return c.ingredientnutritioncontroller.Save(nutritions)
}

// GetRecipeNutrition get recipe nutrition
func (c *NutritionController) GetRecipeNutrition(recipe string) (*models.RecipeNutrition, error) {
	return c.recipenutritioncontroller.Get(recipe)
}

// SaveRecipeNutrition save recipe nutrition
func (c *NutritionController) SaveRecipeNutrition(nutritions []*models.RecipeNutrition) error {
	return c.recipenutritioncontroller.Save(nutritions)
}
