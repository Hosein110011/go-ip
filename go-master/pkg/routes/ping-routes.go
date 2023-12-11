package routes

import (
	"github.com/gorilla/mux"
	"github.com/Hosein110011/go-master/pkg/controllers"
)


var RegisterBookStoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/api/v1", controllers.GetDataCenter).Methods("GET")
	router.HandleFunc("/api/v1/curl", controllers.GetUrls).Methods("GET")
	// router.HandleFunc("/book/{bookId}", controllers.GetBookById).Methods("GET")
	// router.HandleFunc("/book/{bookId}", controllers.UpdateBook).Methods("PUT")
	// router.HandleFunc("/book/{bookId}", controllers.DeleteBook).Methods("DELETE")
}