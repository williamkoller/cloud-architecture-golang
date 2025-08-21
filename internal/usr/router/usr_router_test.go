package usr_router

import (
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/cloud-architecture-golang/internal/usr/handler"
)

func TestRegisterUserRoutes_RegistersAllExpectedRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	api := r.Group("/api/v1")

	h := &handler.UserHandler{}

	RegisterUserRoutes(api, h)

	routes := r.Routes()

	type exp struct {
		method string
		path   string
		wantFn string 
	}
	expected := []exp{
		{method: "POST", path: "/api/v1/users", wantFn: ".CreateUser"},
		{method: "GET", path: "/api/v1/users", wantFn: ".ListUsers"},
		{method: "GET", path: "/api/v1/users/:email", wantFn: ".GetUser"},
		{method: "PATCH", path: "/api/v1/users/:email", wantFn: ".UpdateUser"},
		{method: "DELETE", path: "/api/v1/users/:email", wantFn: ".DeleteUser"},
	}

	for _, e := range expected {
		ri, ok := findRoute(routes, e.method, e.path)
		if !ok {
			t.Fatalf("route not found: %s %s", e.method, e.path)
		}
		
		if !strings.Contains(ri.Handler, e.wantFn) {
			t.Fatalf("handler mismatch for %s %s:\n got: %q\nwant to contain: %q",
				e.method, e.path, ri.Handler, e.wantFn)
		}
	}
}

func findRoute(routes []gin.RouteInfo, method, path string) (gin.RouteInfo, bool) {
	for _, r := range routes {
		if r.Method == method && r.Path == path {
			return r, true
		}
	}
	return gin.RouteInfo{}, false
}
