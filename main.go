package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

// Controller defn
type Controller struct {
	params []string
}

// ControllerInterface provide the interface for the controller
type ControllerInterface interface {
	Run(w http.ResponseWriter, r *http.Request)
	SetParams(p []string)
}

// IndexController Handle displaying all blog entries
type IndexController struct {
	Controller
}

// ViewController Handle displaying individual entry
type ViewController struct {
	Controller
}

// ViewData Data for view
type ViewData struct {
	Params []string
}

// SetParams on the controller
func (c *Controller) SetParams(p []string) {
	c.params = p
}

// Run IndexController
func (c *IndexController) Run(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./src/journal/views/_layout/header.tmpl", "./src/journal/views/_layout/footer.tmpl", "./src/journal/views/index.tmpl")
	t.Execute(w, nil)
}

// Run ViewController
func (c *ViewController) Run(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./src/journal/views/view.gtpl")
	v := ViewData{Params: c.params}
	t.Execute(w, v)
}

type route struct {
	method     string
	uri        string
	matchable  bool
	controller ControllerInterface
}

type mux struct {
	routes []route
}

func (m *mux) add(t string, u string, a bool, c ControllerInterface) {
	r := route{t, u, a, c}
	m.routes = append(m.routes, r)
}

func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Printf("%s: %s", r.Method, r.URL.Path)
	for _, route := range m.routes {
		fmt.Println(route.uri, r.URL.Path)
		if r.URL.Path == route.uri && (r.Method == route.method || (r.Method == "" && route.method == "GET")) {
			route.controller.Run(w, r)
			return
		}

		// Attempt regex match
		if route.matchable {
			matched, _ := regexp.MatchString(route.uri, r.URL.Path)
			if matched && (r.Method == route.method || (r.Method == "" && route.method == "GET")) {
				re := regexp.MustCompile(route.uri)
				route.controller.SetParams(re.FindAllString(r.URL.Path, -1))
				route.controller.Run(w, r)
				return
			}
		}
	}

	log.Printf("%s: %s 404 Not Found", r.Method, r.URL.Path)
	http.NotFound(w, r)
	return
}

func checkErr(err error) {
	if err != nil {
		log.Fatal("Error reported: ", err)
	}
}

func main() {
	const version = "0.1"

	// Command line flags
	var (
		mode = flag.String("mode", "run", "Run or create database file")
		port = flag.String("port", "3000", "Port to run web server on")
	)
	flag.Parse()

	// Load database
	db, err := sql.Open("sqlite3", "./data/journal.db")
	checkErr(err)
	fmt.Printf("Journal v%s...\n-------------------\n\n", version)

	if *mode == "create" {

		_, err := db.Exec("CREATE TABLE `journal` (" +
			"`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
			"`slug` VARCHAR(255) NOT NULL, " +
			"`title` VARCHAR(255) NOT NULL, " +
			"`date` DATE NOT NULL, " +
			"`content` TEXT NOT NULL" +
			")")
		checkErr(err)
		db.Close()
		log.Println("Database created")

	} else {

		m := &mux{}
		m.add("GET", "/", false, &IndexController{})
		m.add("GET", "\\/([\\w\\-]+)", true, &ViewController{})

		log.Printf("Listening on port %s\n", *port)
		log.Fatal(http.ListenAndServe(":"+*port, m))

	}
}
