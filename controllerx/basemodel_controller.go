package controllerx

import (
	"errors"
	"fmt"
	"sync"

	"abmp.cc/appserver/pkg/entity/filter"
	"abmp.cc/webserver/controller"
	"github.com/abmpio/entity"
	"github.com/abmpio/mongodbr"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	webapp "abmp.cc/webserver/app"
)

type ModelController[T mongodbr.IEntity] struct {
	RouterPath    string
	EntityService entity.IEntityService[T]

	options BaseEntityControllerOptions
	once    sync.Once
}

func (c *ModelController[T]) RegistRouter(webapp *webapp.Application, opts ...BaseEntityControllerOption) router.Party {
	for _, eachOpt := range opts {
		eachOpt(&(c.options))
	}
	handlerList := make([]context.Handler, 0)
	if !c.options.AuthenticatedDisabled {
		// handler auth
		handlerList = append(handlerList, getCasdoorMiddleware().Serve)
	}
	routerParty := webapp.Party(c.RouterPath, handlerList...)

	if !c.options.AllDisabled {
		routerParty.Get("/all", c.All)
	}
	if !c.options.ListDisabled {
		routerParty.Get("/", c.GetList)
	}
	if !c.options.GetByIdDisabled {
		routerParty.Get("/{id}", c.GetById)
	}
	if !c.options.CreateDisabled {
		routerParty.Post("/", c.Create)
	}
	if !c.options.UpdateDisabled {
		routerParty.Put("/{id}", c.Update)
	}
	if !c.options.DeleteDisabled {
		routerParty.Delete("/{id}", c.Delete)
	}
	if !c.options.DeleteListDisabled {
		routerParty.Delete("/", c.DeleteList)
	}

	return routerParty
}

func (c *ModelController[T]) getEntityService() entity.IEntityService[T] {
	c.once.Do(func() {
		if c.EntityService != nil {
			return
		}
		c.EntityService = GetEntityService[T]()
	})
	return c.EntityService
}

func (c *ModelController[T]) All(ctx iris.Context) {
	filter := map[string]interface{}{}
	// auto filter current userId
	AddUserIdFilterIfNeed(filter, new(T), ctx)

	var list []T
	var err error
	if len(filter) > 0 {
		list, err = c.getEntityService().FindList(filter)
	} else {
		list, err = c.getEntityService().FindAll()
	}
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	controller.HandleSuccessWithListData(ctx, list, int64(len(list)))
}

func (c *ModelController[T]) GetList(ctx iris.Context) {
	all := filter.MustGetFilterAll(ctx.FormValue)
	if all {
		c.All(ctx)
		return
	}

	// params
	pagination := MustGetPagination(ctx)
	query := filter.MustGetFilterQuery(ctx.FormValue)
	sort := filter.MustGetSortOption(ctx.FormValue)

	// auto filter current userId
	AddUserIdFilterIfNeed(query, new(T), ctx)
	service := c.getEntityService()
	list, err := service.FindList(query, mongodbr.FindOptionWithSort(sort),
		mongodbr.FindOptionWithPage(int64(pagination.Page), int64(pagination.Size)))
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}

	count, err := service.Count(query)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	controller.HandleSuccessWithListData(ctx, list, count)
}

// get by id
func (c *ModelController[T]) GetById(ctx iris.Context) {
	idValue := ctx.Params().Get("id")
	if len(idValue) <= 0 {
		controller.HandleErrorBadRequest(ctx, errors.New("id must not be empty"))
		return
	}

	id, err := primitive.ObjectIDFromHex(idValue)
	if err != nil {
		controller.HandleErrorBadRequest(ctx, fmt.Errorf("invalid id,id must be bson id format,id:%s", idValue))
		return
	}
	item, err := c.getEntityService().FindById(id)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	if item == nil {
		controller.HandleErrorInternalServerError(ctx, fmt.Errorf("invalid id,id:%s", idValue))
		return
	}
	// filter user is current user
	if !filterMustIsCurrentUserId(item, ctx) {
		controller.HandleErrorInternalServerError(ctx, fmt.Errorf("invalid id,id:%s", idValue))
		return
	}
	controller.HandleSuccessWithData(ctx, item)
}

// create
func (c *ModelController[T]) Create(ctx iris.Context) {
	input := new(T)
	err := ctx.ReadJSON(&input)
	if err != nil {
		controller.HandleErrorBadRequest(ctx, err)
		return
	}
	err = mongodbr.Validate(input)
	if err != nil {
		controller.HandleErrorBadRequest(ctx, err)
		return
	}

	// handler user info
	c.setUserInfo(ctx, input)

	newItem, err := c.getEntityService().Create(input)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	controller.HandleSuccessWithData(ctx, newItem)
}

// update
func (c *ModelController[T]) Update(ctx iris.Context) {
	idValue := ctx.Params().Get("id")
	if len(idValue) <= 0 {
		controller.HandleErrorBadRequest(ctx, errors.New("id must not be empty"))
		return
	}
	id, err := primitive.ObjectIDFromHex(idValue)
	if err != nil {
		controller.HandleErrorBadRequest(ctx, fmt.Errorf("invalid id,id must be bson id format,id:%s", idValue))
		return
	}
	service := c.getEntityService()
	item, err := service.FindById(id)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	if item == nil {
		controller.HandleErrorBadRequest(ctx, fmt.Errorf("not found item,id:%s", idValue))
		return
	}
	// filter user is current user
	if !filterMustIsCurrentUserId(item, ctx) {
		controller.HandleErrorInternalServerError(ctx, fmt.Errorf("invalid id,id:%s", idValue))
		return
	}

	input := make(map[string]interface{})
	err = ctx.ReadJSON(&input)
	if err != nil {
		controller.HandleErrorBadRequest(ctx, err)
		return
	}

	err = service.UpdateFields(id, input)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	controller.HandleSuccess(ctx)
}

// delete
func (c *ModelController[T]) Delete(ctx iris.Context) {
	idValue := ctx.Params().Get("id")
	if len(idValue) <= 0 {
		controller.HandleErrorBadRequest(ctx, errors.New("id must not be empty"))
		return
	}
	oid, err := primitive.ObjectIDFromHex(idValue)
	if err != nil {
		controller.HandleErrorBadRequest(ctx, fmt.Errorf("invalid id format,err:%s", err.Error()))
		return
	}
	service := c.getEntityService()
	item, err := service.FindById(oid)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	if item == nil {
		controller.HandleErrorBadRequest(ctx, fmt.Errorf("not found item,id:%s", idValue))
		return
	}
	// filter user is current user
	if !filterMustIsCurrentUserId(item, ctx) {
		controller.HandleErrorInternalServerError(ctx, fmt.Errorf("invalid id,id:%s", idValue))
		return
	}

	err = c.getEntityService().Delete(oid)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	controller.HandleSuccess(ctx)
}

// delete
func (c *ModelController[T]) DeleteList(ctx iris.Context) {
	payload, err := GetBatchRequestPayload(ctx)
	if err != nil {
		controller.HandleErrorBadRequest(ctx, err)
		return
	}
	if len(payload.Ids) <= 0 {
		controller.HandleSuccess(ctx)
		return
	}
	filter := bson.M{
		"_id": bson.M{"$in": payload.Ids},
	}
	// auto filter current userId
	AddUserIdFilterIfNeed(filter, new(T), ctx)

	_, err = c.getEntityService().DeleteMany(filter)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	controller.HandleSuccess(ctx)
}

func (c *ModelController[T]) setUserInfo(ctx iris.Context, entityValue interface{}) {
	userinfoProvider, ok := entityValue.(entity.IEntityWithUser)
	if !ok {
		return
	}
	userId := GetUserId(ctx)
	if userId != "" {
		userinfoProvider.SetUserCreator(userId)
	}
}
