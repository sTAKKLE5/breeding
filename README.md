# Plant Breeding Calculator

A Go program that calculates trait inheritance probabilities in breeding projects. This tool helps breeders predict the distribution of traits across generations and plan their breeding programs effectively.

## Command Line Flags

- **-motherName**: Name of the mother plant
- **-motherTraits**: Traits of the mother plant in the format `Name:Dominant:GeneLabel,Name:Dominant:GeneLabel,...`
- **-fatherName**: Name of the father plant
- **-fatherTraits**: Traits of the father plant in the same format as mother traits
- **-totalPlants**: Total number of plants (default: 64)
- **-targetGenotypes**: (Optional) List of target genotypes to filter for (e.g., "ll,cc" for mutant leaves and purple foliage)

## Core Features

- Calculates F2 generation probabilities for multiple traits
- Handles both dominant and recessive inheritance patterns
- Provides probability ratios and expected plant numbers
- Filters combinations based on target genotypes
- Generates human-readable trait descriptions

## Genetic Notation

The calculator supports two types of genotype notation for filtering:

### Recessive Traits
Use doubled lowercase letters:
- `ll` for mutant traits
- `cc` for recessive color traits
- `ff` for recessive form traits

### Dominant Traits
Use single uppercase letters:
- `L` for regular traits
- `C` for dominant color traits
- `F` for dominant form traits

### Results Notation
- `L_` indicates presence of dominant allele (LL or Ll)
- `ll` indicates homozygous recessive
- Space separated for multiple traits (e.g., `L_ cc F_`)

## Sample Command Line Execution

```bash
go run main.go \
  -motherName="Purple Flash" \
  -motherTraits="Regular Leave Shape:true:L,Purple Foliage:false:C,Round Fruit Shape:false:F" \
  -fatherName="Candlelight Mutant" \
  -fatherTraits="Mutant Leave Shape:false:L,Green Foliage:true:C,Long Fruit Shape:true:F" \
  -totalPlants=64 \
  -targetGenotypes="ll,cc"  # Filter for mutant leaves and purple foliage
```

### Additional Examples

Filter for dominant traits:
```bash
-targetGenotypes="L,F"  # Filter for regular leaves and long fruit
```

Mix of dominant and recessive:
```bash
-targetGenotypes="L,cc"  # Filter for regular leaves and purple foliage
```

## Example Output

```
F2 Generation Probabilities for Purple Flash × Candlelight Mutant
Total plants: 64

Target traits:
- Regular Leave Shape (L)
- Purple Foliage (cc)

Matching Combinations:
=====================
9/64 (14.1%) = Regular Leave Shape, Purple Foliage, Long Fruit Shape
    Genotype: L_ cc F_
    Expected number of plants: 9.0

Summary Statistics:
==================
Total Probability: 12/64
Percentage: 18.8%
Expected Total Plants with Target Traits: 12.0
```

## Program Structure

### Key Types
```go
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
```

## Limitations

- Assumes simple Mendelian inheritance
- Does not account for:
  - Incomplete dominance
  - Codominance
  - Complex gene interactions
  - Environmental effects

## Future Improvements

- Support for incomplete dominance
- Visualization of results
- F3 generation predictions
- Confidence intervals
- Complex filtering patterns
- Visual trait distributions
- Batch processing capabilities

## License

MIT License