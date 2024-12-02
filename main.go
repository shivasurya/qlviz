package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/awalterschulze/gographviz"
)

type ClassNode struct {
	Name     string
	Parent   string
	Children []string
}

func main() {
	classes := make(map[string]*ClassNode)
	rootDir := "./ql/java"

	traverseFiles(rootDir, func(filePath string) {
		parseQLFile(filePath, classes)
	})

	graph := createGraph(classes)

	dotFile, err := os.Create("inheritance.dot")
	if err != nil {
		fmt.Println("Error creating DOT file:", err)
		return
	}
	defer dotFile.Close()

	dotFile.WriteString(graph.String())
}

func traverseFiles(rootDir string, processFile func(string)) {
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (filepath.Ext(path) == ".ql" || filepath.Ext(path) == ".qll") {
			processFile(path)
		}
		return nil
	})
}

var classRegex = regexp.MustCompile(`^\s*(?:abstract\s+)?class\s+(\w+)(?:\s+extends\s+([\w\.]+))?`)

func parseQLFile(filePath string, classes map[string]*ClassNode) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	content := string(data)
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		matches := classRegex.FindStringSubmatch(line)
		if matches != nil {
			className := matches[1]
			parentName := matches[2]
			classes[className] = &ClassNode{
				Name:   className,
				Parent: parentName,
			}
			// Link parent to child
			if parentName != "" {
				if parent, exists := classes[parentName]; exists {
					parent.Children = append(parent.Children, className)
				} else {
					classes[parentName] = &ClassNode{
						Name:     parentName,
						Children: []string{className},
					}
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
	}
}

func createGraph(classes map[string]*ClassNode) *gographviz.Graph {
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	gographviz.Analyse(graphAst, graph)

	for _, classNode := range classes {
		graph.AddNode("G", classNode.Name, nil)
		if classNode.Parent != "" {
			graph.AddEdge(classNode.Parent, classNode.Name, true, nil)
		}
	}
	return graph
}
