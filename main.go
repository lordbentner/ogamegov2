package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/alaingilbert/ogame/pkg/device"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/wrapper"
	"github.com/alaingilbert/ogame/pkg/wrapper/solvers"
)

func printFleets(fleet ogame.Fleet) {
	fmt.Print(fleet.Ships.LargeCargo)
	fmt.Print(" GT,  ")
	fmt.Print(fleet.Ships.SmallCargo)
	fmt.Print(" PT, ")
	fmt.Print(fleet.Ships.Pathfinder)
	fmt.Println(" Eclaireur")
}

func printCurrentconstruction(id ogame.CelestialID, bot *wrapper.OGame) {
	buildingID, buildingCountdown, researchID, researchCountdown, lfBuildingID,
		lfBuildingCountdown, lfResearchID, lfResearchCountdown := bot.ConstructionsBeingBuilt(id)

	print_str := "Construction en cours  ==> "
	if buildingCountdown > 0 {
		print_str += fmt.Sprintf("buildingID = %s, buildingCountdown = %d ", buildingID, buildingCountdown)
	}
	if researchCountdown > 0 {
		print_str += fmt.Sprintf("researchID = %s, researchCountdown = %d ", researchID, researchCountdown)
	}
	if lfBuildingCountdown > 0 {
		print_str += fmt.Sprintf("lfBuildingID = %s, lfBuildingCountdown = %d ", lfBuildingID, lfBuildingCountdown)
	}
	if lfResearchCountdown > 0 {
		print_str += fmt.Sprintf("lfResearchID = %s, lfResearchCountdown = %d ", lfResearchID, lfResearchCountdown)
	}

	fmt.Println(print_str)
}

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

	/*lfBonuses, _ := bot.GetCachedLfBonuses()
	multiplier := float64(bot.GetServerData().CargoHyperspaceTechMultiplier) / 100.0
	cargo := empire[0].Ships.Cargo(bot.GetCachedResearch(), lfBonuses, bot.CharacterClass(), multiplier, bot.GetServer().ProbeRaidsEnabled())
	cargoExpe := cargo / slots.ExpTotal
	cargoGT := ogame.LargeCargo.GetCargoCapacity(bot.GetCachedResearch(), lfBonuses, bot.CharacterClass(), multiplier, bot.GetServer().ProbeRaidsEnabled())
	fmt.Printf("cargo total =%d, cargo par expe = %d\n, cargo GT = %d", cargo, cargoExpe, cargoGT)*/
}

func main() {
	universe := os.Getenv("UNIVERSE")
	username := os.Getenv("USERNAME") // eg: email@gmail.com
	password := os.Getenv("PASSWORD")
	language := os.Getenv("LANGUAGE")

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
	//sendMail(bot)
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
