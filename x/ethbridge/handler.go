package ethbridge

import (
	"fmt"

	"github.com/dawn-protocol/dawn/x/ethbridge/types"
	"github.com/dawn-protocol/dawn/x/oracle"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

// NewHandler returns a handler for "ethbridge" type messages.
func NewHandler(oracleKeeper oracle.Keeper, supplyKeeper supply.Keeper, codespace sdk.CodespaceType, cdc *codec.Codec) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateEthBridgeClaim:
			return handleMsgCreateEthBridgeClaim(ctx, cdc, oracleKeeper, supplyKeeper, msg, codespace)
		default:
			errMsg := fmt.Sprintf("unrecognized ethbridge message type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to create a bridge claim
func handleMsgCreateEthBridgeClaim(ctx sdk.Context, cdc *codec.Codec,
	oracleKeeper oracle.Keeper, supplyKeeper supply.Keeper, msg MsgCreateEthBridgeClaim,
	codespace sdk.CodespaceType) sdk.Result {
	oracleClaim, err := types.CreateOracleClaimFromEthClaim(cdc, types.EthBridgeClaim(msg))
	if err != nil {
		return types.ErrJSONMarshalling(codespace).Result()
	}
	status, sdkErr := oracleKeeper.ProcessClaim(ctx, oracleClaim)
	if sdkErr != nil {
		return sdkErr.Result()
	}

	if status.Text == oracle.SuccessStatusText {
		sdkErr = processSuccessfulClaim(ctx, supplyKeeper, status.FinalClaim)
		if sdkErr != nil {
			return sdkErr.Result()
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
		),
		sdk.NewEvent(
			types.EventTypeCreateClaim,
			sdk.NewAttribute(types.AttributeKeyEthereumSender, msg.EthereumSender.String()),
			sdk.NewAttribute(types.AttributeKeyCosmosReceiver, msg.CosmosReceiver.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			types.EventTypeProphecyStatus,
			sdk.NewAttribute(types.AttributeKeyStatus, status.Text.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func processSuccessfulClaim(ctx sdk.Context, supplyKeeper supply.Keeper, claim string) sdk.Error {
	oracleClaim, err := types.CreateOracleClaimFromOracleString(claim)
	if err != nil {
		return err
	}

	receiverAddress := oracleClaim.CosmosReceiver
	err = supplyKeeper.MintCoins(ctx, ModuleName, oracleClaim.Amount)
	if err != nil {
		return err
	}
	err = supplyKeeper.SendCoinsFromModuleToAccount(ctx, ModuleName, receiverAddress, oracleClaim.Amount)
	if err != nil {
		panic(err)
	}

	return nil
}