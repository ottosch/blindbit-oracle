package common

var (
	MongoDBURI  = "mongodb://root:example@127.0.0.1:27017/"
	RpcEndpoint = "http://umbrel.local:8332/"
	RpcUser     = ""
	RpcPass     = ""
)

// control vars

var (
	CatchUp uint32 = 833_155 // random block where BIP352 was not merged yet. todo change to actual number
	// MinHeightToProcess No block below this number will be processed
	// todo is this actually needed
	//MinHeightToProcess uint32 = 833_000

	// SyncHeadersMaxPerCall how many headers will maximally be requested in one batched RPC call
	SyncHeadersMaxPerCall uint32 = 2_000
	MaxParallelRequests   uint8  = 6
)

// NumsH = 0x50929b74c1a04954b78b4b6035e97a5e078a5a0f28ec96d547bfee9ace803ac0
var NumsH = []byte{80, 146, 155, 116, 193, 160, 73, 84, 183, 139, 75, 96, 53, 233, 122, 94, 7, 138, 90, 15, 40, 236, 150, 213, 71, 191, 238, 154, 206, 128, 58, 192}
