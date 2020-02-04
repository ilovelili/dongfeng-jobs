package repository

import "github.com/ilovelili/dongfeng-jobs/core/model"

// RecipeNutrition recipe nutrition repository
type RecipeNutrition struct{}

// NewRecipeNutritionRepository init recipe nutrition repository
func NewRecipeNutritionRepository() *RecipeNutrition {
	db().AutoMigrate(&model.Menu{}, &model.Recipe{}, &model.RecipeNutrition{}, &model.Ingredient{}, &model.IngredientCategory{})
	return new(RecipeNutrition)
}

// FirstOrCreate first or create nutrition
func (r *RecipeNutrition) FirstOrCreate(nutrition *model.RecipeNutrition) (*model.RecipeNutrition, error) {
	err := db().Where("recipe = ?", nutrition.Recipe).FirstOrCreate(&nutrition).Error
	return nutrition, err
}
