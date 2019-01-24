package models

import "strings"

// Menu menu entry
type Menu struct {
	ID                 int64  `dapper:"id,primarykey,autoincrement,table=menus"`
	Date               string `dapper:"date"`
	Recipe             string `dapper:"recipe"`
	BreakfastOrLunch   int64  `dapper:"breakfast_or_lunch"`
	JuniorOrSenior     int64  `dapper:"junior_or_senior_class"`
	JuniorBreakfastRaw string `dapper:"-"`
	JuniorLunchRaw     string `dapper:"-"`
	JuniorSnackRaw     string `dapper:"-"`
	SeniorBreakfastRaw string `dapper:"-"`
	SeniorLunchRaw     string `dapper:"-"`
	SeniorSnackRaw     string `dapper:"-"`
}

// ParseRecipes parse recipes
func (m *Menu) ParseRecipes() []*Menu {
	result := make([]*Menu, 0)
	recipes := parse(m.JuniorBreakfastRaw)
	for _, r := range recipes {
		result = append(result, &Menu{
			Date:             m.Date,
			BreakfastOrLunch: 0, // 0: breakfast 1: lunch 2: snack
			JuniorOrSenior:   0, // 0: junior 1: senior
			Recipe:           r,
		})
	}

	recipes = parse(m.JuniorLunchRaw)
	for _, r := range recipes {
		result = append(result, &Menu{
			Date:             m.Date,
			BreakfastOrLunch: 1, // 0: breakfast 1: lunch 2: snack
			JuniorOrSenior:   0, // 0: junior 1: senior
			Recipe:           r,
		})
	}

	recipes = parse(m.JuniorSnackRaw)
	for _, r := range recipes {
		result = append(result, &Menu{
			Date:             m.Date,
			BreakfastOrLunch: 2, // 0: breakfast 1: lunch 2: snack
			JuniorOrSenior:   0, // 0: junior 1: senior
			Recipe:           r,
		})
	}

	recipes = parse(m.SeniorBreakfastRaw)
	for _, r := range recipes {
		result = append(result, &Menu{
			Date:             m.Date,
			BreakfastOrLunch: 0, // 0: breakfast 1: lunch 2: snack
			JuniorOrSenior:   1, // 0: junior 1: senior
			Recipe:           r,
		})
	}

	recipes = parse(m.SeniorLunchRaw)
	for _, r := range recipes {
		result = append(result, &Menu{
			Date:             m.Date,
			BreakfastOrLunch: 1, // 0: breakfast 1: lunch 2: snack
			JuniorOrSenior:   1, // 0: junior 1: senior
			Recipe:           r,
		})
	}

	recipes = parse(m.SeniorSnackRaw)
	for _, r := range recipes {
		result = append(result, &Menu{
			Date:             m.Date,
			BreakfastOrLunch: 2, // 0: breakfast 1: lunch 2: snack
			JuniorOrSenior:   1, // 0: junior 1: senior
			Recipe:           r,
		})
	}

	return result
}

func parse(raw string) []string {
	return strings.Split(raw, "|")
}
