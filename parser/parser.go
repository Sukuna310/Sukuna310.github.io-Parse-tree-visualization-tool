package parser

import (
	"fmt"
	"strings"
)

// Parser implements a recursive descent parser for arithmetic expressions
type Parser struct {
	grammar    *Grammar
	tokens     []Token
	pos        int
	nodeID     int
	steps      []Step
	recordSteps bool
}

// NewParser creates a new parser with the given grammar
func NewParser(grammar *Grammar) *Parser {
	return &Parser{
		grammar:     grammar,
		tokens:      []Token{},
		pos:         0,
		nodeID:      0,
		steps:       []Step{},
		recordSteps: false,
	}
}

// Parse parses the input string and returns the parse result
func (p *Parser) Parse(input string, recordSteps bool) *ParseResult {
	// Tokenize input
	lexer := NewLexer(input)
	p.tokens = lexer.Tokenize()
	p.pos = 0
	p.nodeID = 0
	p.steps = []Step{}
	p.recordSteps = recordSteps

	result := &ParseResult{
		Success: true,
		Tokens:  p.tokens,
		Steps:   []Step{},
	}

	// Start parsing from the start symbol
	if p.grammar.StartSymbol == "" {
		result.Success = false
		result.Error = "Grammar has no start symbol"
		return result
	}

	tree, err := p.parseNonTerminal(p.grammar.StartSymbol, -1)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result
	}

	// Check if all tokens were consumed
	if p.current().Type != TOKEN_EOF {
		result.Success = false
		result.Error = fmt.Sprintf("Unexpected token '%s' at position %d", 
			p.current().Value, p.current().Position)
		return result
	}

	result.Tree = tree
	result.Steps = p.steps
	return result
}

// parseNonTerminal parses a non-terminal symbol
func (p *Parser) parseNonTerminal(symbol string, parentID int) (*TreeNode, error) {
	prod, ok := p.grammar.Productions[symbol]
	if !ok {
		return nil, fmt.Errorf("undefined non-terminal: %s", symbol)
	}

	// Create node for this non-terminal
	node := p.createNode(symbol, false, parentID)

	// Try each alternative
	for _, alt := range prod.Body {
		// Save current position for backtracking
		savedPos := p.pos
		savedStepsLen := len(p.steps)

		children, err := p.parseAlternative(alt, node.ID)
		if err == nil {
			node.Children = children
			return node, nil
		}

		// Backtrack
		p.pos = savedPos
		p.steps = p.steps[:savedStepsLen]
	}

	return nil, fmt.Errorf("no matching production for %s at position %d (found '%s')", 
		symbol, p.pos, p.current().Value)
}

// parseAlternative parses a single production alternative
func (p *Parser) parseAlternative(symbols []string, parentID int) ([]*TreeNode, error) {
	children := []*TreeNode{}

	for _, symbol := range symbols {
		// Handle epsilon
		if symbol == "ε" || symbol == "epsilon" {
			epsilonNode := p.createNode("ε", true, parentID)
			children = append(children, epsilonNode)
			continue
		}

		// Check if it's a terminal
		if p.isTerminal(symbol) {
			termNode, err := p.matchTerminal(symbol, parentID)
			if err != nil {
				return nil, err
			}
			children = append(children, termNode)
		} else {
			// It's a non-terminal
			subtree, err := p.parseNonTerminal(symbol, parentID)
			if err != nil {
				return nil, err
			}
			children = append(children, subtree)
		}
	}

	return children, nil
}

// matchTerminal tries to match the current token with a terminal symbol
func (p *Parser) matchTerminal(symbol string, parentID int) (*TreeNode, error) {
	token := p.current()

	matched := false
	switch symbol {
	case "number":
		matched = token.Type == TOKEN_NUMBER
	case "+":
		matched = token.Type == TOKEN_PLUS
	case "-":
		matched = token.Type == TOKEN_MINUS
	case "*":
		matched = token.Type == TOKEN_MULT
	case "/":
		matched = token.Type == TOKEN_DIV
	case "(":
		matched = token.Type == TOKEN_LPAREN
	case ")":
		matched = token.Type == TOKEN_RPAREN
	default:
		// Try exact value match for identifiers
		matched = token.Value == symbol
	}

	if !matched {
		return nil, fmt.Errorf("expected '%s', got '%s' at position %d", 
			symbol, token.Value, token.Position)
	}

	// Create terminal node with actual value
	displayValue := token.Value
	if symbol == "number" {
		displayValue = token.Value
	}

	node := p.createNode(displayValue, true, parentID)
	p.advance()
	return node, nil
}

// isTerminal checks if a symbol is a terminal
func (p *Parser) isTerminal(symbol string) bool {
	if p.grammar.Terminals[symbol] {
		return true
	}
	// Built-in terminals
	terminals := map[string]bool{
		"number": true, "+": true, "-": true, "*": true, "/": true,
		"(": true, ")": true, "ε": true, "epsilon": true,
	}
	return terminals[symbol]
}

// current returns the current token
func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TOKEN_EOF, Value: "", Position: p.pos}
	}
	return p.tokens[p.pos]
}

// advance moves to the next token
func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

// createNode creates a new tree node and records a step if needed
func (p *Parser) createNode(label string, isTerminal bool, parentID int) *TreeNode {
	p.nodeID++
	node := &TreeNode{
		ID:         p.nodeID,
		Label:      label,
		Children:   []*TreeNode{},
		IsTerminal: isTerminal,
	}

	if p.recordSteps {
		step := Step{
			Action:      "add",
			Description: p.buildStepDescription(label, isTerminal),
			NodeID:      node.ID,
			ParentID:    parentID,
		}
		p.steps = append(p.steps, step)
	}

	return node
}

// buildStepDescription creates a human-readable description for a step
func (p *Parser) buildStepDescription(label string, isTerminal bool) string {
	if isTerminal {
		if label == "ε" {
			return "Match epsilon (empty string)"
		}
		return fmt.Sprintf("Match terminal '%s'", label)
	}
	return fmt.Sprintf("Expand non-terminal <%s>", label)
}

// ParseWithDefaultGrammar parses using the default arithmetic expression grammar
func ParseWithDefaultGrammar(input string, recordSteps bool) *ParseResult {
	grammarText := GetDefaultArithmeticGrammar()
	grammar, err := ParseGrammar(grammarText)
	if err != nil {
		return &ParseResult{
			Success: false,
			Error:   "Failed to parse default grammar: " + err.Error(),
		}
	}

	parser := NewParser(grammar)
	return parser.Parse(input, recordSteps)
}

// TreeToString returns a string representation of the parse tree
func TreeToString(node *TreeNode, indent string) string {
	if node == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(indent)
	
	if node.IsTerminal {
		sb.WriteString(fmt.Sprintf("'%s'\n", node.Label))
	} else {
		sb.WriteString(fmt.Sprintf("<%s>\n", node.Label))
	}

	for _, child := range node.Children {
		sb.WriteString(TreeToString(child, indent+"  "))
	}

	return sb.String()
}
