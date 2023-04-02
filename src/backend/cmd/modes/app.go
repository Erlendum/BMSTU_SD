package modes

import (
	"backend/cmd/modes/flags"
	"backend/internal/repository"
	mock_repository "backend/internal/repository/mocks"
	"backend/internal/services"
	servicesImplementation "backend/internal/services/implementation"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang/mock/gomock"
	"log"
	"time"
)

type Config struct {
	Postgres flags.PostgresFlags `mapstructure:"postgres"`
	Address  string              `mapstructure:"address"`
	Port     string              `mapstructure:"port"`
}

type App struct {
	Config       Config
	repositories *appRepositoryFields
	services     *appServiceFields
	ctrl         *gomock.Controller
	app          *fiber.App
}

type appRepositoryFields struct {
	comparisonListRepository repository.ComparisonListRepository
	discountRepository       repository.DiscountRepository
	instrumentRepository     repository.InstrumentRepository
	userRepository           repository.UserRepository
}

type appServiceFields struct {
	calcDiscountService   services.CalcDiscountService
	comparisonListService services.ComparisonListService
	discountService       services.DiscountService
	instrumentService     services.InstrumentService
	userService           services.UserService
}

func (a *App) initRepositories(ctrl *gomock.Controller) *appRepositoryFields {
	f := &appRepositoryFields{
		comparisonListRepository: mock_repository.NewMockComparisonListRepository(ctrl),
		discountRepository:       mock_repository.NewMockDiscountRepository(ctrl),
		instrumentRepository:     mock_repository.NewMockInstrumentRepository(ctrl),
		userRepository:           mock_repository.NewMockUserRepository(ctrl),
	}

	return f
}

func (a *App) initServices(r *appRepositoryFields) *appServiceFields {
	calcDiscountService := servicesImplementation.NewCalcDiscountServiceImplementation(r.discountRepository)
	u := &appServiceFields{
		calcDiscountService:   calcDiscountService,
		comparisonListService: servicesImplementation.NewComparisonListServiceImplementation(r.comparisonListRepository, r.instrumentRepository),
		discountService:       servicesImplementation.NewDiscountServiceImplementation(r.discountRepository, r.userRepository),
		instrumentService:     servicesImplementation.NewInstrumentServiceImplementation(r.instrumentRepository, r.userRepository),
		userService:           servicesImplementation.NewUserServiceImplementation(r.userRepository, r.comparisonListRepository, calcDiscountService),
	}

	return u
}

//func (a *App) ParseConfig(pathToConfig string, configFileName string) error {
//	v := viper.New()
//	v.SetConfigName(configFileName)
//	v.SetConfigType("json")
//	v.AddConfigPath(pathToConfig)
//
//	err := v.ReadInConfig()
//	if err != nil {
//		return err
//	}
//
//	err = v.Unmarshal(&a.Config)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func (a *App) Init() {
	a.ctrl = gomock.NewController(nil)
	a.repositories = a.initRepositories(a.ctrl)
	a.services = a.initServices(a.repositories)

	a.app = fiber.New()
	a.app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUID,
	}))
	a.app.Use(recover.New())
}

func (a *App) Run() {
	a.Init()

	log.Fatal(a.app.Listen(a.Config.Address + a.Config.Port))
}
