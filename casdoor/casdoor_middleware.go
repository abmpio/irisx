package casdoor

import (
	"fmt"
	"strings"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/configurationx/options/casdoor"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/kataras/iris/v12"
)

// A function called whenever an error is encountered
type errorHandler func(iris.Context, error)

// TokenExtractor is a function that takes a context as input and returns
// either a token or an error.  An error should only be returned if an attempt
// to specify a token was found, but the information was somehow incorrectly
// formed.  In the case where a token is simply not present, this should not
// be treated as an error.  An empty string should be returned in that case.
type TokenExtractor func(iris.Context) (string, error)

type CasdoorOptions struct {
	casdoor.CasdoorOptions

	// The function that will be called when there's an error validating the token
	// Default value:
	ErrorHandler errorHandler
	// A function that extracts the token from the request
	// Default: FromAuthHeader (i.e., from Authorization header as bearer token)
	Extractor TokenExtractor
}

// set useId to context
type CasdoorMiddleware struct {
	Options CasdoorOptions
}

func NewCasdoorMiddleware(opts ...CasdoorOptions) *CasdoorMiddleware {
	var options CasdoorOptions
	if len(opts) == 0 {
		options = CasdoorOptions{}
	} else {
		options = opts[0]
	}
	options.Normalize()
	if !options.Disabled {
		casdoorsdk.InitConfig(options.Endpoint,
			options.ClientId,
			options.ClientSecret,
			options.Certificate,
			options.OrganizationName,
			options.ApplicationName)
	}

	if options.ErrorHandler == nil {
		options.ErrorHandler = OnError
	}

	if options.Extractor == nil {
		options.Extractor = FromAuthHeader
	}

	return &CasdoorMiddleware{
		Options: options,
	}
}

// OnError is the default error handler.
// Use it to change the behavior for each error.
// See `Config.ErrorHandler`.
func OnError(ctx iris.Context, err error) {
	if err == nil {
		return
	}

	ctx.StopExecution()
	ctx.StatusCode(iris.StatusUnauthorized)
	ctx.WriteString(err.Error())
}

// FromAuthHeader is a "TokenExtractor" that takes a give context and extracts
// the JWT token from the Authorization header.
func FromAuthHeader(ctx iris.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

// FromHeader is a "TokenExtractor" that takes a give context and extracts
// the specified key value from header.
func FromHeader(key string) TokenExtractor {
	return func(ctx iris.Context) (string, error) {
		headerValue := ctx.GetHeader(key)
		if headerValue == "" {
			return "", nil // No error, just no token
		}
		authHeaderParts := strings.Split(headerValue, " ")
		if len(authHeaderParts) > 1 && strings.ToLower(authHeaderParts[0]) == "bearer" {
			return authHeaderParts[1], nil
		}
		return headerValue, nil
	}
}

// FromParameter returns a function that extracts the token from the specified
// query string parameter
func FromParameter(param string) TokenExtractor {
	return func(ctx iris.Context) (string, error) {
		return ctx.URLParam(param), nil
	}
}

// FromFirst returns a function that runs multiple token extractors and takes the
// first token it finds
func FromFirst(extractors ...TokenExtractor) TokenExtractor {
	return func(ctx iris.Context) (string, error) {
		for _, ex := range extractors {
			token, err := ex(ctx)
			if err != nil {
				return "", err
			}
			if token != "" {
				return token, nil
			}
		}
		return "", nil
	}
}

func logf(ctx iris.Context, format string, args ...interface{}) {
	ctx.Application().Logger().Debugf(format, args...)
}

// Get returns the user (&token) information for this client/request
func (m *CasdoorMiddleware) Get(ctx iris.Context) *casdoorsdk.Claims {
	v := ctx.Values().Get(m.Options.Jwt.ContextKey)
	if v == nil {
		return nil
	}
	return v.(*casdoorsdk.Claims)
}

// Serve the middleware's action
func (m *CasdoorMiddleware) Serve(ctx iris.Context) {
	if err := m.CheckJWT(ctx); err != nil {
		m.Options.ErrorHandler(ctx, err)
		return
	}
	// If everything ok then call next.
	ctx.Next()
}

func (m *CasdoorMiddleware) CheckJWT(ctx iris.Context) error {
	// is authenticated by other middleware?
	user := m.GetUserClaims(ctx)
	if user != nil {
		return nil
	}
	// Use the specified token extractor to extract a token from the request
	token, err := m.Options.Extractor(ctx)
	// If debugging is turned on, log the outcome
	if err != nil {
		log.Logger.Warn(fmt.Sprintf("Error extracting JWT: %v", err))
		return err
	}

	logf(ctx, "Token extracted: %s", token)

	// If the token is empty...
	if token == "" {
		return nil
	}

	// Check if it was required
	if m.Options.Jwt.CredentialsOptional {
		log.Logger.Debug("No credentials found (CredentialsOptional=true)")
		// No error, just no token (and that is ok given that CredentialsOptional is true)
		return nil
	}

	// Now parse the token
	claim, err := casdoorsdk.ParseJwtToken(token)
	// Check if there was an error in parsing...
	if err != nil {
		log.Logger.Warn(fmt.Sprintf("Error parsing token: %v", err))
		return err
	}

	logf(ctx, "claim: %v", claim)

	// If we get here, everything worked and we can set the
	// user property in context.
	ctx.Values().Set(m.Options.Jwt.ContextKey, claim)
	if claim != nil {
		ctx.Values().Set("userId", claim.Id)
	}

	return nil
}

func (m *CasdoorMiddleware) GetUserClaims(ctx iris.Context) *casdoorsdk.Claims {
	v := ctx.Value(m.Options.Jwt.ContextKey)
	if v == nil {
		return nil
	}
	claims, ok := v.(*casdoorsdk.Claims)
	if !ok {
		return nil
	}
	return claims
}
