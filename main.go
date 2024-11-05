package main

import (
	"flag"
	"fmt"
	"math"
	"sort"
	"strings"
)

// At the top of your file, add these constants:
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"

	// Bold versions
	colorBoldGreen  = "\033[1;32m"
	colorBoldYellow = "\033[1;33m"
	colorBoldBlue   = "\033[1;34m"
	colorBoldPurple = "\033[1;35m"
	colorBoldCyan   = "\033[1;36m"
	colorBoldWhite  = "\033[1;37m"
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

func colored(text string, color string) string {
	return color + text + colorReset
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

func getResults(motherPlant Plant, fatherPlant Plant, totalPlants int, targetGenotypes []TargetGenotype) {
	allCombinations := calculateF2Probabilities(motherPlant, fatherPlant, totalPlants)
	filteredCombinations, summary := filterCombinations(allCombinations, targetGenotypes)

	// Title and basic info
	fmt.Printf("\n%s\n", colored(fmt.Sprintf("F2 Generation Probabilities for %s × %s",
		motherPlant.Name, fatherPlant.Name), colorBoldCyan))
	fmt.Printf("%s: %d\n", colored("Total plants", colorBoldWhite), totalPlants)

	// Target traits section if specified
	if len(targetGenotypes) > 0 {
		fmt.Printf("\n%s:\n", colored("Target traits", colorBoldYellow))
		for _, target := range targetGenotypes {
			fmt.Printf("- %s (%s)\n",
				colored(target.Description, colorYellow),
				colored(target.Genotype, colorYellow))
		}
		fmt.Println()

		fmt.Printf("%s:\n", colored("Matching Combinations", colorBoldGreen))
		fmt.Println(colored("=====================", colorGreen))
		for _, combo := range filteredCombinations {
			percentage := float64(combo.Probability) / float64(combo.Denominator) * 100
			fmt.Printf("%s %s = %s\n",
				colored(fmt.Sprintf("%d/%d", combo.Probability, combo.Denominator), colorBoldWhite),
				colored(fmt.Sprintf("(%0.1f%%)", percentage), colorCyan),
				colored(combo.Description, colorGreen))
			fmt.Printf("    Genotype: %s\n", colored(combo.GeneNotation, colorPurple))
			fmt.Printf("    Expected number of plants: %s\n\n",
				colored(fmt.Sprintf("%.1f", combo.Expected), colorYellow))
		}

		fmt.Printf("%s:\n", colored("Target Traits Summary", colorBoldBlue))
		fmt.Println(colored("=====================", colorBlue))
		fmt.Printf("%s: %s\n",
			colored("Total Probability", colorBoldWhite),
			colored(fmt.Sprintf("%d/%d", summary.TotalProbabilityNum, summary.TotalProbabilityDenom), colorWhite))
		fmt.Printf("%s: %s\n",
			colored("Percentage", colorBoldWhite),
			colored(fmt.Sprintf("%.1f%%", summary.Percentage), colorCyan))
		fmt.Printf("%s: %s\n\n",
			colored("Expected Total Plants with Target Traits", colorBoldWhite),
			colored(fmt.Sprintf("%.1f", summary.ExpectedPlants), colorYellow))
	}

	// All possible combinations section
	fmt.Printf("%s:\n", colored("All Possible Combinations", colorBoldPurple))
	fmt.Println(colored("=======================", colorPurple))
	for _, combo := range allCombinations {
		percentage := float64(combo.Probability) / float64(combo.Denominator) * 100

		// Determine if this combination matches target traits
		isTargetMatch := false
		if len(targetGenotypes) > 0 {
			isTargetMatch = true
			for _, target := range targetGenotypes {
				if !isGenotypeMatch(target.Genotype, combo.GeneNotation) {
					isTargetMatch = false
					break
				}
			}
		}

		// Use different colors for matching combinations
		descColor := colorWhite
		if isTargetMatch {
			descColor = colorGreen
		}

		fmt.Printf("%s %s = %s\n",
			colored(fmt.Sprintf("%d/%d", combo.Probability, combo.Denominator), colorBoldWhite),
			colored(fmt.Sprintf("(%0.1f%%)", percentage), colorCyan),
			colored(combo.Description, descColor))
		fmt.Printf("    Genotype: %s\n", colored(combo.GeneNotation, colorPurple))
		fmt.Printf("    Expected number of plants: %s\n",
			colored(fmt.Sprintf("%.1f", combo.Expected), colorYellow))
		if isTargetMatch {
			fmt.Printf("    %s\n", colored("★ Matches target traits", colorBoldYellow))
		}
		fmt.Println()
	}
}
func parseTraits(traitsStr string) []Trait {
	traits := []Trait{}
	traitsArr := strings.Split(traitsStr, ",")
	for _, traitStr := range traitsArr {
		parts := strings.Split(traitStr, ":")
		if len(parts) != 3 {
			continue
		}
		dominant := parts[1] == "true"
		traits = append(traits, Trait{
			Name:      parts[0],
			Dominant:  dominant,
			GeneLabel: parts[2],
		})
	}
	return traits
}

func parseTargetGenotypes(genotypesStr string, motherPlant, fatherPlant Plant) []TargetGenotype {
	if genotypesStr == "" {
		return []TargetGenotype{}
	}

	targetGenotypes := []TargetGenotype{}
	genotypesArr := strings.Split(genotypesStr, ",")

	for _, genotype := range genotypesArr {
		targetGenotypes = append(targetGenotypes, TargetGenotype{
			Genotype:    genotype,
			Description: getGenotypeDescription(genotype, motherPlant, fatherPlant),
		})
	}

	return targetGenotypes
}

func main() {
	motherName := flag.String("motherName", "", "Name of the mother plant")
	motherTraits := flag.String("motherTraits", "", "Traits of the mother plant in the format 'Name:Dominant:GeneLabel,Name:Dominant:GeneLabel,...'")
	fatherName := flag.String("fatherName", "", "Name of the father plant")
	fatherTraits := flag.String("fatherTraits", "", "Traits of the father plant in the format 'Name:Dominant:GeneLabel,Name:Dominant:GeneLabel,...'")
	totalPlants := flag.Int("totalPlants", 64, "Total number of plants")
	targetGenotypesStr := flag.String("targetGenotypes", "", "Comma-separated list of target genotypes (e.g., 'll,cc' for mutant leaves and purple foliage)")

	flag.Parse()

	if *motherName == "" || *motherTraits == "" || *fatherName == "" || *fatherTraits == "" {
		fmt.Println("All plant names and traits must be provided")
		return
	}

	motherPlant := Plant{
		Name:   *motherName,
		Traits: parseTraits(*motherTraits),
	}

	fatherPlant := Plant{
		Name:   *fatherName,
		Traits: parseTraits(*fatherTraits),
	}

	targetGenotypes := parseTargetGenotypes(*targetGenotypesStr, motherPlant, fatherPlant)

	getResults(motherPlant, fatherPlant, *totalPlants, targetGenotypes)
}
