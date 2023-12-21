package handlers

import "fmt"

var (
	deroWalletAddress    string
	desiredRole          string
	registeredRole       string
	deroMembershipAmount int
	resChan              string
)

func SetRegisteredRole(registerationRole string) {
	registeredRole = registerationRole
}

func SetDesiredRole(membershipRole string) {
	desiredRole = membershipRole
}

func SetDeroAddress(walletAddress string) {
	deroWalletAddress = walletAddress
}

func SetMembershipAmount(membershipAmount int) {
	deroMembershipAmount = membershipAmount
	fmt.Println(deroMembershipAmount)
}

func SetResultsChannel(resultsChannel string) {
	resChan = resultsChannel
}
