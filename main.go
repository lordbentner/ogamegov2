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

var botToken string
var chatID string
var boot *wrapper.OGame

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

func getFlottePourExpe(bot *wrapper.OGame) {
	att, _ := bot.IsUnderAttack()
	if att {
		sendTelegramMessage(botToken, chatID, "ATTAQUE EN COURS!")
	}

	/*getMaxExpeDebris(4)
	getMaxExpeDebris(5)
	getMaxExpeDebris(6)
	os.Exit(0)*/

	fleets, slots := bot.GetFleets()
	fmt.Println("=====================Flottes=======================")
	fmt.Printf("%s slots: ", time.Now().Format(time.RFC850))
	fmt.Println(slots)
	for i, fleet := range fleets {
		fmt.Printf("flotte %d ==> ", i)
		printStructFields(fleet.Ships)
	}
	fmt.Println("====================================================")
	empire, _ := bot.GetEmpire(ogame.PlanetType)
	if len(empire) == 0 {
		fmt.Println(empire)
		return
	}

	first_planet := empire[0]

	empireMoon, _ := bot.GetEmpire(ogame.MoonType)
	empire = append(empire, empireMoon...)

	sort.Slice(empire, func(i int, j int) bool {
		cargoTotali := getCargoGT()*empire[i].Ships.LargeCargo + getCargoPT()*empire[i].Ships.SmallCargo
		cargoTotali += getCargoPathFinder() * empire[i].Ships.Pathfinder
		cargoTotalj := getCargoGT()*empire[j].Ships.LargeCargo + getCargoPT()*empire[j].Ships.SmallCargo
		cargoTotalj += getCargoPathFinder() * empire[j].Ships.Pathfinder
		return cargoTotali > cargoTotalj
	})

	expeMes := gestionMessagesExpe(bot)
	coordExpe := expeMes.Coordinate
	fmt.Println(coordExpe)
	if changeSystemeExploration(expeMes.Content) {
		coordMain := empire[0].Coordinate
		sys := coordExpe.System + 1
		if sys > coordMain.System+10 {
			sys = coordMain.System
		}

		coordExpe = ogame.Coordinate{Galaxy: coordExpe.Galaxy, System: sys, Position: 16}
	}

	fmt.Println("================================================================")
	Researches(first_planet, bot, slots)
	//HasMoonRes := false
	for _, planete := range empire {
		fmt.Printf("======================= planete %s(%s) =========================\n", planete.Name, planete.Coordinate)
		if planete.Facilities.ResearchLab < 12 {
			bot.BuildBuilding(planete.ID, ogame.ResearchLabID)
		}

		if planete.Type == ogame.MoonType {
			buildMoon(planete, bot)
			if slots.ExpInUse >= slots.ExpTotal && slots.InUse < slots.Total {
				if sendFleetFromMoonToPlanet(planete) {
					//HasMoonRes = true
				}
			}
		} else if planete.Fields.Built < planete.Fields.Total-2 {
			buildResources(planete)
		} else {
			bot.BuildBuilding(planete.ID, ogame.TerraformerID)
		}

		buildFormeVie(planete)
		SetExpedition(planete, bot, coordExpe)
		printCurrentconstruction(planete.ID, bot)
	}

	/*if slots.ExpInUse >= slots.ExpTotal && slots.InUse < slots.Total && !HasMoonRes {
		setExploVie(planetLife.ID, planetLife.Coordinate, bot)
	}*/
	//setExploVie(planetLife.ID, planetLife.Coordinate, bot)
}

func main() {
	universe := os.Getenv("UNIVERSE")
	username := os.Getenv("USERNAME") // eg: email@gmail.com
	password := os.Getenv("PASSWORD")
	language := os.Getenv("LANGUAGE")
	botToken = os.Getenv("BOTTOKEN")
	chatID = os.Getenv("CHATID") // Exemple : "123456789"
	startDate := time.Now()
	fmt.Println(startDate)
	fmt.Println(startDate.Add(5 * time.Minute))
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
	bot.Logout()
}

func connect(bot *wrapper.OGame) bool {
	fmt.Printf("%s Connexion", time.Now().Format(time.RFC850))
	bot.LoginWithExistingCookies()
	boot = bot
	time.Sleep(5000)
	getFlottePourExpe(bot)
	bot.Logout()
	return connect(bot)
}
