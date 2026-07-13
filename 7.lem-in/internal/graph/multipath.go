package graph

type capGraph map[string]map[string]int

func (g capGraph) cap(u, v string) int {
	if g[u] == nil {
		return 0
	}
	return g[u][v]
}

func (g capGraph) addFlow(u, v string, delta int) {
	if u == "" || v == "" {
		return
	}
	if g[u] == nil {
		g[u] = make(map[string]int)
	}
	g[u][v] += delta
}

func inNode(c *Colony, name string) string {
	if name == c.End.Name {
		return "end_in"
	}
	if name == c.Start.Name {
		return "" 
	}
	return name + "_in"
}

func outNode(c *Colony, name string) string {
	if name == c.Start.Name {
		return "start_out"
	}
	if name == c.End.Name {
		return "" 
	}
	return name + "_out"
}

type linkEdge struct{ u, v string }

func buildFlowNetwork(c *Colony) (capGraph, []linkEdge) {
	g := make(capGraph)
	var linkEdges []linkEdge

	for name := range c.Rooms {
		if name == c.Start.Name || name == c.End.Name {
			continue
		}
		g.addFlow(name+"_in", name+"_out", 1)
	}

	for name, room := range c.Rooms {
		for _, neighbor := range room.Links {
			u := outNode(c, name)
			v := inNode(c, neighbor.Name)
			if u == "" || v == "" {
				continue
			}
			g.addFlow(u, v, 1)
			linkEdges = append(linkEdges, linkEdge{u, v})
		}
	}

	return g, linkEdges
}

func bfsAugmentingPath(g capGraph, source, sink string) []string {
	visited := map[string]bool{source: true}
	cameFrom := map[string]string{}
	queue := []string{source}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == sink {
			var path []string
			for n := sink; ; n = cameFrom[n] {
				path = append([]string{n}, path...)
				if n == source {
					break
				}
			}
			return path
		}
		for next, capacity := range g[cur] {
			if capacity > 0 && !visited[next] {
				visited[next] = true
				cameFrom[next] = cur
				queue = append(queue, next)
			}
		}
	}
	return nil
}

func FindAllPaths(c *Colony) [][]*Room {
	if c.Start == nil || c.End == nil {
		return nil
	}
	g, linkEdges := buildFlowNetwork(c)
	source, sink := "start_out", "end_in"

	for {
		path := bfsAugmentingPath(g, source, sink)
		if path == nil {
			break
		}
		for i := 0; i < len(path)-1; i++ {
			u, v := path[i], path[i+1]
			g.addFlow(u, v, -1) 
			g.addFlow(v, u, 1)  
		}
	}

	nodeToRoom := func(node string) string {
		switch node {
		case "start_out":
			return c.Start.Name
		case "end_in":
			return c.End.Name
		default:
			if len(node) > 3 && node[len(node)-3:] == "_in" {
				return node[:len(node)-3]
			}
			return node[:len(node)-4] 
		}
	}

	netAdj := make(map[string][]string) 
	for _, e := range linkEdges {
		if g.cap(e.u, e.v) == 0 { 
			netAdj[nodeToRoom(e.u)] = append(netAdj[nodeToRoom(e.u)], nodeToRoom(e.v))
		}
	}

	var allPaths [][]*Room
	for len(netAdj[c.Start.Name]) > 0 {
		roomNames := []string{c.Start.Name}
		cur := c.Start.Name
		for cur != c.End.Name {
			next := netAdj[cur][0]
			netAdj[cur] = netAdj[cur][1:]
			roomNames = append(roomNames, next)
			cur = next
		}
		rooms := make([]*Room, len(roomNames))
		for i, n := range roomNames {
			rooms[i] = c.Rooms[n]
		}
		allPaths = append(allPaths, rooms)
	}

	return allPaths
}