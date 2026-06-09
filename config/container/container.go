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
	app.EventService // Додано сервіс для подій
}

type Controllers struct {
	AuthController         controllers.AuthController
	UserController         controllers.UserController
	OrganizationController controllers.OrganizationController
	RoomController         controllers.RoomController
	DeviceController       controllers.DeviceController
	MeasurementController  controllers.MeasurementController
	EventController        controllers.EventController // Додано контролер для подій
}

func New(conf config.Configuration) Container {
	tknAuth := jwtauth.New("HS256", []byte(conf.JwtSecret), nil)
	sess := getDbSess(conf)

	// Ініціалізація репозиторіїв
	sessionRepository := database.NewSessRepository(sess)
	userRepository := database.NewUserRepository(sess)
	organizationRepository := database.NewOrganizationRepository(sess)
	roomRepository := database.NewRoomRepository(sess)
	deviceRepository := database.NewDeviceRepository(sess)
	measurementRepository := database.NewMeasurementRepository(sess)
	eventRepository := database.NewEventRepository(sess) // Додано репозиторій подій

	// Ініціалізація бізнес-сервісів
	userService := app.NewUserService(userRepository)
	authService := app.NewAuthService(sessionRepository, userRepository, tknAuth, conf.JwtTTL)
	organizationService := app.NewOrganizationService(organizationRepository)
	roomService := app.NewRoomService(roomRepository)
	deviceService := app.NewDeviceService(deviceRepository)
	measurementService := app.NewMeasurementService(measurementRepository)
	eventService := app.NewEventService(eventRepository) // Додано сервіс подій

	// Ініціалізація HTTP-контролерів
	authController := controllers.NewAuthController(authService, userService)
	userController := controllers.NewUserController(userService, authService)
	organizationController := controllers.NewOrganizationController(organizationService)
	roomController := controllers.NewRoomController(roomService)
	deviceController := controllers.NewDeviceController(deviceService)
	measurementController := controllers.NewMeasurementController(measurementService)
	eventController := controllers.NewEventController(eventService) // Додано контролер подій

	authMiddleware := middlewares.AuthMiddleware(tknAuth, authService, userService)

	return Container{
		Middlewares: Middlewares{
			AuthMw: authMiddleware,
		},
		Services: Services{
			AuthService:         authService,
			UserService:         userService,
			OrganizationService: organizationService,
			RoomService:         roomService,
			DeviceService:       deviceService,
			MeasurementService:  measurementService,
			EventService:        eventService, // Передаємо у сервіси
		},
		Controllers: Controllers{
			AuthController:         authController,
			UserController:         userController,
			OrganizationController: organizationController,
			RoomController:         roomController,
			DeviceController:       deviceController,
			MeasurementController:  measurementController,
			EventController:        eventController, // Передаємо у контролери
		},
	}
}

func getDbSess(conf config.Configuration) db.Session {
	sess, err := postgresql.Open(
		postgresql.ConnectionURL{
			User:     conf.DatabaseUser,
			Host:     conf.DatabaseHost,
			Password: conf.DatabasePassword,
			Database: conf.DatabaseName,
		})
	if err != nil {
		log.Fatalf("Unable to create new DB session: %q\n", err)
	}
	return sess
}
