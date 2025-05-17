package main

import (
	"fmt"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/wrapper"
)

var tentaTiveExplovie int = 0
var maxCargo int = 120000000

func setExploVie(id ogame.CelestialID, coord ogame.Coordinate, bot *wrapper.OGame) int {
	slots, _ := bot.GetSlots()
	fmt.Println("Gestion exploration forme de vie")
	if slots.InUse >= slots.Total || slots.ExpInUse < slots.ExpTotal {
		fmt.Println("Exploration impossible car pas de slots disponible")
		return 0
	}
	att, _ := bot.IsUnderAttack()
	slotsDispo := int64(slots.Total - slots.InUse)
	nbError := 0
	co := coord
	if !att {
		for i := coord.Position; i < slotsDispo+int64(coord.Position); i++ {
			pos := int64(i + 1)
			gal := coord.Galaxy
			sys := coord.System

			co = getCorrectCoord(ogame.Coordinate{Galaxy: gal, System: sys, Position: pos})

			err := bot.SendDiscoveryFleet(id, co)
			if err != nil {
				nbError++
				fmt.Printf("%s: Erreur d'envoie explo vie : %s", co, err)
			} else {
				fmt.Printf("fleet send to life discovery from %s to %s\n", coord.String(), co)
			}
		}
	}

	if nbError > 0 {
		fmt.Printf("%d erreurs d'envoie d'explo vie détectés\n", nbError)
		tentaTiveExplovie++
		time.Sleep(5 * time.Second)
		if tentaTiveExplovie < 10 {
			return setExploVie(id, co, bot)
		} else {
			tentaTiveExplovie = 0
		}
	}

	return 0
}

func getFleetCompositionForExplo(sh ogame.ShipsInfos) ogame.ShipsInfos {
	var shipsInfos ogame.ShipsInfos
	if sh.EspionageProbe < 10 {
		shipsInfos.EspionageProbe = sh.EspionageProbe
	} else {
		shipsInfos.EspionageProbe = 10
	}

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

	shipsInfos.LargeCargo = int64(maxCargo) / getCargoGT()
	if shipsInfos.LargeCargo > sh.LargeCargo {
		shipsInfos.LargeCargo = sh.LargeCargo
		shipsInfos.SmallCargo = (int64(maxCargo) - shipsInfos.LargeCargo*getCargoGT()) / getCargoPT()
		if shipsInfos.SmallCargo > sh.SmallCargo {
			shipsInfos.SmallCargo = sh.SmallCargo
			cap_gt := int64(maxCargo) - shipsInfos.LargeCargo*getCargoGT()
			shipsInfos.Pathfinder = (cap_gt - shipsInfos.SmallCargo*getCargoPT()) / getCargoPathFinder()
			if shipsInfos.Pathfinder > sh.Pathfinder {
				shipsInfos.Pathfinder = sh.Pathfinder
			}
		}
	}

	if shipsInfos.SmallCargo == 0 && sh.SmallCargo > 0 {
		shipsInfos.SmallCargo++
	}

	if shipsInfos.LargeCargo == 0 && sh.LargeCargo > 0 {
		shipsInfos.LargeCargo++
	}

	if shipsInfos.Pathfinder == 0 && sh.Pathfinder > 0 {
		shipsInfos.Pathfinder++
	}

	return shipsInfos
}

func SetExpedition(planete ogame.EmpireCelestial, bot *wrapper.OGame, coord ogame.Coordinate) bool {
	sh, _ := bot.GetShips(planete.ID)
	slots, _ := bot.GetSlots()
	totalCargo := (sh.LargeCargo * getCargoGT()) + (sh.SmallCargo * getCargoPT()) + (sh.Pathfinder * getCargoPathFinder())
	if slots.ExpInUse >= slots.ExpTotal || totalCargo < int64(maxCargo) || slots.InUse >= slots.Total {
		if totalCargo < int64(maxCargo) {
			fmt.Printf("Pas assez de capacité ======================> %d\n", totalCargo)
		}
		fmt.Println("Pas d'envoie de flotte pour l'instant ======================>")
		return false
	}

	fmt.Println("preparation envoi de flotte ===========================================>")
	shipsInfos := getFleetCompositionForExplo(sh)
	co := ogame.Coordinate{Galaxy: coord.Galaxy, System: coord.System, Position: 16}
	if planete.Coordinate.Galaxy > 4 || planete.Coordinate.System > 12 {
		co = planete.Coordinate
	}
	//bot.SendFleet(planete.ID, shipsInfos, 100, co, ogame.Expedition, ogame.Resources{}, 0, 0)
	_, err := bot.SendFleet(planete.ID, shipsInfos, 100, co, ogame.Expedition, ogame.Resources{}, 0, 0)
	if err != nil {
		pl := ogame.Coordinate{Galaxy: planete.Coordinate.Galaxy, System: planete.Coordinate.System, Position: 16}
		bot.SendFleet(planete.ID, shipsInfos, 100, pl, ogame.Expedition, ogame.Resources{}, 0, 0)
	}
	fmt.Printf("fleet send to expedition from %s with this fleet: ", coord.String())
	printStructFields(shipsInfos)
	time.Sleep(5000)
	return true
}

func gestionMessagesExpe(bot *wrapper.OGame) ogame.ExpeditionMessage {
	expMes, err := bot.GetExpeditionMessages(1)
	if err == nil {
		for i := 2; i < 20; i++ {
			exptest, er := bot.GetExpeditionMessages(int64(i))
			if er != nil {
				break
			} else {
				expMes = exptest
			}
		}
	}

	fmt.Println(expMes[0].Content)
	return expMes[0]
}
