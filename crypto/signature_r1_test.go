package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	pubKey    string = "03b434054a968479e6d1adb7b6185d1373c5b8f9cdd0813028327e6a342d702df6"
	dataHash  string = "989a647219cb0c3de61ec045ea197b8c48e8e40bc3fda8b93033b96b109a222a"
	signature string = "3045022074bcdc20e53b9342b8dad74aa65dfc8c0b80c3963f596440452316c762fa4b81022100f693c4b11ca20c3cffe96a4d0f404fcfeaf4b954d3d4dd162fe156381de654a6"
)

func TestVerifySignatureR1(t *testing.T) {
	pkh, _ := hex.DecodeString(pubKey)
	dhh, _ := hex.DecodeString(dataHash)
	sh, _ := hex.DecodeString(signature)

	ok := VerifySignatureR1(pkh, dhh, sh)
	assert.Equal(t, []byte{1}, ok)

	ok = VerifySignatureR1(pkh, sh, dhh)
	assert.Equal(t, []byte{0}, ok)
}
