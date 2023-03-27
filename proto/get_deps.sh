#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"

TM_VERSION="rollkit/v0.37.0"
TM_PROTO_URL=https://raw.githubusercontent.com/rollkit/cometbft/$TM_VERSION/proto/tendermint

TM_PROTO_FILES=(
  abci/types.proto
  version/types.proto
  types/types.proto
  types/evidence.proto
  types/params.proto
  types/validator.proto
  state/types.proto
  crypto/proof.proto
  crypto/keys.proto
  libs/bits/types.proto
  p2p/types.proto
)

echo Fetching protobuf dependencies from Tendermint $TM_VERSION
for FILE in "${TM_PROTO_FILES[@]}"; do
  echo Fetching "$FILE"
  mkdir -p "tendermint/$(dirname $FILE)"
  curl -sSL "$TM_PROTO_URL/$FILE" > "tendermint/$FILE"
done
