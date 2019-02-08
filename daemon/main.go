package main

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/songxinjianqwe/scheduler/daemon/handler"
	"net/http"
	"os"
	"strings"
)

func init() {
	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	log.SetFormatter(&log.TextFormatter{})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	log.SetOutput(os.Stdout)
	//设置最低loglevel
	log.SetLevel(log.InfoLevel)
}

func main() {
	// register router
	router := mux.NewRouter().StrictSlash(true)
	subRouter := router.PathPrefix("/api/tasks").Subrouter()

	subRouter.HandleFunc("", handler.GetAllTasksHandler).Methods("GET")
	subRouter.HandleFunc("/{id}", handler.GetTaskInfoHandler).Methods("GET")
	subRouter.HandleFunc("", handler.SubmitTask).Methods("POST")
	subRouter.HandleFunc("/{id}", handler.StopTask).Methods("PUT")
	subRouter.HandleFunc("/{id}", handler.DeleteTask).Methods("DELETE")

	// 打印一下handlers
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})
	// start handler listening
	err := http.ListenAndServe(":8865", router)
	if err != nil {
		log.Fatalln("ListenAndServe err:", err)
	}
	log.Println("Server end")
}
