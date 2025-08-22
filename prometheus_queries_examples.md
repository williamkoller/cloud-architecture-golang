# 📊 Prometheus Queries - Exemplos Práticos

Este arquivo contém exemplos de queries PromQL que você pode usar para monitorar a aplicação **Cloud Architecture Golang**.

## 🎯 Como Usar

1. Acesse o Prometheus: http://localhost:9090
2. Vá para a aba **"Graph"**
3. Cole qualquer query abaixo no campo de busca
4. Clique **"Execute"**
5. Escolha **"Graph"** ou **"Table"** para visualizar

---

## 📈 1. Métricas HTTP Básicas

### Taxa de Requisições

```promql
# Total de requisições por segundo
rate(http_requests_total[1m])

# Requisições por minuto
rate(http_requests_total[1m]) * 60

# Por método e rota
sum(rate(http_requests_total[1m])) by (method, route)

# Por status code
sum(rate(http_requests_total[1m])) by (status)

# Top 5 rotas mais acessadas
topk(5, sum(rate(http_requests_total[5m])) by (route))
```

### Requisições em Andamento

```promql
# Número atual de requisições sendo processadas
http_in_flight_requests

# Histórico de requisições em andamento
http_in_flight_requests[10m]
```

---

## ⏱️ 2. Métricas de Latência

### Percentis de Latência

```promql
# P50 (mediana) global
histogram_quantile(0.50, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))

# P95 global
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))

# P99 global
histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))

# P95 por rota
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, route))

# P95 apenas da API de usuários
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{route="/api/users"}[5m])) by (le))
```

### Latência Média

```promql
# Latência média global
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])

# Latência média por rota
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m]) by (route)
```

### Top 5 Rotas Mais Lentas

```promql
# Rotas com maior P95
topk(5, histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[10m])) by (le, route)))
```

---

## ❌ 3. Métricas de Erro

### Taxa de Erro 5xx

```promql
# Porcentagem de erros 5xx
(sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))) * 100

# Taxa de erro por rota
(sum(rate(http_requests_total{status=~"5.."}[5m])) by (route) / sum(rate(http_requests_total[5m])) by (route)) * 100

# Contador de erros 500 específicos
sum(rate(http_requests_total{status="500"}[5m]))
```

### Taxa de Erro 4xx

```promql
# Porcentagem de erros 4xx
(sum(rate(http_requests_total{status=~"4.."}[5m])) / sum(rate(http_requests_total[5m]))) * 100

# Contador de 404s
sum(rate(http_requests_total{status="404"}[5m]))
```

### Distribuição por Status

```promql
# Requisições por classe (2xx, 4xx, 5xx)
sum(rate(http_requests_class_total[5m])) by (class)

# Porcentagem de sucessos (2xx)
(sum(rate(http_requests_total{status=~"2.."}[5m])) / sum(rate(http_requests_total[5m]))) * 100
```

### Top 5 Rotas com Mais Erro

```promql
# Rotas com mais erros 5xx
topk(5, sum(rate(http_requests_total{status=~"5.."}[10m])) by (route))
```

---

## 👥 4. Métricas de Domínio (Usuários)

### Operações de Usuários

```promql
# Taxa de criação de usuários por minuto
rate(users_created_total[5m]) * 60

# Taxa de atualização de usuários por minuto
rate(users_updated_total[5m]) * 60

# Taxa de exclusão de usuários por minuto
rate(users_deleted_total[5m]) * 60

# Total acumulado de usuários criados
users_created_total

# Total acumulado de usuários atualizados
users_updated_total

# Total acumulado de usuários deletados
users_deleted_total
```

### Atividade de Usuários

```promql
# Total de operações de usuários por minuto
(rate(users_created_total[5m]) + rate(users_updated_total[5m]) + rate(users_deleted_total[5m])) * 60

# Proporção de criações vs atualizações
rate(users_created_total[5m]) / (rate(users_created_total[5m]) + rate(users_updated_total[5m]))
```

---

## ⚙️ 5. Métricas de Sistema (Go Runtime)

### Goroutines

```promql
# Número atual de goroutines
go_goroutines

# Crescimento de goroutines
delta(go_goroutines[10m])

# Alerta se muitas goroutines (>1000)
go_goroutines > 1000
```

### Memória

```promql
# Heap alocado em MB
go_memstats_alloc_bytes / 1024 / 1024

# Heap total em MB
go_memstats_heap_alloc_bytes / 1024 / 1024

# Memória total do sistema em MB
go_memstats_sys_bytes / 1024 / 1024

# Taxa de crescimento da memória
rate(go_memstats_alloc_bytes[5m])

# Alerta se heap > 100MB
go_memstats_heap_alloc_bytes / 1024 / 1024 > 100
```

### Garbage Collector

```promql
# Duração do GC (P95)
histogram_quantile(0.95, rate(go_gc_duration_seconds_bucket[5m]))

# Frequência do GC (por minuto)
rate(go_gc_duration_seconds_count[5m]) * 60

# Tempo total gasto em GC por minuto
rate(go_gc_duration_seconds_sum[5m]) * 60
```

---

## 🚨 6. Métricas de Panic e Problemas

### Panics Recuperados

```promql
# Taxa de panics por minuto
rate(panic_recovered_total[5m]) * 60

# Total de panics
panic_recovered_total

# Aumento de panics nos últimos 5 minutos
increase(panic_recovered_total[5m])

# Alerta se houve panic recente
increase(panic_recovered_total[5m]) > 0
```

---

## 📏 7. Tamanho de Requisições e Respostas

### Tamanho das Respostas

```promql
# P95 do tamanho das respostas em bytes
histogram_quantile(0.95, sum(rate(http_response_size_bytes_bucket[5m])) by (le))

# Tamanho médio das respostas
rate(http_response_size_bytes_sum[5m]) / rate(http_response_size_bytes_count[5m])

# P95 do tamanho por rota
histogram_quantile(0.95, sum(rate(http_response_size_bytes_bucket[5m])) by (le, route))
```

### Tamanho das Requisições

```promql
# P95 do tamanho das requisições
histogram_quantile(0.95, sum(rate(http_request_size_bytes_bucket[5m])) by (le))

# Tamanho médio das requisições
rate(http_request_size_bytes_sum[5m]) / rate(http_request_size_bytes_count[5m])
```

---

## 🎯 8. Queries para SLIs (Service Level Indicators)

### Disponibilidade

```promql
# Disponibilidade (% de requisições não-5xx)
(1 - (sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m])))) * 100

# Uptime do serviço
up{job="usr-api-host"} * 100
```

### Performance

```promql
# SLI: 95% das requisições < 300ms
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le)) < 0.3

# Porcentagem de requisições rápidas (<100ms)
(sum(rate(http_request_duration_seconds_bucket{le="0.1"}[5m])) / sum(rate(http_request_duration_seconds_count[5m]))) * 100
```

### Throughput

```promql
# QPS atual
sum(rate(http_requests_total[1m]))

# QPS por hora
sum(rate(http_requests_total[1h])) * 3600

# Peak QPS nos últimos 24h
max_over_time(sum(rate(http_requests_total[1m]))[24h:1m])
```

---

## 🔄 9. Monitoramento da Infraestrutura

### Status dos Targets

```promql
# Targets que estão UP
up

# Targets que estão DOWN
up == 0

# Número de targets UP por job
sum(up) by (job)

# Último tempo de scrape bem-sucedido
time() - up
```

### Informações da Aplicação

```promql
# Info da aplicação
app_info

# Verificar versão da aplicação
app_info{version!=""}
```

---

## 📊 10. Queries para Business Intelligence

### Análise de Padrões

```promql
# Horário de pico (requisições por hora)
sum(rate(http_requests_total[1h])) * 3600

# Comparação com ontem (mesmo horário)
sum(rate(http_requests_total[5m])) / sum(rate(http_requests_total[5m] offset 24h))

# Tendência de criação de usuários (últimas 4 horas)
rate(users_created_total[4h])
```

### Top N Análises

```promql
# Top 10 rotas por volume
topk(10, sum(rate(http_requests_total[1h])) by (route))

# Top 5 métodos HTTP mais usados
topk(5, sum(rate(http_requests_total[1h])) by (method))
```

---

## 🚨 11. Queries de Alerta (Usar nos Alertas)

### Alertas Críticos

```promql
# Taxa de erro alta (>5%)
(sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))) > 0.05

# Latência alta (P95 > 300ms)
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le)) > 0.3

# Serviço down
up{job="usr-api-host"} == 0

# Panic detectado
increase(panic_recovered_total[5m]) > 0
```

### Alertas de Warning

```promql
# Taxa de erro moderada (>1%)
(sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))) > 0.01

# Muitas goroutines (>500)
go_goroutines > 500

# Alto uso de memória (>50MB)
go_memstats_heap_alloc_bytes / 1024 / 1024 > 50

# QPS muito baixo (<0.1/s)
sum(rate(http_requests_total[5m])) < 0.1
```

---

## 🔍 12. Queries de Debugging

### Investigação de Problemas

```promql
# Todas as métricas de uma requisição específica
{route="/api/users", method="POST"}

# Requisições que falharam nos últimos 10 minutos
increase(http_requests_total{status=~"5.."}[10m])

# Latência máxima registrada
max(http_request_duration_seconds)

# Variação de latência (desvio padrão)
stddev(rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m]))
```

### Health Checks

```promql
# Verificar se aplicação está respondendo
up{job="usr-api-host"}

# Última requisição bem-sucedida
time() - timestamp(http_requests_total{status="200"})

# Verificar se métricas estão sendo atualizadas
time() - timestamp(http_requests_total)
```

---

## 📝 Exemplos de Uso Prático

### Cenário 1: Investigar Lentidão

```promql
# 1. Verificar latência atual
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))

# 2. Identificar rota mais lenta
topk(3, histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, route)))

# 3. Verificar se há erros relacionados
sum(rate(http_requests_total{status=~"5.."}[5m])) by (route)
```

### Cenário 2: Investigar Aumento de Erros

```promql
# 1. Taxa de erro atual
(sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))) * 100

# 2. Rota com mais erros
topk(3, sum(rate(http_requests_total{status=~"5.."}[5m])) by (route))

# 3. Verificar se houve panics
increase(panic_recovered_total[15m])
```

### Cenário 3: Análise de Capacidade

```promql
# 1. QPS atual vs capacidade
sum(rate(http_requests_total[5m]))

# 2. Utilização de recursos
go_memstats_heap_alloc_bytes / 1024 / 1024

# 3. Requisições em fila
http_in_flight_requests
```

---

## 🎯 Dicas de Uso

1. **Intervalos de Tempo**: Use `[5m]` para alertas, `[1h]` para análises
2. **Rate vs Increase**: Use `rate()` para taxa por segundo, `increase()` para total no período
3. **Percentis**: P95 é melhor que média para SLAs
4. **Labels**: Sempre use labels para filtrar (ex: `{route="/api/users"}`)
5. **Agregação**: Use `sum()`, `avg()`, `max()` para agregar métricas

---

## 📚 Recursos Adicionais

- **Prometheus Documentation**: https://prometheus.io/docs/prometheus/latest/querying/basics/
- **PromQL Cheat Sheet**: https://promlabs.com/promql-cheat-sheet/
- **Grafana Dashboard**: http://localhost:3000

---

_Arquivo gerado automaticamente para o projeto Cloud Architecture Golang_
