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


func CreateAccount(addr sdk.AccAddress){
	pylonSDK.NewMsgCreateAccount(addr)
}

//recipes

func CreateRecipe(cookbookId string, addr sdk.AccAddress) {
	id := ""
	name := "createRecipe"
	coinInputs := types.CoinInputList(1)
	itemInputs := types.ItemInputList(1)
	entries := types.EntriesList(1)
	outputs := types.WeightedOutputsList(1)
	blockInterval := int64(1)
	description := ""
	pylonSDK.NewMsgCreateRecipe(name, cookbookId, id,
		description, coinInputs, itemInputs, entries,
		outputs, blockInterval, addr)
}

func CreateItem(cookbookID string,
		doubles []types.DoubleKeyValue,
		longs []types.LongKeyValue,
		strings []types.StringKeyValue,
		sender sdk.AccAddress,
		blockHeight int64,
		transferFee int64) {
	types.NewItem(cookbookID,doubles, longs, strings, sender, blockHeight, transferFee)
}

func ExecuteRecipe(recipeId string, addr sdk.AccAddress, itemIds []string) {

	/*
		t := GetTestingT()
		if len(rcpName) == 0 {
			return "", errors.New("Recipe Name does not exist")
		}
		rcpID, ok := RcpIDs[rcpName]
		if !ok {
			return "", errors.New("RecipeID does not exist for rcpName=" + rcpName)
		}
		addr := pylonSDK.GetAccountAddr(user.GetUserName(), GetTestingT())
		sdkAddr, _ := sdk.AccAddressFromBech32(addr)
		execMsg := msgs.NewMsgExecuteRecipe(rcpID, sdkAddr, itemIDs)
		txhash, err := pylonSDK.TestTxWithMsgWithNonce(t, execMsg, user.GetUserName(), false)
		if err != nil {
			return "", fmt.Errorf("error sending transaction; %s: %+v", txhash, err)
		}
		user.SetLastTransaction(txhash, rcpName)
		return txhash, nil
	 */
	pylonSDK.NewMsgExecuteRecipe(recipeId, addr, itemIds)
	pylonSDK.msg
}