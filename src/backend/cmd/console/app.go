package console

import (
	"backend/config"
	"backend/internal/cli/CLIhandlers"
	"backend/internal/models"
	"backend/internal/pkg/logger"
	"backend/internal/repository"
	"backend/internal/repository/mongo_repository"
	"backend/internal/repository/postgres_repository"
	"backend/internal/services"
	servicesImplementation "backend/internal/services/implementation"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/lukewarlow/GoConsoleMenu"
)

type App struct {
	repositories *appRepositoryFields
	services     *appServiceFields
	handlers     *CLIhandlers.Handlers
	config       *config.Config
	logger       *logger.Logger
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
	handlerServices := CLIhandlers.HandlersServicesFields{}
	copier.Copy(&handlerServices, a.services)
	a.handlers = CLIhandlers.NewHandlers(handlerServices)
}

func (a *App) Run() {
	a.Init()

	defer a.logger.Close()

	var user *models.User

	guestFuncs := true
	userFuncs := false
	adminFuncs := false
	Menu := GoConsoleMenu.NewUpdatableMenu("Music Store", func(menu *GoConsoleMenu.Menu) {
		if guestFuncs {
			menu.ShowMenuItem(1)
			menu.ShowMenuItem(2)
		} else {
			menu.HideMenuItem(1)
			menu.HideMenuItem(2)
		}
		if userFuncs {
			menu.ShowMenuItem(5)
			menu.ShowMenuItem(6)
			menu.ShowMenuItem(7)
			menu.ShowMenuItem(15)
			menu.ShowMenuItem(16)
		}

		if adminFuncs {
			menu.ShowMenuItem(8)
			menu.ShowMenuItem(9)
			menu.ShowMenuItem(10)
			menu.ShowMenuItem(11)
			menu.ShowMenuItem(12)
			menu.ShowMenuItem(13)
			menu.ShowMenuItem(14)
			menu.ShowMenuItem(17)
		}
	})
	Menu.AddMenuItem(GoConsoleMenu.NewActionItem(0, "Exit menu", func() {}).SetAsExitOption())

	instrumentsPerPage := 10
	instrumentPageNumber := 0
	instrumentSubMenu := GoConsoleMenu.NewMenu("Navigating through pages")
	instrumentSubMenu.AddMenuItem(GoConsoleMenu.NewActionItem(0, "Exit instruments list", func() { instrumentPageNumber = 0 }).SetAsExitOption())
	instrumentSubMenu.AddMenuItem(GoConsoleMenu.NewActionItem(1, "Show list", func() {
		instruments := a.handlers.InstrumentHandler.GetList(instrumentPageNumber*instrumentsPerPage, instrumentsPerPage)
		fmt.Println(instruments)
	}))
	instrumentSubMenu.AddMenuItem(GoConsoleMenu.NewActionItem(2, "Previous page", func() {
		instrumentPageNumber--
		if instrumentPageNumber < 0 {
			instrumentPageNumber = 0
		}
	}))
	instrumentSubMenu.AddMenuItem(GoConsoleMenu.NewActionItem(3, "Next page", func() {
		instrumentPageNumber++
		instruments, _ := a.services.InstrumentService.GetList()
		instrumentsLen := len(instruments)
		instrumentMaxPageNumber := instrumentsLen / instrumentsPerPage
		if instrumentPageNumber > instrumentMaxPageNumber {
			instrumentPageNumber = instrumentMaxPageNumber
		}
	}))
	instrumentSubMenu.AddMenuItem(GoConsoleMenu.NewActionItem(4, "First page", func() {
		instrumentPageNumber = 0
	}))
	instrumentSubMenu.AddMenuItem(GoConsoleMenu.NewActionItem(5, "Last page", func() {
		instruments, _ := a.services.InstrumentService.GetList()
		instrumentsLen := len(instruments)
		instrumentPageNumber = instrumentsLen / instrumentsPerPage
	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(1, "Register", func() {
		Register(a)
	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(2, "Login", func() {
		user = Login(a)
		guestFuncs = user == nil
		userFuncs = user != nil
		adminFuncs = user != nil && user.IsAdmin
	}))

	Menu.AddMenuItem(GoConsoleMenu.NewSubmenuItem(3, "Instruments List", instrumentSubMenu))

	discountsPerPage := 10
	discountPageNumber := 0
	discountSubMenu := GoConsoleMenu.NewMenu("Navigating through pages")
	discountSubMenu.AddMenuItem(GoConsoleMenu.NewActionItem(0, "Exit discounts list", func() { discountPageNumber = 0 }).SetAsExitOption())
	discountSubMenu.AddMenuItem(GoConsoleMenu.NewActionItem(1, "Show list", func() {
		discounts := a.handlers.DiscountHandler.GetList(discountPageNumber*discountsPerPage, discountsPerPage)
		fmt.Println(discounts)
	}))
	discountSubMenu.AddMenuItem(GoConsoleMenu.NewActionItem(2, "Previous page", func() {
		discountPageNumber--
		if discountPageNumber < 0 {
			discountPageNumber = 0
		}
	}))
	discountSubMenu.AddMenuItem(GoConsoleMenu.NewActionItem(3, "Next page", func() {
		discountPageNumber++
		discounts, _ := a.services.DiscountService.GetList()
		discountsLen := len(discounts)
		discountMaxPageNumber := discountsLen / discountsPerPage
		if discountPageNumber > discountMaxPageNumber {
			discountPageNumber = discountMaxPageNumber
		}
	}))
	discountSubMenu.AddMenuItem(GoConsoleMenu.NewActionItem(4, "First page", func() {
		discountPageNumber = 0
	}))
	discountSubMenu.AddMenuItem(GoConsoleMenu.NewActionItem(5, "Last page", func() {
		discounts, _ := a.services.DiscountService.GetList()
		discountsLen := len(discounts)
		discountPageNumber = discountsLen / discountsPerPage
	}))
	Menu.AddMenuItem(GoConsoleMenu.NewSubmenuItem(4, "Discounts List", discountSubMenu))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(5, "Comparison List", func() {
		comparisonList := a.handlers.UserHandler.GetComparisonList(user.UserId)
		fmt.Println(comparisonList)
	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(6, "Add instrument to Comparison List", func() {
		AddInstrumentToComparisonList(a, user.UserId)
	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(7, "Delete instrument from Comparison List", func() {
		DeleteInstrumentFromComparisonList(a, user.UserId)
	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(8, "Add instrument to Data Base", func() {
		CreateInstrument(a, user.Login)
	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(9, "Delete instrument from Data Base", func() {
		DeleteInstrument(a, user.Login)
	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(10, "Update instrument in Data Base", func() {
		UpdateInstrument(a, user.Login)
	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(11, "Add discount in Data Base", func() {
		CreateDiscount(a, user.Login)
	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(12, "Add discount for all users in Data Base", func() {
		CreateDiscountForAllUsers(a, user.Login)
	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(13, "Delete discount from Data Base", func() {
		DeleteDiscount(a, user.Login)
	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(14, "Update discount in Data Base", func() {
		UpdateDiscount(a, user.Login)
	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(15, "Checkout", func() {
		Checkout(a, user.UserId)
	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(16, "Orders List", func() {
		var orderStr string
		if user.IsAdmin {
			orderStr = a.handlers.OrderHandler.GetListForAll()
		} else {
			orderStr = a.handlers.OrderHandler.GetList(user.UserId)
		}
		var ordersMap map[string][]models.Order
		err := json.Unmarshal([]byte(orderStr), &ordersMap)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		orders := ordersMap["orders"]
		fmt.Println(orderStr)
		for i := range orders {
			orderElements := a.handlers.OrderHandler.GetOrderElements(orders[i].OrderId)
			fmt.Println(orderElements)
		}
	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(17, "Update Order Status", func() {
		UpdateOrderStatus(a, user.Login)
	}))

	Menu.Display()
}
