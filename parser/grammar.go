package parser

import (
	"regexp"
	"strings"
)

// ParseGrammar parses a BNF-style grammar definition
// Format: NonTerminal -> production1 | production2 | ...
// Example:
//   E  -> T E'
//   E' -> + T E' | ε
//   T  -> F T'
//   T' -> * F T' | ε
//   F  -> ( E ) | number
func ParseGrammar(input string) (*Grammar, error) {
	grammar := &Grammar{
		Productions:  make(map[string]*Production),
		Terminals:    make(map[string]bool),
		NonTerminals: make(map[string]bool),
	}

	lines := strings.Split(input, "\n")
	isFirst := true

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			continue
		}

		// Split by -> or →
		parts := regexp.MustCompile(`\s*->|→\s*`).Split(line, 2)
		if len(parts) != 2 {
			continue
		}

		head := strings.TrimSpace(parts[0])
		if head == "" {
			continue
		}

		// Mark as non-terminal
		grammar.NonTerminals[head] = true

		// Set start symbol (first production)
		if isFirst {
			grammar.StartSymbol = head
			isFirst = false
		}

		// Split alternatives by |
		alternatives := strings.Split(parts[1], "|")
		bodies := [][]string{}

		for _, alt := range alternatives {
			alt = strings.TrimSpace(alt)
			if alt == "" {
				continue
			}

			// Parse symbols in this alternative
			symbols := parseSymbols(alt)
			if len(symbols) > 0 {
				bodies = append(bodies, symbols)
			}
		}

		if existing, ok := grammar.Productions[head]; ok {
			existing.Body = append(existing.Body, bodies...)
		} else {
			grammar.Productions[head] = &Production{
				Head: head,
				Body: bodies,
			}
		}
	}

	// Identify terminals (symbols that are not non-terminals)
	for _, prod := range grammar.Productions {
		for _, alt := range prod.Body {
			for _, symbol := range alt {
				if !grammar.NonTerminals[symbol] && symbol != "ε" && symbol != "epsilon" {
					grammar.Terminals[symbol] = true
				}
			}
		}
	}

	return grammar, nil
}

// parseSymbols parses a production body into individual symbols
func parseSymbols(body string) []string {
	symbols := []string{}
	
	// Handle special tokens
	body = strings.ReplaceAll(body, "ε", " ε ")
	body = strings.ReplaceAll(body, "epsilon", " ε ")
	
	// Split by whitespace
	parts := strings.Fields(body)
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			symbols = append(symbols, part)
		}
	}
	
	return symbols
}

// GetDefaultArithmeticGrammar returns the default grammar for arithmetic expressions
// This is an LL(1) compatible grammar (left-recursion removed)
func GetDefaultArithmeticGrammar() string {
	return `E  -> T E'
E' -> + T E' | - T E' | ε
T  -> F T'
T' -> * F T' | / F T' | ε
F  -> ( E ) | number`
}

// ValidateGrammar checks if the grammar is valid for LL(1) parsing
func ValidateGrammar(grammar *Grammar) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
	}

	if grammar.StartSymbol == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "No start symbol defined")
		return result
	}

	if _, ok := grammar.Productions[grammar.StartSymbol]; !ok {
		result.Valid = false
		result.Errors = append(result.Errors, "Start symbol has no productions")
		return result
	}

	// Check for undefined non-terminals
	for _, prod := range grammar.Productions {
		for _, alt := range prod.Body {
			for _, symbol := range alt {
				if grammar.NonTerminals[symbol] {
					if _, ok := grammar.Productions[symbol]; !ok {
						result.Valid = false
						result.Errors = append(result.Errors, "Undefined non-terminal: "+symbol)
					}
				}
			}
		}
	}

	// Warning for potential left recursion (simple check)
	for head, prod := range grammar.Productions {
		for _, alt := range prod.Body {
			if len(alt) > 0 && alt[0] == head {
				result.Warnings = append(result.Warnings, 
					"Potential left recursion in production: "+head+" -> "+strings.Join(alt, " "))
			}
		}
	}

	return result
}

// IsTerminal checks if a symbol is a terminal
func (g *Grammar) IsTerminal(symbol string) bool {
	return g.Terminals[symbol] || isOperatorOrLiteral(symbol)
}

// IsNonTerminal checks if a symbol is a non-terminal
func (g *Grammar) IsNonTerminal(symbol string) bool {
	return g.NonTerminals[symbol]
}

// isOperatorOrLiteral checks if a symbol is an operator or literal
func isOperatorOrLiteral(symbol string) bool {
	operators := map[string]bool{
		"+": true, "-": true, "*": true, "/": true,
		"(": true, ")": true,
		"number": true, "ε": true, "epsilon": true,
	}
	return operators[symbol]
}
