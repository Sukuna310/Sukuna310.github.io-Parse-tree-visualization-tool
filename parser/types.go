package parser

// TokenType represents the type of a token
type TokenType string

const (
	TOKEN_NUMBER  TokenType = "NUMBER"
	TOKEN_PLUS    TokenType = "PLUS"
	TOKEN_MINUS   TokenType = "MINUS"
	TOKEN_MULT    TokenType = "MULT"
	TOKEN_DIV     TokenType = "DIV"
	TOKEN_LPAREN  TokenType = "LPAREN"
	TOKEN_RPAREN  TokenType = "RPAREN"
	TOKEN_EOF     TokenType = "EOF"
	TOKEN_EPSILON TokenType = "EPSILON"
	TOKEN_IDENT   TokenType = "IDENT"
	TOKEN_UNKNOWN TokenType = "UNKNOWN"
)

// Token represents a lexical token
type Token struct {
	Type     TokenType `json:"type"`
	Value    string    `json:"value"`
	Position int       `json:"position"`
}

// TreeNode represents a node in the parse tree
type TreeNode struct {
	ID         int         `json:"id"`
	Label      string      `json:"label"`
	Children   []*TreeNode `json:"children"`
	IsTerminal bool        `json:"isTerminal"`
	Value      string      `json:"value,omitempty"`
}

// Step represents a single step in the parsing process for animation
type Step struct {
	Action      string    `json:"action"`
	Description string    `json:"description"`
	NodeID      int       `json:"nodeId"`
	ParentID    int       `json:"parentId,omitempty"`
	Tree        *TreeNode `json:"tree"`
}

// ParseResult represents the result of parsing
type ParseResult struct {
	Success bool       `json:"success"`
	Tree    *TreeNode  `json:"tree"`
	Steps   []Step     `json:"steps"`
	Error   string     `json:"error,omitempty"`
	Tokens  []Token    `json:"tokens"`
}

// Production represents a single grammar production
type Production struct {
	Head string     `json:"head"`
	Body [][]string `json:"body"` // Each alternative is a slice of symbols
}

// Grammar represents a context-free grammar
type Grammar struct {
	Productions map[string]*Production `json:"productions"`
	StartSymbol string                 `json:"startSymbol"`
	Terminals   map[string]bool        `json:"terminals"`
	NonTerminals map[string]bool       `json:"nonTerminals"`
}

// ValidationResult represents the result of grammar validation
type ValidationResult struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}
