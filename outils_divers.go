package main

import (
	"fmt"

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

func getCorrectCoord(coord ogame.Coordinate) ogame.Coordinate {
	pos := coord.Position
	gal := coord.Galaxy
	sys := coord.System
	if pos > 15 {
		pos = 1
		sys++
	}

	if sys > 499 {
		sys = 1
		gal++
	}

	if gal > 6 {
		gal = 1
	}

	return ogame.Coordinate{Galaxy: gal, System: sys, Position: pos}
}
