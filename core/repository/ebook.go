package repository

import (
	"github.com/ilovelili/dongfeng-jobs/core/model"
)

// Ebook ebook repository
type Ebook struct{}

// NewEbookRepository init ebook repository
func NewEbookRepository() *Ebook {
	db().AutoMigrate(&model.Ebook{})
	return new(Ebook)
}

// FindAll find all unconverted
func (r *Ebook) FindAll() ([]*model.Ebook, error) {
	ebooks := []*model.Ebook{}
	err := db().Where("ebooks.converted = 0").Find(&ebooks).Error
	return ebooks, err
}

// SetConverted set converted flag to true
func (r *Ebook) SetConverted(ebook *model.Ebook) error {
	_ebook := new(model.Ebook)
	if err := db().Where("ebooks.date = ? AND ebooks.pupil_id = ?", ebook.Date, ebook.PupilID).Find(&_ebook).Error; err != nil {
		return err
	}

	_ebook.Converted = true
	return db().Model(&model.Ebook{}).Save(_ebook).Error
}
