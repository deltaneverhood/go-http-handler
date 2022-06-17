package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type dish struct {
	id    uint    `json: "id"`
	name  string  `json: "name"`
	price float32 `json: "price"`
}

func getDishes() []dish {
	return []dish{
		{
			id:    1,
			name:  "sushi",
			price: 500,
		},
		{
			id:    2,
			name:  "pizza",
			price: 600,
		},
		{
			id:    3,
			name:  "burger",
			price: 450,
		},
		{
			id:    4,
			name:  "doner",
			price: 200,
		},
		{
			id:    5,
			name:  "coffee",
			price: 120,
		},
	}
}

// Обработчик главной странице.
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello!"))
}

func showAllDishes(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method "+r.Method+" is not allowed in this request", 405)
		return
	}

	// по какой-то причине структура не маршаллится в json, поэтому временно преобразуем к string
	dishes := getDishes()
	var s []string
	for _, v := range dishes {
		s = append(s, fmt.Sprintf("%d", v.id), v.name, fmt.Sprintf("%f", v.price))
	}
	jsonOut, err := json.Marshal(s)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonOut)

}

func showDish(w http.ResponseWriter, r *http.Request) {

	dishes := getDishes()

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method "+r.Method+" is not allowed in this request", 405)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 || id > len(dishes) {
		http.NotFound(w, r)
		return
	}
	id -= 1

	fmt.Fprintf(w, fmt.Sprintf("%d", dishes[id].id), dishes[id].name, fmt.Sprintf("%f", dishes[id].price))
}

func addDish(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method "+r.Method+" is not allowed in this request", 405)
		return
	}

	// здесь распарсим json в теле запроса и добавим в список блюд
	w.Write([]byte("Adding new dish..."))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/dishes/all", showAllDishes)
	mux.HandleFunc("/dishes", showDish)
	mux.HandleFunc("/snippet/create", addDish)

	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
