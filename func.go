package main

func buildGraph(numbers []string) map[string][]string {
	graph := make(map[string][]string)
	for _, numStr := range numbers {
		if len(numStr) < 2 {
			continue
		}
		start := numStr[:2]
		graph[start] = append(graph[start], numStr)
	}
	return graph
}

func dfs(currNum string, graph map[string][]string, visited map[string]bool, memo map[string][]string) []string {
	if len(currNum) < 2 {
		return []string{}
	}

	if result, exists := memo[currNum]; exists {
		return result
	}

	end := currNum[len(currNum)-2:]
	if nextNums, exists := graph[end]; exists {
		var maxPath []string
		for _, nextNum := range nextNums {
			if !visited[nextNum] {
				visited[nextNum] = true
				newPath := append([]string{currNum}, dfs(nextNum, graph, visited, memo)...)
				if len(newPath) > len(maxPath) {
					maxPath = newPath
				}
				visited[nextNum] = false
			}
		}
		memo[currNum] = maxPath
		return maxPath
	}

	memo[currNum] = []string{currNum}
	return memo[currNum]
}
