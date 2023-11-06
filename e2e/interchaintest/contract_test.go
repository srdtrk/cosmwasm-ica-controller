package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"

	mysuite "github.com/srdtrk/cw-ica-controller/interchaintest/v2/testsuite"
	"github.com/srdtrk/cw-ica-controller/interchaintest/v2/types"
)

type ContractTestSuite struct {
	mysuite.TestSuite

	Contract   *types.IcaContract
	IcaAddress string
}

// SetupContractTestSuite starts the chains, relayer, creates the user accounts, creates the ibc clients and connections,
// sets up the contract and does the channel handshake for the contract test suite.
func (s *ContractTestSuite) SetupContractTestSuite(ctx context.Context, encoding string) {
	s.SetupSuite(ctx, chainSpecs)

	codeId, err := s.ChainA.StoreContract(ctx, s.UserA.KeyName(), "../../artifacts/cw_ica_controller.wasm")
	s.Require().NoError(err)

	// Instantiate the contract with channel:
	instantiateMsg := types.NewInstantiateMsgWithChannelInitOptions(nil, s.ChainAConnID, s.ChainBConnID, nil, &encoding)

	contractAddr, err := s.ChainA.InstantiateContract(ctx, s.UserA.KeyName(), codeId, instantiateMsg, true, "--gas", "500000")
	s.Require().NoError(err)

	s.Contract = types.NewIcaContract(types.NewContract(contractAddr, codeId, s.ChainA))

	// Wait for the channel to get set up
	err = testutil.WaitForBlocks(ctx, 5, s.ChainA, s.ChainB)
	s.Require().NoError(err)

	contractState, err := s.Contract.QueryContractState(ctx)
	s.Require().NoError(err)

	s.IcaAddress = contractState.IcaInfo.IcaAddress
	s.Contract.SetIcaAddress(s.IcaAddress)
}

func TestWithContractTestSuite(t *testing.T) {
	suite.Run(t, new(ContractTestSuite))
}

func (s *ContractTestSuite) TestIcaContractChannelHandshake() {
	ctx := context.Background()

	// This starts the chains, relayer, creates the user accounts, creates the ibc clients and connections,
	// sets up the contract and does the channel handshake for the contract test suite.
	s.SetupContractTestSuite(ctx, icatypes.EncodingProto3JSON)
	wasmd, simd := s.ChainA, s.ChainB
	wasmdUser := s.UserA

	s.Run("TestChannelHandshakeSuccess", func() {
		// Test if the handshake was successful
		wasmdChannels, err := s.Relayer.GetChannels(ctx, s.ExecRep, wasmd.Config().ChainID)
		s.Require().NoError(err)
		s.Require().Equal(1, len(wasmdChannels))

		wasmdChannel := wasmdChannels[0]
		s.T().Logf("wasmd channel: %s", toJSONString(wasmdChannel))
		s.Require().Equal(s.Contract.Port(), wasmdChannel.PortID)
		s.Require().Equal(icatypes.HostPortID, wasmdChannel.Counterparty.PortID)
		s.Require().Equal(channeltypes.OPEN.String(), wasmdChannel.State)

		simdChannels, err := s.Relayer.GetChannels(ctx, s.ExecRep, simd.Config().ChainID)
		s.Require().NoError(err)
		// I don't know why sometimes an extra channel is created in simd.
		// this is not related to the localhost connection, and is a failed
		// clone of the successful channel at index 0. I will log it for now.
		s.Require().Greater(len(simdChannels), 0)
		if len(simdChannels) > 1 {
			s.T().Logf("extra simd channels detected: %s", toJSONString(simdChannels))
		}

		simdChannel := simdChannels[0]
		s.T().Logf("simd channel state: %s", toJSONString(simdChannel.State))
		s.Require().Equal(icatypes.HostPortID, simdChannel.PortID)
		s.Require().Equal(s.Contract.Port(), simdChannel.Counterparty.PortID)
		s.Require().Equal(channeltypes.OPEN.String(), simdChannel.State)

		// Check contract's channel state
		contractChannelState, err := s.Contract.QueryChannelState(ctx)
		s.Require().NoError(err)

		s.T().Logf("contract's channel store after handshake: %s", toJSONString(contractChannelState))

		s.Require().Equal(wasmdChannel.State, contractChannelState.ChannelStatus)
		s.Require().Equal(wasmdChannel.Version, contractChannelState.Channel.Version)
		s.Require().Equal(wasmdChannel.ConnectionHops[0], contractChannelState.Channel.ConnectionID)
		s.Require().Equal(wasmdChannel.ChannelID, contractChannelState.Channel.Endpoint.ChannelID)
		s.Require().Equal(wasmdChannel.PortID, contractChannelState.Channel.Endpoint.PortID)
		s.Require().Equal(wasmdChannel.Counterparty.ChannelID, contractChannelState.Channel.CounterpartyEndpoint.ChannelID)
		s.Require().Equal(wasmdChannel.Counterparty.PortID, contractChannelState.Channel.CounterpartyEndpoint.PortID)
		s.Require().Equal(wasmdChannel.Ordering, contractChannelState.Channel.Order)

		// Check contract state
		contractState, err := s.Contract.QueryContractState(ctx)
		s.Require().NoError(err)

		s.Require().Equal(wasmdUser.FormattedAddress(), contractState.Admin)
		s.Require().Equal(wasmdChannel.ChannelID, contractState.IcaInfo.ChannelID)
	})
}

func (s *ContractTestSuite) TestIcaRelayerInstantiatedChannelHandshake() {
	ctx := context.Background()

	// This starts the chains, relayer, creates the user accounts, and creates the ibc clients and connections.
	s.SetupSuite(ctx, chainSpecs)
	wasmd, simd := s.ChainA, s.ChainB
	wasmdUser := s.UserA

	var err error
	// Upload and Instantiate the contract on wasmd:
	s.Contract, err = types.StoreAndInstantiateNewIcaContract(ctx, wasmd, wasmdUser.KeyName(), "../../artifacts/cw_ica_controller.wasm")
	s.Require().NoError(err)

	version := fmt.Sprintf(`{"version":"%s","controller_connection_id":"%s","host_connection_id":"%s","address":"","encoding":"%s","tx_type":"%s"}`, icatypes.Version, s.ChainAConnID, s.ChainBConnID, icatypes.EncodingProtobuf, icatypes.TxTypeSDKMultiMsg)
	err = s.Relayer.CreateChannel(ctx, s.ExecRep, s.PathName, ibc.CreateChannelOptions{
		SourcePortName: s.Contract.Port(),
		DestPortName:   icatypes.HostPortID,
		Order:          ibc.Ordered,
		// cannot use an empty version here, see README
		Version: version,
	})
	s.Require().NoError(err)

	// Wait for the channel to get set up
	err = testutil.WaitForBlocks(ctx, 5, s.ChainA, s.ChainB)
	s.Require().NoError(err)

	contractState, err := s.Contract.QueryContractState(ctx)
	s.Require().NoError(err)
	s.IcaAddress = contractState.IcaInfo.IcaAddress

	s.Run("TestChannelHandshakeSuccess", func() {
		// Test if the handshake was successful
		wasmdChannels, err := s.Relayer.GetChannels(ctx, s.ExecRep, wasmd.Config().ChainID)
		s.Require().NoError(err)
		s.Require().Equal(1, len(wasmdChannels))

		wasmdChannel := wasmdChannels[0]
		s.T().Logf("wasmd channel: %s", toJSONString(wasmdChannel))
		s.Require().Equal(s.Contract.Port(), wasmdChannel.PortID)
		s.Require().Equal(icatypes.HostPortID, wasmdChannel.Counterparty.PortID)
		s.Require().Equal(channeltypes.OPEN.String(), wasmdChannel.State)

		simdChannels, err := s.Relayer.GetChannels(ctx, s.ExecRep, simd.Config().ChainID)
		s.Require().NoError(err)
		// I don't know why sometimes an extra channel is created in simd.
		// this is not related to the localhost connection, and is a failed
		// clone of the successful channel at index 0. I will log it for now.
		s.Require().Greater(len(simdChannels), 0)
		if len(simdChannels) > 1 {
			s.T().Logf("extra simd channels detected: %s", toJSONString(simdChannels))
		}

		simdChannel := simdChannels[0]
		s.T().Logf("simd channel state: %s", toJSONString(simdChannel.State))
		s.Require().Equal(icatypes.HostPortID, simdChannel.PortID)
		s.Require().Equal(s.Contract.Port(), simdChannel.Counterparty.PortID)
		s.Require().Equal(channeltypes.OPEN.String(), simdChannel.State)

		// Check contract's channel state
		contractChannelState, err := s.Contract.QueryChannelState(ctx)
		s.Require().NoError(err)

		s.T().Logf("contract's channel store after handshake: %s", toJSONString(contractChannelState))

		s.Require().Equal(wasmdChannel.State, contractChannelState.ChannelStatus)
		s.Require().Equal(wasmdChannel.Version, contractChannelState.Channel.Version)
		s.Require().Equal(wasmdChannel.ConnectionHops[0], contractChannelState.Channel.ConnectionID)
		s.Require().Equal(wasmdChannel.ChannelID, contractChannelState.Channel.Endpoint.ChannelID)
		s.Require().Equal(wasmdChannel.PortID, contractChannelState.Channel.Endpoint.PortID)
		s.Require().Equal(wasmdChannel.Counterparty.ChannelID, contractChannelState.Channel.CounterpartyEndpoint.ChannelID)
		s.Require().Equal(wasmdChannel.Counterparty.PortID, contractChannelState.Channel.CounterpartyEndpoint.PortID)
		s.Require().Equal(wasmdChannel.Ordering, contractChannelState.Channel.Order)

		// Check contract state
		contractState, err := s.Contract.QueryContractState(ctx)
		s.Require().NoError(err)

		s.Require().Equal(wasmdUser.FormattedAddress(), contractState.Admin)
		s.Require().Equal(wasmdChannel.ChannelID, contractState.IcaInfo.ChannelID)
	})
}

func (s *ContractTestSuite) TestRecoveredIcaContractInstantiatedChannelHandshake() {
	ctx := context.Background()

	s.SetupSuite(ctx, chainSpecs)
	wasmd, simd := s.ChainA, s.ChainB
	wasmdUser := s.UserA

	codeId, err := wasmd.StoreContract(ctx, wasmdUser.KeyName(), "../../artifacts/cw_ica_controller.wasm")
	s.Require().NoError(err)

	s.Run("TestChannelHandshakeFail: invalid connection id", func() {
		// Instantiate the contract with channel:
		instantiateMsg := types.NewInstantiateMsgWithChannelInitOptions(nil, "invalid", s.ChainBConnID, nil, nil)

		_, err = wasmd.InstantiateContract(ctx, wasmdUser.KeyName(), codeId, instantiateMsg, true, "--gas", "500000")
		s.Require().ErrorContains(err, "submessages: invalid connection hop ID")
	})

	s.Run("TestChannelHandshakeFail: invalid counterparty connection id", func() {
		// Instantiate the contract with channel:
		instantiateMsg := types.NewInstantiateMsgWithChannelInitOptions(nil, s.ChainAConnID, "connection-123", nil, nil)

		// unfortunately, this doesn't error out because the connection id is in the counterparty.
		// instead, the handshake never completes. A new channel may be created by the relayer.
		contractAddr, err := wasmd.InstantiateContract(ctx, wasmdUser.KeyName(), codeId, instantiateMsg, true, "--gas", "500000")
		s.Require().NoError(err)

		s.Contract = types.NewIcaContract(types.NewContract(contractAddr, codeId, wasmd))
	})

	s.Run("TestChannelHandshakeSuccessAfterFail", func() {
		err = s.Contract.ExecCreateChannel(ctx, wasmdUser.KeyName(), s.ChainAConnID, s.ChainBConnID, nil, nil, "--gas", "500000")
		s.Require().NoError(err)

		// Wait for the channel to get set up
		err = testutil.WaitForBlocks(ctx, 8, s.ChainA, s.ChainB)
		s.Require().NoError(err)

		// Test if the handshake was successful
		wasmdChannels, err := s.Relayer.GetChannels(ctx, s.ExecRep, wasmd.Config().ChainID)
		s.Require().NoError(err)
		s.Require().Equal(2, len(wasmdChannels))

		wasmdChannel := wasmdChannels[1]
		s.T().Logf("wasmd channel: %s", toJSONString(wasmdChannel))
		s.Require().Equal(s.Contract.Port(), wasmdChannel.PortID)
		s.Require().Equal(icatypes.HostPortID, wasmdChannel.Counterparty.PortID)
		s.Require().Equal(channeltypes.OPEN.String(), wasmdChannel.State)

		simdChannels, err := s.Relayer.GetChannels(ctx, s.ExecRep, simd.Config().ChainID)
		s.Require().NoError(err)
		// I don't know why sometimes an extra channel is created in simd.
		// this is not related to the localhost connection, and is a failed
		// clone of the successful channel at index 0. I will log it for now.
		s.Require().Greater(len(simdChannels), 0)
		if len(simdChannels) > 1 {
			s.T().Logf("extra simd channels detected: %s", toJSONString(simdChannels))
		}

		simdChannel := simdChannels[0]
		s.T().Logf("simd channel state: %s", toJSONString(simdChannel.State))
		s.Require().Equal(icatypes.HostPortID, simdChannel.PortID)
		s.Require().Equal(s.Contract.Port(), simdChannel.Counterparty.PortID)
		s.Require().Equal(channeltypes.OPEN.String(), simdChannel.State)

		// Check contract's channel state
		contractChannelState, err := s.Contract.QueryChannelState(ctx)
		s.Require().NoError(err)

		s.T().Logf("contract's channel store after handshake: %s", toJSONString(contractChannelState))

		s.Require().Equal(wasmdChannel.State, contractChannelState.ChannelStatus)
		s.Require().Equal(wasmdChannel.Version, contractChannelState.Channel.Version)
		s.Require().Equal(wasmdChannel.ConnectionHops[0], contractChannelState.Channel.ConnectionID)
		s.Require().Equal(wasmdChannel.ChannelID, contractChannelState.Channel.Endpoint.ChannelID)
		s.Require().Equal(wasmdChannel.PortID, contractChannelState.Channel.Endpoint.PortID)
		s.Require().Equal(wasmdChannel.Counterparty.ChannelID, contractChannelState.Channel.CounterpartyEndpoint.ChannelID)
		s.Require().Equal(wasmdChannel.Counterparty.PortID, contractChannelState.Channel.CounterpartyEndpoint.PortID)
		s.Require().Equal(wasmdChannel.Ordering, contractChannelState.Channel.Order)

		// Check contract state
		contractState, err := s.Contract.QueryContractState(ctx)
		s.Require().NoError(err)

		s.Require().Equal(wasmdUser.FormattedAddress(), contractState.Admin)
		s.Require().Equal(wasmdChannel.ChannelID, contractState.IcaInfo.ChannelID)
	})
}

func (s *ContractTestSuite) TestIcaContractExecutionProto3JsonEncoding() {
	s.IcaContractExecutionTestWithEncoding(icatypes.EncodingProto3JSON)
}

func (s *ContractTestSuite) TestIcaContractExecutionProtobufEncoding() {
	s.IcaContractExecutionTestWithEncoding(icatypes.EncodingProtobuf)
}

func (s *ContractTestSuite) IcaContractExecutionTestWithEncoding(encoding string) {
	ctx := context.Background()

	// This starts the chains, relayer, creates the user accounts, creates the ibc clients and connections,
	// sets up the contract and does the channel handshake for the contract test suite.
	s.SetupContractTestSuite(ctx, encoding)
	wasmd, simd := s.ChainA, s.ChainB
	wasmdUser := s.UserA

	// Fund the ICA address:
	s.FundAddressChainB(ctx, s.IcaAddress)

	s.Run(fmt.Sprintf("TestSendCustomIcaMessagesSuccess-%s", encoding), func() {
		// Send custom ICA messages through the contract:
		// Let's create a governance proposal on simd and deposit some funds to it.
		testProposal := govtypes.TextProposal{
			Title:       "IBC Gov Proposal",
			Description: "tokens for all!",
		}
		protoAny, err := codectypes.NewAnyWithValue(&testProposal)
		s.Require().NoError(err)
		proposalMsg := &govtypes.MsgSubmitProposal{
			Content:        protoAny,
			InitialDeposit: sdk.NewCoins(sdk.NewCoin(simd.Config().Denom, sdkmath.NewInt(5000))),
			Proposer:       s.IcaAddress,
		}

		// Create deposit message:
		depositMsg := &govtypes.MsgDeposit{
			ProposalId: 1,
			Depositor:  s.IcaAddress,
			Amount:     sdk.NewCoins(sdk.NewCoin(simd.Config().Denom, sdkmath.NewInt(10000000))),
		}

		// Execute the contract:
		err = s.Contract.ExecCustomIcaMessages(ctx, wasmdUser.KeyName(), []proto.Message{proposalMsg, depositMsg}, encoding, nil, nil)
		s.Require().NoError(err)

		err = testutil.WaitForBlocks(ctx, 5, wasmd, simd)
		s.Require().NoError(err)

		// Check if contract callbacks were executed:
		callbackCounter, err := s.Contract.QueryCallbackCounter(ctx)
		s.Require().NoError(err)

		s.Require().Equal(uint64(1), callbackCounter.Success)
		s.Require().Equal(uint64(0), callbackCounter.Error)

		// Check if the proposal was created:
		proposal, err := simd.QueryProposal(ctx, "1")
		s.Require().NoError(err)
		s.Require().Equal(simd.Config().Denom, proposal.TotalDeposit[0].Denom)
		s.Require().Equal(fmt.Sprint(10000000+5000), proposal.TotalDeposit[0].Amount)
		// We do not check title and description of the proposal because this is a legacy proposal.
	})

	s.Run(fmt.Sprintf("TestSendCustomIcaMessagesError-%s", encoding), func() {
		// Test erroneous callback:
		// Send incorrect custom ICA messages through the contract:
		badMessage := base64.StdEncoding.EncodeToString([]byte("bad message"))
		badCustomMsg := `{"send_custom_ica_messages":{"messages":"` + badMessage + `"}}`

		// Execute the contract:
		err := s.Contract.Execute(ctx, wasmdUser.KeyName(), badCustomMsg)
		s.Require().NoError(err)

		err = testutil.WaitForBlocks(ctx, 5, wasmd, simd)
		s.Require().NoError(err)

		// Check if contract callbacks were executed:
		callbackCounter, err := s.Contract.QueryCallbackCounter(ctx)
		s.Require().NoError(err)
		s.Require().Equal(uint64(1), callbackCounter.Success)
		s.Require().Equal(uint64(1), callbackCounter.Error)
		s.Require().Equal(uint64(0), callbackCounter.Timeout)
	})
}

func (s *ContractTestSuite) TestIcaContractTimeoutPacket() {
	ctx := context.Background()

	// This starts the chains, relayer, creates the user accounts, creates the ibc clients and connections,
	// sets up the contract and does the channel handshake for the contract test suite.
	s.SetupContractTestSuite(ctx, icatypes.EncodingProto3JSON)
	wasmd, simd := s.ChainA, s.ChainB
	wasmdUser, _ := s.UserA, s.UserB

	// Fund the ICA address:
	s.FundAddressChainB(ctx, s.IcaAddress)

	contractState, err := s.Contract.QueryContractState(ctx)
	s.Require().NoError(err)

	var simdChannelsLen int
	s.Run("TestTimeout", func() {
		// We will send a message to the host that will timeout after 3 seconds.
		// You cannot use 0 seconds because block timestamp will be greater than the timeout timestamp which is not allowed.
		// Host will not be able to respond to this message in time.

		// Stop the relayer so that the host cannot respond to the message:
		err := s.Relayer.StopRelayer(ctx, s.ExecRep)
		s.Require().NoError(err)

		time.Sleep(5 * time.Second)

		timeout := uint64(3)
		// Execute the contract:
		err = s.Contract.ExecCustomIcaMessages(ctx, wasmdUser.KeyName(), []proto.Message{}, icatypes.EncodingProto3JSON, nil, &timeout)
		s.Require().NoError(err)

		// Wait until timeout:
		err = testutil.WaitForBlocks(ctx, 5, wasmd, simd)
		s.Require().NoError(err)

		err = s.Relayer.StartRelayer(ctx, s.ExecRep)
		s.Require().NoError(err)

		// Wait until timeout acknoledgement is received:
		err = testutil.WaitForBlocks(ctx, 2, wasmd, simd)
		s.Require().NoError(err)

		// Flush to make sure the channel is closed in simd:
		err = s.Relayer.Flush(ctx, s.ExecRep, s.PathName, contractState.IcaInfo.ChannelID)
		s.Require().NoError(err)

		err = testutil.WaitForBlocks(ctx, 2, wasmd, simd)
		s.Require().NoError(err)

		// Check if channel was closed:
		wasmdChannels, err := s.Relayer.GetChannels(ctx, s.ExecRep, wasmd.Config().ChainID)
		s.Require().NoError(err)
		s.Require().Equal(1, len(wasmdChannels))
		s.Require().Equal(channeltypes.CLOSED.String(), wasmdChannels[0].State)

		simdChannels, err := s.Relayer.GetChannels(ctx, s.ExecRep, simd.Config().ChainID)
		s.Require().NoError(err)
		// sometimes there is a redundant channel for unknown reasons
		simdChannelsLen = len(simdChannels)
		s.Require().Greater(simdChannelsLen, 0)
		s.Require().Equal(channeltypes.CLOSED.String(), simdChannels[0].State)

		// Check if contract callbacks were executed:
		callbackCounter, err := s.Contract.QueryCallbackCounter(ctx)
		s.Require().NoError(err)
		s.Require().Equal(uint64(0), callbackCounter.Success)
		s.Require().Equal(uint64(0), callbackCounter.Error)
		s.Require().Equal(uint64(1), callbackCounter.Timeout)

		// Check if contract channel state was updated:
		contractChannelState, err := s.Contract.QueryChannelState(ctx)
		s.Require().NoError(err)
		s.Require().Equal(channeltypes.CLOSED.String(), contractChannelState.ChannelStatus)
	})

	s.Run("TestChannelReopening", func() {
		// Reopen the channel:
		txEncoding := icatypes.EncodingProto3JSON
		err := s.Contract.ExecCreateChannel(ctx, wasmdUser.KeyName(), s.ChainAConnID, s.ChainBConnID, nil, &txEncoding, "--gas", "500000")
		s.Require().NoError(err)

		// Wait for the channel to get set up
		err = testutil.WaitForBlocks(ctx, 10, s.ChainA, s.ChainB)
		s.Require().NoError(err)

		// Check if a new channel was opened in simd
		simdChannels, err := s.Relayer.GetChannels(ctx, s.ExecRep, simd.Config().ChainID)
		s.Require().NoError(err)
		// An extra channel may be created in simd for unknown reasons.
		s.Require().Greater(len(simdChannels), simdChannelsLen)
		s.Require().Equal(channeltypes.OPEN.String(), simdChannels[simdChannelsLen].State)
		simdChannelsLen = len(simdChannels)

		// Check if a new channel was opened in wasmd:
		wasmdChannels, err := s.Relayer.GetChannels(ctx, s.ExecRep, wasmd.Config().ChainID)
		s.Require().NoError(err)
		s.Require().Equal(2, len(wasmdChannels))
		wasmdChannel := wasmdChannels[1]
		s.Require().Equal(channeltypes.OPEN.String(), wasmdChannel.State)

		// Check if contract channel state was updated:
		contractChannelState, err := s.Contract.QueryChannelState(ctx)
		s.Require().NoError(err)
		s.Require().Equal(channeltypes.OPEN.String(), contractChannelState.ChannelStatus)
		// The version string is wrapped by the fee middleware. We we cannot check it directly here.
		// s.Require().Equal(wasmdChannel.Version, contractChannelState.Channel.Version)
		s.Require().Equal(wasmdChannel.ConnectionHops[0], contractChannelState.Channel.ConnectionID)
		s.Require().Equal(wasmdChannel.ChannelID, contractChannelState.Channel.Endpoint.ChannelID)
		s.Require().Equal(wasmdChannel.PortID, contractChannelState.Channel.Endpoint.PortID)
		s.Require().Equal(wasmdChannel.Counterparty.ChannelID, contractChannelState.Channel.CounterpartyEndpoint.ChannelID)
		s.Require().Equal(wasmdChannel.Counterparty.PortID, contractChannelState.Channel.CounterpartyEndpoint.PortID)
		s.Require().Equal(wasmdChannel.Ordering, contractChannelState.Channel.Order)

		contractState, err := s.Contract.QueryContractState(ctx)
		s.Require().NoError(err)
		s.Require().Equal(wasmdChannel.ChannelID, contractState.IcaInfo.ChannelID)
		s.Require().Equal(s.IcaAddress, contractState.IcaInfo.IcaAddress)

		callbackCounter, err := s.Contract.QueryCallbackCounter(ctx)
		s.Require().NoError(err)

		s.Require().Equal(uint64(0), callbackCounter.Success)
		s.Require().Equal(uint64(0), callbackCounter.Error)
		s.Require().Equal(uint64(1), callbackCounter.Timeout)
	})

	s.Run("TestSendCustomIcaMessagesAfterReopen", func() {
		// Send custom ICA message through the contract:
		sendMsg := &banktypes.MsgSend{
			FromAddress: s.IcaAddress,
			ToAddress:   s.UserB.FormattedAddress(),
			Amount:      sdk.NewCoins(sdk.NewCoin(simd.Config().Denom, sdkmath.NewInt(100))),
		}

		// Execute the contract:
		err = s.Contract.ExecCustomIcaMessages(ctx, wasmdUser.KeyName(), []proto.Message{sendMsg}, icatypes.EncodingProto3JSON, nil, nil)
		s.Require().NoError(err)

		err = testutil.WaitForBlocks(ctx, 10, wasmd, simd)
		s.Require().NoError(err)

		icaBalance, err := simd.GetBalance(ctx, s.IcaAddress, simd.Config().Denom)
		s.Require().NoError(err)
		s.Require().Equal(sdkmath.NewInt(1000000000-100), icaBalance)

		// Check if contract callbacks were executed:
		callbackCounter, err := s.Contract.QueryCallbackCounter(ctx)
		s.Require().NoError(err)

		s.Require().Equal(uint64(1), callbackCounter.Success)
		s.Require().Equal(uint64(0), callbackCounter.Error)
		s.Require().Equal(uint64(1), callbackCounter.Timeout)
	})
}

// toJSONString returns a string representation of the given value
// by marshaling it to JSON. It panics if marshaling fails.
func toJSONString(v any) string {
	bz, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bz)
}
