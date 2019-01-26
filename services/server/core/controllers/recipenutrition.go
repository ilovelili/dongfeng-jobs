package controllers

import (
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/repositories"
)

// RecipeNutritionController recipe nutrition controller
type RecipeNutritionController struct {
	repository *repositories.RecipeNutritionRepository
}

// NewRecipeNutritionController new controller
func NewRecipeNutritionController() *RecipeNutritionController {
	return &RecipeNutritionController{
		repository: repositories.NewRecipeNutritionRepository(),
	}
}

// Get get recipe nutrition
func (c *RecipeNutritionController) Get(recipe string) (*models.RecipeNutrition, error) {
	return c.repository.Select(recipe)
}

// Save save recipe nutrition
func (c *RecipeNutritionController) Save(recipes []*models.RecipeNutrition) error {
	return c.repository.Upsert(recipes)
}
