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
	Researches(planete, bot)
	time.Sleep(10000)
	fmt.Println(planete.Supplies)
	resDetails, _ := bot.GetResourcesDetails(planete.ID)
	fmt.Println("Metal Storage Capacity:")
	fmt.Println(resDetails.Metal.StorageCapacity)
	//bot.BuildBuilding(planete.ID, ogame.RoboticsFactoryID)
	if planete.Resources.Energy < 0 {
		bot.BuildBuilding(planete.ID, ogame.SolarPlantID)
		fmt.Println("construction solar plant ")
	} else if planete.Supplies.DeuteriumSynthesizer < int64(planete.Supplies.CrystalMine-3) {
		err := bot.BuildBuilding(planete.ID, 3)
		fmt.Printf("construction synthethiseur deuterium err =%s\n", err)
	} else if planete.Supplies.CrystalMine < int64(planete.Supplies.MetalMine-3) {
		err := bot.BuildBuilding(planete.ID, 2)
		fmt.Printf("construction crystal mine err =%s\n", err)
	} else {
		err := bot.BuildBuilding(planete.ID, 1)
		fmt.Printf("construction metal mine err =%s\n", err)
		/*bot.BuildBuilding(planete.ID, ogame.MetalStorageID)
		bot.BuildBuilding(planete.ID, ogame.CrystalStorageID)
		bot.BuildBuilding(planete.ID, ogame.DeuteriumTankID)*/
	}

	if resDetails.Crystal.StorageCapacity-resDetails.Crystal.StorageCapacity/10 < planete.Resources.Crystal {
		bot.BuildBuilding(planete.ID, ogame.CrystalStorageID)
	} else if resDetails.Deuterium.StorageCapacity-resDetails.Deuterium.StorageCapacity/10 < planete.Resources.Deuterium {
		bot.BuildBuilding(planete.ID, ogame.DeuteriumTankID)
	} else if resDetails.Metal.StorageCapacity-resDetails.Metal.StorageCapacity/10 < planete.Resources.Metal {
		bot.BuildBuilding(planete.ID, ogame.MetalStorageID)
	}

}

func Researches(planete ogame.EmpireCelestial, bot *wrapper.OGame) {
	res, _ := bot.GetResearch()
	fmt.Println(res)
	id := planete.ID
	fac, _ := bot.GetFacilities(id)

	if fac.ResearchLab < 12 {
		bot.BuildBuilding(id, ogame.ResearchLabID)
	}

	bot.BuildTechnology(id, ogame.AstrophysicsID)
	fmt.Println("Recherche...")

	if res.ImpulseDrive < 4 {
		bot.BuildTechnology(id, ogame.ImpulseDriveID)
	}

	if res.EspionageTechnology < 5 {
		bot.BuildTechnology(id, ogame.EspionageTechnologyID)
	}

	//bot.BuildTechnology(id, ogame.ComputerTechnologyID)
	bot.BuildTechnology(id, ogame.IntergalacticResearchNetworkID)
	bot.BuildTechnology(id, ogame.CombustionDriveID)
	bot.BuildTechnology(id, ogame.ArmourTechnologyID)

	/*if res.EnergyTechnology < 12 {
		bot.BuildTechnology(id, ogame.EnergyTechnologyID)
	}*/

	/*if res.LaserTechnology < 10 {
		bot.BuildTechnology(id, ogame.LaserTechnologyID)
	}

	if res.IonTechnology < 5 {
		bot.BuildTechnology(id, ogame.IonTechnologyID)
	}

	if res.HyperspaceTechnology < 8 {
		bot.BuildTechnology(id, ogame.HyperspaceTechnologyID)
	}

	bot.BuildTechnology(id, ogame.PlasmaTechnologyID)
	bot.BuildTechnology(id, ogame.WeaponsTechnologyID)
	bot.BuildTechnology(id, ogame.ShieldingTechnology.ID)*/
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
		buildingID, buildingCountdown, researchID, researchCountdown, lfBuildingID,
			lfBuildingCountdown, lfResearchID, lfResearchCountdown := bot.ConstructionsBeingBuilt(planete.ID)

		fmt.Printf("buildingID = %s, buildingCountdown = %d, researchID = %s, researchCountdown = %d, lfBuildingID = %s, lfBuildingCountdown = %d, lfResearchID = %s, lfResearchCountdown = %d\n", buildingID, buildingCountdown, researchID, researchCountdown, lfBuildingID, lfBuildingCountdown, lfResearchID, lfResearchCountdown)
		buildResources(planete, bot)
		SetExpedition(planete.ID, planete.Coordinate, bot)
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
