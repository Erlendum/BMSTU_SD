package console

import (
	"backend/internal/models"
	"bufio"
	"fmt"
	"github.com/howeyc/gopass"
	"os"
	"strconv"
	"time"
)

func Register(a *App) {
	scanner := bufio.NewScanner(os.Stdin)

	tmpUser := models.User{}
	fmt.Print("Input login: ")

	if !scanner.Scan() {
		fmt.Printf("Invalid login")
		return
	}
	tmpUser.Login = scanner.Text()

	fmt.Print("Input password:")
	password, err := gopass.GetPasswd()
	if err != nil {
		fmt.Printf("Invalid password")
		return
	}
	tmpUser.Password = string(password)

	fmt.Print("Input your lastname, firstname and middlename: ")
	if !scanner.Scan() {
		fmt.Printf("Invalid lastname, firstname and middlename")
		return
	}
	tmpUser.Fio = scanner.Text()

	fmt.Print("Input your gender: ")
	if !scanner.Scan() {
		fmt.Printf("Invalid gender")
		return
	}
	var gender models.UserGender
	if scanner.Text() == "Мужской" || scanner.Text() == "Женский" {
		gender = models.UserGender(scanner.Text())
		tmpUser.Gender = gender
	} else {
		fmt.Printf("Invalid gender")
		return
	}

	fmt.Print("Input your date of birth (format: year-month-day): ")
	if !scanner.Scan() {
		fmt.Printf("Invalid date of birth")
		return
	}
	layOut := "2006-01-02"
	dateBirth, err := time.Parse(layOut, scanner.Text())
	if err != nil {
		fmt.Printf("Invalid date of birth")
		return
	}

	tmpUser.DateBirth = dateBirth

	res := a.handlers.UserHandler.Create(tmpUser)
	if res == "0" {
		fmt.Print("Success")
	} else {
		fmt.Print(res)
	}
}

func Login(a *App) *models.User {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Input login: ")

	if !scanner.Scan() {
		fmt.Printf("Invalid login")
		return nil
	}
	login := scanner.Text()

	fmt.Print("Input password:")
	password, err := gopass.GetPasswd()
	if err != nil {
		fmt.Printf("Invalid password")
		return nil
	}
	passwordStr := string(password)

	user, res := a.handlers.UserHandler.Get(login, passwordStr)
	if res == "0" {
		fmt.Print("Success")
	} else {
		fmt.Print(res)
	}

	return user
}

func CreateInstrument(a *App, login string) {
	scanner := bufio.NewScanner(os.Stdin)

	tmpInstrument := models.Instrument{}
	fmt.Print("Input name: ")
	if !scanner.Scan() {
		fmt.Printf("Invalid name")
		return
	}
	tmpInstrument.Name = scanner.Text()

	fmt.Print("Input price:")
	if !scanner.Scan() {
		fmt.Printf("Invalid price")
		return
	}
	price, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Printf("Invalid price")
		return
	}
	tmpInstrument.Price = uint64(price)

	fmt.Print("Input material: ")
	if !scanner.Scan() {
		fmt.Printf("Invalid material")
		return
	}
	tmpInstrument.Material = scanner.Text()

	fmt.Print("Input type: ")
	if !scanner.Scan() {
		fmt.Printf("Invalid type")
		return
	}
	tmpInstrument.Type = scanner.Text()

	fmt.Print("Input brand: ")
	if !scanner.Scan() {
		fmt.Printf("Invalid brand")
		return
	}
	tmpInstrument.Brand = scanner.Text()

	fmt.Print("Input Img: ")
	if !scanner.Scan() {
		fmt.Printf("Invalid Img")
		return
	}
	tmpInstrument.Img = scanner.Text()

	res := a.handlers.InstrumentHandler.Create(tmpInstrument, login)
	if res == "0" {
		fmt.Print("Success")
	} else {
		fmt.Print(res)
	}
}

func DeleteInstrument(a *App, login string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Input id:")
	if !scanner.Scan() {
		fmt.Printf("Invalid id")
		return
	}
	id, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Printf("Invalid id")
		return
	}
	res := a.handlers.InstrumentHandler.Delete(uint64(id), login)
	if res == "0" {
		fmt.Print("Success")
	} else {
		fmt.Print(res)
	}
}

func UpdateInstrument(a *App, login string) {
	fields := models.InstrumentFieldsToUpdate{}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Input id:")
	if !scanner.Scan() {
		fmt.Printf("Invalid id")
		return
	}
	id, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Printf("Invalid id")
		return
	}

	fmt.Print("Input name to update (or press enter): ")
	if !scanner.Scan() {
		fmt.Printf("Invalid name")
		return
	}
	name := scanner.Text()
	if name != "" {
		fields[models.InstrumentFieldName] = name
	}

	fmt.Print("Input price (or press enter): ")
	if !scanner.Scan() {
		fmt.Printf("Invalid price")
		return
	}
	if scanner.Text() != "" {
		price, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Printf("Invalid price")
			return
		}
		fields[models.InstrumentFieldPrice] = uint64(price)
	}

	fmt.Print("Input material (or press enter): ")
	if !scanner.Scan() {
		fmt.Printf("Invalid material")
		return
	}
	material := scanner.Text()
	if material != "" {
		fields[models.InstrumentFieldMaterial] = material
	}

	fmt.Print("Input type (or press enter): ")
	if !scanner.Scan() {
		fmt.Printf("Invalid type")
		return
	}
	typeI := scanner.Text()
	if typeI != "" {
		fields[models.InstrumentFieldType] = typeI
	}

	fmt.Print("Input brand (or press enter): ")
	if !scanner.Scan() {
		fmt.Printf("Invalid brand")
		return
	}
	brand := scanner.Text()
	if brand != "" {
		fields[models.InstrumentFieldBrand] = brand
	}

	fmt.Print("Input Img (or press enter): ")
	if !scanner.Scan() {
		fmt.Printf("Invalid Img")
		return
	}
	img := scanner.Text()
	if img != "" {
		fields[models.InstrumentFieldImg] = img
	}

	res := a.handlers.InstrumentHandler.Update(uint64(id), login, fields)
	if res == "0" {
		fmt.Print("Success")
	} else {
		fmt.Print(res)
	}
}

func AddInstrumentToComparisonList(a *App, userId uint64) {
	scanner := bufio.NewScanner(os.Stdin)
	comparisonList, _, err := a.services.UserService.GetComparisonList(userId)
	if err != nil {
		fmt.Print(err)
		return
	}
	comparisonListId := comparisonList.ComparisonListId

	fmt.Print("Input instrument id:")
	if !scanner.Scan() {
		fmt.Printf("Invalid id")
		return
	}

	id, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Printf("Invalid id")
		return
	}

	res := a.handlers.ComparisonListHandler.AddInstrument(comparisonListId, uint64(id))
	if res == "0" {
		fmt.Print("Success")
	} else {
		fmt.Print(res)
	}
}

func DeleteInstrumentFromComparisonList(a *App, userId uint64) {
	scanner := bufio.NewScanner(os.Stdin)
	comparisonList, _, err := a.services.UserService.GetComparisonList(userId)
	if err != nil {
		fmt.Print(err)
		return
	}
	comparisonListId := comparisonList.ComparisonListId

	fmt.Print("Input instrument id:")
	if !scanner.Scan() {
		fmt.Printf("Invalid id")
		return
	}

	id, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Printf("Invalid id")
		return
	}

	res := a.handlers.ComparisonListHandler.DeleteInstrument(comparisonListId, uint64(id))
	if res == "0" {
		fmt.Print("Success")
	} else {
		fmt.Print(res)
	}
}

func CreateDiscount(a *App, login string) {
	scanner := bufio.NewScanner(os.Stdin)

	tmpDiscount := models.Discount{}

	fmt.Print("Input instrument id:")
	if !scanner.Scan() {
		fmt.Printf("Invalid instrument id")
		return
	}
	instrumentId, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Printf("Invalid instrument id")
		return
	}
	tmpDiscount.InstrumentId = uint64(instrumentId)

	fmt.Print("Input user id:")
	if !scanner.Scan() {
		fmt.Printf("Invalid user id")
		return
	}
	userId, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Printf("Invalid user id")
		return
	}
	tmpDiscount.UserId = uint64(userId)

	fmt.Print("Input amount:")
	if !scanner.Scan() {
		fmt.Printf("Invalid amount")
		return
	}
	amount, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Printf("Invalid amount")
		return
	}
	tmpDiscount.Amount = uint64(amount)

	fmt.Print("Input type: ")
	if !scanner.Scan() {
		fmt.Printf("Invalid type")
		return
	}
	tmpDiscount.Type = scanner.Text()

	layOut := "2006-01-02"
	fmt.Print("Input begin date (format: year-month-day): ")
	if !scanner.Scan() {
		fmt.Printf("Invalid begin date")
		return
	}
	dateBegin, err := time.Parse(layOut, scanner.Text())
	if err != nil {
		fmt.Printf("Invalid begin date")
		return
	}
	tmpDiscount.DateBegin = dateBegin

	fmt.Print("Input end date (format: year-month-day): ")
	if !scanner.Scan() {
		fmt.Printf("Invalid end date")
		return
	}
	dateEnd, err := time.Parse(layOut, scanner.Text())
	if err != nil {
		fmt.Printf("Invalid end date")
		return
	}
	tmpDiscount.DateEnd = dateEnd

	res := a.handlers.DiscountHandler.Create(tmpDiscount, login)
	if res == "0" {
		fmt.Print("Success")
	} else {
		fmt.Print(res)
	}
}

func DeleteDiscount(a *App, login string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Input id:")
	if !scanner.Scan() {
		fmt.Printf("Invalid id")
		return
	}
	id, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Printf("Invalid id")
		return
	}
	res := a.handlers.DiscountHandler.Delete(uint64(id), login)
	if res == "0" {
		fmt.Print("Success")
	} else {
		fmt.Print(res)
	}
}

func UpdateDiscount(a *App, login string) {
	fields := models.DiscountFieldsToUpdate{}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Input id:")
	if !scanner.Scan() {
		fmt.Printf("Invalid id")
		return
	}
	id, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Printf("Invalid id")
		return
	}

	fmt.Print("Input instrument id (or press enter): ")
	if !scanner.Scan() {
		fmt.Printf("Invalid instrument id")
		return
	}
	if scanner.Text() != "" {
		instrumentId, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Printf("Invalid instrument id")
			return
		}
		fields[models.DiscountFieldInstrumentId] = uint64(instrumentId)
	}

	fmt.Print("Input user id (or press enter): ")
	if !scanner.Scan() {
		fmt.Printf("Invalid user id")
		return
	}
	if scanner.Text() != "" {
		userId, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Printf("Invalid user id")
			return
		}
		fields[models.DiscountFieldUserId] = uint64(userId)
	}

	fmt.Print("Input amount (or press enter): ")
	if !scanner.Scan() {
		fmt.Printf("Invalid amount")
		return
	}
	if scanner.Text() != "" {
		amount, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Printf("Invalid amount")
			return
		}
		fields[models.DiscountFieldAmount] = uint64(amount)
	}

	fmt.Print("Input type (or press enter): ")
	if !scanner.Scan() {
		fmt.Printf("Invalid type")
		return
	}
	typeD := scanner.Text()
	if typeD != "" {
		fields[models.DiscountFieldType] = typeD
	}

	layOut := "2006-01-02"
	fmt.Print("Input begin date (format: year-month-day): ")
	if !scanner.Scan() {
		fmt.Printf("Invalid begin date")
		return
	}
	if scanner.Text() != "" {
		dateBegin, err := time.Parse(layOut, scanner.Text())
		if err != nil {
			fmt.Printf("Invalid begin date")
			return
		}
		fields[models.DiscountFieldDateBegin] = dateBegin
	}

	fmt.Print("Input end date (format: year-month-day): ")
	if !scanner.Scan() {
		fmt.Printf("Invalid end date")
		return
	}
	if scanner.Text() != "" {
		dateEnd, err := time.Parse(layOut, scanner.Text())
		if err != nil {
			fmt.Printf("Invalid end date")
			return
		}
		fields[models.DiscountFieldDateEnd] = dateEnd
	}

	res := a.handlers.DiscountHandler.Update(uint64(id), login, fields)
	if res == "0" {
		fmt.Print("Success")
	} else {
		fmt.Print(res)
	}
}
