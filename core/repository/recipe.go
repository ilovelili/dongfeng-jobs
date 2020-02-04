package repository

import "github.com/ilovelili/dongfeng-jobs/core/model"

// Recipe recipe repository
type Recipe struct{}

// NewRecipeRepository new recipe repository
func NewRecipeRepository() *Recipe {
	db().AutoMigrate(&model.Menu{}, &model.Recipe{}, &model.RecipeNutrition{}, &model.Ingredient{}, &model.IngredientCategory{})
	return new(Recipe)
}

// Find find first or create recipe
func (r *Recipe) Find(recipe *model.Recipe) (*model.Recipe, error) {
	err := db().Where("name = ?", recipe.Name).Find(&recipe).Error
	return recipe, err
}

// FirstOrCreate find first or create recipe
func (r *Recipe) FirstOrCreate(recipe *model.Recipe) (*model.Recipe, error) {
	err := db().Where("name = ?", recipe.Name).FirstOrCreate(&recipe).Error
	return recipe, err
}

// Save save recipe
func (r *Recipe) Save(recipe *model.Recipe) error {
	if err := db().Where("name = ?", recipe.Name).First(&recipe); err == nil {
		// recipe found, skip (no need to insert save menu)
		return nil
	}
	return db().Save(recipe).Error
}

// Update update recipe
func (r *Recipe) Update(recipe *model.Recipe) error {
	return db().Save(&recipe).Error
}
