package main

import (
	"fmt"
	"sort"

	"github.com/alaingilbert/ogame/pkg/ogame"
)

func sliceEmpireCargo(empire []ogame.EmpireCelestial) []ogame.EmpireCelestial {
	sort.Slice(empire, func(i int, j int) bool {
		cargoTotali := getCargoGT()*empire[i].Ships.LargeCargo + getCargoPT()*empire[i].Ships.SmallCargo
		cargoTotali += getCargoPathFinder() * empire[i].Ships.Pathfinder
		cargoTotalj := getCargoGT()*empire[j].Ships.LargeCargo + getCargoPT()*empire[j].Ships.SmallCargo
		cargoTotalj += getCargoPathFinder() * empire[j].Ships.Pathfinder
		return cargoTotali > cargoTotalj
	})

	return empire
}

func sendFleetFromMoonToPlanet(moon ogame.EmpireCelestial) bool {
	if moon.Ships.LargeCargo <= 10 && moon.Ships.SmallCargo <= 50 {
		return false
	}
	if moon.Resources.Metal <= 0 && moon.Resources.Crystal <= 0 && moon.Resources.Deuterium < 1500000 {
		return false
	}
	r := moon.Resources
	var sh ogame.ShipsInfos
	sh.LargeCargo = moon.Ships.LargeCargo
	sh.SmallCargo = moon.Ships.SmallCargo
	sh.Pathfinder = moon.Ships.Pathfinder
	if r.Deuterium > 1500000 {
		r.Deuterium -= 1500000
	} else {
		r.Deuterium = 0
	}

	resources := ogame.Resources{Metal: r.Metal, Crystal: r.Crystal, Deuterium: r.Deuterium}
	totalRes := resources.Metal + resources.Crystal + resources.Deuterium
	if sh.LargeCargo*getCargoGT()+sh.SmallCargo*getCargoPT()+sh.Pathfinder*getCargoPathFinder() > totalRes {
		_, err := boot.SendFleet(moon.ID, sh, 100, moon.Coordinate.Planet(), ogame.Transport, r, 0, 0)
		if err != nil {
			fmt.Printf("err sendFleetFromMoonToPlanet : %s\n", err)
		}
		return true
	}

	if r.Crystal > sh.LargeCargo*getCargoGT() {
		resources.Crystal = sh.LargeCargo * getCargoGT()
	}

	if r.Metal > sh.SmallCargo*getCargoPT() {
		resources.Metal = sh.SmallCargo * getCargoPT()
	}

	if r.Deuterium > sh.Pathfinder*getCargoPathFinder() {
		resources.Deuterium = sh.Pathfinder * getCargoPathFinder()
	}

	_, err := boot.SendFleet(moon.ID, sh, 100, moon.Coordinate.Planet(), ogame.Transport, resources, 0, 0)
	if err != nil {
		fmt.Printf("err sendFleetFromMoonToPlanet : %s\n", err)
	}

	return true
}

func getMaxExpeDebris(g int) {
	champDebris, errs := boot.GalaxyInfos(1, 2)
	if errs != nil {
		return
	}

	//for g := 1; g < 7; g++ {
	var list []ogame.SystemInfos
	for s := 2; s < 500; s++ {
		infos, err := boot.GalaxyInfos(int64(g), int64(s))
		if err == nil {
			list = append(list, infos)
			somme := infos.ExpeditionDebris.Metal + infos.ExpeditionDebris.Crystal + infos.ExpeditionDebris.Deuterium
			max_debris := champDebris.ExpeditionDebris.Metal + infos.ExpeditionDebris.Crystal + champDebris.ExpeditionDebris.Deuterium
			//fmt.Printf("somme = %d max_debris = %d\n", somme, max_debris)
			if somme > max_debris {
				champDebris = infos
			}
		}
	}
	//}

	sort.Slice(list, func(i int, j int) bool {
		somme_i := list[i].ExpeditionDebris.Metal + list[i].ExpeditionDebris.Crystal + list[i].ExpeditionDebris.Deuterium
		somme_j := list[j].ExpeditionDebris.Metal + list[j].ExpeditionDebris.Crystal + list[j].ExpeditionDebris.Deuterium
		return somme_i > somme_j
	})

	for i := 0; i < 8; i++ {
		fmt.Println(list[i])
	}
}

func getCompoFlotteExpe(planete ogame.EmpireCelestial) ogame.ShipsInfos {
	var shipsInfos ogame.ShipsInfos
	shipsInfos.LargeCargo = planete.Ships.LargeCargo
	shipsInfos.SmallCargo = planete.Ships.SmallCargo
	if planete.Ships.Pathfinder > 0 {
		shipsInfos.Pathfinder = 1
	}
	if planete.Ships.EspionageProbe >= 10 {
		shipsInfos.EspionageProbe = 10
	} else {
		shipsInfos.EspionageProbe = planete.Ships.EspionageProbe
	}

	return shipsInfos
}

func sendFleetToMoon(moon ogame.EmpireCelestial) {
	empire, _ := boot.GetEmpire(ogame.PlanetType)
	for _, pl := range empire {
		sys := pl.Coordinate.System == moon.Coordinate.System
		if pl.Coordinate.Galaxy == moon.Coordinate.Galaxy && sys && pl.Coordinate.Position == moon.Coordinate.Position {
			//boot.SendFleet(moon.ID, sh, 100, moon.Coordinate.Planet(), ogame.Transport, resources, 0, 0)
			res := ogame.Resources{}
			if moon.Resources.Deuterium < 200000 && pl.Resources.Deuterium > 200000 {
				res.Deuterium = 200000
			}

			_, err := boot.SendFleet(pl.ID, pl.Ships, 100, moon.Coordinate, ogame.Park, res, 0, 0)
			if err != nil {
				fmt.Printf("Error sendFleetToMoon : %s\n", err)
			}
			break
		}
	}
}
