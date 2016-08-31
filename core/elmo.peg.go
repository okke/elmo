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
	rulePipedOutput
	ruleEndOfLine
	ruleShortcut
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
	ruleCOLON
	ruleDOT
	rulePIPE
	ruleEOT

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"Script",
	"Line",
	"PipedOutput",
	"EndOfLine",
	"Shortcut",
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
	"COLON",
	"DOT",
	"PIPE",
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
	rules  [30]func() bool
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
		/* 1 Line <- <(NewLine? Identifier Shortcut? Argument* (PipedOutput / EndOfLine)?)> */
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
				{
					position8, tokenIndex8, depth8 := position, tokenIndex, depth
					if !_rules[ruleShortcut]() {
						goto l8
					}
					goto l9
				l8:
					position, tokenIndex, depth = position8, tokenIndex8, depth8
				}
			l9:
			l10:
				{
					position11, tokenIndex11, depth11 := position, tokenIndex, depth
					if !_rules[ruleArgument]() {
						goto l11
					}
					goto l10
				l11:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
				}
				{
					position12, tokenIndex12, depth12 := position, tokenIndex, depth
					{
						position14, tokenIndex14, depth14 := position, tokenIndex, depth
						if !_rules[rulePipedOutput]() {
							goto l15
						}
						goto l14
					l15:
						position, tokenIndex, depth = position14, tokenIndex14, depth14
						if !_rules[ruleEndOfLine]() {
							goto l12
						}
					}
				l14:
					goto l13
				l12:
					position, tokenIndex, depth = position12, tokenIndex12, depth12
				}
			l13:
				depth--
				add(ruleLine, position5)
			}
			return true
		l4:
			position, tokenIndex, depth = position4, tokenIndex4, depth4
			return false
		},
		/* 2 PipedOutput <- <(PIPE Line)> */
		func() bool {
			position16, tokenIndex16, depth16 := position, tokenIndex, depth
			{
				position17 := position
				depth++
				if !_rules[rulePIPE]() {
					goto l16
				}
				if !_rules[ruleLine]() {
					goto l16
				}
				depth--
				add(rulePipedOutput, position17)
			}
			return true
		l16:
			position, tokenIndex, depth = position16, tokenIndex16, depth16
			return false
		},
		/* 3 EndOfLine <- <(PCOMMA / NewLine)> */
		func() bool {
			position18, tokenIndex18, depth18 := position, tokenIndex, depth
			{
				position19 := position
				depth++
				{
					position20, tokenIndex20, depth20 := position, tokenIndex, depth
					if !_rules[rulePCOMMA]() {
						goto l21
					}
					goto l20
				l21:
					position, tokenIndex, depth = position20, tokenIndex20, depth20
					if !_rules[ruleNewLine]() {
						goto l18
					}
				}
			l20:
				depth--
				add(ruleEndOfLine, position19)
			}
			return true
		l18:
			position, tokenIndex, depth = position18, tokenIndex18, depth18
			return false
		},
		/* 4 Shortcut <- <(COLON / DOT)> */
		func() bool {
			position22, tokenIndex22, depth22 := position, tokenIndex, depth
			{
				position23 := position
				depth++
				{
					position24, tokenIndex24, depth24 := position, tokenIndex, depth
					if !_rules[ruleCOLON]() {
						goto l25
					}
					goto l24
				l25:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if !_rules[ruleDOT]() {
						goto l22
					}
				}
			l24:
				depth--
				add(ruleShortcut, position23)
			}
			return true
		l22:
			position, tokenIndex, depth = position22, tokenIndex22, depth22
			return false
		},
		/* 5 Argument <- <(Identifier / StringLiteral / DecimalConstant / FunctionCall / Block)> */
		func() bool {
			position26, tokenIndex26, depth26 := position, tokenIndex, depth
			{
				position27 := position
				depth++
				{
					position28, tokenIndex28, depth28 := position, tokenIndex, depth
					if !_rules[ruleIdentifier]() {
						goto l29
					}
					goto l28
				l29:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if !_rules[ruleStringLiteral]() {
						goto l30
					}
					goto l28
				l30:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if !_rules[ruleDecimalConstant]() {
						goto l31
					}
					goto l28
				l31:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if !_rules[ruleFunctionCall]() {
						goto l32
					}
					goto l28
				l32:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if !_rules[ruleBlock]() {
						goto l26
					}
				}
			l28:
				depth--
				add(ruleArgument, position27)
			}
			return true
		l26:
			position, tokenIndex, depth = position26, tokenIndex26, depth26
			return false
		},
		/* 6 FunctionCall <- <(LPAR Line RPAR)> */
		func() bool {
			position33, tokenIndex33, depth33 := position, tokenIndex, depth
			{
				position34 := position
				depth++
				if !_rules[ruleLPAR]() {
					goto l33
				}
				if !_rules[ruleLine]() {
					goto l33
				}
				if !_rules[ruleRPAR]() {
					goto l33
				}
				depth--
				add(ruleFunctionCall, position34)
			}
			return true
		l33:
			position, tokenIndex, depth = position33, tokenIndex33, depth33
			return false
		},
		/* 7 Block <- <(LCURLY NewLine* Line* RCURLY)> */
		func() bool {
			position35, tokenIndex35, depth35 := position, tokenIndex, depth
			{
				position36 := position
				depth++
				if !_rules[ruleLCURLY]() {
					goto l35
				}
			l37:
				{
					position38, tokenIndex38, depth38 := position, tokenIndex, depth
					if !_rules[ruleNewLine]() {
						goto l38
					}
					goto l37
				l38:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
				}
			l39:
				{
					position40, tokenIndex40, depth40 := position, tokenIndex, depth
					if !_rules[ruleLine]() {
						goto l40
					}
					goto l39
				l40:
					position, tokenIndex, depth = position40, tokenIndex40, depth40
				}
				if !_rules[ruleRCURLY]() {
					goto l35
				}
				depth--
				add(ruleBlock, position36)
			}
			return true
		l35:
			position, tokenIndex, depth = position35, tokenIndex35, depth35
			return false
		},
		/* 8 Spacing <- <(WhiteSpace / LongComment / LineComment)*> */
		func() bool {
			{
				position42 := position
				depth++
			l43:
				{
					position44, tokenIndex44, depth44 := position, tokenIndex, depth
					{
						position45, tokenIndex45, depth45 := position, tokenIndex, depth
						if !_rules[ruleWhiteSpace]() {
							goto l46
						}
						goto l45
					l46:
						position, tokenIndex, depth = position45, tokenIndex45, depth45
						if !_rules[ruleLongComment]() {
							goto l47
						}
						goto l45
					l47:
						position, tokenIndex, depth = position45, tokenIndex45, depth45
						if !_rules[ruleLineComment]() {
							goto l44
						}
					}
				l45:
					goto l43
				l44:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
				}
				depth--
				add(ruleSpacing, position42)
			}
			return true
		},
		/* 9 WhiteSpace <- <(' ' / '\t')> */
		func() bool {
			position48, tokenIndex48, depth48 := position, tokenIndex, depth
			{
				position49 := position
				depth++
				{
					position50, tokenIndex50, depth50 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l51
					}
					position++
					goto l50
				l51:
					position, tokenIndex, depth = position50, tokenIndex50, depth50
					if buffer[position] != rune('\t') {
						goto l48
					}
					position++
				}
			l50:
				depth--
				add(ruleWhiteSpace, position49)
			}
			return true
		l48:
			position, tokenIndex, depth = position48, tokenIndex48, depth48
			return false
		},
		/* 10 LongComment <- <('/' '*' (!('*' '/') .)* ('*' '/'))> */
		func() bool {
			position52, tokenIndex52, depth52 := position, tokenIndex, depth
			{
				position53 := position
				depth++
				if buffer[position] != rune('/') {
					goto l52
				}
				position++
				if buffer[position] != rune('*') {
					goto l52
				}
				position++
			l54:
				{
					position55, tokenIndex55, depth55 := position, tokenIndex, depth
					{
						position56, tokenIndex56, depth56 := position, tokenIndex, depth
						if buffer[position] != rune('*') {
							goto l56
						}
						position++
						if buffer[position] != rune('/') {
							goto l56
						}
						position++
						goto l55
					l56:
						position, tokenIndex, depth = position56, tokenIndex56, depth56
					}
					if !matchDot() {
						goto l55
					}
					goto l54
				l55:
					position, tokenIndex, depth = position55, tokenIndex55, depth55
				}
				if buffer[position] != rune('*') {
					goto l52
				}
				position++
				if buffer[position] != rune('/') {
					goto l52
				}
				position++
				depth--
				add(ruleLongComment, position53)
			}
			return true
		l52:
			position, tokenIndex, depth = position52, tokenIndex52, depth52
			return false
		},
		/* 11 LineComment <- <('#' (!'\n' .)*)> */
		func() bool {
			position57, tokenIndex57, depth57 := position, tokenIndex, depth
			{
				position58 := position
				depth++
				if buffer[position] != rune('#') {
					goto l57
				}
				position++
			l59:
				{
					position60, tokenIndex60, depth60 := position, tokenIndex, depth
					{
						position61, tokenIndex61, depth61 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l61
						}
						position++
						goto l60
					l61:
						position, tokenIndex, depth = position61, tokenIndex61, depth61
					}
					if !matchDot() {
						goto l60
					}
					goto l59
				l60:
					position, tokenIndex, depth = position60, tokenIndex60, depth60
				}
				depth--
				add(ruleLineComment, position58)
			}
			return true
		l57:
			position, tokenIndex, depth = position57, tokenIndex57, depth57
			return false
		},
		/* 12 NewLine <- <(('\n' / '\r') Spacing)+> */
		func() bool {
			position62, tokenIndex62, depth62 := position, tokenIndex, depth
			{
				position63 := position
				depth++
				{
					position66, tokenIndex66, depth66 := position, tokenIndex, depth
					if buffer[position] != rune('\n') {
						goto l67
					}
					position++
					goto l66
				l67:
					position, tokenIndex, depth = position66, tokenIndex66, depth66
					if buffer[position] != rune('\r') {
						goto l62
					}
					position++
				}
			l66:
				if !_rules[ruleSpacing]() {
					goto l62
				}
			l64:
				{
					position65, tokenIndex65, depth65 := position, tokenIndex, depth
					{
						position68, tokenIndex68, depth68 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l69
						}
						position++
						goto l68
					l69:
						position, tokenIndex, depth = position68, tokenIndex68, depth68
						if buffer[position] != rune('\r') {
							goto l65
						}
						position++
					}
				l68:
					if !_rules[ruleSpacing]() {
						goto l65
					}
					goto l64
				l65:
					position, tokenIndex, depth = position65, tokenIndex65, depth65
				}
				depth--
				add(ruleNewLine, position63)
			}
			return true
		l62:
			position, tokenIndex, depth = position62, tokenIndex62, depth62
			return false
		},
		/* 13 Identifier <- <(IdNondigit IdChar* Spacing)> */
		func() bool {
			position70, tokenIndex70, depth70 := position, tokenIndex, depth
			{
				position71 := position
				depth++
				if !_rules[ruleIdNondigit]() {
					goto l70
				}
			l72:
				{
					position73, tokenIndex73, depth73 := position, tokenIndex, depth
					if !_rules[ruleIdChar]() {
						goto l73
					}
					goto l72
				l73:
					position, tokenIndex, depth = position73, tokenIndex73, depth73
				}
				if !_rules[ruleSpacing]() {
					goto l70
				}
				depth--
				add(ruleIdentifier, position71)
			}
			return true
		l70:
			position, tokenIndex, depth = position70, tokenIndex70, depth70
			return false
		},
		/* 14 IdNondigit <- <([a-z] / [A-Z] / '_')> */
		func() bool {
			position74, tokenIndex74, depth74 := position, tokenIndex, depth
			{
				position75 := position
				depth++
				{
					position76, tokenIndex76, depth76 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l77
					}
					position++
					goto l76
				l77:
					position, tokenIndex, depth = position76, tokenIndex76, depth76
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l78
					}
					position++
					goto l76
				l78:
					position, tokenIndex, depth = position76, tokenIndex76, depth76
					if buffer[position] != rune('_') {
						goto l74
					}
					position++
				}
			l76:
				depth--
				add(ruleIdNondigit, position75)
			}
			return true
		l74:
			position, tokenIndex, depth = position74, tokenIndex74, depth74
			return false
		},
		/* 15 IdChar <- <([a-z] / [A-Z] / [0-9] / '_')> */
		func() bool {
			position79, tokenIndex79, depth79 := position, tokenIndex, depth
			{
				position80 := position
				depth++
				{
					position81, tokenIndex81, depth81 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l82
					}
					position++
					goto l81
				l82:
					position, tokenIndex, depth = position81, tokenIndex81, depth81
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l83
					}
					position++
					goto l81
				l83:
					position, tokenIndex, depth = position81, tokenIndex81, depth81
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l84
					}
					position++
					goto l81
				l84:
					position, tokenIndex, depth = position81, tokenIndex81, depth81
					if buffer[position] != rune('_') {
						goto l79
					}
					position++
				}
			l81:
				depth--
				add(ruleIdChar, position80)
			}
			return true
		l79:
			position, tokenIndex, depth = position79, tokenIndex79, depth79
			return false
		},
		/* 16 StringLiteral <- <('"' StringChar* '"' Spacing)> */
		func() bool {
			position85, tokenIndex85, depth85 := position, tokenIndex, depth
			{
				position86 := position
				depth++
				if buffer[position] != rune('"') {
					goto l85
				}
				position++
			l87:
				{
					position88, tokenIndex88, depth88 := position, tokenIndex, depth
					if !_rules[ruleStringChar]() {
						goto l88
					}
					goto l87
				l88:
					position, tokenIndex, depth = position88, tokenIndex88, depth88
				}
				if buffer[position] != rune('"') {
					goto l85
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l85
				}
				depth--
				add(ruleStringLiteral, position86)
			}
			return true
		l85:
			position, tokenIndex, depth = position85, tokenIndex85, depth85
			return false
		},
		/* 17 StringChar <- <(Escape / (!('"' / '\n' / '\\') .))> */
		func() bool {
			position89, tokenIndex89, depth89 := position, tokenIndex, depth
			{
				position90 := position
				depth++
				{
					position91, tokenIndex91, depth91 := position, tokenIndex, depth
					if !_rules[ruleEscape]() {
						goto l92
					}
					goto l91
				l92:
					position, tokenIndex, depth = position91, tokenIndex91, depth91
					{
						position93, tokenIndex93, depth93 := position, tokenIndex, depth
						{
							position94, tokenIndex94, depth94 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l95
							}
							position++
							goto l94
						l95:
							position, tokenIndex, depth = position94, tokenIndex94, depth94
							if buffer[position] != rune('\n') {
								goto l96
							}
							position++
							goto l94
						l96:
							position, tokenIndex, depth = position94, tokenIndex94, depth94
							if buffer[position] != rune('\\') {
								goto l93
							}
							position++
						}
					l94:
						goto l89
					l93:
						position, tokenIndex, depth = position93, tokenIndex93, depth93
					}
					if !matchDot() {
						goto l89
					}
				}
			l91:
				depth--
				add(ruleStringChar, position90)
			}
			return true
		l89:
			position, tokenIndex, depth = position89, tokenIndex89, depth89
			return false
		},
		/* 18 Escape <- <('\\' ('\'' / '"' / '?' / '\\' / 'a' / 'b' / 'f' / 'n' / 'r' / 't' / 'v'))> */
		func() bool {
			position97, tokenIndex97, depth97 := position, tokenIndex, depth
			{
				position98 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l97
				}
				position++
				{
					position99, tokenIndex99, depth99 := position, tokenIndex, depth
					if buffer[position] != rune('\'') {
						goto l100
					}
					position++
					goto l99
				l100:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('"') {
						goto l101
					}
					position++
					goto l99
				l101:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('?') {
						goto l102
					}
					position++
					goto l99
				l102:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('\\') {
						goto l103
					}
					position++
					goto l99
				l103:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('a') {
						goto l104
					}
					position++
					goto l99
				l104:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('b') {
						goto l105
					}
					position++
					goto l99
				l105:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('f') {
						goto l106
					}
					position++
					goto l99
				l106:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('n') {
						goto l107
					}
					position++
					goto l99
				l107:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('r') {
						goto l108
					}
					position++
					goto l99
				l108:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('t') {
						goto l109
					}
					position++
					goto l99
				l109:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('v') {
						goto l97
					}
					position++
				}
			l99:
				depth--
				add(ruleEscape, position98)
			}
			return true
		l97:
			position, tokenIndex, depth = position97, tokenIndex97, depth97
			return false
		},
		/* 19 DecimalConstant <- <('-'? [0-9] [0-9]* Spacing)> */
		func() bool {
			position110, tokenIndex110, depth110 := position, tokenIndex, depth
			{
				position111 := position
				depth++
				{
					position112, tokenIndex112, depth112 := position, tokenIndex, depth
					if buffer[position] != rune('-') {
						goto l112
					}
					position++
					goto l113
				l112:
					position, tokenIndex, depth = position112, tokenIndex112, depth112
				}
			l113:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l110
				}
				position++
			l114:
				{
					position115, tokenIndex115, depth115 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l115
					}
					position++
					goto l114
				l115:
					position, tokenIndex, depth = position115, tokenIndex115, depth115
				}
				if !_rules[ruleSpacing]() {
					goto l110
				}
				depth--
				add(ruleDecimalConstant, position111)
			}
			return true
		l110:
			position, tokenIndex, depth = position110, tokenIndex110, depth110
			return false
		},
		/* 20 LPAR <- <('(' Spacing)> */
		func() bool {
			position116, tokenIndex116, depth116 := position, tokenIndex, depth
			{
				position117 := position
				depth++
				if buffer[position] != rune('(') {
					goto l116
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l116
				}
				depth--
				add(ruleLPAR, position117)
			}
			return true
		l116:
			position, tokenIndex, depth = position116, tokenIndex116, depth116
			return false
		},
		/* 21 RPAR <- <(')' Spacing)> */
		func() bool {
			position118, tokenIndex118, depth118 := position, tokenIndex, depth
			{
				position119 := position
				depth++
				if buffer[position] != rune(')') {
					goto l118
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l118
				}
				depth--
				add(ruleRPAR, position119)
			}
			return true
		l118:
			position, tokenIndex, depth = position118, tokenIndex118, depth118
			return false
		},
		/* 22 LCURLY <- <('{' Spacing)> */
		func() bool {
			position120, tokenIndex120, depth120 := position, tokenIndex, depth
			{
				position121 := position
				depth++
				if buffer[position] != rune('{') {
					goto l120
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l120
				}
				depth--
				add(ruleLCURLY, position121)
			}
			return true
		l120:
			position, tokenIndex, depth = position120, tokenIndex120, depth120
			return false
		},
		/* 23 RCURLY <- <('}' Spacing)> */
		func() bool {
			position122, tokenIndex122, depth122 := position, tokenIndex, depth
			{
				position123 := position
				depth++
				if buffer[position] != rune('}') {
					goto l122
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l122
				}
				depth--
				add(ruleRCURLY, position123)
			}
			return true
		l122:
			position, tokenIndex, depth = position122, tokenIndex122, depth122
			return false
		},
		/* 24 PCOMMA <- <(';' Spacing)> */
		func() bool {
			position124, tokenIndex124, depth124 := position, tokenIndex, depth
			{
				position125 := position
				depth++
				if buffer[position] != rune(';') {
					goto l124
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l124
				}
				depth--
				add(rulePCOMMA, position125)
			}
			return true
		l124:
			position, tokenIndex, depth = position124, tokenIndex124, depth124
			return false
		},
		/* 25 COLON <- <(':' Spacing)> */
		func() bool {
			position126, tokenIndex126, depth126 := position, tokenIndex, depth
			{
				position127 := position
				depth++
				if buffer[position] != rune(':') {
					goto l126
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l126
				}
				depth--
				add(ruleCOLON, position127)
			}
			return true
		l126:
			position, tokenIndex, depth = position126, tokenIndex126, depth126
			return false
		},
		/* 26 DOT <- <('.' Spacing)> */
		func() bool {
			position128, tokenIndex128, depth128 := position, tokenIndex, depth
			{
				position129 := position
				depth++
				if buffer[position] != rune('.') {
					goto l128
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l128
				}
				depth--
				add(ruleDOT, position129)
			}
			return true
		l128:
			position, tokenIndex, depth = position128, tokenIndex128, depth128
			return false
		},
		/* 27 PIPE <- <('|' Spacing)> */
		func() bool {
			position130, tokenIndex130, depth130 := position, tokenIndex, depth
			{
				position131 := position
				depth++
				if buffer[position] != rune('|') {
					goto l130
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l130
				}
				depth--
				add(rulePIPE, position131)
			}
			return true
		l130:
			position, tokenIndex, depth = position130, tokenIndex130, depth130
			return false
		},
		/* 28 EOT <- <!.> */
		func() bool {
			position132, tokenIndex132, depth132 := position, tokenIndex, depth
			{
				position133 := position
				depth++
				{
					position134, tokenIndex134, depth134 := position, tokenIndex, depth
					if !matchDot() {
						goto l134
					}
					goto l132
				l134:
					position, tokenIndex, depth = position134, tokenIndex134, depth134
				}
				depth--
				add(ruleEOT, position133)
			}
			return true
		l132:
			position, tokenIndex, depth = position132, tokenIndex132, depth132
			return false
		},
	}
	p.rules = _rules
}
