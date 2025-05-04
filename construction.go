package main

import (
	"fmt"
	"math"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/wrapper"
)

func satProduction(planete ogame.EmpireCelestial) {
	satprod := ogame.SolarSatellite.Production(planete.Temperature, 1, true)
	cenprice := 20 * math.Pow(1.1, float64(planete.Supplies.SolarPlant))
	fmt.Printf("%s(%s) cout centrale de solaire : %f production solaire: %d\n", planete.Name, planete.Coordinate, cenprice, satprod)
	if cenprice > float64(satprod) {
		boot.BuildShips(planete.ID, ogame.SolarSatelliteID, 1)
	}
}

func getFastestResearch(planete ogame.EmpireCelestial) {
	min, _ := boot.TechnologyDetails(planete.ID, ogame.EspionageTechnologyID)
	min_duration := int64(min.ProductionDuration.Seconds()) * (min.Level + 1)
	fmt.Println(min_duration)
	for i := 107; i < 125; i++ {
		id := ogame.ID(i)
		if !id.IsValid() {
			continue
		}

		if id == ogame.LaserTechnologyID || id == ogame.IonTechnologyID || id == ogame.EnergyTechnologyID {
			continue
		}

		if id == ogame.EspionageTechnologyID && planete.Researches.EspionageTechnology > 7 {
			continue
		}

		tech, _ := boot.TechnologyDetails(planete.ID, id)
		tech_duration := int64(tech.ProductionDuration.Seconds()) * (tech.Level + 1)
		if min_duration > tech_duration {
			min = tech
			min_duration = tech_duration
		}
	}

	boot.BuildTechnology(planete.ID, min.TechnologyID)
}

func buildFormeVieHumans(planete ogame.EmpireCelestial) {
	boot.BuildBuilding(planete.ID, ogame.NeuroCalibrationCentreID)
	boot.BuildBuilding(planete.ID, ogame.MetropolisID)
	boot.BuildBuilding(planete.ID, ogame.FusionPoweredProductionID)
	boot.BuildBuilding(planete.ID, ogame.AcademyOfSciencesID)
	boot.BuildBuilding(planete.ID, ogame.SkyscraperID)
}

func buildFormeVieRocktas(planete ogame.EmpireCelestial) {
	if planete.Resources.Energy < 0 {
		boot.BuildBuilding(planete.ID, ogame.DisruptionChamberID)
	}

	boot.BuildBuilding(planete.ID, ogame.MegalithID)
	boot.BuildBuilding(planete.ID, ogame.OriktoriumID)
	boot.BuildBuilding(planete.ID, ogame.RuneForgeID)
	boot.BuildBuilding(planete.ID, ogame.MagmaForgeID)
	boot.BuildBuilding(planete.ID, ogame.CrystalRefineryID)
	boot.BuildBuilding(planete.ID, ogame.DeuteriumSynthesiserID)
}

func buildFormeVieKaelesh(planete ogame.EmpireCelestial) {
	/*	Sanctuary                  int64 // 14101 // Lifeform (kaelesh)
		AntimatterCondenser        int64 // 14102
		VortexChamber              int64 // 14103
		HallsOfRealisation         int64 // 14104
		ForumOfTranscendence       int64 // 14105
		AntimatterConvector        int64 // 14106
		CloningLaboratory          int64 // 14107
		ChrysalisAccelerator       int64 // 14108
		BioModifier                int64 // 14109
		PsionicModulator           int64 // 14110
		ShipManufacturingHall      int64 // 14111
		SupraRefractor             int64 // 14112*/

	boot.BuildBuilding(planete.ID, ogame.HallsOfRealisationID)
	boot.BuildBuilding(planete.ID, ogame.ForumOfTranscendenceID)
	if planete.LfBuildings.AntimatterConvector+3 < planete.LfBuildings.CloningLaboratory {
		boot.BuildBuilding(planete.ID, ogame.AntimatterConvectorID)
	} else {
		boot.BuildBuilding(planete.ID, ogame.CloningLaboratoryID)
	}

	boot.BuildBuilding(planete.ID, ogame.ChrysalisAcceleratorID)
	boot.BuildBuilding(planete.ID, ogame.PsionicModulatorID)
}

func buildFormeVie(planete ogame.EmpireCelestial) {
	buildFormeVieHumans(planete)
	buildFormeVieRocktas(planete)
	buildFormeVieKaelesh(planete)
	boot.BuildBuilding(planete.ID, ogame.RuneForgeID)
	if planete.LfBuildings.ResearchCentre < 5 || planete.LfBuildings.VortexChamber < 5 || planete.LfBuildings.RuneTechnologium < 5 {
		boot.BuildBuilding(planete.ID, ogame.ResearchCentreID)
		boot.BuildBuilding(planete.ID, ogame.RuneTechnologiumID)
		boot.BuildBuilding(planete.ID, ogame.VortexChamberID)
	} else {
		boot.BuildBuilding(planete.ID, ogame.HighEnergySmeltingID)
		boot.BuildBuilding(planete.ID, ogame.MagmaForgeID)
	}

	res := planete.LfBuildings.ResidentialSector < 41 && planete.LfBuildings.ResidentialSector > 0
	med := planete.LfBuildings.MeditationEnclave < 41 && planete.LfBuildings.MeditationEnclave > 0
	san := planete.LfBuildings.Sanctuary < 41 && planete.LfBuildings.Sanctuary > 0
	if res || med || san {
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
			boot.BuildBuilding(planete.ID, ogame.ResidentialSectorID)
			boot.BuildBuilding(planete.ID, ogame.MeditationEnclaveID)
			boot.BuildBuilding(planete.ID, ogame.SanctuaryID)
		}
	}

	bio := planete.LfBuildings.BiosphereFarm < 42 && planete.LfBuildings.BiosphereFarm > 0
	cry := planete.LfBuildings.CrystalFarm < 42 && planete.LfBuildings.CrystalFarm > 0
	ant := planete.LfBuildings.AntimatterCondenser < 42 && planete.LfBuildings.AntimatterCondenser > 0
	if bio || cry || ant {
		boot.BuildBuilding(planete.ID, ogame.BiosphereFarmID)
		boot.BuildBuilding(planete.ID, ogame.CrystalFarmID)
		boot.BuildBuilding(planete.ID, ogame.AntimatterCondenserID)
	}

	boot.BuildTechnology(planete.ID, resFastestLifeForm(planete, boot))
	boot.BuildTechnology(planete.ID, resFastestLifeFormKaelesh(planete, boot))
	boot.BuildTechnology(planete.ID, ogame.VolcanicBatteriesID)
	boot.BuildBuilding(planete.ID, ogame.CargoHoldExpansionCivilianShipsID)
	boot.BuildTechnology(planete.ID, ogame.HighEnergyPumpSystemsID)
}

func resFastestLifeForm(planete ogame.EmpireCelestial, bot *wrapper.OGame) ogame.ID {
	fast := ogame.AutomatedTransportLinesID
	a, _ := bot.TechnologyDetails(planete.ID, ogame.AutomatedTransportLinesID)
	h, _ := bot.TechnologyDetails(planete.ID, ogame.HighPerformanceExtractorsID)
	m, _ := bot.TechnologyDetails(planete.ID, ogame.MagmaPoweredProductionID)
	e, _ := bot.TechnologyDetails(planete.ID, ogame.EnhancedProductionTechnologiesID)
	s, _ := bot.TechnologyDetails(planete.ID, ogame.DepthSoundingID)
	p, _ := bot.TechnologyDetails(planete.ID, ogame.PsychoharmoniserID)
	t, _ := bot.TechnologyDetails(planete.ID, ogame.HardenedDiamondDrillHeadsID)
	list := []ogame.TechnologyDetails{a, h, m, e, s, p, t}
	min := int64(a.ProductionDuration.Seconds()) * (a.Level + 1)

	for _, elem := range list {
		basecost := int64(elem.ProductionDuration.Seconds()) * (elem.Level + 1)
		if basecost < min {
			min = basecost
			fast = elem.TechnologyID
		}
	}

	return fast
}

func resFastestLifeFormKaelesh(planete ogame.EmpireCelestial, bot *wrapper.OGame) ogame.ID {
	fast := ogame.EnhancedSensorTechnologyID
	a, _ := bot.TechnologyDetails(planete.ID, ogame.PsionicNetworkID)
	h, _ := bot.TechnologyDetails(planete.ID, ogame.EnhancedSensorTechnologyID)
	m, _ := bot.TechnologyDetails(planete.ID, ogame.TelekineticTractorBeamID)
	//6/443/9
	list := []ogame.TechnologyDetails{a, h, m}
	min := int64(a.ProductionDuration.Seconds()) * (a.Level + 1)

	for _, elem := range list {
		basecost := int64(elem.ProductionDuration.Seconds()) * (elem.Level + 1)
		if basecost < min {
			min = basecost
			fast = elem.TechnologyID
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

func buildResources(planete ogame.EmpireCelestial) {
	time.Sleep(10000)
	resDetails, _ := boot.GetResourcesDetails(planete.ID)
	fmt.Println("Detail des ressources Métal ====>")
	fmt.Printf("Available: %d , Storage cap.: %d Current production: %d\n", resDetails.Metal.Available, resDetails.Metal.StorageCapacity, resDetails.Metal.CurrentProduction)
	if planete.Facilities.RoboticsFactory < 10 {
		boot.BuildBuilding(planete.ID, ogame.RoboticsFactoryID)
	} else if planete.Facilities.NaniteFactory < 7 {
		boot.BuildBuilding(planete.ID, ogame.NaniteFactoryID)
	}

	if resDetails.Crystal.StorageCapacity-resDetails.Crystal.StorageCapacity/10 < planete.Resources.Crystal {
		boot.BuildBuilding(planete.ID, ogame.CrystalStorageID)
	} else if resDetails.Deuterium.StorageCapacity-resDetails.Deuterium.StorageCapacity/10 < planete.Resources.Deuterium {
		boot.BuildBuilding(planete.ID, ogame.DeuteriumTankID)
	} else if resDetails.Metal.StorageCapacity-resDetails.Metal.StorageCapacity/10 < planete.Resources.Metal {
		boot.BuildBuilding(planete.ID, ogame.MetalStorageID)
	}

	printStructFields(planete.Supplies)
	if resDetails.Metal.CurrentProduction > 0 || resDetails.Crystal.CurrentProduction > 0 {
		if planete.Resources.Energy < 0 {
			err := boot.BuildBuilding(planete.ID, ogame.SolarPlantID)
			if err != nil {
				satProduction(planete)
			} else {
				fmt.Println("construction solar plant")
			}
		} else if planete.Supplies.DeuteriumSynthesizer < int64(planete.Supplies.CrystalMine-3) {
			err := boot.BuildBuilding(planete.ID, 3)
			fmt.Printf("construction synthethiseur deuterium err =%s\n", err)
		} else if planete.Supplies.CrystalMine < int64(planete.Supplies.MetalMine-3) {
			err := boot.BuildBuilding(planete.ID, 2)
			fmt.Printf("construction crystal mine err =%s\n", err)
		} else {
			err := boot.BuildBuilding(planete.ID, 1)
			fmt.Printf("construction metal mine err =%s\n", err)
		}
	}

	if planete.Facilities.Shipyard < 12 {
		boot.BuildBuilding(planete.ID, ogame.ShipyardID)
	}
}

func Researches(planete ogame.EmpireCelestial, bot *wrapper.OGame, slots ogame.Slots) {
	res, _ := bot.GetResearch()
	id := planete.ID

	if res.EnergyTechnology < 12 {
		bot.BuildTechnology(id, ogame.EnergyTechnologyID)
	}

	if res.ImpulseDrive < 5 {
		bot.BuildTechnology(id, ogame.ImpulseDriveID)
	}

	if res.EspionageTechnology < 8 {
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

	getFastestResearch(planete)

	bot.BuildTechnology(id, ogame.PlasmaTechnologyID)
	bot.BuildTechnology(id, ogame.WeaponsTechnologyID)
	bot.BuildTechnology(id, ogame.ArmourTechnologyID)
	bot.BuildTechnology(id, ogame.AstrophysicsID)
	bot.BuildTechnology(id, ogame.IntergalacticResearchNetworkID)
}
