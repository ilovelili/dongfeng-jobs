package jobs

// HeIsDeadJim always returns error for testing failure
func HeIsDeadJim() int {
	operationName := "HeIsDeadJim"
	errorlog("He's dead, Jim", operationName)
	return 1
}
