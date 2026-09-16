package service

import (
	"fmt"
	"sort"

	"orchestr/internal/model"
)

// TopoSort 对给定步骤和依赖做拓扑排序，返回排序后的步骤ID列表；若存在环则返回错误。
func TopoSort(steps []*model.Step, deps []*model.Dependency) ([]string, error) {
	stepIDs := make(map[string]bool)
	for _, s := range steps {
		stepIDs[s.ID] = true
	}

	adj := make(map[string][]string)
	inDegree := make(map[string]int)
	for _, s := range steps {
		adj[s.ID] = []string{}
		inDegree[s.ID] = 0
	}

	for _, d := range deps {
		if !stepIDs[d.StepID] || !stepIDs[d.DependsOnStepID] {
			continue
		}
		adj[d.DependsOnStepID] = append(adj[d.DependsOnStepID], d.StepID)
		inDegree[d.StepID]++
	}

	queue := make([]string, 0)
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}

	result := make([]string, 0, len(steps))
	for len(queue) > 0 {
		sort.Strings(queue)
		id := queue[0]
		queue = queue[1:]
		result = append(result, id)
		for _, next := range adj[id] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	if len(result) != len(steps) {
		return nil, fmt.Errorf("依赖存在循环")
	}
	return result, nil
}

// BuildDependencyGraph 构建依赖图并检测指定步骤是否已满足依赖。
func BuildDependencyGraph(deps []*model.Dependency) map[string][]string {
	graph := make(map[string][]string)
	for _, d := range deps {
		graph[d.StepID] = append(graph[d.StepID], d.DependsOnStepID)
	}
	return graph
}

// AreDependenciesSatisfied 检查某步骤的所有依赖是否都已完成。
func AreDependenciesSatisfied(stepID string, graph map[string][]string, completed map[string]bool) bool {
	deps, ok := graph[stepID]
	if !ok {
		return true
	}
	for _, dep := range deps {
		if !completed[dep] {
			return false
		}
	}
	return true
}
