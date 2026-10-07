package views

import (
	"meal-planner/core"
	"meal-planner/meals"
	"slices"
	"strconv"
)

// TagSelect is the data for the "tag-select" partial: one checkbox per tag,
// submitted as the repeated form field Name with the tag ids as values.
type TagSelect struct {
	Name    string
	Options []TagOption
}

type TagOption struct {
	Id       int64
	Name     string
	Color    string
	Selected bool
}

func newTagSelect(name string, tags []*meals.Tag, selectedIds []int64) TagSelect {
	options := make([]TagOption, len(tags))
	for i, tag := range tags {
		options[i] = TagOption{
			Id:       tag.Id,
			Name:     tag.Name,
			Color:    tag.Color,
			Selected: slices.Contains(selectedIds, tag.Id),
		}
	}
	return TagSelect{Name: name, Options: options}
}

// formTagIds returns the tag ids submitted by a "tag-select" with the given name.
func formTagIds(ctx *core.WebContext, name string) []int64 {
	values := ctx.FormValues(name)
	ids := make([]int64, 0, len(values))
	for _, value := range values {
		if id, err := strconv.ParseInt(value, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}
