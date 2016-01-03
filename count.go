package genex

import (
	"math"
	"regexp/syntax"
)

// Count computes the total number of matches the `input` regex would generate after whitelisting `charset`.
// The `infinite` argument caps the maximum boundary of repetition operators.
func Count(input, charset *syntax.Regexp, infinite int) float64 {
	var count func(input, charset *syntax.Regexp, infinite int) float64

	count = func(input, charset *syntax.Regexp, infinite int) float64 {
		result := float64(0)

		switch input.Op {
			case syntax.OpStar, syntax.OpPlus, syntax.OpQuest, syntax.OpRepeat:
				value := float64(1)

				for _, sub := range input.Sub {
					value *= count(sub, charset, infinite)
				}

				switch input.Op {
					case syntax.OpStar:
						input.Min = 0
						input.Max = -1

					case syntax.OpPlus:
						input.Min = 1
						input.Max = -1

					case syntax.OpQuest:
						input.Min = 0
						input.Max = 1
				}

				if input.Max == -1 && infinite >= 0 {
					input.Max = input.Min + infinite
				}

				if input.Max == -1 {
					result = math.Inf(1)
				} else if value > 1 {
					if input.Min == input.Max {
						result = math.Pow(value, float64(input.Min))
					} else {
						result = (math.Pow(value, float64(input.Max) + 1) - 1) / (value - 1)

						if input.Min > 0 {
							result -= (math.Pow(value, float64(input.Min) + 0) - 1) / (value - 1)
						}
					}
				} else {
					result = float64(input.Max - input.Min) + 1
				}

			case syntax.OpCharClass, syntax.OpAnyCharNotNL, syntax.OpAnyChar:
				if input.Op != syntax.OpCharClass {
					input = charset
				}

				for i := 0; i < len(input.Rune); i += 2 {
					for j := 0; j < len(charset.Rune); j += 2 {
						bounds := []float64{
							math.Max(float64(input.Rune[i]), float64(charset.Rune[j])),
							math.Min(float64(input.Rune[i + 1]), float64(charset.Rune[j + 1])),
						}

						if bounds[0] <= bounds[1] {
							result += bounds[1] - bounds[0] + 1
						}
					}
				}

			case syntax.OpCapture, syntax.OpConcat:
				result = 1

				for _, sub := range input.Sub {
					result *= count(sub, charset, infinite)
				}

			case syntax.OpAlternate:
				for _, sub := range input.Sub {
					result += count(sub, charset, infinite)
				}

			default:
				result = 1
		}

		if math.IsNaN(result) {
			result = math.Inf(1)
		}

		return math.Max(1, result)
	}

	if charset.Op != syntax.OpCharClass {
		charset, _ = syntax.Parse(`[[:print:]]`, syntax.Perl)
	}

	return count(input, charset, infinite)
}
