package repositories

import (
	"fmt"

	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/olivere/dapper"
)

// EbookRepository ebook repository
type EbookRepository struct {
	session *dapper.Session
}

// NewEbookRepository init ebook repository
func NewEbookRepository() *EbookRepository {
	return &EbookRepository{
		session: session(coreclient()),
	}
}

// Select select ebooks
func (r *EbookRepository) Select() (ebooks []*models.Ebook, err error) {
	query := Table("ebooks").Alias("e").
		Where().
		Eq("e.converted", 0).
		Sql()

	if err = r.session.Find(query, nil).All(&ebooks); err != nil && norows(err) {
		err = nil
	}

	return
}

// Update set converted flag to true
func (r *EbookRepository) Update(ebook *models.Ebook) (err error) {
	query := Table("ebooks").Alias("e").Where().
		Eq("e.year", ebook.Year).
		Eq("e.class", ebook.Class).
		Eq("e.name", ebook.Name).
		Eq("e.date", ebook.Date).
		Sql()

	var _ebook models.Ebook
	err = r.session.Find(query, nil).Single(&_ebook)
	if err != nil || 0 == _ebook.ID {
		return fmt.Errorf("not found")
	}

	ebook.ID = _ebook.ID
	ebook.Hash = _ebook.Hash
	ebook.Converted = true
	err = r.session.Update(ebook)
	return
}
