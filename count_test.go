package genex

import (
	"regexp/syntax"
	"testing"
)

func TestCountCharacters(t *testing.T) {
	charset, _ := syntax.Parse(`[0-9a-zA-Z-._]`, syntax.Perl)
	expected := []struct{
		in	string
		out	float64
	}{
		{``,						1},
		{`a`,						1},
		{`abc`,						1},
	}

	for _, pair := range expected {
		in, _ := syntax.Parse(pair.in, syntax.Perl)
		value := Count(in, charset, 3)

		if value != pair.out {
			t.Error("For", pair.in, "expected", pair.out, "but got", value)
		}
	}
}

func TestCountCharacterSets(t *testing.T) {
	charset, _ := syntax.Parse(`[0-9a-zA-Z-._]`, syntax.Perl)
	expected := []struct{
		in	string
		out	float64
	}{
		{`.`,						65},
		{`[0-4]`,					5},
		{`[0-9a-zA-Z]`,				62},
		{`[\d]`,					10},
		{`[\D]`,					55},
		{`[\s]`,					1},
		{`[\S]`,					65},
		{`[\W]`,					2},
		{`[\w]`,					63},
		{`[^\D\d]`,					1},
		{`[^\D]`,					10},
		{`[^\d]`,					55},
		{`[^\S\s]`,					1},
		{`[^\S]`,					1},
		{`[^\s]`,					65},
		{`[^\W\w]`,					1},
		{`[^\w]`,					2},
		{`[^\W]`,					63},
		{`[^AN]BC`,					63},
		{`[a-z]`,					26},
		{`[abc]`,					3},
		{`\d`,						10},
		{`\D`,						55},
		{`\s`,						1},
		{`\S`,						65},
		{`\W`,						2},
		{`\w`,						63},
	}

	for _, pair := range expected {
		in, _ := syntax.Parse(pair.in, syntax.Perl)
		value := Count(in, charset, 3)

		if value != pair.out {
			t.Error("For", pair.in, "expected", pair.out, "but got", value)
		}
	}
}

func TestCountRepetition(t *testing.T) {
	charset, _ := syntax.Parse(`[0-9a-zA-Z-._]`, syntax.Perl)
	expected := []struct{
		in	string
		out	float64
	}{
		{`[ab]{0,2}?`,				7},
		{`[ab]{0,2}`,				7},
		{`[ab]{1,2}?`,				6},
		{`[ab]{1,2}`,				6},
		{`\d{0,4}`,					11111},
		{`\d{1,4}`,					11110},
		{`\d{2,4}`,					11100},
		{`\d{3,4}`,					11000},
		{`\d{4,4}`,					10000},
		{`\d{5}`,					100000},
		{`\s{0,}`,					4},
		{`\s{2,}`,					4},
		{`a\b{0,}`,					4},
		{`ab*`,						4},
		{`ab+`,						4},
		{`ab?`,						2},
	}

	for _, pair := range expected {
		in, _ := syntax.Parse(pair.in, syntax.Perl)
		value := Count(in, charset, 3)

		if value != pair.out {
			t.Error("For", pair.in, "expected", pair.out, "but got", value)
		}
	}
}

func TestCountRepetitionWithHigherInfinity(t *testing.T) {
	charset, _ := syntax.Parse(`[0-9a-zA-Z-._]`, syntax.Perl)
	expected := []struct{
		in	string
		out	float64
	}{
		{`[ab]{0,2}?`,				7},
		{`[ab]{0,2}`,				7},
		{`[ab]{1,2}?`,				6},
		{`[ab]{1,2}`,				6},
		{`\d{0,4}`,					11111},
		{`\d{1,4}`,					11110},
		{`\d{2,4}`,					11100},
		{`\d{3,4}`,					11000},
		{`\d{4,4}`,					10000},
		{`\d{5}`,					100000},
		{`\s{0,}`,					513},
		{`\s{2,}`,					513},
		{`a\b{0,}`,					513},
		{`ab*`,						513},
		{`ab+`,						513},
		{`ab?`,						2},
	}

	for _, pair := range expected {
		in, _ := syntax.Parse(pair.in, syntax.Perl)
		value := Count(in, charset, 512)

		if value != pair.out {
			t.Error("For", pair.in, "expected", pair.out, "but got", value)
		}
	}
}

func TestCountAlternationAndGrouping(t *testing.T) {
	charset, _ := syntax.Parse(`[0-9a-zA-Z-._]`, syntax.Perl)
	expected := []struct{
		in	string
		out	float64
	}{
		{`(a)`,						1},
		{`[0101]?[0-3]?`,			15},
		{`[01]?[01]?`,				9},
		{`[01]?[01]`,				6},
		{`forever|(old|young)?`,	4},
		{`forever|(old|young)`,		3},
		{`forever|young`,			2},
		{`|\b`,						2},
		{`|`,						1},
		{`|a`,						2},
	}

	for _, pair := range expected {
		in, _ := syntax.Parse(pair.in, syntax.Perl)
		value := Count(in, charset, 3)

		if value != pair.out {
			t.Error("For", pair.in, "expected", pair.out, "but got", value)
		}
	}
}
