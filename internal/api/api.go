package api

import (
	"net/http"
	"net/http/pprof"

	"l4.5/internal/api/controllers"
)

func RouteController(c *controllers.Controller) http.Handler {
	mux := http.NewServeMux()
	// Для простоты сервер будет раздавать и статику
	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/info/", http.StripPrefix("/info/", fs))
	mux.HandleFunc("/order/", c.GetOrder)

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return mux
}
