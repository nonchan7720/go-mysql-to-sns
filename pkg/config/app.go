package config

import (
	"runtime/debug"
	"strings"

	"github.com/nonchan7720/go-storage-to-messenger/pkg/utils"
)

type App struct {
	AppEnv      string `yaml:"appEnv" default:"dev"`
	TrackingEnv string `yaml:"trackingEnv" default:"dev"`
	ServiceName string `yaml:"serviceName"`
	ServiceRoot string `yaml:"serviceRoot"`
}

func (app *App) SetDefaults() {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		path := bi.Path
		app.ServiceRoot = path
		paths := strings.Split(path, "/")
		app.ServiceName = paths[len(paths)-1]
	}
}

func AppEnv() string {
	return DefaultConfig().App.AppEnv
}

func IsProduction() bool {
	return strings.EqualFold(AppEnv(), "production")
}

func IsStaging() bool {
	return strings.EqualFold(AppEnv(), "staging")
}

func IsDev() bool {
	return strings.EqualFold(AppEnv(), "dev")
}

func IsTest() bool {
	return strings.EqualFold(AppEnv(), "test")
}

func IsCIorTest() bool {
	return utils.IsCI() || IsTest()
}

func IsDevOrCIorTest() bool {
	return IsDev() || utils.IsCI() || IsTest()
}
