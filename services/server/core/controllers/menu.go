package controllers

import (
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/repositories"
)

// MenuController menu controller
type MenuController struct {
	repository *repositories.MenuRepository
}

// NewMenuController new controller
func NewMenuController() *MenuController {
	return &MenuController{
		repository: repositories.NewMenuRepository(),
	}
}

// Save save attendence
func (c *MenuController) Save(menus []*models.Menu) error {
	return c.repository.Upsert(menus)
}
