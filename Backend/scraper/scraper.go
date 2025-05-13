package scraper

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "strings"
    "github.com/gocolly/colly"
    "github.com/PuerkitoBio/goquery"
)

type Recipe struct {
    Product     string     `json:"product"`
    Ingredients [][]string `json:"ingredients"`
}

type Block struct {
    Title   string   `json:"title"`
    Recipes []Recipe `json:"recipes"`
}

func Scrape() {
    var blocks []Block
    var total int
    c := colly.NewCollector()

    c.OnHTML("div#mw-content-text", func(e *colly.HTMLElement) {
        var current *Block

        e.DOM.Find("div.mw-parser-output").Children().Each(func(i int, s *goquery.Selection) {
            tag := goquery.NodeName(s)

            if tag == "h3" {
                span := s.Find("span.mw-headline")
                if span.Length() > 0 {
                    if current != nil {
                        blocks = append(blocks, *current)
                    }
                    current = &Block{
                        Title:   strings.TrimSpace(span.Text()),
                        Recipes: []Recipe{},
                    }
                }
            } else if tag == "table" && current != nil {
                s.Find("tbody tr").Each(func(j int, tr *goquery.Selection) {
                    tds := tr.Find("td")
                    if tds.Length() >= 2 {
                        product := strings.TrimSpace(tds.Eq(0).Text())
                        var ingredients [][]string

                        ul := tds.Eq(1).Find("ul")
                        if ul.Length() > 0 {
                            ul.Find("li").Each(func(k int, li *goquery.Selection) {
                                text := strings.TrimSpace(li.Text())
                                if text != "" {
                                    separated := splitIngredients(text)
                                    ingredients = append(ingredients, separated)
                                }
                            })
                            total++
                        } else {
                            tds.Eq(1).Find("a").Each(func(k int, a *goquery.Selection) {
                                text := strings.TrimSpace(a.Text())
                                if text != "" {
                                    separated := splitIngredients(text)
                                    ingredients = append(ingredients, separated)
                                }
                            })
                            total++
                        }

                        current.Recipes = append(current.Recipes, Recipe{
                            Product:     product,
                            Ingredients: ingredients,
                        })
                    } else if tds.Length() >= 1 {
                        product := strings.TrimSpace(tds.Eq(0).Text())
                        if product != "" {
                            current.Recipes = append(current.Recipes, Recipe{
                                Product:     product,
                                Ingredients: [][]string{},
                            })
                        }
                    }
                })
            }
        })

        if current != nil {
            blocks = append(blocks, *current)
        }
    })

    fmt.Println("Memulai scraping")

    err := c.Visit("https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Scraping selesai, total produk: ",total)

    file, err := os.Create("output.json")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(blocks); err != nil {
        log.Fatal(err)
    }

    fmt.Println("Berhasil save ke JSON")
}

func splitIngredients(text string) []string {
    parts := strings.Split(text, "+")
    var result []string
    for _, part := range parts {
        trimmed := strings.TrimSpace(part)
        if trimmed != "" {
            result = append(result, trimmed)
        }
    }
    return result
}