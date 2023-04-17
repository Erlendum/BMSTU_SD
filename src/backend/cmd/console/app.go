package console

import (
	"backend/internal/cli/handlers"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/repository/postgres_repository"
	"backend/internal/services"
	servicesImplementation "backend/internal/services/implementation"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/lukewarlow/GoConsoleMenu"
)

type App struct {
	repositories *appRepositoryFields
	services     *appServiceFields
	handlers     *handlers.Handlers
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
	f := &appRepositoryFields{
		comparisonListRepository: postgres_repository.CreateComparisonListPostgresRepository(fields),
		discountRepository:       postgres_repository.CreateDiscountPostgresRepository(fields),
		instrumentRepository:     postgres_repository.CreateInstrumentPostgresRepository(fields),
		userRepository:           postgres_repository.CreateUserPostgresRepository(fields),
	}

	return f
}

func (a *App) initServices(r *appRepositoryFields) *appServiceFields {
	calcDiscountService := servicesImplementation.NewCalcDiscountServiceImplementation(r.discountRepository)
	u := &appServiceFields{
		CalcDiscountService:   calcDiscountService,
		ComparisonListService: servicesImplementation.NewComparisonListServiceImplementation(r.comparisonListRepository, r.instrumentRepository),
		DiscountService:       servicesImplementation.NewDiscountServiceImplementation(r.discountRepository, r.userRepository),
		InstrumentService:     servicesImplementation.NewInstrumentServiceImplementation(r.instrumentRepository, r.userRepository),
		UserService:           servicesImplementation.NewUserServiceImplementation(r.userRepository, r.comparisonListRepository, calcDiscountService),
	}

	return u
}

func (a *App) Init() {
	a.repositories = a.initRepositories()
	a.services = a.initServices(a.repositories)
	handlerServices := handlers.HandlersServicesFields{}
	copier.Copy(&handlerServices, a.services)
	a.handlers = handlers.NewHandlers(handlerServices)
}

func (a *App) Run() {
	a.Init()

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
		}

		if adminFuncs {
			menu.ShowMenuItem(8)
			menu.ShowMenuItem(9)
			menu.ShowMenuItem(10)
			menu.ShowMenuItem(11)
			menu.ShowMenuItem(12)
			menu.ShowMenuItem(13)
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
	discountSubMenu.AddMenuItem(GoConsoleMenu.NewActionItem(0, "Exit discounts list", func() { discountsPerPage = 0 }).SetAsExitOption())
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

	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(6, "Add instrument to Comparison List", func() {
	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(7, "Delete instrument from Comparison List", func() {
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

	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(12, "Delete discount from Data Base", func() {

	}))

	Menu.AddHiddenMenuItem(GoConsoleMenu.NewActionItem(13, "Update discount in Data Base", func() {

	}))
	Menu.Display()
}
