package main

import (
	"fmt"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
)

var maxCargo int = 120000000
var validCoordLF ogame.Coordinate = ogame.Coordinate{Galaxy: 5, System: 140, Position: 1}

func setExploVie(id ogame.CelestialID, coord ogame.Coordinate) int {
	fmt.Println("Gestion exploration forme de vie")
	nbError := 0
	if validCoordLF.Galaxy != 0 {
		coord = validCoordLF
	}
	err := boot.SendDiscoveryFleet(id, coord)
	if err != nil {
		nbError++
		fmt.Printf("%s: Erreur d'envoie explo vie : %s", coord, err)
	} else {
		fmt.Printf("fleet send to life discovery from %s to %s\n", coord.String(), coord)
		coord.Position = coord.Position + 1
		validCoordLF = getCorrectCoord(validCoordLF)
		time.Sleep(1 * time.Second)
		_, slots := boot.GetFleets()
		if slots.InUse < slots.Total {
			return setExploVie(id, validCoordLF)
		}
	}

	if nbError > 0 {
		fmt.Printf("%d erreurs d'envoie d'explo vie détectés\n", nbError)
		time.Sleep(1 * time.Second)
		validCoordLF = getCorrectCoord(ogame.Coordinate{Galaxy: coord.Galaxy, System: coord.System, Position: coord.Position + 1})
		return -1
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
		} else if sh.LightFighter > 0 {
			shipsInfos.LightFighter = 1
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

func SetExpedition(planete ogame.EmpireCelestial, coord ogame.Coordinate) bool {
	sh, _ := boot.GetShips(planete.ID)
	slots, _ := boot.GetSlots()
	totalCargo := (sh.LargeCargo * getCargoGT()) + (sh.SmallCargo * getCargoPT()) + (sh.Pathfinder * getCargoPathFinder())
	if slots.ExpInUse >= slots.ExpTotal || totalCargo < int64(maxCargo) || slots.InUse >= slots.Total {
		if totalCargo < int64(maxCargo) {
			fmt.Printf("Flotte non envoyée : Pas assez de capacité ======================> %d\n", totalCargo)
			return false
		}
		fmt.Println("Pas d'envoie de flotte pour l'instant ======================>")
		return true
	}

	fmt.Println("preparation envoi de flotte ===========================================>")
	shipsInfos := getFleetCompositionForExplo(sh)
	co := ogame.Coordinate{Galaxy: coord.Galaxy, System: coord.System, Position: 16}
	if planete.Coordinate.Galaxy > 4 || planete.Coordinate.System > 12 {
		co = planete.Coordinate
	}
	//bot.SendFleet(planete.ID, shipsInfos, 100, co, ogame.Expedition, ogame.Resources{}, 0, 0)
	_, err := boot.SendFleet(planete.ID, shipsInfos, 100, co, ogame.Expedition, ogame.Resources{}, 0, 0)
	if err != nil {
		pl := ogame.Coordinate{Galaxy: planete.Coordinate.Galaxy, System: planete.Coordinate.System, Position: 16}
		boot.SendFleet(planete.ID, shipsInfos, 100, pl, ogame.Expedition, ogame.Resources{}, 0, 0)
	}
	fmt.Printf("fleet send to expedition from %s with this fleet: ", coord.String())
	printStructFields(shipsInfos)
	time.Sleep(5000)
	return true
}

func gestionMessagesExpe() ogame.ExpeditionMessage {
	expMes, err := boot.GetExpeditionMessages(1)
	if err == nil {
		for i := 2; i < 20; i++ {
			exptest, er := boot.GetExpeditionMessages(int64(i))
			if er != nil {
				break
			} else {
				expMes = exptest
			}
		}
	}

	if expMes == nil || len(expMes) == 0 {
		var ff ogame.ExpeditionMessage
		return ff
	}

	fmt.Println(expMes[0].Content)
	return expMes[0]
}
