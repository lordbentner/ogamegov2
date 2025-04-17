package main

import (
	"fmt"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/wrapper"
)

type Fleet struct {
	SmallCargo     int // Petit transporteur
	LargeCargo     int // Grand transporteur
	LightFighter   int // Chasseur léger
	EspionageProbe int // Sonde d'espionnage
	Pathfinder     int // Éclaireur
}

type Planet struct {
	Name           string
	Coordinates    string
	Fleet          Fleet
	MaxExpeditions int
}

type Expedition struct {
	Planet string
	Target string
	Fleet  Fleet
}

/*type Slots struct {
	InUse    int64
	Total    int64
	ExpInUse int64
	ExpTotal int64
}*/

func GetFleetsForCargo(bot *wrapper.OGame) (int, int, int) {
	GT := 0
	PT := 0
	Eclaireur := 0
	fleets, _ := bot.GetFleets()

	for _, fleet := range fleets {
		ships := fleet.Ships
		fmt.Println(fleet.Destination)
		printFleets(fleet)
		GT = GT + int(ships.LargeCargo)
		PT = PT + int(ships.SmallCargo)
		Eclaireur = Eclaireur + int(ships.Pathfinder)
	}

	return GT, PT, Eclaireur
}

func GetTotalForExpeShips(empire []ogame.EmpireCelestial) (int, int, int) {
	GT := 0
	PT := 0
	Eclaireur := 0

	for _, planete := range empire {
		ships := planete.Ships
		printShips(planete)
		GT = GT + int(ships.LargeCargo)
		PT = PT + int(ships.SmallCargo)
		Eclaireur = Eclaireur + int(ships.Pathfinder)
	}

	return GT, PT, Eclaireur
}

func printShips(planete ogame.EmpireCelestial) {
	fmt.Println(planete.Name + " : " + planete.Coordinate.String())
	fmt.Print(planete.Ships.LargeCargo)
	fmt.Print(" GT,  ")
	fmt.Print(planete.Ships.SmallCargo)
	fmt.Print(" PT, ")
	fmt.Print(planete.Ships.Pathfinder)
	fmt.Println(" Eclaireur")
}

func printShipsInfos(ships ogame.ShipsInfos) {
	fmt.Printf("GT = %d, PT = %d, Eclaireur = %d, Sonde = %d\n", ships.LargeCargo, ships.SmallCargo, ships.Pathfinder, ships.EspionageProbe)
}

func setExploVie(id ogame.CelestialID, coord ogame.Coordinate, bot *wrapper.OGame, index int) bool {
	att, _ := bot.IsUnderAttack()
	slots, _ := bot.GetSlots()
	slotsDispo := int(slots.Total - slots.InUse)
	if !att {
		for i := 0; i < slotsDispo; i++ {
			pos := int64(i + 1 + index)
			gal := coord.Galaxy
			sys := coord.System
			if pos > 15 {
				pos = 1
				gal++
				sys++
			}
			if gal > 6 {
				gal = 1
			}
			if sys > 499 {
				sys = 1
			}

			co := ogame.Coordinate{Galaxy: gal, System: sys, Position: pos}

			bot.SendDiscoveryFleet(id, co)
			fmt.Printf("fleet send to life discovery from %s to %s\n", coord.String(), co)
		}
	}

	time.Sleep(5000)
	newslots, _ := bot.GetSlots()
	slotsDispo = int(newslots.Total - newslots.InUse)
	if newslots.Total-newslots.InUse > 0 {
		return setExploVie(id, coord, bot, slotsDispo)
	}

	return false
}

func SetExpedition(id ogame.CelestialID, coord ogame.Coordinate, bot *wrapper.OGame) {
	sh, _ := bot.GetShips(id)
	slots, _ := bot.GetSlots()
	if slots.ExpInUse >= slots.ExpTotal || sh.EspionageProbe == 0 || sh.Pathfinder == 0 {
		return
	}

	var shipsInfos ogame.ShipsInfos
	if sh.EspionageProbe < 10 {
		shipsInfos.EspionageProbe = sh.EspionageProbe
	} else {
		shipsInfos.EspionageProbe = 10
	}

	shipsInfos.LargeCargo = sh.LargeCargo
	shipsInfos.SmallCargo = sh.SmallCargo
	shipsInfos.Pathfinder = sh.Pathfinder
	if sh.Destroyer > 0 {
		shipsInfos.Destroyer = 1
	} else if sh.Battlecruiser > 0 {
		shipsInfos.Battlecruiser = 1
	}

	co := ogame.Coordinate{Galaxy: coord.Galaxy, System: coord.System, Position: 16}
	bot.SendFleet(id, shipsInfos, 100, co, ogame.Expedition, ogame.Resources{}, 0, 0)
	fmt.Printf("fleet send to expedition from %s\n", coord.String())
	printShipsInfos(shipsInfos)
	time.Sleep(5000)
}

func getTotalGTCombine(empire []ogame.EmpireCelestial, empireMoon []ogame.EmpireCelestial, bot *wrapper.OGame) int {
	GT, PT, Eclaireur := GetTotalForExpeShips(empire)
	GTM, PTM, EclaireurM := GetTotalForExpeShips(empireMoon)
	GTF, PTF, ECLF := GetFleetsForCargo(bot)
	GT += GTM + GTF
	PT += PTM + PTF
	Eclaireur += EclaireurM + ECLF

	fmt.Printf("GT = %d PT = %d, Eclaireur = %d\n", GT, PT, Eclaireur)
	total := GT + (PT / 5) + Eclaireur/(5/2)
	fmt.Printf("GT Total = %d\n", total)
	return total
}

/*func testExpe() {
	total := getTotalGTCombine
}*/
