package main

import (
	"context"

	"parse-tree-viz/parser"
)

// App struct holds the application state
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ParseString parses an input string using the provided grammar
// and returns the complete parse tree
func (a *App) ParseString(grammarText string, input string) *parser.ParseResult {
	// Parse the grammar
	grammar, err := parser.ParseGrammar(grammarText)
	if err != nil {
		return &parser.ParseResult{
			Success: false,
			Error:   "Failed to parse grammar: " + err.Error(),
		}
	}

	// Validate the grammar
	validation := parser.ValidateGrammar(grammar)
	if !validation.Valid {
		errMsg := "Grammar validation failed: "
		for i, e := range validation.Errors {
			if i > 0 {
				errMsg += "; "
			}
			errMsg += e
		}
		return &parser.ParseResult{
			Success: false,
			Error:   errMsg,
		}
	}

	// Parse the input
	p := parser.NewParser(grammar)
	return p.Parse(input, false)
}

// ParseStepByStep parses an input string and returns steps for animation
func (a *App) ParseStepByStep(grammarText string, input string) *parser.ParseResult {
	// Parse the grammar
	grammar, err := parser.ParseGrammar(grammarText)
	if err != nil {
		return &parser.ParseResult{
			Success: false,
			Error:   "Failed to parse grammar: " + err.Error(),
		}
	}

	// Validate the grammar
	validation := parser.ValidateGrammar(grammar)
	if !validation.Valid {
		errMsg := "Grammar validation failed: "
		for i, e := range validation.Errors {
			if i > 0 {
				errMsg += "; "
			}
			errMsg += e
		}
		return &parser.ParseResult{
			Success: false,
			Error:   errMsg,
		}
	}

	// Parse the input with step recording
	p := parser.NewParser(grammar)
	return p.Parse(input, true)
}

// ValidateGrammar validates a grammar definition
func (a *App) ValidateGrammar(grammarText string) *parser.ValidationResult {
	grammar, err := parser.ParseGrammar(grammarText)
	if err != nil {
		return &parser.ValidationResult{
			Valid:  false,
			Errors: []string{"Failed to parse grammar: " + err.Error()},
		}
	}

	return parser.ValidateGrammar(grammar)
}

// GetDefaultGrammar returns the default arithmetic expression grammar
func (a *App) GetDefaultGrammar() string {
	return parser.GetDefaultArithmeticGrammar()
}

// GetTokens tokenizes an input string and returns the tokens
func (a *App) GetTokens(input string) []parser.Token {
	lexer := parser.NewLexer(input)
	return lexer.Tokenize()
}
