package wizard

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand/v2"
	"meal-planner/core"
	"meal-planner/meals"
	"meal-planner/planner"
)

type Wizard struct {
	mealsRepo   *meals.Repository
	plannerRepo *planner.Repository
}

func NewWizard(c core.Context) *Wizard {
	return &Wizard{
		mealsRepo:   meals.NewRepository(c),
		plannerRepo: planner.NewRepository(c),
	}
}

func (w *Wizard) Generate(wizardWeek Week) (*planner.Week, error) {
	allMealTags, err := w.mealsRepo.GetAllMealTags()
	if err != nil {
		return nil, err
	}

	mealIdsWithTags := make(map[int64][]*meals.MealTag)
	for _, tag := range allMealTags {
		tagsOfMeal, exists := mealIdsWithTags[tag.MealId]
		if !exists {
			tagsOfMeal = make([]*meals.MealTag, 0)
		}

		tagsOfMeal = append(tagsOfMeal, tag)
		mealIdsWithTags[tag.MealId] = tagsOfMeal
	}

	plannerWeek := &planner.Week{
		Start:  wizardWeek.Start,
		End:    wizardWeek.End,
		Offset: wizardWeek.Offset,
		Meals:  make([]*planner.MealOfTheDay, 0),
	}

	for _, day := range wizardWeek.Days {
		mod, err := w.getMealOfTheDay(day, mealIdsWithTags)
		if err != nil {
			logrus.Warnf("could not find meal of the day: %v", err)
		}

		plannerWeek.Meals = append(plannerWeek.Meals, mod)
	}

	return plannerWeek, nil
}

func (w *Wizard) getMealOfTheDay(day *Day, mealIdsWithTags map[int64][]*meals.MealTag) (*planner.MealOfTheDay, error) {
	taggedMealIds := w.findMealIdsByTags(mealIdsWithTags, day.Tags)

	if len(taggedMealIds) <= 0 {
		return nil, fmt.Errorf("could to find meals for tags %v", day.Tags)
	}

	//pick random mealId
	pickedMealId := taggedMealIds[rand.IntN(len(taggedMealIds))]

	pickedMeal, err := w.mealsRepo.GetMeal(pickedMealId)
	if err != nil {
		return nil, fmt.Errorf("could get meal: %v", err)
	}

	mod := planner.NewMeal(day.Date)
	mod.SetMeal(pickedMeal)

	if err := w.plannerRepo.UpsertMealOfDay(mod); err != nil {
		return nil, fmt.Errorf("error while saving meal of the day: %v", err)
	}

	return mod, nil
}

func (w *Wizard) findMealIdsByTags(mealIdsWithTags map[int64][]*meals.MealTag, tags []*WeekdayTag) []int64 {
	result := make([]int64, 0)

	for mealId, tagsOfMeal := range mealIdsWithTags {
		if containsAllWeekDayTags(tags, tagsOfMeal) {
			result = append(result, mealId)
		}
	}

	return result
}

func containsAllWeekDayTags(tags []*WeekdayTag, tagsOfMeal []*meals.MealTag) bool {
	for _, weekdayTag := range tags {
		hasWeekdayTag := false
		for _, mealTag := range tagsOfMeal {
			if mealTag.Id == weekdayTag.TagId {
				// contains tag
				hasWeekdayTag = true
				continue
			}
		}

		if !hasWeekdayTag {
			return false
		}
	}

	return true
}
