package main

import (
	"fmt"
	"strings"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/wrapper"
)

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

func changeSystemeExploration(content string) bool {

	if strings.Contains(content, "avons découvert") || strings.Contains(content, "avons failli") {
		return true
	}

	if strings.Contains(content, "avons fêté") {
		return true
	}

	return false
}
