// descriptor_test.go - Descriptor s11n tests.
// Copyright (C) 2017  Yawning Angel
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package pki

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/katzenpost/katzenpost/core/crypto/cert"
	"github.com/katzenpost/katzenpost/core/crypto/ecdh"
	"github.com/katzenpost/katzenpost/core/crypto/rand"
	"github.com/katzenpost/katzenpost/core/wire"
)

const debugTestEpoch = 0xFFFFFFFF

func TestDescriptor(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	require := require.New(t)

	d := new(MixDescriptor)
	d.Epoch = debugTestEpoch

	err := IsDescriptorWellFormed(d, debugTestEpoch)
	assert.Error(err, "IsDescriptorWellFormed(bad)")

	// Build a well formed descriptor.
	d.Name = "hydra-dominatus.example.net"
	d.Addresses = map[Transport][]string{
		TransportTCPv4:     []string{"192.0.2.1:4242", "192.0.2.1:1234", "198.51.100.2:4567"},
		TransportTCPv6:     []string{"[2001:DB8::1]:8901"},
		Transport("torv2"): []string{"thisisanoldonion.onion:2323"},
		TransportTCP:       []string{"example.com:4242"},
	}
	d.Provider = true
	d.LoadWeight = 23
	identityPriv, identityPub := cert.Scheme.NewKeypair()
	d.IdentityKey = identityPub
	scheme := wire.DefaultScheme
	_, d.LinkKey = scheme.GenerateKeypair(rand.Reader)
	d.MixKeys = make(map[uint64]*ecdh.PublicKey)
	for e := debugTestEpoch; e < debugTestEpoch+3; e++ {
		mPriv, err := ecdh.NewKeypair(rand.Reader)
		require.NoError(err, "[%d]: ecdh.NewKeypair()", e)
		d.MixKeys[uint64(e)] = mPriv.PublicKey()
	}
	d.Kaetzchen = make(map[string]map[string]interface{})
	d.Kaetzchen["miau"] = map[string]interface{}{
		"endpoint":  "+miau",
		"miauCount": 23,
	}
	err = IsDescriptorWellFormed(d, debugTestEpoch)
	require.NoError(err, "IsDescriptorWellFormed(good)")

	// Sign the descriptor.
	signed, err := SignDescriptor(identityPriv, identityPub, d)
	require.NoError(err, "SignDescriptor()")

	// Verify and deserialize the signed descriptor.
	dd := new(MixDescriptor)
	err = dd.UnmarshalBinary(signed)
	require.NoError(err)

	// Ensure the base and de-serialized descriptors match.
	assert.Equal(d.Name, dd.Name, "Name")
	assert.Equal(d.Addresses, dd.Addresses, "Addresses")
	assert.Equal(d.Provider, dd.Provider, "Provider")
	assert.Equal(d.LoadWeight, dd.LoadWeight, "LoadWeight")
	assert.Equal(d.IdentityKey.Bytes(), dd.IdentityKey.Bytes(), "IdentityKey")
	assert.Equal(d.LinkKey.Bytes(), dd.LinkKey.Bytes(), "LinkKey")
	require.Equal(len(d.MixKeys), len(dd.MixKeys), "len(MixKeys)")
	for k, v := range d.MixKeys {
		vv := dd.MixKeys[k]
		require.NotNil(vv)
		require.Equal(v.Bytes(), vv.Bytes(), "MixKeys[%v]", k)
	}
}
