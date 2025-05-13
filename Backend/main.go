package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	Element     string        `json:"element"`
	Ingredients []*RecipeNode `json:"ingredients,omitempty"`
}

type SearchRequest struct {
	Target string `json:"target"`
	Method string `json:"method"`
	Mode   string `json:"mode"`
	Limit  int    `json:"limit"`
}

type SearchResponse struct {
	Trees        []*RecipeNode `json:"trees"`
	TimeMs       float64       `json:"timeMs"`
	NodesVisited int           `json:"nodesVisited"`
}

var baseElements = map[string]bool{
	"air":   true,
	"earth": true,
	"fire":  true,
	"water": true,
}

var recipesIndex map[string][][2]string

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

func cloneVisited(original map[string]bool) map[string]bool {
	copy := make(map[string]bool)
	for k, v := range original {
		copy[k] = v
	}
	return copy
}

func buildSingleTreeDFS(element string, index map[string][][2]string, visited map[string]bool, counter *int) *RecipeNode {
	(*counter)++
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
		(*counter)++
		mu.Unlock()
	}

	var dfs func(string, map[string]bool)
	dfs = func(e string, visited map[string]bool) {
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
			mu.Unlock()
			break
		}
		mu.Unlock()
	}

	if len(trees) == 0 {
		trees = append(trees, &RecipeNode{Element: element})
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
		(*counter)++
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
						discovered[result] = &RecipeNode{Element: result, Ingredients: []*RecipeNode{left, right}}
						queue = append(queue, result)
						(*counter)++
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
	var wg sync.WaitGroup

	results := make([][2]*RecipeNode, len(ingredients))

	for i, pair := range ingredients {
		wg.Add(1)
		go func(i int, pair [2]string) {
			defer wg.Done()
			left := pair[0]
			right := pair[1]
			leftSubTree := buildRecipeTreeBFS(left, index, counter)
			rightSubTree := buildRecipeTreeBFS(right, index, counter)
			results[i] = [2]*RecipeNode{leftSubTree, rightSubTree}
		}(i, pair)
	}

	wg.Wait()
	for _, pair := range results {
		leftSubTree := pair[0]
		rightSubTree := pair[1]

		nodes[leftSubTree.Element] = leftSubTree
		nodes[rightSubTree.Element] = rightSubTree
		root.Ingredients = append(root.Ingredients, leftSubTree, rightSubTree)
	}

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
					(*counter)++
				}
				if !visited[right] {
					queue = append(queue, right)
					visited[right] = true
					(*counter)++
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

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	start := time.Now()
	counter := 0
	var trees []*RecipeNode
	target := strings.ToLower(req.Target)

	if baseElements[target] {
		json.NewEncoder(w).Encode(SearchResponse{
			Trees:        []*RecipeNode{{Element: target}},
			TimeMs:       0,
			NodesVisited: 1,
		})
		return
	}

	if _, ok := recipesIndex[target]; !ok {
		http.Error(w, "Element not available", http.StatusNotFound)
		return
	}

	switch req.Mode {
	case "single":
		if req.Method == "dfs" {
			trees = []*RecipeNode{buildSingleTreeDFS(target, recipesIndex, make(map[string]bool), &counter)}
		} else {
			trees = []*RecipeNode{bfsSingleTree(target, recipesIndex, &counter)}
		}
	case "multiple":
		if req.Method == "dfs" {
			trees = buildNRecipesDFS(target, recipesIndex, make(map[string]bool), &counter, req.Limit)
		} else {
			trees = []*RecipeNode{multiBFS(target, recipesIndex, &counter)}
		}
	default:
		http.Error(w, "Invalid mode", http.StatusBadRequest)
		return
	}
	elapsed := float64(time.Since(start).Nanoseconds()) / 1e6
	log.Printf("Returning %d trees for %s, visited %d nodes in %.2fms\n", len(trees), req.Target, counter, elapsed)
	json.NewEncoder(w).Encode(SearchResponse{
		Trees:        trees,
		TimeMs:       elapsed,
		NodesVisited: counter,
	})
}

func main() {
	recipes, err := loadRecipes("scraping/flattened.json")
	if err != nil {
		log.Fatal("Failed to load recipes:", err)
	}
	recipesIndex = indexRecipes(recipes)
	http.HandleFunc("/api/search", searchHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("✅ Server running at http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}