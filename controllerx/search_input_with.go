package controllerx

import (
	"github.com/abmpio/webserver/controller"
	"github.com/kataras/iris/v12"
)

type SearchInputWith[T any] struct {
	controller.Pagination
	Filter T `json:"Filter"`

	SortInput
}

type SearchInputWithMap = SearchInputWith[interface{}]

func (i *SearchInputWith[T]) ReadFromBody(ctx iris.Context) error {
	err := ctx.ReadBody(i)
	if err != nil {
		return err
	}
	return nil
}
