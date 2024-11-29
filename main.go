package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func createCustomer(w http.ResponseWriter, r *http.Request) {
	db, _ := connectToDB()
	var customer CustomerJSON
	json.NewDecoder(r.Body).Decode(&customer)
	err := addCustomer(customer, db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	db, _ := connectToDB()
	customers, _ := fetchCustomers(db)
	json.NewEncoder(w).Encode(customers)
}

func getMaterialTypes(w http.ResponseWriter, r *http.Request) {
	materialTypes := fetchMaterialTypes()
	json.NewEncoder(w).Encode(materialTypes)
}

func sendMaterials(w http.ResponseWriter, r *http.Request) {
	db, _ := connectToDB()
	var material MaterialTemplate
	json.NewDecoder(r.Body).Decode(&material)
	err := sendMaterial(material, db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(material)
}

func main() {
	router := mux.NewRouter()
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	router.HandleFunc("/customers", createCustomer).Methods("POST")
	router.HandleFunc("/customers", getCustomers).Methods("GET")

	router.HandleFunc("/material_types", getMaterialTypes).Methods("GET")

	router.HandleFunc("/csr_materials", sendMaterials).Methods("POST")

	fmt.Println("Server running...")
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(origins, methods, headers)(router)))
}
