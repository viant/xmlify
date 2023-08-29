package xmlify

type (
	NodeIndex map[interface{}]map[interface{}]bool
	Node      interface {
		ID() interface{}
		ParentID() interface{}
		AddChildren(node Node)
	}
)

func ParentOf(objects []*Object) (*Object, bool) {
	nodes := make([]Node, 0, len(objects))
	for i, _ := range objects {
		nodes = append(nodes, objects[i])
	}

	parents := BuildTree(nodes)
	for _, parent := range parents {
		if parent.ID() == "" {
			asObj, ok := parent.(*Object)
			if ok {
				return asObj, true
			}
		}
	}

	return nil, false
}

func (i NodeIndex) Get(id interface{}) map[interface{}]bool {
	index, ok := i[id]
	if !ok {
		index = map[interface{}]bool{}
		i[id] = index
	}

	return index
}

func BuildTree(nodes []Node) []Node {
	if len(nodes) == 0 {
		return []Node{}
	}

	var parents []Node
	index := map[interface{}]int{}

	for i, node := range nodes {
		id := node.ID()
		if id == nil {
			continue
		}
		index[id] = i
	}

	for i, node := range nodes {
		if node.ParentID() == nil {
			parents = append(parents, nodes[i])
			continue
		}

		nodeParentIndex, ok := index[node.ParentID()]
		if ok {
			parent := nodes[nodeParentIndex]
			parent.AddChildren(node)
		}
	}

	return parents
}
