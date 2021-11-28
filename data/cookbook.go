package data

import (

	pylonSDK "github.com/Pylons-tech/pylons_sdk/x/pylons/msgs"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func CreateCookbook() {
	id := "artbotCookbook"
	name := "artbotCookbook"
	desc := ""
	ver := types.SemVer("0.0.1")
	dev := "armen"
	email := types.Email("support@artbot.tv")
	level := types.Level(1)
	address := sdk.AccAddress("address")

	pylonSDK.NewMsgCreateCookbook(id, name, desc, dev, ver, email, level, 0, address)
}

func CreateRecipe(cookbookId string, addr sdk.AccAddress) {
	id := ""
	name := "createRecipe"
	coinInputs := types.CoinInputList(1)
	itemInputs := types.ItemInputList(1)
	entries := types.EntriesList(1)
	outputs := types.WeightedOutputsList(1)
	blockInterval := 1
	description := ""
	pylonSDK.NewMsgCreateRecipe()
}

func CreateAccount(){
	pylonSDK.NewMsgCreateAccount()
}

func ExecuteRecipe() {
	pylonSDK.NewMsgExecuteRecipe()
}