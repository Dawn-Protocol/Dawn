package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/gorilla/mux"

	"github.com/dawn-protocol/dawn/x/ethbridge/querier"
	"github.com/dawn-protocol/dawn/x/ethbridge/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	restEthereumChainID = "ethereumChainID"
	restBridgeContract  = "bridgeContract"
	restNonce           = "nonce"
	restSymbol          = "symbol"
	restTokenContract   = "tokenContract"
	restEthereumSender  = "ethereumSender"
)

type createEthClaimReq struct {
	BaseReq               rest.BaseReq `json:"base_req"`
	EthereumChainID       int          `json:"ethereum_chain_id"`
	BridgeContractAddress string       `json:"bridge_contract_address"`
	Nonce                 int          `json:"nonce"`
	Symbol                string       `json:"symbol"`
	TokenContractAddress  string       `json:"token_contract_address"`
	EthereumSender        string       `json:"ethereum_sender"`
	CosmosReceiver        string       `json:"cosmos_receiver"`
	Validator             string       `json:"validator"`
	Amount                string       `json:"amount"`
}

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRESTRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/prophecies", storeName), createClaimHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/prophecies/{%s}/{%s}/{%s}/{%s}/{%s}/{%s}", storeName, restEthereumChainID, restBridgeContract, restNonce, restSymbol, restTokenContract, restEthereumSender), getProphecyHandler(cliCtx, storeName)).Methods("GET")
}

func createClaimHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createEthClaimReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		bridgeContractAddress := types.NewEthereumAddress(req.BridgeContractAddress)

		tokenContractAddress := types.NewEthereumAddress(req.TokenContractAddress)

		ethereumSender := types.NewEthereumAddress(req.EthereumSender)

		cosmosReceiver, err := sdk.AccAddressFromBech32(req.CosmosReceiver)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		validator, err := sdk.ValAddressFromBech32(req.Validator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		amount, err := sdk.ParseCoins(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		ethBridgeClaim := types.NewEthBridgeClaim(req.EthereumChainID, bridgeContractAddress, req.Nonce, req.Symbol, tokenContractAddress, ethereumSender, cosmosReceiver, validator, amount)
		msg := types.NewMsgCreateEthBridgeClaim(ethBridgeClaim)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func getProphecyHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		ethereumChainID := vars[restEthereumChainID]
		ethereumChainIDString, err := strconv.Atoi(ethereumChainID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		bridgeContract := types.NewEthereumAddress(vars[restBridgeContract])

		nonce := vars[restNonce]
		nonceString, err := strconv.Atoi(nonce)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		tokenContract := types.NewEthereumAddress(vars[restTokenContract])

		symbol := vars[restSymbol]
		if strings.TrimSpace(symbol) == "" {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		ethereumSender := types.NewEthereumAddress(vars[restEthereumSender])

		bz, err := cliCtx.Codec.MarshalJSON(types.NewQueryEthProphecyParams(ethereumChainIDString, bridgeContract, nonceString, symbol, tokenContract, ethereumSender))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", storeName, querier.QueryEthProphecy)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}