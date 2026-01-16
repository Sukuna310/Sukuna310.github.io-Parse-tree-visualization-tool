package parser

import (
	"strings"
	"unicode"
)

// Lexer tokenizes input strings for arithmetic expressions
type Lexer struct {
	input   string
	pos     int
	tokens  []Token
}

// NewLexer creates a new lexer for the given input
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		tokens: []Token{},
	}
}

// Tokenize converts the input string into a slice of tokens
func (l *Lexer) Tokenize() []Token {
	l.tokens = []Token{}
	l.pos = 0

	for l.pos < len(l.input) {
		ch := l.input[l.pos]

		// Skip whitespace
		if unicode.IsSpace(rune(ch)) {
			l.pos++
			continue
		}

		// Number (integer or decimal)
		if unicode.IsDigit(rune(ch)) {
			l.tokens = append(l.tokens, l.readNumber())
			continue
		}

		// Identifier (for grammar terminals like 'number', 'id')
		if unicode.IsLetter(rune(ch)) {
			l.tokens = append(l.tokens, l.readIdentifier())
			continue
		}

		// Single character tokens
		token := Token{Position: l.pos}
		switch ch {
		case '+':
			token.Type = TOKEN_PLUS
			token.Value = "+"
		case '-':
			token.Type = TOKEN_MINUS
			token.Value = "-"
		case '*':
			token.Type = TOKEN_MULT
			token.Value = "*"
		case '/':
			token.Type = TOKEN_DIV
			token.Value = "/"
		case '(':
			token.Type = TOKEN_LPAREN
			token.Value = "("
		case ')':
			token.Type = TOKEN_RPAREN
			token.Value = ")"
		default:
			token.Type = TOKEN_UNKNOWN
			token.Value = string(ch)
		}
		l.tokens = append(l.tokens, token)
		l.pos++
	}

	// Add EOF token
	l.tokens = append(l.tokens, Token{
		Type:     TOKEN_EOF,
		Value:    "",
		Position: l.pos,
	})

	return l.tokens
}

// readNumber reads a number (integer or decimal) from the input
func (l *Lexer) readNumber() Token {
	start := l.pos
	hasDecimal := false

	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if unicode.IsDigit(rune(ch)) {
			l.pos++
		} else if ch == '.' && !hasDecimal {
			hasDecimal = true
			l.pos++
		} else {
			break
		}
	}

	return Token{
		Type:     TOKEN_NUMBER,
		Value:    l.input[start:l.pos],
		Position: start,
	}
}

// readIdentifier reads an identifier from the input
func (l *Lexer) readIdentifier() Token {
	start := l.pos

	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_' || ch == '\'' {
			l.pos++
		} else {
			break
		}
	}

	value := l.input[start:l.pos]
	
	// Check for special keywords that map to tokens
	tokenType := TOKEN_IDENT
	switch strings.ToLower(value) {
	case "number":
		tokenType = TOKEN_NUMBER
	}

	return Token{
		Type:     tokenType,
		Value:    value,
		Position: start,
	}
}

// GetTokenTypeName returns a human-readable name for a token type
func GetTokenTypeName(t TokenType) string {
	switch t {
	case TOKEN_NUMBER:
		return "number"
	case TOKEN_PLUS:
		return "+"
	case TOKEN_MINUS:
		return "-"
	case TOKEN_MULT:
		return "*"
	case TOKEN_DIV:
		return "/"
	case TOKEN_LPAREN:
		return "("
	case TOKEN_RPAREN:
		return ")"
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_EPSILON:
		return "ε"
	default:
		return string(t)
	}
}
