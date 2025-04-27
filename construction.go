package main

import (
	"fmt"
	"math"
	"sort"
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

func getFastestResearch(planete ogame.EmpireCelestial, bot *wrapper.OGame) {
	level := int64(planete.Researches.Astrophysics) + 1
	universeSpeed := int64(8)
	hasTechnocrat := false
	lFBonuses, _ := bot.GetLfBonuses()
	fac := planete.Facilities
	empire, _ := bot.GetEmpire(ogame.PlanetType)

	fac.ResearchLab = 0
	sort.Slice(empire, func(i int, j int) bool {
		return empire[i].Facilities.ResearchLab > empire[j].Facilities.ResearchLab
	})
	for index, pl := range empire {
		if index-1 < int(planete.Researches.IntergalacticResearchNetwork) {
			fac.ResearchLab += pl.Facilities.ResearchLab
		}
	}

	id_tech := ogame.EspionageTechnologyID
	duration := ogame.Objs.ByID(ogame.EspionageTechnologyID).ConstructionTime(level, universeSpeed, fac, lFBonuses, bot.CharacterClass(), hasTechnocrat)
	for i := 107; i < 125; i++ {
		id := ogame.ID(i)
		if !id.IsValid() {
			continue
		}

		level = int64(planete.Researches.Astrophysics) + 1
		/*if id == ogame.LaserTechnologyID || id == ogame.IonTechnologyID || id == ogame.EnergyTechnologyID {
			continue
		}*/

		if id == ogame.EspionageTechnologyID && planete.Researches.EspionageTechnology > 7 {
			continue
		}

		temp_duration := ogame.Objs.ByID(id).ConstructionTime(level, universeSpeed, fac, lFBonuses, bot.CharacterClass(), hasTechnocrat)
		if duration > temp_duration {
			duration = temp_duration
			id_tech = id
		}
	}

	fmt.Println(id_tech)
	//bot.BuildTechnology(planete.ID, id_tech)
	//os.Exit(0)
}

func buildFormeVie(planete ogame.EmpireCelestial, bot *wrapper.OGame) {
	bot.BuildBuilding(planete.ID, ogame.CrystalRefineryID)
	bot.BuildBuilding(planete.ID, ogame.FusionPoweredProductionID)
	bot.BuildBuilding(planete.ID, ogame.AcademyOfSciencesID)
	bot.BuildBuilding(planete.ID, ogame.RuneForgeID)
	if planete.LfBuildings.ResearchCentre < 5 || planete.LfBuildings.VortexChamber < 5 || planete.LfBuildings.RuneTechnologium < 5 {
		bot.BuildBuilding(planete.ID, ogame.ResearchCentreID)
		bot.BuildBuilding(planete.ID, ogame.RuneTechnologiumID)
		bot.BuildBuilding(planete.ID, ogame.VortexChamberID)
	} else {
		bot.BuildBuilding(planete.ID, ogame.HighEnergySmeltingID)
		bot.BuildBuilding(planete.ID, ogame.MagmaForgeID)
	}

	if planete.LfBuildings.ResidentialSector < 41 || planete.LfBuildings.MeditationEnclave < 41 || planete.LfBuildings.Sanctuary < 41 {
		hab := planete.LfBuildings.ResidentialSector
		if planete.LfBuildings.MeditationEnclave > 0 {
			hab = planete.LfBuildings.MeditationEnclave
		} else if planete.LfBuildings.Sanctuary > 0 {
			hab = planete.LfBuildings.Sanctuary
		}

		food := planete.LfBuildings.BiosphereFarm
		if planete.LfBuildings.CrystalFarm > 0 {
			food = planete.LfBuildings.CrystalFarm
		} else if planete.LfBuildings.AntimatterCondenser > 0 {
			food = planete.LfBuildings.AntimatterCondenser
		}
		if hab < food-1 {
			bot.BuildBuilding(planete.ID, ogame.ResidentialSectorID)
			bot.BuildBuilding(planete.ID, ogame.MeditationEnclaveID)
			bot.BuildBuilding(planete.ID, ogame.SanctuaryID)
		}
	}

	if planete.LfBuildings.BiosphereFarm < 42 || planete.LfBuildings.CrystalFarm < 42 || planete.LfBuildings.AntimatterCondenser < 42 {
		bot.BuildBuilding(planete.ID, ogame.BiosphereFarmID)
		bot.BuildBuilding(planete.ID, ogame.CrystalFarmID)
		bot.BuildBuilding(planete.ID, ogame.AntimatterCondenserID)
	}

	bot.BuildTechnology(planete.ID, resFastestLifeForm(planete, bot))
	bot.BuildTechnology(planete.ID, resFastestLifeFormKaelesh(planete, bot))
	bot.BuildTechnology(planete.ID, ogame.VolcanicBatteriesID)
	bot.BuildBuilding(planete.ID, ogame.CargoHoldExpansionCivilianShipsID)
	bot.BuildTechnology(planete.ID, ogame.HighEnergyPumpSystemsID)
	ff, _ := bot.TechnologyDetails(planete.ID, ogame.AutomatedTransportLinesID)
	fmt.Println(ff.ProductionDuration)
}

func resFastestLifeForm(planete ogame.EmpireCelestial, bot *wrapper.OGame) ogame.ID {
	fast := ogame.AutomatedTransportLinesID
	a, _ := bot.TechnologyDetails(planete.ID, ogame.AutomatedTransportLinesID)
	h, _ := bot.TechnologyDetails(planete.ID, ogame.HighPerformanceExtractorsID)
	m, _ := bot.TechnologyDetails(planete.ID, ogame.MagmaPoweredProductionID)
	e, _ := bot.TechnologyDetails(planete.ID, ogame.EnhancedProductionTechnologiesID)
	list := []ogame.TechnologyDetails{a, h, m, e}
	min := a.ProductionDuration

	for _, elem := range list {
		if elem.ProductionDuration > min {
			min = elem.ProductionDuration
			fast = elem.TechnologyID
		}
	}

	/*if a.ProductionDuration > h.ProductionDuration || a.ProductionDuration > m.ProductionDuration || a.ProductionDuration > e.ProductionDuration {
		if h.ProductionDuration < m.ProductionDuration {
			fast = ogame.HighPerformanceExtractorsID
		} else {
			fast = ogame.MagmaPoweredProductionID
		}
	}*/

	return fast
}

func resFastestLifeFormKaelesh(planete ogame.EmpireCelestial, bot *wrapper.OGame) ogame.ID {
	fast := ogame.EnhancedSensorTechnologyID
	a, _ := bot.TechnologyDetails(planete.ID, ogame.PsionicNetworkID)
	h, _ := bot.TechnologyDetails(planete.ID, ogame.EnhancedSensorTechnologyID)
	m, _ := bot.TechnologyDetails(planete.ID, ogame.TelekineticTractorBeamID)
	if a.ProductionDuration > h.ProductionDuration || a.ProductionDuration > m.ProductionDuration {
		if h.ProductionDuration < m.ProductionDuration {
			fast = ogame.EnhancedSensorTechnologyID
		} else {
			fast = ogame.TelekineticTractorBeamID
		}
	}

	return fast
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
	if planete.Facilities.RoboticsFactory < 10 {
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

	//getFastestResearch(planete, bot)

	if res.EnergyTechnology < 12 {
		bot.BuildTechnology(id, ogame.EnergyTechnologyID)
	}

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
	bot.BuildTechnology(id, ogame.AstrophysicsID)
	bot.BuildTechnology(id, ogame.IntergalacticResearchNetworkID)
}
