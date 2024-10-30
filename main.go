package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type Trait struct {
	Name      string
	Dominant  bool
	GeneLabel string
}

type Plant struct {
	Name   string
	Traits []Trait
}

type TraitCombination struct {
	Traits       []bool
	Probability  int
	Denominator  int
	Description  string
	GeneNotation string
	Expected     float64
}

type SummaryStats struct {
	TotalProbabilityNum   int
	TotalProbabilityDenom int
	Percentage            float64
	ExpectedPlants        float64
}

type TargetGenotype struct {
	Genotype    string
	Description string
}

func isGenotypeMatch(genotype string, geneNotation string) bool {
	// Handle recessive case (e.g., "ll")
	if len(genotype) == 2 && genotype[0] == genotype[1] {
		return strings.Contains(geneNotation, genotype)
	}

	// Handle dominant case (e.g., "L")
	// Match either "L_" or "LL"
	if len(genotype) == 1 {
		dominantPattern := genotype + "_"
		homozygousPattern := strings.ToUpper(genotype + genotype)
		return strings.Contains(geneNotation, dominantPattern) ||
			strings.Contains(geneNotation, homozygousPattern)
	}

	return false
}

func getGenotypeDescription(genotype string, plant1, plant2 Plant) string {
	for i, trait := range plant1.Traits {
		// Handle recessive case
		lowerGenotype := strings.ToLower(trait.GeneLabel + trait.GeneLabel)
		if genotype == lowerGenotype {
			if plant1.Traits[i].Dominant {
				return plant2.Traits[i].Name
			} else {
				return plant1.Traits[i].Name
			}
		}

		// Handle dominant case
		if genotype == trait.GeneLabel {
			if plant1.Traits[i].Dominant {
				return plant1.Traits[i].Name
			} else {
				return plant2.Traits[i].Name
			}
		}
	}
	return genotype
}

func calculateDenominator(numTraits int) int {
	return int(math.Pow(2, float64(2*numTraits)))
}

func calculateProbability(numTraits, numRecessive int) int {
	dominantTraits := numTraits - numRecessive
	return int(math.Pow(3, float64(dominantTraits)))
}

func calculateF2Probabilities(plant1, plant2 Plant, totalPlants int) []TraitCombination {
	numTraits := len(plant1.Traits)
	if numTraits != len(plant2.Traits) {
		panic("Plants must have same number of traits")
	}

	denominator := calculateDenominator(numTraits)
	numCombinations := int(math.Pow(2, float64(numTraits)))
	combinations := make([]TraitCombination, 0)

	for i := 0; i < numCombinations; i++ {
		traits := make([]bool, numTraits)
		description := make([]string, 0)
		geneNotation := make([]string, 0)
		numRecessive := 0

		for j := 0; j < numTraits; j++ {
			traits[j] = (i & (1 << j)) == 0
			if !traits[j] {
				numRecessive++
			}

			var traitName string
			if traits[j] {
				if plant1.Traits[j].Dominant {
					traitName = plant1.Traits[j].Name
				} else {
					traitName = plant2.Traits[j].Name
				}
				geneNotation = append(geneNotation, plant1.Traits[j].GeneLabel+"_")
			} else {
				if plant1.Traits[j].Dominant {
					traitName = plant2.Traits[j].Name
				} else {
					traitName = plant1.Traits[j].Name
				}
				geneNotation = append(geneNotation, strings.ToLower(plant1.Traits[j].GeneLabel+plant1.Traits[j].GeneLabel))
			}
			description = append(description, traitName)
		}

		probability := calculateProbability(numTraits, numRecessive)
		expected := float64(probability) * float64(totalPlants) / float64(denominator)

		combinations = append(combinations, TraitCombination{
			Traits:       traits,
			Probability:  probability,
			Denominator:  denominator,
			Description:  strings.Join(description, ", "),
			GeneNotation: strings.Join(geneNotation, " "),
			Expected:     expected,
		})
	}

	sort.Slice(combinations, func(i, j int) bool {
		return combinations[i].Probability > combinations[j].Probability
	})

	return combinations
}

func filterCombinations(combinations []TraitCombination, targetGenotypes []TargetGenotype) ([]TraitCombination, SummaryStats) {
	// Get denominator from the first combination
	totalCombinations := combinations[0].Denominator

	if len(targetGenotypes) == 0 {
		return combinations, SummaryStats{
			TotalProbabilityNum:   totalCombinations,
			TotalProbabilityDenom: totalCombinations,
			Percentage:            100.0,
			ExpectedPlants:        float64(totalCombinations),
		}
	}

	filtered := make([]TraitCombination, 0)
	totalProb := 0

	for _, combo := range combinations {
		matches := true
		for _, target := range targetGenotypes {
			if !isGenotypeMatch(target.Genotype, combo.GeneNotation) {
				matches = false
				break
			}
		}
		if matches {
			filtered = append(filtered, combo)
			totalProb += combo.Probability
		}
	}

	summary := SummaryStats{
		TotalProbabilityNum:   totalProb,
		TotalProbabilityDenom: totalCombinations,
		Percentage:            float64(totalProb) / float64(totalCombinations) * 100.0,
		ExpectedPlants:        float64(totalProb) / float64(totalCombinations) * (combinations[0].Expected * float64(totalCombinations) / float64(combinations[0].Probability)),
	}

	return filtered, summary
}

func main() {
	purpleFlash := Plant{
		Name: "Purple Flash",
		Traits: []Trait{
			{
				Name:      "Regular Leave Shape",
				Dominant:  true,
				GeneLabel: "L",
			},
			{
				Name:      "Purple Foliage",
				Dominant:  false,
				GeneLabel: "C",
			},
			{
				Name:      "Round Fruit Shape",
				Dominant:  false,
				GeneLabel: "F",
			},
		},
	}

	candlelight := Plant{
		Name: "Candlelight Mutant",
		Traits: []Trait{
			{
				Name:      "Mutant Leave Shape",
				Dominant:  false,
				GeneLabel: "L",
			},
			{
				Name:      "Green Foliage",
				Dominant:  true,
				GeneLabel: "C",
			},
			{
				Name:      "Long Fruit Shape",
				Dominant:  true,
				GeneLabel: "F",
			},
		},
	}

	totalPlants := 64

	// Example of filtering for both dominant and recessive traits
	targetGenotypes := []TargetGenotype{
		{Genotype: "ll", Description: getGenotypeDescription("ll", purpleFlash, candlelight)},
		{Genotype: "cc", Description: getGenotypeDescription("cc", purpleFlash, candlelight)},
	}

	allCombinations := calculateF2Probabilities(purpleFlash, candlelight, totalPlants)
	filteredCombinations, summary := filterCombinations(allCombinations, targetGenotypes)

	fmt.Printf("\nF2 Generation Probabilities for %s × %s\n", purpleFlash.Name, candlelight.Name)
	fmt.Printf("Total plants: %d\n", totalPlants)
	fmt.Println("\nTarget traits:")
	for _, target := range targetGenotypes {
		fmt.Printf("- %s (%s)\n", target.Description, target.Genotype)
	}
	fmt.Println()

	fmt.Println("Matching Combinations:")
	fmt.Println("=====================")
	for _, combo := range filteredCombinations {
		percentage := float64(combo.Probability) / float64(combo.Denominator) * 100
		fmt.Printf("%d/%d (%0.1f%%) = %s\n",
			combo.Probability,
			combo.Denominator,
			percentage,
			combo.Description)
		fmt.Printf("    Genotype: %s\n", combo.GeneNotation)
		fmt.Printf("    Expected number of plants: %.1f\n\n", combo.Expected)
	}

	fmt.Println("\nSummary Statistics:")
	fmt.Println("==================")
	fmt.Printf("Total Probability: %d/%d\n", summary.TotalProbabilityNum, summary.TotalProbabilityDenom)
	fmt.Printf("Percentage: %.1f%%\n", summary.Percentage)
	fmt.Printf("Expected Total Plants with Target Traits: %.1f\n", summary.ExpectedPlants)

	fmt.Println("\nAll Combinations:")
	fmt.Println("===============")
	for _, combo := range allCombinations {
		percentage := float64(combo.Probability) / float64(combo.Denominator) * 100
		fmt.Printf("%d/%d (%0.1f%%) = %s\n",
			combo.Probability,
			combo.Denominator,
			percentage,
			combo.Description)
		fmt.Printf("    Genotype: %s\n", combo.GeneNotation)
		fmt.Printf("    Expected number of plants: %.1f\n\n", combo.Expected)
	}
}
