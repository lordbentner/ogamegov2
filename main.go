package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alaingilbert/ogame/pkg/device"
	"github.com/alaingilbert/ogame/pkg/gameforge/solvers"
	"github.com/alaingilbert/ogame/pkg/wrapper"
)

var botToken string
var chatID string

func gestionAttack() {
	attacks, _ := boot.GetAttacks()
	//empire, _ := boot.GetEmpire()
	sendTelegramMessage(botToken, chatID, "ATTAQUE EN COURS!")
	sendTelegramMessage(botToken, chatID, fmt.Sprint(attacks))

	/*for _, attack := range attacks {
		for _, planete := range empire {
			if planete.Coordinate == attack.Destination {
				fmt.Println("attack ")
				break
			}
		}
	}*/
}

func main() {
	universe := os.Getenv("UNIVERSE")
	username := os.Getenv("USERNAME") // eg: email@gmail.com
	password := os.Getenv("PASSWORD")
	language := os.Getenv("LANGUAGE")
	botToken = os.Getenv("BOTTOKEN")
	chatID = os.Getenv("CHATID") // Exemple : "123456789"
	fmt.Printf("Paramètres utilisateur récupéré => univers: %s, username: %s, mdp:%s, language: %s\n", universe, username, password, language)

	deviceInst, err := device.NewBuilder("device_name").
		SetOsName(device.Windows).
		SetBrowserName(device.Chrome).
		SetMemory(8).
		SetHardwareConcurrency(16).
		ScreenColorDepth(24).
		SetScreenWidth(1900).
		SetScreenHeight(900).
		SetTimezone("America/Los_Angeles").
		SetLanguages("en-US,en").
		Build()

	if err != nil {
		panic(err)
	}

	b, err := wrapper.NewWithParams(wrapper.Params{
		Universe:      universe,
		Username:      username,
		Password:      password,
		Lang:          language,
		AutoLogin:     true,
		Device:        deviceInst,
		CaptchaSolver: solvers.ManualSolver(),
	})

	if err != nil {
		panic(err)
	}

	ff, _, err := b.LoginWithExistingCookies()
	if !ff {
		fmt.Println("Lgin Cookiees failed : ")
		fmt.Println(err)
		b.Login()
	}

	connect(b)
	boot.Logout()
}

func connect(bot *wrapper.OGame) bool {
	fmt.Printf("%s Connexion", time.Now().Format(time.RFC850))
	bot.LoginWithExistingCookies()
	boot = bot
	if bot.IsConnected() && bot.IsLoggedIn() {
		time.Sleep(5000)
		getFlottePourExpe()
		boot.Logout()
	}
	return connect(bot)
}
