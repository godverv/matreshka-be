package sqlite

import (
	"context"

	"go.redsock.ru/evon"
	errors "go.redsock.ru/rerrors"
)

func (p *Provider) GetConfigNodes(ctx context.Context, serviceName string) (*evon.Node, error) {
	row, err := p.conn.QueryContext(ctx, `
			SELECT 
				cv.key,
				cv.value
			FROM 		configs_values AS cv
			INNER JOIN 	configs 	   AS c  ON c.id = cv.config_id
			AND 		c.name 				  = $1
`, serviceName)
	if err != nil {
		return nil, errors.Wrap(err, "error getting config values")
	}
	defer row.Close()

	var rootNodes []*evon.Node

	for row.Next() {
		var node evon.Node
		err = row.Scan(&node.Name, &node.Value)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning node")
		}

		rootNodes = append(rootNodes, &node)
	}

	if len(rootNodes) == 0 {
		return nil, nil
	}

	ns := evon.NodesToStorage(rootNodes)
	return ns[""], nil
}
