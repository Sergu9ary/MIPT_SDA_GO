//go:build !solution

package hogwarts

import (
	"sort"
)

func GetCourseList(prereqs map[string][]string) []string {
	visited := make(map[string]bool)
	onStack := make(map[string]bool)
	var courses []string
	var DFS func(course string) bool
	DFS = func(course string) bool {
		if visited[course] {
			return false
		}
		if onStack[course] {
			panic("Cycle detected")
		}
		onStack[course] = true
		for _, prereq := range prereqs[course] {
			if DFS(prereq) {
				return true
			}
		}
		onStack[course] = false
		visited[course] = true
		courses = append(courses, course)
		return false
	}
	var allCourse []string
	for course := range prereqs {
		allCourse = append(allCourse, course)
	}
	sort.Strings(allCourse)
	for _, course := range allCourse {
		if DFS(course) {
			panic("Cycle detected")
		}
	}
	return courses
}
