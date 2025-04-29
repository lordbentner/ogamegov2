package main

import (
	"fmt"
	"strings"
	"time"

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

func convertSecToTime(seconds int64) string {
	duration := time.Duration(seconds) * time.Second

	// Extraire les heures, minutes et secondes
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	secs := int(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d\n", hours, minutes, secs)
}

func getCargoGT(bot *wrapper.OGame) int64 {
	lfBonuses, _ := bot.GetCachedLfBonuses()
	multiplier := float64(bot.GetServerData().CargoHyperspaceTechMultiplier) / 100.0
	return ogame.LargeCargo.GetCargoCapacity(bot.GetCachedResearch(), lfBonuses, bot.CharacterClass(), multiplier, bot.GetServer().ProbeRaidsEnabled())
}

func getCargoPT(bot *wrapper.OGame) int64 {
	lfBonuses, _ := bot.GetCachedLfBonuses()
	multiplier := float64(bot.GetServerData().CargoHyperspaceTechMultiplier) / 100.0
	return ogame.SmallCargo.GetCargoCapacity(bot.GetCachedResearch(), lfBonuses, bot.CharacterClass(), multiplier, bot.GetServer().ProbeRaidsEnabled())
}

func getCargoPathFinder(bot *wrapper.OGame) int64 {
	lfBonuses, _ := bot.GetCachedLfBonuses()
	multiplier := float64(bot.GetServerData().CargoHyperspaceTechMultiplier) / 100.0
	return ogame.Pathfinder.GetCargoCapacity(bot.GetCachedResearch(), lfBonuses, bot.CharacterClass(), multiplier, bot.GetServer().ProbeRaidsEnabled())
}
