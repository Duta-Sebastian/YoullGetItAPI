package middleware

import "net/http"

func ChainMiddleware(router *http.ServeMux, middlewares ...func(*http.ServeMux) http.Handler) http.Handler {
	var handler http.Handler = router

	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		if middleware != nil {
			middlewareFunc := func(next http.Handler) http.Handler {
				return middleware(router)
			}
			handler = middlewareFunc(handler)
		}
	}

	return handler
}
