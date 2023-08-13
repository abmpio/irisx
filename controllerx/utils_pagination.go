package controllerx

import (
	"github.com/abmpio/entity"
	"github.com/kataras/iris/v12"
)

func GetDefaultPagination() (p *entity.Pagination) {
	return &entity.Pagination{
		Page: entity.PaginationDefaultPage,
		Size: entity.PaginationDefaultSize,
	}
}

func GetPagination(ctx iris.Context) (p *entity.Pagination, err error) {
	var _p entity.Pagination

	if err := ctx.ReadQuery(&_p); err != nil {
		return GetDefaultPagination(), err
	}
	return &_p, nil
}

func MustGetPagination(ctx iris.Context) (p *entity.Pagination) {
	p, err := GetPagination(ctx)
	if err != nil || p == nil {
		return GetDefaultPagination()
	}
	return p
}

func GetBatchRequestPayload(ctx iris.Context) (payload entity.BatchRequestPayload, err error) {
	if err := ctx.ReadJSON(&payload); err != nil {
		return payload, err
	}
	return payload, err
}
