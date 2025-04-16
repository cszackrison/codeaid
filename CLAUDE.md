# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands
- Build: `go build`
- Run: `go run main.go`
- Test: `go test ./...`
- Run single test: `go test -v -run TestName`
- Lint: `go vet ./...`
- Format: `gofmt -w .`

## Code Style
- Imports: Group standard library imports first, then third-party packages
- Formatting: Follow Go standard formatting with gofmt
- Types: Use explicit types, prefer interfaces for flexibility
- Naming: Use camelCase for variables, PascalCase for exported identifiers
- Error handling: Check errors explicitly, avoid silent failures
- Project structure: Maintain a clean separation of concerns
- Comments: Use godoc style comments for exported identifiers