package native

// RGB represents an RGB color
type RGB struct {
	R, G, B byte
}

// OctreeNode represents a node in the octree color quantization tree
type OctreeNode struct {
	isLeaf      bool
	pixelCount  int
	redSum      int
	greenSum    int
	blueSum     int
	children    [8]*OctreeNode
	paletteIndex int
	level       int
}

// OctreeQuantizer implements Octree color quantization algorithm
type OctreeQuantizer struct {
	root           *OctreeNode
	maxColors      int
	leafCount      int
	reducibleNodes [9][]*OctreeNode // One list per level (0-8)
	palette        []RGB
}

// priorityQueue for reducible nodes
type reducibleQueue struct {
	nodes []*OctreeNode
}

func (pq *reducibleQueue) Len() int { return len(pq.nodes) }

func (pq *reducibleQueue) Less(i, j int) bool {
	// Prefer nodes with fewer children and at deeper levels
	if pq.nodes[i].level != pq.nodes[j].level {
		return pq.nodes[i].level > pq.nodes[j].level
	}
	return pq.nodes[i].pixelCount < pq.nodes[j].pixelCount
}

func (pq *reducibleQueue) Swap(i, j int) {
	pq.nodes[i], pq.nodes[j] = pq.nodes[j], pq.nodes[i]
}

func (pq *reducibleQueue) Push(x interface{}) {
	pq.nodes = append(pq.nodes, x.(*OctreeNode))
}

func (pq *reducibleQueue) Pop() interface{} {
	old := pq.nodes
	n := len(old)
	node := old[n-1]
	pq.nodes = old[0 : n-1]
	return node
}

// NewOctreeQuantizer creates a new octree quantizer
func NewOctreeQuantizer(maxColors int) *OctreeQuantizer {
	return &OctreeQuantizer{
		root:           &OctreeNode{level: 0},
		maxColors:      maxColors,
		leafCount:      0,
		reducibleNodes: [9][]*OctreeNode{},
	}
}

// getColorIndex returns the octree index (0-7) for a color at a given level
func getColorIndex(r, g, b byte, level int) int {
	shift := 7 - level
	index := 0
	if (r>>shift)&1 == 1 {
		index |= 4
	}
	if (g>>shift)&1 == 1 {
		index |= 2
	}
	if (b>>shift)&1 == 1 {
		index |= 1
	}
	return index
}

// AddColor adds a color to the octree
func (oq *OctreeQuantizer) AddColor(r, g, b byte) {
	node := oq.root

	for level := 0; level < 8; level++ {
		index := getColorIndex(r, g, b, level)

		if node.children[index] == nil {
			newNode := &OctreeNode{level: level + 1}
			node.children[index] = newNode

			// Add to reducible nodes if not at leaf level
			if level < 7 {
				oq.reducibleNodes[level] = append(oq.reducibleNodes[level], newNode)
			} else {
				newNode.isLeaf = true
				oq.leafCount++
			}
		}

		node = node.children[index]
	}

	// Accumulate color values
	node.pixelCount++
	node.redSum += int(r)
	node.greenSum += int(g)
	node.blueSum += int(b)

	// Reduce tree if necessary
	for oq.leafCount > oq.maxColors {
		oq.reduceTree()
	}
}

// reduceTree reduces the tree by merging nodes at the deepest level
func (oq *OctreeQuantizer) reduceTree() {
	// Find deepest level with reducible nodes
	for level := 7; level >= 0; level-- {
		if len(oq.reducibleNodes[level]) > 0 {
			// Get and remove last node from this level
			node := oq.reducibleNodes[level][len(oq.reducibleNodes[level])-1]
			oq.reducibleNodes[level] = oq.reducibleNodes[level][:len(oq.reducibleNodes[level])-1]

			// Merge children
			oq.mergeChildren(node)
			return
		}
	}
}

// mergeChildren merges all children of a node into the node itself
func (oq *OctreeQuantizer) mergeChildren(node *OctreeNode) {
	// Sum up all children
	for i := 0; i < 8; i++ {
		if node.children[i] != nil {
			child := node.children[i]

			node.pixelCount += child.pixelCount
			node.redSum += child.redSum
			node.greenSum += child.greenSum
			node.blueSum += child.blueSum

			if child.isLeaf {
				oq.leafCount--
			}

			node.children[i] = nil
		}
	}

	// Convert node to leaf
	node.isLeaf = true
	oq.leafCount++
}

// GeneratePalette generates the color palette from the octree
func (oq *OctreeQuantizer) GeneratePalette() []RGB {
	oq.palette = make([]RGB, 0, oq.maxColors)
	oq.generatePaletteRecursive(oq.root, 0)

	// Fill remaining palette entries with black if needed
	for len(oq.palette) < oq.maxColors {
		oq.palette = append(oq.palette, RGB{0, 0, 0})
	}

	return oq.palette
}

// generatePaletteRecursive recursively generates palette colors
func (oq *OctreeQuantizer) generatePaletteRecursive(node *OctreeNode, paletteIndex int) int {
	if node.isLeaf {
		if node.pixelCount > 0 {
			r := byte(node.redSum / node.pixelCount)
			g := byte(node.greenSum / node.pixelCount)
			b := byte(node.blueSum / node.pixelCount)

			node.paletteIndex = len(oq.palette)
			oq.palette = append(oq.palette, RGB{R: r, G: g, B: b})
		} else {
			node.paletteIndex = 0
		}
		return node.paletteIndex
	}

	// Recurse through children
	for i := 0; i < 8; i++ {
		if node.children[i] != nil {
			oq.generatePaletteRecursive(node.children[i], paletteIndex)
		}
	}

	return paletteIndex
}

// GetPaletteIndex returns the palette index for a given color
func (oq *OctreeQuantizer) GetPaletteIndex(r, g, b byte) byte {
	node := oq.root

	for level := 0; level < 8; level++ {
		if node.isLeaf {
			return byte(node.paletteIndex)
		}

		index := getColorIndex(r, g, b, level)
		if node.children[index] == nil {
			// No exact match, find closest leaf
			return oq.findClosestLeaf(node, r, g, b)
		}

		node = node.children[index]
	}

	if node.isLeaf {
		return byte(node.paletteIndex)
	}

	// Shouldn't reach here, but return 0 as fallback
	return 0
}

// findClosestLeaf finds the closest leaf node from current node
func (oq *OctreeQuantizer) findClosestLeaf(node *OctreeNode, r, g, b byte) byte {
	// Simple approach: check all children and find closest
	minDist := uint32(0xFFFFFFFF)
	closestIndex := 0

	for i := 0; i < 8; i++ {
		if node.children[i] != nil {
			idx := oq.findClosestInSubtree(node.children[i], r, g, b)
			color := oq.palette[idx]

			dr := int(r) - int(color.R)
			dg := int(g) - int(color.G)
			db := int(b) - int(color.B)
			dist := uint32(dr*dr + dg*dg + db*db)

			if dist < minDist {
				minDist = dist
				closestIndex = idx
			}
		}
	}

	return byte(closestIndex)
}

// findClosestInSubtree finds the first leaf in a subtree
func (oq *OctreeQuantizer) findClosestInSubtree(node *OctreeNode, r, g, b byte) int {
	if node.isLeaf {
		return node.paletteIndex
	}

	// Find first leaf by traversing tree
	for i := 0; i < 8; i++ {
		if node.children[i] != nil {
			return oq.findClosestInSubtree(node.children[i], r, g, b)
		}
	}

	return 0
}

// QuantizeImageOctree quantizes an image using Octree algorithm with optimization
func QuantizeImageOctree(pixels []RGB, maxColors int) ([]byte, []RGB) {
	// Infer dimensions (assume reasonable aspect ratio)
	width, height := inferDimensions(len(pixels))
	return QuantizeImageOctreeWithDimensions(pixels, maxColors, width, height)
}

// QuantizeImageOctreeWithDimensions quantizes with known dimensions for better dithering
func QuantizeImageOctreeWithDimensions(pixels []RGB, maxColors, width, height int) ([]byte, []RGB) {
	// Fast path: if image has fewer unique colors than maxColors, use them directly
	uniqueColors := make(map[uint32]struct{})
	for _, pixel := range pixels {
		colorKey := (uint32(pixel.R) << 16) | (uint32(pixel.G) << 8) | uint32(pixel.B)
		uniqueColors[colorKey] = struct{}{}

		// Early exit if we already have more than maxColors
		if len(uniqueColors) > maxColors*2 {
			break
		}
	}

	// If few unique colors, use them directly without complex quantization
	if len(uniqueColors) <= maxColors {
		return quantizeSimpleWithDithering(pixels, uniqueColors, width, height)
	}

	// Otherwise use Octree algorithm with histogram optimization
	return quantizeOctreeOptimizedWithDithering(pixels, maxColors, width, height)
}

// inferDimensions tries to guess reasonable image dimensions
func inferDimensions(pixelCount int) (width, height int) {
	// Try to find dimensions that give reasonable aspect ratio
	width = 1
	for width*width < pixelCount {
		width++
	}

	// Check if it's exactly square
	if width*width == pixelCount {
		return width, width
	}

	// Try to find divisor that gives aspect ratio close to 16:9 or 4:3
	bestWidth := width
	bestHeight := pixelCount / width

	for w := width - 10; w <= width+10 && w > 0; w++ {
		if pixelCount%w == 0 {
			h := pixelCount / w
			ratio := float64(w) / float64(h)

			// Prefer ratios closer to common aspect ratios
			if ratio > 0.5 && ratio < 2.5 {
				// Check if this ratio is better than current best
				if (ratio >= 1.3 && ratio <= 1.8) || (ratio >= 0.55 && ratio <= 0.8) {
					bestWidth = w
					bestHeight = h
				}
			}
		}
	}

	return bestWidth, bestHeight
}

// quantizeSimpleWithDithering handles images with few unique colors efficiently with dithering
func quantizeSimpleWithDithering(pixels []RGB, uniqueColors map[uint32]struct{}, width, height int) ([]byte, []RGB) {
	// Build palette from unique colors
	palette := make([]RGB, 0, len(uniqueColors))
	colorToIndex := make(map[uint32]byte)

	for colorKey := range uniqueColors {
		r := byte(colorKey >> 16)
		g := byte(colorKey >> 8)
		b := byte(colorKey)

		palette = append(palette, RGB{R: r, G: g, B: b})
		colorToIndex[colorKey] = byte(len(palette) - 1)
	}

	// Fill to 256 colors if needed
	for len(palette) < 256 {
		palette = append(palette, RGB{0, 0, 0})
	}

	// Simple quantization WITHOUT dithering (like Pillow with optimize=False)
	indexed := make([]byte, len(pixels))
	for i, p := range pixels {
		colorKey := (uint32(p.R) << 16) | (uint32(p.G) << 8) | uint32(p.B)
		if idx, exists := colorToIndex[colorKey]; exists {
			indexed[i] = idx
		} else {
			// Find closest
			indexed[i] = findClosestPaletteColor(palette, p.R, p.G, p.B)
		}
	}

	return indexed, palette
}

// quantizeOctreeOptimizedWithDithering uses histogram to reduce redundant processing with dithering
func quantizeOctreeOptimizedWithDithering(pixels []RGB, maxColors, width, height int) ([]byte, []RGB) {
	// Build histogram of colors
	histogram := make(map[uint32]int)
	for _, pixel := range pixels {
		colorKey := (uint32(pixel.R) << 16) | (uint32(pixel.G) << 8) | uint32(pixel.B)
		histogram[colorKey]++
	}

	// Build octree from unique colors only (much faster)
	quantizer := NewOctreeQuantizer(maxColors)
	for colorKey, count := range histogram {
		r := byte(colorKey >> 16)
		g := byte(colorKey >> 8)
		b := byte(colorKey)

		// Add color with its frequency
		node := quantizer.root
		for level := 0; level < 8; level++ {
			index := getColorIndex(r, g, b, level)

			if node.children[index] == nil {
				newNode := &OctreeNode{level: level + 1}
				node.children[index] = newNode

				if level < 7 {
					quantizer.reducibleNodes[level] = append(quantizer.reducibleNodes[level], newNode)
				} else {
					newNode.isLeaf = true
					quantizer.leafCount++
				}
			}
			node = node.children[index]
		}

		// Accumulate with count (weighted)
		node.pixelCount += count
		node.redSum += int(r) * count
		node.greenSum += int(g) * count
		node.blueSum += int(b) * count
	}

	// Reduce tree to target colors
	for quantizer.leafCount > maxColors {
		quantizer.reduceTree()
	}

	// Generate palette
	palette := quantizer.GeneratePalette()

	// Build fast lookup table for unique colors
	colorToIndex := make(map[uint32]byte, len(histogram))
	for colorKey := range histogram {
		r := byte(colorKey >> 16)
		g := byte(colorKey >> 8)
		b := byte(colorKey)
		colorToIndex[colorKey] = quantizer.GetPaletteIndex(r, g, b)
	}

	// Apply Floyd-Steinberg dithering for better quality
	indexed := applyFloydSteinberg(pixels, palette, colorToIndex, width, height)

	return indexed, palette
}

// applyFloydSteinberg applies Floyd-Steinberg dithering algorithm
func applyFloydSteinberg(pixels []RGB, palette []RGB, colorToIndex map[uint32]byte, width, height int) []byte {
	indexed := make([]byte, len(pixels))

	// Create working buffer for error diffusion
	workPixels := make([]struct{ r, g, b int }, len(pixels))
	for i, p := range pixels {
		workPixels[i].r = int(p.R)
		workPixels[i].g = int(p.G)
		workPixels[i].b = int(p.B)
	}

	pixelCount := len(pixels)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x

			// Safety check: ensure idx is within bounds
			if idx >= pixelCount {
				break
			}

			// Clamp values to valid range [0, 255]
			oldR := clampInt(workPixels[idx].r)
			oldG := clampInt(workPixels[idx].g)
			oldB := clampInt(workPixels[idx].b)

			// Find closest palette color using lookup table
			colorKey := (uint32(oldR) << 16) | (uint32(oldG) << 8) | uint32(oldB)
			paletteIdx, exists := colorToIndex[colorKey]

			// If not in lookup, find closest color
			if !exists {
				paletteIdx = findClosestPaletteColor(palette, oldR, oldG, oldB)
			}

			indexed[idx] = paletteIdx
			newColor := palette[paletteIdx]

			// Calculate quantization error
			errR := int(oldR) - int(newColor.R)
			errG := int(oldG) - int(newColor.G)
			errB := int(oldB) - int(newColor.B)

			// Distribute error to neighboring pixels (Floyd-Steinberg matrix)
			//        X    7/16
			// 3/16  5/16  1/16

			if x+1 < width && idx+1 < len(workPixels) {
				// Right pixel: 7/16 of error
				workPixels[idx+1].r += errR * 7 / 16
				workPixels[idx+1].g += errG * 7 / 16
				workPixels[idx+1].b += errB * 7 / 16
			}

			if y+1 < height {
				nextRowIdx := idx + width
				// Bottom-left pixel: 3/16 of error
				if x > 0 && nextRowIdx-1 < len(workPixels) {
					workPixels[nextRowIdx-1].r += errR * 3 / 16
					workPixels[nextRowIdx-1].g += errG * 3 / 16
					workPixels[nextRowIdx-1].b += errB * 3 / 16
				}

				// Bottom pixel: 5/16 of error
				if nextRowIdx < len(workPixels) {
					workPixels[nextRowIdx].r += errR * 5 / 16
					workPixels[nextRowIdx].g += errG * 5 / 16
					workPixels[nextRowIdx].b += errB * 5 / 16
				}

				// Bottom-right pixel: 1/16 of error
				if x+1 < width && nextRowIdx+1 < len(workPixels) {
					workPixels[nextRowIdx+1].r += errR * 1 / 16
					workPixels[nextRowIdx+1].g += errG * 1 / 16
					workPixels[nextRowIdx+1].b += errB * 1 / 16
				}
			}
		}
	}

	return indexed
}

// clampInt clamps an integer value to byte range [0, 255]
func clampInt(val int) byte {
	if val < 0 {
		return 0
	}
	if val > 255 {
		return 255
	}
	return byte(val)
}

// findClosestPaletteColor finds the closest color in palette
func findClosestPaletteColor(palette []RGB, r, g, b byte) byte {
	minDist := uint32(0xFFFFFFFF)
	closest := byte(0)

	for i, c := range palette {
		dr := int(r) - int(c.R)
		dg := int(g) - int(c.G)
		db := int(b) - int(c.B)
		dist := uint32(dr*dr + dg*dg + db*db)

		if dist < minDist {
			minDist = dist
			closest = byte(i)
		}
	}

	return closest
}
