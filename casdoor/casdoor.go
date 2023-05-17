package casdoor

import (
	"errors"
	"fmt"
	"strings"

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

type Middleware struct {
	Options CasdoorOptions
}

func New(opts ...CasdoorOptions) *Middleware {
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

	return &Middleware{
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
			return "", nil
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

var (
	// ErrTokenMissing is the error value that it's returned when
	// a token is not found based on the token extractor.
	ErrTokenMissing = errors.New("required authorization token not found")
)

// Get returns the user (&token) information for this client/request
func (m *Middleware) Get(ctx iris.Context) *casdoorsdk.Claims {
	v := ctx.Values().Get(m.Options.Jwt.ContextKey)
	if v == nil {
		return nil
	}
	return v.(*casdoorsdk.Claims)
}

// Serve the middleware's action
func (m *Middleware) Serve(ctx iris.Context) {
	if err := m.CheckJWT(ctx); err != nil {
		m.Options.ErrorHandler(ctx, err)
		return
	}
	// If everything ok then call next.
	ctx.Next()
}

func (m *Middleware) CheckJWT(ctx iris.Context) error {

	// Use the specified token extractor to extract a token from the request
	token, err := m.Options.Extractor(ctx)
	// If debugging is turned on, log the outcome
	if err != nil {
		logf(ctx, "Error extracting JWT: %v", err)
		return err
	}

	logf(ctx, "Token extracted: %s", token)

	// If the token is empty...
	if token == "" {
		// Check if it was required
		if m.Options.Jwt.CredentialsOptional {
			logf(ctx, "No credentials found (CredentialsOptional=true)")
			// No error, just no token (and that is ok given that CredentialsOptional is true)
			return nil
		}

		// If we get here, the required token is missing
		logf(ctx, "Error: No credentials found (CredentialsOptional=false)")
		return ErrTokenMissing
	}

	// Now parse the token

	claim, err := casdoorsdk.ParseJwtToken(token)
	// Check if there was an error in parsing...
	if err != nil {
		logf(ctx, "Error parsing token: %v", err)
		return err
	}

	logf(ctx, "claim: %v", claim)

	// If we get here, everything worked and we can set the
	// user property in context.
	ctx.Values().Set(m.Options.Jwt.ContextKey, claim)

	return nil
}

func (m *Middleware) GetUserClaims(ctx iris.Context) *casdoorsdk.Claims {
	claims, ok := ctx.Value(m.Options.Jwt.ContextKey).(*casdoorsdk.Claims)
	if !ok {
		return nil
	}
	return claims
}
