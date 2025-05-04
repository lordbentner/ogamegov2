package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/wrapper"
)

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

	if strings.Contains(content, "avons fêté") || strings.Contains(content, "Si cela continue comme ca") {
		return true
	}

	if strings.Contains(content, "Il serait peut être plus judicieux") {
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
