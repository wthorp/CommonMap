package commonmap

const (
	MinX = 0
	MinY = 1
	MaxX = 2
	MaxY = 3
)

type Box [4]float64

type qixTree struct {
	root        *qixNode
	numFeatures int32
}

type qixNode struct {
	bbox        Box
	numFeatures int32
	FeatureIds  []int32
	numSubNodes int32
	SubNodes    [4]*qixNode
}

// Returns true if b contains a
func (b *Box) Contains(a *Box) bool {
	return a[MinX] >= b[MinX] && a[MaxX] <= b[MaxX] && a[MinY] >= b[MinY] && a[MaxY] <= b[MaxY]
}

// converts one bboxangle into two reactangles
func (in *Box) Split() (x1, x2 *Box) {
	var out1, out2 = *in, *in                          // The output bounds will be very similar to the input bounds
	if (in[MaxX] - in[MinX]) > (in[MaxY] - in[MinY]) { //Split in X direction
		width := in[MaxX] - in[MinX]
		out1[MaxX] = out1[MinX] + width*0.55
		out2[MinX] = out2[MaxX] - width*0.55
	} else { // Otherwise split in Y direction
		height := in[MaxY] - in[MinY]
		out1[MaxY] = out1[MinY] + height*0.55
		out2[MinY] = out2[MaxY] - height*0.55
	}
	return &out1, &out2
}

// Create a new quadtree
func CreateQixTree() *qixTree {
	return &qixTree{root: &qixNode{bbox: Box([4]float64{-181, -90, 181, 90})}}
}

// Insert a feature into a quadtree
func (tree *qixTree) Insert(feature int32, bbox *Box) {
	tree.numFeatures++
	qixNodeAddFeature(tree, tree.root, feature, bbox, 10)
}

// Insert a feature at a node
func qixNodeAddFeature(tree *qixTree, node *qixNode, feature int32, bbox *Box, nMaxDepth int) {
	// If there are SubNodes, then consider whether this object will fit in them.
	if nMaxDepth > 1 && node.numSubNodes > int32(0) {
		for i := int32(0); i < node.numSubNodes; i++ {
			if node.SubNodes[i].bbox.Contains(bbox) {
				qixNodeAddFeature(tree, node.SubNodes[i], feature, bbox, nMaxDepth-1)
				return
			}
		}
	} else if nMaxDepth > 1 && node.numSubNodes == int32(0) {
		// Otherwise, consider creating four SubNodes if could fit into
		// them, and adding to the appropriate SubNodes
		var half1, half2, quad1, quad2, quad3, quad4 *Box
		half1, half2 = (&node.bbox).Split()
		quad1, quad2 = half1.Split()
		quad3, quad4 = half2.Split()

		if quad1.Contains(bbox) || quad2.Contains(bbox) ||
			quad3.Contains(bbox) || quad4.Contains(bbox) {
			node.numSubNodes = 4
			node.SubNodes[0] = &qixNode{bbox: *quad1}
			node.SubNodes[1] = &qixNode{bbox: *quad2}
			node.SubNodes[2] = &qixNode{bbox: *quad3}
			node.SubNodes[3] = &qixNode{bbox: *quad4}

			/* recurse back on this node now that it has SubNodes */
			qixNodeAddFeature(tree, node, feature, bbox, nMaxDepth)
			return
		}
	}
	// If none of that worked, just add it to this nodes list.
	node.numFeatures++
	node.FeatureIds = append(node.FeatureIds, feature)
}
