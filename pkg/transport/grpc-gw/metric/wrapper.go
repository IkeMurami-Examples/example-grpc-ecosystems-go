// Package metric is a helper package to get a standard `runtime.HandlerFunc` compatible middleware.
package metric

import (
	"fmt"
	"net/http"
	"path"
	"reflect"
	"unsafe"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/slok/go-http-metrics/middleware"
)

// ServMuxWrapper wrap handlers from runtime.ServeMux functions inplace
func ServMuxWrapper(mux *runtime.ServeMux, mdlw middleware.Middleware, pfx string) (*runtime.ServeMux, error) {
	val := reflect.ValueOf(mux).Elem()
	handlersMap := val.FieldByName("handlers")
	if handlersMap.Kind() != reflect.Map {
		return nil, fmt.Errorf("Error getting handlers")
	}
	for _, key := range handlersMap.MapKeys() {
		handlers := handlersMap.MapIndex(key)
		if handlers.Kind() != reflect.Slice {
			return nil, fmt.Errorf("Error type handlers")
		}
		for i := 0; i < handlers.Len(); i++ {
			h := handlers.Index(i)
			refl := h.FieldByName("pat")
			refl = reflect.NewAt(refl.Type(), unsafe.Pointer(refl.UnsafeAddr()))
			pattern, ok := refl.Elem().Interface().(runtime.Pattern)
			if !ok {
				return nil, fmt.Errorf("Error cast pattern for handlers")
			}
			refl = h.FieldByName("h")
			refl = reflect.NewAt(refl.Type(), unsafe.Pointer(refl.UnsafeAddr()))
			handlerFunc, ok := refl.Elem().Interface().(runtime.HandlerFunc)
			if !ok {
				return nil, fmt.Errorf("Error cast handlerfunc")
			}
			handlerFunc = HandlerWrapper(path.Join(pfx, pattern.String()), mdlw, handlerFunc)
			refl.Elem().Set(reflect.ValueOf(handlerFunc))
		}
	}
	return mux, nil
}

// HandlerWrapper returns an measuring standard runtime.HandlerFunc.
func HandlerWrapper(handlerID string, m middleware.Middleware, h runtime.HandlerFunc) runtime.HandlerFunc {
	return runtime.HandlerFunc(func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		wi := &responseWriterInterceptor{
			statusCode:     http.StatusOK,
			ResponseWriter: w,
		}
		reporter := &stdReporter{
			w: wi,
			r: r,
		}

		m.Measure(handlerID, reporter, func() {
			h(wi, r, pathParams)
		})
	})
}
