package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"prj-ar-26-go-back/config"
	"prj-ar-26-go-back/config/container"
	"prj-ar-26-go-back/internal/app"
	"prj-ar-26-go-back/internal/infra/http/controllers"
	"prj-ar-26-go-back/internal/infra/http/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5/middleware"
)

func Router(cont container.Container) http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.RedirectSlashes, middleware.Logger, cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "capacitor://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Route("/api", func(apiRouter chi.Router) {
		// Health
		apiRouter.Route("/ping", func(healthRouter chi.Router) {
			healthRouter.Get("/", PingHandler())
			healthRouter.Handle("/*", NotFoundJSON())
		})

		apiRouter.Route("/v1", func(apiRouter chi.Router) {
			// Public routes
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Route("/auth", func(apiRouter chi.Router) {
					AuthRouter(apiRouter, cont.AuthController, cont.AuthMw)
				})
			})

			// Protected routes
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Use(cont.AuthMw)

				UserRouter(apiRouter, cont.UserController)
				OrganizationRouter(
					apiRouter,
					cont.OrganizationController,
					cont.OrganizationService)
				RoomRouter(
					apiRouter,
					cont.RoomController,
					cont.RoomService,
				)
				DeviceRouter(
					apiRouter,
					cont.DeviceController,
					cont.DeviceService,
				)
				MeasurementRouter(
					apiRouter,
					cont.MeasurementController,
					cont.MeasurementService,
				)
				EventRouter(
					apiRouter,
					cont.EventController,
					cont.EventService,
				)
				apiRouter.Handle("/*", NotFoundJSON())
			})
		})
	})

	router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, config.GetConfiguration().FileStorageLocation))
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
		fs.ServeHTTP(w, r)
	})

	return router
}

func AuthRouter(r chi.Router, ac controllers.AuthController, amw func(http.Handler) http.Handler) {
	r.Route("/", func(apiRouter chi.Router) {
		apiRouter.Post("/register", ac.Register())
		apiRouter.Post(
			"/login",
			ac.Login(),
		)
		apiRouter.With(amw).Post(
			"/logout",
			ac.Logout(),
		)
	})
}

func UserRouter(r chi.Router, uc controllers.UserController) {
	r.Route("/users", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/",
			uc.FindMe(),
		)
		apiRouter.Put(
			"/",
			uc.Update(),
		)
		apiRouter.Delete(
			"/",
			uc.Delete(),
		)
	})
}

func OrganizationRouter(r chi.Router, oc controllers.OrganizationController, os app.OrganizationService) {
	opom := middlewares.PathObject("orgId", controllers.OrgKey, os)
	r.Route("/organizations", func(apiRouter chi.Router) {
		apiRouter.Post("/", oc.Save())
		apiRouter.Get("/", oc.FindList())
		apiRouter.With(opom).Get("/{orgId}", oc.Find())
		apiRouter.With(opom).Put("/{orgId}", oc.Update())
		apiRouter.With(opom).Delete("/{orgId}", oc.Delete())
	})
}

func RoomRouter(r chi.Router, rc controllers.RoomController, rs app.RoomService) {

	rpom := middlewares.PathObject("roomId", controllers.RoomKey, rs)

	r.Route("/rooms", func(apiRouter chi.Router) {

		apiRouter.Post("/", rc.Save())
		apiRouter.Get("/", rc.FindList())
		apiRouter.With(rpom).Get("/{roomId}", rc.Find())
		apiRouter.With(rpom).Put("/{roomId}", rc.Update())
		apiRouter.With(rpom).Delete("/{roomId}", rc.Delete())
	})
}
func DeviceRouter(r chi.Router, dc controllers.DeviceController, ds app.DeviceService) {

	dpom := middlewares.PathObject("deviceId", controllers.DevKey, ds)

	r.Route("/devices", func(apiRouter chi.Router) {

		apiRouter.Post("/", dc.Save())
		apiRouter.Get("/", dc.FindList())
		apiRouter.With(dpom).Get("/{deviceId}", dc.Find())
		apiRouter.With(dpom).Put("/{deviceId}", dc.Update())
		apiRouter.With(dpom).Delete("/{deviceId}", dc.Delete())
	})
}
func MeasurementRouter(r chi.Router, mc controllers.MeasurementController, ms app.MeasurementService) {
	mpom := middlewares.PathObject("measId", controllers.MesKey, ms)

	r.Route("/measurements", func(apiRouter chi.Router) {
		// 1. Ендпоінт для пристроїв (запис вимірювань сенсорів у БД)
		apiRouter.Post("/", mc.Save())

		// 2. Ендпоінт для адміністратора (перегляд за пристроєм за день/тиждень/місяць)
		// Доступний за URL: GET /api/v1/measurements/report?device_id=X&period=week
		apiRouter.Get("/report", mc.GetAdminReport())

		// Standard CRUD
		apiRouter.Get("/", mc.FindList())
		apiRouter.With(mpom).Get("/{measId}", mc.Find())
		apiRouter.With(mpom).Put("/{measId}", mc.Update())
		apiRouter.With(mpom).Delete("/{measId}", mc.Delete())
	})
}
func EventRouter(r chi.Router, ec controllers.EventController, es app.EventService) {
	r.Route("/events", func(apiRouter chi.Router) {
		apiRouter.Post("/", ec.Save())                        // Прийом подій від пристроїв
		apiRouter.Get("/energy-report", ec.GetEnergyReport()) // Звіт для адміна
	})
}

func NotFoundJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode("Resource Not Found")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}

func PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode("Ok")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}
