package controllerx

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/abmpio/entity"
	"github.com/abmpio/entity/filter"
	"github.com/abmpio/mongodbr"
	"github.com/abmpio/webserver/controller"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	webapp "github.com/abmpio/webserver/app"
)

type EntityController[T mongodbr.IEntity] struct {
	EntityService entity.IEntityService[T]

	// 针对本路由节点级别的中间件
	handlerList []context.Handler

	Options BaseEntityControllerOptions
	once    sync.Once
}

func NewEntityController[T mongodbr.IEntity](opts ...BaseEntityControllerOption) *EntityController[T] {
	options := BaseEntityControllerOptions{
		AllDisabled: true,
	}
	entityController := &EntityController[T]{
		Options: options,
	}
	for _, eachOpt := range opts {
		eachOpt(&(entityController.Options))
	}
	return entityController
}

func (c *EntityController[T]) RegistRouter(webapp *webapp.Application, opts ...BaseEntityControllerOption) router.Party {
	for _, eachOpt := range opts {
		eachOpt(&(c.Options))
	}

	c.handlerList = defaultContextHandlers(&c.Options.BaseControllerOptions)
	routerParty := webapp.Party(c.Options.RouterPath, c.handlerList...)

	if !c.Options.AllDisabled {
		routerParty.Get("/all", c.All)
	}
	if !c.Options.ListDisabled {
		routerParty.Get("/", c.GetList)
	}
	if !c.Options.SearchDiabled {
		routerParty.Post("/search", c.Search)
	}
	if !c.Options.GetByIdDisabled {
		routerParty.Get("/{id}", c.GetById)
	}
	if !c.Options.CreateDisabled {
		routerParty.Post("/", c.Create)
	}
	if !c.Options.UpdateDisabled {
		routerParty.Put("/{id}", c.Update)
	}
	if !c.Options.DeleteDisabled {
		routerParty.Delete("/{id}", c.Delete)
	}
	if !c.Options.DeleteListDisabled {
		routerParty.Delete("/", c.DeleteList)
	}

	return routerParty
}

func (c *EntityController[T]) MergeAuthenticatedContextIfNeed(authenticatedDisabled bool, handlers ...context.Handler) []context.Handler {
	return MergeAuthenticatedContextIfNeed(authenticatedDisabled, handlers...)
}

func (c *EntityController[T]) GetEntityService() entity.IEntityService[T] {
	c.once.Do(func() {
		if c.EntityService != nil {
			return
		}
		c.EntityService = GetEntityService[T]()
	})
	return c.EntityService
}

func (c *EntityController[T]) All(ctx iris.Context) {
	filter := map[string]interface{}{}

	if c.Options.ListFilterFunc != nil {
		c.Options.ListFilterFunc(new(T), filter, ctx)
	}
	var list []*T
	var err error
	if len(filter) > 0 {
		list, err = c.GetEntityService().FindList(filter)
	} else {
		list, err = c.GetEntityService().FindAll()
	}
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	controller.HandleSuccessWithListData(ctx, list, int64(len(list)))
}

func (c *EntityController[T]) GetList(ctx iris.Context) {
	all := filter.MustGetFilterAll(ctx.FormValue)
	if all {
		c.All(ctx)
		return
	}

	// params
	pagination := MustGetPagination(ctx)
	query := filter.MustGetFilterQuery(ctx.FormValue)
	sort := filter.MustGetSortOption(ctx.FormValue)

	service := c.GetEntityService()
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

func (c *EntityController[T]) Search(ctx iris.Context) {
	input := &SearchInput{}
	err := ctx.ReadJSON(input)
	if err != nil {
		controller.HandleErrorBadRequest(ctx, err)
		return
	}
	err = mongodbr.Validate(input)
	if err != nil {
		controller.HandleErrorBadRequest(ctx, err)
		return
	}

	findOptions := make([]mongodbr.FindOption, 0)
	findOptions = append(findOptions, mongodbr.FindOptionWithPage(int64(input.CurrentPage), int64(input.PageSize)))
	findOptions = append(findOptions, SetupFindOptionsWithSort(input.SortInput)...)
	service := c.GetEntityService()
	list, err := service.FindList(input.Filter, findOptions...)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}

	count, err := service.Count(input.Filter)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	controller.HandleSuccessWithTableData(ctx, list, count,
		controller.TableDataWithCurrentPage(input.CurrentPage),
		controller.TableDataWithPageSize(input.PageSize))
}

// get by id
func (c *EntityController[T]) GetById(ctx iris.Context) {
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
	item, err := c.GetEntityService().FindById(id)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	if item == nil {
		controller.HandleErrorInternalServerError(ctx, fmt.Errorf("invalid id,id:%s", idValue))
		return
	}
	controller.HandleSuccessWithData(ctx, item)
}

// create
func (c *EntityController[T]) Create(ctx iris.Context) {
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
	c.SetUserInfo(ctx, input)

	newItem, err := c.GetEntityService().Create(input)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	controller.HandleSuccessWithData(ctx, newItem)
}

// update
func (c *EntityController[T]) Update(ctx iris.Context) {
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
	service := c.GetEntityService()
	item, err := service.FindById(id)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	if item == nil {
		controller.HandleErrorBadRequest(ctx, fmt.Errorf("not found item,id:%s", idValue))
		return
	}

	input := make(map[string]interface{})
	err = ctx.ReadJSON(&input)
	if err != nil {
		controller.HandleErrorBadRequest(ctx, err)
		return
	}

	c.hookUpdate(ctx, input)
	err = service.UpdateFields(id, input)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	controller.HandleSuccess(ctx)
}

// delete
func (c *EntityController[T]) Delete(ctx iris.Context) {
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
	service := c.GetEntityService()
	item, err := service.FindById(oid)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	if item == nil {
		controller.HandleErrorBadRequest(ctx, fmt.Errorf("not found item,id:%s", idValue))
		return
	}

	err = c.GetEntityService().Delete(oid)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	controller.HandleSuccess(ctx)
}

// delete
func (c *EntityController[T]) DeleteList(ctx iris.Context) {
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

	_, err = c.GetEntityService().DeleteMany(filter)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	controller.HandleSuccess(ctx)
}

func (c *EntityController[T]) SetUserInfo(ctx iris.Context, entityValue interface{}) {
	userinfoProvider, ok := entityValue.(entity.IEntityWithUser)
	if !ok {
		return
	}
	userId := GetUserId(ctx)
	if userId != "" {
		userinfoProvider.SetUserCreator(userId)
	}
}

func (c *EntityController[T]) hookUpdate(ctx iris.Context, updated map[string]interface{}) {
	if len(updated) <= 0 {
		return
	}
	t := new(T)
	_, ok := interface{}(t).(mongodbr.IModificationEntity)
	if !ok {
		return
	}
	now := time.Now()
	updated["lastModificationTime"] = &now
	userId := GetUserId(ctx)
	if userId != "" {
		updated["lastModifierId"] = userId
	}
}

func SetupFindOptionsWithSort(i SortInput) []mongodbr.FindOption {
	opts := make([]mongodbr.FindOption, 0)
	if len(i.Sorts) <= 0 {
		return opts
	}

	for _, eachSort := range i.Sorts {
		if len(eachSort.Key) <= 0 || len(eachSort.Direction) <= 0 {
			continue
		}
		var isAsc bool
		if eachSort.Direction == entity.ASCENDING || eachSort.Direction == entity.ASC {
			isAsc = true
		} else if eachSort.Direction == entity.DESCENDING || eachSort.Direction == entity.DESC {
			isAsc = false
		} else {
			// invalid direction
			continue
		}
		opts = append(opts, mongodbr.FindOptionWithFieldSort(eachSort.Key, isAsc))
	}
	return opts
}
