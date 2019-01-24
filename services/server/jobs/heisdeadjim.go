package jobs

import "github.com/micro/cli"

// HeIsDeadJim always returns error for testing failure
func HeIsDeadJim(ctx *cli.Context) int {
	operationName := "HeIsDeadJim"
	errorlog("He's dead, Jim", operationName)
	return 1
}
