package core

import (
	"net/http"
)

// RegisterItemOrderRoutes registers the routes for item and order endpoints
func RegisterItemOrderRoutes(mux *http.ServeMux) {
	i := &items{}
	o := &orders{}

	// Register item routes
	mux.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			i.create(w, r)
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				i.get(w, r)
			} else {
				i.list(w, r)
			}
		case http.MethodPut:
			i.update(w, r)
		case http.MethodDelete:
			i.delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Register order routes
	mux.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			o.create(w, r)
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				o.get(w, r)
			} else {
				o.list(w, r)
			}
		case http.MethodPut:
			o.update(w, r)
		case http.MethodDelete:
			o.delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}