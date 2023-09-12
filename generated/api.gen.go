// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.14.0 DO NOT EDIT.
package generated

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Message string `json:"message"`
}

// HelloResponse defines model for HelloResponse.
type HelloResponse struct {
	Message string `json:"message"`
}

// LoginRequest defines model for LoginRequest.
type LoginRequest struct {
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	Id    int    `json:"id"`
	Token string `json:"token"`
}

// RegisterRequest defines model for RegisterRequest.
type RegisterRequest struct {
	FullName    string `json:"full_name"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

// RegisterResponse defines model for RegisterResponse.
type RegisterResponse struct {
	Id int `json:"id"`
}

// UpdateProfileRequest defines model for UpdateProfileRequest.
type UpdateProfileRequest struct {
	FullName    *string `json:"full_name,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}

// UpdateProfileResponse defines model for UpdateProfileResponse.
type UpdateProfileResponse struct {
	Message string `json:"message"`
}

// UserDetailResponse defines model for UserDetailResponse.
type UserDetailResponse struct {
	FullName    string `json:"full_name"`
	Id          *int   `json:"id,omitempty"`
	PhoneNumber string `json:"phone_number"`
}

// HelloParams defines parameters for Hello.
type HelloParams struct {
	Id int `form:"id" json:"id"`
}

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody = LoginRequest

// RegistrationJSONRequestBody defines body for Registration for application/json ContentType.
type RegistrationJSONRequestBody = RegisterRequest

// UpdateProfileJSONRequestBody defines body for UpdateProfile for application/json ContentType.
type UpdateProfileJSONRequestBody = UpdateProfileRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// This is just a test endpoint to get you started. Please delete this endpoint.
	// (GET /hello)
	Hello(ctx echo.Context, params HelloParams) error
	// Login endpoint.
	// (POST /login)
	Login(ctx echo.Context) error
	// User Registration Endpoint.
	// (POST /registration)
	Registration(ctx echo.Context) error
	// Update Profile Endpoint.
	// (PUT /users)
	UpdateProfile(ctx echo.Context) error
	// Get My Profile endpoint. Return User Data By ID.
	// (GET /users/{id})
	GetUserDetailByID(ctx echo.Context, id int) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// Hello converts echo context to params.
func (w *ServerInterfaceWrapper) Hello(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params HelloParams
	// ------------- Required query parameter "id" -------------

	err = runtime.BindQueryParameter("form", true, true, "id", ctx.QueryParams(), &params.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Hello(ctx, params)
	return err
}

// Login converts echo context to params.
func (w *ServerInterfaceWrapper) Login(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Login(ctx)
	return err
}

// Registration converts echo context to params.
func (w *ServerInterfaceWrapper) Registration(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Registration(ctx)
	return err
}

// UpdateProfile converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateProfile(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{"users:w"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UpdateProfile(ctx)
	return err
}

// GetUserDetailByID converts echo context to params.
func (w *ServerInterfaceWrapper) GetUserDetailByID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{"users:w"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetUserDetailByID(ctx, id)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/hello", wrapper.Hello)
	router.POST(baseURL+"/login", wrapper.Login)
	router.POST(baseURL+"/registration", wrapper.Registration)
	router.PUT(baseURL+"/users", wrapper.UpdateProfile)
	router.GET(baseURL+"/users/:id", wrapper.GetUserDetailByID)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xX32/bNhD+V4jbHg3LXfsyvS1N22VYh8JJsIcgKBjxbDGTSOZ4TKEZ+t8HUv6lRI7R",
	"Ijb8sKfI0pH38buPd18WUNjaWYOGPeQL8EWJtUyPH4gsTdE7azzGF46sQ2KN6XON3st5+sCNQ8jBM2kz",
	"h7YdAeFD0IQK8pt14O1oFWjv7rFgaEfwO1aVPXCOP+1cmyk+BPT8PIWT3n+zpAZyjMCV1uBXE+o7pP0g",
	"etGjzc4vgNp1cL2NRxvGOVJcxvYfNPuRaAWr2KHkU5xrz0g7SZmFqvpqZI3DrByMsk3ePfRtTvB9DD7n",
	"aWj3a6ck4xeyM13hj5K0l4d9aQ96Ka490jmy1NXuRC8fcJdEv08AvYpvr3wOuh2BxyKQ5uYydqkO5B1K",
	"QvotcLn59dFSLRly+OPvKxh1PS3u1H2F9c4ls4M2bqzNzMb1lS5wyUV3cPh8cZUunuYq/oy8iUukR11E",
	"yI9IXlsDObwZT8aTGGkdGuk05PA2vYpK5jJhzcrY7uLTHJOgItuStTUXCvKuGaZ4kjUykof8ZgE6bv8Q",
	"kBoYrVClC77hkSng8pxyUPa3Mbqrc0Lyy2QS/xTWMJoERTpX6SKBye69NZthEJ9+JpxBDj9lm2mRLUdF",
	"1u/hiU6FviDtuKPmCj0LQg5kIkHvJu9eLXd/Rg3k/suymNlgVKefUNeSmoip1F5oL+6DZyEFR4holLPa",
	"sGAr5siisUF4lsSoxuJLhdKjUFgho+C4fBU/TntnVWzp6RZZP1Dd1PGXVUPPZ1Y1r0ZDb8S1/TsWtdEe",
	"sPz9STZQgstQFOh9V/rJ8Up/JpVYc3IisktkPRUOpWHWSWW3fqbbUYeR0VNfcGQlPRvq/4vpRTGlWbQt",
	"C/GhL6zg0wxZgAsDgurZjQMpatBJHVlWw7Zqn7beHK++10ZyaUn/i6rL/fZ4uT9autNKoRHvrcIu/a/H",
	"S//emlmln7i75Hq2fd0NJCXn3+A2OpmtG5BKK5a1HZR/ttCq3Wm5PiFvnPBZc3G+w35FC3ea7mvAyJ+W",
	"tMPJaPuHJPYJWXxu1hJbj24xTYZWpCZ8LlmKs0ZcnI87IB7pcSWgQNXyf408yypbyKqMA769bf8LAAD/",
	"/1pb+c19EQAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
