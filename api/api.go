package api

import (
	"fmt"
	"net/http"

	"github.com/pshevtsov/verigo"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func gae(next verigo.ContextHandler) verigo.ContextHandler {
	return verigo.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		ctx = appengine.NewContext(req)
		next.ServeHTTPContext(ctx, rw, req)
	})
}

func token(next verigo.ContextHandler) verigo.ContextHandler {
	return verigo.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		log.Infof(ctx, "Token Middleware")
		ctx = context.WithValue(ctx, "token", "ok")
		next.ServeHTTPContext(ctx, rw, req)
	})
}

func acl(next verigo.ContextHandler) verigo.ContextHandler {
	return verigo.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		log.Infof(ctx, "ACL Middleware")
		ctx = context.WithValue(ctx, "acl", 1)
		next.ServeHTTPContext(ctx, rw, req)
	})
}

func validation(next verigo.ContextHandler) verigo.ContextHandler {
	return verigo.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		log.Infof(ctx, "Validation Middleware")
		if req.FormValue("ok") == "false" {
			http.Error(rw, http.StatusText(400), http.StatusBadRequest)
			return
		}
		ctx = context.WithValue(ctx, "valid", true)
		next.ServeHTTPContext(ctx, rw, req)
	})
}

func app(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(rw, req)
		return
	}
	fmt.Fprintf(
		rw,
		"Hello, world!\n\ntoken: %#v (%[1]T)\nacl: %#v (%[2]T)\nvalid: %#v (%[3]T)\n",
		ctx.Value("token"),
		ctx.Value("acl"),
		ctx.Value("valid"),
	)
}

func init() {
	chain := verigo.New(gae, token, acl, validation).Then(app)
	http.Handle("/", chain)
}
