package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB // database

var routes = flag.Bool("routes", false, "Generate router documentation") // routes

// Item Type
type Item struct {
	Name            string  `json:"name"`
	AlternativeName string  `json:"alternative_name"`
	Barcode         string  `json:"barcode"`
	Price           float64 `json:"price"`
	CategoryID      int     `json:"category_id"`
	Quantity        int     `json:"quantity"`
}

// Items List
func allItems(w http.ResponseWriter, r *http.Request) {
	items := []Item{}

	rows, err := db.Query("select name, price, quantity,barcode from items")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var r Item
		err = rows.Scan(&r.Name, &r.Price, &r.Quantity, &r.Barcode)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, r)
	}
	defer rows.Close()

	json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	// var item Item
	// barcode := chi.URLParam(r, "barcode")

	// stament, err := db.Prepare("SELECT name, price, category_id, barcode from items WHERE barcode=? ")

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// rows, er := statement.Exec(barcode)

	// if er != nil {
	// 	log.Fatal(er)
	// }

	// for rows.Next() {
	// 	var r Item
	// 	e = rows.Scan(&r.Name, &r.Price, &r.CategoryID, &r.Barcode)
	// 	if e != nil {
	// 		log.Fatal(e)
	// 	}
	// 	items = append(items, r)
	// }

	// json.NewEncoder(w).Encode(item)

}


func postItem(w http.ResponseWriter, r *http.Request) {
	var item Item

	json.NewDecoder(r.Body).Decode(&item)

	statement, err := db.Prepare("INSERT INTO items (name,price,sku,barcode,category_id) VALUES(?,?,112323213211,78548737238,0)")
	if err != nil {
		log.Fatal(err)
	}
	_, er := statement.Exec(item.Name, item.Price)

	if er != nil {
		log.Fatal(er)
	}

}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	json.NewDecoder(r.Body).Decode(&item)

	statement, err := db.Prepare("DELETE FROM items WHERE barcode=?")
	if err != nil {
		log.Fatal(err)
	}
	_, er := statement.Exec(item.Barcode)

	if er != nil {
		log.Fatal(er)
	}

}

func updateItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	barcode := chi.URLParam(r, "barcode")
	json.NewDecoder(r.Body).Decode(&item)

	statement, err := db.Prepare("UPDATE items SET name=?,price=? WHERE barcode=?")
	if err != nil {
		log.Fatal(err)
	}
	_, er := statement.Exec(item.Name, item.Price, barcode)

	if er != nil {
		log.Fatal(er)
	}

}


func init() {
	conn, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	db = conn
}

func main() {

	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:     []string{"*"},
		AllowOriginFunc:    func(r *http.Request, origin string) bool { return true },
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		MaxAge:             3599, // Maximum value not ignored by any of major browsers
	})

	// r.Group(func(r chi.Router) {
	// 	r.Use(cors.Handler)
	// })

	flag.Parse()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Timeout(60 * time.Second))

	r.Group(func(r chi.Router) {
		r.Use(cors.Handler)
		r.Route("/api", func(r chi.Router) {
			r.Get("/items", allItems)
			// r.Route("/items/{barcode}", func(r chi.Router) {
			// 	r.Get("/", getItem)
			// })
			r.Post("/items", postItem)
			r.Delete("/items", deleteItem)
			r.Route("/items/{barcode}", func(r chi.Router) {
				r.Put("/", updateItem)
			})

		})
	})

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	FileServer(r, "/", http.Dir(filesDir))

	http.ListenAndServe(":8080", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
