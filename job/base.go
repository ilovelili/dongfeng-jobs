package job

import (
	"github.com/ilovelili/dongfeng-jobs/core/repository"
	"github.com/ilovelili/dongfeng-jobs/util"
	logger "github.com/ilovelili/logger"
)

var (
	config         = util.LoadConfig()
	ebooksRepo     = repository.NewEbookRepository()
	ingredientRepo = repository.NewIngredientRepository()
	recipeRepo     = repository.NewRecipeRepository()
	nutritionRepo = repository.NewRecipeNutritionRepository()
	menuRepo       = repository.NewMenuRepository()
)

func systemLog(msg, operation string) {
	logger.Type("system").WithFields(logger.Fields{
		"operation": operation,
	}).Infoln(msg)
}

func errorLog(msg, operation string) {
	logger.Type("error").WithFields(logger.Fields{
		"operation": operation,
	}).Errorln(msg)
}
