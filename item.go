package staticbackend

import (
	"net/http"
	"time"

	"github.com/staticbackendhq/core/backend"
	"github.com/staticbackendhq/core/logger"
	"github.com/staticbackendhq/core/middleware"
	"github.com/staticbackendhq/core/model"
)

// Item represents a product that is sold
type Item struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"accountId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Inventory   int       `json:"inventory"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

// items handles all the CRUD operations for Item models
type items struct {
	log *logger.Logger
}

// RegisterItemRoutes registers the item routes with the provided mux
func RegisterItemRoutes(mux *http.ServeMux, log *logger.Logger) {
	h := &items{log: log}
	
	// Create and get all items
	mux.HandleFunc("/v1/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.create(w, r)
			return
		}
		h.list(w, r)
	})

	// Get, update, delete by ID
	mux.HandleFunc("/v1/items/", func(w http.ResponseWriter, r *http.Request) {
		// Extract ID from path
		id := r.URL.Path[len("/v1/items/"):]
		
		switch r.Method {
		case http.MethodGet:
			h.getByID(w, r, id)
		case http.MethodPut:
			h.update(w, r, id)
		case http.MethodDelete:
			h.delete(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

// create handles creating a new item
func (i *items) create(w http.ResponseWriter, r *http.Request) {
	conf, auth, err := middleware.Extract(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var item Item
	if err := parseBody(r.Body, &item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set timestamp fields
	item.Created = time.Now()
	item.Updated = time.Now()

	// Create the item using the generic database
	col := backend.Collection[Item](auth, conf, "items")
	inserted, err := col.Create(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusCreated, inserted)
}

// list handles listing all items with pagination
func (i *items) list(w http.ResponseWriter, r *http.Request) {
	conf, auth, err := middleware.Extract(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse query params for pagination
	page, size := 1, 10 // Default values
	if r.URL.Query().Get("page") != "" {
		if p, err := parseInt(r.URL.Query().Get("page")); err == nil && p > 0 {
			page = p
		}
	}
	if r.URL.Query().Get("size") != "" {
		if s, err := parseInt(r.URL.Query().Get("size")); err == nil && s > 0 {
			size = s
		}
	}

	lp := model.ListParams{
		Page: int64(page),
		Size: int64(size),
	}

	col := backend.Collection[Item](auth, conf, "items")
	result, err := col.List(lp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusOK, result)
}

// getByID handles retrieving an item by ID
func (i *items) getByID(w http.ResponseWriter, r *http.Request, id string) {
	conf, auth, err := middleware.Extract(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	col := backend.Collection[Item](auth, conf, "items")
	item, err := col.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusOK, item)
}

// update handles updating an existing item
func (i *items) update(w http.ResponseWriter, r *http.Request, id string) {
	conf, auth, err := middleware.Extract(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updates Item
	if err := parseBody(r.Body, &updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set updated timestamp
	updates.Updated = time.Now()

	col := backend.Collection[Item](auth, conf, "items")
	updated, err := col.Update(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusOK, updated)
}

// delete handles removing an item
func (i *items) delete(w http.ResponseWriter, r *http.Request, id string) {
	conf, auth, err := middleware.Extract(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	col := backend.Collection[Item](auth, conf, "items")
	count, err := col.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]int64{"deleted": count}
	respond(w, http.StatusOK, response)
}