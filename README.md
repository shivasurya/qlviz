# QLViz

A Go tool that generates inheritance diagrams for CodeQL classes by analyzing `.ql` and `.qll` files. The output is a visual representation of class hierarchies in GraphViz DOT format.

## Features

- Scans CodeQL source files (`.ql` and `.qll`) recursively
- Extracts class definitions and inheritance relationships
- Generates DOT format output
- Supports abstract classes and inheritance chains
- Creates visual class hierarchy diagrams

## Requirements

- Go 1.11+
- GraphViz (for rendering)