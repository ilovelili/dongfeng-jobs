package repositories

import (
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/olivere/dapper"
)

// IngredientNutritionRepository ingredient nutrition repository
type IngredientNutritionRepository struct {
	session *dapper.Session
}

// NewIngredientNutritionRepository init ingredient nutrition repository
func NewIngredientNutritionRepository() *IngredientNutritionRepository {
	return &IngredientNutritionRepository{
		session: session(coreclient()),
	}
}

// Select get ingredient nutrition by name
func (r *IngredientNutritionRepository) Select(ingredient string) (nutrition *models.IngredientNutrition, err error) {
	var n models.IngredientNutrition
	query := Table("ingredient_nutritions").Alias("n").Where().Eq("n.ingredient", ingredient).Sql()
	// no rows is actually not an error
	if err = r.session.Find(query, nil).Single(&n); err != nil && norows(err) {
		err = nil
	}
	nutrition = &n
	return
}

// Upsert upsert nutritions
func (r *IngredientNutritionRepository) Upsert(nutritions []*models.IngredientNutrition) (err error) {
	tx, err := r.session.Begin()
	if err != nil {
		return
	}

	// upsert by loop
	for _, nutrition := range nutritions {
		query := Table("ingredient_nutritions").Alias("n").Project("n.id").
			Where().
			Eq("n.ingredient", nutrition.Ingredient).
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
