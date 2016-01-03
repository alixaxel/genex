package genex

import (
	"reflect"
	"regexp/syntax"
	"testing"
)

func TestGenerateCharacters(t *testing.T) {
	charset, _ := syntax.Parse(`[0-9a-zA-Z-._]`, syntax.Perl)
	expected := []struct{
		in	string
		out	[]string
	}{
		{``,                        []string{``}},
		{`a`,                       []string{`a`}},
		{`abc`,                     []string{`abc`}},
	}

	for _, pair := range expected {
		in, _ := syntax.Parse(pair.in, syntax.Perl)
		values := []string{}

		Generate(in, charset, 3, func(output string) {
			values = append(values, output)
		})

		if reflect.DeepEqual(values, pair.out) != true {
			t.Error("For", pair.in, "expected", pair.out, "but got", values)
		}
	}
}

func TestGenerateCharacterSets(t *testing.T) {
	charset, _ := syntax.Parse(`[0-9a-zA-Z-._]`, syntax.Perl)
	expected := []struct{
		in	string
		out	[]string
	}{
		{`.`,                       []string{`-`, `.`, `0`, `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `A`, `B`, `C`, `D`, `E`, `F`, `G`, `H`, `I`, `J`, `K`, `L`, `M`, `N`, `O`, `P`, `Q`, `R`, `S`, `T`, `U`, `V`, `W`, `X`, `Y`, `Z`, `_`, `a`, `b`, `c`, `d`, `e`, `f`, `g`, `h`, `i`, `j`, `k`, `l`, `m`, `n`, `o`, `p`, `q`, `r`, `s`, `t`, `u`, `v`, `w`, `x`, `y`, `z`}},
		{`[0-4]`,                   []string{`0`, `1`, `2`, `3`, `4`}},
		{`[\s]`,                    []string{}},
		{`[\W]`,                    []string{`-`, `.`}},
		{`[^\D\d]`,                 []string{}},
		{`[^\S\s]`,                 []string{}},
		{`[^\S]`,                   []string{}},
		{`[^\W\w]`,                 []string{}},
		{`[^\w]`,                   []string{`-`, `.`}},
		{`[abc]`,                   []string{`a`, `b`, `c`}},
		{`\W`,                      []string{`-`, `.`}},
	}

	for _, pair := range expected {
		in, _ := syntax.Parse(pair.in, syntax.Perl)
		values := []string{}

		Generate(in, charset, 3, func(output string) {
			values = append(values, output)
		})

		if reflect.DeepEqual(values, pair.out) != true {
			t.Error("For", pair.in, "expected", pair.out, "but got", values)
		}
	}
}

func TestGenerateRepetition(t *testing.T) {
	charset, _ := syntax.Parse(`[0-9a-zA-Z-._]`, syntax.Perl)
	expected := []struct{
		in	string
		out	[]string
	}{
		{`[ab]{0,2}?`,              []string{``, `a`, `aa`, `ab`, `b`, `ba`, `bb`}},
		{`[ab]{0,2}`,               []string{``, `a`, `aa`, `ab`, `b`, `ba`, `bb`}},
		{`[ab]{1,2}?`,              []string{`a`, `aa`, `ab`, `b`, `ba`, `bb`}},
		{`[ab]{1,2}`,               []string{`a`, `aa`, `ab`, `b`, `ba`, `bb`}},
		{`\s{0,}`,                  []string{``}},
		{`\s{2,}`,                  []string{}},
		{`a\b{0,}`,                 []string{`a`, `a`, `a`, `a`}},
		{`ab*`,                     []string{`a`, `ab`, `abb`, `abbb`}},
		{`ab+`,                     []string{`ab`, `abb`, `abbb`, `abbbb`}},
		{`ab?`,                     []string{`a`, `ab`}},
	}

	for _, pair := range expected {
		in, _ := syntax.Parse(pair.in, syntax.Perl)
		values := []string{}

		Generate(in, charset, 3, func(output string) {
			values = append(values, output)
		})

		if reflect.DeepEqual(values, pair.out) != true {
			t.Error("For", pair.in, "expected", pair.out, "but got", values)
		}
	}
}

func TestGenerateAlternationAndGrouping(t *testing.T) {
	charset, _ := syntax.Parse(`[0-9a-zA-Z-._]`, syntax.Perl)
	expected := []struct{
		in	string
		out	[]string
	}{
		{`(a)`,                     []string{`a`}},
		{`[0101]?[0-3]?`,           []string{``, `0`, `1`, `2`, `3`, `0`, `00`, `01`, `02`, `03`, `1`, `10`, `11`, `12`, `13`}},
		{`[01]?[01]?`,              []string{``, `0`, `1`, `0`, `00`, `01`, `1`, `10`, `11`}},
		{`[01]?[01]`,               []string{`0`, `1`, `00`, `01`, `10`, `11`}},
		{`forever|(old|young)?`,    []string{`forever`, ``, `old`, `young`}},
		{`forever|(old|young)`,     []string{`forever`, `old`, `young`}},
		{`forever|young`,           []string{`forever`, `young`}},
		{`|`,                       []string{``}},
		{`|a`,                      []string{``, `a`}},
	}

	for _, pair := range expected {
		in, _ := syntax.Parse(pair.in, syntax.Perl)
		values := []string{}

		Generate(in, charset, 3, func(output string) {
			values = append(values, output)
		})

		if reflect.DeepEqual(values, pair.out) != true {
			t.Error("For", pair.in, "expected", pair.out, "but got", values)
		}
	}
}

func TestInvalidCharset(t *testing.T) {
	charset, _ := syntax.Parse(`^$`, syntax.Perl)
	expected := []struct{
		in	string
		out	float64
	}{
		{``,                        1},
	}

	for _, pair := range expected {
		in, _ := syntax.Parse(pair.in, syntax.Perl)
		value := Count(in, charset, 3)

		if value != pair.out {
			t.Error("For", pair.in, "expected", pair.out, "but got", value)
		}
	}
}
