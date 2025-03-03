// +build ctidh
// ctidh_test.go - Adapts ctidh module to our NIKE interface.
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

package ctidh

import (
	"testing"

	ctidh "git.xx.network/elixxir/ctidh_cgo"
	"github.com/stretchr/testify/require"
)

func TestCtidhNike(t *testing.T) {
	ctidhNike := new(CtidhNike)

	alicePrivateKey, alicePublicKey := ctidhNike.NewKeypair()

	tmp := ctidh.DerivePublicKey(alicePrivateKey.(*ctidh.PrivateKey))
	require.Equal(t, alicePublicKey.Bytes(), tmp.Bytes())

	bobPrivKey, bobPubKey := ctidh.GenerateKeyPair()

	aliceS := ctidhNike.DeriveSecret(alicePrivateKey, bobPubKey)

	bobS := ctidh.DeriveSecret(bobPrivKey, alicePublicKey.(*ctidh.PublicKey))
	require.Equal(t, bobS, aliceS)
}
