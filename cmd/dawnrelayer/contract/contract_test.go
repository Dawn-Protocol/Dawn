package contract

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLoadABI : test that contract containing named event is successfully loaded
func TestLoadABI(t *testing.T) {

	const AbiPath = "/src/github.com/dawn-protocol/dawn/cmd/dawnrelayer/contract/abi/Peggy.abi"

	//Get the ABI ready
	abi := LoadABI()

	require.NotNil(t, abi.Events["LogLock"])
}