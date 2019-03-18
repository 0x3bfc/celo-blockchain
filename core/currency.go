// Copyright 2017 The Celo Authors
// This file is part of the celo library.
//
// The celo library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The celo library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the celo library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var (
	cgExchangeRateNum = big.NewInt(1)
	cgExchangeRateDen = big.NewInt(1)
)

type exchangeRate struct {
	Numerator   *big.Int
	Denominator *big.Int
}

type PriceComparator struct {
	exchangeRates map[common.Address]*exchangeRate // indexedCurrency:CeloGold exchange rate
}

func (pc *PriceComparator) Cmp(val1 *big.Int, currency1 *common.Address, val2 *big.Int, currency2 *common.Address) int {
	if currency1 == currency2 {
		return val1.Cmp(val2)
	}

	var exchangeRate1Num *big.Int
	var exchangeRate1Den *big.Int
	if currency1 == nil {
		exchangeRate1Num = cgExchangeRateNum
		exchangeRate1Den = cgExchangeRateDen
	} else {
		exchangeRate1 := pc.exchangeRates[*currency1]
		exchangeRate1Num = exchangeRate1.Numerator
		exchangeRate1Den = exchangeRate1.Denominator
	}

	var exchangeRate2Num *big.Int
	var exchangeRate2Den *big.Int
	if currency2 == nil {
		exchangeRate2Num = cgExchangeRateNum
		exchangeRate2Den = cgExchangeRateDen
	} else {
		exchangeRate2 := pc.exchangeRates[*currency2]
		exchangeRate2Num = exchangeRate2.Numerator
		exchangeRate2Den = exchangeRate2.Denominator
	}

	// Below code block is basically evaluating this comparison:
	// val1 * exchangeRate1Num/exchangeRate1Den < val2 * exchangeRate2Num/exchangeRate2Den
	// It Will transform that comparison to this, to remove having to deal with fractional values.
	// val1 * exchangeRate1Num * exchangeRate2Den < val2 * exchangeRate2Num * exchangeRate1Den
	leftSide := new(big.Int).Mul(val1, new(big.Int).Mul(exchangeRate1Num, exchangeRate2Den))
	rightSide := new(big.Int).Mul(val2, new(big.Int).Mul(exchangeRate2Num, exchangeRate1Den))
	return leftSide.Cmp(rightSide)
}

func NewPriceComparator() *PriceComparator {
	// For now, hard coding in some exchange rates.  Will modify this to retrieve the
	// exchange rates from the Celo's exchange smart contract.
	// C$ will have a 2:1 exchange rate with CG
	exchangeRates := make(map[common.Address]*exchangeRate)
	exchangeRates[common.HexToAddress("0x0000000000000000000000000000000ce10d011a")] = &exchangeRate{Numerator: big.NewInt(2), Denominator: big.NewInt(1)}

	return &PriceComparator{
		exchangeRates: exchangeRates,
	}
}
