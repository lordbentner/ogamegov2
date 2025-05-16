package main

import (
	"fmt"
	"sort"

	"github.com/alaingilbert/ogame/pkg/ogame"
)

func sendFleetFromMoonToPlanet(moon ogame.EmpireCelestial) {
	if moon.Resources.Metal <= 0 && moon.Resources.Crystal <= 0 {
		return
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
		return
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
