package main

import (
	"fmt"
	"reflect"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/wrapper"
)

func printFleets(fleet ogame.Fleet) {
	fmt.Print(fleet.Ships.LargeCargo)
	fmt.Print(" GT,  ")
	fmt.Print(fleet.Ships.SmallCargo)
	fmt.Print(" PT, ")
	fmt.Print(fleet.Ships.Pathfinder)
	fmt.Println(" Eclaireur")
}

func printStructFields(s interface{}) string {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	// S'assurer qu'on travaille avec une struct
	if typ.Kind() != reflect.Struct {
		fmt.Println("Ce n'est pas une struct")
		return ""
	}

	nextLine := 0
	print_str := ""
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		if value.Kind() == reflect.Int || value.Kind() == reflect.Int64 {
			if int(value.Int()) > 0 {
				print_str += fmt.Sprintf("%s: %v, ", field.Name, value.Interface())
				nextLine++
			}
		} else if value.Interface() != "0" {
			print_str += fmt.Sprintf("%s: %v, ", field.Name, value.Interface())
			nextLine++
		}

		if nextLine == 3 {
			print_str += "\n"
			nextLine = 0
		}
	}

	fmt.Println(print_str)
	return print_str
}

func printCurrentconstruction(id ogame.CelestialID, bot *wrapper.OGame) string {
	buildingID, buildingCountdown, researchID, researchCountdown, lfBuildingID,
		lfBuildingCountdown, lfResearchID, lfResearchCountdown := bot.ConstructionsBeingBuilt(id)

	print_str := "Construction en cours  ==> "
	if buildingCountdown > 0 {
		print_str += fmt.Sprintf("buildingID = %s, buildingCountdown = %s ", buildingID, convertSecToTime(buildingCountdown))
	}
	if researchCountdown > 0 {
		print_str += fmt.Sprintf("researchID = %s, researchCountdown = %s ", researchID, convertSecToTime(researchCountdown))
	}
	if lfBuildingCountdown > 0 {
		print_str += fmt.Sprintf("lfBuildingID = %s, lfBuildingCountdown = %s ", lfBuildingID, convertSecToTime(lfBuildingCountdown))
	}
	if lfResearchCountdown > 0 {
		print_str += fmt.Sprintf("lfResearchID = %s, lfResearchCountdown = %s ", lfResearchID, convertSecToTime(lfResearchCountdown))
	}

	fmt.Println(print_str)
	return print_str
}

func printShipsInfos(ships ogame.ShipsInfos) {
	fmt.Printf("GT = %d, PT = %d, Eclaireur = %d, Sonde = %d\n", ships.LargeCargo, ships.SmallCargo, ships.Pathfinder, ships.EspionageProbe)
}
