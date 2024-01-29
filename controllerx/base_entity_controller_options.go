package controllerx

import "github.com/kataras/iris/v12"

type BaseControllerOptions struct {
	AuthenticatedDisabled bool
}

type BaseEntityControllerOptions struct {
	AllDisabled        bool
	ListDisabled       bool
	GetByIdDisabled    bool
	CreateDisabled     bool
	UpdateDisabled     bool
	DeleteDisabled     bool
	DeleteListDisabled bool

	ListFilterFunc                   func(entityType interface{}, filter map[string]interface{}, ctx iris.Context)
	FilterCurrentUserForListDisabled bool

	BaseControllerOptions
}

type BaseEntityControllerOption func(*BaseEntityControllerOptions)

func BaseEntityControllerWithAllEndpointDisabled(v bool) BaseEntityControllerOption {
	return func(rro *BaseEntityControllerOptions) {
		rro.AllDisabled = v
		rro.ListDisabled = v
		rro.GetByIdDisabled = v
		rro.CreateDisabled = v
		rro.UpdateDisabled = v
		rro.DeleteDisabled = v
		rro.DeleteListDisabled = v
	}
}

func BaseEntityControllerWithAllDisabled(v bool) BaseEntityControllerOption {
	return func(rro *BaseEntityControllerOptions) {
		rro.AllDisabled = v
	}
}

func BaseEntityControllerWithListDisabled(v bool) BaseEntityControllerOption {
	return func(rro *BaseEntityControllerOptions) {
		rro.ListDisabled = v
	}
}

func BaseEntityControllerWithGetByIdDisabled(v bool) BaseEntityControllerOption {
	return func(rro *BaseEntityControllerOptions) {
		rro.GetByIdDisabled = v
	}
}

func BaseEntityControllerWithCreateDisabled(v bool) BaseEntityControllerOption {
	return func(rro *BaseEntityControllerOptions) {
		rro.CreateDisabled = v
	}
}

func BaseEntityControllerWithUpdateDisabled(v bool) BaseEntityControllerOption {
	return func(rro *BaseEntityControllerOptions) {
		rro.UpdateDisabled = v
	}
}

func BaseEntityControllerWithDeleteDisabled(v bool) BaseEntityControllerOption {
	return func(rro *BaseEntityControllerOptions) {
		rro.DeleteDisabled = v
	}
}

func BaseEntityControllerWithDeleteListDisabled(v bool) BaseEntityControllerOption {
	return func(rro *BaseEntityControllerOptions) {
		rro.DeleteListDisabled = v
	}
}

func BaseEntityControllerWithAuthenticatedDisabled(v bool) BaseEntityControllerOption {
	return func(rro *BaseEntityControllerOptions) {
		rro.AuthenticatedDisabled = v
	}
}
