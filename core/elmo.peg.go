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
	ruleBlockWithoutSpacing
	ruleList
	ruleSpacing
	ruleWhiteSpace
	ruleLongComment
	ruleLineComment
	ruleNewLine
	ruleIdentifier
	ruleIdNondigit
	ruleIdChar
	ruleIdEnd
	ruleStringLiteral
	ruleQuote
	ruleStringChar
	ruleEscape
	ruleLongStringLiteral
	ruleBackTick
	ruleLongStringChar
	ruleLongEscape
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
	"BlockWithoutSpacing",
	"List",
	"Spacing",
	"WhiteSpace",
	"LongComment",
	"LineComment",
	"NewLine",
	"Identifier",
	"IdNondigit",
	"IdChar",
	"IdEnd",
	"StringLiteral",
	"Quote",
	"StringChar",
	"Escape",
	"LongStringLiteral",
	"BackTick",
	"LongStringChar",
	"LongEscape",
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
	rules  [43]func() bool
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
		/* 7 BlockWithoutSpacing <- <(LCURLY NewLine* Line* '}')> */
		func() bool {
			position53, tokenIndex53, depth53 := position, tokenIndex, depth
			{
				position54 := position
				depth++
				if !_rules[ruleLCURLY]() {
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
			l57:
				{
					position58, tokenIndex58, depth58 := position, tokenIndex, depth
					if !_rules[ruleLine]() {
						goto l58
					}
					goto l57
				l58:
					position, tokenIndex, depth = position58, tokenIndex58, depth58
				}
				if buffer[position] != rune('}') {
					goto l53
				}
				position++
				depth--
				add(ruleBlockWithoutSpacing, position54)
			}
			return true
		l53:
			position, tokenIndex, depth = position53, tokenIndex53, depth53
			return false
		},
		/* 8 List <- <(LBRACKET NewLine* (Argument / NewLine)? (((COMMA NewLine?)? Argument) / NewLine)* RBRACKET)> */
		func() bool {
			position59, tokenIndex59, depth59 := position, tokenIndex, depth
			{
				position60 := position
				depth++
				if !_rules[ruleLBRACKET]() {
					goto l59
				}
			l61:
				{
					position62, tokenIndex62, depth62 := position, tokenIndex, depth
					if !_rules[ruleNewLine]() {
						goto l62
					}
					goto l61
				l62:
					position, tokenIndex, depth = position62, tokenIndex62, depth62
				}
				{
					position63, tokenIndex63, depth63 := position, tokenIndex, depth
					{
						position65, tokenIndex65, depth65 := position, tokenIndex, depth
						if !_rules[ruleArgument]() {
							goto l66
						}
						goto l65
					l66:
						position, tokenIndex, depth = position65, tokenIndex65, depth65
						if !_rules[ruleNewLine]() {
							goto l63
						}
					}
				l65:
					goto l64
				l63:
					position, tokenIndex, depth = position63, tokenIndex63, depth63
				}
			l64:
			l67:
				{
					position68, tokenIndex68, depth68 := position, tokenIndex, depth
					{
						position69, tokenIndex69, depth69 := position, tokenIndex, depth
						{
							position71, tokenIndex71, depth71 := position, tokenIndex, depth
							if !_rules[ruleCOMMA]() {
								goto l71
							}
							{
								position73, tokenIndex73, depth73 := position, tokenIndex, depth
								if !_rules[ruleNewLine]() {
									goto l73
								}
								goto l74
							l73:
								position, tokenIndex, depth = position73, tokenIndex73, depth73
							}
						l74:
							goto l72
						l71:
							position, tokenIndex, depth = position71, tokenIndex71, depth71
						}
					l72:
						if !_rules[ruleArgument]() {
							goto l70
						}
						goto l69
					l70:
						position, tokenIndex, depth = position69, tokenIndex69, depth69
						if !_rules[ruleNewLine]() {
							goto l68
						}
					}
				l69:
					goto l67
				l68:
					position, tokenIndex, depth = position68, tokenIndex68, depth68
				}
				if !_rules[ruleRBRACKET]() {
					goto l59
				}
				depth--
				add(ruleList, position60)
			}
			return true
		l59:
			position, tokenIndex, depth = position59, tokenIndex59, depth59
			return false
		},
		/* 9 Spacing <- <(WhiteSpace / LongComment / LineComment)*> */
		func() bool {
			{
				position76 := position
				depth++
			l77:
				{
					position78, tokenIndex78, depth78 := position, tokenIndex, depth
					{
						position79, tokenIndex79, depth79 := position, tokenIndex, depth
						if !_rules[ruleWhiteSpace]() {
							goto l80
						}
						goto l79
					l80:
						position, tokenIndex, depth = position79, tokenIndex79, depth79
						if !_rules[ruleLongComment]() {
							goto l81
						}
						goto l79
					l81:
						position, tokenIndex, depth = position79, tokenIndex79, depth79
						if !_rules[ruleLineComment]() {
							goto l78
						}
					}
				l79:
					goto l77
				l78:
					position, tokenIndex, depth = position78, tokenIndex78, depth78
				}
				depth--
				add(ruleSpacing, position76)
			}
			return true
		},
		/* 10 WhiteSpace <- <(' ' / '\t')> */
		func() bool {
			position82, tokenIndex82, depth82 := position, tokenIndex, depth
			{
				position83 := position
				depth++
				{
					position84, tokenIndex84, depth84 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l85
					}
					position++
					goto l84
				l85:
					position, tokenIndex, depth = position84, tokenIndex84, depth84
					if buffer[position] != rune('\t') {
						goto l82
					}
					position++
				}
			l84:
				depth--
				add(ruleWhiteSpace, position83)
			}
			return true
		l82:
			position, tokenIndex, depth = position82, tokenIndex82, depth82
			return false
		},
		/* 11 LongComment <- <('/' '*' (!('*' '/') .)* ('*' '/'))> */
		func() bool {
			position86, tokenIndex86, depth86 := position, tokenIndex, depth
			{
				position87 := position
				depth++
				if buffer[position] != rune('/') {
					goto l86
				}
				position++
				if buffer[position] != rune('*') {
					goto l86
				}
				position++
			l88:
				{
					position89, tokenIndex89, depth89 := position, tokenIndex, depth
					{
						position90, tokenIndex90, depth90 := position, tokenIndex, depth
						if buffer[position] != rune('*') {
							goto l90
						}
						position++
						if buffer[position] != rune('/') {
							goto l90
						}
						position++
						goto l89
					l90:
						position, tokenIndex, depth = position90, tokenIndex90, depth90
					}
					if !matchDot() {
						goto l89
					}
					goto l88
				l89:
					position, tokenIndex, depth = position89, tokenIndex89, depth89
				}
				if buffer[position] != rune('*') {
					goto l86
				}
				position++
				if buffer[position] != rune('/') {
					goto l86
				}
				position++
				depth--
				add(ruleLongComment, position87)
			}
			return true
		l86:
			position, tokenIndex, depth = position86, tokenIndex86, depth86
			return false
		},
		/* 12 LineComment <- <('#' (!'\n' .)*)> */
		func() bool {
			position91, tokenIndex91, depth91 := position, tokenIndex, depth
			{
				position92 := position
				depth++
				if buffer[position] != rune('#') {
					goto l91
				}
				position++
			l93:
				{
					position94, tokenIndex94, depth94 := position, tokenIndex, depth
					{
						position95, tokenIndex95, depth95 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l95
						}
						position++
						goto l94
					l95:
						position, tokenIndex, depth = position95, tokenIndex95, depth95
					}
					if !matchDot() {
						goto l94
					}
					goto l93
				l94:
					position, tokenIndex, depth = position94, tokenIndex94, depth94
				}
				depth--
				add(ruleLineComment, position92)
			}
			return true
		l91:
			position, tokenIndex, depth = position91, tokenIndex91, depth91
			return false
		},
		/* 13 NewLine <- <(('\n' / '\r') Spacing)+> */
		func() bool {
			position96, tokenIndex96, depth96 := position, tokenIndex, depth
			{
				position97 := position
				depth++
				{
					position100, tokenIndex100, depth100 := position, tokenIndex, depth
					if buffer[position] != rune('\n') {
						goto l101
					}
					position++
					goto l100
				l101:
					position, tokenIndex, depth = position100, tokenIndex100, depth100
					if buffer[position] != rune('\r') {
						goto l96
					}
					position++
				}
			l100:
				if !_rules[ruleSpacing]() {
					goto l96
				}
			l98:
				{
					position99, tokenIndex99, depth99 := position, tokenIndex, depth
					{
						position102, tokenIndex102, depth102 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l103
						}
						position++
						goto l102
					l103:
						position, tokenIndex, depth = position102, tokenIndex102, depth102
						if buffer[position] != rune('\r') {
							goto l99
						}
						position++
					}
				l102:
					if !_rules[ruleSpacing]() {
						goto l99
					}
					goto l98
				l99:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
				}
				depth--
				add(ruleNewLine, position97)
			}
			return true
		l96:
			position, tokenIndex, depth = position96, tokenIndex96, depth96
			return false
		},
		/* 14 Identifier <- <((IdNondigit IdChar* ((IdEnd? Spacing) / (IdEnd Spacing?))) / (IdEnd Spacing))> */
		func() bool {
			position104, tokenIndex104, depth104 := position, tokenIndex, depth
			{
				position105 := position
				depth++
				{
					position106, tokenIndex106, depth106 := position, tokenIndex, depth
					if !_rules[ruleIdNondigit]() {
						goto l107
					}
				l108:
					{
						position109, tokenIndex109, depth109 := position, tokenIndex, depth
						if !_rules[ruleIdChar]() {
							goto l109
						}
						goto l108
					l109:
						position, tokenIndex, depth = position109, tokenIndex109, depth109
					}
					{
						position110, tokenIndex110, depth110 := position, tokenIndex, depth
						{
							position112, tokenIndex112, depth112 := position, tokenIndex, depth
							if !_rules[ruleIdEnd]() {
								goto l112
							}
							goto l113
						l112:
							position, tokenIndex, depth = position112, tokenIndex112, depth112
						}
					l113:
						if !_rules[ruleSpacing]() {
							goto l111
						}
						goto l110
					l111:
						position, tokenIndex, depth = position110, tokenIndex110, depth110
						if !_rules[ruleIdEnd]() {
							goto l107
						}
						{
							position114, tokenIndex114, depth114 := position, tokenIndex, depth
							if !_rules[ruleSpacing]() {
								goto l114
							}
							goto l115
						l114:
							position, tokenIndex, depth = position114, tokenIndex114, depth114
						}
					l115:
					}
				l110:
					goto l106
				l107:
					position, tokenIndex, depth = position106, tokenIndex106, depth106
					if !_rules[ruleIdEnd]() {
						goto l104
					}
					if !_rules[ruleSpacing]() {
						goto l104
					}
				}
			l106:
				depth--
				add(ruleIdentifier, position105)
			}
			return true
		l104:
			position, tokenIndex, depth = position104, tokenIndex104, depth104
			return false
		},
		/* 15 IdNondigit <- <([a-z] / [A-Z] / '_')> */
		func() bool {
			position116, tokenIndex116, depth116 := position, tokenIndex, depth
			{
				position117 := position
				depth++
				{
					position118, tokenIndex118, depth118 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l119
					}
					position++
					goto l118
				l119:
					position, tokenIndex, depth = position118, tokenIndex118, depth118
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l120
					}
					position++
					goto l118
				l120:
					position, tokenIndex, depth = position118, tokenIndex118, depth118
					if buffer[position] != rune('_') {
						goto l116
					}
					position++
				}
			l118:
				depth--
				add(ruleIdNondigit, position117)
			}
			return true
		l116:
			position, tokenIndex, depth = position116, tokenIndex116, depth116
			return false
		},
		/* 16 IdChar <- <([a-z] / [A-Z] / [0-9] / '_')> */
		func() bool {
			position121, tokenIndex121, depth121 := position, tokenIndex, depth
			{
				position122 := position
				depth++
				{
					position123, tokenIndex123, depth123 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l124
					}
					position++
					goto l123
				l124:
					position, tokenIndex, depth = position123, tokenIndex123, depth123
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l125
					}
					position++
					goto l123
				l125:
					position, tokenIndex, depth = position123, tokenIndex123, depth123
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l126
					}
					position++
					goto l123
				l126:
					position, tokenIndex, depth = position123, tokenIndex123, depth123
					if buffer[position] != rune('_') {
						goto l121
					}
					position++
				}
			l123:
				depth--
				add(ruleIdChar, position122)
			}
			return true
		l121:
			position, tokenIndex, depth = position121, tokenIndex121, depth121
			return false
		},
		/* 17 IdEnd <- <('?' / '!')> */
		func() bool {
			position127, tokenIndex127, depth127 := position, tokenIndex, depth
			{
				position128 := position
				depth++
				{
					position129, tokenIndex129, depth129 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l130
					}
					position++
					goto l129
				l130:
					position, tokenIndex, depth = position129, tokenIndex129, depth129
					if buffer[position] != rune('!') {
						goto l127
					}
					position++
				}
			l129:
				depth--
				add(ruleIdEnd, position128)
			}
			return true
		l127:
			position, tokenIndex, depth = position127, tokenIndex127, depth127
			return false
		},
		/* 18 StringLiteral <- <(Quote StringChar* Quote Spacing)> */
		func() bool {
			position131, tokenIndex131, depth131 := position, tokenIndex, depth
			{
				position132 := position
				depth++
				if !_rules[ruleQuote]() {
					goto l131
				}
			l133:
				{
					position134, tokenIndex134, depth134 := position, tokenIndex, depth
					if !_rules[ruleStringChar]() {
						goto l134
					}
					goto l133
				l134:
					position, tokenIndex, depth = position134, tokenIndex134, depth134
				}
				if !_rules[ruleQuote]() {
					goto l131
				}
				if !_rules[ruleSpacing]() {
					goto l131
				}
				depth--
				add(ruleStringLiteral, position132)
			}
			return true
		l131:
			position, tokenIndex, depth = position131, tokenIndex131, depth131
			return false
		},
		/* 19 Quote <- <'"'> */
		func() bool {
			position135, tokenIndex135, depth135 := position, tokenIndex, depth
			{
				position136 := position
				depth++
				if buffer[position] != rune('"') {
					goto l135
				}
				position++
				depth--
				add(ruleQuote, position136)
			}
			return true
		l135:
			position, tokenIndex, depth = position135, tokenIndex135, depth135
			return false
		},
		/* 20 StringChar <- <(Escape / (!('"' / '\n' / '\\') .))> */
		func() bool {
			position137, tokenIndex137, depth137 := position, tokenIndex, depth
			{
				position138 := position
				depth++
				{
					position139, tokenIndex139, depth139 := position, tokenIndex, depth
					if !_rules[ruleEscape]() {
						goto l140
					}
					goto l139
				l140:
					position, tokenIndex, depth = position139, tokenIndex139, depth139
					{
						position141, tokenIndex141, depth141 := position, tokenIndex, depth
						{
							position142, tokenIndex142, depth142 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l143
							}
							position++
							goto l142
						l143:
							position, tokenIndex, depth = position142, tokenIndex142, depth142
							if buffer[position] != rune('\n') {
								goto l144
							}
							position++
							goto l142
						l144:
							position, tokenIndex, depth = position142, tokenIndex142, depth142
							if buffer[position] != rune('\\') {
								goto l141
							}
							position++
						}
					l142:
						goto l137
					l141:
						position, tokenIndex, depth = position141, tokenIndex141, depth141
					}
					if !matchDot() {
						goto l137
					}
				}
			l139:
				depth--
				add(ruleStringChar, position138)
			}
			return true
		l137:
			position, tokenIndex, depth = position137, tokenIndex137, depth137
			return false
		},
		/* 21 Escape <- <('\\' (BlockWithoutSpacing / .))> */
		func() bool {
			position145, tokenIndex145, depth145 := position, tokenIndex, depth
			{
				position146 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l145
				}
				position++
				{
					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					if !_rules[ruleBlockWithoutSpacing]() {
						goto l148
					}
					goto l147
				l148:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
					if !matchDot() {
						goto l145
					}
				}
			l147:
				depth--
				add(ruleEscape, position146)
			}
			return true
		l145:
			position, tokenIndex, depth = position145, tokenIndex145, depth145
			return false
		},
		/* 22 LongStringLiteral <- <(BackTick LongStringChar* BackTick Spacing)> */
		func() bool {
			position149, tokenIndex149, depth149 := position, tokenIndex, depth
			{
				position150 := position
				depth++
				if !_rules[ruleBackTick]() {
					goto l149
				}
			l151:
				{
					position152, tokenIndex152, depth152 := position, tokenIndex, depth
					if !_rules[ruleLongStringChar]() {
						goto l152
					}
					goto l151
				l152:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
				}
				if !_rules[ruleBackTick]() {
					goto l149
				}
				if !_rules[ruleSpacing]() {
					goto l149
				}
				depth--
				add(ruleLongStringLiteral, position150)
			}
			return true
		l149:
			position, tokenIndex, depth = position149, tokenIndex149, depth149
			return false
		},
		/* 23 BackTick <- <'`'> */
		func() bool {
			position153, tokenIndex153, depth153 := position, tokenIndex, depth
			{
				position154 := position
				depth++
				if buffer[position] != rune('`') {
					goto l153
				}
				position++
				depth--
				add(ruleBackTick, position154)
			}
			return true
		l153:
			position, tokenIndex, depth = position153, tokenIndex153, depth153
			return false
		},
		/* 24 LongStringChar <- <(LongEscape / (!'`' .))> */
		func() bool {
			position155, tokenIndex155, depth155 := position, tokenIndex, depth
			{
				position156 := position
				depth++
				{
					position157, tokenIndex157, depth157 := position, tokenIndex, depth
					if !_rules[ruleLongEscape]() {
						goto l158
					}
					goto l157
				l158:
					position, tokenIndex, depth = position157, tokenIndex157, depth157
					{
						position159, tokenIndex159, depth159 := position, tokenIndex, depth
						if buffer[position] != rune('`') {
							goto l159
						}
						position++
						goto l155
					l159:
						position, tokenIndex, depth = position159, tokenIndex159, depth159
					}
					if !matchDot() {
						goto l155
					}
				}
			l157:
				depth--
				add(ruleLongStringChar, position156)
			}
			return true
		l155:
			position, tokenIndex, depth = position155, tokenIndex155, depth155
			return false
		},
		/* 25 LongEscape <- <('`' (BlockWithoutSpacing / '`'))> */
		func() bool {
			position160, tokenIndex160, depth160 := position, tokenIndex, depth
			{
				position161 := position
				depth++
				if buffer[position] != rune('`') {
					goto l160
				}
				position++
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if !_rules[ruleBlockWithoutSpacing]() {
						goto l163
					}
					goto l162
				l163:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
					if buffer[position] != rune('`') {
						goto l160
					}
					position++
				}
			l162:
				depth--
				add(ruleLongEscape, position161)
			}
			return true
		l160:
			position, tokenIndex, depth = position160, tokenIndex160, depth160
			return false
		},
		/* 26 Number <- <(Numbers Spacing)> */
		func() bool {
			position164, tokenIndex164, depth164 := position, tokenIndex, depth
			{
				position165 := position
				depth++
				if !_rules[ruleNumbers]() {
					goto l164
				}
				if !_rules[ruleSpacing]() {
					goto l164
				}
				depth--
				add(ruleNumber, position165)
			}
			return true
		l164:
			position, tokenIndex, depth = position164, tokenIndex164, depth164
			return false
		},
		/* 27 Numbers <- <('-'? [0-9] [0-9]* ('.' [0-9] [0-9]*)?)> */
		func() bool {
			position166, tokenIndex166, depth166 := position, tokenIndex, depth
			{
				position167 := position
				depth++
				{
					position168, tokenIndex168, depth168 := position, tokenIndex, depth
					if buffer[position] != rune('-') {
						goto l168
					}
					position++
					goto l169
				l168:
					position, tokenIndex, depth = position168, tokenIndex168, depth168
				}
			l169:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l166
				}
				position++
			l170:
				{
					position171, tokenIndex171, depth171 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l171
					}
					position++
					goto l170
				l171:
					position, tokenIndex, depth = position171, tokenIndex171, depth171
				}
				{
					position172, tokenIndex172, depth172 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l172
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l172
					}
					position++
				l174:
					{
						position175, tokenIndex175, depth175 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l175
						}
						position++
						goto l174
					l175:
						position, tokenIndex, depth = position175, tokenIndex175, depth175
					}
					goto l173
				l172:
					position, tokenIndex, depth = position172, tokenIndex172, depth172
				}
			l173:
				depth--
				add(ruleNumbers, position167)
			}
			return true
		l166:
			position, tokenIndex, depth = position166, tokenIndex166, depth166
			return false
		},
		/* 28 LPAR <- <('(' Spacing)> */
		func() bool {
			position176, tokenIndex176, depth176 := position, tokenIndex, depth
			{
				position177 := position
				depth++
				if buffer[position] != rune('(') {
					goto l176
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l176
				}
				depth--
				add(ruleLPAR, position177)
			}
			return true
		l176:
			position, tokenIndex, depth = position176, tokenIndex176, depth176
			return false
		},
		/* 29 RPAR <- <(')' Spacing)> */
		func() bool {
			position178, tokenIndex178, depth178 := position, tokenIndex, depth
			{
				position179 := position
				depth++
				if buffer[position] != rune(')') {
					goto l178
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l178
				}
				depth--
				add(ruleRPAR, position179)
			}
			return true
		l178:
			position, tokenIndex, depth = position178, tokenIndex178, depth178
			return false
		},
		/* 30 LCURLY <- <('{' Spacing)> */
		func() bool {
			position180, tokenIndex180, depth180 := position, tokenIndex, depth
			{
				position181 := position
				depth++
				if buffer[position] != rune('{') {
					goto l180
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l180
				}
				depth--
				add(ruleLCURLY, position181)
			}
			return true
		l180:
			position, tokenIndex, depth = position180, tokenIndex180, depth180
			return false
		},
		/* 31 RCURLY <- <('}' Spacing)> */
		func() bool {
			position182, tokenIndex182, depth182 := position, tokenIndex, depth
			{
				position183 := position
				depth++
				if buffer[position] != rune('}') {
					goto l182
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l182
				}
				depth--
				add(ruleRCURLY, position183)
			}
			return true
		l182:
			position, tokenIndex, depth = position182, tokenIndex182, depth182
			return false
		},
		/* 32 LBRACKET <- <('[' Spacing)> */
		func() bool {
			position184, tokenIndex184, depth184 := position, tokenIndex, depth
			{
				position185 := position
				depth++
				if buffer[position] != rune('[') {
					goto l184
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l184
				}
				depth--
				add(ruleLBRACKET, position185)
			}
			return true
		l184:
			position, tokenIndex, depth = position184, tokenIndex184, depth184
			return false
		},
		/* 33 RBRACKET <- <(']' Spacing)> */
		func() bool {
			position186, tokenIndex186, depth186 := position, tokenIndex, depth
			{
				position187 := position
				depth++
				if buffer[position] != rune(']') {
					goto l186
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l186
				}
				depth--
				add(ruleRBRACKET, position187)
			}
			return true
		l186:
			position, tokenIndex, depth = position186, tokenIndex186, depth186
			return false
		},
		/* 34 COMMA <- <(',' Spacing)> */
		func() bool {
			position188, tokenIndex188, depth188 := position, tokenIndex, depth
			{
				position189 := position
				depth++
				if buffer[position] != rune(',') {
					goto l188
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l188
				}
				depth--
				add(ruleCOMMA, position189)
			}
			return true
		l188:
			position, tokenIndex, depth = position188, tokenIndex188, depth188
			return false
		},
		/* 35 PCOMMA <- <(';' Spacing)> */
		func() bool {
			position190, tokenIndex190, depth190 := position, tokenIndex, depth
			{
				position191 := position
				depth++
				if buffer[position] != rune(';') {
					goto l190
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l190
				}
				depth--
				add(rulePCOMMA, position191)
			}
			return true
		l190:
			position, tokenIndex, depth = position190, tokenIndex190, depth190
			return false
		},
		/* 36 COLON <- <(':' Spacing)> */
		func() bool {
			position192, tokenIndex192, depth192 := position, tokenIndex, depth
			{
				position193 := position
				depth++
				if buffer[position] != rune(':') {
					goto l192
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l192
				}
				depth--
				add(ruleCOLON, position193)
			}
			return true
		l192:
			position, tokenIndex, depth = position192, tokenIndex192, depth192
			return false
		},
		/* 37 DOT <- <('.' Spacing)> */
		func() bool {
			position194, tokenIndex194, depth194 := position, tokenIndex, depth
			{
				position195 := position
				depth++
				if buffer[position] != rune('.') {
					goto l194
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l194
				}
				depth--
				add(ruleDOT, position195)
			}
			return true
		l194:
			position, tokenIndex, depth = position194, tokenIndex194, depth194
			return false
		},
		/* 38 PIPE <- <('|' Spacing)> */
		func() bool {
			position196, tokenIndex196, depth196 := position, tokenIndex, depth
			{
				position197 := position
				depth++
				if buffer[position] != rune('|') {
					goto l196
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l196
				}
				depth--
				add(rulePIPE, position197)
			}
			return true
		l196:
			position, tokenIndex, depth = position196, tokenIndex196, depth196
			return false
		},
		/* 39 DOLLAR <- <('$' Spacing)> */
		func() bool {
			position198, tokenIndex198, depth198 := position, tokenIndex, depth
			{
				position199 := position
				depth++
				if buffer[position] != rune('$') {
					goto l198
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l198
				}
				depth--
				add(ruleDOLLAR, position199)
			}
			return true
		l198:
			position, tokenIndex, depth = position198, tokenIndex198, depth198
			return false
		},
		/* 40 AMPERSAND <- <('&' Spacing)> */
		func() bool {
			position200, tokenIndex200, depth200 := position, tokenIndex, depth
			{
				position201 := position
				depth++
				if buffer[position] != rune('&') {
					goto l200
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l200
				}
				depth--
				add(ruleAMPERSAND, position201)
			}
			return true
		l200:
			position, tokenIndex, depth = position200, tokenIndex200, depth200
			return false
		},
		/* 41 EOT <- <!.> */
		func() bool {
			position202, tokenIndex202, depth202 := position, tokenIndex, depth
			{
				position203 := position
				depth++
				{
					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					if !matchDot() {
						goto l204
					}
					goto l202
				l204:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
				}
				depth--
				add(ruleEOT, position203)
			}
			return true
		l202:
			position, tokenIndex, depth = position202, tokenIndex202, depth202
			return false
		},
	}
	p.rules = _rules
}
