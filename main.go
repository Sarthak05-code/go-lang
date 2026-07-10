package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Seat struct {
	Id string `json:"id"`
	Booked bool `json:"booked"`
	Passenger string `json:"passenger"`
}

var (
	mu sync.Mutex
	seats []Seat
)

func initSeats() {
	rows := 6
	cols := []string{"A","B","C","D"}
	seats = make([]Seat , 0 , rows*len(cols))
	for r := 1; r <= rows ; r++ {
		for _,c := range cols {
			seats = append(seats , Seat{
				Id : fmt.Sprintf("%d%s",r,c),
			})
		}
	}
}

func enableCORS (w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin" , "*")
	w.Header().Set("Acess-Control-Allow-Methods", "GET , POST , OPTIONS")
	w.Header().Set("Allow-Control-Allow-Headers" , "Content-Type")
}

func getSeatHandler(w http.ResponseWriter , r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w , "method not allowed" , http.StatusMethodNotAllowed)
		return
	}
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type" , "application/json")
	json.NewEncoder(w).Encode(seats)
}

type bookRequest struct {
	SeatId string `json:"seatId"`
	Passenger string `json:"passenger"`
}

func bookSeatHandler(w http.ResponseWriter , r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w  , "Method not allowed" , http.StatusMethodNotAllowed)
		return
	}
	var req bookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w , "Invalid request body", http.StatusBadRequest)
			return
	}
	if req.SeatId == "" || req.Passenger == "" {
		http.Error(w , "SeatId and passenger are required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i , s := range seats {
		if s.Id == req.SeatId {
			if s.Booked {
				http.Error(w , "Seat is already taken" , http.StatusConflict)
				return
			}
			seats[i].Booked = true
			seats[i].Passenger = req.Passenger

			w.Header().Set("Content-Type", "application/json")	
			json.NewEncoder(w).Encode(seats[i])
			return
		}
	}
	http.Error(w , "Seat not found" , http.StatusNotFound)
}

func resetHandler(w http.ResponseWriter , r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	mu.Lock()
	defer mu.Unlock()
	initSeats()
	w.Header().Set("Content-Type" ,"application/json" )
	json.NewEncoder(w).Encode(seats)
}


func main() {
	initSeats()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/seats", getSeatHandler)
	mux.HandleFunc("/api/book" , bookSeatHandler)
	mux.HandleFunc("/api/reset" , resetHandler)

	mux.Handle("/" , http.FileServer(http.Dir("./static")))

	fmt.Println("System Running on server : http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080" , mux))
}