package main

import (
	"fmt"
	"sort"

	"github.com/alaingilbert/ogame/pkg/ogame"
)

func getFastestResearch(planete ogame.EmpireCelestial) {
	min, _ := boot.TechnologyDetails(planete.ID, ogame.EspionageTechnologyID)
	min_duration := int64(min.ProductionDuration.Seconds()) * (min.Level + 1)
	fmt.Println(min_duration)
	var list []ogame.TechnologyDetails
	for i := 107; i < 125; i++ {
		id := ogame.ID(i)
		if !id.IsValid() {
			continue
		}

		if id == ogame.LaserTechnologyID || id == ogame.IonTechnologyID || id == ogame.EnergyTechnologyID || id == ogame.ImpulseDriveID {
			continue
		}

		if id == ogame.EspionageTechnologyID && planete.Researches.EspionageTechnology > 7 {
			continue
		}

		tech, _ := boot.TechnologyDetails(planete.ID, id)
		list = append(list, tech)
	}

	sort.Slice(list, func(i int, j int) bool {
		pricecost_i := list[i].Price.Metal + list[i].Price.Crystal + list[i].Price.Deuterium
		pricecost_j := list[j].Price.Metal + list[j].Price.Crystal + list[j].Price.Deuterium
		basecost := pricecost_i * (list[i].Level + 1)
		j_basecost := pricecost_j * (list[j].Level + 1)
		return basecost < j_basecost
	})

	for _, elem := range list {
		if elem.ProductionDuration.Seconds() > 0 || elem.ProductionDuration.Minutes() > 0 || elem.ProductionDuration.Hours() > 0 {
			fmt.Println(elem)
			boot.BuildTechnology(planete.ID, elem.TechnologyID)
		}
	}
}

func Researches(planete ogame.EmpireCelestial, slots ogame.Slots) {
	res, _ := boot.GetResearch()
	id := planete.ID

	if res.ComputerTechnology < 10 {
		boot.BuildTechnology(id, ogame.ComputerTechnologyID)
	}

	if res.Astrophysics == 0 && res.EspionageTechnology < 4 && res.ImpulseDrive < 3 {
		if res.EspionageTechnology < 4 {
			boot.BuildTechnology(id, ogame.EspionageTechnologyID)
		}

		if res.ImpulseDrive < 3 {
			boot.BuildTechnology(id, ogame.ImpulseDriveID)
		}

		return
	}

	boot.BuildTechnology(id, ogame.AstrophysicsID)
	if res.ImpulseDrive < 3 {
		boot.BuildTechnology(id, ogame.ImpulseDriveID)
	}

	if res.EnergyTechnology < 12 {
		boot.BuildTechnology(id, ogame.EnergyTechnologyID)
	}

	if res.EspionageTechnology < 8 {
		boot.BuildTechnology(id, ogame.EspionageTechnologyID)
	}

	if res.LaserTechnology < 10 {
		boot.BuildTechnology(id, ogame.LaserTechnologyID)
	}

	if res.IonTechnology < 5 {
		boot.BuildTechnology(id, ogame.IonTechnologyID)
	}

	getFastestResearch(planete)
}

func resFastestLifeForm(planete ogame.EmpireCelestial) {
	a, _ := boot.TechnologyDetails(planete.ID, ogame.AutomatedTransportLinesID)
	h, _ := boot.TechnologyDetails(planete.ID, ogame.HighPerformanceExtractorsID)
	m, _ := boot.TechnologyDetails(planete.ID, ogame.MagmaPoweredProductionID)
	e, _ := boot.TechnologyDetails(planete.ID, ogame.EnhancedProductionTechnologiesID)
	s, _ := boot.TechnologyDetails(planete.ID, ogame.DepthSoundingID)
	p, _ := boot.TechnologyDetails(planete.ID, ogame.PsychoharmoniserID)
	t, _ := boot.TechnologyDetails(planete.ID, ogame.HardenedDiamondDrillHeadsID)
	i, _ := boot.TechnologyDetails(planete.ID, ogame.ArtificialSwarmIntelligenceID)
	pl, _ := boot.TechnologyDetails(planete.ID, ogame.ImprovedStellaratorID)
	list := []ogame.TechnologyDetails{a, h, m, e, s, p, t, i, pl}

	sort.Slice(list, func(i int, j int) bool {
		basecost := int64(list[i].ProductionDuration.Seconds()) * (list[i].Level + 1)
		j_basecost := int64(list[j].ProductionDuration.Seconds()) * (list[j].Level + 1)
		return basecost < j_basecost
	})

	for _, elem := range list {
		boot.BuildTechnology(planete.ID, elem.TechnologyID)
	}
}

func resFastestLifeFormKaelesh(planete ogame.EmpireCelestial) {
	a, _ := boot.TechnologyDetails(planete.ID, ogame.PsionicNetworkID)
	h, _ := boot.TechnologyDetails(planete.ID, ogame.EnhancedSensorTechnologyID)
	m, _ := boot.TechnologyDetails(planete.ID, ogame.TelekineticTractorBeamID)
	list := []ogame.TechnologyDetails{a, h, m}
	sort.Slice(list, func(i int, j int) bool {
		basecost := int64(list[i].ProductionDuration.Seconds()) * (list[i].Level + 1)
		j_basecost := int64(list[j].ProductionDuration.Seconds()) * (list[j].Level + 1)
		return basecost < j_basecost
	})

	for _, elem := range list {
		boot.BuildTechnology(planete.ID, elem.TechnologyID)
	}
}
