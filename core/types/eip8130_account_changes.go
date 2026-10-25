// Copyright 2026 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
)

// EIP-8130 account_changes entry type bytes. On the wire each entry is a single
// flat RLP list whose first element is the type byte: rlp([type_byte, fields...]).
const (
	accountChangeTypeCreate     = 0x00
	accountChangeTypeConfig     = 0x01
	accountChangeTypeDelegation = 0x02
)

// ActorChangeType is the operation performed by an ActorChange. The value is the
// on-wire op byte; in RLP it encodes as a bare uint, in JSON as a string.
type ActorChangeType uint8

const (
	ActorChangeAuthorize ActorChangeType = 0x01
	ActorChangeRevoke    ActorChangeType = 0x02
)

func (t ActorChangeType) MarshalJSON() ([]byte, error) {
	switch t {
	case ActorChangeAuthorize:
		return []byte(`"Authorize"`), nil
	case ActorChangeRevoke:
		return []byte(`"Revoke"`), nil
	default:
		return nil, fmt.Errorf("eip8130: invalid actor change type %d", uint8(t))
	}
}

func (t *ActorChangeType) UnmarshalJSON(input []byte) error {
	switch string(input) {
	case `"Authorize"`:
		*t = ActorChangeAuthorize
	case `"Revoke"`:
		*t = ActorChangeRevoke
	default:
		return fmt.Errorf("eip8130: invalid actor change type %s", input)
	}
	return nil
}

// DecodeRLP reads the op byte and rejects any value other than Authorize (0x01)
// or Revoke (0x02), matching the Rust strict decode (ActorChangeType::from_op_byte).
// Without this, the derived uint8 decoder would accept any byte 0x00..0xff.
func (t *ActorChangeType) DecodeRLP(s *rlp.Stream) error {
	b, err := s.Uint8()
	if err != nil {
		return err
	}
	switch ActorChangeType(b) {
	case ActorChangeAuthorize, ActorChangeRevoke:
		*t = ActorChangeType(b)
		return nil
	default:
		return fmt.Errorf("eip8130: invalid actor change type byte 0x%x", b)
	}
}

// InitialActor is an actor installed on a newly-created account. Wire form is
// rlp([actorId, authenticator, scope, policyData]); scope is stored verbatim
// (0x00 = unrestricted admin) and policyData is empty unless scope sets the
// POLICY bit.
type InitialActor struct {
	ActorID       common.Hash    `json:"actorId"`
	Authenticator common.Address `json:"authenticator"`
	Scope         uint8          `json:"scope"`
	PolicyData    hexutil.Bytes  `json:"policyData"`
}

// ActorChange is a single actor authorize/revoke operation inside a ConfigChange.
// Wire form is rlp([changeType, actorId, data]); data is opaque at this layer.
type ActorChange struct {
	ChangeType ActorChangeType `json:"changeType"`
	ActorID    common.Hash     `json:"actorId"`
	Data       hexutil.Bytes   `json:"data"`
}

// CreateEntry is the body of an AccountChange create entry. Wire form is
// rlp([userSalt, code, [InitialActor, ...]]).
type CreateEntry struct {
	UserSalt      common.Hash    `json:"userSalt"`
	Code          hexutil.Bytes  `json:"code"`
	InitialActors []InitialActor `json:"initialActors"`
}

// ConfigChange is the body of an AccountChange config-change entry. Wire form is
// rlp([chainId, sequence, [ActorChange, ...], auth]).
type ConfigChange struct {
	ChainID      uint64        `json:"chainId"`
	Sequence     uint64        `json:"sequence"`
	ActorChanges []ActorChange `json:"actorChanges"`
	Auth         hexutil.Bytes `json:"auth"`
}

// Delegation is the body of an AccountChange delegation entry. Wire form is
// rlp([target]); a zero target clears the existing delegation.
type Delegation struct {
	Target common.Address `json:"target"`
}

// AccountChange is a tagged-union entry inside Eip8130Tx.AccountChanges. Exactly
// one of the body pointers is set. On the wire each entry is a single flat RLP
// list whose first element is the type byte, followed by the body fields inline:
// rlp([type_byte, body fields...]). The type byte is a genuine list element (not
// an EIP-2718-style type_byte || rlp(...) prefix), so each entry is one
// self-contained RLP item. In JSON it is the body object with an added "type"
// discriminator ("create" / "configChange" / "delegation").
type AccountChange struct {
	Create       *CreateEntry
	ConfigChange *ConfigChange
	Delegation   *Delegation
}

// resolveBody returns the wire type byte, the JSON discriminator and the body
// value of the single set body pointer. It enforces the tagged-union invariant
// that exactly one of Create/ConfigChange/Delegation is set, matching the Rust
// enum, which can only ever hold a single variant. The returned body has its
// nil nested slices normalized to empty slices so JSON marshalling emits [] (to
// match serde's Vec) rather than null; the original AccountChange is never
// mutated.
func (a AccountChange) resolveBody() (typeByte byte, typ string, body interface{}, err error) {
	n := 0
	if a.Create != nil {
		n++
		cpy := *a.Create
		if cpy.InitialActors == nil {
			cpy.InitialActors = []InitialActor{}
		}
		typeByte, typ, body = accountChangeTypeCreate, "create", &cpy
	}
	if a.ConfigChange != nil {
		n++
		cpy := *a.ConfigChange
		if cpy.ActorChanges == nil {
			cpy.ActorChanges = []ActorChange{}
		}
		typeByte, typ, body = accountChangeTypeConfig, "configChange", &cpy
	}
	if a.Delegation != nil {
		n++
		typeByte, typ, body = accountChangeTypeDelegation, "delegation", a.Delegation
	}
	if n != 1 {
		return 0, "", nil, errors.New("eip8130: account change must set exactly one body")
	}
	return typeByte, typ, body, nil
}

// EncodeRLP writes the entry as a single flat list rlp([type_byte, body
// fields...]). The type byte is an in-list element encoded as an RLP uint (e.g.
// 0x00 -> 0x80), followed by the set body's fields inline, mirroring the Rust
// encoder.
func (a AccountChange) EncodeRLP(w io.Writer) error {
	typeByte, _, body, err := a.resolveBody()
	if err != nil {
		return err
	}

	buf := rlp.NewEncoderBuffer(w)
	list := buf.List()
	buf.WriteUint64(uint64(typeByte))

	// The body's fields are written inline (no inner list header) so the whole
	// entry is one flat list.
	switch b := body.(type) {
	case *CreateEntry:
		buf.WriteBytes(b.UserSalt[:])
		buf.WriteBytes(b.Code)
		if err := rlp.Encode(buf, b.InitialActors); err != nil {
			return err
		}
	case *ConfigChange:
		buf.WriteUint64(b.ChainID)
		buf.WriteUint64(b.Sequence)
		if err := rlp.Encode(buf, b.ActorChanges); err != nil {
			return err
		}
		buf.WriteBytes(b.Auth)
	case *Delegation:
		buf.WriteBytes(b.Target[:])
	default:
		return fmt.Errorf("eip8130: unexpected account change body %T", body)
	}

	buf.ListEnd(list)
	return buf.Flush()
}

// DecodeRLP reads the single flat list, reads the type byte as an RLP uint,
// dispatches on it and decodes the body fields positionally, then enforces that
// the list is fully consumed. Rejects any unknown type byte. Returns rlp.EOL at
// the end of the enclosing list so it composes with slice decoding.
func (a *AccountChange) DecodeRLP(s *rlp.Stream) error {
	if _, err := s.List(); err != nil {
		return err
	}
	typeByte, err := s.Uint8()
	if err != nil {
		return err
	}

	switch typeByte {
	case accountChangeTypeCreate:
		body := new(CreateEntry)
		if err := s.Decode(&body.UserSalt); err != nil {
			return err
		}
		var code []byte
		if err := s.Decode(&code); err != nil {
			return err
		}
		body.Code = code
		if err := s.Decode(&body.InitialActors); err != nil {
			return err
		}
		a.Create = body
	case accountChangeTypeConfig:
		body := new(ConfigChange)
		if err := s.Decode(&body.ChainID); err != nil {
			return err
		}
		if err := s.Decode(&body.Sequence); err != nil {
			return err
		}
		if err := s.Decode(&body.ActorChanges); err != nil {
			return err
		}
		var auth []byte
		if err := s.Decode(&auth); err != nil {
			return err
		}
		body.Auth = auth
		a.ConfigChange = body
	case accountChangeTypeDelegation:
		body := new(Delegation)
		if err := s.Decode(&body.Target); err != nil {
			return err
		}
		a.Delegation = body
	default:
		return fmt.Errorf("eip8130: invalid account change type byte 0x%x", typeByte)
	}

	// Enforce that the list is fully consumed (no trailing elements).
	return s.ListEnd()
}

// MarshalJSON encodes the entry as its body object plus a "type" discriminator.
func (a AccountChange) MarshalJSON() ([]byte, error) {
	_, typ, body, err := a.resolveBody()
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	fields := map[string]json.RawMessage{}
	if err := json.Unmarshal(raw, &fields); err != nil {
		return nil, err
	}
	fields["type"], _ = json.Marshal(typ)
	return json.Marshal(fields)
}

// UnmarshalJSON dispatches on the "type" discriminator.
func (a *AccountChange) UnmarshalJSON(input []byte) error {
	var tag struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(input, &tag); err != nil {
		return err
	}
	switch tag.Type {
	case "create":
		a.Create = new(CreateEntry)
		return json.Unmarshal(input, a.Create)
	case "configChange":
		a.ConfigChange = new(ConfigChange)
		return json.Unmarshal(input, a.ConfigChange)
	case "delegation":
		a.Delegation = new(Delegation)
		return json.Unmarshal(input, a.Delegation)
	default:
		return fmt.Errorf("eip8130: unknown account change type %q", tag.Type)
	}
}
