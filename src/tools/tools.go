/*This file is part of Algebror.

Algebror is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

Algebror is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with Algebror. If not, see <https://www.gnu.org/licenses/>.

*/

package tools

import (
	"math/rand"
	"strconv"
)

func GetRandomInteger(r, d int) int {
	n := rand.Intn(r * d)
	if n == 0 {
		return n + 1
	} else {
		return n
	}
}

func GetRandomDecimal(r, d int) float64 {
	precision := 0
	if d > 4 {
		precision = 100
	} else {
		precision = 10
	}
	n := float64(rand.Intn(r*precision*d)) / float64(precision)
	if n == 0 {
		return n + 1/float64(precision)
	} else {
		return n
	}
}

func FormatNumber(num float64, p int) string {
	return strconv.FormatFloat(num, 'f', p, 64)
}
