package main

import (
	"extender/controller"
	//"github.com/cnych/sample-scheduler-extender/controller"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	router := httprouter.New()
	router.GET("/", controller.Index)
	router.POST("/filter", controller.Filter)
	router.POST("/prioritize", controller.Prioritize)

	log.Printf("extender 已经开启!!!!\n")
	log.Fatal(http.ListenAndServe(":8888", router))
}
