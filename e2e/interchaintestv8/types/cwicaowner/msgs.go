/* Code generated by github.com/srdtrk/go-codegen, DO NOT EDIT. */
package cwicaowner

type InstantiateMsg struct {
	Admin *string `json:"admin,omitempty"`
	IcaControllerCodeId int `json:"ica_controller_code_id"`
}

type ExecuteMsg struct {
	CreateIcaContract *ExecuteMsg_CreateIcaContract `json:"create_ica_contract,omitempty"`
	/*
	   SendPredefinedAction sends a predefined action from the ICA controller to the ICA host. This demonstration is useful for contracts that have predefined actions such as DAOs.

	   In this example, the predefined action is a `MsgSend` message which sends 100 "stake" tokens.
	*/
	SendPredefinedAction *ExecuteMsg_SendPredefinedAction `json:"send_predefined_action,omitempty"`
	// The callback message from `cw-ica-controller`. The handler for this variant should verify that this message comes from an expected legitimate source.
	ReceiveIcaCallback *ExecuteMsg_ReceiveIcaCallback `json:"receive_ica_callback,omitempty"`
}

type QueryMsg struct {
	// GetContractState returns the contact's state.
	GetContractState *QueryMsg_GetContractState `json:"get_contract_state,omitempty"`
	// GetIcaState returns the ICA state for the given ICA ID.
	GetIcaContractState *QueryMsg_GetIcaContractState `json:"get_ica_contract_state,omitempty"`
	// GetIcaCount returns the number of ICAs.
	GetIcaCount *QueryMsg_GetIcaCount `json:"get_ica_count,omitempty"`
}
type ExecuteMsg_ReceiveIcaCallback IcaControllerCallbackMsg

/*
A thin wrapper around u64 that is using strings for JSON encoding/decoding, such that the full u64 range can be used for clients that convert JSON numbers to floats, like JavaScript and jq.

# Examples

Use `from` to create instances of this and `u64` to get the value out:

``` # use cosmwasm_std::Uint64; let a = Uint64::from(42u64); assert_eq!(a.u64(), 42);

let b = Uint64::from(70u32); assert_eq!(b.u64(), 70); ```
*/
type Uint64 string

// IcaControllerCallbackMsg is the type of message that this contract can send to other contracts.
type IcaControllerCallbackMsg struct {
	// OnAcknowledgementPacketCallback is the callback that this contract makes to other contracts when it receives an acknowledgement packet.
	OnAcknowledgementPacketCallback *IcaControllerCallbackMsg_OnAcknowledgementPacketCallback `json:"on_acknowledgement_packet_callback,omitempty"`
	// OnTimeoutPacketCallback is the callback that this contract makes to other contracts when it receives a timeout packet.
	OnTimeoutPacketCallback *IcaControllerCallbackMsg_OnTimeoutPacketCallback `json:"on_timeout_packet_callback,omitempty"`
	// OnChannelOpenAckCallback is the callback that this contract makes to other contracts when it receives a channel open acknowledgement.
	OnChannelOpenAckCallback *IcaControllerCallbackMsg_OnChannelOpenAckCallback `json:"on_channel_open_ack_callback,omitempty"`
}

// In IBC each package must set at least one type of timeout: the timestamp or the block height. Using this rather complex enum instead of two timeout fields we ensure that at least one timeout is set.
type IbcTimeout struct {
	Block *IbcTimeoutBlock `json:"block,omitempty"`
	Timestamp *Timestamp `json:"timestamp,omitempty"`
}

type QueryMsg_GetIcaCount struct{}

// IcaContractState is the state of the cw-ica-controller contract.
type IcaContractState struct {
	ContractAddr Addr `json:"contract_addr"`
	IcaState *IcaState `json:"ica_state,omitempty"`
}

type ExecuteMsg_CreateIcaContract struct {
	ChannelOpenInitOptions ChannelOpenInitOptions `json:"channel_open_init_options"`
	Salt *string `json:"salt,omitempty"`
}

// IbcOrder defines if a channel is ORDERED or UNORDERED Values come from https://github.com/cosmos/cosmos-sdk/blob/v0.40.0/proto/ibc/core/channel/v1/channel.proto#L69-L80 Naming comes from the protobuf files and go translations.
type IbcOrder string

const (
	IbcOrder_OrderUnordered IbcOrder = "ORDER_UNORDERED"
	IbcOrder_OrderOrdered   IbcOrder = "ORDER_ORDERED"
)

type QueryMsg_GetContractState struct{}

// ContractState is the state of the IBC application.
type ContractState struct {
	// The admin of this contract.
	Admin Addr `json:"admin"`
	// The code ID of the cw-ica-controller contract.
	IcaControllerCodeId int `json:"ica_controller_code_id"`
}

// State is the state of the IBC application's channel. This application only supports one channel.
type State struct {
	// The IBC channel, as defined by cosmwasm.
	Channel IbcChannel `json:"channel"`
	// The status of the channel.
	ChannelStatus Status `json:"channel_status"`
}

type IbcPacket struct {
	// The raw data sent from the other side in the packet
	Data Binary `json:"data"`
	// identifies the channel and port on the receiving chain.
	Dest IbcEndpoint `json:"dest"`
	// The sequence number of the packet on the given channel
	Sequence int `json:"sequence"`
	// identifies the channel and port on the sending chain.
	Src IbcEndpoint `json:"src"`
	Timeout IbcTimeout `json:"timeout"`
}

// IBCTimeoutHeight Height is a monotonically increasing data type that can be compared against another Height for the purposes of updating and freezing clients. Ordering is (revision_number, timeout_height)
type IbcTimeoutBlock struct {
	// block height after which the packet times out. the height within the given revision
	Height int `json:"height"`
	// the version that the client is currently on (e.g. after resetting the chain this could increment 1 as height drops to 0)
	Revision int `json:"revision"`
}

// `TxEncoding` is the encoding of the transactions sent to the ICA host.
type TxEncoding string

const (
	// `Protobuf` is the protobuf serialization of the CosmosSDK's Any.
	TxEncoding_Proto3 TxEncoding = "proto3"
	// `Proto3Json` is the json serialization of the CosmosSDK's Any.
	TxEncoding_Proto3Json TxEncoding = "proto3json"
)

// IbcChannel defines all information on a channel. This is generally used in the hand-shake process, but can be queried directly.
type IbcChannel struct {
	// The connection upon which this channel was created. If this is a multi-hop channel, we only expose the first hop.
	ConnectionId string `json:"connection_id"`
	CounterpartyEndpoint IbcEndpoint `json:"counterparty_endpoint"`
	Endpoint IbcEndpoint `json:"endpoint"`
	Order IbcOrder `json:"order"`
	// Note: in ibcv3 this may be "", in the IbcOpenChannel handshake messages
	Version string `json:"version"`
}

/*
A human readable address.

In Cosmos, this is typically bech32 encoded. But for multi-chain smart contracts no assumptions should be made other than being UTF-8 encoded and of reasonable length.

This type represents a validated address. It can be created in the following ways 1. Use `Addr::unchecked(input)` 2. Use `let checked: Addr = deps.api.addr_validate(input)?` 3. Use `let checked: Addr = deps.api.addr_humanize(canonical_addr)?` 4. Deserialize from JSON. This must only be done from JSON that was validated before such as a contract's state. `Addr` must not be used in messages sent by the user because this would result in unvalidated instances.

This type is immutable. If you really need to mutate it (Really? Are you sure?), create a mutable copy using `let mut mutable = Addr::to_string()` and operate on that `String` instance.
*/
type Addr string

/*
A point in time in nanosecond precision.

This type can represent times from 1970-01-01T00:00:00Z to 2554-07-21T23:34:33Z.

## Examples

``` # use cosmwasm_std::Timestamp; let ts = Timestamp::from_nanos(1_000_000_202); assert_eq!(ts.nanos(), 1_000_000_202); assert_eq!(ts.seconds(), 1); assert_eq!(ts.subsec_nanos(), 202);

let ts = ts.plus_seconds(2); assert_eq!(ts.nanos(), 3_000_000_202); assert_eq!(ts.seconds(), 3); assert_eq!(ts.subsec_nanos(), 202); ```
*/
type Timestamp Uint64

type QueryMsg_GetIcaContractState struct {
	IcaId int `json:"ica_id"`
}

// Status is the status of an IBC channel.
type Status string

const (
	// Uninitialized is the default state of the channel.
	Status_StateUninitializedUnspecified Status = "STATE_UNINITIALIZED_UNSPECIFIED"
	// Init is the state of the channel when it is created.
	Status_StateInit Status = "STATE_INIT"
	// TryOpen is the state of the channel when it is trying to open.
	Status_StateTryopen Status = "STATE_TRYOPEN"
	// Open is the state of the channel when it is open.
	Status_StateOpen Status = "STATE_OPEN"
	// Closed is the state of the channel when it is closed.
	Status_StateClosed Status = "STATE_CLOSED"
	// The channel has just accepted the upgrade handshake attempt and is flushing in-flight packets. Added in `ibc-go` v8.1.0.
	Status_StateFlushing Status = "STATE_FLUSHING"
	// The channel has just completed flushing any in-flight packets. Added in `ibc-go` v8.1.0.
	Status_StateFlushcomplete Status = "STATE_FLUSHCOMPLETE"
)

type ExecuteMsg_SendPredefinedAction struct {
	// The ICA ID.
	IcaId int `json:"ica_id"`
	// The recipient's address, on the counterparty chain, to send the tokens to from ICA host.
	ToAddress string `json:"to_address"`
}

type IbcEndpoint struct {
	ChannelId string `json:"channel_id"`
	PortId string `json:"port_id"`
}

// `Data` is the response to an ibc packet. It either contains a result or an error.
type Data struct {
	// Result is the result of a successful transaction.
	Result *Data_Result `json:"result,omitempty"`
	// Error is the error message of a failed transaction. It is a string of the error message (not base64 encoded).
	Error *Data_Error `json:"error,omitempty"`
}

/*
Binary is a wrapper around Vec<u8> to add base64 de/serialization with serde. It also adds some helper methods to help encode inline.

This is only needed as serde-json-{core,wasm} has a horrible encoding for Vec<u8>. See also <https://github.com/CosmWasm/cosmwasm/blob/main/docs/MESSAGE_TYPES.md>.
*/
type Binary string

// The options needed to initialize the IBC channel.
type ChannelOpenInitOptions struct {
	// The order of the channel. If not specified, [`IbcOrder::Ordered`] is used. [`IbcOrder::Unordered`] is only supported if the counterparty chain is using `ibc-go` v8.1.0 or later.
	ChannelOrdering *IbcOrder `json:"channel_ordering,omitempty"`
	// The connection id on this chain.
	ConnectionId string `json:"connection_id"`
	// The counterparty connection id on the counterparty chain.
	CounterpartyConnectionId string `json:"counterparty_connection_id"`
	// The counterparty port id. If not specified, [`crate::ibc::types::keys::HOST_PORT_ID`] is used. Currently, this contract only supports the host port.
	CounterpartyPortId *string `json:"counterparty_port_id,omitempty"`
}

// IcaState is the state of the ICA.
type IcaState struct {
	ChannelState State `json:"channel_state"`
	IcaAddr string `json:"ica_addr"`
	IcaId int `json:"ica_id"`
	TxEncoding TxEncoding `json:"tx_encoding"`
}

type IcaControllerCallbackMsg_OnAcknowledgementPacketCallback struct {
	// The deserialized ICA acknowledgement data
	IcaAcknowledgement Data `json:"ica_acknowledgement"`
	// The original packet that was sent
	OriginalPacket IbcPacket `json:"original_packet"`
	// The relayer that submitted acknowledgement packet
	Relayer Addr `json:"relayer"`
}

type IcaControllerCallbackMsg_OnTimeoutPacketCallback struct {
	// The original packet that was sent
	OriginalPacket IbcPacket `json:"original_packet"`
	// The relayer that submitted acknowledgement packet
	Relayer Addr `json:"relayer"`
}

type IcaControllerCallbackMsg_OnChannelOpenAckCallback struct {
	// The channel that was opened.
	Channel IbcChannel `json:"channel"`
	// The address of the interchain account that was created.
	IcaAddress string `json:"ica_address"`
	// The tx encoding this ICA channel uses.
	TxEncoding TxEncoding `json:"tx_encoding"`
}
type Data_Result Binary

type Data_Error string
