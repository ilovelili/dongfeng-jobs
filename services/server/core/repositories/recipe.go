package repositories

import (
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/olivere/dapper"
)

// RecipeRepository recipe repository
type RecipeRepository struct {
	session *dapper.Session
}

// NewRecipeRepository init recipe repository
func NewRecipeRepository() *RecipeRepository {
	return &RecipeRepository{
		session: session(coreclient()),
	}
}

// Select get recipe / ingredients by name
func (r *RecipeRepository) Select(name string) (recipes []*models.Recipe, err error) {
	query := Table("recipes").Alias("r").Where().Eq("r.name", name).Sql()
	// no rows is actually not an error
	if err = r.session.Find(query, nil).All(&recipes); err != nil && norows(err) {
		err = nil
	}

	return
}

// Upsert upsert recipes
func (r *RecipeRepository) Upsert(recipes []*models.Recipe) (err error) {
	tx, err := r.session.Begin()
	if err != nil {
		return
	}

	// upsert by loop
	for _, recipe := range recipes {
		query := Table("recipes").Alias("r").
			Where().
			Eq("r.name", recipe.Name).
			Eq("r.ingredient_id", recipe.Ingredient).
			Sql()

		var _r models.Recipe
		err = r.session.Find(query, nil).Single(&_r)
		if err != nil || 0 == _r.ID {
			err = r.session.InsertTx(tx, recipe)
		} else {
			recipe.ID = _r.ID
			recipe.UnitAmount = _r.UnitAmount
			err = r.session.UpdateTx(tx, recipe)
		}

		if err != nil {
			r.session.Rollback(tx)
			return err
		}
	}

	return r.session.Commit(tx)
}
