package main

import (
		"fmt"
		"time"
		"os"
		"os/signal"
		"net"
		"net/http"
		"strconv"
		"encoding/json"
		"io/ioutil"
		"context"
		"syscall"
		"math/rand"

		"github.com/joho/godotenv"
		"github.com/gorilla/mux"
		"github.com/jaswdr/faker"

		entity "github/snipet-go-http-server-mock/main/entity"

)

var my_pod entity.Pod
var my_setup entity.Setup 

func envVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("*** WARNING .env file NOT FOUND, using the os.env")
	}
	return os.Getenv(key)
}

func main(){
	fmt.Println("Starting Server Mock 1.1")

	// catch the ip-address
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Erro Fatal")
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				my_pod.Ip = ipnet.IP.String()
			}
		}
	}
	// catch the port in the ENV
	my_pod.Port = envVariable("PORT")
	if my_pod.Port == ""{
		my_pod.Port = "8900"
	}
	my_pod.Name = envVariable("NAME_POD")
	if my_pod.Name == ""{
		my_pod.Name = "no-name"
	}

	my_pod.PID = strconv.Itoa(os.Getpid())

	my_setup.ResponseTime = 0
	my_setup.ResponseStatusCode = 200
	my_setup.IsRandomTime = false
	my_setup.Count = 100
	
	my_pod.Setup = my_setup

	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", health)
    myRouter.HandleFunc("/setup", setup).Methods("POST")

	stress_cpu_route := myRouter.Methods(http.MethodPost).Subrouter()
    stress_cpu_route.HandleFunc("/stress_cpu", stress_cpu).Methods("POST")
	stress_cpu_route.Use(MiddleWareHandler)

	customer_fake_route := myRouter.Methods(http.MethodGet).Subrouter()
    customer_fake_route.HandleFunc("/customer_fake", customer_fake).Methods("GET")
	customer_fake_route.Use(MiddleWareHandler)
	
	list_customer_fake_route := myRouter.Methods(http.MethodGet).Subrouter()
    list_customer_fake_route.HandleFunc("/list_customer_fake", list_customer_fake).Methods("GET")
	list_customer_fake_route.Use(MiddleWareHandler)

	list_account_fake_route := myRouter.Methods(http.MethodGet).Subrouter()
    list_account_fake_route.HandleFunc("/list_account_fake", list_account_fake).Methods("GET")
	list_account_fake_route.Use(MiddleWareHandler)

	s := http.Server{
		Addr:         ":" + my_pod.Port,      	
		Handler:      myRouter,                	          
		ReadTimeout:  time.Duration(600) * time.Second,   
		WriteTimeout: time.Duration(600) * time.Second,  
		IdleTimeout:  time.Duration(600) * time.Second, 
	}

	go func() {
		fmt.Printf("Server Running -> name: %v | pid: %v | ip: %v | port: %v) \n", my_pod.Name , my_pod.PID ,my_pod.Ip, my_pod.Port)
		err := s.ListenAndServe()
		if err != nil {
			fmt.Printf("Fatal error starting server: %s\n", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//signal.Notify(c, os.Interrupt)
	fmt.Println("app: received a shutdown signal:", <-c)
	
	ctx , cancel := context.WithTimeout(context.Background(), time.Duration(60) * time.Second)
	defer cancel()
	fmt.Println("Shutting down...")
	if err := s.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		fmt.Println("app: dirty shutdown:", err)
		return
	}
	fmt.Println("Shutdown DONE !!!")
}

func setup(w http.ResponseWriter, r *http.Request) {
	/*
	http://localhost:8900/setup
	{
		"response_time":4,
		"response_status_code":200,
		"is_random_time": true
	}
	*/
	fmt.Println("/setup")

	reqBody, _ := ioutil.ReadAll(r.Body)
    var setup entity.Setup 
    json.Unmarshal(reqBody, &setup)

	my_setup.ResponseTime = setup.ResponseTime
	my_setup.IsRandomTime = setup.IsRandomTime
	my_setup.ResponseStatusCode = setup.ResponseStatusCode
	my_setup.Count = setup.Count

	my_pod.Setup = my_setup

    json.NewEncoder(w).Encode(my_pod.Setup)
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/")

	json.NewEncoder(w).Encode(my_pod)
}

func customer_fake(w http.ResponseWriter, r *http.Request) {
	/* 
	http://localhost:8900/customer_fake
	*/
	fmt.Println("/customer_fake")

	f := faker.New()
	fake_data := entity.FakeData{}
	fake_data.Id                = f.IntBetween(0, 10000000)
	fake_data.Name              = f.Person().Name()
	fake_data.Address           = f.Address().Address()
	fake_data.Email             = f.Internet().Email()
	fake_data.CreditCardNumber  = f.Payment().CreditCardNumber()
	var phone_arr []entity.PhoneArray
	for i:=0; i < 2; i++ {
		f_phone := entity.PhoneArray{Number: f.Phone().Number()}
		phone_arr = append(phone_arr, f_phone)
	}
	fake_data.Phones 			= phone_arr
	fake_data.Obs               = f.Lorem().Sentence(10)

	json.NewEncoder(w).Encode(fake_data)
}

func list_customer_fake(w http.ResponseWriter, r *http.Request) {
	/* 
	http://localhost:8900/list_customer_fake
	*/
	fmt.Println("/list_customer_fake")

	fmt.Println("Inicio...")
	f := faker.New()
	
	var list_fake_data []entity.FakeData

	for i:=0; i < 50; i++ {
		fake_data := entity.FakeData{}
		fake_data.Id                = f.IntBetween(0, 10000000)
		fake_data.Name              = f.Person().Name()
		fake_data.Address           = f.Address().Address()
		fake_data.Email             = f.Internet().Email()
		fake_data.CreditCardNumber  = f.Payment().CreditCardNumber()
		var phone_arr []entity.PhoneArray
		for i:=0; i < 2; i++ {
			f_phone := entity.PhoneArray{Number: f.Phone().Number()}
			phone_arr = append(phone_arr, f_phone)
		}
		fake_data.Phones 			= phone_arr
		fake_data.Obs               = f.Lorem().Sentence(20)

		list_fake_data = append(list_fake_data, fake_data)
	}

	fmt.Println("Fim...")
	json.NewEncoder(w).Encode(list_fake_data)
}

func stress_cpu(w http.ResponseWriter, r *http.Request) {
	/*
		{
			"count":200
		}
	*/
	fmt.Println("/stress_cpu")

	start := time.Now()

	reqBody, _ := ioutil.ReadAll(r.Body)
    var setup entity.Setup 
    json.Unmarshal(reqBody, &setup)

	count := setup.Count
	fmt.Println("count", count)

	for n := 0; n <= count; n++ {
		f := make([]int, count+1, count+2)
		if count < 2 {
			f = f[0:2]
		}
		f[0] = 0
		f[1] = 1
		for i := 2; i <= count; i++ {
			f[i] = f[i-1] + f[i-2]
		}
    }
	
	t := time.Now()
	elapsed := t.Sub(start)

	fmt.Println("elapsed : ", elapsed)

	w.Write([]byte("Done in " + elapsed.String()))
}

func list_account_fake(w http.ResponseWriter, r *http.Request) {
	/* 
	http://localhost:8900/list_account_fake
	*/
	fmt.Println("/list_account_fake")

	fmt.Println("Inicio...")

	f := faker.New()
	rand.Seed(time.Now().UnixNano())
	randon_i := rand.Intn(100000)

	var array_balance []entity.Balance

	for i:=0; i < 50; i++ {
		_balance :=  entity.Balance{}

		_balance.Account 	= f.Phone().Number()
		_balance.Amount 	= int32(randon_i)
		_balance.DateBalance = time.Now()
		_balance.Description =  f.Lorem().Sentence(8)

		array_balance = append(array_balance, _balance)
	}

	fmt.Println("Fim...")
	json.NewEncoder(w).Encode(array_balance)
}