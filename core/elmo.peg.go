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
	"Shortcut",
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
		/* 1 Line <- <(NewLine? Argument Shortcut? Argument? ((COMMA NewLine?)? Argument)* (PipedOutput / EndOfLine)?)> */
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
					if !_rules[ruleShortcut]() {
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
		/* 4 Shortcut <- <(COLON / DOT)> */
		func() bool {
			position28, tokenIndex28, depth28 := position, tokenIndex, depth
			{
				position29 := position
				depth++
				{
					position30, tokenIndex30, depth30 := position, tokenIndex, depth
					if !_rules[ruleCOLON]() {
						goto l31
					}
					goto l30
				l31:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if !_rules[ruleDOT]() {
						goto l28
					}
				}
			l30:
				depth--
				add(ruleShortcut, position29)
			}
			return true
		l28:
			position, tokenIndex, depth = position28, tokenIndex28, depth28
			return false
		},
		/* 5 Argument <- <(Identifier / StringLiteral / Number / FunctionCall / Block / List)> */
		func() bool {
			position32, tokenIndex32, depth32 := position, tokenIndex, depth
			{
				position33 := position
				depth++
				{
					position34, tokenIndex34, depth34 := position, tokenIndex, depth
					if !_rules[ruleIdentifier]() {
						goto l35
					}
					goto l34
				l35:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
					if !_rules[ruleStringLiteral]() {
						goto l36
					}
					goto l34
				l36:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
					if !_rules[ruleNumber]() {
						goto l37
					}
					goto l34
				l37:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
					if !_rules[ruleFunctionCall]() {
						goto l38
					}
					goto l34
				l38:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
					if !_rules[ruleBlock]() {
						goto l39
					}
					goto l34
				l39:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
					if !_rules[ruleList]() {
						goto l32
					}
				}
			l34:
				depth--
				add(ruleArgument, position33)
			}
			return true
		l32:
			position, tokenIndex, depth = position32, tokenIndex32, depth32
			return false
		},
		/* 6 FunctionCall <- <((LPAR Line RPAR) / (DOLLAR Argument (DOT Argument)?))> */
		func() bool {
			position40, tokenIndex40, depth40 := position, tokenIndex, depth
			{
				position41 := position
				depth++
				{
					position42, tokenIndex42, depth42 := position, tokenIndex, depth
					if !_rules[ruleLPAR]() {
						goto l43
					}
					if !_rules[ruleLine]() {
						goto l43
					}
					if !_rules[ruleRPAR]() {
						goto l43
					}
					goto l42
				l43:
					position, tokenIndex, depth = position42, tokenIndex42, depth42
					if !_rules[ruleDOLLAR]() {
						goto l40
					}
					if !_rules[ruleArgument]() {
						goto l40
					}
					{
						position44, tokenIndex44, depth44 := position, tokenIndex, depth
						if !_rules[ruleDOT]() {
							goto l44
						}
						if !_rules[ruleArgument]() {
							goto l44
						}
						goto l45
					l44:
						position, tokenIndex, depth = position44, tokenIndex44, depth44
					}
				l45:
				}
			l42:
				depth--
				add(ruleFunctionCall, position41)
			}
			return true
		l40:
			position, tokenIndex, depth = position40, tokenIndex40, depth40
			return false
		},
		/* 7 Block <- <(LCURLY NewLine* Line* RCURLY)> */
		func() bool {
			position46, tokenIndex46, depth46 := position, tokenIndex, depth
			{
				position47 := position
				depth++
				if !_rules[ruleLCURLY]() {
					goto l46
				}
			l48:
				{
					position49, tokenIndex49, depth49 := position, tokenIndex, depth
					if !_rules[ruleNewLine]() {
						goto l49
					}
					goto l48
				l49:
					position, tokenIndex, depth = position49, tokenIndex49, depth49
				}
			l50:
				{
					position51, tokenIndex51, depth51 := position, tokenIndex, depth
					if !_rules[ruleLine]() {
						goto l51
					}
					goto l50
				l51:
					position, tokenIndex, depth = position51, tokenIndex51, depth51
				}
				if !_rules[ruleRCURLY]() {
					goto l46
				}
				depth--
				add(ruleBlock, position47)
			}
			return true
		l46:
			position, tokenIndex, depth = position46, tokenIndex46, depth46
			return false
		},
		/* 8 List <- <(LBRACKET NewLine* (Argument / NewLine)? (((COMMA NewLine?)? Argument) / NewLine)* RBRACKET)> */
		func() bool {
			position52, tokenIndex52, depth52 := position, tokenIndex, depth
			{
				position53 := position
				depth++
				if !_rules[ruleLBRACKET]() {
					goto l52
				}
			l54:
				{
					position55, tokenIndex55, depth55 := position, tokenIndex, depth
					if !_rules[ruleNewLine]() {
						goto l55
					}
					goto l54
				l55:
					position, tokenIndex, depth = position55, tokenIndex55, depth55
				}
				{
					position56, tokenIndex56, depth56 := position, tokenIndex, depth
					{
						position58, tokenIndex58, depth58 := position, tokenIndex, depth
						if !_rules[ruleArgument]() {
							goto l59
						}
						goto l58
					l59:
						position, tokenIndex, depth = position58, tokenIndex58, depth58
						if !_rules[ruleNewLine]() {
							goto l56
						}
					}
				l58:
					goto l57
				l56:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
				}
			l57:
			l60:
				{
					position61, tokenIndex61, depth61 := position, tokenIndex, depth
					{
						position62, tokenIndex62, depth62 := position, tokenIndex, depth
						{
							position64, tokenIndex64, depth64 := position, tokenIndex, depth
							if !_rules[ruleCOMMA]() {
								goto l64
							}
							{
								position66, tokenIndex66, depth66 := position, tokenIndex, depth
								if !_rules[ruleNewLine]() {
									goto l66
								}
								goto l67
							l66:
								position, tokenIndex, depth = position66, tokenIndex66, depth66
							}
						l67:
							goto l65
						l64:
							position, tokenIndex, depth = position64, tokenIndex64, depth64
						}
					l65:
						if !_rules[ruleArgument]() {
							goto l63
						}
						goto l62
					l63:
						position, tokenIndex, depth = position62, tokenIndex62, depth62
						if !_rules[ruleNewLine]() {
							goto l61
						}
					}
				l62:
					goto l60
				l61:
					position, tokenIndex, depth = position61, tokenIndex61, depth61
				}
				if !_rules[ruleRBRACKET]() {
					goto l52
				}
				depth--
				add(ruleList, position53)
			}
			return true
		l52:
			position, tokenIndex, depth = position52, tokenIndex52, depth52
			return false
		},
		/* 9 Spacing <- <(WhiteSpace / LongComment / LineComment)*> */
		func() bool {
			{
				position69 := position
				depth++
			l70:
				{
					position71, tokenIndex71, depth71 := position, tokenIndex, depth
					{
						position72, tokenIndex72, depth72 := position, tokenIndex, depth
						if !_rules[ruleWhiteSpace]() {
							goto l73
						}
						goto l72
					l73:
						position, tokenIndex, depth = position72, tokenIndex72, depth72
						if !_rules[ruleLongComment]() {
							goto l74
						}
						goto l72
					l74:
						position, tokenIndex, depth = position72, tokenIndex72, depth72
						if !_rules[ruleLineComment]() {
							goto l71
						}
					}
				l72:
					goto l70
				l71:
					position, tokenIndex, depth = position71, tokenIndex71, depth71
				}
				depth--
				add(ruleSpacing, position69)
			}
			return true
		},
		/* 10 WhiteSpace <- <(' ' / '\t')> */
		func() bool {
			position75, tokenIndex75, depth75 := position, tokenIndex, depth
			{
				position76 := position
				depth++
				{
					position77, tokenIndex77, depth77 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l78
					}
					position++
					goto l77
				l78:
					position, tokenIndex, depth = position77, tokenIndex77, depth77
					if buffer[position] != rune('\t') {
						goto l75
					}
					position++
				}
			l77:
				depth--
				add(ruleWhiteSpace, position76)
			}
			return true
		l75:
			position, tokenIndex, depth = position75, tokenIndex75, depth75
			return false
		},
		/* 11 LongComment <- <('/' '*' (!('*' '/') .)* ('*' '/'))> */
		func() bool {
			position79, tokenIndex79, depth79 := position, tokenIndex, depth
			{
				position80 := position
				depth++
				if buffer[position] != rune('/') {
					goto l79
				}
				position++
				if buffer[position] != rune('*') {
					goto l79
				}
				position++
			l81:
				{
					position82, tokenIndex82, depth82 := position, tokenIndex, depth
					{
						position83, tokenIndex83, depth83 := position, tokenIndex, depth
						if buffer[position] != rune('*') {
							goto l83
						}
						position++
						if buffer[position] != rune('/') {
							goto l83
						}
						position++
						goto l82
					l83:
						position, tokenIndex, depth = position83, tokenIndex83, depth83
					}
					if !matchDot() {
						goto l82
					}
					goto l81
				l82:
					position, tokenIndex, depth = position82, tokenIndex82, depth82
				}
				if buffer[position] != rune('*') {
					goto l79
				}
				position++
				if buffer[position] != rune('/') {
					goto l79
				}
				position++
				depth--
				add(ruleLongComment, position80)
			}
			return true
		l79:
			position, tokenIndex, depth = position79, tokenIndex79, depth79
			return false
		},
		/* 12 LineComment <- <('#' (!'\n' .)*)> */
		func() bool {
			position84, tokenIndex84, depth84 := position, tokenIndex, depth
			{
				position85 := position
				depth++
				if buffer[position] != rune('#') {
					goto l84
				}
				position++
			l86:
				{
					position87, tokenIndex87, depth87 := position, tokenIndex, depth
					{
						position88, tokenIndex88, depth88 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l88
						}
						position++
						goto l87
					l88:
						position, tokenIndex, depth = position88, tokenIndex88, depth88
					}
					if !matchDot() {
						goto l87
					}
					goto l86
				l87:
					position, tokenIndex, depth = position87, tokenIndex87, depth87
				}
				depth--
				add(ruleLineComment, position85)
			}
			return true
		l84:
			position, tokenIndex, depth = position84, tokenIndex84, depth84
			return false
		},
		/* 13 NewLine <- <(('\n' / '\r') Spacing)+> */
		func() bool {
			position89, tokenIndex89, depth89 := position, tokenIndex, depth
			{
				position90 := position
				depth++
				{
					position93, tokenIndex93, depth93 := position, tokenIndex, depth
					if buffer[position] != rune('\n') {
						goto l94
					}
					position++
					goto l93
				l94:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('\r') {
						goto l89
					}
					position++
				}
			l93:
				if !_rules[ruleSpacing]() {
					goto l89
				}
			l91:
				{
					position92, tokenIndex92, depth92 := position, tokenIndex, depth
					{
						position95, tokenIndex95, depth95 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l96
						}
						position++
						goto l95
					l96:
						position, tokenIndex, depth = position95, tokenIndex95, depth95
						if buffer[position] != rune('\r') {
							goto l92
						}
						position++
					}
				l95:
					if !_rules[ruleSpacing]() {
						goto l92
					}
					goto l91
				l92:
					position, tokenIndex, depth = position92, tokenIndex92, depth92
				}
				depth--
				add(ruleNewLine, position90)
			}
			return true
		l89:
			position, tokenIndex, depth = position89, tokenIndex89, depth89
			return false
		},
		/* 14 Identifier <- <(IdNondigit IdChar* Spacing)> */
		func() bool {
			position97, tokenIndex97, depth97 := position, tokenIndex, depth
			{
				position98 := position
				depth++
				if !_rules[ruleIdNondigit]() {
					goto l97
				}
			l99:
				{
					position100, tokenIndex100, depth100 := position, tokenIndex, depth
					if !_rules[ruleIdChar]() {
						goto l100
					}
					goto l99
				l100:
					position, tokenIndex, depth = position100, tokenIndex100, depth100
				}
				if !_rules[ruleSpacing]() {
					goto l97
				}
				depth--
				add(ruleIdentifier, position98)
			}
			return true
		l97:
			position, tokenIndex, depth = position97, tokenIndex97, depth97
			return false
		},
		/* 15 IdNondigit <- <([a-z] / [A-Z] / ('_' / '?'))> */
		func() bool {
			position101, tokenIndex101, depth101 := position, tokenIndex, depth
			{
				position102 := position
				depth++
				{
					position103, tokenIndex103, depth103 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l104
					}
					position++
					goto l103
				l104:
					position, tokenIndex, depth = position103, tokenIndex103, depth103
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l105
					}
					position++
					goto l103
				l105:
					position, tokenIndex, depth = position103, tokenIndex103, depth103
					{
						position106, tokenIndex106, depth106 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l107
						}
						position++
						goto l106
					l107:
						position, tokenIndex, depth = position106, tokenIndex106, depth106
						if buffer[position] != rune('?') {
							goto l101
						}
						position++
					}
				l106:
				}
			l103:
				depth--
				add(ruleIdNondigit, position102)
			}
			return true
		l101:
			position, tokenIndex, depth = position101, tokenIndex101, depth101
			return false
		},
		/* 16 IdChar <- <([a-z] / [A-Z] / [0-9] / ('_' / '?'))> */
		func() bool {
			position108, tokenIndex108, depth108 := position, tokenIndex, depth
			{
				position109 := position
				depth++
				{
					position110, tokenIndex110, depth110 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l111
					}
					position++
					goto l110
				l111:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l112
					}
					position++
					goto l110
				l112:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l113
					}
					position++
					goto l110
				l113:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					{
						position114, tokenIndex114, depth114 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l115
						}
						position++
						goto l114
					l115:
						position, tokenIndex, depth = position114, tokenIndex114, depth114
						if buffer[position] != rune('?') {
							goto l108
						}
						position++
					}
				l114:
				}
			l110:
				depth--
				add(ruleIdChar, position109)
			}
			return true
		l108:
			position, tokenIndex, depth = position108, tokenIndex108, depth108
			return false
		},
		/* 17 StringLiteral <- <('"' StringChar* '"' Spacing)> */
		func() bool {
			position116, tokenIndex116, depth116 := position, tokenIndex, depth
			{
				position117 := position
				depth++
				if buffer[position] != rune('"') {
					goto l116
				}
				position++
			l118:
				{
					position119, tokenIndex119, depth119 := position, tokenIndex, depth
					if !_rules[ruleStringChar]() {
						goto l119
					}
					goto l118
				l119:
					position, tokenIndex, depth = position119, tokenIndex119, depth119
				}
				if buffer[position] != rune('"') {
					goto l116
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l116
				}
				depth--
				add(ruleStringLiteral, position117)
			}
			return true
		l116:
			position, tokenIndex, depth = position116, tokenIndex116, depth116
			return false
		},
		/* 18 StringChar <- <(Escape / (!('"' / '\n' / '\\') .))> */
		func() bool {
			position120, tokenIndex120, depth120 := position, tokenIndex, depth
			{
				position121 := position
				depth++
				{
					position122, tokenIndex122, depth122 := position, tokenIndex, depth
					if !_rules[ruleEscape]() {
						goto l123
					}
					goto l122
				l123:
					position, tokenIndex, depth = position122, tokenIndex122, depth122
					{
						position124, tokenIndex124, depth124 := position, tokenIndex, depth
						{
							position125, tokenIndex125, depth125 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l126
							}
							position++
							goto l125
						l126:
							position, tokenIndex, depth = position125, tokenIndex125, depth125
							if buffer[position] != rune('\n') {
								goto l127
							}
							position++
							goto l125
						l127:
							position, tokenIndex, depth = position125, tokenIndex125, depth125
							if buffer[position] != rune('\\') {
								goto l124
							}
							position++
						}
					l125:
						goto l120
					l124:
						position, tokenIndex, depth = position124, tokenIndex124, depth124
					}
					if !matchDot() {
						goto l120
					}
				}
			l122:
				depth--
				add(ruleStringChar, position121)
			}
			return true
		l120:
			position, tokenIndex, depth = position120, tokenIndex120, depth120
			return false
		},
		/* 19 Escape <- <('\\' ('\'' / '"' / '?' / '\\' / 'a' / 'b' / 'f' / 'n' / 'r' / 't' / 'v'))> */
		func() bool {
			position128, tokenIndex128, depth128 := position, tokenIndex, depth
			{
				position129 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l128
				}
				position++
				{
					position130, tokenIndex130, depth130 := position, tokenIndex, depth
					if buffer[position] != rune('\'') {
						goto l131
					}
					position++
					goto l130
				l131:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('"') {
						goto l132
					}
					position++
					goto l130
				l132:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('?') {
						goto l133
					}
					position++
					goto l130
				l133:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('\\') {
						goto l134
					}
					position++
					goto l130
				l134:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('a') {
						goto l135
					}
					position++
					goto l130
				l135:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('b') {
						goto l136
					}
					position++
					goto l130
				l136:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('f') {
						goto l137
					}
					position++
					goto l130
				l137:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('n') {
						goto l138
					}
					position++
					goto l130
				l138:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('r') {
						goto l139
					}
					position++
					goto l130
				l139:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('t') {
						goto l140
					}
					position++
					goto l130
				l140:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('v') {
						goto l128
					}
					position++
				}
			l130:
				depth--
				add(ruleEscape, position129)
			}
			return true
		l128:
			position, tokenIndex, depth = position128, tokenIndex128, depth128
			return false
		},
		/* 20 Number <- <('-'? [0-9] [0-9]* ('.' [0-9] [0-9]*)? Spacing)> */
		func() bool {
			position141, tokenIndex141, depth141 := position, tokenIndex, depth
			{
				position142 := position
				depth++
				{
					position143, tokenIndex143, depth143 := position, tokenIndex, depth
					if buffer[position] != rune('-') {
						goto l143
					}
					position++
					goto l144
				l143:
					position, tokenIndex, depth = position143, tokenIndex143, depth143
				}
			l144:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l141
				}
				position++
			l145:
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l146
					}
					position++
					goto l145
				l146:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
				}
				{
					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l147
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l147
					}
					position++
				l149:
					{
						position150, tokenIndex150, depth150 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l150
						}
						position++
						goto l149
					l150:
						position, tokenIndex, depth = position150, tokenIndex150, depth150
					}
					goto l148
				l147:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
				}
			l148:
				if !_rules[ruleSpacing]() {
					goto l141
				}
				depth--
				add(ruleNumber, position142)
			}
			return true
		l141:
			position, tokenIndex, depth = position141, tokenIndex141, depth141
			return false
		},
		/* 21 LPAR <- <('(' Spacing)> */
		func() bool {
			position151, tokenIndex151, depth151 := position, tokenIndex, depth
			{
				position152 := position
				depth++
				if buffer[position] != rune('(') {
					goto l151
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l151
				}
				depth--
				add(ruleLPAR, position152)
			}
			return true
		l151:
			position, tokenIndex, depth = position151, tokenIndex151, depth151
			return false
		},
		/* 22 RPAR <- <(')' Spacing)> */
		func() bool {
			position153, tokenIndex153, depth153 := position, tokenIndex, depth
			{
				position154 := position
				depth++
				if buffer[position] != rune(')') {
					goto l153
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l153
				}
				depth--
				add(ruleRPAR, position154)
			}
			return true
		l153:
			position, tokenIndex, depth = position153, tokenIndex153, depth153
			return false
		},
		/* 23 LCURLY <- <('{' Spacing)> */
		func() bool {
			position155, tokenIndex155, depth155 := position, tokenIndex, depth
			{
				position156 := position
				depth++
				if buffer[position] != rune('{') {
					goto l155
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l155
				}
				depth--
				add(ruleLCURLY, position156)
			}
			return true
		l155:
			position, tokenIndex, depth = position155, tokenIndex155, depth155
			return false
		},
		/* 24 RCURLY <- <('}' Spacing)> */
		func() bool {
			position157, tokenIndex157, depth157 := position, tokenIndex, depth
			{
				position158 := position
				depth++
				if buffer[position] != rune('}') {
					goto l157
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l157
				}
				depth--
				add(ruleRCURLY, position158)
			}
			return true
		l157:
			position, tokenIndex, depth = position157, tokenIndex157, depth157
			return false
		},
		/* 25 LBRACKET <- <('[' Spacing)> */
		func() bool {
			position159, tokenIndex159, depth159 := position, tokenIndex, depth
			{
				position160 := position
				depth++
				if buffer[position] != rune('[') {
					goto l159
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l159
				}
				depth--
				add(ruleLBRACKET, position160)
			}
			return true
		l159:
			position, tokenIndex, depth = position159, tokenIndex159, depth159
			return false
		},
		/* 26 RBRACKET <- <(']' Spacing)> */
		func() bool {
			position161, tokenIndex161, depth161 := position, tokenIndex, depth
			{
				position162 := position
				depth++
				if buffer[position] != rune(']') {
					goto l161
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l161
				}
				depth--
				add(ruleRBRACKET, position162)
			}
			return true
		l161:
			position, tokenIndex, depth = position161, tokenIndex161, depth161
			return false
		},
		/* 27 COMMA <- <(',' Spacing)> */
		func() bool {
			position163, tokenIndex163, depth163 := position, tokenIndex, depth
			{
				position164 := position
				depth++
				if buffer[position] != rune(',') {
					goto l163
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l163
				}
				depth--
				add(ruleCOMMA, position164)
			}
			return true
		l163:
			position, tokenIndex, depth = position163, tokenIndex163, depth163
			return false
		},
		/* 28 PCOMMA <- <(';' Spacing)> */
		func() bool {
			position165, tokenIndex165, depth165 := position, tokenIndex, depth
			{
				position166 := position
				depth++
				if buffer[position] != rune(';') {
					goto l165
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l165
				}
				depth--
				add(rulePCOMMA, position166)
			}
			return true
		l165:
			position, tokenIndex, depth = position165, tokenIndex165, depth165
			return false
		},
		/* 29 COLON <- <(':' Spacing)> */
		func() bool {
			position167, tokenIndex167, depth167 := position, tokenIndex, depth
			{
				position168 := position
				depth++
				if buffer[position] != rune(':') {
					goto l167
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l167
				}
				depth--
				add(ruleCOLON, position168)
			}
			return true
		l167:
			position, tokenIndex, depth = position167, tokenIndex167, depth167
			return false
		},
		/* 30 DOT <- <('.' Spacing)> */
		func() bool {
			position169, tokenIndex169, depth169 := position, tokenIndex, depth
			{
				position170 := position
				depth++
				if buffer[position] != rune('.') {
					goto l169
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l169
				}
				depth--
				add(ruleDOT, position170)
			}
			return true
		l169:
			position, tokenIndex, depth = position169, tokenIndex169, depth169
			return false
		},
		/* 31 PIPE <- <('|' Spacing)> */
		func() bool {
			position171, tokenIndex171, depth171 := position, tokenIndex, depth
			{
				position172 := position
				depth++
				if buffer[position] != rune('|') {
					goto l171
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l171
				}
				depth--
				add(rulePIPE, position172)
			}
			return true
		l171:
			position, tokenIndex, depth = position171, tokenIndex171, depth171
			return false
		},
		/* 32 DOLLAR <- <('$' Spacing)> */
		func() bool {
			position173, tokenIndex173, depth173 := position, tokenIndex, depth
			{
				position174 := position
				depth++
				if buffer[position] != rune('$') {
					goto l173
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l173
				}
				depth--
				add(ruleDOLLAR, position174)
			}
			return true
		l173:
			position, tokenIndex, depth = position173, tokenIndex173, depth173
			return false
		},
		/* 33 EOT <- <!.> */
		func() bool {
			position175, tokenIndex175, depth175 := position, tokenIndex, depth
			{
				position176 := position
				depth++
				{
					position177, tokenIndex177, depth177 := position, tokenIndex, depth
					if !matchDot() {
						goto l177
					}
					goto l175
				l177:
					position, tokenIndex, depth = position177, tokenIndex177, depth177
				}
				depth--
				add(ruleEOT, position176)
			}
			return true
		l175:
			position, tokenIndex, depth = position175, tokenIndex175, depth175
			return false
		},
	}
	p.rules = _rules
}
