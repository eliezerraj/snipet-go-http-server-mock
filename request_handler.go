package main

import (
	"fmt"
	"net/http"
	"math/rand"
	"time"

)

func MiddleWareHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("MiddleWareHandler (INICIO)")

		var seed = my_setup.ResponseTime
		if (my_setup.IsRandomTime && seed > 0){
			seed = rand.Intn(my_setup.ResponseTime)
		}
		time.Sleep(time.Second * time.Duration(seed))

		switch {
			case my_setup.ResponseStatusCode >= 500:
				w.WriteHeader(my_setup.ResponseStatusCode)
			case my_setup.ResponseStatusCode >= 400:	
				w.WriteHeader(my_setup.ResponseStatusCode)
			default:	
				w.WriteHeader(http.StatusOK)
		}

		fmt.Println("MiddleWareHandler (FIM)")
		next.ServeHTTP(w, r)
	})
}