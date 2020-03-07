package endpoint

import (
	"github.com/threeq/docs/swagger"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// Builder uses the builder pattern to generate swagger endpoint definitions
type Builder struct {
	Endpoint *swagger.Endpoint
}

// Option represents a functional option to customize the swagger endpoint
type Option func(builder *Builder)

// Apply improves the readability of applied options
func (o Option) Apply(builder *Builder) {
	o(builder)
}

// Handler allows an instance of the web handler to be associated with the endpoint.  This can be especially useful when
// using swag to bind the endpoints to the web router.  See the examples package for how the Handler can be used in
// conjunction with Walk to simplify binding endpoints to a router
func Handler(handler interface{}) Option {
	return func(b *Builder) {
		if v, ok := handler.(func(w http.ResponseWriter, r *http.Request)); ok {
			handler = http.HandlerFunc(v)
		}
		b.Endpoint.Handler = handler
	}
}

// Description sets the endpoint's description
func Description(v string) Option {
	return func(b *Builder) {
		b.Endpoint.Description = v
	}
}

// OperationID sets the endpoint's operationId
func OperationID(v string) Option {
	return func(b *Builder) {
		b.Endpoint.OperationID = v
	}
}

// Produces sets the endpoint's produces; by default this will be set to application/json
func Produces(v ...string) Option {
	return func(b *Builder) {
		b.Endpoint.Produces = v
	}
}

// Consumes sets the endpoint's produces; by default this will be set to application/json
func Consumes(v ...string) Option {
	return func(b *Builder) {
		b.Endpoint.Consumes = v
	}
}

func parameter(p swagger.Parameter) Option {
	return func(b *Builder) {
		if b.Endpoint.Parameters == nil {
			b.Endpoint.Parameters = []swagger.Parameter{}
		}

		b.Endpoint.Parameters = append(b.Endpoint.Parameters, p)
	}
}

// Path defines a path parameter for the endpoint; name, typ, description, and required correspond to the matching
// swagger fields
func Path(name, typ, description string, required bool) Option {
	p := swagger.Parameter{
		Name:        name,
		In:          "path",
		Type:        typ,
		Description: description,
		Required:    required,
	}
	return parameter(p)
}

// Query defines a query parameter for the endpoint; name, typ, description, and required correspond to the matching
// swagger fields
func Query(name, typ, description string, required bool) Option {
	p := swagger.Parameter{
		Name:        name,
		In:          "query",
		Type:        typ,
		Description: description,
		Required:    required,
	}
	return parameter(p)
}

// BodyType defines a body parameter for the swagger endpoint as would commonly be used for the POST, PUT, and PATCH methods
// prototype should be a struct or a pointer to struct that swag can use to reflect upon the return type
// t represents the Type of the body
func BodyType(t reflect.Type, description string, required bool) Option {
	p := swagger.Parameter{
		In:          "body",
		Name:        "body",
		Description: description,
		Schema:      swagger.MakeSchema("", t),
		Required:    required,
	}
	return parameter(p)
}

// Body defines a body parameter for the swagger endpoint as would commonly be used for the POST, PUT, and PATCH methods
// prototype should be a struct or a pointer to struct that swag can use to reflect upon the return type
func Body(prototype interface{}, description string, required bool) Option {
	return BodyType(reflect.TypeOf(prototype), description, required)
}

// Tags allows one or more tags to be associated with the endpoint
func Tags(tags ...string) Option {
	return func(b *Builder) {
		if b.Endpoint.Tags == nil {
			b.Endpoint.Tags = []string{}
		}

		b.Endpoint.Tags = append(b.Endpoint.Tags, tags...)
	}
}

// Security allows a security scheme to be associated with the endpoint.
func Security(scheme string, scopes ...string) Option {
	return func(b *Builder) {
		if b.Endpoint.Security == nil {
			b.Endpoint.Security = &swagger.SecurityRequirement{}
		}

		if b.Endpoint.Security.Requirements == nil {
			b.Endpoint.Security.Requirements = []map[string][]string{}
		}

		b.Endpoint.Security.Requirements = append(b.Endpoint.Security.Requirements, map[string][]string{scheme: scopes})
	}
}

// NoSecurity explicitly sets the endpoint to have no security requirements.
func NoSecurity() Option {
	return func(b *Builder) {
		b.Endpoint.Security = &swagger.SecurityRequirement{DisableSecurity: true}
	}
}

// ResponseOption allows for additional configurations on responses like header information
type ResponseOption func(response *swagger.Response)

// Apply improves the readability of applied options
func (o ResponseOption) Apply(response *swagger.Response) {
	o(response)
}

// Header adds header definitions to swagger responses
func Header(name, typ, format, description string) ResponseOption {
	return func(response *swagger.Response) {
		if response.Headers == nil {
			response.Headers = map[string]swagger.Header{}
		}

		response.Headers[name] = swagger.Header{
			Type:        typ,
			Format:      format,
			Description: description,
		}
	}
}

// ResponseType sets the endpoint response for the specified code; may be used multiple times with different status codes
// t represents the Type of the response
func ResponseType(code int, t reflect.Type, name string, description string, opts ...ResponseOption) Option {
	return func(b *Builder) {
		if b.Endpoint.Responses == nil {
			b.Endpoint.Responses = map[string]swagger.Response{}
		}

		r := swagger.Response{
			Description: description,
			Schema:      swagger.MakeSchema(name, t),
		}

		for _, opt := range opts {
			opt.Apply(&r)
		}

		b.Endpoint.Responses[strconv.Itoa(code)] = r
	}
}

// Response sets the endpoint response for the specified code; may be used multiple times with different status codes
func Response(code int, prototype interface{}, name string, description string, opts ...ResponseOption) Option {
	return func(b *Builder) {
		if b.Endpoint.Responses == nil {
			b.Endpoint.Responses = map[string]swagger.Response{}
		}

		r := swagger.Response{
			Description: description,
			Schema:      swagger.MakeSchema(name, prototype),
		}

		for _, opt := range opts {
			opt.Apply(&r)
		}

		b.Endpoint.Responses[strconv.Itoa(code)] = r
	}
}

// New constructs a new swagger endpoint using the fields and functional options provided
func New(method, path, summary string, options ...Option) *swagger.Endpoint {
	method = strings.ToUpper(method)
	e := &Builder{
		Endpoint: &swagger.Endpoint{
			Method:      method,
			Path:        path,
			Summary:     summary,
			OperationID: strings.ToLower(method) + camel(path),
			Produces:    []string{"application/json"},
			Consumes:    []string{"application/json"},
			Tags:        []string{},
		},
	}

	for _, opt := range options {
		opt.Apply(e)
	}

	return e.Endpoint
}

// Get constructs a new swagger [get] endpoint using the fields and functional options provided
func Get(path, summary string, options ...Option) *swagger.Endpoint {
	return New("get", path, summary, options...)
}

// Post constructs a new swagger [Post] endpoint using the fields and functional options provided
func Post(path, summary string, options ...Option) *swagger.Endpoint {
	return New("post", path, summary, options...)
}

// Put constructs a new swagger [Put] endpoint using the fields and functional options provided
func Put(path, summary string, options ...Option) *swagger.Endpoint {
	return New("put", path, summary, options...)
}

// Delete constructs a new swagger [Delete] endpoint using the fields and functional options provided
func Delete(path, summary string, options ...Option) *swagger.Endpoint {
	return New("delete", path, summary, options...)
}

// Patch constructs a new swagger [Patch] endpoint using the fields and functional options provided
func Patch(path, summary string, options ...Option) *swagger.Endpoint {
	return New("patch", path, summary, options...)
}

// Options constructs a new swagger [Options] endpoint using the fields and functional options provided
func Options(path, summary string, options ...Option) *swagger.Endpoint {
	return New("options", path, summary, options...)
}
