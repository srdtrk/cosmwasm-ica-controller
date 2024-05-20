/* Code generated by github.com/srdtrk/go-codegen, DO NOT EDIT. */
package cwicacontroller

import (
	"context"
	"encoding/json"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	grpc "google.golang.org/grpc"
	insecure "google.golang.org/grpc/credentials/insecure"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// GetContractState is the client API for the QueryMsg_GetContractState query message
	GetContractState(ctx context.Context, req *QueryMsg_GetContractState, opts ...grpc.CallOption) (*State, error)
	// Ownership is the client API for the QueryMsg_Ownership query message
	Ownership(ctx context.Context, req *QueryMsg_Ownership, opts ...grpc.CallOption) (*Ownership_for_String, error)
	// GetChannel is the client API for the QueryMsg_GetChannel query message
	GetChannel(ctx context.Context, req *QueryMsg_GetChannel, opts ...grpc.CallOption) (*ChannelState, error)
}

type queryClient struct {
	cc      *grpc.ClientConn
	address string
}

var _ QueryClient = (*queryClient)(nil)

// NewQueryClient creates a new QueryClient
func NewQueryClient(gRPCAddress, contractAddress string, opts ...grpc.DialOption) (QueryClient, error) {
	if len(opts) == 0 {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Create a connection to the gRPC server
	grpcConn, err := grpc.Dial(gRPCAddress, opts...)
	if err != nil {
		return nil, err
	}

	return &queryClient{
		address: contractAddress,
		cc:      grpcConn,
	}, nil
}

// Close closes the gRPC connection to the server
func (q *queryClient) Close() error {
	return q.cc.Close()
}

// queryContract is a helper function to query the contract with raw query data
func (q *queryClient) queryContract(ctx context.Context, rawQueryData []byte, opts ...grpc.CallOption) ([]byte, error) {
	in := &wasmtypes.QuerySmartContractStateRequest{
		Address:   q.address,
		QueryData: rawQueryData,
	}
	out := new(wasmtypes.QuerySmartContractStateResponse)
	err := q.cc.Invoke(ctx, "/cosmwasm.wasm.v1.Query/SmartContractState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out.Data, nil
}

func (q *queryClient) Ownership(ctx context.Context, req *QueryMsg_Ownership, opts ...grpc.CallOption) (*Ownership_for_String, error) {
	rawQueryData, err := json.Marshal(&QueryMsg{Ownership: req})
	if err != nil {
		return nil, err
	}

	rawResponseData, err := q.queryContract(ctx, rawQueryData, opts...)
	if err != nil {
		return nil, err
	}

	var response Ownership_for_String
	if err := json.Unmarshal(rawResponseData, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (q *queryClient) GetChannel(ctx context.Context, req *QueryMsg_GetChannel, opts ...grpc.CallOption) (*ChannelState, error) {
	rawQueryData, err := json.Marshal(&QueryMsg{GetChannel: req})
	if err != nil {
		return nil, err
	}

	rawResponseData, err := q.queryContract(ctx, rawQueryData, opts...)
	if err != nil {
		return nil, err
	}

	var response ChannelState
	if err := json.Unmarshal(rawResponseData, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (q *queryClient) GetContractState(ctx context.Context, req *QueryMsg_GetContractState, opts ...grpc.CallOption) (*State, error) {
	rawQueryData, err := json.Marshal(&QueryMsg{GetContractState: req})
	if err != nil {
		return nil, err
	}

	rawResponseData, err := q.queryContract(ctx, rawQueryData, opts...)
	if err != nil {
		return nil, err
	}

	var response State
	if err := json.Unmarshal(rawResponseData, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
