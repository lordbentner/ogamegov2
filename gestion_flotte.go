package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/wrapper"
)

func setGhostFleet(bot *wrapper.OGame) {

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
	os.Exit(0)
}
