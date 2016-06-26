package elmo

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const end_symbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleScript
	ruleLine
	ruleEndOfLine
	ruleArgument
	ruleFunctionCall
	ruleBlock
	ruleSpacing
	ruleWhiteSpace
	ruleLongComment
	ruleLineComment
	ruleNewLine
	ruleIdentifier
	ruleIdNondigit
	ruleIdChar
	ruleStringLiteral
	ruleStringChar
	ruleEscape
	ruleDecimalConstant
	ruleLPAR
	ruleRPAR
	ruleLCURLY
	ruleRCURLY
	rulePCOMMA
	ruleEOT

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"Script",
	"Line",
	"EndOfLine",
	"Argument",
	"FunctionCall",
	"Block",
	"Spacing",
	"WhiteSpace",
	"LongComment",
	"LineComment",
	"NewLine",
	"Identifier",
	"IdNondigit",
	"IdChar",
	"StringLiteral",
	"StringChar",
	"Escape",
	"DecimalConstant",
	"LPAR",
	"RPAR",
	"LCURLY",
	"RCURLY",
	"PCOMMA",
	"EOT",

	"Pre_",
	"_In_",
	"_Suf",
}

type tokenTree interface {
	Print()
	PrintSyntax()
	PrintSyntaxTree(buffer string)
	Add(rule pegRule, begin, end, next uint32, depth int)
	Expand(index int) tokenTree
	Tokens() <-chan token32
	AST() *node32
	Error() []token32
	trim(length int)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(string(([]rune(buffer)[node.begin:node.end]))))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (ast *node32) Print(buffer string) {
	ast.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next uint32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: uint32(t.begin), end: uint32(t.end), next: uint32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = uint32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i, _ := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, uint32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre_, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: rule_Suf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(string(([]rune(buffer)[token.begin:token.end]))))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth uint32, index int) {
	t.tree[index] = token32{pegRule: rule, begin: uint32(begin), end: uint32(end), next: uint32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i, _ := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

/*func (t *tokens16) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2 * len(tree))
		for i, v := range tree {
			expanded[i] = v.getToken32()
		}
		return &tokens32{tree: expanded}
	}
	return nil
}*/

func (t *tokens32) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	return nil
}

type ElmoGrammar struct {
	Buffer string
	buffer []rune
	rules  [25]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	Pretty bool
	tokenTree
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *ElmoGrammar
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *ElmoGrammar) PrintSyntaxTree() {
	p.tokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *ElmoGrammar) Highlighter() {
	p.tokenTree.PrintSyntax()
}

func (p *ElmoGrammar) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != end_symbol {
		p.buffer = append(p.buffer, end_symbol)
	}

	var tree tokenTree = &tokens32{tree: make([]token32, math.MaxInt16)}
	var max token32
	position, depth, tokenIndex, buffer, _rules := uint32(0), uint32(0), 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokenTree = tree
		if matches {
			p.tokenTree.trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin uint32) {
		if t := tree.Expand(tokenIndex); t != nil {
			tree = t
		}
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position, depth}
		}
	}

	matchDot := func() bool {
		if buffer[position] != end_symbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Script <- <(Spacing Line* EOT)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				if !_rules[ruleSpacing]() {
					goto l0
				}
			l2:
				{
					position3, tokenIndex3, depth3 := position, tokenIndex, depth
					if !_rules[ruleLine]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex, depth = position3, tokenIndex3, depth3
				}
				if !_rules[ruleEOT]() {
					goto l0
				}
				depth--
				add(ruleScript, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 Line <- <(NewLine? Identifier Argument* EndOfLine?)> */
		func() bool {
			position4, tokenIndex4, depth4 := position, tokenIndex, depth
			{
				position5 := position
				depth++
				{
					position6, tokenIndex6, depth6 := position, tokenIndex, depth
					if !_rules[ruleNewLine]() {
						goto l6
					}
					goto l7
				l6:
					position, tokenIndex, depth = position6, tokenIndex6, depth6
				}
			l7:
				if !_rules[ruleIdentifier]() {
					goto l4
				}
			l8:
				{
					position9, tokenIndex9, depth9 := position, tokenIndex, depth
					if !_rules[ruleArgument]() {
						goto l9
					}
					goto l8
				l9:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
				}
				{
					position10, tokenIndex10, depth10 := position, tokenIndex, depth
					if !_rules[ruleEndOfLine]() {
						goto l10
					}
					goto l11
				l10:
					position, tokenIndex, depth = position10, tokenIndex10, depth10
				}
			l11:
				depth--
				add(ruleLine, position5)
			}
			return true
		l4:
			position, tokenIndex, depth = position4, tokenIndex4, depth4
			return false
		},
		/* 2 EndOfLine <- <(PCOMMA / NewLine)> */
		func() bool {
			position12, tokenIndex12, depth12 := position, tokenIndex, depth
			{
				position13 := position
				depth++
				{
					position14, tokenIndex14, depth14 := position, tokenIndex, depth
					if !_rules[rulePCOMMA]() {
						goto l15
					}
					goto l14
				l15:
					position, tokenIndex, depth = position14, tokenIndex14, depth14
					if !_rules[ruleNewLine]() {
						goto l12
					}
				}
			l14:
				depth--
				add(ruleEndOfLine, position13)
			}
			return true
		l12:
			position, tokenIndex, depth = position12, tokenIndex12, depth12
			return false
		},
		/* 3 Argument <- <(Identifier / StringLiteral / DecimalConstant / FunctionCall / Block)> */
		func() bool {
			position16, tokenIndex16, depth16 := position, tokenIndex, depth
			{
				position17 := position
				depth++
				{
					position18, tokenIndex18, depth18 := position, tokenIndex, depth
					if !_rules[ruleIdentifier]() {
						goto l19
					}
					goto l18
				l19:
					position, tokenIndex, depth = position18, tokenIndex18, depth18
					if !_rules[ruleStringLiteral]() {
						goto l20
					}
					goto l18
				l20:
					position, tokenIndex, depth = position18, tokenIndex18, depth18
					if !_rules[ruleDecimalConstant]() {
						goto l21
					}
					goto l18
				l21:
					position, tokenIndex, depth = position18, tokenIndex18, depth18
					if !_rules[ruleFunctionCall]() {
						goto l22
					}
					goto l18
				l22:
					position, tokenIndex, depth = position18, tokenIndex18, depth18
					if !_rules[ruleBlock]() {
						goto l16
					}
				}
			l18:
				depth--
				add(ruleArgument, position17)
			}
			return true
		l16:
			position, tokenIndex, depth = position16, tokenIndex16, depth16
			return false
		},
		/* 4 FunctionCall <- <(LPAR Line RPAR)> */
		func() bool {
			position23, tokenIndex23, depth23 := position, tokenIndex, depth
			{
				position24 := position
				depth++
				if !_rules[ruleLPAR]() {
					goto l23
				}
				if !_rules[ruleLine]() {
					goto l23
				}
				if !_rules[ruleRPAR]() {
					goto l23
				}
				depth--
				add(ruleFunctionCall, position24)
			}
			return true
		l23:
			position, tokenIndex, depth = position23, tokenIndex23, depth23
			return false
		},
		/* 5 Block <- <(LCURLY NewLine* Line* RCURLY)> */
		func() bool {
			position25, tokenIndex25, depth25 := position, tokenIndex, depth
			{
				position26 := position
				depth++
				if !_rules[ruleLCURLY]() {
					goto l25
				}
			l27:
				{
					position28, tokenIndex28, depth28 := position, tokenIndex, depth
					if !_rules[ruleNewLine]() {
						goto l28
					}
					goto l27
				l28:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
				}
			l29:
				{
					position30, tokenIndex30, depth30 := position, tokenIndex, depth
					if !_rules[ruleLine]() {
						goto l30
					}
					goto l29
				l30:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
				}
				if !_rules[ruleRCURLY]() {
					goto l25
				}
				depth--
				add(ruleBlock, position26)
			}
			return true
		l25:
			position, tokenIndex, depth = position25, tokenIndex25, depth25
			return false
		},
		/* 6 Spacing <- <(WhiteSpace / LongComment / LineComment)*> */
		func() bool {
			{
				position32 := position
				depth++
			l33:
				{
					position34, tokenIndex34, depth34 := position, tokenIndex, depth
					{
						position35, tokenIndex35, depth35 := position, tokenIndex, depth
						if !_rules[ruleWhiteSpace]() {
							goto l36
						}
						goto l35
					l36:
						position, tokenIndex, depth = position35, tokenIndex35, depth35
						if !_rules[ruleLongComment]() {
							goto l37
						}
						goto l35
					l37:
						position, tokenIndex, depth = position35, tokenIndex35, depth35
						if !_rules[ruleLineComment]() {
							goto l34
						}
					}
				l35:
					goto l33
				l34:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
				}
				depth--
				add(ruleSpacing, position32)
			}
			return true
		},
		/* 7 WhiteSpace <- <(' ' / '\t')> */
		func() bool {
			position38, tokenIndex38, depth38 := position, tokenIndex, depth
			{
				position39 := position
				depth++
				{
					position40, tokenIndex40, depth40 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l41
					}
					position++
					goto l40
				l41:
					position, tokenIndex, depth = position40, tokenIndex40, depth40
					if buffer[position] != rune('\t') {
						goto l38
					}
					position++
				}
			l40:
				depth--
				add(ruleWhiteSpace, position39)
			}
			return true
		l38:
			position, tokenIndex, depth = position38, tokenIndex38, depth38
			return false
		},
		/* 8 LongComment <- <('/' '*' (!('*' '/') .)* ('*' '/'))> */
		func() bool {
			position42, tokenIndex42, depth42 := position, tokenIndex, depth
			{
				position43 := position
				depth++
				if buffer[position] != rune('/') {
					goto l42
				}
				position++
				if buffer[position] != rune('*') {
					goto l42
				}
				position++
			l44:
				{
					position45, tokenIndex45, depth45 := position, tokenIndex, depth
					{
						position46, tokenIndex46, depth46 := position, tokenIndex, depth
						if buffer[position] != rune('*') {
							goto l46
						}
						position++
						if buffer[position] != rune('/') {
							goto l46
						}
						position++
						goto l45
					l46:
						position, tokenIndex, depth = position46, tokenIndex46, depth46
					}
					if !matchDot() {
						goto l45
					}
					goto l44
				l45:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
				}
				if buffer[position] != rune('*') {
					goto l42
				}
				position++
				if buffer[position] != rune('/') {
					goto l42
				}
				position++
				depth--
				add(ruleLongComment, position43)
			}
			return true
		l42:
			position, tokenIndex, depth = position42, tokenIndex42, depth42
			return false
		},
		/* 9 LineComment <- <('#' (!'\n' .)*)> */
		func() bool {
			position47, tokenIndex47, depth47 := position, tokenIndex, depth
			{
				position48 := position
				depth++
				if buffer[position] != rune('#') {
					goto l47
				}
				position++
			l49:
				{
					position50, tokenIndex50, depth50 := position, tokenIndex, depth
					{
						position51, tokenIndex51, depth51 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l51
						}
						position++
						goto l50
					l51:
						position, tokenIndex, depth = position51, tokenIndex51, depth51
					}
					if !matchDot() {
						goto l50
					}
					goto l49
				l50:
					position, tokenIndex, depth = position50, tokenIndex50, depth50
				}
				depth--
				add(ruleLineComment, position48)
			}
			return true
		l47:
			position, tokenIndex, depth = position47, tokenIndex47, depth47
			return false
		},
		/* 10 NewLine <- <(('\n' / '\r') Spacing)+> */
		func() bool {
			position52, tokenIndex52, depth52 := position, tokenIndex, depth
			{
				position53 := position
				depth++
				{
					position56, tokenIndex56, depth56 := position, tokenIndex, depth
					if buffer[position] != rune('\n') {
						goto l57
					}
					position++
					goto l56
				l57:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					if buffer[position] != rune('\r') {
						goto l52
					}
					position++
				}
			l56:
				if !_rules[ruleSpacing]() {
					goto l52
				}
			l54:
				{
					position55, tokenIndex55, depth55 := position, tokenIndex, depth
					{
						position58, tokenIndex58, depth58 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l59
						}
						position++
						goto l58
					l59:
						position, tokenIndex, depth = position58, tokenIndex58, depth58
						if buffer[position] != rune('\r') {
							goto l55
						}
						position++
					}
				l58:
					if !_rules[ruleSpacing]() {
						goto l55
					}
					goto l54
				l55:
					position, tokenIndex, depth = position55, tokenIndex55, depth55
				}
				depth--
				add(ruleNewLine, position53)
			}
			return true
		l52:
			position, tokenIndex, depth = position52, tokenIndex52, depth52
			return false
		},
		/* 11 Identifier <- <(IdNondigit IdChar* Spacing)> */
		func() bool {
			position60, tokenIndex60, depth60 := position, tokenIndex, depth
			{
				position61 := position
				depth++
				if !_rules[ruleIdNondigit]() {
					goto l60
				}
			l62:
				{
					position63, tokenIndex63, depth63 := position, tokenIndex, depth
					if !_rules[ruleIdChar]() {
						goto l63
					}
					goto l62
				l63:
					position, tokenIndex, depth = position63, tokenIndex63, depth63
				}
				if !_rules[ruleSpacing]() {
					goto l60
				}
				depth--
				add(ruleIdentifier, position61)
			}
			return true
		l60:
			position, tokenIndex, depth = position60, tokenIndex60, depth60
			return false
		},
		/* 12 IdNondigit <- <([a-z] / [A-Z] / '_')> */
		func() bool {
			position64, tokenIndex64, depth64 := position, tokenIndex, depth
			{
				position65 := position
				depth++
				{
					position66, tokenIndex66, depth66 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l67
					}
					position++
					goto l66
				l67:
					position, tokenIndex, depth = position66, tokenIndex66, depth66
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l68
					}
					position++
					goto l66
				l68:
					position, tokenIndex, depth = position66, tokenIndex66, depth66
					if buffer[position] != rune('_') {
						goto l64
					}
					position++
				}
			l66:
				depth--
				add(ruleIdNondigit, position65)
			}
			return true
		l64:
			position, tokenIndex, depth = position64, tokenIndex64, depth64
			return false
		},
		/* 13 IdChar <- <([a-z] / [A-Z] / [0-9] / '_')> */
		func() bool {
			position69, tokenIndex69, depth69 := position, tokenIndex, depth
			{
				position70 := position
				depth++
				{
					position71, tokenIndex71, depth71 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l72
					}
					position++
					goto l71
				l72:
					position, tokenIndex, depth = position71, tokenIndex71, depth71
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l73
					}
					position++
					goto l71
				l73:
					position, tokenIndex, depth = position71, tokenIndex71, depth71
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l74
					}
					position++
					goto l71
				l74:
					position, tokenIndex, depth = position71, tokenIndex71, depth71
					if buffer[position] != rune('_') {
						goto l69
					}
					position++
				}
			l71:
				depth--
				add(ruleIdChar, position70)
			}
			return true
		l69:
			position, tokenIndex, depth = position69, tokenIndex69, depth69
			return false
		},
		/* 14 StringLiteral <- <('"' StringChar* '"' Spacing)+> */
		func() bool {
			position75, tokenIndex75, depth75 := position, tokenIndex, depth
			{
				position76 := position
				depth++
				if buffer[position] != rune('"') {
					goto l75
				}
				position++
			l79:
				{
					position80, tokenIndex80, depth80 := position, tokenIndex, depth
					if !_rules[ruleStringChar]() {
						goto l80
					}
					goto l79
				l80:
					position, tokenIndex, depth = position80, tokenIndex80, depth80
				}
				if buffer[position] != rune('"') {
					goto l75
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l75
				}
			l77:
				{
					position78, tokenIndex78, depth78 := position, tokenIndex, depth
					if buffer[position] != rune('"') {
						goto l78
					}
					position++
				l81:
					{
						position82, tokenIndex82, depth82 := position, tokenIndex, depth
						if !_rules[ruleStringChar]() {
							goto l82
						}
						goto l81
					l82:
						position, tokenIndex, depth = position82, tokenIndex82, depth82
					}
					if buffer[position] != rune('"') {
						goto l78
					}
					position++
					if !_rules[ruleSpacing]() {
						goto l78
					}
					goto l77
				l78:
					position, tokenIndex, depth = position78, tokenIndex78, depth78
				}
				depth--
				add(ruleStringLiteral, position76)
			}
			return true
		l75:
			position, tokenIndex, depth = position75, tokenIndex75, depth75
			return false
		},
		/* 15 StringChar <- <(Escape / (!('"' / '\n' / '\\') .))> */
		func() bool {
			position83, tokenIndex83, depth83 := position, tokenIndex, depth
			{
				position84 := position
				depth++
				{
					position85, tokenIndex85, depth85 := position, tokenIndex, depth
					if !_rules[ruleEscape]() {
						goto l86
					}
					goto l85
				l86:
					position, tokenIndex, depth = position85, tokenIndex85, depth85
					{
						position87, tokenIndex87, depth87 := position, tokenIndex, depth
						{
							position88, tokenIndex88, depth88 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l89
							}
							position++
							goto l88
						l89:
							position, tokenIndex, depth = position88, tokenIndex88, depth88
							if buffer[position] != rune('\n') {
								goto l90
							}
							position++
							goto l88
						l90:
							position, tokenIndex, depth = position88, tokenIndex88, depth88
							if buffer[position] != rune('\\') {
								goto l87
							}
							position++
						}
					l88:
						goto l83
					l87:
						position, tokenIndex, depth = position87, tokenIndex87, depth87
					}
					if !matchDot() {
						goto l83
					}
				}
			l85:
				depth--
				add(ruleStringChar, position84)
			}
			return true
		l83:
			position, tokenIndex, depth = position83, tokenIndex83, depth83
			return false
		},
		/* 16 Escape <- <('\\' ('\'' / '"' / '?' / '\\' / 'a' / 'b' / 'f' / 'n' / 'r' / 't' / 'v'))> */
		func() bool {
			position91, tokenIndex91, depth91 := position, tokenIndex, depth
			{
				position92 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l91
				}
				position++
				{
					position93, tokenIndex93, depth93 := position, tokenIndex, depth
					if buffer[position] != rune('\'') {
						goto l94
					}
					position++
					goto l93
				l94:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('"') {
						goto l95
					}
					position++
					goto l93
				l95:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('?') {
						goto l96
					}
					position++
					goto l93
				l96:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('\\') {
						goto l97
					}
					position++
					goto l93
				l97:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('a') {
						goto l98
					}
					position++
					goto l93
				l98:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('b') {
						goto l99
					}
					position++
					goto l93
				l99:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('f') {
						goto l100
					}
					position++
					goto l93
				l100:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('n') {
						goto l101
					}
					position++
					goto l93
				l101:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('r') {
						goto l102
					}
					position++
					goto l93
				l102:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('t') {
						goto l103
					}
					position++
					goto l93
				l103:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('v') {
						goto l91
					}
					position++
				}
			l93:
				depth--
				add(ruleEscape, position92)
			}
			return true
		l91:
			position, tokenIndex, depth = position91, tokenIndex91, depth91
			return false
		},
		/* 17 DecimalConstant <- <([0-9] [0-9]* Spacing)> */
		func() bool {
			position104, tokenIndex104, depth104 := position, tokenIndex, depth
			{
				position105 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l104
				}
				position++
			l106:
				{
					position107, tokenIndex107, depth107 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l107
					}
					position++
					goto l106
				l107:
					position, tokenIndex, depth = position107, tokenIndex107, depth107
				}
				if !_rules[ruleSpacing]() {
					goto l104
				}
				depth--
				add(ruleDecimalConstant, position105)
			}
			return true
		l104:
			position, tokenIndex, depth = position104, tokenIndex104, depth104
			return false
		},
		/* 18 LPAR <- <('(' Spacing)> */
		func() bool {
			position108, tokenIndex108, depth108 := position, tokenIndex, depth
			{
				position109 := position
				depth++
				if buffer[position] != rune('(') {
					goto l108
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l108
				}
				depth--
				add(ruleLPAR, position109)
			}
			return true
		l108:
			position, tokenIndex, depth = position108, tokenIndex108, depth108
			return false
		},
		/* 19 RPAR <- <(')' Spacing)> */
		func() bool {
			position110, tokenIndex110, depth110 := position, tokenIndex, depth
			{
				position111 := position
				depth++
				if buffer[position] != rune(')') {
					goto l110
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l110
				}
				depth--
				add(ruleRPAR, position111)
			}
			return true
		l110:
			position, tokenIndex, depth = position110, tokenIndex110, depth110
			return false
		},
		/* 20 LCURLY <- <('{' Spacing)> */
		func() bool {
			position112, tokenIndex112, depth112 := position, tokenIndex, depth
			{
				position113 := position
				depth++
				if buffer[position] != rune('{') {
					goto l112
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l112
				}
				depth--
				add(ruleLCURLY, position113)
			}
			return true
		l112:
			position, tokenIndex, depth = position112, tokenIndex112, depth112
			return false
		},
		/* 21 RCURLY <- <('}' Spacing)> */
		func() bool {
			position114, tokenIndex114, depth114 := position, tokenIndex, depth
			{
				position115 := position
				depth++
				if buffer[position] != rune('}') {
					goto l114
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l114
				}
				depth--
				add(ruleRCURLY, position115)
			}
			return true
		l114:
			position, tokenIndex, depth = position114, tokenIndex114, depth114
			return false
		},
		/* 22 PCOMMA <- <(';' Spacing)> */
		func() bool {
			position116, tokenIndex116, depth116 := position, tokenIndex, depth
			{
				position117 := position
				depth++
				if buffer[position] != rune(';') {
					goto l116
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l116
				}
				depth--
				add(rulePCOMMA, position117)
			}
			return true
		l116:
			position, tokenIndex, depth = position116, tokenIndex116, depth116
			return false
		},
		/* 23 EOT <- <!.> */
		func() bool {
			position118, tokenIndex118, depth118 := position, tokenIndex, depth
			{
				position119 := position
				depth++
				{
					position120, tokenIndex120, depth120 := position, tokenIndex, depth
					if !matchDot() {
						goto l120
					}
					goto l118
				l120:
					position, tokenIndex, depth = position120, tokenIndex120, depth120
				}
				depth--
				add(ruleEOT, position119)
			}
			return true
		l118:
			position, tokenIndex, depth = position118, tokenIndex118, depth118
			return false
		},
	}
	p.rules = _rules
}
