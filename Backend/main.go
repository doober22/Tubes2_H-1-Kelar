package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type FlatRecipe struct {
	Element     string `json:"Element"`
	Ingredient1 string `json:"Ingredient1"`
	Ingredient2 string `json:"Ingredient2"`
	Tier        int    `json:"Tier"`
}

type RecipeNode struct {
	Element     string
	Ingredients []*RecipeNode
}

var baseElements = map[string]bool{
	"air":   true,
	"earth": true,
	"fire":  true,
	"water": true,
}

func loadRecipes(path string) ([]FlatRecipe, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var recipes []FlatRecipe
	err = json.Unmarshal(data, &recipes)
	if err != nil {
		return nil, err
	}

	for i := range recipes {
		recipes[i].Element = strings.ToLower(recipes[i].Element)
		recipes[i].Ingredient1 = strings.ToLower(recipes[i].Ingredient1)
		recipes[i].Ingredient2 = strings.ToLower(recipes[i].Ingredient2)
	}
	return recipes, nil
}

func indexRecipes(recipes []FlatRecipe) map[string][][2]string {
	index := make(map[string][][2]string)
	for _, r := range recipes {
		index[r.Element] = append(index[r.Element], [2]string{r.Ingredient1, r.Ingredient2})
	}
	return index
}

func buildSingleTreeDFS(element string, index map[string][][2]string, visited map[string]bool, counter *int) *RecipeNode {
	*counter++
	if baseElements[element] {
		return &RecipeNode{Element: element}
	}
	if visited[element] {
		return nil
	}
	visited[element] = true

	var tree *RecipeNode
	if recipes, ok := index[element]; ok {
		for _, pair := range recipes {
			left := buildSingleTreeDFS(pair[0], index, cloneVisited(visited), counter)
			right := buildSingleTreeDFS(pair[1], index, cloneVisited(visited), counter)
			if left != nil && right != nil {
				tree = &RecipeNode{Element: element, Ingredients: []*RecipeNode{left, right}}
				break
			}
		}
	}
	if tree == nil {
		tree = &RecipeNode{Element: element}
	}
	return tree
}

func buildNRecipesDFS(element string, index map[string][][2]string, visited map[string]bool, counter *int, maxRecipes int) []*RecipeNode {
	var trees []*RecipeNode
	var mu sync.Mutex
	var wg sync.WaitGroup
	resultChan := make(chan *RecipeNode, maxRecipes)
	done := make(chan struct{})
	var once sync.Once

	if counter != nil {
		mu.Lock()
		*counter++
		mu.Unlock()
	}

	var dfs func(string, map[string]bool) = func(e string, visited map[string]bool) {
		defer wg.Done()

		select {
		case <-done:
			return
		default:
		}

		if baseElements[e] || visited[e] {
			select {
			case resultChan <- &RecipeNode{Element: e}:
			case <-done:
			}
			return
		}

		visited[e] = true

		if recipes, ok := index[e]; ok {
			for _, pair := range recipes {
				select {
				case <-done:
					return
				default:
				}

				leftCh := make(chan []*RecipeNode, 1)
				rightCh := make(chan []*RecipeNode, 1)

				wg.Add(2)
				go func(p string, v map[string]bool) {
					defer wg.Done()
					leftCh <- buildNRecipesDFS(p, index, cloneVisited(v), counter, maxRecipes)
				}(pair[0], visited)

				go func(p string, v map[string]bool) {
					defer wg.Done()
					rightCh <- buildNRecipesDFS(p, index, cloneVisited(v), counter, maxRecipes)
				}(pair[1], visited)

				lefts := <-leftCh
				rights := <-rightCh

				for _, l := range lefts {
					for _, r := range rights {
						select {
						case <-done:
							return
						default:
							select {
							case resultChan <- &RecipeNode{Element: e, Ingredients: []*RecipeNode{l, r}}:
							case <-done:
								return
							}
						}
					}
				}
			}
		}
	}

	wg.Add(1)
	go dfs(element, cloneVisited(visited))

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for recipe := range resultChan {
		mu.Lock()
		if len(trees) < maxRecipes {
			trees = append(trees, recipe)
		}
		if len(trees) >= maxRecipes {
			once.Do(func() { close(done) })
		}
		mu.Unlock()
	}

	return trees
}

func bfsSingleTree(element string, index map[string][][2]string, counter *int) *RecipeNode {
	element = strings.ToLower(element)
	discovered := make(map[string]*RecipeNode)

	queue := []string{}
	for base := range baseElements {
		discovered[base] = &RecipeNode{Element: base}
		queue = append(queue, base)
		*counter++
	}

	for len(queue) > 0 {
		currentSize := len(queue)
		for i := 0; i < currentSize; i++ {
			for result, recipes := range index {
				if _, already := discovered[result]; already {
					continue
				}
				for _, pair := range recipes {
					left, okL := discovered[pair[0]]
					right, okR := discovered[pair[1]]
					if okL && okR {
						discovered[result] = &RecipeNode{
							Element:     result,
							Ingredients: []*RecipeNode{left, right},
						}
						queue = append(queue, result)
						*counter++
						break
					}
				}
			}
		}
		queue = queue[currentSize:]
		if _, found := discovered[element]; found {
			break
		}
	}

	if tree, ok := discovered[element]; ok {
		return tree
	}
	return &RecipeNode{Element: element}
}




func multiBFS(element string, index map[string][][2]string, counter *int) *RecipeNode {

	root := &RecipeNode{Element: element}
	nodes := map[string]*RecipeNode{element: root}
	ingredients := make([][2]string, 0)

	if recipes, ok := index[element]; ok {
		ingredients = append(ingredients, recipes...)
	}

	buildSubTree := func(pair [2]string, wg *sync.WaitGroup) {
		defer wg.Done()

		left := pair[0]
		right := pair[1]

		leftSubTree := buildRecipeTreeBFS(left, index, counter)
		rightSubTree := buildRecipeTreeBFS(right, index, counter)

		nodes[left] = leftSubTree
		nodes[right] = rightSubTree

		if parentNode, exists := nodes[element]; exists {
			parentNode.Ingredients = append(parentNode.Ingredients, leftSubTree, rightSubTree)
		}
	}
	var wg sync.WaitGroup

	for _, pair := range ingredients {
		wg.Add(1)
		go buildSubTree(pair, &wg)
	}

	wg.Wait()
	return root
}



func buildRecipeTreeBFS(element string, index map[string][][2]string, counter *int) *RecipeNode {
	queue := []string{element}
	visited := make(map[string]bool)
	visited[element] = true

	root := &RecipeNode{Element: element}
	nodes := map[string]*RecipeNode{element: root}

	for len(queue) > 0 {
		currentElement := queue[0]
		queue = queue[1:]
		if recipes, ok := index[currentElement]; ok {
			for _, pair := range recipes {
				left := pair[0]
				right := pair[1]

				if !visited[left] {
					queue = append(queue, left)
					visited[left] = true
					*counter++
				}
				if !visited[right] {
					queue = append(queue, right)
					visited[right] = true
					*counter++ 
				}
				leftNode, leftExists := nodes[left]
				if !leftExists {
					leftNode = &RecipeNode{Element: left}
					nodes[left] = leftNode
				}

				rightNode, rightExists := nodes[right]
				if !rightExists {
					rightNode = &RecipeNode{Element: right}
					nodes[right] = rightNode
				}
				nodes[currentElement].Ingredients = append(nodes[currentElement].Ingredients, leftNode, rightNode)
			}
		}
	}
	return root
}


func cloneVisited(original map[string]bool) map[string]bool {
	copy := make(map[string]bool)
	for k, v := range original {
		copy[k] = v
	}
	return copy
}

func printTree(node *RecipeNode, prefix, childPrefix string) {
	fmt.Println(prefix + node.Element)
	for i, child := range node.Ingredients {
		if i == len(node.Ingredients)-1 {
			printTree(child, childPrefix+"└── ", childPrefix+"    ")
		} else {
			printTree(child, childPrefix+"├── ", childPrefix+"│   ")
		}
	}
}

func main() {
	recipes, err := loadRecipes("scraping/flattened.json")
	if err != nil {
		fmt.Println("Error loading recipes:", err)
		return
	}

	index := indexRecipes(recipes)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter target element: ")
	scanner.Scan()
	target := strings.ToLower(strings.TrimSpace(scanner.Text()))

	fmt.Print("Choose recipe type (single/multiple): ")
	scanner.Scan()
	recipeType := strings.ToLower(strings.TrimSpace(scanner.Text()))

	var algo string
	var maxRecipes int
	switch recipeType {
	case "single", "multiple":
		fmt.Print("Choose algorithm (dfs/bfs): ")
		scanner.Scan()
		algo = strings.ToLower(strings.TrimSpace(scanner.Text()))
	default:
		fmt.Println("Invalid recipe type. Use 'single' or 'multiple'.")
		return
	}

	if recipeType == "multiple" && algo == "dfs" {
		fmt.Print("Enter the number of recipes to find (N): ")
		scanner.Scan()
		fmt.Sscanf(scanner.Text(), "%d", &maxRecipes)
	}

	switch recipeType {
	case "single":
		switch algo {
		case "bfs":
			fmt.Println("\n[Single Recipe Tree - BFS]")
			start := time.Now()
			counter := 0
			rootTree := bfsSingleTree(target, index, &counter)
			duration := time.Since(start)

			fmt.Printf("\nRecipe Tree for %s:\n", target)
			printTree(rootTree, "", "")

			fmt.Printf("\nNodes visited: %d\n", counter)
			fmt.Printf("Time taken: %dms\n", duration.Milliseconds())

		case "dfs":
			fmt.Println("\n[Single Recipe Tree - DFS]")
			start := time.Now()
			counter := 0
			rootTree := buildSingleTreeDFS(target, index, make(map[string]bool), &counter)
			duration := time.Since(start)

			fmt.Printf("\nRecipe Tree for %s:\n", target)
			printTree(rootTree, "", "")

			fmt.Printf("\nNodes visited: %d\n", counter)
			fmt.Printf("Time taken: %dms\n", duration.Milliseconds())

		default:
			fmt.Println("Invalid algorithm. Use 'dfs' or 'bfs'.")
		}

	case "multiple":
		switch algo {
		case "bfs":
			fmt.Println("\n[Multiple Recipe Tree - BFS]")
			start := time.Now()
			counter := 0
			rootTree := multiBFS(target, index, &counter)
			duration := time.Since(start)

			fmt.Printf("\nRecipe Tree for %s:\n", target)
			printTree(rootTree, "", "")

			fmt.Printf("\nNodes visited: %d\n", counter)
			fmt.Printf("Time taken: %dms\n", duration.Milliseconds())

		case "dfs":
			fmt.Println("\n[Multiple Recipe Trees - DFS]")
			start := time.Now()
			counter := 0
			trees := buildNRecipesDFS(target, index, make(map[string]bool), &counter, maxRecipes)
			duration := time.Since(start)

			for i, tree := range trees {
				fmt.Printf("\nRecipe #%d:\n", i+1)
				printTree(tree, "", "")
			}

			fmt.Printf("\nTotal recipes: %d\n", len(trees))
			fmt.Printf("Nodes visited: %d\n", counter)
			fmt.Printf("Time taken: %dms\n", duration.Milliseconds())

		default:
			fmt.Println("Invalid algorithm. Use 'dfs' or 'bfs'.")
		}
	}
}
