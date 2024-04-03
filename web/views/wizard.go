package views

import (
	"meal-planner/core"
	"meal-planner/meals"
	"meal-planner/planner"
	"meal-planner/wizard"
	"net/http"
	"time"
)

func Wizard(ctx *core.WebContext) error {
	offset := ctx.ParamAsInt("offset")

	tagsRepo := meals.NewRepository(ctx.Ctx)

	now := time.Now().AddDate(0, 0, offset*7)
	week := wizard.Wizard{
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
			Tags:    getWeekdayTags(int(date.Weekday()), weekdayTags),
		}
		date = date.AddDate(0, 0, 1)
	}

	return ctx.RenderTemplate(http.StatusOK, "wizard.html", core.TemplateData{
		"week": week,
		"tags": toSelectOptions(tags, tagsConverter),
	})
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
