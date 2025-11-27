#!/bin/bash

# Извлекает версии из Taskfile.yml и устанавливает их как переменные окружения

# Извлекаем GO_VERSION
GO_VERSION=$(grep -E "^\s*GO_VERSION:" Taskfile.yml | head -1 | sed "s/.*GO_VERSION:\s*['\"]*//" | sed "s/['\"].*//")
echo "GO_VERSION=$GO_VERSION"
echo "GO_VERSION=$GO_VERSION" >> $GITHUB_OUTPUT

# Извлекаем GOLANGCI_LINT_VERSION
GOLANGCI_LINT_VERSION=$(grep -E "^\s*GOLANGCI_LINT_VERSION:" Taskfile.yml | head -1 | sed "s/.*GOLANGCI_LINT_VERSION:\s*['\"]*//" | sed "s/['\"].*//")
echo "GOLANGCI_LINT_VERSION=$GOLANGCI_LINT_VERSION"
echo "GOLANGCI_LINT_VERSION=$GOLANGCI_LINT_VERSION" >> $GITHUB_OUTPUT

# Извлекаем MODULES
MODULES=$(grep -E "^\s*MODULES:" Taskfile.yml | head -1 | sed "s/.*MODULES:\s*//" | tr -d '\r')
echo "MODULES=$MODULES"
echo "MODULES=$MODULES" >> $GITHUB_OUTPUT
