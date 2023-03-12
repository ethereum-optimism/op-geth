package vm

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

const (
	pubKey    string = "03b434054a968479e6d1adb7b6185d1373c5b8f9cdd0813028327e6a342d702df6"
	dataHash  string = "989a647219cb0c3de61ec045ea197b8c48e8e40bc3fda8b93033b96b109a222a"
	signature string = "3045022074bcdc20e53b9342b8dad74aa65dfc8c0b80c3963f596440452316c762fa4b81022100f693c4b11ca20c3cffe96a4d0f404fcfeaf4b954d3d4dd162fe156381de654a6"
)

var input []byte

func init() {
	in := make([]byte, 256)

	pkh, _ := hex.DecodeString(pubKey)
	dhh, _ := hex.DecodeString(dataHash)
	sh, _ := hex.DecodeString(signature)
	sh = common.LeftPadBytes(sh, 72)

	copy(in[0:34], pkh)
	copy(in[34:66], dhh)
	copy(in[66:138], sh)

	input = in
}

func TestRunEcverify(t *testing.T) {
	ecv := &ecverify{}
	ok, _ := ecv.Run(input)
	assert.Equal(t, []byte{1}, ok)
}
