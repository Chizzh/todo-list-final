package server

import (
    "log"
    "net/http"

    "todo-list-final/pkg/api"

)

func Run(port string) error {
    api.Init()

    fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)
    
    log.Printf("Server starting on :%s", port)
    return http.ListenAndServe(":"+port, nil)
}
