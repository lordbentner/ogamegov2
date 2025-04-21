package main

import (
	"fmt"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/wrapper"
)

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

func setExploVie(id ogame.CelestialID, coord ogame.Coordinate, bot *wrapper.OGame, index int) int {
	slots, _ := bot.GetSlots()
	if slots.InUse >= slots.ExpTotal {
		return 0
	}
	att, _ := bot.IsUnderAttack()
	slotsDispo := int(slots.Total - slots.InUse)
	if !att {
		fmt.Println("lklkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk")
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
	return int(newslots.Total - newslots.InUse)
}

func getFleetCompositionForExplo(sh ogame.ShipsInfos, slotDispo int64) ogame.ShipsInfos {
	var shipsInfos ogame.ShipsInfos
	if sh.EspionageProbe < 10 {
		shipsInfos.EspionageProbe = sh.EspionageProbe
	} else {
		shipsInfos.EspionageProbe = 10
	}

	shipsInfos.LargeCargo = sh.LargeCargo / slotDispo
	shipsInfos.SmallCargo = sh.SmallCargo / slotDispo
	shipsInfos.Pathfinder = sh.Pathfinder / slotDispo
	if slotDispo == 1 {
		shipsInfos.LargeCargo = sh.LargeCargo
		shipsInfos.SmallCargo = sh.SmallCargo
		shipsInfos.Pathfinder = sh.Pathfinder
	}
	shipsInfos.Pathfinder = sh.Pathfinder
	if sh.Destroyer > 0 {
		shipsInfos.Destroyer = 1
	} else if sh.Battlecruiser > 0 {
		shipsInfos.Battlecruiser = 1
	} else if sh.Battleship > 0 {
		shipsInfos.Battleship = 1
	} else if sh.Pathfinder < 1 {
		if sh.Cruiser > 0 {
			shipsInfos.Cruiser = 1
		} else if sh.HeavyFighter > 0 {
			shipsInfos.HeavyFighter = 1
		}
	}

	if shipsInfos.SmallCargo == 0 && sh.SmallCargo > 0 {
		shipsInfos.SmallCargo++
	}

	if shipsInfos.LargeCargo == 0 && sh.LargeCargo > 0 {
		shipsInfos.LargeCargo++
	}

	return shipsInfos
}

func SetExpedition(id ogame.CelestialID, coord ogame.Coordinate, bot *wrapper.OGame) {
	sh, _ := bot.GetShips(id)
	slots, _ := bot.GetSlots()
	if slots.ExpInUse >= slots.ExpTotal || sh.SmallCargo == 0 /*|| sh.EspionageProbe == 0 || sh.Pathfinder == 0*/ || slots.InUse >= slots.Total {
		return
	}

	slotDispo := slots.ExpTotal - slots.ExpInUse
	shipsInfos := getFleetCompositionForExplo(sh, slotDispo)

	co := ogame.Coordinate{Galaxy: coord.Galaxy, System: coord.System + 3, Position: 16}
	bot.SendFleet(id, shipsInfos, 100, co, ogame.Expedition, ogame.Resources{}, 0, 0)
	fmt.Printf("fleet send to expedition from %s with this fleet: ", coord.String())
	printStructFields(shipsInfos)
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
