package jobs

import (
	"fmt"
	"os"

	gocsv "github.com/gocarina/gocsv"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/controllers"
	"github.com/ilovelili/dongfeng-jobs/services/server/core/models"
	"github.com/micro/cli"
)

// NutritionUpload nutrition file upload
func NutritionUpload(ctx *cli.Context) int {
	operationname := "upload nutrition csv file"

	if len(ctx.String("nutrition_file_path")) == 0 {
		errorlog(fmt.Sprintf("Error: %s", "must assign nutrition_file_path"), operationname)
		return 1
	}

	filepath := ctx.String("nutrition_file_path")
	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		errorlog(fmt.Sprintf("Error: failed to open %s", filepath), operationname)
		return 1
	}
	defer file.Close()

	nutritions := []*models.Nutrition{}
	if err := gocsv.UnmarshalFile(file, &nutritions); err != nil {
		errorlog(fmt.Sprintf("Error: invalid nutrition file %s", filepath), operationname)
		return 1
	}

	nutritioncontroller := controllers.NewNutritionController()
	err = nutritioncontroller.Save(nutritions)
	if err != nil {
		errorlog(fmt.Sprintf("Error on saving nutritions: %s", err), operationname)
		return 1
	}

	systemlog("job ended", operationname)
	return 0
}
