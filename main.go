package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/alaingilbert/ogame/pkg/device"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/wrapper"
	"github.com/alaingilbert/ogame/pkg/wrapper/solvers"
)

var botToken string
var chatID string
var boot *wrapper.OGame

func getFlottePourExpe(bot *wrapper.OGame) {
	att, _ := bot.IsUnderAttack()
	if att {
		info, _ := bot.GetAttacks()
		sendTelegramMessage(botToken, chatID, "ATTAQUE EN COURS!")
		sendTelegramMessage(botToken, chatID, fmt.Sprint(info))
	}

	/*getMaxExpeDebris(4)
	getMaxExpeDebris(5)
	getMaxExpeDebris(6)
	os.Exit(0)*/

	CargoExpeInsuffisant := 0

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
	empire = sliceEmpireCargo(empire)
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
	//validCoordLF = readJSONCoordFdV()
	//fmt.Println(validCoordLF)
	//os.Exit(0)
	//HasMoonRes := false
	for i, planete := range empire {
		fmt.Printf("======================= planete %s(%s) =========================\n", planete.Name, planete.Coordinate)

		if planete.Type == ogame.MoonType {
			buildMoon(planete, bot)
			if slots.ExpInUse >= slots.ExpTotal && slots.InUse < slots.Total {
				sendFleetToMoon(planete)
				if sendFleetFromMoonToPlanet(planete) {
					//HasMoonRes = true
				}
			}
		} else if planete.Fields.Built < planete.Fields.Total-2 {
			buildResources(planete)
		} else {
			bot.BuildBuilding(planete.ID, ogame.TerraformerID)
		}

		if planete.Facilities.ResearchLab < 12 && i == 0 {
			bot.BuildBuilding(planete.ID, ogame.ResearchLabID)
		}

		buildFormeVie(planete)
		if !SetExpedition(planete, bot, coordExpe) {
			CargoExpeInsuffisant++
		}
		printCurrentconstruction(planete.ID, bot)
	}

	Researches(first_planet, bot, slots)

	_, slots = bot.GetFleets()
	if slots.ExpInUse < slots.ExpTotal && slots.InUse < slots.Total {
		empire = sliceEmpireCargo(empire)
		for _, planete := range empire {
			co := planete.Coordinate
			co.Position = 16
			_, err := bot.SendFleet(planete.ID, getCompoFlotteExpe(planete), 100, co, ogame.Expedition, ogame.Resources{}, 0, 0)
			if err != nil {
				fmt.Printf("Erreur envoie expe restant : %s\n", err)
				if strings.Contains(err.Error(), "all slots are in use") {
					break
				}
			} else {

			}

			time.Sleep(4 * time.Second)
		}

		CargoExpeInsuffisant = 0
	}

	if slots.ExpInUse >= slots.ExpTotal && slots.InUse < slots.Total /*&& !HasMoonRes*/ {
		sort.Slice(empire, func(i int, j int) bool {
			resources_i := empire[i].Resources.Metal + empire[i].Resources.Crystal + empire[i].Resources.Deuterium
			resources_j := empire[j].Resources.Metal + empire[j].Resources.Crystal + empire[j].Resources.Deuterium
			HasValidForExplo := empire[i].Resources.Metal > 5000 && empire[i].Resources.Crystal > 5000 && empire[i].Resources.Deuterium > 5000
			return resources_i > resources_j && HasValidForExplo
		})
		setExploVie(empire[0].ID, empire[0].Coordinate, bot)
	}
}

func main() {
	universe := os.Getenv("UNIVERSE")
	username := os.Getenv("USERNAME") // eg: email@gmail.com
	password := os.Getenv("PASSWORD")
	language := os.Getenv("LANGUAGE")
	botToken = os.Getenv("BOTTOKEN")
	chatID = os.Getenv("CHATID") // Exemple : "123456789"
	fmt.Printf("Paramètres utilisateur récupéré => univers: %s username: %s mdp:%s language: %s\n", universe, username, password, language)
	universe = "Quasi-Stellar"

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
