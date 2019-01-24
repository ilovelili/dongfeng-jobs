package repositories

import (
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/olivere/dapper"
)

// MenuRepository menu repository
type MenuRepository struct {
	session *dapper.Session
}

// NewMenuRepository init menu repository
func NewMenuRepository() *MenuRepository {
	return &MenuRepository{
		session: session(coreclient()),
	}
}

// Upsert upsert Menus
func (r *MenuRepository) Upsert(menus []*models.Menu) (err error) {
	tx, err := r.session.Begin()
	if err != nil {
		return
	}

	// upsert by loop
	for _, menu := range menus {
		query := Table("menus").Alias("m").Project("m.id").
			Where().
			Eq("r.date", menu.Date).
			Eq("r.recipe", menu.Recipe).
			Sql()

		var id int64
		err = r.session.Find(query, nil).Scalar(&id)
		if err != nil || 0 == id {
			err = r.session.InsertTx(tx, menu)
		} else {
			menu.ID = id
			err = r.session.UpdateTx(tx, menu)
		}

		if err != nil {
			r.session.Rollback(tx)
			return err
		}
	}

	return r.session.Commit(tx)
}
