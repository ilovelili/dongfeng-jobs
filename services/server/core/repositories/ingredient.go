package repositories

import (
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/olivere/dapper"
)

// IngredientRepository Ingredient repository
type IngredientRepository struct {
	session *dapper.Session
}

// NewIngredientRepository init Ingredient repository
func NewIngredientRepository() *IngredientRepository {
	return &IngredientRepository{
		session: session(coreclient()),
	}
}

// Count get Ingredient / ingredients by name
func (r *IngredientRepository) Count(name string) (count int64, err error) {
	query := Table("ingredients").Alias("i").Project(`count("i.*")`).Where().Eq("i.name", name).Sql()
	// no rows is actually not an error
	if count, err = r.session.Count(query, nil); err != nil && norows(err) {
		err = nil
	}

	return
}

// Select select Ingredients
func (r *IngredientRepository) Select(name string) (ingredient *models.Ingredient, err error) {
	var i models.Ingredient
	query := Table("ingredients").Alias("i").Where().Eq("i.name", name).Sql()
	if err = r.session.Find(query, nil).Single(&i); err != nil && norows(err) {
		err = nil
	}

	if err == nil {
		ingredient = &i
	}
	return
}

// Upsert upsert Ingredients
func (r *IngredientRepository) Upsert(ingredient *models.Ingredient) (err error) {
	tx, err := r.session.Begin()
	if err != nil {
		return
	}

	query := Table("ingredients").Alias("i").Project("i.id").
		Where().Eq("i.name", ingredient.Name).
		Sql()

	var id int64
	err = r.session.Find(query, nil).Scalar(&id)
	if err != nil || 0 == id {
		err = r.session.InsertTx(tx, ingredient)
	} else {
		ingredient.ID = id
		err = r.session.UpdateTx(tx, ingredient)
	}

	if err != nil {
		r.session.Rollback(tx)
		return err
	}

	return r.session.Commit(tx)
}
