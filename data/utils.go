package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	logrus "github.com/sirupsen/logrus"

	pylonSDK "github.com/Pylons-tech/pylons_sdk/cmd/test_utils"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/msgs"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetInitialPylons is a function to get initial pylons from faucet
func GetInitialPylons(username string) (string, error) {
	t := GetTestingT()
	addr := pylonSDK.GetAccountAddr(username, GetTestingT())
	sdkAddr, err := sdk.AccAddressFromBech32(addr)
	log.WithFields(log.Fields{
		"sdk_addr": sdkAddr.String(),
		"error":    err,
	}).Debugln("sdkAddr get result")

	// run create-account cli command
	result, logstr, err := pylonSDK.CreateChainAccount(username)
	log.WithFields(log.Fields{
		"result": result,
		"logstr": logstr,
		"error":  err,
	}).Info("creating account on chain result")
	if err != nil {
		log.WithFields(log.Fields{
			"result": result,
			"logstr": logstr,
			"error":  err,
		}).Fatal("error creating account on chain")
	}

	caTxHash := pylonSDK.GetTxHashFromLog(result)
	t.MustTrue(caTxHash != "", "error fetching txhash from result")
	t.WithFields(testing.Fields{
		"txhash": caTxHash,
	}).Info("started waiting for create account transaction")
	txResponseBytes, err := pylonSDK.WaitAndGetTxData(caTxHash, pylonSDK.GetMaxWaitBlock(), t)
	if err != nil {
		log.WithFields(log.Fields{
			"result": string(txResponseBytes),
			"error":  err,
		}).Fatal("error waiting for create account transaction")
	}
	pylonSDK.GetAccountInfoFromAddr(addr, t)

	getPylonsMsg := msgs.NewMsgGetPylons(types.PremiumTier.Fee, sdkAddr)
	txhash, err := pylonSDK.TestTxWithMsgWithNonce(t, getPylonsMsg, username, false)
	if err != nil {
		return "", fmt.Errorf("error sending transaction; %s: %+v", txhash, err)
	}
	return txhash, nil
}

// InitPylonAccount initialize an account on local and get initial balance from faucet
func InitPylonAccount(username string) string {

	return ""
}
