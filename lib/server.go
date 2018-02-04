package lib

import (
	"database/sql"
	"journal/controller"
	"journal/controller/apiv1"
	"journal/model"
	"log"
	"net/http"
)

// Server Contain the server
type Server struct {
	db     *sql.DB
	router Router
}

// CheckErr Check and fatal if applicable
func (s Server) CheckErr(err error) {
	if err != nil {
		log.Fatal("Error reported: ", err)
	}
}

func (s *Server) createDb() {
	err := model.JournalCreateTable()
	s.CheckErr(err)
	log.Println("Database created")
}

func (s *Server) initRouter() {
	s.router = NewRouter(s, &controller.Error{})
	s.router.Add("GET", "/", false, &controller.Index{})
	s.router.Add("GET", "/new", false, &controller.New{})
	s.router.Add("POST", "/new", false, &controller.New{})
	s.router.Add("GET", "/api/v1/post", false, &apiv1.List{})
	s.router.Add("GET", "\\/api\\/v1\\/post\\/([\\w\\-]+)", true, &apiv1.Single{})
	s.router.Add("GET", "\\/([\\w\\-]+)", true, &controller.View{})
}

// Run Run the server
func (s *Server) Run(mode string, port string) {
	if mode == "create" {
		s.createDb()
	} else {
		s.initRouter()
		s.serve(port)
	}
	s.db.Close()
}

func (s *Server) serve(port string) {
	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, &s.router))
}

// NewServer Create an instance of the server
func NewServer() Server {
	return Server{}
}
