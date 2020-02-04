package repository

import "github.com/ilovelili/dongfeng-jobs/core/model"

// Ingredient ingredient repository
type Ingredient struct{}

// NewIngredientRepository init Ingredient repository
func NewIngredientRepository() *Ingredient {
	db().AutoMigrate(&model.Menu{}, &model.Recipe{}, &model.RecipeNutrition{}, &model.Ingredient{}, &model.IngredientCategory{})
	return new(Ingredient)
}

// FirstOrCreate find first or create ingredient
func (r *Ingredient) FirstOrCreate(ingredient *model.Ingredient) (*model.Ingredient, error) {
	err := db().Where("ingredients.ingredient = ?", ingredient.Ingredient).FirstOrCreate(&ingredient).Error
	return ingredient, err
}
