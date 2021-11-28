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

	logrus "github.com/sirupsen/logrus

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
	log.Debugln("InitPylonAccount has started")
	var privKey string
	// "pylonscli keys add ${username}"
	addResult, _, err := pylonSDK.RunPylonsCli([]string{
		"keys", "add", username,
	}, "")

	log.WithFields(log.Fields{
		"addResult": string(addResult),
		"error":     err,
	}).Debugln("debug log")

	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			log.Warnln("pylonscli is not globally installed on your machine")
			SomethingWentWrongMsg = "pylonscli is not globally installed on your machine"
		} else {
			log.WithFields(log.Fields{
				"username": username,
			}).Infoln("using existing account")
			usr, _ := user.Current()
			pylonsDir := filepath.Join(usr.HomeDir, ".pylons")
			err = os.MkdirAll(pylonsDir, os.ModePerm)
			if err != nil {
				log.WithFields(log.Fields{
					"dir_path": pylonsDir,
				}).Fatal("create dir error")
			}
			userKeyFileName := username + ".json"
			keyFilePath := filepath.Join(pylonsDir, userKeyFileName)
			addResult, err = ioutil.ReadFile(keyFilePath)
			if err != nil && AutomateInput {
				log.WithFields(log.Fields{
					"key_file": userKeyFileName,
				}).Fatal("get private key error")
			}
			addedKeyResInterface := make(map[string]string)
			err = json.Unmarshal(addResult, &addedKeyResInterface)
			if err != nil && AutomateInput {
				log.WithFields(log.Fields{
					"key_file": userKeyFileName,
					"error":    err,
				}).Fatal("parse file error")
			}
			privKey = addedKeyResInterface["privkey"]
			log.WithFields(log.Fields{
				"privKey": privKey,
			}).Debugln("debug log")
		}
	} else {
		addedKeyResInterface := make(map[string]string)
		err = json.Unmarshal(addResult, &addedKeyResInterface)
		if err != nil {
			log.Fatal("Error unmarshalling into key result interface")
		}

		// mnemonic key from the pylonscli add result
		mnemonic := addedKeyResInterface["mnemonic"]
		log.WithFields(log.Fields{
			"mnemonic": mnemonic,
		}).Debugln("using mnemonic")

		privKey, _ = ComputePrivKeyFromMnemonic(mnemonic) // get privKey and cosmosAddr

		addResult, err = json.Marshal(addedKeyResInterface)
		if err != nil {
			log.Fatal("marshal added keys result error")
		}

		usr, _ := user.Current()
		pylonsDir := filepath.Join(usr.HomeDir, ".pylons")
		err = os.MkdirAll(pylonsDir, os.ModePerm)
		if err != nil {
			log.WithFields(log.Fields{
				"dir_path": pylonsDir,
			}).Fatal("create directory error")
		}
		userKeyFileName := username + ".json"
		keyFilePath := filepath.Join(pylonsDir, userKeyFileName)
		if ioutil.WriteFile(keyFilePath, addResult, 0644) != nil {
			log.WithFields(log.Fields{
				"dir_path": pylonsDir,
			}).Fatal("error writing file to directory")
		}
		log.WithFields(log.Fields{
			"privKey": privKey,
		}).Debugln("debug log")
		log.WithFields(log.Fields{
			"username":  username,
			"file_path": pylonsDir + "/" + userKeyFileName,
		}).Infoln("created new account")
	}
	addr := pylonSDK.GetAccountAddr(username, GetTestingT())
	accBytes, _, err := pylonSDK.RunPylonsCli([]string{"query", "account", addr}, "")
	log.WithFields(log.Fields{
		"address": addr,
		"result":  string(accBytes),
		"error":   err,
	}).Debugln("query account")
	if err != nil {
		if strings.Contains(string(accBytes), "dial tcp [::1]:26657: connect: connection refused") { // Daemon is off
			log.WithFields(log.Fields{
				"error": "daemon connection refuse",
			}).Fatalln("please check daemon is running!")
		} else { // account does not exist
			txhash, err := GetInitialPylons(username)
			if err != nil {
				log.WithFields(log.Fields{
					"txhash": txhash,
					"error":  err,
				}).Fatalln("GetInitialPylons result")
			}
			log.WithFields(log.Fields{
				"address": addr,
			}).Debugln("ran command for new account on remote chain and waiting for next block ...")
			if pylonSDK.WaitForNextBlock() != nil {
				return "error waiting for block"
			}
		}
	} else {
		log.WithFields(log.Fields{
			"address": addr,
		}).Infoln("using existing account on remote chain")
	}

	// Remove nonce file
	log.Debugln("start removing nonce file")
	nonceRootDir := "./"
	nonceFile := filepath.Join(nonceRootDir, "nonce.json")
	err = os.Remove(nonceFile)
	log.WithFields(log.Fields{
		"error": err,
	}).Debugln("remove nonce file result")

	log.WithFields(log.Fields{
		"privKey": privKey,
	}).Debugln("function ended")
	return privKey
}
