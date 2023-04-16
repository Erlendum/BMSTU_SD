package console

import (
	"backend/internal/models"
	"bufio"
	"fmt"
	"github.com/howeyc/gopass"
	"os"
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
