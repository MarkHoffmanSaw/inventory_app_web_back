package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type ResponseJSON struct {
	Message string `json:"message"`
}

func main() {
	router := mux.NewRouter()
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	// Routes
	router.HandleFunc("/customers", createCustomerHandler).Methods("POST")
	router.HandleFunc("/customers", getCustomersHandler).Methods("GET")

	router.HandleFunc("/materials", createMaterialHandler).Methods("POST")
	router.HandleFunc("/materials", getMaterialsHandler).Methods("GET")
	router.HandleFunc("/material_types", getMaterialTypesHandler).Methods("GET")
	router.HandleFunc("/materials/move-to-location", moveMaterialHandler).Methods("PATCH")
	router.HandleFunc("/materials/remove-from-location", removeMaterialHandler).Methods("PATCH")

	router.HandleFunc("/incoming_materials", sendMaterialHandler).Methods("POST")
	router.HandleFunc("/incoming_materials", getIncomingMaterialsHandler).Methods("GET")

	router.HandleFunc("/warehouses", createWarehouseHandler).Methods("POST")
	router.HandleFunc("/locations", getLocationsHandler).Methods("GET")

	fmt.Println("Server running...")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(origins, methods, headers)(router)))
}

// Controllers
func createCustomerHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := connectToDB()
	defer db.Close()
	var customer CustomerJSON
	json.NewDecoder(r.Body).Decode(&customer)
	err := createCustomer(customer, db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func getCustomersHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := connectToDB()
	defer db.Close()
	customers, _ := fetchCustomers(db)
	json.NewEncoder(w).Encode(customers)
}

func getMaterialTypesHandler(w http.ResponseWriter, r *http.Request) {
	materialTypes := fetchMaterialTypes()
	json.NewEncoder(w).Encode(materialTypes)
}

func sendMaterialHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := connectToDB()
	defer db.Close()
	var material IncomingMaterialJSON
	json.NewDecoder(r.Body).Decode(&material)
	err := sendMaterial(material, db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(material)
}

func getIncomingMaterialsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := connectToDB()
	defer db.Close()
	materials, err := getIncomingMaterials(db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(materials)
}

func createMaterialHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := connectToDB()
	defer db.Close()
	var material MaterialJSON
	json.NewDecoder(r.Body).Decode(&material)
	err := createMaterial(material, db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(material)
}

func getMaterialsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := connectToDB()
	defer db.Close()
	materials, err := getMaterials(db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(materials)
}

func moveMaterialHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := connectToDB()
	defer db.Close()
	var material MaterialJSON
	json.NewDecoder(r.Body).Decode(&material)
	err := moveMaterial(material, db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(material)
}

func removeMaterialHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := connectToDB()
	defer db.Close()
	var material MaterialToRemoveJSON
	json.NewDecoder(r.Body).Decode(&material)
	err := removeMaterial(material, db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(material)
}

func createWarehouseHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := connectToDB()
	defer db.Close()

	var warehouse WarehouseJSON
	json.NewDecoder(r.Body).Decode(&warehouse)
	err := createWarehouse(warehouse, db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(warehouse)
}

func getLocationsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := connectToDB()
	defer db.Close()
	query1 := r.URL.Query().Get("query1")
	locations, _ := fetchLocations(db, LocationFilter{query1: query1})
	json.NewEncoder(w).Encode(locations)
}
