package genex

import (
	"math"
	"regexp/syntax"
)

// Generates all the strings that match the `input` regex after whitelisting `charset`.
// The `infinite` argument caps the maximum boundary of repetition operators.
func Generate(input, charset *syntax.Regexp, infinite int, callback func(string)) {
	var generate func(input, charset *syntax.Regexp, infinite int) _Iterator

	generate = func(input, charset *syntax.Regexp, infinite int) _Iterator {
		result := []_Iterator{}

		switch input.Op {
			case syntax.OpStar, syntax.OpPlus, syntax.OpQuest, syntax.OpRepeat:
				value := []_Iterator{}

				for _, sub := range input.Sub {
					value = append(value, generate(sub, charset, infinite))
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

				result = append(result, _NewRepeat(_NewStack(value), input.Min, input.Max))

			case syntax.OpCharClass, syntax.OpAnyCharNotNL, syntax.OpAnyChar:
				if input.Op != syntax.OpCharClass {
					input = charset
				}

				data := []string{}

				for i := 0; i < len(input.Rune); i += 2 {
					for j := 0; j < len(charset.Rune); j += 2 {
						bounds := []rune{
							rune(math.Max(float64(input.Rune[i]), float64(charset.Rune[j]))),
							rune(math.Min(float64(input.Rune[i + 1]), float64(charset.Rune[j + 1]))),
						}

						if bounds[0] <= bounds[1] {
							for char := bounds[0]; char <= bounds[1]; char++ {
								data = append(data, string(char))
							}
						}
					}
				}

				result = append(result, _NewSet(data))

			case syntax.OpCapture, syntax.OpConcat:
				for _, sub := range input.Sub {
					result = append(result, generate(sub, charset, infinite))
				}

			case syntax.OpAlternate:
				options := []_Iterator{}

				for _, sub := range input.Sub {
					options = append(options, generate(sub, charset, infinite))
				}

				result = append(result, _NewOption(options))

			case syntax.OpLiteral:
				result = append(result, _NewSet([]string{string(input.Rune)}))

			default:
				result = append(result, _NewSet([]string{""}))
		}

		return _NewStack(result)
	}

	if charset.Op != syntax.OpCharClass {
		charset, _ = syntax.Parse(`[[:print:]]`, syntax.Perl)
	}

	iterator := generate(input, charset, infinite)

	for iterator.rewind(); iterator.valid(); iterator.next() {
		callback(iterator.current())
	}
}

type _Iterator interface {
	rewind()
	next()
	valid() bool
	current() string
	clone() _Iterator
}

type _Set struct {
	i int
	data []string
}

func (this *_Set) rewind() {
	this.i = 0
}

func (this *_Set) valid() bool {
	return this.i < len(this.data)
}

func (this *_Set) current() string {
	return string(this.data[this.i])
}

func (this *_Set) next() {
	this.i++
}

func (this *_Set) clone() _Iterator {
	clone := &_Set{data: make([]string, len(this.data))}

	for key, value := range this.data {
		clone.data[key] = value
	}

	return clone
}

func _NewSet(data []string) *_Set {
	return &_Set{data: data}
}

type _Stack struct {
	data []_Iterator
}

func (this *_Stack) rewind() {
	for i := range this.data {
		this.data[i].rewind()
	}
}

func (this *_Stack) valid() bool {
	return this.data[0].valid();
}

func (this *_Stack) current() string {
	result := ""

	for i := range this.data {
		result += this.data[i].current()
	}

	return result
}

func (this *_Stack) next() {
	if this.valid() {
		i := len(this.data)
		i--; this.data[i].next()

		for i > 0 && !this.data[i].valid() {
			this.data[i].rewind()
			i--; this.data[i].next()
		}
	}
}

func (this *_Stack) clone() _Iterator {
	clone := &_Stack{data: make([]_Iterator, len(this.data))}

	for key, value := range this.data {
		clone.data[key] = value.clone()
	}

	return clone
}

func _NewStack(data []_Iterator) *_Stack {
	return &_Stack{data: data}
}

type _Option struct {
	i int
	data []_Iterator
}

func (this *_Option) rewind() {
	this.i = 0

	for i := range this.data {
		this.data[i].rewind()
	}
}

func (this *_Option) valid() bool {
	return this.i < len(this.data)
}

func (this *_Option) current() string {
	return this.data[this.i].current()
}

func (this *_Option) next() {
	if this.valid() {
		this.data[this.i].next()

		for this.valid() && !this.data[this.i].valid() {
			this.i++
		}
	}
}

func (this *_Option) clone() _Iterator {
	clone := &_Option{data: make([]_Iterator, len(this.data))}

	for key, value := range this.data {
		clone.data[key] = value.clone()
	}

	return clone
}

func _NewOption(data []_Iterator) *_Option {
	return &_Option{data: data}
}

func _NewRepeat(data _Iterator, min int, max int) *_Stack {
	stack := []_Iterator{}

	for i := 0; i < min; i++ {
		stack = append(stack, data.clone())
	}

	if max > min {
		stack = append(stack, _NewOption([]_Iterator{
			_NewSet([]string{""}), _NewRepeat(data, 1, max - min),
		}))
	}

	return _NewStack(stack)
}
