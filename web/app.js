const API_KEY = 'orchestr-api-key';

async function api(path) {
  const res = await fetch(path, { headers: { 'X-API-Key': API_KEY } });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text);
  }
  return res.json();
}

function statusClass(status) {
  if (status === 'pending') return 'status-pending';
  if (status === 'running') return 'status-running';
  if (status === 'completed') return 'status-completed';
  if (status === 'failed') return 'status-failed';
  if (status === 'timeout') return 'status-timeout';
  return '';
}

async function loadStats() {
  try {
    const body = await api('/api/stats/overview');
    const data = body.data || {};
    const container = document.getElementById('stats');
    container.innerHTML = '';
    const items = [
      { label: '总执行数', value: data.total_executions ?? 0 },
      { label: '已完成', value: data.completed_executions ?? 0 },
      { label: '失败', value: data.failed_executions ?? 0 },
      { label: '超时', value: data.timeout_executions ?? 0 },
      { label: '成功率(%)', value: (data.success_rate ?? 0).toFixed(1) },
      { label: '平均耗时(ms)', value: data.avg_duration_ms ?? 0 },
    ];
    items.forEach(item => {
      const div = document.createElement('div');
      div.className = 'stat-card';
      div.innerHTML = `<div class="value">${item.value}</div><div class="label">${item.label}</div>`;
      container.appendChild(div);
    });
  } catch (e) {
    document.getElementById('stats').textContent = '统计加载失败: ' + e.message;
  }
}

async function loadExecutions() {
  try {
    const body = await api('/api/executions');
    const items = (body.data && body.data.items) || [];
    const tbody = document.querySelector('#list tbody');
    tbody.innerHTML = '';
    items.forEach(e => {
      const tr = document.createElement('tr');
      tr.innerHTML = `<td>${e.id}</td><td>${e.flow_id}</td><td class="${statusClass(e.status)}">${e.status}</td><td>${e.created_at}</td>`;
      tbody.appendChild(tr);
    });
  } catch (e) {
    const tbody = document.querySelector('#list tbody');
    tbody.innerHTML = `<tr><td colspan="4">加载失败: ${e.message}</td></tr>`;
  }
}

async function loadFlows() {
  try {
    const body = await api('/api/flows');
    const items = (body.data && body.data.items) || [];
    const tbody = document.querySelector('#flow-list tbody');
    tbody.innerHTML = '';
    items.forEach(f => {
      const tr = document.createElement('tr');
      tr.innerHTML = `<td>${f.id}</td><td>${f.name}</td><td class="${statusClass(f.status)}">${f.status}</td><td>${f.created_at}</td>`;
      tbody.appendChild(tr);
    });
  } catch (e) {
    const tbody = document.querySelector('#flow-list tbody');
    tbody.innerHTML = `<tr><td colspan="4">加载失败: ${e.message}</td></tr>`;
  }
}

async function load() {
  await Promise.all([loadStats(), loadExecutions(), loadFlows()]);
}

load();
