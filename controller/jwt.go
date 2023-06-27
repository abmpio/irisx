package controller

import (
	"sync"

	"github.com/abmpio/configurationx"
	optCasdoor "github.com/abmpio/configurationx/options/casdoor"
	"github.com/abmpio/irisx/casdoor"
)

var (
	_casdoorOptions casdoor.CasdoorOptions
	_cm             *casdoor.Middleware
	_sync           sync.Once
)

func getCasdoorMiddleware() *casdoor.Middleware {
	_sync.Do(func() {
		casdoorOpt := &optCasdoor.CasdoorOptions{}
		configurationx.GetInstance().UnmarshalPropertiesTo(optCasdoor.ConfigurationKey, casdoorOpt)
		_casdoorOptions = casdoor.CasdoorOptions{
			CasdoorOptions: *casdoorOpt,
			Extractor: casdoor.FromFirst(casdoor.FromHeader("Authorization"),
				casdoor.FromAuthHeader),
		}
		_cm = casdoor.New(_casdoorOptions)
	})
	return _cm
}
