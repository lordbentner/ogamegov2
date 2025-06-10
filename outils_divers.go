package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
)

func sendTelegramMessage(token, chatID, message string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id": {chatID},
		"text":    {message},
	})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("✅ Message envoyé avec succès.")
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

func getCargoGT() int64 {
	lfBonuses, _ := boot.GetCachedLfBonuses()
	multiplier := float64(boot.GetServerData().CargoHyperspaceTechMultiplier) / 100.0
	return ogame.LargeCargo.GetCargoCapacity(boot.GetCachedResearch(), lfBonuses, boot.CharacterClass(), multiplier, boot.GetServer().ProbeRaidsEnabled())
}

func getCargoPT() int64 {
	lfBonuses, _ := boot.GetCachedLfBonuses()
	multiplier := float64(boot.GetServerData().CargoHyperspaceTechMultiplier) / 100.0
	return ogame.SmallCargo.GetCargoCapacity(boot.GetCachedResearch(), lfBonuses, boot.CharacterClass(), multiplier, boot.GetServer().ProbeRaidsEnabled())
}

func getCargoPathFinder() int64 {
	lfBonuses, _ := boot.GetCachedLfBonuses()
	multiplier := float64(boot.GetServerData().CargoHyperspaceTechMultiplier) / 100.0
	return ogame.Pathfinder.GetCargoCapacity(boot.GetCachedResearch(), lfBonuses, boot.CharacterClass(), multiplier, boot.GetServer().ProbeRaidsEnabled())
}

func readJSONCoordFdV() ogame.Coordinate {
	file, err := os.Open("C:\\Users\\Utilisateur\\Documents\\ogameBot\\data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var coord ogame.Coordinate
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&coord)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Nom: %s\n", coord)
	return coord
}
