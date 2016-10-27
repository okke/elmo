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
	ruleArgument
	ruleFunctionCall
	ruleBlock
	ruleList
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
	ruleLongStringLiteral
	ruleEscape
	ruleNumber
	ruleLPAR
	ruleRPAR
	ruleLCURLY
	ruleRCURLY
	ruleLBRACKET
	ruleRBRACKET
	ruleCOMMA
	rulePCOMMA
	ruleCOLON
	ruleDOT
	rulePIPE
	ruleDOLLAR
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
	"Argument",
	"FunctionCall",
	"Block",
	"List",
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
	"LongStringLiteral",
	"Escape",
	"Number",
	"LPAR",
	"RPAR",
	"LCURLY",
	"RCURLY",
	"LBRACKET",
	"RBRACKET",
	"COMMA",
	"PCOMMA",
	"COLON",
	"DOT",
	"PIPE",
	"DOLLAR",
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
	rules  [35]func() bool
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
		/* 1 Line <- <(NewLine? Argument COLON? Argument? ((COMMA NewLine?)? Argument)* (PipedOutput / EndOfLine)?)> */
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
				if !_rules[ruleArgument]() {
					goto l4
				}
				{
					position8, tokenIndex8, depth8 := position, tokenIndex, depth
					if !_rules[ruleCOLON]() {
						goto l8
					}
					goto l9
				l8:
					position, tokenIndex, depth = position8, tokenIndex8, depth8
				}
			l9:
				{
					position10, tokenIndex10, depth10 := position, tokenIndex, depth
					if !_rules[ruleArgument]() {
						goto l10
					}
					goto l11
				l10:
					position, tokenIndex, depth = position10, tokenIndex10, depth10
				}
			l11:
			l12:
				{
					position13, tokenIndex13, depth13 := position, tokenIndex, depth
					{
						position14, tokenIndex14, depth14 := position, tokenIndex, depth
						if !_rules[ruleCOMMA]() {
							goto l14
						}
						{
							position16, tokenIndex16, depth16 := position, tokenIndex, depth
							if !_rules[ruleNewLine]() {
								goto l16
							}
							goto l17
						l16:
							position, tokenIndex, depth = position16, tokenIndex16, depth16
						}
					l17:
						goto l15
					l14:
						position, tokenIndex, depth = position14, tokenIndex14, depth14
					}
				l15:
					if !_rules[ruleArgument]() {
						goto l13
					}
					goto l12
				l13:
					position, tokenIndex, depth = position13, tokenIndex13, depth13
				}
				{
					position18, tokenIndex18, depth18 := position, tokenIndex, depth
					{
						position20, tokenIndex20, depth20 := position, tokenIndex, depth
						if !_rules[rulePipedOutput]() {
							goto l21
						}
						goto l20
					l21:
						position, tokenIndex, depth = position20, tokenIndex20, depth20
						if !_rules[ruleEndOfLine]() {
							goto l18
						}
					}
				l20:
					goto l19
				l18:
					position, tokenIndex, depth = position18, tokenIndex18, depth18
				}
			l19:
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
			position22, tokenIndex22, depth22 := position, tokenIndex, depth
			{
				position23 := position
				depth++
				if !_rules[rulePIPE]() {
					goto l22
				}
				if !_rules[ruleLine]() {
					goto l22
				}
				depth--
				add(rulePipedOutput, position23)
			}
			return true
		l22:
			position, tokenIndex, depth = position22, tokenIndex22, depth22
			return false
		},
		/* 3 EndOfLine <- <(PCOMMA / NewLine)> */
		func() bool {
			position24, tokenIndex24, depth24 := position, tokenIndex, depth
			{
				position25 := position
				depth++
				{
					position26, tokenIndex26, depth26 := position, tokenIndex, depth
					if !_rules[rulePCOMMA]() {
						goto l27
					}
					goto l26
				l27:
					position, tokenIndex, depth = position26, tokenIndex26, depth26
					if !_rules[ruleNewLine]() {
						goto l24
					}
				}
			l26:
				depth--
				add(ruleEndOfLine, position25)
			}
			return true
		l24:
			position, tokenIndex, depth = position24, tokenIndex24, depth24
			return false
		},
		/* 4 Argument <- <((Identifier (DOT Identifier)?) / StringLiteral / LongStringLiteral / Number / FunctionCall / Block / List)> */
		func() bool {
			position28, tokenIndex28, depth28 := position, tokenIndex, depth
			{
				position29 := position
				depth++
				{
					position30, tokenIndex30, depth30 := position, tokenIndex, depth
					if !_rules[ruleIdentifier]() {
						goto l31
					}
					{
						position32, tokenIndex32, depth32 := position, tokenIndex, depth
						if !_rules[ruleDOT]() {
							goto l32
						}
						if !_rules[ruleIdentifier]() {
							goto l32
						}
						goto l33
					l32:
						position, tokenIndex, depth = position32, tokenIndex32, depth32
					}
				l33:
					goto l30
				l31:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if !_rules[ruleStringLiteral]() {
						goto l34
					}
					goto l30
				l34:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if !_rules[ruleLongStringLiteral]() {
						goto l35
					}
					goto l30
				l35:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if !_rules[ruleNumber]() {
						goto l36
					}
					goto l30
				l36:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if !_rules[ruleFunctionCall]() {
						goto l37
					}
					goto l30
				l37:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if !_rules[ruleBlock]() {
						goto l38
					}
					goto l30
				l38:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if !_rules[ruleList]() {
						goto l28
					}
				}
			l30:
				depth--
				add(ruleArgument, position29)
			}
			return true
		l28:
			position, tokenIndex, depth = position28, tokenIndex28, depth28
			return false
		},
		/* 5 FunctionCall <- <((LPAR Line RPAR) / (DOLLAR Argument (DOT Argument)?))> */
		func() bool {
			position39, tokenIndex39, depth39 := position, tokenIndex, depth
			{
				position40 := position
				depth++
				{
					position41, tokenIndex41, depth41 := position, tokenIndex, depth
					if !_rules[ruleLPAR]() {
						goto l42
					}
					if !_rules[ruleLine]() {
						goto l42
					}
					if !_rules[ruleRPAR]() {
						goto l42
					}
					goto l41
				l42:
					position, tokenIndex, depth = position41, tokenIndex41, depth41
					if !_rules[ruleDOLLAR]() {
						goto l39
					}
					if !_rules[ruleArgument]() {
						goto l39
					}
					{
						position43, tokenIndex43, depth43 := position, tokenIndex, depth
						if !_rules[ruleDOT]() {
							goto l43
						}
						if !_rules[ruleArgument]() {
							goto l43
						}
						goto l44
					l43:
						position, tokenIndex, depth = position43, tokenIndex43, depth43
					}
				l44:
				}
			l41:
				depth--
				add(ruleFunctionCall, position40)
			}
			return true
		l39:
			position, tokenIndex, depth = position39, tokenIndex39, depth39
			return false
		},
		/* 6 Block <- <(LCURLY NewLine* Line* RCURLY)> */
		func() bool {
			position45, tokenIndex45, depth45 := position, tokenIndex, depth
			{
				position46 := position
				depth++
				if !_rules[ruleLCURLY]() {
					goto l45
				}
			l47:
				{
					position48, tokenIndex48, depth48 := position, tokenIndex, depth
					if !_rules[ruleNewLine]() {
						goto l48
					}
					goto l47
				l48:
					position, tokenIndex, depth = position48, tokenIndex48, depth48
				}
			l49:
				{
					position50, tokenIndex50, depth50 := position, tokenIndex, depth
					if !_rules[ruleLine]() {
						goto l50
					}
					goto l49
				l50:
					position, tokenIndex, depth = position50, tokenIndex50, depth50
				}
				if !_rules[ruleRCURLY]() {
					goto l45
				}
				depth--
				add(ruleBlock, position46)
			}
			return true
		l45:
			position, tokenIndex, depth = position45, tokenIndex45, depth45
			return false
		},
		/* 7 List <- <(LBRACKET NewLine* (Argument / NewLine)? (((COMMA NewLine?)? Argument) / NewLine)* RBRACKET)> */
		func() bool {
			position51, tokenIndex51, depth51 := position, tokenIndex, depth
			{
				position52 := position
				depth++
				if !_rules[ruleLBRACKET]() {
					goto l51
				}
			l53:
				{
					position54, tokenIndex54, depth54 := position, tokenIndex, depth
					if !_rules[ruleNewLine]() {
						goto l54
					}
					goto l53
				l54:
					position, tokenIndex, depth = position54, tokenIndex54, depth54
				}
				{
					position55, tokenIndex55, depth55 := position, tokenIndex, depth
					{
						position57, tokenIndex57, depth57 := position, tokenIndex, depth
						if !_rules[ruleArgument]() {
							goto l58
						}
						goto l57
					l58:
						position, tokenIndex, depth = position57, tokenIndex57, depth57
						if !_rules[ruleNewLine]() {
							goto l55
						}
					}
				l57:
					goto l56
				l55:
					position, tokenIndex, depth = position55, tokenIndex55, depth55
				}
			l56:
			l59:
				{
					position60, tokenIndex60, depth60 := position, tokenIndex, depth
					{
						position61, tokenIndex61, depth61 := position, tokenIndex, depth
						{
							position63, tokenIndex63, depth63 := position, tokenIndex, depth
							if !_rules[ruleCOMMA]() {
								goto l63
							}
							{
								position65, tokenIndex65, depth65 := position, tokenIndex, depth
								if !_rules[ruleNewLine]() {
									goto l65
								}
								goto l66
							l65:
								position, tokenIndex, depth = position65, tokenIndex65, depth65
							}
						l66:
							goto l64
						l63:
							position, tokenIndex, depth = position63, tokenIndex63, depth63
						}
					l64:
						if !_rules[ruleArgument]() {
							goto l62
						}
						goto l61
					l62:
						position, tokenIndex, depth = position61, tokenIndex61, depth61
						if !_rules[ruleNewLine]() {
							goto l60
						}
					}
				l61:
					goto l59
				l60:
					position, tokenIndex, depth = position60, tokenIndex60, depth60
				}
				if !_rules[ruleRBRACKET]() {
					goto l51
				}
				depth--
				add(ruleList, position52)
			}
			return true
		l51:
			position, tokenIndex, depth = position51, tokenIndex51, depth51
			return false
		},
		/* 8 Spacing <- <(WhiteSpace / LongComment / LineComment)*> */
		func() bool {
			{
				position68 := position
				depth++
			l69:
				{
					position70, tokenIndex70, depth70 := position, tokenIndex, depth
					{
						position71, tokenIndex71, depth71 := position, tokenIndex, depth
						if !_rules[ruleWhiteSpace]() {
							goto l72
						}
						goto l71
					l72:
						position, tokenIndex, depth = position71, tokenIndex71, depth71
						if !_rules[ruleLongComment]() {
							goto l73
						}
						goto l71
					l73:
						position, tokenIndex, depth = position71, tokenIndex71, depth71
						if !_rules[ruleLineComment]() {
							goto l70
						}
					}
				l71:
					goto l69
				l70:
					position, tokenIndex, depth = position70, tokenIndex70, depth70
				}
				depth--
				add(ruleSpacing, position68)
			}
			return true
		},
		/* 9 WhiteSpace <- <(' ' / '\t')> */
		func() bool {
			position74, tokenIndex74, depth74 := position, tokenIndex, depth
			{
				position75 := position
				depth++
				{
					position76, tokenIndex76, depth76 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l77
					}
					position++
					goto l76
				l77:
					position, tokenIndex, depth = position76, tokenIndex76, depth76
					if buffer[position] != rune('\t') {
						goto l74
					}
					position++
				}
			l76:
				depth--
				add(ruleWhiteSpace, position75)
			}
			return true
		l74:
			position, tokenIndex, depth = position74, tokenIndex74, depth74
			return false
		},
		/* 10 LongComment <- <('/' '*' (!('*' '/') .)* ('*' '/'))> */
		func() bool {
			position78, tokenIndex78, depth78 := position, tokenIndex, depth
			{
				position79 := position
				depth++
				if buffer[position] != rune('/') {
					goto l78
				}
				position++
				if buffer[position] != rune('*') {
					goto l78
				}
				position++
			l80:
				{
					position81, tokenIndex81, depth81 := position, tokenIndex, depth
					{
						position82, tokenIndex82, depth82 := position, tokenIndex, depth
						if buffer[position] != rune('*') {
							goto l82
						}
						position++
						if buffer[position] != rune('/') {
							goto l82
						}
						position++
						goto l81
					l82:
						position, tokenIndex, depth = position82, tokenIndex82, depth82
					}
					if !matchDot() {
						goto l81
					}
					goto l80
				l81:
					position, tokenIndex, depth = position81, tokenIndex81, depth81
				}
				if buffer[position] != rune('*') {
					goto l78
				}
				position++
				if buffer[position] != rune('/') {
					goto l78
				}
				position++
				depth--
				add(ruleLongComment, position79)
			}
			return true
		l78:
			position, tokenIndex, depth = position78, tokenIndex78, depth78
			return false
		},
		/* 11 LineComment <- <('#' (!'\n' .)*)> */
		func() bool {
			position83, tokenIndex83, depth83 := position, tokenIndex, depth
			{
				position84 := position
				depth++
				if buffer[position] != rune('#') {
					goto l83
				}
				position++
			l85:
				{
					position86, tokenIndex86, depth86 := position, tokenIndex, depth
					{
						position87, tokenIndex87, depth87 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l87
						}
						position++
						goto l86
					l87:
						position, tokenIndex, depth = position87, tokenIndex87, depth87
					}
					if !matchDot() {
						goto l86
					}
					goto l85
				l86:
					position, tokenIndex, depth = position86, tokenIndex86, depth86
				}
				depth--
				add(ruleLineComment, position84)
			}
			return true
		l83:
			position, tokenIndex, depth = position83, tokenIndex83, depth83
			return false
		},
		/* 12 NewLine <- <(('\n' / '\r') Spacing)+> */
		func() bool {
			position88, tokenIndex88, depth88 := position, tokenIndex, depth
			{
				position89 := position
				depth++
				{
					position92, tokenIndex92, depth92 := position, tokenIndex, depth
					if buffer[position] != rune('\n') {
						goto l93
					}
					position++
					goto l92
				l93:
					position, tokenIndex, depth = position92, tokenIndex92, depth92
					if buffer[position] != rune('\r') {
						goto l88
					}
					position++
				}
			l92:
				if !_rules[ruleSpacing]() {
					goto l88
				}
			l90:
				{
					position91, tokenIndex91, depth91 := position, tokenIndex, depth
					{
						position94, tokenIndex94, depth94 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l95
						}
						position++
						goto l94
					l95:
						position, tokenIndex, depth = position94, tokenIndex94, depth94
						if buffer[position] != rune('\r') {
							goto l91
						}
						position++
					}
				l94:
					if !_rules[ruleSpacing]() {
						goto l91
					}
					goto l90
				l91:
					position, tokenIndex, depth = position91, tokenIndex91, depth91
				}
				depth--
				add(ruleNewLine, position89)
			}
			return true
		l88:
			position, tokenIndex, depth = position88, tokenIndex88, depth88
			return false
		},
		/* 13 Identifier <- <(IdNondigit IdChar* Spacing)> */
		func() bool {
			position96, tokenIndex96, depth96 := position, tokenIndex, depth
			{
				position97 := position
				depth++
				if !_rules[ruleIdNondigit]() {
					goto l96
				}
			l98:
				{
					position99, tokenIndex99, depth99 := position, tokenIndex, depth
					if !_rules[ruleIdChar]() {
						goto l99
					}
					goto l98
				l99:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
				}
				if !_rules[ruleSpacing]() {
					goto l96
				}
				depth--
				add(ruleIdentifier, position97)
			}
			return true
		l96:
			position, tokenIndex, depth = position96, tokenIndex96, depth96
			return false
		},
		/* 14 IdNondigit <- <([a-z] / [A-Z] / ('_' / '?'))> */
		func() bool {
			position100, tokenIndex100, depth100 := position, tokenIndex, depth
			{
				position101 := position
				depth++
				{
					position102, tokenIndex102, depth102 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l103
					}
					position++
					goto l102
				l103:
					position, tokenIndex, depth = position102, tokenIndex102, depth102
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l104
					}
					position++
					goto l102
				l104:
					position, tokenIndex, depth = position102, tokenIndex102, depth102
					{
						position105, tokenIndex105, depth105 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l106
						}
						position++
						goto l105
					l106:
						position, tokenIndex, depth = position105, tokenIndex105, depth105
						if buffer[position] != rune('?') {
							goto l100
						}
						position++
					}
				l105:
				}
			l102:
				depth--
				add(ruleIdNondigit, position101)
			}
			return true
		l100:
			position, tokenIndex, depth = position100, tokenIndex100, depth100
			return false
		},
		/* 15 IdChar <- <([a-z] / [A-Z] / [0-9] / ('_' / '?'))> */
		func() bool {
			position107, tokenIndex107, depth107 := position, tokenIndex, depth
			{
				position108 := position
				depth++
				{
					position109, tokenIndex109, depth109 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l110
					}
					position++
					goto l109
				l110:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l111
					}
					position++
					goto l109
				l111:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l112
					}
					position++
					goto l109
				l112:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
					{
						position113, tokenIndex113, depth113 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l114
						}
						position++
						goto l113
					l114:
						position, tokenIndex, depth = position113, tokenIndex113, depth113
						if buffer[position] != rune('?') {
							goto l107
						}
						position++
					}
				l113:
				}
			l109:
				depth--
				add(ruleIdChar, position108)
			}
			return true
		l107:
			position, tokenIndex, depth = position107, tokenIndex107, depth107
			return false
		},
		/* 16 StringLiteral <- <('"' StringChar* '"' Spacing)> */
		func() bool {
			position115, tokenIndex115, depth115 := position, tokenIndex, depth
			{
				position116 := position
				depth++
				if buffer[position] != rune('"') {
					goto l115
				}
				position++
			l117:
				{
					position118, tokenIndex118, depth118 := position, tokenIndex, depth
					if !_rules[ruleStringChar]() {
						goto l118
					}
					goto l117
				l118:
					position, tokenIndex, depth = position118, tokenIndex118, depth118
				}
				if buffer[position] != rune('"') {
					goto l115
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l115
				}
				depth--
				add(ruleStringLiteral, position116)
			}
			return true
		l115:
			position, tokenIndex, depth = position115, tokenIndex115, depth115
			return false
		},
		/* 17 StringChar <- <(Escape / (!('"' / '\n' / '\\') .))> */
		func() bool {
			position119, tokenIndex119, depth119 := position, tokenIndex, depth
			{
				position120 := position
				depth++
				{
					position121, tokenIndex121, depth121 := position, tokenIndex, depth
					if !_rules[ruleEscape]() {
						goto l122
					}
					goto l121
				l122:
					position, tokenIndex, depth = position121, tokenIndex121, depth121
					{
						position123, tokenIndex123, depth123 := position, tokenIndex, depth
						{
							position124, tokenIndex124, depth124 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l125
							}
							position++
							goto l124
						l125:
							position, tokenIndex, depth = position124, tokenIndex124, depth124
							if buffer[position] != rune('\n') {
								goto l126
							}
							position++
							goto l124
						l126:
							position, tokenIndex, depth = position124, tokenIndex124, depth124
							if buffer[position] != rune('\\') {
								goto l123
							}
							position++
						}
					l124:
						goto l119
					l123:
						position, tokenIndex, depth = position123, tokenIndex123, depth123
					}
					if !matchDot() {
						goto l119
					}
				}
			l121:
				depth--
				add(ruleStringChar, position120)
			}
			return true
		l119:
			position, tokenIndex, depth = position119, tokenIndex119, depth119
			return false
		},
		/* 18 LongStringLiteral <- <('`' (('`' '`') / (!'`' .))* '`' Spacing)> */
		func() bool {
			position127, tokenIndex127, depth127 := position, tokenIndex, depth
			{
				position128 := position
				depth++
				if buffer[position] != rune('`') {
					goto l127
				}
				position++
			l129:
				{
					position130, tokenIndex130, depth130 := position, tokenIndex, depth
					{
						position131, tokenIndex131, depth131 := position, tokenIndex, depth
						if buffer[position] != rune('`') {
							goto l132
						}
						position++
						if buffer[position] != rune('`') {
							goto l132
						}
						position++
						goto l131
					l132:
						position, tokenIndex, depth = position131, tokenIndex131, depth131
						{
							position133, tokenIndex133, depth133 := position, tokenIndex, depth
							if buffer[position] != rune('`') {
								goto l133
							}
							position++
							goto l130
						l133:
							position, tokenIndex, depth = position133, tokenIndex133, depth133
						}
						if !matchDot() {
							goto l130
						}
					}
				l131:
					goto l129
				l130:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
				}
				if buffer[position] != rune('`') {
					goto l127
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l127
				}
				depth--
				add(ruleLongStringLiteral, position128)
			}
			return true
		l127:
			position, tokenIndex, depth = position127, tokenIndex127, depth127
			return false
		},
		/* 19 Escape <- <('\\' .)> */
		func() bool {
			position134, tokenIndex134, depth134 := position, tokenIndex, depth
			{
				position135 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l134
				}
				position++
				if !matchDot() {
					goto l134
				}
				depth--
				add(ruleEscape, position135)
			}
			return true
		l134:
			position, tokenIndex, depth = position134, tokenIndex134, depth134
			return false
		},
		/* 20 Number <- <('-'? [0-9] [0-9]* ('.' [0-9] [0-9]*)? Spacing)> */
		func() bool {
			position136, tokenIndex136, depth136 := position, tokenIndex, depth
			{
				position137 := position
				depth++
				{
					position138, tokenIndex138, depth138 := position, tokenIndex, depth
					if buffer[position] != rune('-') {
						goto l138
					}
					position++
					goto l139
				l138:
					position, tokenIndex, depth = position138, tokenIndex138, depth138
				}
			l139:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l136
				}
				position++
			l140:
				{
					position141, tokenIndex141, depth141 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l141
					}
					position++
					goto l140
				l141:
					position, tokenIndex, depth = position141, tokenIndex141, depth141
				}
				{
					position142, tokenIndex142, depth142 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l142
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l142
					}
					position++
				l144:
					{
						position145, tokenIndex145, depth145 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l145
						}
						position++
						goto l144
					l145:
						position, tokenIndex, depth = position145, tokenIndex145, depth145
					}
					goto l143
				l142:
					position, tokenIndex, depth = position142, tokenIndex142, depth142
				}
			l143:
				if !_rules[ruleSpacing]() {
					goto l136
				}
				depth--
				add(ruleNumber, position137)
			}
			return true
		l136:
			position, tokenIndex, depth = position136, tokenIndex136, depth136
			return false
		},
		/* 21 LPAR <- <('(' Spacing)> */
		func() bool {
			position146, tokenIndex146, depth146 := position, tokenIndex, depth
			{
				position147 := position
				depth++
				if buffer[position] != rune('(') {
					goto l146
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l146
				}
				depth--
				add(ruleLPAR, position147)
			}
			return true
		l146:
			position, tokenIndex, depth = position146, tokenIndex146, depth146
			return false
		},
		/* 22 RPAR <- <(')' Spacing)> */
		func() bool {
			position148, tokenIndex148, depth148 := position, tokenIndex, depth
			{
				position149 := position
				depth++
				if buffer[position] != rune(')') {
					goto l148
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l148
				}
				depth--
				add(ruleRPAR, position149)
			}
			return true
		l148:
			position, tokenIndex, depth = position148, tokenIndex148, depth148
			return false
		},
		/* 23 LCURLY <- <('{' Spacing)> */
		func() bool {
			position150, tokenIndex150, depth150 := position, tokenIndex, depth
			{
				position151 := position
				depth++
				if buffer[position] != rune('{') {
					goto l150
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l150
				}
				depth--
				add(ruleLCURLY, position151)
			}
			return true
		l150:
			position, tokenIndex, depth = position150, tokenIndex150, depth150
			return false
		},
		/* 24 RCURLY <- <('}' Spacing)> */
		func() bool {
			position152, tokenIndex152, depth152 := position, tokenIndex, depth
			{
				position153 := position
				depth++
				if buffer[position] != rune('}') {
					goto l152
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l152
				}
				depth--
				add(ruleRCURLY, position153)
			}
			return true
		l152:
			position, tokenIndex, depth = position152, tokenIndex152, depth152
			return false
		},
		/* 25 LBRACKET <- <('[' Spacing)> */
		func() bool {
			position154, tokenIndex154, depth154 := position, tokenIndex, depth
			{
				position155 := position
				depth++
				if buffer[position] != rune('[') {
					goto l154
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l154
				}
				depth--
				add(ruleLBRACKET, position155)
			}
			return true
		l154:
			position, tokenIndex, depth = position154, tokenIndex154, depth154
			return false
		},
		/* 26 RBRACKET <- <(']' Spacing)> */
		func() bool {
			position156, tokenIndex156, depth156 := position, tokenIndex, depth
			{
				position157 := position
				depth++
				if buffer[position] != rune(']') {
					goto l156
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l156
				}
				depth--
				add(ruleRBRACKET, position157)
			}
			return true
		l156:
			position, tokenIndex, depth = position156, tokenIndex156, depth156
			return false
		},
		/* 27 COMMA <- <(',' Spacing)> */
		func() bool {
			position158, tokenIndex158, depth158 := position, tokenIndex, depth
			{
				position159 := position
				depth++
				if buffer[position] != rune(',') {
					goto l158
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l158
				}
				depth--
				add(ruleCOMMA, position159)
			}
			return true
		l158:
			position, tokenIndex, depth = position158, tokenIndex158, depth158
			return false
		},
		/* 28 PCOMMA <- <(';' Spacing)> */
		func() bool {
			position160, tokenIndex160, depth160 := position, tokenIndex, depth
			{
				position161 := position
				depth++
				if buffer[position] != rune(';') {
					goto l160
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l160
				}
				depth--
				add(rulePCOMMA, position161)
			}
			return true
		l160:
			position, tokenIndex, depth = position160, tokenIndex160, depth160
			return false
		},
		/* 29 COLON <- <(':' Spacing)> */
		func() bool {
			position162, tokenIndex162, depth162 := position, tokenIndex, depth
			{
				position163 := position
				depth++
				if buffer[position] != rune(':') {
					goto l162
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l162
				}
				depth--
				add(ruleCOLON, position163)
			}
			return true
		l162:
			position, tokenIndex, depth = position162, tokenIndex162, depth162
			return false
		},
		/* 30 DOT <- <('.' Spacing)> */
		func() bool {
			position164, tokenIndex164, depth164 := position, tokenIndex, depth
			{
				position165 := position
				depth++
				if buffer[position] != rune('.') {
					goto l164
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l164
				}
				depth--
				add(ruleDOT, position165)
			}
			return true
		l164:
			position, tokenIndex, depth = position164, tokenIndex164, depth164
			return false
		},
		/* 31 PIPE <- <('|' Spacing)> */
		func() bool {
			position166, tokenIndex166, depth166 := position, tokenIndex, depth
			{
				position167 := position
				depth++
				if buffer[position] != rune('|') {
					goto l166
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l166
				}
				depth--
				add(rulePIPE, position167)
			}
			return true
		l166:
			position, tokenIndex, depth = position166, tokenIndex166, depth166
			return false
		},
		/* 32 DOLLAR <- <('$' Spacing)> */
		func() bool {
			position168, tokenIndex168, depth168 := position, tokenIndex, depth
			{
				position169 := position
				depth++
				if buffer[position] != rune('$') {
					goto l168
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l168
				}
				depth--
				add(ruleDOLLAR, position169)
			}
			return true
		l168:
			position, tokenIndex, depth = position168, tokenIndex168, depth168
			return false
		},
		/* 33 EOT <- <!.> */
		func() bool {
			position170, tokenIndex170, depth170 := position, tokenIndex, depth
			{
				position171 := position
				depth++
				{
					position172, tokenIndex172, depth172 := position, tokenIndex, depth
					if !matchDot() {
						goto l172
					}
					goto l170
				l172:
					position, tokenIndex, depth = position172, tokenIndex172, depth172
				}
				depth--
				add(ruleEOT, position171)
			}
			return true
		l170:
			position, tokenIndex, depth = position170, tokenIndex170, depth170
			return false
		},
	}
	p.rules = _rules
}
