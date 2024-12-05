package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	// Routes
	router.HandleFunc("/customers", createCustomerHandler).Methods("POST")
	router.HandleFunc("/customers", getCustomersHandler).Methods("GET")

	router.HandleFunc("/materials", createMaterialHandler).Methods("POST")
	router.HandleFunc("/material_types", getMaterialTypesHandler).Methods("GET")

	router.HandleFunc("/incoming_materials", sendMaterialHandler).Methods("POST")
	router.HandleFunc("/incoming_materials", getIncomingMaterialsHandler).Methods("GET")

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

func getLocationsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := connectToDB()
	defer db.Close()
	query1 := r.URL.Query().Get("query1")
	locations, _ := fetchLocations(db, LocationFilter{query1: query1})
	json.NewEncoder(w).Encode(locations)
}
