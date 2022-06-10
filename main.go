package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"hello/blogfunc"
	"log"
	"net/http"
)

func main() {

	router := httprouter.New()
	router.GET("/assets/*filepath", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	router.GET("/", blogfunc.Index)
	router.GET("/page/:name", blogfunc.Indexpage)
	router.GET("/archives/:name", blogfunc.Content)

	router.GET("/login", blogfunc.Loginget)
	router.POST("/login", blogfunc.Loginpost)

	router.GET("/ws", blogfunc.Ws)
	router.GET("/tz", blogfunc.Tz)
	log.Fatal(http.ListenAndServe(":80", router))

	fmt.Scanln()
}
