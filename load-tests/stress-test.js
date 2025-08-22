import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

// Métricas customizadas
const errorRate = new Rate('errors');
const apiLatency = new Trend('api_latency');
const usersCreated = new Counter('users_created');
const usersDeleted = new Counter('users_deleted');

// Configuração dos cenários de teste
export const options = {
  scenarios: {
    // Cenário 1: Smoke Test - Verificação básica
    smoke_test: {
      executor: 'constant-vus',
      vus: 1,
      duration: '30s',
      tags: { test_type: 'smoke' },
    },

    // Cenário 2: Load Test - Carga normal
    load_test: {
      executor: 'constant-vus',
      vus: 10,
      duration: '2m',
      startTime: '30s',
      tags: { test_type: 'load' },
    },

    // Cenário 3: Stress Test - Carga alta
    stress_test: {
      executor: 'ramping-vus',
      startVUs: 10,
      stages: [
        { duration: '1m', target: 50 },
        { duration: '2m', target: 100 },
        { duration: '1m', target: 150 },
        { duration: '2m', target: 200 },
        { duration: '1m', target: 0 },
      ],
      startTime: '2m30s',
      tags: { test_type: 'stress' },
    },

    // Cenário 4: Spike Test - Picos de tráfego
    spike_test: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '30s', target: 20 },
        { duration: '30s', target: 300 }, // Spike
        { duration: '1m', target: 20 },
        { duration: '30s', target: 0 },
      ],
      startTime: '9m30s',
      tags: { test_type: 'spike' },
    },
  },

  thresholds: {
    // SLIs ajustados para ambiente de desenvolvimento
    http_req_duration: ['p(95)<1000'], // P95 < 1s
    http_req_failed: ['rate<0.2'], // Taxa de erro < 20% (inclui erros intencionais)
    errors: ['rate<0.1'], // Erros não-intencionais < 10%
    api_latency: ['p(95)<1500'], // API latency < 1.5s
  },
};

// Base URL da aplicação
const BASE_URL = 'http://localhost:8080';

// Dados de teste
const testUsers = [
  {
    name: 'João Silva',
    email: 'joao@test.com',
    password: '123456',
    userType: 'Admin',
  },
  {
    name: 'Maria Santos',
    email: 'maria@test.com',
    password: '123456',
    userType: 'User',
  },
  {
    name: 'Pedro Costa',
    email: 'pedro@test.com',
    password: '123456',
    userType: 'Admin',
  },
  {
    name: 'Ana Oliveira',
    email: 'ana@test.com',
    password: '123456',
    userType: 'User',
  },
  {
    name: 'Carlos Lima',
    email: 'carlos@test.com',
    password: '123456',
    userType: 'Manager',
  },
];

// Função para gerar email único
function generateUniqueEmail(baseEmail) {
  const timestamp = new Date().getTime();
  const random = Math.floor(Math.random() * 1000);
  return baseEmail.replace('@', `_${timestamp}_${random}@`);
}

// Função principal de teste
export default function () {
  const testScenario = Math.random();

  if (testScenario < 0.3) {
    // 30% - Cenário de leitura (GET requests)
    testReadOperations();
  } else if (testScenario < 0.6) {
    // 30% - Cenário de escrita (POST/PUT requests)
    testWriteOperations();
  } else if (testScenario < 0.85) {
    // 25% - Cenário misto (CRUD completo)
    testCrudOperations();
  } else {
    // 15% - Cenário de erro (força erros 404/500)
    testErrorScenarios();
  }

  // Pausa aleatória entre 0.5s e 2s
  sleep(Math.random() * 1.5 + 0.5);
}

// Cenário 1: Operações de Leitura
function testReadOperations() {
  const responses = [];

  // 1. Health Check
  let response = http.get(`${BASE_URL}/health`, {
    tags: { operation: 'health_check' },
  });
  responses.push(response);

  // 2. Listar usuários
  response = http.get(`${BASE_URL}/api/users`, {
    tags: { operation: 'list_users' },
  });
  responses.push(response);

  // 3. Buscar usuário específico (pode dar 404)
  const randomEmail = 'user' + Math.floor(Math.random() * 100) + '@test.com';
  response = http.get(`${BASE_URL}/api/users/${randomEmail}`, {
    tags: { operation: 'get_user' },
  });
  responses.push(response);

  // Validações
  responses.forEach((res, index) => {
    const isHealthCheck = index === 0;
    const expectedStatus = isHealthCheck ? 200 : [200, 404];

    check(res, {
      [`Read operation ${index + 1} status ok`]: (r) =>
        Array.isArray(expectedStatus)
          ? expectedStatus.includes(r.status)
          : r.status === expectedStatus,
      [`Read operation ${index + 1} response time < 1s`]: (r) =>
        r.timings.duration < 1000,
    });

    errorRate.add(res.status >= 400 && res.status !== 404);
    apiLatency.add(res.timings.duration);
  });
}

// Cenário 2: Operações de Escrita
function testWriteOperations() {
  const user = testUsers[Math.floor(Math.random() * testUsers.length)];
  user.email = generateUniqueEmail(user.email);

  // 1. Criar usuário
  let response = http.post(`${BASE_URL}/api/users`, JSON.stringify(user), {
    headers: { 'Content-Type': 'application/json' },
    tags: { operation: 'create_user' },
  });

  const createSuccess = check(response, {
    'User creation status is 201': (r) => r.status === 201,
    'User creation response time < 2s': (r) => r.timings.duration < 2000,
    'User creation returns user data': (r) => {
      try {
        const userData = JSON.parse(r.body);
        return userData.email === user.email;
      } catch {
        return false;
      }
    },
  });

  if (createSuccess) {
    usersCreated.add(1);

    // 2. Atualizar usuário criado
    const updateData = { name: user.name + ' Updated' };
    response = http.patch(
      `${BASE_URL}/api/users/${user.email}`,
      JSON.stringify(updateData),
      {
        headers: { 'Content-Type': 'application/json' },
        tags: { operation: 'update_user' },
      }
    );

    check(response, {
      'User update status is 200': (r) => r.status === 200,
      'User update response time < 2s': (r) => r.timings.duration < 2000,
    });
  }

  errorRate.add(response.status >= 400);
  apiLatency.add(response.timings.duration);
}

// Cenário 3: CRUD Completo
function testCrudOperations() {
  const user = testUsers[Math.floor(Math.random() * testUsers.length)];
  user.email = generateUniqueEmail(user.email);

  // 1. Criar
  let response = http.post(`${BASE_URL}/api/users`, JSON.stringify(user), {
    headers: { 'Content-Type': 'application/json' },
    tags: { operation: 'crud_create' },
  });

  if (response.status === 201) {
    usersCreated.add(1);

    // 2. Ler
    response = http.get(`${BASE_URL}/api/users/${user.email}`, {
      tags: { operation: 'crud_read' },
    });

    check(response, {
      'CRUD read status is 200': (r) => r.status === 200,
    });

    // 3. Atualizar
    const updateData = {
      name: user.name + ' CRUD Updated',
      active: !user.active,
    };
    response = http.patch(
      `${BASE_URL}/api/users/${user.email}`,
      JSON.stringify(updateData),
      {
        headers: { 'Content-Type': 'application/json' },
        tags: { operation: 'crud_update' },
      }
    );

    check(response, {
      'CRUD update status is 200': (r) => r.status === 200,
    });

    // 4. Deletar
    response = http.del(`${BASE_URL}/api/users/${user.email}`, null, {
      tags: { operation: 'crud_delete' },
    });

    const deleteSuccess = check(response, {
      'CRUD delete status is 204': (r) => r.status === 204,
    });

    if (deleteSuccess) {
      usersDeleted.add(1);
    }
  }

  errorRate.add(response.status >= 400);
  apiLatency.add(response.timings.duration);
}

// Cenário 4: Cenários de Erro
function testErrorScenarios() {
  // 1. Tentar acessar usuário inexistente (404)
  let response = http.get(`${BASE_URL}/api/users/inexistente@test.com`, {
    tags: { operation: 'error_404' },
  });

  check(response, {
    'Expected 404 for non-existent user': (r) => r.status === 404,
  });

  // 2. Tentar criar usuário com dados inválidos (400)
  const invalidUser = { name: '', email: 'invalid-email', password: '123' };
  response = http.post(`${BASE_URL}/api/users`, JSON.stringify(invalidUser), {
    headers: { 'Content-Type': 'application/json' },
    tags: { operation: 'error_400' },
  });

  check(response, {
    'Expected 400 for invalid user data': (r) =>
      r.status === 400 || r.status === 422,
  });

  // 3. Forçar erro 500 (se endpoint existir)
  response = http.get(`${BASE_URL}/test/error500`, {
    tags: { operation: 'error_500' },
  });

  check(response, {
    'Expected 500 for error endpoint': (r) => r.status === 500,
  });

  // 4. Forçar panic (se endpoint existir)
  response = http.get(`${BASE_URL}/test/panic`, {
    tags: { operation: 'panic_test' },
  });

  check(response, {
    'Panic endpoint handled': (r) => r.status === 500,
  });

  // Não conta 404 e erros esperados como falhas
  errorRate.add(response.status >= 500);
  apiLatency.add(response.timings.duration);
}

// Função executada no setup (antes dos testes)
export function setup() {
  console.log('🚀 Iniciando testes de estresse...');
  console.log('📊 Verifique as métricas em:');
  console.log('   - Prometheus: http://localhost:9090');
  console.log('   - Grafana: http://localhost:3000');

  // Verificar se a aplicação está rodando
  const response = http.get(`${BASE_URL}/health`);
  if (response.status !== 200) {
    throw new Error(
      '❌ Aplicação não está rodando! Execute: LOCAL=true go run app/main.go'
    );
  }

  console.log('✅ Aplicação está rodando');
  return { baseUrl: BASE_URL };
}

// Função executada no teardown (após os testes)
export function teardown(data) {
  console.log('✅ Testes de estresse concluídos!');
  console.log('📈 Verificar dashboards para análise dos resultados');
}
