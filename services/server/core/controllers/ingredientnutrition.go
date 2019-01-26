package controllers

import (
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/repositories"
)

// IngredientNutritionController nutrition controller
type IngredientNutritionController struct {
	repository *repositories.IngredientNutritionRepository
}

// NewIngredientNutritionController new controller
func NewIngredientNutritionController() *IngredientNutritionController {
	return &IngredientNutritionController{
		repository: repositories.NewIngredientNutritionRepository(),
	}
}

// Get get ingredient nutrition
func (c *IngredientNutritionController) Get(ingredient string) (*models.IngredientNutrition, error) {
	return c.repository.Select(ingredient)
}

// Save save ingredient nutrition
func (c *IngredientNutritionController) Save(nutritions []*models.IngredientNutrition) error {
	return c.repository.Upsert(nutritions)
}
