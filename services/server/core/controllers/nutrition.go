package controllers

import (
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/repositories"
)

// NutritionController nutrition controller
type NutritionController struct {
	ingredientcontroller *IngredientController
	repository           *repositories.NutritionRepository
}

// NewNutritionController new controller
func NewNutritionController() *NutritionController {
	return &NutritionController{
		ingredientcontroller: NewIngredientController(),
		repository:           repositories.NewNutritionRepository(),
	}
}

// Get get nutrition
func (c *NutritionController) Get(ingredient string) (*models.Nutrition, error) {
	return c.repository.Select(ingredient)
}

// Save save nutrition
func (c *NutritionController) Save(nutritions []*models.Nutrition) error {
	return c.repository.Upsert(nutritions)
}
