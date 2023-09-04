package precompiles

/*
#cgo LDFLAGS: -L../target/release/ -lprecompiles -lm


#include <stdbool.h>
#include <stdint.h>

uint64_t __precompile_anonymous_verify_gas(const void* data_ptr, const uint32_t data_len);
uint8_t __precompile_anonymous_verify(const void* data_ptr, const uint32_t data_len);

uint64_t __precompile_anemoi_gas(const void* data_ptr, const uint32_t data_len);
uint8_t __precompile_anemoi(const void* data_ptr, const uint32_t data_len, void* ret_val);

uint64_t __precompile_mental_pokey_verify_gas(const void* data_ptr, const uint32_t data_len);
uint8_t __precompile_mental_pokey_verify(const void* data_ptr, const uint32_t data_len);

uint64_t __precompile_mental_pokey_exec_gas(const void* data_ptr, const uint32_t data_len);
uint8_t __precompile_mental_pokey_exec(const void* data_ptr, const uint32_t data_len, void* ret_val);

*/
import "C"
import "unsafe"

type Anonymous struct{}

func (a *Anonymous) RequiredGas(input []byte) uint64 {
	cstr := unsafe.Pointer(&input[0])
	len := C.uint(len(input))

	gas := C.__precompile_anonymous_verify_gas(cstr, len)

	return uint64(gas)
}

func (a *Anonymous) Run(input []byte) ([]byte, error) {
	cstr := unsafe.Pointer(&input[0])
	len := C.uint(len(input))

	res := C.__precompile_anonymous_verify(cstr, len)

	output := make([]byte, 32)

	output[31] = byte(res)

	return output, nil
}

type Anemoi struct{}

func (a *Anemoi) RequiredGas(input []byte) uint64 {
	cstr := unsafe.Pointer(&input[0])
	len := C.uint(len(input))

	gas := C.__precompile_anemoi_gas(cstr, len)

	return uint64(gas)
}

func (a *Anemoi) Run(input []byte) ([]byte, error) {
	output := make([]byte, 64)
	cout := unsafe.Pointer(&output[0])

	cstr := unsafe.Pointer(&input[0])
	len := C.uint(len(input))

	res := C.__precompile_anemoi(cstr, len, cout)

	output[63] = byte(res)

	return output, nil
}

type MentalPokeyVerify struct{}

func (a *MentalPokeyVerify) RequiredGas(input []byte) uint64 {
	cstr := unsafe.Pointer(&input[0])
	len := C.uint(len(input))

	gas := C.__precompile_mental_pokey_verify_gas(cstr, len)

	return uint64(gas)
}

func (a *MentalPokeyVerify) Run(input []byte) ([]byte, error) {

	cstr := unsafe.Pointer(&input[0])
	len := C.uint(len(input))

	res := C.__precompile_mental_pokey_verify(cstr, len)

	output := make([]byte, 32)

	output[31] = byte(res)

	return output, nil
}

type MentalPokeyExec struct{}

func (a *MentalPokeyExec) RequiredGas(input []byte) uint64 {
	cstr := unsafe.Pointer(&input[0])
	len := C.uint(len(input))

	gas := C.__precompile_mental_pokey_exec_gas(cstr, len)

	return uint64(gas)
}

func (a *MentalPokeyExec) Run(input []byte) ([]byte, error) {
	output := make([]byte, 64)
	cout := unsafe.Pointer(&output[0])

	cstr := unsafe.Pointer(&input[0])
	len := C.uint(len(input))

	res := C.__precompile_mental_pokey_exec(cstr, len, cout)

	output[63] = byte(res)

	return output, nil
}
