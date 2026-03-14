package gorouter

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

type ResponseWriter struct {
	http.ResponseWriter
	component templ.Component
}

func (w *ResponseWriter) Trigger(event string) {
	w.Header().Set("HX-Trigger", event)
}

func (w *ResponseWriter) Render(ctx context.Context) {
	w.component.Render(ctx, w.ResponseWriter)
}

// func (w ResponseWriter) Wrap(wrapper Wrapper)
