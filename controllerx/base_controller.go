package controllerx

import (
	webapp "github.com/abmpio/webserver/app"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
)

type BaseController struct {
	Options BaseControllerOptions
}

func NewBaseController(opts ...BaseControllerOption) *BaseController {
	baseController := &BaseController{}
	for _, eachOpt := range opts {
		eachOpt(&(baseController.Options))
	}
	return baseController
}

func (c *BaseController) RegistRouter(webapp *webapp.Application, opts ...BaseControllerOption) router.Party {
	for _, eachOpt := range opts {
		eachOpt(&(c.Options))
	}
	handlerList := defaultContextHandlers(&c.Options)
	routerParty := webapp.Party(c.Options.RouterPath, handlerList...)

	return routerParty
}

func defaultContextHandlers(o *BaseControllerOptions) []context.Handler {
	handlerList := make([]context.Handler, 0)
	// add authenticate middleware
	if !o.AuthenticatedDisabled {
		// handler auth
		handlerList = append(handlerList, GetCasdoorMiddleware().Serve)
	}
	return handlerList
}
