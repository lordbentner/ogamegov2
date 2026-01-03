package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/wrapper"
)

var boot *wrapper.OGame

func getFlottePourExpe() {
	att, _ := boot.IsUnderAttack()
	if att {
		gestionAttack()
	}

	/*getMaxExpeDebris(4)
	getMaxExpeDebris(5)
	os.Exit(0)*/

	CargoExpeInsuffisant := 0

	fleets, slots, err := boot.GetFleets()
	if err != nil {
		fmt.Println("1 Errors impossible d'obtenir les données sur les flottes: %s", err.Error())
		return
	}

	fmt.Println("=====================Flottes=======================")
	fmt.Printf("%s slots: ", time.Now().Format(time.RFC850))
	fmt.Println(slots)
	for i, fleet := range fleets {
		fmt.Printf("flotte %d ==> ", i)
		printStructFields(fleet.Ships)
	}
	fmt.Println("====================================================")
	empire, _ := boot.GetEmpire(ogame.PlanetType)
	if len(empire) == 0 {
		fmt.Println(empire)
		return
	}

	first_planet := empire[0]

	empireMoon, _ := boot.GetEmpire(ogame.MoonType)
	empire = append(empire, empireMoon...)
	empire = sliceEmpireCargo(empire)
	expeMes := gestionMessagesExpe()
	coordExpe := expeMes.Coordinate
	fmt.Println(coordExpe)
	if changeSystemeExploration(expeMes.Content) {
		coordMain := empire[0].Coordinate
		sys := coordExpe.System + 1
		if sys > coordMain.System+10 {
			sys = coordMain.System
		}

		coordExpe = ogame.Coordinate{Galaxy: coordExpe.Galaxy, System: sys, Position: 16}
	}

	fmt.Println("================================================================")
	//validCoordLF = readJSONCoordFdV()
	//fmt.Println(validCoordLF)
	//os.Exit(0)
	//HasMoonRes := false
	for i, planete := range empire {
		fmt.Printf("======================= planete %s(%s) =========================\n", planete.Name, planete.Coordinate)

		if planete.Type == ogame.MoonType {
			buildMoon(planete)
			if slots.ExpInUse >= slots.ExpTotal && slots.InUse < slots.Total {
				sendFleetToMoon(planete)
				if sendFleetFromMoonToPlanet(planete) {
					//HasMoonRes = true
				}
			}
		} else if planete.Fields.Built < planete.Fields.Total-2 {
			buildResources(planete)
		} else {
			boot.BuildBuilding(planete.ID, ogame.TerraformerID)
		}

		if planete.Facilities.ResearchLab < 12 && i == 0 {
			boot.BuildBuilding(planete.ID, ogame.ResearchLabID)
		}

		buildFormeVie(planete)
		if !SetExpedition(planete, coordExpe) {
			CargoExpeInsuffisant++
		}
		printCurrentconstruction(planete.ID, boot)
	}

	Researches(first_planet, slots)

	_, slots, err = boot.GetFleets()
	if err != nil {
		fmt.Println("2 Errors impossible d'obtenir les données sur les flottes: %s", err.Error())
		return
	}

	if slots.ExpInUse < slots.ExpTotal && slots.InUse < slots.Total {
		empire = sliceEmpireCargo(empire)
		for _, planete := range empire {
			co := planete.Coordinate
			co.Position = 16
			_, err := boot.SendFleet(planete.ID, getCompoFlotteExpe(planete), 100, co, ogame.Expedition, ogame.Resources{}, 0, 0)
			if err != nil {
				fmt.Printf("Erreur envoie expe restant : %s\n", err)
				if strings.Contains(err.Error(), "all slots are in use") {
					break
				}
			} else {

			}

			time.Sleep(4 * time.Second)
		}

		CargoExpeInsuffisant = 0
	}

	if slots.ExpInUse >= slots.ExpTotal && slots.InUse < slots.Total /*&& !HasMoonRes*/ {
		sort.Slice(empire, func(i int, j int) bool {
			resources_i := empire[i].Resources.Metal + empire[i].Resources.Crystal + empire[i].Resources.Deuterium
			resources_j := empire[j].Resources.Metal + empire[j].Resources.Crystal + empire[j].Resources.Deuterium
			HasValidForExplo := empire[i].Resources.Metal > 5000 && empire[i].Resources.Crystal > 5000 && empire[i].Resources.Deuterium > 5000
			return resources_i > resources_j && HasValidForExplo
		})
		setExploVie(empire[0].ID, empire[0].Coordinate)
	}
}
