package main

import (
	"fmt"
	"os"

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
	for s := 2; s < 500; s++ {
		infos, err := boot.GalaxyInfos(int64(g), int64(s))
		if err == nil {
			somme := infos.ExpeditionDebris.Metal + infos.ExpeditionDebris.Crystal + infos.ExpeditionDebris.Deuterium
			max_debris := champDebris.ExpeditionDebris.Metal + infos.ExpeditionDebris.Crystal + champDebris.ExpeditionDebris.Deuterium
			fmt.Printf("somme = %d max_debris = %d\n", somme, max_debris)
			if somme > max_debris {
				champDebris = infos
			}
		}
	}
	//}

	fmt.Println(champDebris)
	os.Exit(0)
}
