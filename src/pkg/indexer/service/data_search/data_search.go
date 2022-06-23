// Package data_search keys index
package data_search

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/keystore"
)

type DataSearchIndex struct {
	index indexer.Indexer
	node  *Node
}

func New(ks *keystore.KeyStore) *DataSearchIndex {
	return &DataSearchIndex{
		index: indexer.Init("data_search"),
		node:  newNode("INDEX", ""),
	}
}

func (d *DataSearchIndex) Add(key string, value interface{}) error {
	return nil
}

type DataEntry struct {
	GroupId       *uuid.UUID
	ValueSet      map[string]interface{}
	DocStorIdEncr []byte
}

type Element struct {
	ItemType    string
	Type        string
	NodeId      string
	Name        string
	DataEntries []*DataEntry
}

type Node struct {
	Type   string
	NodeId string
	Next   map[string]*Node
	Items  map[string]*Element // nodeId -> Element
}

func newNode(_type, nodeId string) *Node {
	return &Node{
		Type:   _type,
		NodeId: nodeId,
		Next:   make(map[string]*Node),
	}
}

func (n *Node) Dump() {
	data, err := json.MarshalIndent(n, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}

func (i *DataSearchIndex) UpdateIndexWithNewContent(content interface{}, groupId *uuid.UUID, docStorId []byte) {
	if i.node == nil {
		i.node = newNode("INDEX", "")
	}

	var iterate func(items interface{}, node *Node)

	iterate = func(items interface{}, node *Node) {
		for _, item := range items.([]interface{}) {
			item := item.(map[string]interface{})

			_type := item["_type"].(string)
			itemNodeId := item["archetype_node_id"].(string)

			switch _type {
			case "SECTION":
				iterate(item["items"].([]interface{}), node)
			case "EVALUATION", "OBSERVATION", "INSTRUCTION":
				if node.Next[_type] == nil {
					node.Next[_type] = newNode(_type, itemNodeId)
				}
				nodeType := node.Next[_type]
				for _, key := range []string{"data", "protocol"} {
					if item[key] == nil {
						continue
					}

					itemsKey := item[key].(map[string]interface{})
					itemsKeyType := itemsKey["_type"].(string)
					itemsKeyNodeId := itemsKey["archetype_node_id"].(string)

					if nodeType.Next[key] == nil {
						nodeType.Next[key] = newNode(itemsKeyType, itemsKeyNodeId)
					}
					nodeCurrent := nodeType.Next[key]

					if nodeCurrent.Next[itemsKeyNodeId] == nil {
						nodeCurrent.Next[itemsKeyNodeId] = newNode(itemsKeyType, itemsKeyNodeId)
					}
					nodeCurrent = nodeCurrent.Next[itemsKeyNodeId]

					if itemsKey["items"] != nil {
						if nodeCurrent.Next["items"] == nil {
							nodeCurrent.Next["items"] = newNode("items", "")
						}
						nodeCurrent = nodeCurrent.Next["items"]
						iterate(itemsKey["items"].([]interface{}), nodeCurrent)
					}

					if itemsKey["events"] != nil {
						if nodeCurrent.Next["events"] == nil {
							nodeCurrent.Next["events"] = newNode("events", "")
						}
						nodeCurrent = nodeCurrent.Next["events"]
						iterate(itemsKey["events"].([]interface{}), nodeCurrent)
					}
				}

				if item["activities"] != nil {
					if nodeType.Next["activities"] == nil {
						nodeType.Next["activities"] = newNode("activities", "")
					}
					nodeCurrent := nodeType.Next["activities"]
					iterate(item["activities"].([]interface{}), nodeCurrent)
				}
			case "ACTION":
				if node.Next["ACTION"] == nil {
					node.Next["ACTION"] = newNode("ACTION", itemNodeId)
				}
				nodeCurrent := node.Next["ACTION"]

				if nodeCurrent.Next[itemNodeId] == nil {
					nodeCurrent.Next[itemNodeId] = newNode(_type, itemNodeId)
				}

				for _, key := range []string{"protocol", "description"} {
					if item[key] == nil {
						continue
					}
					itemsKey := item[key].(map[string]interface{})
					itemsKeyType := itemsKey["_type"].(string)
					itemsKeyNodeId := itemsKey["archetype_node_id"].(string)
					if nodeCurrent.Next[key] == nil {
						nodeCurrent.Next[key] = newNode(itemsKeyType, itemsKeyNodeId)
					}
					nodeCurrent = nodeCurrent.Next[key]

					if nodeCurrent.Next[itemsKeyNodeId] == nil {
						nodeCurrent.Next[itemsKeyNodeId] = newNode(itemsKeyType, itemsKeyNodeId)
					}
					nodeCurrent = nodeCurrent.Next[itemsKeyNodeId]

					iterate(itemsKey["items"].([]interface{}), nodeCurrent)
				}
			case "CLUSTER":
				if node.Next["CLUSTER"] == nil {
					node.Next["CLUSTER"] = newNode("CLUSTER", itemNodeId)
				}
				nodeCluster := node.Next["CLUSTER"]

				if nodeCluster.Next[itemNodeId] == nil {
					nodeCluster.Next[itemNodeId] = newNode(_type, itemNodeId)
				}
				iterate(item["items"].([]interface{}), nodeCluster.Next[itemNodeId])
			case "ACTIVITY":
				itemsDescription := item["description"].(map[string]interface{})
				itemsDescriptionType := itemsDescription["_type"].(string)
				itemsDescriptionNodeId := itemsDescription["archetype_node_id"].(string)
				if node.Next["description"] == nil {
					node.Next["description"] = newNode("description", itemsDescriptionNodeId)
				}
				nodeCurrent := node.Next["description"]

				if nodeCurrent.Next[itemsDescriptionNodeId] == nil {
					nodeCurrent.Next[itemsDescriptionNodeId] = newNode(itemsDescriptionType, itemsDescriptionNodeId)
				}
				nodeCurrent = nodeCurrent.Next[itemsDescriptionNodeId]
				iterate(itemsDescription["items"].([]interface{}), nodeCurrent)
			case "POINT_EVENT":
				itemsData := item["data"].(map[string]interface{})
				itemsDataType := itemsData["_type"].(string)
				itemsDataNodeId := itemsData["archetype_node_id"].(string)
				if node.Next["data"] == nil {
					node.Next["data"] = newNode("data", itemsDataNodeId)
				}
				nodeData := node.Next["data"]

				if nodeData.Next[itemsDataNodeId] == nil {
					nodeData.Next[itemsDataNodeId] = newNode(itemsDataType, itemsDataNodeId)
				}
				iterate(itemsData["items"].([]interface{}), nodeData.Next[itemsDataNodeId])
			case "ITEM_TREE":
				iterate(item["items"].([]interface{}), node)
			case "HISTORY":
				iterate(item["events"].([]interface{}), node)
			case "ELEMENT":
				itemValue := item["value"].(map[string]interface{})
				itemName := item["name"].(map[string]interface{})
				valueType := itemValue["_type"].(string)
				var valueSet map[string]interface{}
				switch valueType {
				case "DV_TEXT":
					valueSet = map[string]interface{}{
						"value": itemValue["value"],
					}
				case "DV_CODED_TEXT":
					defCode := itemValue["defining_code"].(map[string]interface{})
					valueSet = map[string]interface{}{
						"value":       itemValue["value"],
						"code_string": defCode["code_string"],
					}
				case "DV_IDENTIFIER":
					valueSet = map[string]interface{}{
						"id": itemValue["id"],
					}
				case "DV_MULTIMEDIA":
					valueSet = map[string]interface{}{
						"uri": itemValue["uri"],
					}
				case "DV_DATE_TIME", "DV_DATE", "DV_TIME":
					valueSet = map[string]interface{}{
						"value": itemValue["value"],
					}
				case "DV_QUANTITY":
					valueSet = map[string]interface{}{
						"magnitude": itemValue["magnitude"],
						"units":     itemValue["units"],
						"precision": itemValue["precision"],
					}
				case "DV_COUNT":
					valueSet = map[string]interface{}{
						"magnitude": itemValue["magnitude"],
					}
				case "DV_PROPORTION":
					valueSet = map[string]interface{}{
						"numerator":   itemValue["numerator"],
						"denominator": itemValue["denominator"],
						"type":        itemValue["type"],
					}
				case "DV_URI":
					valueSet = map[string]interface{}{
						"uri": itemValue["uri"],
					}
				case "DV_BOOLEAN":
					valueSet = map[string]interface{}{
						"value": itemValue["value"],
					}
				case "DV_DURATION":
					valueSet = map[string]interface{}{
						"value": itemValue["value"],
					}
				}

				if node.Items == nil {
					node.Items = make(map[string]*Element)
				}

				element, ok := node.Items[itemNodeId]
				if !ok {
					element = &Element{
						ItemType:    _type,
						Type:        valueType,
						NodeId:      itemNodeId,
						Name:        itemName["value"].(string),
						DataEntries: []*DataEntry{},
					}
					node.Items[itemNodeId] = element
				}
				dataEntry := &DataEntry{
					GroupId:       groupId,
					ValueSet:      valueSet,
					DocStorIdEncr: docStorId,
				}
				element.DataEntries = append(element.DataEntries, dataEntry)
			}
		}
	}

	iterate(content, i.node)

	i.index.Replace("INDEX", i.node)
}
