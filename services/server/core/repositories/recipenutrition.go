package repositories

import (
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/olivere/dapper"
)

// RecipeNutritionRepository recipe nutrition repository
type RecipeNutritionRepository struct {
	session *dapper.Session
}

// NewRecipeNutritionRepository init recipe nutrition repository
func NewRecipeNutritionRepository() *RecipeNutritionRepository {
	return &RecipeNutritionRepository{
		session: session(coreclient()),
	}
}

// Select get recipe nutrition by name
func (r *RecipeNutritionRepository) Select(recipe string) (nutrition *models.RecipeNutrition, err error) {
	var n models.RecipeNutrition
	query := Table("recipe_nutritions").Alias("n").Where().Eq("n.recipe", recipe).Sql()
	// no rows is actually not an error
	if err = r.session.Find(query, nil).Single(&n); err != nil && norows(err) {
		err = nil
	}
	nutrition = &n
	return
}

// Upsert upsert nutritions
func (r *RecipeNutritionRepository) Upsert(nutritions []*models.RecipeNutrition) (err error) {
	tx, err := r.session.Begin()
	if err != nil {
		return
	}

	// upsert by loop
	for _, nutrition := range nutritions {
		query := Table("recipe_nutritions").Alias("n").Project("n.id").
			Where().
			Eq("n.recipe", nutrition.Recipe).
			Sql()

		var id int64
		err = r.session.Find(query, nil).Scalar(&id)
		if err != nil || 0 == id {
			err = r.session.InsertTx(tx, nutrition)
		} else {
			nutrition.ID = id
			err = r.session.UpdateTx(tx, nutrition)
		}

		if err != nil {
			r.session.Rollback(tx)
			return err
		}
	}

	return r.session.Commit(tx)
}
