package graph

type edgeKey struct {
	a, b *Room
}

func makeEdgeKey(a, b *Room) edgeKey {
	if a.Name > b.Name {
		a, b = b, a
	}
	return edgeKey{a, b}
}

func findPathAvoiding(c *Colony, used map[edgeKey]bool) []*Room {
	visited := make(map[*Room]bool)
	cameFrom := make(map[*Room]*Room)

	queue := []*Room{c.Start}
	visited[c.Start] = true

	found := false
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == c.End {
			found = true
			break
		}

		for _, neighbor := range current.Links {
			if visited[neighbor] {
				continue
			}
			if used[makeEdgeKey(current, neighbor)] {
				continue
			}
			visited[neighbor] = true
			cameFrom[neighbor] = current
			queue = append(queue, neighbor)
		}
	}

	if !found {
		return nil
	}

	var path []*Room
	for r := c.End; r != nil; r = cameFrom[r] {
		path = append(path, r)
		if r == c.Start {
			break
		}
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func FindAllPaths(c *Colony) [][]*Room {
	used := make(map[edgeKey]bool)
	var allPaths [][]*Room

	for {
		path := findPathAvoiding(c, used)
		if path == nil {
			break
		}
		allPaths = append(allPaths, path)

		for i := 0; i < len(path)-1; i++ {
			used[makeEdgeKey(path[i], path[i+1])] = true
		}
	}

	return allPaths
}
