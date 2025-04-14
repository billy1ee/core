package staticbackend

import (
	"encoding/json"
net/http"
	"strconv"
	"time

	"github.com/staticbackendhq/core/backend"
logger"
	"github.com/staticbackendhq/core/middleware"
	"github.com/staticbackendhq/core/model
)

// Item represents a product that is being sold
type Item struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"accountId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Category    string    `json:"category"`
	ImageURL    string    `json:"imageUrl,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type items struct {
	log *logger.Logger
}

// Collection name for items
// respond writes a JSON response with the given status code and data
func respond(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}



// Create a new item
func (i *items) create(w http.ResponseWriter, r *http.Request) {
	conf, auth, err := middleware.Extract(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Set defaults
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	// Convert to document map
	doc := map[string]interface{}{
		"name":        item.Name,
		"description": item.Description,
		"price":       item.Price,
		"stock":       item.Stock,
		"category":    item.Category,
		"imageUrl":    item.ImageURL,
		"createdAt":   item.CreatedAt,
		"updatedAt":   item.UpdatedAt,
	}

	inserted, err := backend.DB.CreateDocument(auth, conf.Name, itemsCollection, doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusCreated, inserted)
}

// Get a single item by ID
func (i *items) get(w http.ResponseWriter, r *http.Request) {
	conf, auth, err := middleware.Extract(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	doc, err := backend.DB.GetDocumentByID(auth, conf.Name, itemsCollection, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusOK, doc)
}

// List items with optional filtering
func (i *items) list(w http.ResponseWriter, r *http.Request) {
	conf, auth, err := middleware.Extract(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Setup paging params
	params := model.ListParams{
		Page: 1,
		Size: 50,
	}

	// Parse query parameters
	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			params.Page = p
		}
	}

	if size := r.URL.Query().Get("size"); size != "" {
		if s, err := strconv.Atoi(size); err == nil && s > 0 && s <= 100 {
			params.Size = s
		}
	}

	// Parse filters if provided
	var filters map[string]interface{}
	if filterParam := r.URL.Query().Get("filter"); filterParam != "" {
		if err := json.Unmarshal([]byte(filterParam), &filters); err != nil {
			http.Error(w, "Invalid filter format", http.StatusBadRequest)
			return
		}
	}

	var result model.PagedResult
	var err2 error

	// Apply filters if provided, otherwise get all
	if filters != nil && len(filters) > 0 {
		result, err2 = backend.DB.QueryDocuments(auth, conf.Name, itemsCollection, filters, params)
	} else {
		result, err2 = backend.DB.ListDocuments(auth, conf.Name, itemsCollection, params)
	}

	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusOK, result)
}

// Update an existing item
func (i *items) update(w http.ResponseWriter, r *http.Request) {
	conf, auth, err := middleware.Extract(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Add updatedAt timestamp
	updates["updatedAt"] = time.Now()

	updated, err := backend.DB.UpdateDocument(auth, conf.Name, itemsCollection, id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusOK, updated)
}

// Delete an item by ID
func (i *items) delete(w http.ResponseWriter, r *http.Request) {
	conf, auth, err := middleware.Extract(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	rowsAffected, err := backend.DB.DeleteDocument(auth, conf.Name, itemsCollection, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusOK, map[string]int64{"deleted": rowsAffected})
}

// UpdateStock updates the stock quantity of an item
func (i *items) updateStock(w http.ResponseWriter, r *http.Request) {
	conf, auth, err := middleware.Extract(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	var stockUpdate struct {
		Quantity int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&stockUpdate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Update the stock using IncrementValue for atomic update
	err = backend.DB.IncrementValue(auth, conf.Name, itemsCollection, id, "stock", stockUpdate.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	doc, err := backend.DB.GetDocumentByID(auth, conf.Name, itemsCollection, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusOK, doc)


// RegisterItemRoutes registers all Item-related routes
func RegisterItemRoutes(mux *http.ServeMux) {
	i := &items{
		log: logger.Get(),
	}
	
	mux.HandleFunc("/api/items/create", i.create)
	mux.HandleFunc("/api/items/get", i.get)
	mux.HandleFunc("/api/items/list", i.list)
	mux.HandleFunc("/api/items/update", i.update)
	mux.HandleFunc("/api/items/delete", i.delete)
	mux.HandleFunc("/api/items/updateStock", i.updateStock)
}