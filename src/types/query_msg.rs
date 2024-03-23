//! This module contains the helpers to convert [`QueryRequest`] to protobuf bytes and vice versa.

use cosmwasm_std::{QueryRequest, Empty};

/// `MsgModuleQuerySafe` defines the query request tx added in ibc-go v8.2
#[derive(::prost::Message)]
pub struct MsgModuleQuerySafe {
    #[prost(string, tag = "1")]
    /// signer is the address of the account that signed the transaction
    pub signer: ::prost::alloc::string::String,
    /// requests is the list of query requests
    #[prost(message, repeated, tag = "2")]
    pub requests: ::prost::alloc::vec::Vec<AbciQueryRequest>,
}

/// `AbciQueryRequest` defines the parameters for a particular query request by an interchain account.
#[derive(::prost::Message)]
pub struct AbciQueryRequest {
    #[prost(string, tag = "1")]
    /// `path` defines the path of the query request as defined by ADR-021.
    /// https://github.com/cosmos/cosmos-sdk/blob/main/docs/architecture/adr-021-protobuf-query-encoding.md#custom-query-registration-and-routing
    pub path: ::prost::alloc::string::String,
    #[prost(bytes = "vec", tag = "2")]
    /// `data` defines the payload of the query request as defined by ADR-021.
    /// https://github.com/cosmos/cosmos-sdk/blob/main/docs/architecture/adr-021-protobuf-query-encoding.md#custom-query-registration-and-routing
    pub data: ::prost::alloc::vec::Vec<u8>,
}

/// `MsgModuleQuerySafeResponse` defines the response for Msg/ModuleQuerySafe
#[derive(::prost::Message)]
pub struct MsgModuleQuerySafeResponse {
    /// responses is the list of query responses as bytes
    /// The responses are in the same order as the requests
    #[prost(bytes = "vec", repeated, tag = "1")]
    pub responses: ::prost::alloc::vec::Vec<Vec<u8>>,
}

/// Converts a [`QueryRequest`] to a gRPC method path and protobuf bytes.
///
/// # Panics
///
/// Panics if the query type is not supported.
#[must_use]
pub fn query_to_protobuf(query: QueryRequest<Empty>) -> (String, Vec<u8>) {
    match query {
        QueryRequest::Bank(bank_query) => convert_to_protobuf::bank(bank_query),
        QueryRequest::Stargate { path, data } => (path, data.0),
        QueryRequest::Wasm(_) => panic!("wasmd queries are not marked module safe (yet)"),
        QueryRequest::Ibc(_) => panic!("ibc-go queries are not marked module safe (yet)"),
        QueryRequest::Custom(_) => panic!("custom queries are not supported"),
        #[cfg(feature = "staking")]
        QueryRequest::Staking(_) => todo!(),
        #[cfg(feature = "staking")]
        QueryRequest::Distribution(_) => todo!(),
        _ => panic!("Unsupported QueryRequest"),
    }
}

mod convert_to_protobuf {
    use cosmos_sdk_proto::{
        prost::Message,
        cosmos::bank::v1beta1::{QueryAllBalancesRequest, QueryBalanceRequest, QueryDenomMetadataRequest, QueryDenomsMetadataRequest},
        cosmos::{bank::v1beta1::QuerySupplyOfRequest, base::query::v1beta1::PageRequest},
    };
    use cosmwasm_std::BankQuery;

    pub fn bank(bank_query: BankQuery) -> (String, Vec<u8>) {
        match bank_query {
            BankQuery::Balance { address, denom } => (
                "/cosmos.bank.v1beta1.Query/Balance".to_string(),
                QueryBalanceRequest {
                    address,
                    denom,
                }.encode_to_vec()
            ),
            BankQuery::AllBalances { address } => (
                "/cosmos.bank.v1beta1.Query/AllBalances".to_string(),
                QueryAllBalancesRequest {
                    address,
                    pagination: None,
                }.encode_to_vec()
            ),
            BankQuery::DenomMetadata { denom } => (
                "/cosmos.bank.v1beta1.Query/DenomMetadata".to_string(),
                QueryDenomMetadataRequest {
                    denom,
                }.encode_to_vec()
            ),
            BankQuery::AllDenomMetadata { pagination } => {
                let pagination = pagination.map(|pagination| PageRequest {
                    key: pagination.key.unwrap_or_default().0,
                    limit: u64::from(pagination.limit),
                    reverse: pagination.reverse,
                    offset: 0,
                    count_total: false,
                });

                (
                    "/cosmos.bank.v1beta1.Query/AllDenomMetadata".to_string(),
                    QueryDenomsMetadataRequest {
                        pagination,
                    }.encode_to_vec(),
                )
            },
            BankQuery::Supply { denom } => (
                "/cosmos.bank.v1beta1.Query/Supply".to_string(),
                QuerySupplyOfRequest {
                    denom,
                }.encode_to_vec()
            ),
            _ => panic!("Unsupported BankQuery"),
        }
    }
}
