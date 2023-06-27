package controller

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

	BaseControllerOptions
}

type BaseEntityControllerOption func(*BaseEntityControllerOptions)

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
