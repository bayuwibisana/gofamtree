package routes

import (
	"gofamtree/handlers"
	"log"
	"net/http"
	"strings"
)

func RegisterRoutes() {
	// Admin routes
	http.HandleFunc("/admin/login", corsMiddleware(methodMiddleware("POST", handlers.AdminLogin)))
	http.HandleFunc("/admin/register", corsMiddleware(methodMiddleware("POST", handlers.AdminRegister)))

	// House routes
	http.HandleFunc("/houses", corsMiddleware(handleHouseRoutes))
	http.HandleFunc("/houses/", corsMiddleware(handleHouseRoutes))

	// Person routes
	http.HandleFunc("/persons", corsMiddleware(handlePersonRoutes))
	http.HandleFunc("/persons/", corsMiddleware(handlePersonRoutes))

	// Relation routes
	http.HandleFunc("/relations", corsMiddleware(handleRelationRoutes))
	http.HandleFunc("/relations/", corsMiddleware(handleRelationRoutes))

	// Family tree route
	http.HandleFunc("/family-tree/", corsMiddleware(methodMiddleware("GET", handlers.GetFamilyTree)))

	log.Println("Routes registered successfully")
}

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Method middleware
func methodMiddleware(allowedMethod string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

// House route handler
func handleHouseRoutes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if strings.HasPrefix(r.URL.Path, "/houses/") && len(strings.TrimPrefix(r.URL.Path, "/houses/")) > 0 {
			handlers.GetHouse(w, r)
		} else {
			handlers.GetHouses(w, r)
		}
	case "POST":
		if r.URL.Path == "/houses" {
			handlers.CreateHouse(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	case "PUT":
		if strings.HasPrefix(r.URL.Path, "/houses/") {
			handlers.UpdateHouse(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	case "DELETE":
		if strings.HasPrefix(r.URL.Path, "/houses/") {
			handlers.DeleteHouse(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Person route handler
func handlePersonRoutes(w http.ResponseWriter, r *http.Request) {
	
	switch r.Method {
	case "GET":
		if strings.HasPrefix(r.URL.Path, "/persons/") && len(strings.TrimPrefix(r.URL.Path, "/persons/")) > 0 {
			handlers.GetPerson(w, r)
		} else {
			handlers.GetPersons(w, r)
		}
	case "POST":
		if r.URL.Path == "/persons" {
			handlers.CreatePerson(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	case "PUT":
		if strings.HasPrefix(r.URL.Path, "/persons/") {
			handlers.UpdatePerson(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	case "DELETE":
		if strings.HasPrefix(r.URL.Path, "/persons/") {
			handlers.DeletePerson(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Relation route handler
func handleRelationRoutes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if strings.HasPrefix(r.URL.Path, "/relations/") && len(strings.TrimPrefix(r.URL.Path, "/relations/")) > 0 {
			handlers.GetRelation(w, r)
		} else {
			handlers.GetRelations(w, r)
		}
	case "POST":
		if r.URL.Path == "/relations" {
			handlers.CreateRelation(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	case "PUT":
		if strings.HasPrefix(r.URL.Path, "/relations/") {
			handlers.UpdateRelation(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	case "DELETE":
		if strings.HasPrefix(r.URL.Path, "/relations/") {
			handlers.DeleteRelation(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
