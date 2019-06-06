package controllers

import (
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/repositories"
)

// EbookController ebook controller
type EbookController struct {
	repository *repositories.EbookRepository
}

// NewEbookController new controller
func NewEbookController() *EbookController {
	return &EbookController{
		repository: repositories.NewEbookRepository(),
	}
}

// SaveEbook save ebook
func (c *EbookController) SaveEbook(ebook *models.Ebook) error {
	return c.repository.Update(ebook)

}

// GetEbooks get ebooks
func (c *EbookController) GetEbooks() ([]*models.Ebook, error) {
	return c.repository.Select()
}
