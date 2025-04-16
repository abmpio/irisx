package casdoor

import (
	"github.com/abmpio/configurationx"
	optCasdoor "github.com/abmpio/configurationx/options/casdoor"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

// init casdoorsdk
func InitCasdoorSdk(opts ...func(*CasdoorOptions)) *CasdoorOptions {
	casdoorOpt := &optCasdoor.CasdoorOptions{}
	configurationx.GetInstance().UnmarshalPropertiesTo(optCasdoor.ConfigurationKey, casdoorOpt)
	casdoorOptions := &CasdoorOptions{
		CasdoorOptions: *casdoorOpt,
		Extractor:      FromFirst(FromAuthHeader, FromHeader("Authorization")),
	}
	for _, eachOpt := range opts {
		eachOpt(casdoorOptions)
	}
	casdoorsdk.InitConfig(casdoorOptions.Endpoint,
		casdoorOptions.ClientId,
		casdoorOptions.ClientSecret,
		casdoorOptions.Certificate,
		casdoorOptions.OrganizationName,
		casdoorOptions.ApplicationName)

	if casdoorOptions.ErrorHandler == nil {
		casdoorOptions.ErrorHandler = OnError
	}

	if casdoorOptions.Extractor == nil {
		casdoorOptions.Extractor = FromAuthHeader
	}
	return casdoorOptions
}
