// Package http handles GraphQL HTTP Requests including WebSocket Upgrades.
package http

import (
	"bytes"
	"fmt"
	"net/http"

	log "github.com/jensneuse/abstractlogger"

	"github.com/jensneuse/graphql-go-tools/pkg/graphql"
)

const (
	httpHeaderContentType          string = "Content-Type"
	httpContentTypeApplicationJson string = "application/json"
)

func (g *GraphQLHTTPRequestHandler) handleHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	var gqlRequest graphql.Request
	if err = graphql.UnmarshalHttpRequest(r, &gqlRequest, g.GetMaxMemory()); err != nil {
		g.log.Error("UnmarshalHttpRequest", log.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isIntrospection, err := gqlRequest.IsIntrospectionQuery()
	if err != nil {
		g.log.Error("IsIntrospectionQuery", log.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if isIntrospection {
		if err = g.schema.IntrospectionResponse(w); err != nil {
			g.log.Error("schema.IsIntrospectionQuery", log.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	resultWriter := graphql.NewEngineResultWriterFromBuffer(buf)

	w.Header().Add(httpHeaderContentType, httpContentTypeApplicationJson)
	w.WriteHeader(http.StatusOK)

	if err = g.engine.Execute(r.Context(), &gqlRequest, &resultWriter); err != nil {
		g.log.Error("engine.Execute", log.Error(err))
		_, _ = fmt.Fprintf(w, `{"errors":[{"message":"%s"}]}`, err)
		return
	}

	if _, err = w.Write(buf.Bytes()); err != nil {
		g.log.Error("write response", log.Error(err))
		return
	}
}
