package modes

import (
	"backend/cmd/modes/flags"
	"backend/config"
	"backend/internal/network"
	"backend/internal/network/handlers"
	"backend/internal/pkg/logger"
	"backend/internal/repository"
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
}

type appServiceFields struct {
	CalcDiscountService   services.CalcDiscountService
	ComparisonListService services.ComparisonListService
	DiscountService       services.DiscountService
	InstrumentService     services.InstrumentService
	UserService           services.UserService
}

func (a *App) initRepositories() *appRepositoryFields {
	fields := postgres_repository.CreatePostgresRepositoryFields("config.json", "./config")
	a.config = &fields.Config
	f := &appRepositoryFields{
		comparisonListRepository: postgres_repository.CreateComparisonListPostgresRepository(fields),
		discountRepository:       postgres_repository.CreateDiscountPostgresRepository(fields),
		instrumentRepository:     postgres_repository.CreateInstrumentPostgresRepository(fields),
		userRepository:           postgres_repository.CreateUserPostgresRepository(fields),
	}

	return f
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
	}

	return u
}

func (a *App) Init() {
	a.repositories = a.initRepositories()
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
