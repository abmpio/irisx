package controllerx

import "github.com/kataras/iris/v12/context"

func MergeAuthenticatedContextIfNeed(authenticatedDisabled bool, handlers ...context.Handler) []context.Handler {
	handlerList := make([]context.Handler, 0)
	if !authenticatedDisabled {
		// handler auth
		handlerList = append(handlerList, GetCasdoorMiddleware().Serve)
	}
	handlerList = append(handlerList, handlers...)
	return handlerList
}
