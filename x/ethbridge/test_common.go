package ethbridge

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/dawn-protocol/dawn/x/ethbridge/types"
	oracle "github.com/dawn-protocol/dawn/x/oracle"
	keeperLib "github.com/dawn-protocol/dawn/x/oracle/keeper"
)

func CreateTestHandler(t *testing.T, consensusNeeded float64, validatorAmounts []int64) (sdk.Context, oracle.Keeper, bank.Keeper, supply.Keeper, []sdk.ValAddress, sdk.Handler) {
	ctx, oracleKeeper, bankKeeper, supplyKeeper, validatorAddresses := oracle.CreateTestKeepers(t, consensusNeeded, validatorAmounts, ModuleName)
	bridgeAccount := supply.NewEmptyModuleAccount(ModuleName, supply.Burner, supply.Minter)
	supplyKeeper.SetModuleAccount(ctx, bridgeAccount)

	cdc := keeperLib.MakeTestCodec()
	handler := NewHandler(oracleKeeper, supplyKeeper, types.DefaultCodespace, cdc)

	return ctx, oracleKeeper, bankKeeper, supplyKeeper, validatorAddresses, handler
}