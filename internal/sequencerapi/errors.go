package sequencerapi

import (
	"github.com/ethereum/go-ethereum/rpc"
)

var _ rpc.Error = new(jsonRpcError)

type jsonRpcError struct {
	message string
	code    int
}

func (j *jsonRpcError) Error() string {
	return j.message
}

func (j *jsonRpcError) ErrorCode() int {
	return j.code
}
