package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"

	"github.com/alaingilbert/ogame/pkg/device"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/wrapper"
	"github.com/alaingilbert/ogame/pkg/wrapper/solvers"
)

func getFlottePourExpe(bot *wrapper.OGame) {

	slots, _ := bot.GetSlots()
	fmt.Printf("%s slots: ", time.Now().Format(time.RFC850))
	fmt.Println(slots)
	empire, _ := bot.GetEmpire(ogame.PlanetType)
	empireMoon, _ := bot.GetEmpire(ogame.MoonType)
	empire = append(empire, empireMoon...)

	sort.Slice(empire, func(i int, j int) bool {
		return empire[i].Ships.LargeCargo > empire[j].Ships.LargeCargo
	})

	if len(empire) == 0 {
		fmt.Println(empire)
		return
	}

	Researches(empire[0], bot, slots)
	buildFormeVie(empire[0], bot)
	//	planetLife := empire[0]
	//trouve := false
	for _, planete := range empire {
		buildResources(planete, bot)
		SetExpedition(planete.ID, planete.Coordinate, bot)
		printCurrentconstruction(planete.ID, bot)
		/*if planete.Type == ogame.PlanetType && !trouve {
			fmt.Println(planete.Resources)
			planetLife = planete
			trouve = true
		}*/
	}

	//setExploVie(planetLife.ID, planetLife.Coordinate, bot, 0)
}

func sendTelegramMessage(token, chatID, message string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id": {chatID},
		"text":    {message},
	})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("✅ Message envoyé avec succès.")
}

var botToken string
var chatID string

func main() {
	universe := os.Getenv("UNIVERSE")
	username := os.Getenv("USERNAME") // eg: email@gmail.com
	password := os.Getenv("PASSWORD")
	language := os.Getenv("LANGUAGE")
	botToken = os.Getenv("BOTTOKEN")
	chatID = os.Getenv("CHATID") // Exemple : "123456789"
	//message := "📣 Notification envoyée depuis un programme Go sous Windows vers Telegram !"
	startDate := time.Now()
	fmt.Println(startDate)
	fmt.Println(startDate.Add(5 * time.Minute))

	//sendTelegramMessage(botToken, chatID, message)
	fmt.Printf("Paramètres utilisateur récupéré => univers: %s username: %s mdp:%s language: %s\n", universe, username, password, language)

	deviceName := "DEVICE-NAME"
	deviceInst, err := device.NewBuilder(deviceName).
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

	params := wrapper.Params{
		Universe:        universe,
		Username:        username,
		Password:        password,
		Lang:            language,
		AutoLogin:       false,
		Device:          deviceInst,
		CaptchaCallback: solvers.ManualSolver(),
	}
	bot, err := wrapper.NewWithParams(params)

	if err != nil {
		panic(err)
	}

	ff, _ := bot.LoginWithExistingCookies()
	if !ff {
		bot.Login()
	}

	connect(bot)
	//	getFlottePourExpe(bot)
	bot.Logout()
}

func connect(bot *wrapper.OGame) bool {
	fmt.Printf("%s Connexion", time.Now().Format(time.RFC850))
	bot.LoginWithExistingCookies()
	time.Sleep(5000)
	//if bot.IsConnected() {
	getFlottePourExpe(bot)
	//}

	bot.Logout()
	time.Sleep(10 * time.Second)
	return connect(bot)
}
