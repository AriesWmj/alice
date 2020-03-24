// Copyright © 2020 AMIS Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package birkhoffinterpolation

import (
	"math/big"
	"testing"

	"github.com/getamis/alice/crypto/matrix"
	"github.com/getamis/alice/crypto/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func TestBirkhoffinterpolation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Birkhoffinterpolation Suite")
}

var _ = Describe("Birkhoff Interpolation", func() {
	var (
		bigNumber   = "115792089237316195423570985008687907852837564279074904382605163141518161494337"
		bigPrime, _ = new(big.Int).SetString(bigNumber, 10)
	)

	Context("getLinearEquationCoefficientMatrix()", func() {
		It("should be ok", func() {
			ps := make(BkParameters, 5)
			ps[0] = NewBkParameter(big.NewInt(1), 0)
			ps[1] = NewBkParameter(big.NewInt(2), 1)
			ps[2] = NewBkParameter(big.NewInt(3), 2)
			ps[3] = NewBkParameter(big.NewInt(4), 3)
			ps[4] = NewBkParameter(big.NewInt(5), 4)
			got, err := ps.getLinearEquationCoefficientMatrix(4, bigPrime)
			Expect(err).Should(BeNil())

			expectedMatrix := [][]*big.Int{
				{big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1)},
				{big.NewInt(0), big.NewInt(1), big.NewInt(4), big.NewInt(12)},
				{big.NewInt(0), big.NewInt(0), big.NewInt(2), big.NewInt(18)},
				{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(6)},
				{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
			}
			expected, err := matrix.NewMatrix(bigPrime, expectedMatrix)
			Expect(got).Should(Equal(expected))
			Expect(err).Should(BeNil())
		})
	})

	It("Getter func", func() {
		x := big.NewInt(1)
		rank := uint32(0)
		bk := NewBkParameter(x, rank)
		Expect(bk.GetX()).Should(Equal(x))
		Expect(bk.GetRank()).Should(Equal(rank))
		Expect(bk.String()).Should(Equal("(x, rank) = (1, 0)"))
	})

	DescribeTable("VerifyEnoughRankCanRecoverSecret func", func(ps BkParameters) {
		err := ps.CheckValid(uint32(3), bigPrime)
		Expect(err).Should(BeNil())
	},
		Entry("BK:(x,rank):(1,0),(2,1),(3,2),(5,4),(4,3)",
			[]*BkParameter{NewBkParameter(big.NewInt(1), 0), NewBkParameter(big.NewInt(2), 1),
				NewBkParameter(big.NewInt(3), 2), NewBkParameter(big.NewInt(5), 4), NewBkParameter(big.NewInt(4), 3)},
		),
		Entry("BK:(x,rank):(1,0),(2,3),(3,0),(5,0),(4,0)",
			[]*BkParameter{NewBkParameter(big.NewInt(1), 0), NewBkParameter(big.NewInt(2), 3),
				NewBkParameter(big.NewInt(3), 0), NewBkParameter(big.NewInt(5), 0), NewBkParameter(big.NewInt(4), 0)},
		),
		Entry("BK:(x,rank):(1,0),(2,0),(3,0),(5,0),(4,0)",
			[]*BkParameter{NewBkParameter(big.NewInt(1), 0), NewBkParameter(big.NewInt(2), 0),
				NewBkParameter(big.NewInt(3), 0), NewBkParameter(big.NewInt(5), 0), NewBkParameter(big.NewInt(4), 0)},
		),
		Entry("BK:(x,rank):(1,1),(2,1),(3,1),(5,0),(4,0)",
			[]*BkParameter{NewBkParameter(big.NewInt(1), 1), NewBkParameter(big.NewInt(2), 1),
				NewBkParameter(big.NewInt(3), 1), NewBkParameter(big.NewInt(5), 0), NewBkParameter(big.NewInt(4), 0)},
		),
		Entry("BK:(x,rank):(1,1),(2,1),(3,1),(5,1),(4,0)",
			[]*BkParameter{NewBkParameter(big.NewInt(1), 1), NewBkParameter(big.NewInt(2), 1),
				NewBkParameter(big.NewInt(3), 1), NewBkParameter(big.NewInt(5), 1), NewBkParameter(big.NewInt(4), 0)},
		),
	)

	It("duplicate Bk", func() {
		ps := make(BkParameters, 5)
		ps[0] = NewBkParameter(big.NewInt(1), 0)
		ps[1] = NewBkParameter(big.NewInt(2), 1)
		ps[2] = NewBkParameter(big.NewInt(3), 2)
		ps[3] = NewBkParameter(big.NewInt(1), 0)
		ps[4] = NewBkParameter(big.NewInt(5), 4)
		err := ps.CheckValid(uint32(3), bigPrime)
		Expect(err).Should(Equal(ErrInvalidBks))
	})

	// The problem is that (1,0) and (2,1) and (3,0) can not recover secret.
	It("Expect no valid bks", func() {
		ps := make(BkParameters, 5)
		ps[3] = NewBkParameter(big.NewInt(4), 2)
		ps[4] = NewBkParameter(big.NewInt(5), 2)
		ps[0] = NewBkParameter(big.NewInt(1), 2)
		ps[1] = NewBkParameter(big.NewInt(2), 2)
		ps[2] = NewBkParameter(big.NewInt(3), 2)
		err := ps.CheckValid(uint32(3), bigPrime)
		Expect(err).Should(Equal(ErrNoValidBks))
	})

	// The problem is that (1,0) and (2,1) and (3,0) can not recover secret.
	It("Expect Enough Rank but not have", func() {
		ps := make(BkParameters, 5)
		ps[3] = NewBkParameter(big.NewInt(4), 0)
		ps[4] = NewBkParameter(big.NewInt(5), 0)
		ps[0] = NewBkParameter(big.NewInt(1), 0)
		ps[1] = NewBkParameter(big.NewInt(2), 1)
		ps[2] = NewBkParameter(big.NewInt(3), 0)
		err := ps.CheckValid(uint32(3), bigPrime)
		Expect(err).Should(Equal(ErrInvalidBks))
	})

	Context("ComputeBkCoefficient()", func() {
		It("should be ok", func() {
			ps := make(BkParameters, 4)
			ps[0] = NewBkParameter(big.NewInt(1), 0)
			ps[1] = NewBkParameter(big.NewInt(2), 1)
			ps[2] = NewBkParameter(big.NewInt(3), 2)
			ps[3] = NewBkParameter(big.NewInt(4), 3)
			expectedStrs := []string{
				"1",
				"115792089237316195423570985008687907852837564279074904382605163141518161494336",
				"57896044618658097711785492504343953926418782139537452191302581570759080747170",
				"0",
			}
			expected := make([]*big.Int, len(expectedStrs))
			for i, s := range expectedStrs {
				expected[i], _ = new(big.Int).SetString(s, 10)
			}
			got, err := ps.ComputeBkCoefficient(3, bigPrime)
			Expect(err).Should(BeNil())
			Expect(got).Should(Equal(expected))
		})

		It("invalid field order", func() {
			ps := make(BkParameters, 4)
			ps[0] = NewBkParameter(big.NewInt(1), 0)
			ps[1] = NewBkParameter(big.NewInt(2), 1)
			ps[2] = NewBkParameter(big.NewInt(3), 2)
			ps[3] = NewBkParameter(big.NewInt(4), 3)
			got, err := ps.ComputeBkCoefficient(3, big.NewInt(2))
			Expect(err).Should(Equal(utils.ErrLessOrEqualBig2))
			Expect(got).Should(BeNil())
		})

		It("larger threshold", func() {
			ps := make(BkParameters, 2)
			ps[0] = NewBkParameter(big.NewInt(1), 0)
			ps[1] = NewBkParameter(big.NewInt(2), 1)
			got, err := ps.ComputeBkCoefficient(3, bigPrime)
			Expect(err).Should(Equal(ErrEqualOrLargerThreshold))
			Expect(got).Should(BeNil())
		})

		It("not invertible matrix #0", func() {
			ps := make(BkParameters, 4)
			ps[0] = NewBkParameter(big.NewInt(1), 2)
			ps[1] = NewBkParameter(big.NewInt(2), 2)
			ps[2] = NewBkParameter(big.NewInt(3), 3)
			ps[3] = NewBkParameter(big.NewInt(4), 0)
			got, err := ps.ComputeBkCoefficient(3, bigPrime)
			Expect(err).Should(Equal(matrix.ErrNotInvertableMatrix))
			Expect(got).Should(BeNil())
		})

		It("not invertible matrix #1", func() {
			ps := make(BkParameters, 5)
			ps[0] = NewBkParameter(big.NewInt(1), 2)
			ps[1] = NewBkParameter(big.NewInt(2), 2)
			ps[2] = NewBkParameter(big.NewInt(3), 3)
			ps[3] = NewBkParameter(big.NewInt(4), 1)
			ps[4] = NewBkParameter(big.NewInt(5), 4)
			got, err := ps.ComputeBkCoefficient(3, bigPrime)
			Expect(err).Should(Equal(matrix.ErrNotInvertableMatrix))
			Expect(got).Should(BeNil())
		})

		It("not invertible matrix #2 - two the same X", func() {
			ps := make(BkParameters, 5)
			ps[0] = NewBkParameter(big.NewInt(1), 1)
			ps[1] = NewBkParameter(big.NewInt(2), 3)
			ps[2] = NewBkParameter(big.NewInt(3), 3)
			ps[3] = NewBkParameter(big.NewInt(1), 1)
			ps[4] = NewBkParameter(big.NewInt(5), 3)
			got, err := ps.ComputeBkCoefficient(3, bigPrime)
			Expect(err).Should(Equal(matrix.ErrNotInvertableMatrix))
			Expect(got).Should(BeNil())
		})
	})
})