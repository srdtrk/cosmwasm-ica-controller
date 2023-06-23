package types

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewInstantiateMsg creates a new InstantiateMsg.
func NewInstantiateMsg(admin *string) string {
	if admin == nil {
		return `{}`
	} else {
		return fmt.Sprintf(`{"admin":"%s"}`, *admin)
	}
}

// NewSendPredefinedActionMsg creates a new SendPredefinedActionMsg.
func NewSendPredefinedActionMsg(to_address string) string {
	return fmt.Sprintf(`{"send_predefined_action":{"to_address":"%s"}}`, to_address)
}

// NewSendCustomIcaMessagesMsg creates a new SendCustomIcaMessagesMsg.
func NewSendCustomIcaMessagesMsg(cdc codec.BinaryCodec, msgs []sdk.Msg, memo *string, timeout *uint64) string {
	type SendCustomIcaMessagesMsg struct {
		Messages       []string `json:"messages"`
		PacketMemo     *string  `json:"packet_memo,omitempty"`
		TimeoutSeconds *uint64  `json:"timeout_seconds,omitempty"`
	}

	type SendCustomIcaMessagesMsgWrapper struct {
		SendCustomIcaMessagesMsg SendCustomIcaMessagesMsg `json:"send_custom_ica_messages"`
	}

	messages := make([]string, len(msgs))

	for _, msg := range msgs {
		bz, err := cdc.(*codec.ProtoCodec).MarshalJSON(msg)
		if err != nil {
			panic(err)
		}
		b64 := base64.StdEncoding.EncodeToString(bz)
		messages = append(messages, b64)
	}

	msg := SendCustomIcaMessagesMsgWrapper{
		SendCustomIcaMessagesMsg: SendCustomIcaMessagesMsg{
			Messages:       messages,
			PacketMemo:     memo,
			TimeoutSeconds: timeout,
		},
	}

	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return string(jsonBytes)
}
