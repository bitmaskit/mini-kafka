package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func DelayedResponse(duration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling request...")
		defer log.Println("request handled...")

		// Simulate a slow response or long-running task
		time.Sleep(duration)

		_, err := fmt.Fprintln(w, "hello")
		if err != nil {
			log.Println("failed to write response:", err)
			return
		}
	}
}

//
//func Slow(w http.ResponseWriter, r *http.Request) {
//	log.Println("handling request...")
//	defer log.Println("request handled...")
//
//	// Simulate a slow response or long-running task
//	time.Sleep(7 * time.Second)
//
//	_, err := fmt.Fprintln(w, "hello")
//	if err != nil {
//		log.Println("failed to write response:", err)
//		return
//	}
//}
//
//func ReallySlow(w http.ResponseWriter, r *http.Request) {
//	log.Println("handling request...")
//	defer log.Println("request handled...")
//
//	// Simulate a slow response or long-running task
//	time.Sleep(20 * time.Second)
//
//	_, err := fmt.Fprintln(w, "hello")
//	if err != nil {
//		log.Println("failed to write response:", err)
//		return
//	}
//}
