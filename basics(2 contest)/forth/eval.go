//go:build !solution

package main

import (
	"regexp/syntax"
	"strconv"
	"strings"
)

func checkArgCount(stack []int, needed int) bool {
	return len(stack) >= needed
}

func sum(stack []int) ([]int, error) {
	if !checkArgCount(stack, 2) {
		return stack, &syntax.Error{}
	}
	f, s := stack[len(stack)-1], stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	stack = append(stack, f+s)
	return stack, nil
}

func div(stack []int) ([]int, error) {
	if !checkArgCount(stack, 2) {
		return stack, &syntax.Error{}
	}
	f, s := stack[len(stack)-1], stack[len(stack)-2]
	if f == 0 {
		return stack, &syntax.Error{}
	}
	stack = stack[:len(stack)-2]
	stack = append(stack, s/f)
	return stack, nil
}

func sub(stack []int) ([]int, error) {
	if !checkArgCount(stack, 2) {
		return stack, &syntax.Error{}
	}
	f, s := stack[len(stack)-1], stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	stack = append(stack, s-f)
	return stack, nil
}

func mult(stack []int) ([]int, error) {
	if !checkArgCount(stack, 2) {
		return stack, &syntax.Error{}
	}
	f, s := stack[len(stack)-1], stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	stack = append(stack, f*s)
	return stack, nil
}

func drop(stack []int) ([]int, error) {
	if !checkArgCount(stack, 1) {
		return stack, &syntax.Error{}
	}
	stack = stack[:len(stack)-1]
	return stack, nil
}

func swap(stack []int) ([]int, error) {
	if !checkArgCount(stack, 2) {
		return stack, &syntax.Error{}
	}
	f, s := stack[len(stack)-1], stack[len(stack)-2]
	stack = append(stack[:len(stack)-2], f, s)
	return stack, nil
}

func dup(stack []int) ([]int, error) {
	if !checkArgCount(stack, 1) {
		return stack, &syntax.Error{}
	}
	f := stack[len(stack)-1]
	stack = append(stack, f)
	return stack, nil
}

func over(stack []int) ([]int, error) {
	if !checkArgCount(stack, 2) {
		return stack, &syntax.Error{}
	}
	s := stack[len(stack)-2]
	stack = append(stack, s)
	return stack, nil
}

type Evaluator struct {
	stack []int
	mp    map[string]func([]int) ([]int, error)
	redis map[string]string
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	eval := Evaluator{
		mp: map[string]func([]int) ([]int, error){
			"+":    sum,
			"-":    sub,
			"*":    mult,
			"/":    div,
			"drop": drop,
			"dup":  dup,
			"swap": swap,
			"over": over,
		},
		redis: map[string]string{},
	}
	return &eval
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.

func (e *Evaluator) AddFunc(key string, processes []string) error {
	if _, err := strconv.Atoi(key); err == nil {
		return &syntax.Error{}
	}
	_, ok := e.mp[key]
	if ok {
		newKey := key + "l1"
		e.redis[key] = newKey
		key = newKey
	}
	e.mp[key] = func(stack []int) ([]int, error) {
		for _, proc := range processes {
			if num, err := strconv.Atoi(proc); err == nil {
				stack = append(stack, num)
			} else {
				newSt, err := e.mp[proc](stack)
				if err != nil {
					return stack, err
				}
				stack = newSt
			}
		}
		return stack, nil
	}
	return nil
}

func (e *Evaluator) Compute(processes []string) error {
	for _, proc := range processes {
		if num, err := strconv.Atoi(proc); err == nil {
			e.stack = append(e.stack, num)
		} else {
			f, ok := e.mp[proc]
			if !ok {
				return &syntax.Error{}
			}
			stack, err := f(e.stack)
			if err != nil {
				return err
			}
			e.stack = stack
		}
	}
	return nil
}

func (e *Evaluator) Redef(processes []string) []string {
	for ind, proc := range processes {
		if newProc, ok := e.redis[proc]; ok {
			processes[ind] = newProc
		}
	}
	return processes
}

func (e *Evaluator) Process(row string) ([]int, error) {
	processes := strings.Split(strings.ToLower(row), " ")
	processes = e.Redef(processes)
	err := error(nil)
	if processes[0] == ":" {
		err = e.AddFunc(processes[1], processes[2:len(processes)-1])
	} else {
		err = e.Compute(processes)
	}
	return e.stack, err
}
