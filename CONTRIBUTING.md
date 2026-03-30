# Contributing to Cue

Thank you for your interest in contributing to **Cue**! This document outlines the process for contributing to the project and helps maintain a high standard of code quality and community collaboration.

## Getting Started

1. **Fork the repository** on GitHub.
2. **Clone your fork locally**:
   ```bash
   git clone https://github.com/YOUR-USERNAME/cue.git
   ```
3. **Set the upstream remote**:
   ```bash
   git remote add upstream https://github.com/GyaneshSamanta/cue.git
   ```

## Development Environment Setup

Cue is built primarily in Go, and we strive to keep development simple.

1. **Ensure Go 1.22+ is installed.**
2. **Run `go mod tidy`** to fetch dependencies.
3. **Compile the binary**:
   ```bash
   go build -o cue.exe
   ```
4. **Test out your local binary**:
   ```bash
   ./cue.exe
   ```

## Adding Features

If you are adding new Stores or Macros to the engine:
1. Navigate to either `internal/macro/builtins/*.go` for new Macros or `internal/store/stacks/*.go` for new specialized environment stores.
2. Be sure to use our custom structured error package `ui.StructuredError` for all panics/crashes to provide robust terminal guidance. 
3. Verify your changes do not break existing test cases and the project compiles.

## Submitting a Pull Request

1. **Create a feature branch:** `git checkout -b feature/my-new-feature`
2. **Commit your changes:** Write a clear, concise commit message.
3. **Push to your fork:** `git push origin feature/my-new-feature`
4. **Open a PR:** Go to the main repository and open a Pull Request. Provide a detailed explanation of why the change is necessary and what it achieves.

## Issues and Bug Reports

If you discover a bug or have a feature request:
1. Search the existing issues to ensure it hasn't already been reported.
2. Create a new issue describing the bug, including clear reproduction steps and the expected behavior.

We appreciate all contributions to make **Cue** the ultimate context-aware developer assistant!
