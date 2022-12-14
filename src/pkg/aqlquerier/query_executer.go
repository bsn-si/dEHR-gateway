package aqlquerier

import (
	"database/sql/driver"
	"fmt"
	"log"
	"os"

	"hms/gateway/pkg/aqlprocessor"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/storage/treeindex"
)

type executer struct {
	query  *aqlprocessor.Query
	params map[string]driver.Value

	index *treeindex.Tree
}

func (exec *executer) run() (*Rows, error) {
	// handle FROM block
	dataSources, err := exec.findSources()
	if err != nil {
		return nil, err
	}

	// handle SELECT block
	rows, err := exec.queryData(dataSources)
	if err != nil {
		return nil, err
	}

	// handle ORDER block
	//TODO: add order logic

	return rows, nil
}

func (exec *executer) findSources() (map[string]dataSource, error) {
	from := exec.query.From
	if len(from.Contains) > 0 || from.Operator != nil {
		return nil, errors.New("not implemented")
	}

	source := map[string]dataSource{}

	if from.Operand != nil {
		switch operand := from.Operand.(type) {
		case aqlprocessor.ClassExpression:
			if operand.PathPredicate != nil {
				return nil, errors.New("not implemented")
			}

			ds := dataSource{
				name: operand.Identifiers[0],
			}

			if len(operand.Identifiers) > 1 {
				ds.alias = operand.Identifiers[1]
			}

			data, err := exec.index.GetDataSourceByName(ds.name)
			if err != nil {
				return nil, errors.Wrap(err, "cannot get data source by name")
			}

			ds.data = data

			source[ds.getName()] = ds
		case aqlprocessor.VersionClassExpr:
			return nil, errors.New("not implemented")
		default:
			return nil, fmt.Errorf("unexpected FROM.Operand type: %T", operand) // nolint
		}
	}

	return source, nil
}

func (exec *executer) queryData(sources map[string]dataSource) (*Rows, error) {
	rows := &Rows{}

	//TODO: add DISTINCT handling
	// exec.query.Select.Distinct

	// for _, source := range sources {
	// log.Println(source.name, source.alias, source.data)

	// row := Row{}

	primitives := Row{}

	for _, selectExpr := range exec.query.Select.SelectExprs {
		switch slct := selectExpr.Value.(type) {
		case *aqlprocessor.IdentifiedPathSelectValue:
			columnValues, err := exec.getDataByIdentifiedPath(slct, sources)
			if err != nil {
				return nil, errors.Wrap(err, "cannot get data fo~r identified data")
			}

			for _, val := range columnValues {
				row := Row{
					values: []interface{}{val},
				}
				rows.rows = append(rows.rows, row)
			}
		case *aqlprocessor.PrimitiveSelectValue:
			val, err := exec.getPrimitiveColumnValue(slct)
			if err != nil {
				return nil, errors.Wrap(err, "cannot get primitive select value")
			}

			primitives.values = append(primitives.values, val)
		case *aqlprocessor.AggregateFunctionCallSelectValue:
			return nil, errors.New("Aggregation function call not implemented")
		case *aqlprocessor.FunctionCallSelectValue:
			return nil, errors.New("Function call not implemented")
		default:
			return nil, errors.New("Unexpected SelectExpr type")
		}
		// }
	}

	rows.rows = append(rows.rows, primitives)

	return exec.fillColumns(rows), nil
}

func (exec *executer) getPrimitiveColumnValue(prim *aqlprocessor.PrimitiveSelectValue) (driver.Value, error) {
	if prim == nil {
		return nil, nil
	}

	return prim.Val.Val, nil
}

func (exec *executer) fillColumns(rows *Rows) *Rows {
	for _, se := range exec.query.Select.SelectExprs {
		rows.columns = append(rows.columns, se.AliasName)
	}

	return rows
}

func (exec *executer) getDataByIdentifiedPath(slctExpr *aqlprocessor.IdentifiedPathSelectValue, sources map[string]dataSource) ([]any, error) {
	result := []any{}
	logger := log.New(os.Stderr, "\t[getDataByIdentifiedPath]\t", log.LstdFlags)
	logger.Println()

	selectPath := slctExpr.Val

	source, ok := sources[selectPath.Identifier]
	if !ok {
		return nil, errors.New("unexpected identifier " + selectPath.Identifier)
	}

	for _, indexNodes := range source.data {
		for _, indexNode := range indexNodes {
			if selectPath.ObjectPath != nil {
				if resultData, ok := getValueForPath(selectPath.ObjectPath, indexNode); ok {
					logger.Printf("%v = %T", resultData, resultData)
					result = append(result, resultData)
				}
			}
		}
	}

	return result, nil
}

func getValueForPath(path *aqlprocessor.ObjectPath, node treeindex.Noder) (any, bool) {
	index := 0
	queue := []treeindex.Noder{node}

	for len(queue) > 0 {
		if index >= len(path.Paths) {
			return nil, false
		}

		path := path.Paths[index]

		node := queue[0]
		queue = queue[1:]

		switch node := node.(type) {
		case *treeindex.ObjectNode:
			{
				nextNode := node.TryGetChild(path.Identifier)
				if nextNode == nil {
					continue
				}

				switch nextNode.GetNodeType() {
				case treeindex.ObjectNodeType:
					if path.PathPredicate != nil && path.PathPredicate.Type == aqlprocessor.NodePathPredicate {
						if np := path.PathPredicate.NodePredicate; np.AtCode != nil && nextNode.GetID() == np.AtCode.ToString() {
							index++
						}
					}
				case treeindex.DataValueNodeType:
					index++
				}

				queue = append(queue, nextNode)
			}
		case *treeindex.SliceNode:
			if path.PathPredicate != nil && path.PathPredicate.Type == aqlprocessor.NodePathPredicate {
				np := path.PathPredicate.NodePredicate

				if np.AtCode != nil {
					nextNode := node.TryGetChild(np.AtCode.ToString())
					if nextNode != nil {
						queue = append(queue, nextNode)
						index++
					}
				}
			}
		case *treeindex.DataValueNode:
			if valueNode := node.TryGetChild(path.Identifier); valueNode != nil {
				queue = append(queue, valueNode)
			}

		case *treeindex.ValueNode:
			return node.GetData(), true
		default:
		}
	}

	return nil, false
}

type dataSource struct {
	name  string
	alias string
	data  treeindex.Container
}

func (ds dataSource) getName() string {
	if ds.alias != "" {
		return ds.alias
	}

	return ds.name
}
