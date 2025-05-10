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
	tier := 0

	for _, block := range blocks {
		// Skip blocks with "Time" or "Special element" in title
		if strings.Contains(strings.ToLower(block.Title), "special") {
			continue
		}

		var validRecipes int
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
				results = append(results, newRecipe(r.Product, pair[0], pair[1], tier))
				validRecipes++
			}
		}

		if validRecipes > 0 {
			tier++
		}
	}

	return results, nil
}
