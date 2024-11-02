package routes

import (
	"net/http"

	"github.com/Akhilbisht798/cloud-text-editor/go-server/internal/handlers"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/register", applyMiddleware(handlers.Register, enableCORS))
	mux.HandleFunc("/api/login", applyMiddleware(handlers.Login, enableCORS))
	mux.HandleFunc("/api/user", applyMiddleware(handlers.GetUser, enableCORS))
	mux.HandleFunc("/api/logout", applyMiddleware(handlers.Logout, enableCORS))

	//If the project is alredy present start that if not create a new one.
	mux.HandleFunc("/api/start", applyMiddleware(handlers.CreateProject, enableCORS))
	mux.HandleFunc("/api/getUserFiles", applyMiddleware(handlers.S3PresignedGetURLHandler, enableCORS))
	mux.HandleFunc("/api/close", applyMiddleware(handlers.CloseTask, enableCORS))
	return mux
}
