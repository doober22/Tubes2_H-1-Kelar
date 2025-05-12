package scraper

import (
	"encoding/json"
	"os"
	"strings"
)

type InputRecipe struct {
	Product     string     `json:"product"`
	Ingredients [][]string `json:"ingredients"`
}

type InputBlock struct {
	Title   string        `json:"title"`
	Recipes []InputRecipe `json:"recipes"`
}

type FlatRecipe struct {
	Element     string `json:"Element"`
	Ingredient1 string `json:"Ingredient1"`
	Ingredient2 string `json:"Ingredient2"`
	Tier        int    `json:"Tier"`
}

func newRecipe(element, ing1, ing2 string, tier int) FlatRecipe {
	return FlatRecipe{
		Element:     element,
		Ingredient1: ing1,
		Ingredient2: ing2,
		Tier:        tier,
	}
}

func FlattenRecipesFromFile(inputPath string) ([]FlatRecipe, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, err
	}

	var blocks []InputBlock
	if err := json.Unmarshal(data, &blocks); err != nil {
		return nil, err
	}

	var results []FlatRecipe
	elementTiers := map[string]int{
		"Air":   0,
		"Fire":  0,
		"Water": 0,
		"Earth": 0,
	}
	tier := 1

	for _, block := range blocks {
		if strings.Contains(strings.ToLower(block.Title), "special") {
			continue
		}

		var validRecipes []FlatRecipe
		for _, r := range block.Recipes {
			if strings.EqualFold(r.Product, "Time") || strings.EqualFold(r.Product, "Archeologist") || r.Ingredients == nil || len(r.Ingredients) == 0 {
				continue
			}

			for _, pair := range r.Ingredients {
				if len(pair) != 2 {
					continue
				}
				if strings.EqualFold(pair[0], "Time") || strings.EqualFold(pair[1], "Time") {
					continue
				}

				t1, ok1 := elementTiers[pair[0]]
				t2, ok2 := elementTiers[pair[1]]
				if !ok1 || !ok2 {
					continue 
				}
				if t1 >= tier || t2 >= tier {
					continue
				}
				validRecipes = append(validRecipes, newRecipe(r.Product, pair[0], pair[1], tier))
			}
		}

		if len(validRecipes) > 0 {
			results = append(results, validRecipes...)
			for _, r := range validRecipes {
				elementTiers[r.Element] = tier
			}
			tier++
		}
	}

	return results, nil
}

