package views

import (
	"encoding/json"
	"fmt"
	"meal-planner/core"
	"meal-planner/meals"
	"meal-planner/planner"
	"meal-planner/wizard"
	"net/http"
	"strconv"
	"time"
)

func Wizard(ctx *core.WebContext) error {
	offset := ctx.ParamAsInt("offset")

	tagsRepo := meals.NewRepository(ctx.Ctx)

	now := time.Now().AddDate(0, 0, offset*7)
	week := wizard.Week{
		Start:  planner.GetStartWeek(now),
		End:    planner.GetEndWeek(now),
		Offset: offset,
	}

	tags, err := tagsRepo.GetTags()
	if err != nil {
		return err
	}

	weekdayTags, err := loadTags(ctx, tags)
	if err != nil {
		return err
	}

	week.Days = make([]*wizard.Day, 7)
	date := week.Start
	for i := 0; i < 7; i++ {
		week.Days[i] = &wizard.Day{
			Weekday: int(date.Weekday()),
			Date:    date,
			Tags:    getWeekdayTags(int(date.Weekday()), weekdayTags),
		}
		date = date.AddDate(0, 0, 1)
	}

	return ctx.RenderTemplate(http.StatusOK, "wizard.html", core.TemplateData{
		"week": map[string]any{
			"start": week.Start,
			"end":   week.End,
			"days": func(days []*wizard.Day) []map[string]any {
				result := make([]map[string]any, len(days))
				for i, day := range days {
					result[i] = map[string]any{
						"weekday": day.Weekday,
						"tags":    toSelectOptions(day.Tags, weekdayTagsConverter),
					}
				}
				return result
			}(week.Days),
		},
		"tags": toSelectOptions(tags, tagsConverter),
	})
}

func Generate(ctx *core.WebContext) error {
	offset := ctx.ParamAsInt("offset")

	tagsRepo := meals.NewRepository(ctx.Ctx)
	wizardRepo := wizard.NewRepository(ctx.Ctx)

	now := time.Now().AddDate(0, 0, offset*7)
	week := wizard.Week{
		Start:  planner.GetStartWeek(now),
		End:    planner.GetEndWeek(now),
		Offset: offset,
	}

	week.Days = make([]*wizard.Day, 7)
	date := week.Start
	for i := 0; i < 7; i++ {
		weekday := int(date.Weekday())

		// get tags for weekday
		weekdayTags := make([]*wizard.WeekdayTag, 0)

		tagsJson := ctx.FormValue(fmt.Sprintf("tags_%d", weekday))
		if tagsJson != "" {
			var options []*SelectOption
			err := json.Unmarshal([]byte(tagsJson), &options)
			if err != nil {
				return err
			}
			for _, option := range options {
				id, _ := strconv.ParseInt(option.Id, 10, 64)
				tag, err := tagsRepo.GetTag(id)
				if err != nil {
					return err
				}

				weekdayTags = append(weekdayTags, &wizard.WeekdayTag{
					Weekday:  weekday,
					TagId:    tag.Id,
					TagName:  tag.Name,
					TagColor: tag.Color,
				})
			}
		}

		// save new weekday-tags
		if err := wizardRepo.SetWeekdayTags(weekday, weekdayTags); err != nil {
			return err
		}

		week.Days[i] = &wizard.Day{
			Weekday: int(date.Weekday()),
			Date:    date,
			Tags:    weekdayTags,
		}
		date = date.AddDate(0, 0, 1)
	}

	wizzard := wizard.NewWizard(ctx)
	plannerWeek, err := wizzard.Generate(week)
	if err != nil {
		return err
	}

	return ctx.Redirect(http.StatusFound, fmt.Sprintf("/offset/%d", plannerWeek.Offset))
}

func loadTags(ctx *core.WebContext, tags []*meals.Tag) ([]*wizard.WeekdayTag, error) {
	wizardRepo := wizard.NewRepository(ctx.Ctx)

	weekdayTags, err := wizardRepo.GetWeekdayTags()
	if err != nil {
		return nil, err
	}

	for _, weekdayTag := range weekdayTags {
		for _, tag := range tags {
			if weekdayTag.TagId == tag.Id {
				weekdayTag.TagName = tag.Name
				weekdayTag.TagColor = tag.Color
			}
		}
	}

	return weekdayTags, nil
}

func getWeekdayTags(weekday int, tags []*wizard.WeekdayTag) []*wizard.WeekdayTag {
	result := make([]*wizard.WeekdayTag, 0)
	for _, tag := range tags {
		if tag.Weekday == weekday {
			result = append(result, tag)
		}
	}

	return result
}
