package main

import (
	"fmt"
	"math"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/wrapper"
)

func satProduction(planete ogame.EmpireCelestial, bot *wrapper.OGame) {
	satprod := ogame.SolarSatellite.Production(planete.Temperature, 1, true)
	cenprice := 20 * math.Pow(1.1, float64(planete.Supplies.SolarPlant))
	fmt.Printf("%s(%s) cout centrale de solaire : %f production solaire: %d\n", planete.Name, planete.Coordinate, cenprice, satprod)
	if cenprice > float64(satprod*2000) {
		bot.BuildShips(planete.ID, ogame.SolarSatelliteID, 1)
	}
}

func buildFormeVie(planete ogame.EmpireCelestial, bot *wrapper.OGame) {
	if planete.LfBuildings.ResidentialSector < 21 {
		bot.BuildBuilding(planete.ID, ogame.ResidentialSectorID)
	}

	if planete.LfBuildings.BiosphereFarm < 22 {
		bot.BuildBuilding(planete.ID, ogame.BiosphereFarmID)
	}
}

func buildResources(planete ogame.EmpireCelestial, bot *wrapper.OGame) {
	time.Sleep(10000)
	resDetails, _ := bot.GetResourcesDetails(planete.ID)
	if resDetails.Crystal.StorageCapacity-resDetails.Crystal.StorageCapacity/10 < planete.Resources.Crystal {
		bot.BuildBuilding(planete.ID, ogame.CrystalStorageID)
	} else if resDetails.Deuterium.StorageCapacity-resDetails.Deuterium.StorageCapacity/10 < planete.Resources.Deuterium {
		bot.BuildBuilding(planete.ID, ogame.DeuteriumTankID)
	} else if resDetails.Metal.StorageCapacity-resDetails.Metal.StorageCapacity/10 < planete.Resources.Metal {
		bot.BuildBuilding(planete.ID, ogame.MetalStorageID)
	}

	printStructFields(planete.Supplies)
	if planete.Facilities.RoboticsFactory < 10 {
		bot.BuildBuilding(planete.ID, ogame.RoboticsFactoryID)
	} else {
		bot.BuildBuilding(planete.ID, ogame.NaniteFactoryID)
	}

	if planete.Resources.Energy < 0 {
		err := bot.BuildBuilding(planete.ID, ogame.SolarPlantID)
		if err != nil {
			satProduction(planete, bot)
		} else {
			fmt.Println("construction solar plant")
		}
	} else if planete.Supplies.DeuteriumSynthesizer < int64(planete.Supplies.CrystalMine-3) {
		err := bot.BuildBuilding(planete.ID, 3)
		fmt.Printf("construction synthethiseur deuterium err =%s\n", err)
	} else if planete.Supplies.CrystalMine < int64(planete.Supplies.MetalMine-3) {
		err := bot.BuildBuilding(planete.ID, 2)
		fmt.Printf("construction crystal mine err =%s\n", err)
	} else {
		err := bot.BuildBuilding(planete.ID, 1)
		fmt.Printf("construction metal mine err =%s\n", err)
	}

	if planete.Facilities.Shipyard < 12 {
		bot.BuildBuilding(planete.ID, ogame.ShipyardID)
	}
}

func Researches(planete ogame.EmpireCelestial, bot *wrapper.OGame, slots ogame.Slots) {
	res, _ := bot.GetResearch()
	id := planete.ID
	fac, _ := bot.GetFacilities(id)

	if fac.ResearchLab < 12 {
		bot.BuildBuilding(id, ogame.ResearchLabID)
	}

	bot.BuildTechnology(id, ogame.AstrophysicsID)

	if slots.Total-slots.ExpTotal < 1 || res.ComputerTechnology < 10 {
		bot.BuildTechnology(id, ogame.ComputerTechnologyID)
	}

	if res.ImpulseDrive < 5 {
		bot.BuildTechnology(id, ogame.ImpulseDriveID)
	}

	if res.EspionageTechnology < 5 {
		bot.BuildTechnology(id, ogame.EspionageTechnologyID)
	}

	bot.BuildTechnology(id, ogame.IntergalacticResearchNetworkID)
	bot.BuildTechnology(id, ogame.CombustionDriveID)
	if res.ShieldingTechnology < 6 {
		bot.BuildTechnology(id, ogame.ShieldingTechnologyID)
	}

	if res.HyperspaceTechnology < 8 {
		bot.BuildTechnology(id, ogame.HyperspaceTechnologyID)
	}

	if res.EnergyTechnology < 5 {
		bot.BuildTechnology(id, ogame.EnergyTechnologyID)
	}

	if slots.Total-slots.ExpTotal < 1 || res.ComputerTechnology < 10 {
		bot.BuildTechnology(id, ogame.ComputerTechnologyID)
	}

	if res.LaserTechnology < 10 {
		bot.BuildTechnology(id, ogame.LaserTechnologyID)
	}

	if res.IonTechnology < 5 {
		bot.BuildTechnology(id, ogame.IonTechnologyID)
	}

	bot.BuildTechnology(id, ogame.PlasmaTechnologyID)
	bot.BuildTechnology(id, ogame.WeaponsTechnologyID)
	bot.BuildTechnology(id, ogame.ArmourTechnologyID)
}
