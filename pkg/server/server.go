package server

import (
    "log"
    "net/http"

    "todo-list-final/pkg/api"

)

func Run(port string) error {
    api.Init()
 
	http.Handle("/", http.FileServer(http.Dir("web")))
    
    log.Printf("Server started on http://localhost:%s\n", port)
    return http.ListenAndServe(":"+port, nil)
}
