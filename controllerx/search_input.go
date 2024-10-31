package controllerx

import (
	"github.com/abmpio/entity"
	"github.com/abmpio/webserver/controller"
	"github.com/kataras/iris/v12"
)

type SearchInput struct {
	controller.Pagination
	Filter map[string]interface{} `json:",inline"`

	SortInput
}

func (i *SearchInput) GetFilterValueAsString(key string) string {
	v, ok := i.Filter[key].(string)
	if ok {
		return v
	}
	return ""
}

func (i *SearchInput) ReadFromQuery(ctx iris.Context) error {
	err := ctx.ReadQuery(&i.Pagination)
	if err != nil {
		return err
	}
	err = ctx.ReadQuery(&i.SortInput)
	if err != nil {
		return err
	}
	return nil
}

type SortInput struct {
	Sorts []entity.Sort `json:"sorts" url:"sorts"`
}
