package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
	"log"
	"scraper/scraper"
)

type PathResult struct {
	Steps        []string
	NodesVisited int
	Time         time.Duration
}

type Recipe struct {
	Element     string `json:"Element"`
	Ingredient1 string `json:"Ingredient1"`
	Ingredient2 string `json:"Ingredient2"`
	Tier        int    `json:"Tier"`
}

type Result struct {
	Found bool     `json:"found"`
	Steps []string `json:"steps"`
}

var (
	recipesMap   map[string][][]string
	tierMap      map[string]int
	baseElements = map[string]bool{
		"fire": true, "water": true, "earth": true, "air": true, "time": true,
	}
	printCount = 0
	maxPrints  = 200
)

func loadRecipes(file string) ([]Recipe, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var recipes []Recipe
	err = json.Unmarshal(data, &recipes)
	return recipes, err
}

func buildRecipeMap(recipes []Recipe) {
	recipesMap = make(map[string][][]string)
	tierMap = make(map[string]int)
	for _, r := range recipes {
		element := strings.ToLower(r.Element)
		ingr1 := strings.ToLower(r.Ingredient1)
		ingr2 := strings.ToLower(r.Ingredient2)
		ingr := []string{ingr1, ingr2}
		recipesMap[element] = append(recipesMap[element], ingr)
		tierMap[element] = r.Tier
	}
}

// bfs normal
func bfsSinglePath(target string) ([]string, bool, int, time.Duration) {
	startTime := time.Now()
	visitCount := 0
	target = strings.ToLower(target)
	discovered := make(map[string]bool)
	recipeUsed := make(map[string][]string)
	for base := range baseElements {
		discovered[base] = true
		visitCount++
	}
	for {
		newDiscoveries := false
		for resultElement, recipes := range recipesMap {
			if discovered[resultElement] {
				continue
			}
			resultTier := tierMap[resultElement]
			for _, ingredients := range recipes {
				ingr1 := ingredients[0]
				ingr2 := ingredients[1]
				if discovered[ingr1] && discovered[ingr2] {
					if tierMap[ingr1] >= resultTier || tierMap[ingr2] >= resultTier {
						if printCount < maxPrints {
							fmt.Printf("Skipping recipe due to tier: %s + %s => %s (tiers: %d,%d>=%d)\n",
								ingr1, ingr2, resultElement, tierMap[ingr1], tierMap[ingr2], resultTier)
							printCount++
						}
						continue
					}

					if printCount < maxPrints {
						fmt.Printf("Discovered: %s + %s => %s\n", ingr1, ingr2, resultElement)
						printCount++
					}
					discovered[resultElement] = true
					recipeUsed[resultElement] = []string{ingr1, ingr2}
					newDiscoveries = true
					visitCount++
					if resultElement == target {
						elapsedTime := time.Since(startTime)
						return reconstructPath(target, recipeUsed), true, visitCount, elapsedTime
					}
					break
				}
			}
		}
		if !newDiscoveries {
			break
		}
	}
	elapsedTime := time.Since(startTime)
	return nil, false, visitCount, elapsedTime
}

// bfs multiple
func bfsSinglePathWithVariation(target string, seed int64) ([]string, bool, int, time.Duration) {
	startTime := time.Now()
	visitCount := 0
	target = strings.ToLower(target)
	discovered := make(map[string]bool)
	recipeUsed := make(map[string][]string)
	for base := range baseElements {
		discovered[base] = true
		visitCount++
	}
	for {
		newDiscoveries := false
		var discoveredElements []string
		for elem := range discovered {
			discoveredElements = append(discoveredElements, elem)
		}
		shuffleElements(discoveredElements, seed)
		for _, currentElement := range discoveredElements {
			for resultElement, recipes := range recipesMap {
				if discovered[resultElement] {
					continue
				}
				resultTier := tierMap[resultElement]
				shuffledRecipes := make([][]string, len(recipes))
				copy(shuffledRecipes, recipes)
				shuffleRecipes(shuffledRecipes, seed)
				for _, ingredients := range shuffledRecipes {
					ingr1 := ingredients[0]
					ingr2 := ingredients[1]
					if (ingr1 == currentElement || ingr2 == currentElement) &&
						discovered[ingr1] && discovered[ingr2] {
						if tierMap[ingr1] >= resultTier || tierMap[ingr2] >= resultTier {
							continue
						}

						discovered[resultElement] = true
						recipeUsed[resultElement] = []string{ingr1, ingr2}
						newDiscoveries = true
						visitCount++

						if resultElement == target {
							elapsedTime := time.Since(startTime)
							return reconstructPath(target, recipeUsed), true, visitCount, elapsedTime
						}

						break
					}
				}
			}
		}

		if !newDiscoveries {
			break
		}
		seed = (seed*17 + 31) % 10000
	}

	elapsedTime := time.Since(startTime)
	return nil, false, visitCount, elapsedTime
}

// dfs normal
func dfsSinglePath(target string) ([]string, bool, int, time.Duration) {
	startTime := time.Now()
	visitCount := 0
	target = strings.ToLower(target)
	discovered := make(map[string]bool)
	recipeUsed := make(map[string][]string)
	var stack []string
	for base := range baseElements {
		discovered[base] = true
		stack = append(stack, base)
		visitCount++
	}
	visited := make(map[string]bool)

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if visited[current] {
			continue
		}
		visited[current] = true
		for resultElement, recipes := range recipesMap {
			if discovered[resultElement] {
				continue
			}
			resultTier := tierMap[resultElement]
			for _, ingredients := range recipes {
				ingr1 := ingredients[0]
				ingr2 := ingredients[1]
				if discovered[ingr1] && discovered[ingr2] {
					if tierMap[ingr1] >= resultTier || tierMap[ingr2] >= resultTier {
						if printCount < maxPrints {
							fmt.Printf("Skipping recipe due to tier: %s + %s => %s (tiers: %d,%d>=%d)\n",
								ingr1, ingr2, resultElement, tierMap[ingr1], tierMap[ingr2], resultTier)
							printCount++
						}
						continue
					}
					if printCount < maxPrints {
						fmt.Printf("Discovered: %s + %s => %s\n", ingr1, ingr2, resultElement)
						printCount++
					}
					discovered[resultElement] = true
					recipeUsed[resultElement] = []string{ingr1, ingr2}
					visitCount++
					stack = append(stack, resultElement)
					if resultElement == target {
						elapsedTime := time.Since(startTime)
						return reconstructPath(target, recipeUsed), true, visitCount, elapsedTime
					}
					break
				}
			}
		}
	}
	elapsedTime := time.Since(startTime)
	return nil, false, visitCount, elapsedTime
}

// dfs multiple
func dfsSinglePathWithVariation(target string, seed int64) ([]string, bool, int, time.Duration) {
	startTime := time.Now()
	visitCount := 0
	target = strings.ToLower(target)
	discovered := make(map[string]bool)
	recipeUsed := make(map[string][]string)
	var stack []string
	for base := range baseElements {
		discovered[base] = true
		stack = append(stack, base)
		visitCount++
	}
	shuffleElements(stack, seed)
	visited := make(map[string]bool)
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if visited[current] {
			continue
		}

		visited[current] = true
		var possibleResults []string
		for resultElement, recipes := range recipesMap {
			if discovered[resultElement] {
				continue
			}
			for _, ingredients := range recipes {
				ingr1 := ingredients[0]
				ingr2 := ingredients[1]

				if discovered[ingr1] && discovered[ingr2] && (ingr1 == current || ingr2 == current) {
					resultTier := tierMap[resultElement]

					if tierMap[ingr1] >= resultTier || tierMap[ingr2] >= resultTier {
						continue
					}

					possibleResults = append(possibleResults, resultElement)
					break
				}
			}
		}
		shuffleElements(possibleResults, seed)
		for _, resultElement := range possibleResults {
			for _, ingredients := range recipesMap[resultElement] {
				ingr1 := ingredients[0]
				ingr2 := ingredients[1]

				if discovered[ingr1] && discovered[ingr2] {
					discovered[resultElement] = true
					recipeUsed[resultElement] = []string{ingr1, ingr2}
					visitCount++
					stack = append(stack, resultElement)

					if resultElement == target {
						elapsedTime := time.Since(startTime)
						return reconstructPath(target, recipeUsed), true, visitCount, elapsedTime
					}

					break
				}
			}
		}
		seed = (seed*17 + 31) % 10000
	}

	elapsedTime := time.Since(startTime)
	return nil, false, visitCount, elapsedTime
}

func findMultipleRecipes(target string, maxRecipes int, searchMethod string) []PathResult {
	target = strings.ToLower(target)
	resultChan := make(chan PathResult, maxRecipes)
	var wg sync.WaitGroup
	var mutex sync.Mutex
	foundPaths := make(map[string]bool)
	foundCount := 0
	numThreads := 8
	if maxRecipes < numThreads {
		numThreads = maxRecipes
	}
	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go func(threadID int) {
			defer wg.Done()
			seed := int64(threadID)
			for {
				mutex.Lock()
				if foundCount >= maxRecipes {
					mutex.Unlock()
					break
				}
				mutex.Unlock()
				var steps []string
				var found bool
				var visitCount int
				var elapsedTime time.Duration
				if searchMethod == "bfs" {
					steps, found, visitCount, elapsedTime = bfsSinglePathWithVariation(target, seed)
				} else {
					steps, found, visitCount, elapsedTime = dfsSinglePathWithVariation(target, seed)
				}
				if found {
					pathKey := strings.Join(steps, "|")

					mutex.Lock()
					if !foundPaths[pathKey] && foundCount < maxRecipes {
						foundPaths[pathKey] = true
						foundCount++
						resultChan <- PathResult{
							Steps:        steps,
							NodesVisited: visitCount,
							Time:         elapsedTime,
						}
					}
					mutex.Unlock()
				}
				seed = (seed*17 + 31) % 10000
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	var results []PathResult
	for result := range resultChan {
		results = append(results, result)
	}
	return results
}

func shuffleElements(elements []string, seed int64) {
	r := rand.New(rand.NewSource(seed))
	r.Shuffle(len(elements), func(i, j int) {
		elements[i], elements[j] = elements[j], elements[i]
	})
}

func shuffleRecipes(recipes [][]string, seed int64) {
	r := rand.New(rand.NewSource(seed))
	r.Shuffle(len(recipes), func(i, j int) {
		recipes[i], recipes[j] = recipes[j], recipes[i]
	})
}

func reconstructPath(target string, recipeUsed map[string][]string) []string {
	var steps []string
	var buildPath func(element string) []string
	buildPath = func(element string) []string {
		if baseElements[element] {
			return []string{}
		}

		ingredients, exists := recipeUsed[element]
		if !exists {
			return []string{}
		}
		path1 := buildPath(ingredients[0])
		path2 := buildPath(ingredients[1])

		result := append(path1, path2...)
		result = append(result, fmt.Sprintf("%s + %s = %s", ingredients[0], ingredients[1], element))
		return result
	}

	steps = buildPath(target)
	return steps
}

func main() {
	fmt.Print("Do you want to scrape (y/n)? ")
	var scrapeInput string
	fmt.Scanln(&scrapeInput)
	if strings.ToLower(scrapeInput) == "y" {
		fmt.Println("Scraping...")
		scraper.Scrape()

		fmt.Println("Scraping completed.")
		fmt.Println("Flattening recipes...")
		recipes, err := scraper.FlattenRecipesFromFile("scraping/output.json")
		if err != nil {
			log.Fatal("Error flattening recipes:", err)
		}

		outFile, err := os.Create("scraping/flattened.json")
		if err != nil {
			log.Fatal("Error creating output file:", err)
		}
		defer outFile.Close()

		encoder := json.NewEncoder(outFile)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(recipes); err != nil {
			log.Fatal("Error encoding JSON:", err)
		}

		fmt.Println("Flattened recipes saved to scraping/flattened.json")
	} else {
		fmt.Println("Skipping scraping.")
	}


	recipes, err := loadRecipes("scraping/flattened.json")
	if err != nil {
		fmt.Println("Failed to load recipes:", err)
		return
	}
	buildRecipeMap(recipes)
	var target string
	fmt.Print("Enter target element to search for: ")
	fmt.Scanln(&target)
	var searchMode string
	fmt.Print("Search for single recipe or multiple recipes? (single/multiple): ")
	fmt.Scanln(&searchMode)
	var useMethod string
	fmt.Print("Enter search method (bfs or dfs): ")
	fmt.Scanln(&useMethod)
	useMethod = strings.ToLower(useMethod)
	if useMethod != "bfs" && useMethod != "dfs" {
		fmt.Println("Invalid search method. Defaulting to bfs.")
		useMethod = "bfs"
	}
	fmt.Printf("Searching for: %s using %s\n", target, useMethod)
	if strings.ToLower(searchMode) == "multiple" {
		var maxRecipes int
		fmt.Print("Enter maximum number of recipes to find: ")
		fmt.Scanln(&maxRecipes)

		if maxRecipes <= 0 {
			fmt.Println("Invalid number. Defaulting to 5.")
			maxRecipes = 5
		}

		fmt.Printf("Searching for up to %d different recipes for %s using %s with multithreading...\n",
			maxRecipes, target, useMethod)

		startTime := time.Now()
		results := findMultipleRecipes(target, maxRecipes, useMethod)
		totalTime := time.Since(startTime)

		fmt.Printf("\nFound %d different recipes for %s in %v\n", len(results), target, totalTime)

		for i, result := range results {
			fmt.Printf("\n--- Recipe %d ---\n", i+1)
			fmt.Printf("Nodes visited: %d\n", result.NodesVisited)
			fmt.Printf("Search time: %v\n", result.Time)
			fmt.Println("Steps:")
			for j, step := range result.Steps {
				fmt.Printf("%d. %s\n", j+1, step)
			}
		}
	} else {
		var steps []string
		var found bool
		var visitCount int
		var elapsedTime time.Duration
		if useMethod == "bfs" {
			steps, found, visitCount, elapsedTime = bfsSinglePath(target)
		} else {
			steps, found, visitCount, elapsedTime = dfsSinglePath(target)
		}
		fmt.Printf("\nSearch statistics (%s):\n", useMethod)
		fmt.Printf("Nodes visited: %d\n", visitCount)
		fmt.Printf("Search time: %v\n", elapsedTime)
		if found {
			fmt.Printf("\nSteps to create %s:\n", target)
			for i, step := range steps {
				fmt.Printf("%d. %s\n", i+1, step)
			}
		} else {
			fmt.Println("\nTarget element not found.")
		}
	}
}
