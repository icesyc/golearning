package main

import (
    "fmt"
    "sort"
    "strings"
)

var prereqs = map[string][]string{
    "algorithms": {"data structures"},
    "calculus": {"linear algebra"},
    "compilers": {
        "data structures",
        "formal languages",
        "computer organization",
    },
    "intro to programming":  {"algorithms"},
    "data structures":       {"discrete math"},
    "databases":             {"data structures"},
    "discrete math":         {"intro to programming"},
    "formal languages":      {"discrete math"},
    "networks":              {"operating systems"},
    "operating systems":     {"data structures", "computer organization"},
    "programming languages": {"data structures", "computer organization"},
}

func main() {
    orders, err := topoSort(prereqs)
    if err != nil {
        fmt.Printf("error: %s\n", err)
        return
    }
    for i, course := range orders{
        fmt.Printf("%d:\t%s\n", i+1, course)
    }
}

func index(item string, arr []string) int{
    for i, v := range arr {
        if v == item {
            return i
        }
    }
    return -1
}

func topoSort(m map[string][]string) ([]string, error){
    var order []string
    resolvedMap := make(map[string]bool)
    var visitAll func(items []string, parents []string) error
    visitAll = func(items []string, parents []string) error {
        for _, item := range items {
            resolved, exists := resolvedMap[item]
            //闭环了
            if exists && !resolved {
                start := index(item, parents)
                err := fmt.Errorf("cycles: %s", strings.Join(append(parents[start:], item), " -> "))
                return err
            }
            if !exists {
                resolvedMap[item] = false
                if err := visitAll(m[item], append(parents, item)); err != nil {
                    return err
                }
                resolvedMap[item] = true
                order = append(order, item)
            }
        }
        return nil
    }
    var keys []string
    for key := range m {
        keys = append(keys, key)
    }
    sort.Strings(keys)
    err := visitAll(keys, nil)
    return order, err
}
