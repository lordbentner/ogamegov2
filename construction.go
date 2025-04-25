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

/*func getFastestResearch(planete ogame.EmpireCelestial, bot *wrapper.OGame) {
	speed := 8
	hasTechnocrat := false
	lfbonuses, _ := bot.GetLfBonuses()
	d_astro := ogame.Astrophysics.ConstructionTime(planete.Researches.Astrophysics+1, int64(speed), planete.Facilities, lfbonuses, ogame.Discoverer, hasTechnocrat)
	d_hyper_techno := ogame.HyperspaceTechnology.ConstructionTime(planete.Researches.HyperspaceTechnology+1, int64(speed), planete.Facilities, lfbonuses, ogame.Discoverer, hasTechnocrat)
	d_hyper_drive := ogame.HyperspaceTechnology.ConstructionTime(planete.Researches.HyperspaceDrive+1, int64(speed), planete.Facilities, lfbonuses, ogame.Discoverer, hasTechnocrat)
	d_weapon := ogame.WeaponsTechnology.ConstructionTime(planete.Researches.WeaponsTechnology+1, int64(speed), planete.Facilities, lfbonuses, ogame.Discoverer, hasTechnocrat)
	d_shield := ogame.ShieldingTechnology.ConstructionTime(planete.Researches.ShieldingTechnology+1, int64(speed), planete.Facilities, lfbonuses, ogame.Discoverer, hasTechnocrat)

	fmt.Printf("a = %d, ht = %d, hg = %d, wepon = %d", d_astro, d_hyper_techno, d_hyper_drive, d_weapon)
	//(level int64, universeSpeed int64, facilities ogame.BuildAccelerators, lfBonuses ogame.LfBonuses, class ogame.CharacterClass, hasTechnocrat bool)

}*/

func buildFormeVie(planete ogame.EmpireCelestial, bot *wrapper.OGame) {
	//	bot.BuildBuilding(planete.ID, ogame.FusionPoweredProductionID)
	bot.BuildBuilding(planete.ID, ogame.CrystalRefineryID)
	if planete.LfBuildings.ResearchCentre < 5 {
		bot.BuildBuilding(planete.ID, ogame.ResearchCentreID)
		bot.BuildBuilding(planete.ID, ogame.RuneTechnologiumID)
	} else {
		bot.BuildBuilding(planete.ID, ogame.HighEnergySmeltingID)
		bot.BuildBuilding(planete.ID, ogame.MagmaForgeID)
	}

	if planete.LfBuildings.ResidentialSector < 41 || planete.LfBuildings.MeditationEnclave < 41 {
		hab := planete.LfBuildings.ResidentialSector
		if planete.LfBuildings.MeditationEnclave > 0 {
			hab = planete.LfBuildings.MeditationEnclave
		}

		food := planete.LfBuildings.BiosphereFarm
		if planete.LfBuildings.CrystalFarm > 0 {
			food = planete.LfBuildings.CrystalFarm
		}
		if hab < food-1 {
			bot.BuildBuilding(planete.ID, ogame.ResidentialSectorID)
			bot.BuildBuilding(planete.ID, ogame.MeditationEnclaveID)
		}
	}

	if planete.LfBuildings.BiosphereFarm < 42 || planete.LfBuildings.CrystalFarm < 42 {
		bot.BuildBuilding(planete.ID, ogame.BiosphereFarmID)
		bot.BuildBuilding(planete.ID, ogame.CrystalFarmID)
	}

	bot.BuildTechnology(planete.ID, ogame.HighPerformanceExtractorsID)
	bot.BuildTechnology(planete.ID, ogame.VolcanicBatteriesID)
	bot.BuildTechnology(planete.ID, ogame.HighEnergyPumpSystemsID)
	bot.BuildTechnology(planete.ID, ogame.MagmaPoweredProductionID)
	bot.BuildTechnology(planete.ID, ogame.AutomatedTransportLinesID)
}

func buildMoon(moon ogame.EmpireCelestial, bot *wrapper.OGame) {
	if moon.Facilities.JumpGate > 0 {
		return
	}
	if moon.Fields.Built == moon.Fields.Total-2 && moon.Facilities.RoboticsFactory > 8 {
		bot.BuildBuilding(moon.ID, ogame.JumpGateID)
	} else if moon.Fields.Built == moon.Fields.Total-1 {
		bot.BuildBuilding(moon.ID, ogame.LunarBaseID)
		fmt.Println("Construction base lunaire")
	} else {
		fmt.Println("Construction usine de robot")
		bot.BuildBuilding(moon.ID, ogame.RoboticsFactoryID)
	}
}

func buildResources(planete ogame.EmpireCelestial, bot *wrapper.OGame) {
	time.Sleep(10000)
	resDetails, _ := bot.GetResourcesDetails(planete.ID)
	fmt.Println("Detail des ressources Métal ====>")
	fmt.Printf("Available: %d , Storage cap.: %d Current production: %d\n", resDetails.Metal.Available, resDetails.Metal.StorageCapacity, resDetails.Metal.CurrentProduction)
	if planete.Facilities.RoboticsFactory < 12 {
		bot.BuildBuilding(planete.ID, ogame.RoboticsFactoryID)
	} else if planete.Facilities.NaniteFactory < 7 {
		bot.BuildBuilding(planete.ID, ogame.NaniteFactoryID)
	}

	if resDetails.Crystal.StorageCapacity-resDetails.Crystal.StorageCapacity/10 < planete.Resources.Crystal {
		bot.BuildBuilding(planete.ID, ogame.CrystalStorageID)
	} else if resDetails.Deuterium.StorageCapacity-resDetails.Deuterium.StorageCapacity/10 < planete.Resources.Deuterium {
		bot.BuildBuilding(planete.ID, ogame.DeuteriumTankID)
	} else if resDetails.Metal.StorageCapacity-resDetails.Metal.StorageCapacity/10 < planete.Resources.Metal {
		bot.BuildBuilding(planete.ID, ogame.MetalStorageID)
	}

	printStructFields(planete.Supplies)
	if resDetails.Metal.CurrentProduction > 0 || resDetails.Crystal.CurrentProduction > 0 {
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
	}

	if planete.Facilities.Shipyard < 12 {
		bot.BuildBuilding(planete.ID, ogame.ShipyardID)
	}
}

func Researches(planete ogame.EmpireCelestial, bot *wrapper.OGame, slots ogame.Slots) {
	res, _ := bot.GetResearch()
	id := planete.ID

	if slots.Total-slots.ExpTotal < 1 || res.ComputerTechnology < 10 {
		bot.BuildTechnology(id, ogame.ComputerTechnologyID)
	}

	if res.EnergyTechnology < 12 {
		bot.BuildTechnology(id, ogame.EnergyTechnologyID)
	}

	bot.BuildTechnology(id, ogame.AstrophysicsID)

	if res.ImpulseDrive < 5 {
		bot.BuildTechnology(id, ogame.ImpulseDriveID)
	}

	if res.EspionageTechnology < 5 {
		bot.BuildTechnology(id, ogame.EspionageTechnologyID)
	}

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
	bot.BuildTechnology(id, ogame.IntergalacticResearchNetworkID)
}
