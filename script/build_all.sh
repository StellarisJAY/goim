#!/usr/bin/env bash

echo "building API server..."
go build -mod=mod -o ./target/api ./cmd/api
echo "building Gateway server..."
go build -mod=mod -o ./target/gateway ./cmd/gateway
echo "building Friend RPC Service..."
go build -mod=mod -o ./target/rpc_friend ./cmd/rpc/friend
echo "building Group RPC Service..."
go build -mod=mod -o ./target/rpc_group ./cmd/rpc/group
echo "building Message RPC Service..."
go build -mod=mod -o ./target/rpc_message ./cmd/rpc/message
echo "building User RPC Service..."
go build -mod=mod -o ./target/rpc_user ./cmd/rpc/user
echo "building Transfer Service..."
go build -mod=mod -o ./target/transfer ./cmd/transfer

echo "goim build complete"