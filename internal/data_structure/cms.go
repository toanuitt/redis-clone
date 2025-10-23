package data_structure

import (
	"math"

	"github.com/spaolacci/murmur3"
)

type CMS struct {
	width   uint32
	depth   uint32
	counter [][]uint32
}

// Log10PointFive is a precomputed value for log10(0.5).
const Log10PointFive = -0.30102999566

func CreateCMS(width uint32, depth uint32) *CMS {
	CMS := &CMS{
		width: width,
		depth: depth,
	}
	CMS.counter = make([][]uint32, depth)
	for i := uint32(0); i < depth; i++ {
		CMS.counter[i] = make([]uint32, width)
	}
	return CMS
}

func CalcCMSDim(errRate float64, errProb float64) (uint32, uint32) {
	w := uint32(math.Ceil(2.0 / errRate))
	d := uint32(math.Ceil(math.Log10(errProb) / Log10PointFive))
	return w, d
}

func (c *CMS) calcHash(item string, seed uint32) uint32 {
	hasher := murmur3.New32WithSeed(seed)
	hasher.Write([]byte(item))
	return hasher.Sum32()
}

func (c *CMS) IncrBy(item string, value uint32) uint32 {
	var minCount uint32 = math.MaxUint32

	// Loop through each row of the 2D array.
	for i := uint32(0); i < c.depth; i++ {
		// Calculate a new hash for each row using the row index as the seed.
		hash := c.calcHash(item, i)
		// Use the hash to get the column index within the row.
		j := hash % c.width

		// Safely add the value to prevent overflow.
		if math.MaxUint32-c.counter[i][j] < value {
			c.counter[i][j] = math.MaxUint32
		} else {
			c.counter[i][j] += value
		}

		// Keep track of the minimum count across all rows.
		if c.counter[i][j] < minCount {
			minCount = c.counter[i][j]
		}
	}
	return minCount
}

// Count returns the estimated count for an item.
// It retrieves the minimum count across all hash functions to provide the most accurate estimate.
func (c *CMS) Count(item string) uint32 {
	var minCount uint32 = math.MaxUint32

	// Loop through each row of the 2D array.
	for i := uint32(0); i < c.depth; i++ {
		// Calculate the hash for this row.
		hash := c.calcHash(item, i)
		// Determine the column index.
		j := hash % c.width

		// Find the minimum count across all rows.
		if c.counter[i][j] < minCount {
			minCount = c.counter[i][j]
		}
	}
	return minCount
}
