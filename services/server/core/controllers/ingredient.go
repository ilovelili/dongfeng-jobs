package controllers

import (
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/repositories"
)

// IngredientController ingredient controller
type IngredientController struct {
	ingredientcontroller *IngredientController
	repository           *repositories.IngredientRepository
}

// NewIngredientController new controller
func NewIngredientController() *IngredientController {
	return &IngredientController{
		repository: repositories.NewIngredientRepository(),
	}
}

// Exists ingredient exists or not
func (c *IngredientController) Exists(ingredient string) (bool, error) {
	count, err := c.repository.Count(ingredient)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Save save ingredient
func (c *IngredientController) Save(ingredient *models.Ingredient) error {
	return c.repository.Upsert(ingredient)
}

// SelectByName select by name
func (c *IngredientController) SelectByName(ingredient string) (*models.Ingredient, error) {
	return c.repository.Select(ingredient)
}
