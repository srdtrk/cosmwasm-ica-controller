package icacontroller

import "encoding/json"

// ContractState is used to represent its state in Contract's storage
type ContractState struct {
	IcaInfo              IcaInfo `json:"ica_info"`
	AllowChannelOpenInit bool    `json:"allow_channel_open_init"`
}

// IcaInfo is used to represent the ICA info in the contract's state
type IcaInfo struct {
	IcaAddress string `json:"ica_address"`
	ChannelID  string `json:"channel_id"`
}

// CallbackCounter is used to represent the callback counter in the contract's storage
type CallbackCounter struct {
	Success uint64 `json:"success"`
	Error   uint64 `json:"error"`
	Timeout uint64 `json:"timeout"`
}

// ContractChannelState is used to represent the channel state in the contract's storage
type ContractChannelState struct {
	Channel       CwIbcChannel `json:"channel"`
	ChannelStatus string       `json:"channel_status"`
}

// OwnershipResponse is the response type for the OwnershipQueryMsg
type OwnershipResponse struct {
	// The current owner of the contract.
	// This contract must have an owner.
	Owner string `json:"owner"`
	// The pending owner of the contract if one exists.
	PendingOwner *string `json:"pending_owner"`
	// The height at which the pending owner offer expires.
	// Not sure how to represent this, so we'll just use a raw message
	PendingExpiry *json.RawMessage `json:"pending_expiry"`
}

// IsOpen returns true if the channel is open
func (c *ContractChannelState) IsOpen() bool {
	return c.ChannelStatus == "STATE_OPEN"
}

// CwIbcEndpoint is the endpoint of a channel defined in CosmWasm
type CwIbcEndpoint struct {
	PortID    string `json:"port_id"`
	ChannelID string `json:"channel_id"`
}

// CwIbcChannel is the channel defined in CosmWasm
type CwIbcChannel struct {
	Endpoint             CwIbcEndpoint `json:"endpoint"`
	CounterpartyEndpoint CwIbcEndpoint `json:"counterparty_endpoint"`
	// Order is either "ORDER_UNORDERED" or "ORDER_ORDERED"
	Order        string `json:"order"`
	Version      string `json:"version"`
	ConnectionID string `json:"connection_id"`
}
