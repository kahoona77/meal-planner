package wizard

import (
	"fmt"
	"log/slog"
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
	mealIdsWithTags, err := w.loadMealIdsWithTags()
	if err != nil {
		return nil, err
	}

	plannerWeek := &planner.Week{
		Start:  wizardWeek.Start,
		End:    wizardWeek.End,
		Offset: wizardWeek.Offset,
		Meals:  make([]*planner.MealOfTheDay, 0),
	}

	for _, day := range wizardWeek.Days {
		mod, err := w.getMealOfTheDay(day, mealIdsWithTags, plannerWeek.Meals)
		if err != nil {
			// the day keeps its current meal
			slog.Warn("could not find meal of the day", "date", day.Date, "err", err)
			continue
		}

		plannerWeek.Meals = append(plannerWeek.Meals, mod)
	}

	return plannerWeek, nil
}

// Reroll replaces the meal of the day by another random meal with the tags of the day.
func (w *Wizard) Reroll(day *Day, plannedMeals []*planner.MealOfTheDay) (*planner.MealOfTheDay, error) {
	mealIdsWithTags, err := w.loadMealIdsWithTags()
	if err != nil {
		return nil, err
	}

	mod, err := w.getMealOfTheDay(day, mealIdsWithTags, plannedMeals)
	if err != nil && len(day.Tags) > 0 {
		slog.Info("no other meal with the tags of the day, ignoring them", "date", day.Date, "err", err)
		mod, err = w.getMealOfTheDay(&Day{Weekday: day.Weekday, Date: day.Date}, mealIdsWithTags, plannedMeals)
	}
	return mod, err
}

// loadMealIdsWithTags returns the tags of every meal that has tags, by meal id.
func (w *Wizard) loadMealIdsWithTags() (map[int64][]*meals.MealTag, error) {
	allMealTags, err := w.mealsRepo.GetAllMealTags()
	if err != nil {
		return nil, err
	}

	mealIdsWithTags := make(map[int64][]*meals.MealTag)
	for _, tag := range allMealTags {
		mealIdsWithTags[tag.MealId] = append(mealIdsWithTags[tag.MealId], tag)
	}
	return mealIdsWithTags, nil
}

func (w *Wizard) getMealOfTheDay(day *Day, mealIdsWithTags map[int64][]*meals.MealTag, alreadyPickedMeals []*planner.MealOfTheDay) (*planner.MealOfTheDay, error) {
	taggedMealIds := w.findMealIdsByTags(mealIdsWithTags, day.Tags, alreadyPickedMeals)

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

func (w *Wizard) findMealIdsByTags(mealIdsWithTags map[int64][]*meals.MealTag, tags []*WeekdayTag, alreadyPickedMeals []*planner.MealOfTheDay) []int64 {
	result := make([]int64, 0)

	for mealId, tagsOfMeal := range mealIdsWithTags {
		if containsAllWeekDayTags(tags, tagsOfMeal) && !isAlreadyPicked(alreadyPickedMeals, mealId) {
			result = append(result, mealId)
		}
	}

	return result
}

func isAlreadyPicked(alreadyPickedMeals []*planner.MealOfTheDay, mealId int64) bool {
	for _, mod := range alreadyPickedMeals {
		if mod.Meal.Id == mealId {
			return true
		}
	}

	return false
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
