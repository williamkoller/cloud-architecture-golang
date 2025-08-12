#!/bin/bash

# Script para executar testes do projeto Go
# Autor: Sistema de testes automatizado
# Versão: 1.0

set -e  # Parar execução em caso de erro

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configurações
PROJECT_ROOT="$(pwd)"
COVERAGE_DIR="coverage"
COVERAGE_FILE="coverage.out"
COVERAGE_HTML="coverage.html"
BENCHMARK_FILE="benchmark.out"

# Função para imprimir cabeçalho
print_header() {
    echo -e "${CYAN}================================================${NC}"
    echo -e "${CYAN}    🧪 Sistema de Testes - Cloud Architecture${NC}"
    echo -e "${CYAN}================================================${NC}"
    echo
}

# Função para imprimir seção
print_section() {
    echo -e "\n${BLUE}📋 $1${NC}"
    echo -e "${BLUE}$(printf '%.0s-' {1..50})${NC}"
}

# Função para verificar dependências
check_dependencies() {
    print_section "Verificando Dependências"
    
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go não está instalado${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Go $(go version | cut -d' ' -f3) encontrado${NC}"
    
    # Verificar se o módulo Go está inicializado
    if [[ ! -f "go.mod" ]]; then
        echo -e "${RED}❌ go.mod não encontrado${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Módulo Go inicializado${NC}"
}

# Função para baixar dependências
download_dependencies() {
    print_section "Baixando Dependências"
    echo -e "${YELLOW}📦 Executando go mod download...${NC}"
    go mod download
    echo -e "${GREEN}✅ Dependências baixadas${NC}"
}

# Função para criar diretório de cobertura
setup_coverage_dir() {
    if [[ ! -d "$COVERAGE_DIR" ]]; then
        mkdir -p "$COVERAGE_DIR"
        echo -e "${GREEN}✅ Diretório de cobertura criado${NC}"
    fi
}

# Função para executar testes unitários
run_unit_tests() {
    print_section "Executando Testes Unitários"
    echo -e "${YELLOW}🧪 Executando todos os testes...${NC}"
    
    if go test ./... -v; then
        echo -e "${GREEN}✅ Todos os testes passaram${NC}"
    else
        echo -e "${RED}❌ Alguns testes falharam${NC}"
        exit 1
    fi
}

# Função para executar testes com cobertura
run_coverage_tests() {
    print_section "Executando Testes com Cobertura"
    setup_coverage_dir
    
    echo -e "${YELLOW}📊 Gerando relatório de cobertura...${NC}"
    
    # Executar testes com cobertura
    go test ./... -coverprofile="$COVERAGE_DIR/$COVERAGE_FILE" -covermode=atomic
    
    if [[ -f "$COVERAGE_DIR/$COVERAGE_FILE" ]]; then
        # Mostrar cobertura no terminal
        echo -e "\n${PURPLE}📈 Relatório de Cobertura:${NC}"
        go tool cover -func="$COVERAGE_DIR/$COVERAGE_FILE"
        
        # Gerar HTML
        go tool cover -html="$COVERAGE_DIR/$COVERAGE_FILE" -o="$COVERAGE_DIR/$COVERAGE_HTML"
        echo -e "${GREEN}✅ Relatório HTML gerado: $COVERAGE_DIR/$COVERAGE_HTML${NC}"
        
        # Mostrar cobertura total
        local total_coverage=$(go tool cover -func="$COVERAGE_DIR/$COVERAGE_FILE" | tail -1 | awk '{print $3}')
        echo -e "${CYAN}🎯 Cobertura Total: $total_coverage${NC}"
    else
        echo -e "${RED}❌ Arquivo de cobertura não foi gerado${NC}"
    fi
}

# Função para executar benchmarks
run_benchmarks() {
    print_section "Executando Benchmarks"
    
    echo -e "${YELLOW}⚡ Executando benchmarks...${NC}"
    
    # Executar benchmarks
    if go test ./... -bench=. -benchmem -run=^$ > "$COVERAGE_DIR/$BENCHMARK_FILE" 2>&1; then
        echo -e "${GREEN}✅ Benchmarks executados com sucesso${NC}"
        
        echo -e "\n${PURPLE}📈 Resultados dos Benchmarks:${NC}"
        cat "$COVERAGE_DIR/$BENCHMARK_FILE"
    else
        echo -e "${YELLOW}⚠️  Nenhum benchmark encontrado ou erro na execução${NC}"
    fi
}

# Função para executar testes específicos
run_specific_tests() {
    local package_path="$1"
    print_section "Executando Testes Específicos: $package_path"
    
    echo -e "${YELLOW}🎯 Executando testes para $package_path...${NC}"
    
    if go test "$package_path" -v -cover; then
        echo -e "${GREEN}✅ Testes específicos passaram${NC}"
    else
        echo -e "${RED}❌ Testes específicos falharam${NC}"
        exit 1
    fi
}

# Função para linting (se golangci-lint estiver disponível)
run_linting() {
    print_section "Executando Linting"
    
    if command -v golangci-lint &> /dev/null; then
        echo -e "${YELLOW}🔍 Executando golangci-lint...${NC}"
        if golangci-lint run; then
            echo -e "${GREEN}✅ Linting passou${NC}"
        else
            echo -e "${RED}❌ Problemas de linting encontrados${NC}"
            exit 1
        fi
    else
        echo -e "${YELLOW}⚠️  golangci-lint não encontrado, pulando linting${NC}"
    fi
}

# Função para executar race detector
run_race_tests() {
    print_section "Executando Testes com Race Detector"
    
    echo -e "${YELLOW}🏃 Executando testes com race detector...${NC}"
    
    if go test ./... -race; then
        echo -e "${GREEN}✅ Nenhuma race condition detectada${NC}"
    else
        echo -e "${RED}❌ Race conditions detectadas${NC}"
        exit 1
    fi
}

# Função para limpar arquivos temporários
cleanup() {
    print_section "Limpeza"
    
    echo -e "${YELLOW}🧹 Limpando arquivos temporários...${NC}"
    
    # Limpar cache de teste
    go clean -testcache
    
    echo -e "${GREEN}✅ Limpeza concluída${NC}"
}

# Função para mostrar estatísticas finais
show_stats() {
    print_section "Estatísticas Finais"
    
    if [[ -f "$COVERAGE_DIR/$COVERAGE_FILE" ]]; then
        local total_coverage=$(go tool cover -func="$COVERAGE_DIR/$COVERAGE_FILE" | tail -1 | awk '{print $3}')
        echo -e "${CYAN}📊 Cobertura Total: $total_coverage${NC}"
    fi
    
    local test_files=$(find . -name "*_test.go" -type f | wc -l)
    echo -e "${CYAN}📁 Arquivos de Teste: $test_files${NC}"
    
    local go_files=$(find . -name "*.go" -not -name "*_test.go" -type f | wc -l)
    echo -e "${CYAN}📄 Arquivos Go: $go_files${NC}"
    
    echo -e "\n${GREEN}🎉 Execução de testes concluída com sucesso!${NC}"
}

# Função para mostrar ajuda
show_help() {
    echo -e "${CYAN}Uso: $0 [OPÇÃO]${NC}"
    echo -e ""
    echo -e "${YELLOW}Opções:${NC}"
    echo -e "  ${GREEN}all${NC}          Executar todos os testes (padrão)"
    echo -e "  ${GREEN}unit${NC}         Executar apenas testes unitários"
    echo -e "  ${GREEN}coverage${NC}     Executar testes com cobertura"
    echo -e "  ${GREEN}bench${NC}        Executar apenas benchmarks"
    echo -e "  ${GREEN}race${NC}         Executar testes com race detector"
    echo -e "  ${GREEN}lint${NC}         Executar apenas linting"
    echo -e "  ${GREEN}clean${NC}        Limpar cache e arquivos temporários"
    echo -e "  ${GREEN}router${NC}       Executar testes específicos do router"
    echo -e "  ${GREEN}lint${NC}         Executar apenas linting"
    echo -e "  ${GREEN}race${NC}         Executar testes com race detector"
    echo -e "  ${GREEN}help${NC}         Mostrar esta ajuda"
    echo -e ""
    echo -e "${YELLOW}Exemplos:${NC}"
    echo -e "  $0                    # Executar todos os testes"
    echo -e "  $0 coverage           # Apenas cobertura"
    echo -e "  $0 router             # Apenas testes do router"
    echo -e ""
}

# Função principal
main() {
    local command="${1:-all}"
    
    print_header
    
    case "$command" in
        "all")
            check_dependencies
            download_dependencies
            run_unit_tests
            run_coverage_tests
            run_benchmarks
            run_race_tests
            run_linting
            show_stats
            ;;
        "unit")
            check_dependencies
            run_unit_tests
            ;;
        "coverage")
            check_dependencies
            run_coverage_tests
            ;;
        "bench")
            check_dependencies
            run_benchmarks
            ;;
        "race")
            check_dependencies
            run_race_tests
            ;;
        "lint")
            check_dependencies
            run_linting
            ;;
                 "clean")
             cleanup
             ;;
         "router")
             check_dependencies
             run_specific_tests "./internal/usr/router/"
             ;;
         "lint")
             check_dependencies
             run_linting
             ;;
         "race")
             check_dependencies
             run_race_tests
             ;;
         "help"|"-h"|"--help")
             show_help
             ;;
         *)
             echo -e "${RED}❌ Comando inválido: $command${NC}"
             echo -e "${YELLOW}Use '$0 help' para ver as opções disponíveis${NC}"
             exit 1
             ;;
     esac
}

# Executar função principal com todos os argumentos
main "$@" 