// ecdh.go - Adapts ecdh module to our NIKE interface.
// Copyright (C) 2022  David Stainton.
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

package ecdh

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/katzenpost/katzenpost/core/crypto/ecdh"
	"github.com/katzenpost/katzenpost/core/crypto/rand"
)

func TestEcdhNike(t *testing.T) {
	ecdhNike := NewEcdhNike(rand.Reader)

	alicePrivateKey, alicePublicKey := ecdhNike.NewKeypair()

	tmp := alicePrivateKey.(*ecdh.PrivateKey).PublicKey()
	require.Equal(t, alicePublicKey.Bytes(), tmp.Bytes())

	bobKeypair, err := ecdh.NewKeypair(rand.Reader)
	require.NoError(t, err)

	aliceS := ecdhNike.DeriveSecret(alicePrivateKey, bobKeypair.PublicKey())

	bobS := ecdh.Exp(alicePublicKey.Bytes(), bobKeypair.Bytes())
	require.Equal(t, bobS, aliceS)
}
