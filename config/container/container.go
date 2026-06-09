package container

import (
	"log"
	"net/http"

	"prj-ar-26-go-back/config"
	"prj-ar-26-go-back/internal/app"
	"prj-ar-26-go-back/internal/infra/database"
	"prj-ar-26-go-back/internal/infra/http/controllers"
	"prj-ar-26-go-back/internal/infra/http/middlewares"

	"github.com/go-chi/jwtauth/v5"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

type Container struct {
	Middlewares
	Services
	Controllers
}

type Middlewares struct {
	AuthMw func(http.Handler) http.Handler
}

type Services struct {
	app.AuthService
	app.UserService
	app.OrganizationService
	app.RoomService
	app.DeviceService
	app.MeasurementService
}

type Controllers struct {
	AuthController         controllers.AuthController
	UserController         controllers.UserController
	OrganizationController controllers.OrganizationController
	RoomController         controllers.RoomController
	DeviceController       controllers.DeviceController
	MeasurementController  controllers.MeasurementController
}

func New(conf config.Configuration) Container {
	tknAuth := jwtauth.New("HS256", []byte(conf.JwtSecret), nil)
	sess := getDbSess(conf)

	sessionRepository := database.NewSessRepository(sess)
	userRepository := database.NewUserRepository(sess)
	organizationRepository := database.NewOrganizationRepository(sess)
	roomRepository := database.NewRoomRepository(sess)
	deviceRepository := database.NewDeviceRepository(sess)
	measurementRepository := database.NewMeasurementRepository(sess)

	userService := app.NewUserService(userRepository)
	authService := app.NewAuthService(sessionRepository, userRepository, tknAuth, conf.JwtTTL)
	organizationService := app.NewOrganizationService(organizationRepository)
	roomService := app.NewRoomService(roomRepository)
	deviceService := app.NewDeviceService(deviceRepository)
	measurementService := app.NewMeasurementService(measurementRepository)

	authController := controllers.NewAuthController(authService, userService)
	userController := controllers.NewUserController(userService, authService)
	organizationController := controllers.NewOrganizationController(organizationService)
	roomController := controllers.NewRoomController(roomService)
	deviceController := controllers.NewDeviceController(deviceService)
	measurementController := controllers.NewMeasurementController(measurementService)
	authMiddleware := middlewares.AuthMiddleware(tknAuth, authService, userService)

	return Container{
		Middlewares: Middlewares{
			AuthMw: authMiddleware,
		},
		Services: Services{
			authService,
			userService,
			organizationService,
			roomService,
			deviceService,
			measurementService,
		},
		Controllers: Controllers{
			authController,
			userController,

			organizationController,
			roomController,
			deviceController,
			measurementController,
		},
	}
}

func getDbSess(conf config.Configuration) db.Session {
	sess, err := postgresql.Open(
		postgresql.ConnectionURL{
			User:     conf.DatabaseUser,
			Host:     conf.DatabaseHost, // Порт має бути записаний прямо тут, наприклад "localhost:5432"
			Password: conf.DatabasePassword,
			Database: conf.DatabaseName,
		})
	if err != nil {
		log.Fatalf("Unable to create new DB session: %q\n", err)
	}
	return sess
}
