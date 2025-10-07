package native

import (
	"sort"
)

// MedianCutQuantize quantizes an image using the Median Cut algorithm
// This algorithm produces better results for photographic images than Octree
func MedianCutQuantize(pixels []RGB, maxColors int) ([]byte, []RGB) {
	if len(pixels) == 0 {
		return []byte{}, []RGB{}
	}

	// Create initial bucket with all colors
	bucket := &colorBucket{
		pixels: pixels,
	}
	bucket.calculateBounds()

	buckets := []*colorBucket{bucket}

	// Split buckets until we have maxColors buckets
	for len(buckets) < maxColors {
		// Find bucket with largest range
		largest := findLargestBucket(buckets)
		if largest == nil {
			break
		}

		// Split the bucket
		bucket1, bucket2 := largest.split()
		if bucket1 == nil || bucket2 == nil {
			break
		}

		// Replace largest bucket with the two new buckets
		buckets = removeBucket(buckets, largest)
		buckets = append(buckets, bucket1, bucket2)
	}

	// Generate palette from buckets
	palette := make([]RGB, len(buckets))
	for i, bucket := range buckets {
		palette[i] = bucket.averageColor()
	}

	// Map pixels to palette
	indexed := make([]byte, len(pixels))
	for i, pixel := range pixels {
		closestIdx := findClosestColorIndex(palette, pixel)
		indexed[i] = byte(closestIdx)
	}

	return indexed, palette
}

type colorBucket struct {
	pixels    []RGB
	minR, maxR byte
	minG, maxG byte
	minB, maxB byte
}

func (b *colorBucket) calculateBounds() {
	if len(b.pixels) == 0 {
		return
	}

	b.minR, b.maxR = 255, 0
	b.minG, b.maxG = 255, 0
	b.minB, b.maxB = 255, 0

	for _, p := range b.pixels {
		if p.R < b.minR {
			b.minR = p.R
		}
		if p.R > b.maxR {
			b.maxR = p.R
		}
		if p.G < b.minG {
			b.minG = p.G
		}
		if p.G > b.maxG {
			b.maxG = p.G
		}
		if p.B < b.minB {
			b.minB = p.B
		}
		if p.B > b.maxB {
			b.maxB = p.B
		}
	}
}

func (b *colorBucket) largestDimension() int {
	rRange := int(b.maxR) - int(b.minR)
	gRange := int(b.maxG) - int(b.minG)
	bRange := int(b.maxB) - int(b.minB)

	if rRange >= gRange && rRange >= bRange {
		return 0 // Red
	} else if gRange >= bRange {
		return 1 // Green
	}
	return 2 // Blue
}

func (b *colorBucket) split() (*colorBucket, *colorBucket) {
	if len(b.pixels) < 2 {
		return nil, nil
	}

	// Sort by largest dimension
	dim := b.largestDimension()
	sort.Slice(b.pixels, func(i, j int) bool {
		switch dim {
		case 0:
			return b.pixels[i].R < b.pixels[j].R
		case 1:
			return b.pixels[i].G < b.pixels[j].G
		default:
			return b.pixels[i].B < b.pixels[j].B
		}
	})

	// Split at median
	median := len(b.pixels) / 2

	bucket1 := &colorBucket{pixels: b.pixels[:median]}
	bucket2 := &colorBucket{pixels: b.pixels[median:]}

	bucket1.calculateBounds()
	bucket2.calculateBounds()

	return bucket1, bucket2
}

func (b *colorBucket) averageColor() RGB {
	if len(b.pixels) == 0 {
		return RGB{}
	}

	var sumR, sumG, sumB int
	for _, p := range b.pixels {
		sumR += int(p.R)
		sumG += int(p.G)
		sumB += int(p.B)
	}

	count := len(b.pixels)
	return RGB{
		R: byte(sumR / count),
		G: byte(sumG / count),
		B: byte(sumB / count),
	}
}

func (b *colorBucket) range_() int {
	rRange := int(b.maxR) - int(b.minR)
	gRange := int(b.maxG) - int(b.minG)
	bRange := int(b.maxB) - int(b.minB)

	// Weighted by perceptual importance
	return 2*rRange + 4*gRange + 3*bRange
}

func findLargestBucket(buckets []*colorBucket) *colorBucket {
	if len(buckets) == 0 {
		return nil
	}

	largest := buckets[0]
	largestRange := largest.range_()

	for _, bucket := range buckets[1:] {
		r := bucket.range_()
		if r > largestRange {
			largest = bucket
			largestRange = r
		}
	}

	return largest
}

func removeBucket(buckets []*colorBucket, toRemove *colorBucket) []*colorBucket {
	result := make([]*colorBucket, 0, len(buckets)-1)
	for _, bucket := range buckets {
		if bucket != toRemove {
			result = append(result, bucket)
		}
	}
	return result
}

func findClosestColorIndex(palette []RGB, color RGB) int {
	minDist := uint32(0xFFFFFFFF)
	closest := 0

	for i, c := range palette {
		dr := int(color.R) - int(c.R)
		dg := int(color.G) - int(c.G)
		db := int(color.B) - int(c.B)

		// Perceptual distance
		dist := uint32(2*dr*dr + 4*dg*dg + 3*db*db)

		if dist < minDist {
			minDist = dist
			closest = i
			if dist == 0 {
				break
			}
		}
	}

	return closest
}
