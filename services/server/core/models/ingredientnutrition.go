package models

// IngredientNutrition ingredient nutrition entry
type IngredientNutrition struct {
	ID                int64   `dapper:"id,primarykey,autoincrement,table=ingredient_nutritions" csv:"-"`
	Ingredient        string  `dapper:"ingredient" csv:"ingredient"`
	Protein100g       float64 `dapper:"protein_100g" csv:"protein_100g"`
	ProteinDaily      float64 `dapper:"protein_daily" csv:"protein_daily"`
	Fat100g           float64 `dapper:"fat_100g" csv:"fat_100g"`
	FatDaily          float64 `dapper:"fat_daily" csv:"fat_daily"`
	Carbohydrate100g  float64 `dapper:"carbohydrate_100g" csv:"carbohydrate_100g"`
	CarbohydrateDaily float64 `dapper:"carbohydrate_daily" csv:"carbohydrate_daily"`
	Heat100g          float64 `dapper:"heat_100g" csv:"heat_100g"`
	HeatDaily         float64 `dapper:"heat_daily" csv:"heat_daily"`
	Calcium100g       float64 `dapper:"calcium_100g" csv:"calcium_100g"`
	CalciumDaily      float64 `dapper:"calcium_daily" csv:"calcium_daily"`
	Iron100g          float64 `dapper:"iron_100g" csv:"iron_100g"`
	IronDaily         float64 `dapper:"iron_daily" csv:"iron_daily"`
	Zinc100g          float64 `dapper:"zinc_100g" csv:"zinc_100g"`
	ZincDaily         float64 `dapper:"zinc_daily" csv:"zinc_daily"`
	VA100g            float64 `dapper:"va_100g" csv:"va_100g"`
	VADaily           float64 `dapper:"va_daily" csv:"va_daily"`
	VB1100g           float64 `dapper:"vb1_100g" csv:"vb1_100g"`
	VB1Daily          float64 `dapper:"vb1_daily" csv:"vb1_daily"`
	VB2100g           float64 `dapper:"vb2_100g" csv:"vb2_100g"`
	VB2Daily          float64 `dapper:"vb2_daily" csv:"vb2_daily"`
	VC100g            float64 `dapper:"vc_100g" csv:"vc_100g"`
	VCDaily           float64 `dapper:"vc_daily" csv:"vc_daily"`
}
