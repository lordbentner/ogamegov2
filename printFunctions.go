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

func printStructFields(s interface{}) {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	// S'assurer qu'on travaille avec une struct
	if typ.Kind() != reflect.Struct {
		fmt.Println("Ce n'est pas une struct")
		return
	}

	nextLine := 0
	print_str := ""
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		if value.Interface() != "0" && value.Interface() != 0 {
			print_str += fmt.Sprintf("%s: %v, ", field.Name, value.Interface())
			nextLine++
		}

		if nextLine == 3 {
			print_str += "\n"
			nextLine = 0
		}
	}

	fmt.Println(print_str)
}

func printCurrentconstruction(id ogame.CelestialID, bot *wrapper.OGame) {
	buildingID, buildingCountdown, researchID, researchCountdown, lfBuildingID,
		lfBuildingCountdown, lfResearchID, lfResearchCountdown := bot.ConstructionsBeingBuilt(id)

	print_str := "Construction en cours  ==> "
	if buildingCountdown > 0 {
		print_str += fmt.Sprintf("buildingID = %s, buildingCountdown = %d ", buildingID, buildingCountdown)
	}
	if researchCountdown > 0 {
		print_str += fmt.Sprintf("researchID = %s, researchCountdown = %d ", researchID, researchCountdown)
	}
	if lfBuildingCountdown > 0 {
		print_str += fmt.Sprintf("lfBuildingID = %s, lfBuildingCountdown = %d ", lfBuildingID, lfBuildingCountdown)
	}
	if lfResearchCountdown > 0 {
		print_str += fmt.Sprintf("lfResearchID = %s, lfResearchCountdown = %d ", lfResearchID, lfResearchCountdown)
	}

	fmt.Println(print_str)
}
