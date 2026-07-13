package graph

func FindShortestPath(c *Colony) []*Room {
	if c.Start == nil || c.End == nil {
		return nil
	}

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
			if !visited[neighbor] {
				visited[neighbor] = true
				cameFrom[neighbor] = current
				queue = append(queue, neighbor)
			}
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