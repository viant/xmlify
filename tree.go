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
		if id == nil { //TODO MFI how it's possible that id == nil ??? it means path == nil
			continue
		}

		index[id] = i
	}

	//for i := 0; i < len(nodes); i++ {
	//	node := nodes[i]
	//	if node.ParentID() == nil {
	//		continue
	//	}
	//
	//	_, ok := index[node.ParentID()]
	//	if i > 5 {
	//		break
	//	}
	//
	//	if !ok {
	//		nodes = append(nodes, node)
	//	}
	//}

	//indexes := NodeIndex{}

	/*
		// TODO ensure parent Object
		found := false
		for _, o := range s.objects {
			if o.parentID == nil {
				found = true
				break
			}
		}

		if !found {
			o := &Object{
				stringifier: nil, //stringifiers[field.parentType],
				objType:     nil, //???   field.parentType,
				path:        "",  //   field.path,
				parentID:    nil, //  asInterface(field.path, parentID),
				dest:        nil, // dest,
				appender:    nil, // appender,
				index: &Index{
					positionInSlice: map[interface{}]int{},
					data:            nil,
					objectIndex:     ObjectIndex{},
					appenders:       map[interface{}]*xunsafe.Appender{},
				},
				xField: nil,
				xSlice: nil,
				holder: "",
			}
			s.objects = append(s.objects, o)
		}
	*/

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

		//if !ok { // if we didn't find parent for current node, this node becomes parent
		//	parents = append(parents, nodes[i])
		//	continue
		//}

		//for ok {
		//
		//	id := parent.ID()
		//	if id == nil {
		//		break
		//	}
		//
		//	nodeIndex := indexes.Get(id)
		//	nodeID := node.ID()
		//	if nodeID == nil {
		//		break
		//	}
		//
		//	if !nodeIndex[nodeID] {
		//		nodeCopy := nodes[index[nodeID]]
		//		parent.AddChildren(nodeCopy)
		//		nodeIndex[nodeID] = true
		//		node = parent
		//		parentID := node.ParentID()
		//		if parentID == nil {
		//			ok = false
		//		} else {
		//			nodeParentIndex, ok = index[node.ParentID()]
		//		}
		//		continue
		//	}
		//
		//	break
		//}
	}

	return parents
}

func BuildTreeOLD(nodes []Node) []Node { //TODO MFI delete
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

	indexes := NodeIndex{}

	for i, node := range nodes {
		if node.ParentID() == nil {
			parents = append(parents, nodes[i])
			continue
		}

		nodeParentIndex, ok := index[node.ParentID()]
		if !ok {
			parents = append(parents, nodes[i])
			continue
		}

		for ok {
			parent := nodes[nodeParentIndex]
			id := parent.ID()
			if id == nil {
				break
			}

			nodeIndex := indexes.Get(id)
			nodeID := node.ID()
			if nodeID == nil {
				break
			}

			if !nodeIndex[nodeID] {
				nodeCopy := nodes[index[nodeID]]
				parent.AddChildren(nodeCopy)
				nodeIndex[nodeID] = true
				node = parent
				parentID := node.ParentID()
				if parentID == nil {
					ok = false
				} else {
					nodeParentIndex, ok = index[node.ParentID()]
				}
				continue
			}

			break
		}
	}

	return parents
}
