package views

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"meal-planner/core"
	"meal-planner/files"
	"meal-planner/meals"
	"meal-planner/wizard"
	"net/http"
	"strconv"
)

func Meals(ctx *core.WebContext) error {
	repo := meals.NewRepository(ctx)
	mealsList, err := repo.GetMeals()
	if err != nil {
		logrus.Errorf("error loading meals: %v", err)
	}

	mealsWithTags := make([]*MealWithTags, len(mealsList))
	for i, meal := range mealsList {
		tags, err := repo.GetMealTags(meal.Id)
		if err != nil {
			return err
		}

		mealsWithTags[i] = &MealWithTags{
			Meal: meal,
			Tags: tags,
		}
	}

	data := core.TemplateData{
		"meals": mealsWithTags,
	}

	return ctx.RenderTemplate(http.StatusOK, "meals-list.html", data)
}

func MealEdit(ctx *core.WebContext) error {
	meal := &meals.Meal{}
	id := ctx.Param("id")

	repo := meals.NewRepository(ctx.Ctx)

	tags, err := repo.GetTags()
	if err != nil {
		return err
	}

	var mealTags []*meals.MealTag

	if id != "" {
		var err error
		meal, err = repo.GetMeal(int64(ctx.ParamAsInt("id")))
		if err != nil {
			return err
		}

		mealTags, err = repo.GetMealTags(meal.Id)
		if err != nil {
			return err
		}
	}

	return ctx.RenderTemplate(http.StatusOK, "meals-edit.html", core.TemplateData{
		"meal":     meal,
		"mealTags": toSelectOptions(mealTags, mealTagsConverter),
		"tags":     toSelectOptions(tags, tagsConverter),
	})
}

func MealSave(ctx *core.WebContext) error {
	repo := meals.NewRepository(ctx)

	meal := &meals.Meal{}
	isNew := ctx.Param("id") == ""
	if !isNew {
		var err error
		meal, err = meals.NewRepository(ctx.Ctx).GetMeal(int64(ctx.ParamAsInt("id")))
		if err != nil {
			return err
		}
	}

	//update
	meal.Name = ctx.FormValue("name")
	meal.Description = ctx.FormValue("description")

	var tags []*meals.Tag
	tagsJson := ctx.FormValue("tags")
	if tagsJson != "" {
		var options []*SelectOption
		err := json.Unmarshal([]byte(tagsJson), &options)
		if err != nil {
			return err
		}
		tags = make([]*meals.Tag, len(options))
		for i, option := range options {
			id, _ := strconv.ParseInt(option.Id, 10, 64)
			tags[i] = &meals.Tag{Id: id, Name: option.Name}
		}
	}

	imageFile, err := getNewImageFile(ctx, "image")
	if err == nil {
		filesRepo := files.NewRepository(ctx)
		if err := filesRepo.CreateFile(imageFile); err != nil {
			return err
		}

		//first delete previous file
		if meal.ImageFileId.Valid {
			if err := filesRepo.DeleteFile(meal.ImageFileId.Int64); err != nil {
				return err
			}
		}

		meal.ImageFileId = sql.NullInt64{Int64: imageFile.Id, Valid: true}
	}

	if isNew {
		if err := repo.CreateMeal(meal); err != nil {
			logrus.Errorf("error creating meal: %v", err)
			return err
		}
	} else {
		if err := repo.UpdateMeal(meal); err != nil {
			logrus.Errorf("error updating meal: %v", err)
			return err
		}
	}

	err = repo.SetTagsForMeal(meal, tags)
	if err != nil {
		return err
	}

	return ctx.Redirect(http.StatusFound, "/meals")
}

func MealDelete(ctx *core.WebContext) error {
	repo := meals.NewRepository(ctx)
	id := ctx.ParamAsInt("id")

	meal, err := repo.GetMeal(int64(id))
	if err != nil {
		return err
	}

	//first delete file
	if meal.ImageFileId.Valid {
		filesRepo := files.NewRepository(ctx)
		if err := filesRepo.DeleteFile(meal.ImageFileId.Int64); err != nil {
			return err
		}
	}

	if err := repo.DeleteMeal(meal.Id); err != nil {
		return err
	}

	return ctx.Redirect(http.StatusFound, "/meals")
}

func getNewImageFile(ctx *core.WebContext, name string) (*files.File, error) {
	file, err := ctx.FormFile(name)
	if err != nil {
		return nil, err
	}
	src, err := file.Open()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}

	if err := src.Close(); err != nil {
		return nil, err
	}

	imageFile := &files.File{
		Name:        file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Data:        data,
	}
	return imageFile, nil
}

type SelectOption struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type MealWithTags struct {
	*meals.Meal
	Tags []*meals.MealTag
}

type SelectOptionConverterFunc[T any] func(v T) *SelectOption

func toSelectOptions[T any](options []T, converterFunc SelectOptionConverterFunc[T]) []*SelectOption {
	result := make([]*SelectOption, len(options))
	for i, option := range options {
		result[i] = converterFunc(option)
	}

	return result
}

func tagsConverter(tag *meals.Tag) *SelectOption {
	return &SelectOption{
		Id:   fmt.Sprintf("%d", tag.Id),
		Name: tag.Name,
	}
}

func mealTagsConverter(mealTag *meals.MealTag) *SelectOption {
	return &SelectOption{
		Id:   fmt.Sprintf("%d", mealTag.TagId),
		Name: mealTag.Name,
	}
}

func weekdayTagsConverter(tag *wizard.WeekdayTag) *SelectOption {
	return &SelectOption{
		Id:   fmt.Sprintf("%d", tag.TagId),
		Name: tag.TagName,
	}
}
