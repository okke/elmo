package elmo

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

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
	ruleQuote
	ruleStringChar
	ruleLongStringLiteral
	ruleLongString
	ruleEscape
	ruleNumber
	ruleNumbers
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
	ruleAMPERSAND
	ruleEOT

	rulePre
	ruleIn
	ruleSuf
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
	"Quote",
	"StringChar",
	"LongStringLiteral",
	"LongString",
	"Escape",
	"Number",
	"Numbers",
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
	"AMPERSAND",
	"EOT",

	"Pre_",
	"_In_",
	"_Suf",
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

func (node *node32) Print(buffer string) {
	node.print(0, buffer)
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
		for i := range states {
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
							write(token32{pegRule: ruleIn, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre, begin: a.begin, end: b.begin}, true)
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
					write(token32{pegRule: ruleSuf, begin: b.end, end: a.end}, true)
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
	for i := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

func (t *tokens32) Expand(index int) {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
}

type ElmoGrammar struct {
	Buffer string
	buffer []rune
	rules  [39]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	Pretty bool
	tokens32
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
	p.tokens32.PrintSyntaxTree(p.Buffer)
}

func (p *ElmoGrammar) Highlighter() {
	p.PrintSyntax()
}

func (p *ElmoGrammar) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
		p.buffer = append(p.buffer, endSymbol)
	}

	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	var max token32
	position, depth, tokenIndex, buffer, _rules := uint32(0), uint32(0), 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin uint32) {
		tree.Expand(tokenIndex)
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position, depth}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
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
		/* 4 Argument <- <((Identifier (DOT Identifier)*) / StringLiteral / LongStringLiteral / Number / FunctionCall / Block / List)> */
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
				l32:
					{
						position33, tokenIndex33, depth33 := position, tokenIndex, depth
						if !_rules[ruleDOT]() {
							goto l33
						}
						if !_rules[ruleIdentifier]() {
							goto l33
						}
						goto l32
					l33:
						position, tokenIndex, depth = position33, tokenIndex33, depth33
					}
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
		/* 5 FunctionCall <- <((LPAR Line RPAR) / ((DOLLAR / AMPERSAND) Argument (DOT Argument)*))> */
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
					{
						position43, tokenIndex43, depth43 := position, tokenIndex, depth
						if !_rules[ruleDOLLAR]() {
							goto l44
						}
						goto l43
					l44:
						position, tokenIndex, depth = position43, tokenIndex43, depth43
						if !_rules[ruleAMPERSAND]() {
							goto l39
						}
					}
				l43:
					if !_rules[ruleArgument]() {
						goto l39
					}
				l45:
					{
						position46, tokenIndex46, depth46 := position, tokenIndex, depth
						if !_rules[ruleDOT]() {
							goto l46
						}
						if !_rules[ruleArgument]() {
							goto l46
						}
						goto l45
					l46:
						position, tokenIndex, depth = position46, tokenIndex46, depth46
					}
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
			position47, tokenIndex47, depth47 := position, tokenIndex, depth
			{
				position48 := position
				depth++
				if !_rules[ruleLCURLY]() {
					goto l47
				}
			l49:
				{
					position50, tokenIndex50, depth50 := position, tokenIndex, depth
					if !_rules[ruleNewLine]() {
						goto l50
					}
					goto l49
				l50:
					position, tokenIndex, depth = position50, tokenIndex50, depth50
				}
			l51:
				{
					position52, tokenIndex52, depth52 := position, tokenIndex, depth
					if !_rules[ruleLine]() {
						goto l52
					}
					goto l51
				l52:
					position, tokenIndex, depth = position52, tokenIndex52, depth52
				}
				if !_rules[ruleRCURLY]() {
					goto l47
				}
				depth--
				add(ruleBlock, position48)
			}
			return true
		l47:
			position, tokenIndex, depth = position47, tokenIndex47, depth47
			return false
		},
		/* 7 List <- <(LBRACKET NewLine* (Argument / NewLine)? (((COMMA NewLine?)? Argument) / NewLine)* RBRACKET)> */
		func() bool {
			position53, tokenIndex53, depth53 := position, tokenIndex, depth
			{
				position54 := position
				depth++
				if !_rules[ruleLBRACKET]() {
					goto l53
				}
			l55:
				{
					position56, tokenIndex56, depth56 := position, tokenIndex, depth
					if !_rules[ruleNewLine]() {
						goto l56
					}
					goto l55
				l56:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
				}
				{
					position57, tokenIndex57, depth57 := position, tokenIndex, depth
					{
						position59, tokenIndex59, depth59 := position, tokenIndex, depth
						if !_rules[ruleArgument]() {
							goto l60
						}
						goto l59
					l60:
						position, tokenIndex, depth = position59, tokenIndex59, depth59
						if !_rules[ruleNewLine]() {
							goto l57
						}
					}
				l59:
					goto l58
				l57:
					position, tokenIndex, depth = position57, tokenIndex57, depth57
				}
			l58:
			l61:
				{
					position62, tokenIndex62, depth62 := position, tokenIndex, depth
					{
						position63, tokenIndex63, depth63 := position, tokenIndex, depth
						{
							position65, tokenIndex65, depth65 := position, tokenIndex, depth
							if !_rules[ruleCOMMA]() {
								goto l65
							}
							{
								position67, tokenIndex67, depth67 := position, tokenIndex, depth
								if !_rules[ruleNewLine]() {
									goto l67
								}
								goto l68
							l67:
								position, tokenIndex, depth = position67, tokenIndex67, depth67
							}
						l68:
							goto l66
						l65:
							position, tokenIndex, depth = position65, tokenIndex65, depth65
						}
					l66:
						if !_rules[ruleArgument]() {
							goto l64
						}
						goto l63
					l64:
						position, tokenIndex, depth = position63, tokenIndex63, depth63
						if !_rules[ruleNewLine]() {
							goto l62
						}
					}
				l63:
					goto l61
				l62:
					position, tokenIndex, depth = position62, tokenIndex62, depth62
				}
				if !_rules[ruleRBRACKET]() {
					goto l53
				}
				depth--
				add(ruleList, position54)
			}
			return true
		l53:
			position, tokenIndex, depth = position53, tokenIndex53, depth53
			return false
		},
		/* 8 Spacing <- <(WhiteSpace / LongComment / LineComment)*> */
		func() bool {
			{
				position70 := position
				depth++
			l71:
				{
					position72, tokenIndex72, depth72 := position, tokenIndex, depth
					{
						position73, tokenIndex73, depth73 := position, tokenIndex, depth
						if !_rules[ruleWhiteSpace]() {
							goto l74
						}
						goto l73
					l74:
						position, tokenIndex, depth = position73, tokenIndex73, depth73
						if !_rules[ruleLongComment]() {
							goto l75
						}
						goto l73
					l75:
						position, tokenIndex, depth = position73, tokenIndex73, depth73
						if !_rules[ruleLineComment]() {
							goto l72
						}
					}
				l73:
					goto l71
				l72:
					position, tokenIndex, depth = position72, tokenIndex72, depth72
				}
				depth--
				add(ruleSpacing, position70)
			}
			return true
		},
		/* 9 WhiteSpace <- <(' ' / '\t')> */
		func() bool {
			position76, tokenIndex76, depth76 := position, tokenIndex, depth
			{
				position77 := position
				depth++
				{
					position78, tokenIndex78, depth78 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l79
					}
					position++
					goto l78
				l79:
					position, tokenIndex, depth = position78, tokenIndex78, depth78
					if buffer[position] != rune('\t') {
						goto l76
					}
					position++
				}
			l78:
				depth--
				add(ruleWhiteSpace, position77)
			}
			return true
		l76:
			position, tokenIndex, depth = position76, tokenIndex76, depth76
			return false
		},
		/* 10 LongComment <- <('/' '*' (!('*' '/') .)* ('*' '/'))> */
		func() bool {
			position80, tokenIndex80, depth80 := position, tokenIndex, depth
			{
				position81 := position
				depth++
				if buffer[position] != rune('/') {
					goto l80
				}
				position++
				if buffer[position] != rune('*') {
					goto l80
				}
				position++
			l82:
				{
					position83, tokenIndex83, depth83 := position, tokenIndex, depth
					{
						position84, tokenIndex84, depth84 := position, tokenIndex, depth
						if buffer[position] != rune('*') {
							goto l84
						}
						position++
						if buffer[position] != rune('/') {
							goto l84
						}
						position++
						goto l83
					l84:
						position, tokenIndex, depth = position84, tokenIndex84, depth84
					}
					if !matchDot() {
						goto l83
					}
					goto l82
				l83:
					position, tokenIndex, depth = position83, tokenIndex83, depth83
				}
				if buffer[position] != rune('*') {
					goto l80
				}
				position++
				if buffer[position] != rune('/') {
					goto l80
				}
				position++
				depth--
				add(ruleLongComment, position81)
			}
			return true
		l80:
			position, tokenIndex, depth = position80, tokenIndex80, depth80
			return false
		},
		/* 11 LineComment <- <('#' (!'\n' .)*)> */
		func() bool {
			position85, tokenIndex85, depth85 := position, tokenIndex, depth
			{
				position86 := position
				depth++
				if buffer[position] != rune('#') {
					goto l85
				}
				position++
			l87:
				{
					position88, tokenIndex88, depth88 := position, tokenIndex, depth
					{
						position89, tokenIndex89, depth89 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l89
						}
						position++
						goto l88
					l89:
						position, tokenIndex, depth = position89, tokenIndex89, depth89
					}
					if !matchDot() {
						goto l88
					}
					goto l87
				l88:
					position, tokenIndex, depth = position88, tokenIndex88, depth88
				}
				depth--
				add(ruleLineComment, position86)
			}
			return true
		l85:
			position, tokenIndex, depth = position85, tokenIndex85, depth85
			return false
		},
		/* 12 NewLine <- <(('\n' / '\r') Spacing)+> */
		func() bool {
			position90, tokenIndex90, depth90 := position, tokenIndex, depth
			{
				position91 := position
				depth++
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
						goto l90
					}
					position++
				}
			l94:
				if !_rules[ruleSpacing]() {
					goto l90
				}
			l92:
				{
					position93, tokenIndex93, depth93 := position, tokenIndex, depth
					{
						position96, tokenIndex96, depth96 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l97
						}
						position++
						goto l96
					l97:
						position, tokenIndex, depth = position96, tokenIndex96, depth96
						if buffer[position] != rune('\r') {
							goto l93
						}
						position++
					}
				l96:
					if !_rules[ruleSpacing]() {
						goto l93
					}
					goto l92
				l93:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
				}
				depth--
				add(ruleNewLine, position91)
			}
			return true
		l90:
			position, tokenIndex, depth = position90, tokenIndex90, depth90
			return false
		},
		/* 13 Identifier <- <(IdNondigit IdChar* Spacing)> */
		func() bool {
			position98, tokenIndex98, depth98 := position, tokenIndex, depth
			{
				position99 := position
				depth++
				if !_rules[ruleIdNondigit]() {
					goto l98
				}
			l100:
				{
					position101, tokenIndex101, depth101 := position, tokenIndex, depth
					if !_rules[ruleIdChar]() {
						goto l101
					}
					goto l100
				l101:
					position, tokenIndex, depth = position101, tokenIndex101, depth101
				}
				if !_rules[ruleSpacing]() {
					goto l98
				}
				depth--
				add(ruleIdentifier, position99)
			}
			return true
		l98:
			position, tokenIndex, depth = position98, tokenIndex98, depth98
			return false
		},
		/* 14 IdNondigit <- <([a-z] / [A-Z] / ('_' / '?' / '!'))> */
		func() bool {
			position102, tokenIndex102, depth102 := position, tokenIndex, depth
			{
				position103 := position
				depth++
				{
					position104, tokenIndex104, depth104 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l105
					}
					position++
					goto l104
				l105:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l106
					}
					position++
					goto l104
				l106:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
					{
						position107, tokenIndex107, depth107 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l108
						}
						position++
						goto l107
					l108:
						position, tokenIndex, depth = position107, tokenIndex107, depth107
						if buffer[position] != rune('?') {
							goto l109
						}
						position++
						goto l107
					l109:
						position, tokenIndex, depth = position107, tokenIndex107, depth107
						if buffer[position] != rune('!') {
							goto l102
						}
						position++
					}
				l107:
				}
			l104:
				depth--
				add(ruleIdNondigit, position103)
			}
			return true
		l102:
			position, tokenIndex, depth = position102, tokenIndex102, depth102
			return false
		},
		/* 15 IdChar <- <([a-z] / [A-Z] / [0-9] / ('_' / '?' / '!'))> */
		func() bool {
			position110, tokenIndex110, depth110 := position, tokenIndex, depth
			{
				position111 := position
				depth++
				{
					position112, tokenIndex112, depth112 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l113
					}
					position++
					goto l112
				l113:
					position, tokenIndex, depth = position112, tokenIndex112, depth112
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l114
					}
					position++
					goto l112
				l114:
					position, tokenIndex, depth = position112, tokenIndex112, depth112
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l115
					}
					position++
					goto l112
				l115:
					position, tokenIndex, depth = position112, tokenIndex112, depth112
					{
						position116, tokenIndex116, depth116 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l117
						}
						position++
						goto l116
					l117:
						position, tokenIndex, depth = position116, tokenIndex116, depth116
						if buffer[position] != rune('?') {
							goto l118
						}
						position++
						goto l116
					l118:
						position, tokenIndex, depth = position116, tokenIndex116, depth116
						if buffer[position] != rune('!') {
							goto l110
						}
						position++
					}
				l116:
				}
			l112:
				depth--
				add(ruleIdChar, position111)
			}
			return true
		l110:
			position, tokenIndex, depth = position110, tokenIndex110, depth110
			return false
		},
		/* 16 StringLiteral <- <(Quote StringChar* Quote Spacing)> */
		func() bool {
			position119, tokenIndex119, depth119 := position, tokenIndex, depth
			{
				position120 := position
				depth++
				if !_rules[ruleQuote]() {
					goto l119
				}
			l121:
				{
					position122, tokenIndex122, depth122 := position, tokenIndex, depth
					if !_rules[ruleStringChar]() {
						goto l122
					}
					goto l121
				l122:
					position, tokenIndex, depth = position122, tokenIndex122, depth122
				}
				if !_rules[ruleQuote]() {
					goto l119
				}
				if !_rules[ruleSpacing]() {
					goto l119
				}
				depth--
				add(ruleStringLiteral, position120)
			}
			return true
		l119:
			position, tokenIndex, depth = position119, tokenIndex119, depth119
			return false
		},
		/* 17 Quote <- <'"'> */
		func() bool {
			position123, tokenIndex123, depth123 := position, tokenIndex, depth
			{
				position124 := position
				depth++
				if buffer[position] != rune('"') {
					goto l123
				}
				position++
				depth--
				add(ruleQuote, position124)
			}
			return true
		l123:
			position, tokenIndex, depth = position123, tokenIndex123, depth123
			return false
		},
		/* 18 StringChar <- <(Escape / (!('"' / '\n' / '\\') .))> */
		func() bool {
			position125, tokenIndex125, depth125 := position, tokenIndex, depth
			{
				position126 := position
				depth++
				{
					position127, tokenIndex127, depth127 := position, tokenIndex, depth
					if !_rules[ruleEscape]() {
						goto l128
					}
					goto l127
				l128:
					position, tokenIndex, depth = position127, tokenIndex127, depth127
					{
						position129, tokenIndex129, depth129 := position, tokenIndex, depth
						{
							position130, tokenIndex130, depth130 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l131
							}
							position++
							goto l130
						l131:
							position, tokenIndex, depth = position130, tokenIndex130, depth130
							if buffer[position] != rune('\n') {
								goto l132
							}
							position++
							goto l130
						l132:
							position, tokenIndex, depth = position130, tokenIndex130, depth130
							if buffer[position] != rune('\\') {
								goto l129
							}
							position++
						}
					l130:
						goto l125
					l129:
						position, tokenIndex, depth = position129, tokenIndex129, depth129
					}
					if !matchDot() {
						goto l125
					}
				}
			l127:
				depth--
				add(ruleStringChar, position126)
			}
			return true
		l125:
			position, tokenIndex, depth = position125, tokenIndex125, depth125
			return false
		},
		/* 19 LongStringLiteral <- <(LongString Spacing)> */
		func() bool {
			position133, tokenIndex133, depth133 := position, tokenIndex, depth
			{
				position134 := position
				depth++
				if !_rules[ruleLongString]() {
					goto l133
				}
				if !_rules[ruleSpacing]() {
					goto l133
				}
				depth--
				add(ruleLongStringLiteral, position134)
			}
			return true
		l133:
			position, tokenIndex, depth = position133, tokenIndex133, depth133
			return false
		},
		/* 20 LongString <- <('`' (('`' '`') / (!'`' .))* '`')> */
		func() bool {
			position135, tokenIndex135, depth135 := position, tokenIndex, depth
			{
				position136 := position
				depth++
				if buffer[position] != rune('`') {
					goto l135
				}
				position++
			l137:
				{
					position138, tokenIndex138, depth138 := position, tokenIndex, depth
					{
						position139, tokenIndex139, depth139 := position, tokenIndex, depth
						if buffer[position] != rune('`') {
							goto l140
						}
						position++
						if buffer[position] != rune('`') {
							goto l140
						}
						position++
						goto l139
					l140:
						position, tokenIndex, depth = position139, tokenIndex139, depth139
						{
							position141, tokenIndex141, depth141 := position, tokenIndex, depth
							if buffer[position] != rune('`') {
								goto l141
							}
							position++
							goto l138
						l141:
							position, tokenIndex, depth = position141, tokenIndex141, depth141
						}
						if !matchDot() {
							goto l138
						}
					}
				l139:
					goto l137
				l138:
					position, tokenIndex, depth = position138, tokenIndex138, depth138
				}
				if buffer[position] != rune('`') {
					goto l135
				}
				position++
				depth--
				add(ruleLongString, position136)
			}
			return true
		l135:
			position, tokenIndex, depth = position135, tokenIndex135, depth135
			return false
		},
		/* 21 Escape <- <('\\' (Block / .))> */
		func() bool {
			position142, tokenIndex142, depth142 := position, tokenIndex, depth
			{
				position143 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l142
				}
				position++
				{
					position144, tokenIndex144, depth144 := position, tokenIndex, depth
					if !_rules[ruleBlock]() {
						goto l145
					}
					goto l144
				l145:
					position, tokenIndex, depth = position144, tokenIndex144, depth144
					if !matchDot() {
						goto l142
					}
				}
			l144:
				depth--
				add(ruleEscape, position143)
			}
			return true
		l142:
			position, tokenIndex, depth = position142, tokenIndex142, depth142
			return false
		},
		/* 22 Number <- <(Numbers Spacing)> */
		func() bool {
			position146, tokenIndex146, depth146 := position, tokenIndex, depth
			{
				position147 := position
				depth++
				if !_rules[ruleNumbers]() {
					goto l146
				}
				if !_rules[ruleSpacing]() {
					goto l146
				}
				depth--
				add(ruleNumber, position147)
			}
			return true
		l146:
			position, tokenIndex, depth = position146, tokenIndex146, depth146
			return false
		},
		/* 23 Numbers <- <('-'? [0-9] [0-9]* ('.' [0-9] [0-9]*)?)> */
		func() bool {
			position148, tokenIndex148, depth148 := position, tokenIndex, depth
			{
				position149 := position
				depth++
				{
					position150, tokenIndex150, depth150 := position, tokenIndex, depth
					if buffer[position] != rune('-') {
						goto l150
					}
					position++
					goto l151
				l150:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
				}
			l151:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l148
				}
				position++
			l152:
				{
					position153, tokenIndex153, depth153 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l153
					}
					position++
					goto l152
				l153:
					position, tokenIndex, depth = position153, tokenIndex153, depth153
				}
				{
					position154, tokenIndex154, depth154 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l154
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l154
					}
					position++
				l156:
					{
						position157, tokenIndex157, depth157 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l157
						}
						position++
						goto l156
					l157:
						position, tokenIndex, depth = position157, tokenIndex157, depth157
					}
					goto l155
				l154:
					position, tokenIndex, depth = position154, tokenIndex154, depth154
				}
			l155:
				depth--
				add(ruleNumbers, position149)
			}
			return true
		l148:
			position, tokenIndex, depth = position148, tokenIndex148, depth148
			return false
		},
		/* 24 LPAR <- <('(' Spacing)> */
		func() bool {
			position158, tokenIndex158, depth158 := position, tokenIndex, depth
			{
				position159 := position
				depth++
				if buffer[position] != rune('(') {
					goto l158
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l158
				}
				depth--
				add(ruleLPAR, position159)
			}
			return true
		l158:
			position, tokenIndex, depth = position158, tokenIndex158, depth158
			return false
		},
		/* 25 RPAR <- <(')' Spacing)> */
		func() bool {
			position160, tokenIndex160, depth160 := position, tokenIndex, depth
			{
				position161 := position
				depth++
				if buffer[position] != rune(')') {
					goto l160
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l160
				}
				depth--
				add(ruleRPAR, position161)
			}
			return true
		l160:
			position, tokenIndex, depth = position160, tokenIndex160, depth160
			return false
		},
		/* 26 LCURLY <- <('{' Spacing)> */
		func() bool {
			position162, tokenIndex162, depth162 := position, tokenIndex, depth
			{
				position163 := position
				depth++
				if buffer[position] != rune('{') {
					goto l162
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l162
				}
				depth--
				add(ruleLCURLY, position163)
			}
			return true
		l162:
			position, tokenIndex, depth = position162, tokenIndex162, depth162
			return false
		},
		/* 27 RCURLY <- <('}' Spacing)> */
		func() bool {
			position164, tokenIndex164, depth164 := position, tokenIndex, depth
			{
				position165 := position
				depth++
				if buffer[position] != rune('}') {
					goto l164
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l164
				}
				depth--
				add(ruleRCURLY, position165)
			}
			return true
		l164:
			position, tokenIndex, depth = position164, tokenIndex164, depth164
			return false
		},
		/* 28 LBRACKET <- <('[' Spacing)> */
		func() bool {
			position166, tokenIndex166, depth166 := position, tokenIndex, depth
			{
				position167 := position
				depth++
				if buffer[position] != rune('[') {
					goto l166
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l166
				}
				depth--
				add(ruleLBRACKET, position167)
			}
			return true
		l166:
			position, tokenIndex, depth = position166, tokenIndex166, depth166
			return false
		},
		/* 29 RBRACKET <- <(']' Spacing)> */
		func() bool {
			position168, tokenIndex168, depth168 := position, tokenIndex, depth
			{
				position169 := position
				depth++
				if buffer[position] != rune(']') {
					goto l168
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l168
				}
				depth--
				add(ruleRBRACKET, position169)
			}
			return true
		l168:
			position, tokenIndex, depth = position168, tokenIndex168, depth168
			return false
		},
		/* 30 COMMA <- <(',' Spacing)> */
		func() bool {
			position170, tokenIndex170, depth170 := position, tokenIndex, depth
			{
				position171 := position
				depth++
				if buffer[position] != rune(',') {
					goto l170
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l170
				}
				depth--
				add(ruleCOMMA, position171)
			}
			return true
		l170:
			position, tokenIndex, depth = position170, tokenIndex170, depth170
			return false
		},
		/* 31 PCOMMA <- <(';' Spacing)> */
		func() bool {
			position172, tokenIndex172, depth172 := position, tokenIndex, depth
			{
				position173 := position
				depth++
				if buffer[position] != rune(';') {
					goto l172
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l172
				}
				depth--
				add(rulePCOMMA, position173)
			}
			return true
		l172:
			position, tokenIndex, depth = position172, tokenIndex172, depth172
			return false
		},
		/* 32 COLON <- <(':' Spacing)> */
		func() bool {
			position174, tokenIndex174, depth174 := position, tokenIndex, depth
			{
				position175 := position
				depth++
				if buffer[position] != rune(':') {
					goto l174
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l174
				}
				depth--
				add(ruleCOLON, position175)
			}
			return true
		l174:
			position, tokenIndex, depth = position174, tokenIndex174, depth174
			return false
		},
		/* 33 DOT <- <('.' Spacing)> */
		func() bool {
			position176, tokenIndex176, depth176 := position, tokenIndex, depth
			{
				position177 := position
				depth++
				if buffer[position] != rune('.') {
					goto l176
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l176
				}
				depth--
				add(ruleDOT, position177)
			}
			return true
		l176:
			position, tokenIndex, depth = position176, tokenIndex176, depth176
			return false
		},
		/* 34 PIPE <- <('|' Spacing)> */
		func() bool {
			position178, tokenIndex178, depth178 := position, tokenIndex, depth
			{
				position179 := position
				depth++
				if buffer[position] != rune('|') {
					goto l178
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l178
				}
				depth--
				add(rulePIPE, position179)
			}
			return true
		l178:
			position, tokenIndex, depth = position178, tokenIndex178, depth178
			return false
		},
		/* 35 DOLLAR <- <('$' Spacing)> */
		func() bool {
			position180, tokenIndex180, depth180 := position, tokenIndex, depth
			{
				position181 := position
				depth++
				if buffer[position] != rune('$') {
					goto l180
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l180
				}
				depth--
				add(ruleDOLLAR, position181)
			}
			return true
		l180:
			position, tokenIndex, depth = position180, tokenIndex180, depth180
			return false
		},
		/* 36 AMPERSAND <- <('&' Spacing)> */
		func() bool {
			position182, tokenIndex182, depth182 := position, tokenIndex, depth
			{
				position183 := position
				depth++
				if buffer[position] != rune('&') {
					goto l182
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l182
				}
				depth--
				add(ruleAMPERSAND, position183)
			}
			return true
		l182:
			position, tokenIndex, depth = position182, tokenIndex182, depth182
			return false
		},
		/* 37 EOT <- <!.> */
		func() bool {
			position184, tokenIndex184, depth184 := position, tokenIndex, depth
			{
				position185 := position
				depth++
				{
					position186, tokenIndex186, depth186 := position, tokenIndex, depth
					if !matchDot() {
						goto l186
					}
					goto l184
				l186:
					position, tokenIndex, depth = position186, tokenIndex186, depth186
				}
				depth--
				add(ruleEOT, position185)
			}
			return true
		l184:
			position, tokenIndex, depth = position184, tokenIndex184, depth184
			return false
		},
	}
	p.rules = _rules
}
