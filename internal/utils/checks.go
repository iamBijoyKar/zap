package utils

import "github.com/iamBijoyKar/zap/internal/tasks"

func CheckDeps(completed_tasks []tasks.Task, deps_on []string) bool {
	dep_checks := make([]bool, len(deps_on))
	for idx, dep := range deps_on {
		for _, comp := range completed_tasks {
			if dep == comp.Name {
				dep_checks[idx] = true
				break
			}
		}
	}
	result := true
	for _, val := range dep_checks {
		result = result && val
	}
	return result
}
