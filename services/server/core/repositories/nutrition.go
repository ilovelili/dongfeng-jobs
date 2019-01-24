package repositories

import (
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/olivere/dapper"
)

// NutritionRepository nutrition repository
type NutritionRepository struct {
	session *dapper.Session
}

// NewNutritionRepository init Nutrition repository
func NewNutritionRepository() *NutritionRepository {
	return &NutritionRepository{
		session: session(coreclient()),
	}
}

// Select get Nutrition / ingredients by name
func (r *NutritionRepository) Select(ingredient string) (nutrition *models.Nutrition, err error) {
	var n models.Nutrition
	query := Table("nutritions").Alias("n").Where().Eq("n.ingredient", ingredient).Sql()
	// no rows is actually not an error
	if err = r.session.Find(query, nil).Single(&n); err != nil && norows(err) {
		err = nil
	}
	nutrition = &n
	return
}

// Upsert upsert nutritions
func (r *NutritionRepository) Upsert(nutritions []*models.Nutrition) (err error) {
	tx, err := r.session.Begin()
	if err != nil {
		return
	}

	// upsert by loop
	for _, nutrition := range nutritions {
		query := Table("nutritions").Alias("n").Project("n.id").
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
