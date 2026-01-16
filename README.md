# Parse Tree Visualizer

A GUI application that parses strings using context-free grammars and visualizes the parse tree interactively.

Built with **Go (Wails)** backend and **D3.js** frontend for a university System Programming course.

## Features

- 📝 **Grammar Input** - Define BNF-style context-free grammars
- ✅ **Real-time Validation** - Grammar syntax checking with error messages
- 🔤 **Token Display** - Color-coded tokenization preview
- 🌳 **Instant Parse Tree** - Full tree visualization on demand
- 🎬 **Step-by-Step Mode** - Animated tree construction for demos
- ⚡ **Configurable Speed** - Animation speed from 100ms to 2000ms
- 🔍 **Interactive Zoom/Pan** - Navigate large trees easily

## Prerequisites

- Go 1.21+
- Node.js 15+ with npm
- Xcode Command Line Tools (macOS)

## Quick Start

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Run in development mode
wails dev

# Build production app
wails build
```

## Default Grammar

The app includes a default LL(1) arithmetic expression grammar:

```
E  -> T E'
E' -> + T E' | - T E' | ε
T  -> F T'
T' -> * F T' | / F T' | ε
F  -> ( E ) | number
```

## Demo

1. Enter or load the default grammar
2. Type an expression like `3 + 5 * 2`
3. Click **Parse** for instant visualization
4. Click **Step Mode** for animated tree building
5. Use **Auto** to auto-play the animation

## Project Structure

```
├── parser/          # Go parser package
│   ├── types.go     # Data structures
│   ├── lexer.go     # Tokenization
│   ├── grammar.go   # Grammar parsing
│   └── parser.go    # Recursive descent parser
├── app.go           # Wails bindings
├── main.go          # App entry point
└── frontend/        # Web frontend
    ├── index.html
    └── src/
        ├── main.js  # D3.js visualization
        └── style.css
```

## License

MIT
