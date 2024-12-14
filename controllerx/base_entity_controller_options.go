package controllerx

import "github.com/kataras/iris/v12"

type BaseControllerOptions struct {
	RouterPath            string
	AuthenticatedDisabled bool
}

type BaseControllerOption func(*BaseControllerOptions)

// set controller's router path
func BaseControllerWithRouterPath(routerPath string) BaseControllerOption {
	return func(bco *BaseControllerOptions) {
		bco.RouterPath = routerPath
	}
}

// set controller's router path
func BaseControllerWithAuthenticatedDisabled(authenticatedDisabled bool) BaseControllerOption {
	return func(bco *BaseControllerOptions) {
		bco.AuthenticatedDisabled = authenticatedDisabled
	}
}

type BaseEntityControllerOptions struct {
	AllDisabled        bool
	ListDisabled       bool
	SearchDiabled      bool
	GetByIdDisabled    bool
	CreateDisabled     bool
	UpdateDisabled     bool
	DeleteDisabled     bool
	DeleteListDisabled bool

	ListFilterFunc func(entityType interface{}, filter map[string]interface{}, ctx iris.Context)

	BaseControllerOptions
}

type BaseEntityControllerOption func(*BaseEntityControllerOptions)

func BaseEntityControllerWithAllEndpointDisabled(v bool) BaseEntityControllerOption {
	return func(rro *BaseEntityControllerOptions) {
		rro.AllDisabled = v
		rro.ListDisabled = v
		rro.SearchDiabled = v
		rro.GetByIdDisabled = v
		rro.CreateDisabled = v
		rro.UpdateDisabled = v
		rro.DeleteDisabled = v
		rro.DeleteListDisabled = v
	}
}

// set controller's router path
func BaseEntityControllerWithRouterPath(routerPath string) BaseEntityControllerOption {
	return func(beco *BaseEntityControllerOptions) {
		beco.RouterPath = routerPath
	}
}

func BaseEntityControllerWithAuthenticatedDisabled(v bool) BaseEntityControllerOption {
	return func(beco *BaseEntityControllerOptions) {
		beco.AuthenticatedDisabled = v
	}
}
func BaseEntityControllerWithAllDisabled(v bool) BaseEntityControllerOption {
	return func(beco *BaseEntityControllerOptions) {
		beco.AllDisabled = v
	}
}

func BaseEntityControllerWithListDisabled(v bool) BaseEntityControllerOption {
	return func(beco *BaseEntityControllerOptions) {
		beco.ListDisabled = v
	}
}

func BaseEntityControllerWithSearchDisabled(v bool) BaseEntityControllerOption {
	return func(beco *BaseEntityControllerOptions) {
		beco.SearchDiabled = v
	}
}

func BaseEntityControllerWithGetByIdDisabled(v bool) BaseEntityControllerOption {
	return func(beco *BaseEntityControllerOptions) {
		beco.GetByIdDisabled = v
	}
}

func BaseEntityControllerWithCreateDisabled(v bool) BaseEntityControllerOption {
	return func(beco *BaseEntityControllerOptions) {
		beco.CreateDisabled = v
	}
}

func BaseEntityControllerWithUpdateDisabled(v bool) BaseEntityControllerOption {
	return func(beco *BaseEntityControllerOptions) {
		beco.UpdateDisabled = v
	}
}

func BaseEntityControllerWithDeleteDisabled(v bool) BaseEntityControllerOption {
	return func(beco *BaseEntityControllerOptions) {
		beco.DeleteDisabled = v
	}
}

func BaseEntityControllerWithDeleteListDisabled(v bool) BaseEntityControllerOption {
	return func(beco *BaseEntityControllerOptions) {
		beco.DeleteListDisabled = v
	}
}
