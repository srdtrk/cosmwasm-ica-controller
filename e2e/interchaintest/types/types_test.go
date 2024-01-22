package types_test

import (
	"testing"

	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"

	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos/wasm"
)

// This is some boilerplate test code to insert some tests for the types package.
// It is not meant to be executed, but to be used as a way to test some functions when
// debugging developing the types package.
func TestTypes(t *testing.T) {
	const testAddress = "srdtrk"

	t.Parallel()

	// Create deposit message:
	depositMsg := &govtypes.MsgDeposit{
		ProposalId: 1,
		Depositor:  testAddress,
		Amount:     sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(10000000))),
	}

	_, err := icatypes.SerializeCosmosTxWithEncoding(wasm.WasmEncoding().Codec, []proto.Message{depositMsg}, icatypes.EncodingProto3JSON)
	require.NoError(t, err)
}
