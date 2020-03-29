package main

import (
	"fmt"
	"net/http"

	"github.com/shinomontaz/tilesrv/config"
	"github.com/shinomontaz/tilesrv/internal/reader"
	"github.com/shinomontaz/tilesrv/internal/tile"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var env *config.Env

func init() {
	env = config.NewEnv("./config")
}

func main() {
	r := reader.New(env.Config.FileMap)

	// fmt.Println(env.Config.Styles)
	// panic("!")

	gIndex := r.Init()

	tlh := tile.New("/tiles", gIndex, env.Config.Styles)

	http.Handle("/tiles/", tlh)
	http.Handle("/metrics", promhttp.Handler())
	log.Debug("start prometheus handler on port: ", env.Config.ListenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", env.Config.ListenPort), nil))
}
