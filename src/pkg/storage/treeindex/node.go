package treeindex

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
)

type NodeType byte

const (
	NoneNodeType NodeType = iota
	ObjectNodeType
	SliceNodeType
	DataValueNodeType
	ValueNodeType
	EHRNodeType
	CompostionNodeType
	EventContextNodeType
)

type Noder interface {
	GetNodeType() NodeType

	GetID() string

	addAttribute(key string, val Noder)
	TryGetChild(key string) Noder
}

type baseNode struct {
	ID       string        `json:"id,omitempty" msgpack:"id,omitempty"`
	Type     base.ItemType `json:"type,omitempty" msgpack:"type,omitempty"`
	Name     string        `json:"name,omitempty" msgpack:"name,omitempty"`
	NodeType NodeType      `json:"node_type" msgpack:"node_type"`
}

func (node baseNode) GetNodeType() NodeType {
	return node.NodeType
}

func (node baseNode) TryGetChild(key string) Noder {
	switch key {
	case "id":
		return newNode(node.ID)
	default:
		return nil
	}
}

type ObjectNode struct {
	baseNode

	AttributesOrder []string
	Attributes      Attributes `json:"-"`
}

func (node ObjectNode) GetID() string {
	return node.ID
}

func (node ObjectNode) TryGetChild(key string) Noder {
	n := node.baseNode.TryGetChild(key)
	if n != nil {
		return n
	}

	return node.Attributes[key]
}

func (node *ObjectNode) addAttribute(key string, val Noder) {
	node.AttributesOrder = append(node.AttributesOrder, key)
	node.Attributes[key] = val
}

func (node ObjectNode) MarshalJSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	fmt.Fprintf(buffer, "{")
	fmt.Fprintf(buffer, `"id":"%s",`, node.ID)
	fmt.Fprintf(buffer, `"name":"%s",`, node.Name)
	fmt.Fprintf(buffer, `"type":"%s"`, node.Type)

	for _, k := range node.AttributesOrder {
		data, err := json.Marshal(node.Attributes[k])
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(buffer, `,"%s":%s`, k, string(data))
	}

	fmt.Fprintf(buffer, "}")
	return buffer.Bytes(), nil
}

type SliceNode struct {
	baseNode
	Data Attributes
}

func (node SliceNode) GetID() string {
	return ""
}

func (node SliceNode) TryGetChild(key string) Noder {
	n, ok := node.Data[key]
	if !ok {
		return nil
	}

	return n
}

func (node SliceNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(node.Data)
}

func (node *SliceNode) addAttribute(key string, val Noder) {
	node.Data[val.GetID()] = val
}

type DataValueNode struct {
	baseNode
	Values Attributes `json:"values,omitempty"`
}

func (node DataValueNode) GetID() string {
	return node.ID
}

func (node DataValueNode) TryGetChild(key string) Noder {
	n := node.baseNode.TryGetChild(key)
	if n != nil {
		return n
	}

	return node.Values[key]
}

func (node *DataValueNode) addAttribute(key string, val Noder) {
	node.Values[key] = val
}

type ValueNode struct {
	baseNode
	Data any
}

func (node ValueNode) GetData() any {
	return node.Data
}

func (node ValueNode) GetID() string {
	return ""
}

func (node ValueNode) TryGetChild(key string) Noder {
	n := node.baseNode.TryGetChild(key)
	if n != nil {
		return n
	}

	return nil
}

func (node *ValueNode) addAttribute(key string, val Noder) {
	noderInstance, ok := node.Data.(Noder)
	if !ok {
		return
	}

	noderInstance.addAttribute(key, val)
}

func (node ValueNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(node.Data)
}

func newNode(obj any) Noder {
	switch obj := obj.(type) {
	case model.EHR:
		return newEHRNode(obj)
	case model.Composition:
		return newCompositionNode(obj)
	case base.Root:
		return newObjectNode(obj)
	case base.DataValue:
		return newDataValueNode(obj)
	case *base.CodePhrase:
		return nodeForCodePhrase(*obj)
	case base.CodePhrase:
		return nodeForCodePhrase(obj)
	case base.ObjectID:
		return nodeForObjectID(obj)
	case base.UIDBasedID:
		return nodeForObjectID(obj.ObjectID)
	case base.HierObjectID:
		return nodeForObjectID(obj.ObjectID)
	case base.ObjectVersionID:
		return nodeForObjectID(obj.UID.ObjectID)
	default:
		return newValueNode(obj)
	}
}

func newObjectNode(obj base.Root) Noder {
	l := obj.GetLocatable()

	return &ObjectNode{
		baseNode: baseNode{
			ID:       l.ArchetypeNodeID,
			Type:     l.Type,
			Name:     l.Name.Value,
			NodeType: ObjectNodeType,
		},
		Attributes: Attributes{
			"name":              newNode(l.Name),
			"archetype_node_id": newNode(l.ArchetypeNodeID),
		},
	}
}

func newSliceNode() Noder {
	return &SliceNode{
		baseNode: baseNode{
			NodeType: SliceNodeType,
		},
		Data: make(Attributes),
	}
}

func newDataValueNode(dv base.DataValue) Noder {
	return &DataValueNode{
		baseNode: baseNode{
			Type:     dv.GetType(),
			NodeType: DataValueNodeType,
		},
		Values: make(Attributes),
	}
}

func nodeForCodePhrase(cp base.CodePhrase) Noder {
	return &ObjectNode{
		baseNode: baseNode{
			Type:     cp.Type,
			NodeType: ObjectNodeType,
		},
		Attributes: Attributes{
			"terminology_id": nodeForObjectID(cp.TerminologyID),
			"code_string":    newValueNode(cp.CodeString),
			"preferred_term": newValueNode(cp.PreferredTerm),
		},
	}
}

func nodeForObjectID(objectID base.ObjectID) Noder {
	return &ValueNode{
		baseNode: baseNode{
			Type:     objectID.Type,
			NodeType: ValueNodeType,
		},
		Data: objectID.Value,
	}
}

func newValueNode(val any) Noder {
	return &ValueNode{
		baseNode: baseNode{
			NodeType: ValueNodeType,
		},
		Data: val,
	}
}
