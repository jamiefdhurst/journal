package router

import (
	"github.com/jamiefdhurst/journal/internal/app/controller/apiv1"
	"github.com/jamiefdhurst/journal/internal/app/controller/web"
	"github.com/jamiefdhurst/journal/pkg/database"
	pkgrouter "github.com/jamiefdhurst/journal/pkg/router"
)

// NewRouter Define a new router and initialise routes
func NewRouter(db database.Database) *pkgrouter.Router {
	rtr := pkgrouter.Router{}
	rtr.Db = db
	rtr.ErrorController = &web.Error{}

	rtr.Get("/new", &web.New{})
	rtr.Post("/new", &web.New{})
	rtr.Get("/api/v1/post", &apiv1.List{})
	rtr.Post("/api/v1/post", &apiv1.Create{})
	rtr.Get("/api/v1/post/[%s]", &apiv1.Single{})
	rtr.Put("/api/v1/post/[%s]", &apiv1.Update{})
	rtr.Get("/[%s]/edit", &web.Edit{})
	rtr.Post("/[%s]/edit", &web.Edit{})
	rtr.Get("/[%s]", &web.View{})
	rtr.Get("/", &web.Index{})

	return &rtr
}
