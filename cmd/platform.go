package main

import apphttp "github.com/Krokozabra213/sso/internal/platform/app/http"

const (
	configfile = "settings/platform_main.yml"
	envfile    = "platform.env"
)

func main() {
	apphttp.Run(configfile, envfile)
}
