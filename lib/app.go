package lib

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jamiefdhurst/journal/controller"
	"github.com/jamiefdhurst/journal/controller/apiv1"
	"github.com/jamiefdhurst/journal/model"
)

// App Main application, contain the router
type App struct {
	router Router
}

// ExitOnError Check for an error and log/exit it if necessary
func (a App) ExitOnError(err error) {
	if err != nil {
		log.Fatal("Error reported: ", err)
	}
}

// Run Determine the mode and run appropriate app call
func (a *App) Run(mode string, port string) {
	if mode == "createdb" {
		a.createDatabase()
	} else if mode == "giphy" {
		a.giphyAPIKey()
	} else {
		a.initRouter()
		a.serveHTTP(port)
	}

	// Close database once finished
	model.Close()
}

func (a *App) createDatabase() {
	err := model.CreateGiphyTable()
	a.ExitOnError(err)
	err2 := model.CreateJournalTable()
	a.ExitOnError(err2)
	log.Println("Database created")
}

func (a App) giphyAPIKey() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter GIPHY API key: ")
	apiKey, _ := reader.ReadString('\n')
	model.UpdateGiphyAPIKey(strings.Replace(apiKey, "\n", "", -1))
	log.Println("API key saved")
}

func (a *App) initRouter() {
	var routes []Route
	a.router = Router{a, routes, &controller.Error{}}
	a.router.Get("/new", &controller.New{})
	a.router.Post("/new", &controller.New{})
	a.router.Get("/api/v1/post", &apiv1.List{})
	a.router.Post("/api/v1/post", &apiv1.Create{})
	a.router.Get("/api/v1/post/[%s]", &apiv1.Single{})
	a.router.Put("/api/v1/post/[%s]", &apiv1.Update{})
	a.router.Get("/[%s]/edit", &controller.Edit{})
	a.router.Post("/[%s]/edit", &controller.Edit{})
	a.router.Get("/[%s]", &controller.View{})
	a.router.Get("/", &controller.Index{})
}

func (a *App) serveHTTP(port string) {
	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, &a.router))
}
