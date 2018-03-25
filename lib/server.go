package lib

import (
	"bufio"
	"fmt"
	"journal/controller"
	"journal/controller/apiv1"
	"journal/model"
	"log"
	"net/http"
	"os"
	"strings"
)

// Server Contain the server
type Server struct {
	router Router
}

// CheckErr Check and fatal if applicable
func (s Server) CheckErr(err error) {
	if err != nil {
		log.Fatal("Error reported: ", err)
	}
}

func (s *Server) createDb() {
	err := model.GiphyCreateTable()
	s.CheckErr(err)
	err2 := model.JournalCreateTable()
	s.CheckErr(err2)
	log.Println("Database created")
}

func (s *Server) giphyAPIKey() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter GIPHY API key: ")
	apiKey, _ := reader.ReadString('\n')
	gs := model.Giphys{}
	gs.Update(strings.Replace(apiKey, "\n", "", -1))
	log.Println("API key saved")
}

func (s *Server) initRouter() {
	s.router = NewRouter(s, &controller.Error{})
	s.router.Add("GET", "/", false, &controller.Index{})
	s.router.Add("GET", "/new", false, &controller.New{})
	s.router.Add("POST", "/new", false, &controller.New{})
	s.router.Add("GET", "/api/v1/post", false, &apiv1.List{})
	s.router.Add("POST", "/api/v1/post", false, &apiv1.Create{})
	s.router.Add("GET", "\\/api\\/v1\\/post\\/([\\w\\-]+)", true, &apiv1.Single{})
	s.router.Add("PUT", "\\/api\\/v1\\/post\\/([\\w\\-]+)", true, &apiv1.Update{})
	s.router.Add("GET", "\\/([\\w\\-]+)\\/edit", true, &controller.Edit{})
	s.router.Add("POST", "\\/([\\w\\-]+)\\/edit", true, &controller.Edit{})
	s.router.Add("GET", "\\/([\\w\\-]+)", true, &controller.View{})
}

// Run Run the server
func (s *Server) Run(mode string, port string) {
	if mode == "create" {
		s.createDb()
	} else if mode == "giphy" {
		s.giphyAPIKey()
	} else {
		s.initRouter()
		s.serve(port)
	}
	model.Close()
}

func (s *Server) serve(port string) {
	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, &s.router))
}

// NewServer Create an instance of the server
func NewServer() Server {
	return Server{}
}
