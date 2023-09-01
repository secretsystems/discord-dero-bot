package dero

var (
	// dero
	deroServerIP   string
	deroServerPort string
	deroWalletPort string
	deroUser       string
	deroPass       string
	homeDir        string
	pongAmount     int
	pongDir        string
)

func InitializeDERO(serverIP, walletPort, serverPort, user, pass, homedir, utilsdir string, deroMembershipAmount int) {
	deroServerIP = serverIP
	deroWalletPort = walletPort
	deroServerPort = serverPort
	deroUser = user
	deroPass = pass
	homeDir = homedir
	pongDir = utilsdir
	pongAmount = deroMembershipAmount

}

// Getter functions for accessing the variables from other packages
func GetDeroServerIP() string {
	return deroServerIP
}

func GetDeroServerPort() string {
	return deroServerPort
}

func GetDeroWalletPort() string {
	return deroWalletPort
}

func GetDeroUser() string {
	return deroUser
}

func GetDeroPass() string {
	return deroPass
}

func GetHomeDir() string {
	return homeDir
}

func GetPongAmount() int {
	return pongAmount
}

func GetPongDir() string {
	return pongDir
}
