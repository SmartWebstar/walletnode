package rest

import (
	"context"
	"embed"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pastelnetwork/go-commons/log"
	swaggersvr "github.com/pastelnetwork/walletnode/rest/gen/http/swagger/server"
	ticketssvr "github.com/pastelnetwork/walletnode/rest/gen/http/tickets/server"
	tickets "github.com/pastelnetwork/walletnode/rest/gen/tickets"
	"github.com/pastelnetwork/walletnode/rest/services"

	goahttp "goa.design/goa/v3/http"
	goahttpmiddleware "goa.design/goa/v3/http/middleware"
	goamiddleware "goa.design/goa/v3/middleware"
)

//go:embed gen/http/openapi3.json
var openapi embed.FS

func httpHandler() http.Handler {
	var (
		dec        = goahttp.RequestDecoder
		enc        = goahttp.ResponseEncoder
		mux        = goahttp.NewMuxer()
		errHandler = errorHandler()
	)

	ticketsEndpoints := tickets.NewEndpoints(services.NewTickets())
	ticketsServer := ticketssvr.New(ticketsEndpoints, mux, dec, enc, errHandler, nil)
	ticketssvr.Mount(mux, ticketsServer)

	swaggerServer := swaggersvr.New(nil, mux, dec, enc, errHandler, nil)
	mountSwagger(mux, swaggerServer)

	servers := goahttp.Servers{
		ticketsServer,
		swaggerServer,
	}
	servers.Use(goahttpmiddleware.Debug(mux, os.Stdout))

	for _, m := range ticketsServer.Mounts {
		log.Infof("%s HTTP %q mounted on %s %s", logPrefix, m.Method, m.Verb, m.Pattern)
	}
	for _, m := range swaggerServer.Mounts {
		log.Infof("%s HTTP %q mounted on %s %s", logPrefix, m.Method, m.Verb, m.Pattern)
	}

	var handler http.Handler = mux

	handler = Log()(handler)
	handler = goahttpmiddleware.RequestID()(handler)

	return handler
}

func mountSwagger(mux goahttp.Muxer, server *swaggersvr.Server) {
	for _, m := range server.Mounts {
		file, err := openapi.Open(m.Method)
		if err != nil {
			continue
		}

		mux.Handle(m.Verb, m.Pattern, func(w http.ResponseWriter, r *http.Request) {
			http.ServeContent(w, r, m.Method, time.Time{}, file.(io.ReadSeeker))
		})

	}
}

// errorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func errorHandler() func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		id := ctx.Value(goamiddleware.RequestIDKey).(string)
		_, _ = w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		log.Errorf("%s [%s] %s", logPrefix, id, err.Error())
	}
}