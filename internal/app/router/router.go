package router

import (
	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/controller/apiv1"
	"github.com/jamiefdhurst/journal/internal/app/controller/web"
	pkgrouter "github.com/jamiefdhurst/journal/pkg/router"
)

// NewRouter Define a new router and initialise routes
func NewRouter(app *app.Container) *pkgrouter.Router {
	rtr := pkgrouter.Router{}
	rtr.Container = app
	rtr.ErrorController = &web.BadRequest{}
	rtr.StaticPaths = []string{
		app.Configuration.ThemePath + "/" + app.Configuration.Theme,
		app.Configuration.StaticPath,
		"api",
	}

	// API v1
	rtr.Get("/api/v1/stats", &apiv1.Stats{})
	rtr.Get("/api/v1/post", &apiv1.List{})
	rtr.Put("/api/v1/post", &apiv1.Create{})
	rtr.Get("/api/v1/post/random", &apiv1.Random{})
	rtr.Get("/api/v1/post/[%s]", &apiv1.Single{})
	rtr.Post("/api/v1/post/[%s]", &apiv1.Update{})

	// Web
	rtr.Get("/sitemap.xml", &web.Sitemap{})
	rtr.Get("/stats", &web.Stats{})
	rtr.Get("/new", &web.New{})
	rtr.Post("/new", &web.New{})
	rtr.Get("/random", &web.Random{})
	rtr.Get("/calendar/[%s]/[%s]", &web.Calendar{})
	rtr.Get("/calendar/[%s]", &web.Calendar{})
	rtr.Get("/calendar", &web.Calendar{})
	rtr.Get("/[%s]/edit", &web.Edit{})
	rtr.Post("/[%s]/edit", &web.Edit{})
	rtr.Get("/[%s]", &web.View{})
	rtr.Get("/", &web.Index{})

	return &rtr
}
