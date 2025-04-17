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

func buildResources(planete ogame.EmpireCelestial, bot *wrapper.OGame) {
	fmt.Println(planete.Supplies)
	bot.BuildBuilding(planete.ID, ogame.RoboticsFactoryID)
	if planete.Resources.Energy < 0 {
		bot.BuildBuilding(planete.ID, ogame.SolarPlantID)
		fmt.Println("construction solar plant ")
	} else if planete.Supplies.DeuteriumSynthesizer < int64(planete.Supplies.CrystalMine-4) {
		bot.BuildBuilding(planete.ID, ogame.DeuteriumSynthesiserID)
		fmt.Println("construction synthethiseur deuterium ")
	} else if planete.Supplies.CrystalMine < int64(planete.Supplies.MetalMine-4) {
		bot.BuildBuilding(planete.ID, ogame.CrystalMineID)
		fmt.Println("construction crystal mine")
	} else {
		bot.BuildBuilding(planete.ID, ogame.MetalMineID)
		fmt.Println("construction metal mine")
		bot.BuildBuilding(planete.ID, ogame.MetalStorageID)
		bot.BuildBuilding(planete.ID, ogame.CrystalStorageID)
		bot.BuildBuilding(planete.ID, ogame.DeuteriumTankID)
	}
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
	//	planetLife := empire[0]
	//trouve := false
	for _, planete := range empire {
		SetExpedition(planete.ID, planete.Coordinate, bot)
		buildResources(planete, bot)
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

	deviceName := "DEVICE-NAME"
	deviceInst, err := device.NewBuilder(deviceName).
		SetOsName(device.Windows).
		SetBrowserName(device.Chrome).
		SetMemory(8).
		SetHardwareConcurrency(16).
		ScreenColorDepth(24).
		SetScreenWidth(1900).
		SetScreenHeight(900).
		SetTimezone("Europe/Los_Angeles").
		SetLanguages("en-US,en").
		Build()

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
	getFlottePourExpe(bot)

	bot.Logout()
	time.Sleep(1 * time.Minute)
	return connect(bot)
}
