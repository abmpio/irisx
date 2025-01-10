package controllerx

import (
	"sync"

	"github.com/abmpio/configurationx"
	optCasdoor "github.com/abmpio/configurationx/options/casdoor"
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
		casdoorOpt := &optCasdoor.CasdoorOptions{}
		configurationx.GetInstance().UnmarshalPropertiesTo(optCasdoor.ConfigurationKey, casdoorOpt)
		_casdoorOptions = casdoor.CasdoorOptions{
			CasdoorOptions: *casdoorOpt,
			Extractor:      casdoor.FromFirst(casdoor.FromAuthHeader, casdoor.FromHeader("Authorization")),
		}
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
