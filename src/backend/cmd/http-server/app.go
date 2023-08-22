package http_server

import (
	"backend/cmd/flags"
	"backend/config"
	"backend/internal/network"
	"backend/internal/network/handlers"
	"backend/internal/pkg/logger"
	"backend/internal/repository"
	"backend/internal/repository/mongo_repository"
	"backend/internal/repository/postgres_repository"
	"backend/internal/services"
	servicesImplementation "backend/internal/services/implementation"
	"github.com/jinzhu/copier"
	"net/http"
	"os"
)

type Config struct {
	Postgres flags.PostgresFlags `mapstructure:"postgres"`
	Address  string              `mapstructure:"address"`
	Port     string              `mapstructure:"port"`
}

type App struct {
	config       *config.Config
	repositories *appRepositoryFields
	services     *appServiceFields
	handlers     *handlers.Handlers
	logger       *logger.Logger
	mux          *http.ServeMux
}

type appRepositoryFields struct {
	comparisonListRepository repository.ComparisonListRepository
	discountRepository       repository.DiscountRepository
	instrumentRepository     repository.InstrumentRepository
	userRepository           repository.UserRepository
	orderRepository          repository.OrderRepository
}

type appServiceFields struct {
	CalcDiscountService   services.CalcDiscountService
	ComparisonListService services.ComparisonListService
	DiscountService       services.DiscountService
	InstrumentService     services.InstrumentService
	UserService           services.UserService
	OrderService          services.OrderService
}

const (
	ConfigFileName = "config.json"
	ConfigFilePath = "./config"
)

func (a *App) initPostgresRepositories() *appRepositoryFields {
	fields := postgres_repository.CreatePostgresRepositoryFields(ConfigFileName, ConfigFilePath)
	f := &appRepositoryFields{
		comparisonListRepository: postgres_repository.CreateComparisonListPostgresRepository(fields),
		discountRepository:       postgres_repository.CreateDiscountPostgresRepository(fields),
		instrumentRepository:     postgres_repository.CreateInstrumentPostgresRepository(fields),
		userRepository:           postgres_repository.CreateUserPostgresRepository(fields),
		orderRepository:          postgres_repository.CreateOrderPostgresRepository(fields),
	}

	return f
}

func (a *App) initMongoRepositories() *appRepositoryFields {
	fields := mongo_repository.CreateMongoRepositoryFields(ConfigFileName, ConfigFilePath)
	f := &appRepositoryFields{
		instrumentRepository:     mongo_repository.CreateInstrumentMongoRepository(fields),
		userRepository:           mongo_repository.CreateUserMongoRepository(fields),
		comparisonListRepository: mongo_repository.CreateComparisonListMongoRepository(fields),
		orderRepository:          mongo_repository.CreateOrderMongoRepository(fields),
		discountRepository:       mongo_repository.CreateDiscountMongoRepository(fields),
	}

	return f
}

var initRepositoriesMap = map[string]func(*App) *appRepositoryFields{
	"postgres": (*App).initPostgresRepositories,
	"mongo":    (*App).initMongoRepositories,
}

func (a *App) initServices(r *appRepositoryFields) *appServiceFields {
	lg := logger.New(a.config.LogPath, a.config.LogLevel)
	a.logger = lg
	calcDiscountService := servicesImplementation.NewCalcDiscountServiceImplementation(r.discountRepository, lg)
	u := &appServiceFields{
		CalcDiscountService:   calcDiscountService,
		ComparisonListService: servicesImplementation.NewComparisonListServiceImplementation(r.comparisonListRepository, r.instrumentRepository, lg),
		DiscountService:       servicesImplementation.NewDiscountServiceImplementation(r.discountRepository, r.userRepository, lg),
		InstrumentService:     servicesImplementation.NewInstrumentServiceImplementation(r.instrumentRepository, r.userRepository, lg),
		UserService:           servicesImplementation.NewUserServiceImplementation(r.userRepository, r.comparisonListRepository, calcDiscountService, lg),
		OrderService:          servicesImplementation.NewOrderServiceImplementation(r.orderRepository, r.comparisonListRepository, r.userRepository, lg),
	}

	return u
}

func (a *App) Init() {
	var c config.Config
	err := c.ParseConfig(ConfigFileName, ConfigFilePath)
	if err != nil {
		return
	}
	a.config = &c

	a.repositories = initRepositoriesMap[a.config.Db](a)
	a.services = a.initServices(a.repositories)
	handlerServices := handlers.HandlersServicesFields{}
	copier.Copy(&handlerServices, a.services)
	a.handlers = handlers.NewHandlers(handlerServices)

	a.mux = network.NewRouter(a.handlers)
}

func (a *App) Run() {
	a.Init()

	errChan := make(chan error)
	go func() {
		errChan <- http.ListenAndServe(a.config.Address+a.config.Port, a.mux)
	}()
	select {
	case err := <-errChan:
		if err != nil {
			a.logger.Fatal(err)
			os.Exit(1)
		}
	}
}
