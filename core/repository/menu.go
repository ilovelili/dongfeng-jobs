package repository

import "github.com/ilovelili/dongfeng-jobs/core/model"

// Menu menu repository
type Menu struct{}

// NewMenuRepository init menu repository
func NewMenuRepository() *Menu {
	db().AutoMigrate(&model.Menu{}, &model.Recipe{}, &model.RecipeNutrition{}, &model.Ingredient{}, &model.IngredientCategory{})
	return new(Menu)
}

// SaveAll save all menus
func (r *Menu) SaveAll(menus []*model.Menu) error {
	tx := db().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, menu := range menus {
		_menu := new(model.Menu)
		if err := db().Where("recipe_id = ? AND breakfast_or_lunch = ? AND junior_or_senior = ?", menu.RecipeID, menu.BreakfastOrLunch, menu.JuniorOrSenior).First(&_menu); err == nil {
			// menu found, skip (no need to insert save menu)
			continue
		}

		if err := tx.Save(menu).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
