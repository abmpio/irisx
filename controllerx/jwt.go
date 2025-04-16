package controllerx

import (
	"sync"

	"github.com/abmpio/irisx/casdoor"
)

var (
	_casdoorOptions     casdoor.CasdoorOptions
	_casdoorM           *casdoor.CasdoorMiddleware
	_mustAuthenticatedM *casdoor.MustAuthenticated
	_sync               sync.Once
)

func GetCasdoorMiddleware() *casdoor.CasdoorMiddleware {
	_sync.Do(func() {
		_casdoorOptions = *casdoor.InitCasdoorSdk()
		_casdoorM = casdoor.NewCasdoorMiddleware(_casdoorOptions)
	})
	return _casdoorM
}

func GetMustAuthenticatedMiddleware() *casdoor.MustAuthenticated {
	_sync.Do(func() {
		_mustAuthenticatedM = casdoor.NewMustAuthenticated()
	})
	return _mustAuthenticatedM
}
