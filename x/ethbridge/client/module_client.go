package client

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/dawn-protocol/dawn/x/ethbridge/client/cli"
	"github.com/dawn-protocol/dawn/x/ethbridge/client/rest"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	// Group ethbridge queries under a subcommand
	ethBBridgeQueryCmd := &cobra.Command{
		Use:   "ethbridge",
		Short: "Querying commands for the ethbridge module",
	}

	ethBBridgeQueryCmd.AddCommand(client.GetCommands(
		cli.GetCmdGetEthBridgeProphecy(storeKey, cdc),
	)...)

	return ethBBridgeQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	ethBridgeTxCmd := &cobra.Command{
		Use:   "ethbridge",
		Short: "EthBridge transactions subcommands",
	}

	ethBridgeTxCmd.AddCommand(client.PostCommands(
		cli.GetCmdCreateEthBridgeClaim(cdc),
	)...)

	return ethBridgeTxCmd
}

// RegisterRESTRoutes - Central function to define routes that get registered by the main application
func RegisterRESTRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	rest.RegisterRESTRoutes(cliCtx, r, storeName)
}
