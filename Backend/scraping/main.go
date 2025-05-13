package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"scraper/scraper"
)

func main() {


		fmt.Println("Scraping...")


		scraper.Scrape()





		fmt.Println("Scraping completed.")


		fmt.Println("Flattening recipes...")


		recipes, err := scraper.FlattenRecipesFromFile("output.json")


		if err != nil {


			log.Fatal("Error flattening recipes:", err)


		}





		outFile, err := os.Create("flattened.json")


		if err != nil {


			log.Fatal("Error creating output file:", err)


		}


		defer outFile.Close()





		encoder := json.NewEncoder(outFile)


		encoder.SetIndent("", "  ")


		if err := encoder.Encode(recipes); err != nil {


			log.Fatal("Error encoding JSON:", err)


		}





		fmt.Println("Flattened recipes saved to flattened.json")
}