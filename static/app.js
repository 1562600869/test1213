let currentAssignReservationId = null;

document.querySelectorAll('.tab').forEach(tab => {
  tab.addEventListener('click', () => {
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
    document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
    tab.classList.add('active');
    document.getElementById(tab.dataset.tab).classList.add('active');
    if (tab.dataset.tab === 'halls') loadHalls();
    if (tab.dataset.tab === 'guides') loadGuides();
    if (tab.dataset.tab === 'reservations') { loadReservations(); loadHallOptions(); }
    if (tab.dataset.tab === 'today') loadToday();
    if (tab.dataset.tab === 'stats') loadStats();
  });
});

async function api(url, options = {}) {
  const res = await fetch(url, Object.assign({
    headers: { 'Content-Type': 'application/json' }
  }, options));
  const data = await res.json();
  if (!res.ok) {
    alert(data.error || '操作失败');
    throw new Error(data.error || 'error');
  }
  return data;
}

function statusClass(s) {
  if (s === '开放' || s === '在职') return 'status-open';
  if (s === '关闭' || s === '离职') return 'status-closed';
  return 'status-maintenance';
}

async function loadHalls() {
  const halls = await api('/api/halls');
  const tbody = document.getElementById('hallList');
  tbody.innerHTML = halls.map(h => `
    <tr>
      <td>${h.id}</td>
      <td><input value="${h.name}" id="hallName_${h.id}"></td>
      <td>
        <select id="hallTheme_${h.id}">
          <option ${h.theme==='历史人物'?'selected':''}>历史人物</option>
          <option ${h.theme==='影视明星'?'selected':''}>影视明星</option>
          <option ${h.theme==='体育冠军'?'selected':''}>体育冠军</option>
          <option ${h.theme==='世界领袖'?'selected':''}>世界领袖</option>
        </select>
      </td>
      <td><input type="number" value="${h.max_capacity}" id="hallCapacity_${h.id}" min="1"></td>
      <td>
        <select id="hallStatus_${h.id}" class="${statusClass(h.status)}">
          <option value="开放" ${h.status==='开放'?'selected':''}>开放</option>
          <option value="关闭" ${h.status==='关闭'?'selected':''}>关闭</option>
          <option value="维修" ${h.status==='维修'?'selected':''}>维修</option>
        </select>
      </td>
      <td>
        <button onclick="updateHall(${h.id})">保存</button>
        <button class="danger" onclick="deleteHall(${h.id})">删除</button>
      </td>
    </tr>
  `).join('');
}

async function addHall() {
  const name = document.getElementById('hallName').value.trim();
  const theme = document.getElementById('hallTheme').value;
  const capacity = parseInt(document.getElementById('hallCapacity').value);
  const status = document.getElementById('hallStatus').value;
  if (!name || !capacity) { alert('请填写完整信息'); return; }
  await api('/api/halls', { method: 'POST', body: JSON.stringify({ name, theme, max_capacity: capacity, status }) });
  document.getElementById('hallName').value = '';
  document.getElementById('hallCapacity').value = '';
  loadHalls();
}

async function updateHall(id) {
  const name = document.getElementById('hallName_' + id).value.trim();
  const theme = document.getElementById('hallTheme_' + id).value;
  const max_capacity = parseInt(document.getElementById('hallCapacity_' + id).value);
  const status = document.getElementById('hallStatus_' + id).value;
  await api('/api/halls', { method: 'PUT', body: JSON.stringify({ id, name, theme, max_capacity, status }) });
  loadHalls();
}

async function deleteHall(id) {
  if (!confirm('确定删除该展厅？')) return;
  await api('/api/halls?id=' + id, { method: 'DELETE' });
  loadHalls();
}

async function loadGuides() {
  const guides = await api('/api/guides');
  const tbody = document.getElementById('guideList');
  tbody.innerHTML = guides.map(g => `
    <tr>
      <td>${g.id}</td>
      <td><input value="${g.nickname}" id="guideNickname_${g.id}"></td>
      <td><input value="${g.phone}" id="guidePhone_${g.id}"></td>
      <td>
        <select id="guideLanguage_${g.id}">
          <option ${g.language==='普通话'?'selected':''}>普通话</option>
          <option ${g.language==='粤语'?'selected':''}>粤语</option>
          <option ${g.language==='英语'?'selected':''}>英语</option>
          <option ${g.language==='日语'?'selected':''}>日语</option>
        </select>
      </td>
      <td>
        <select id="guideStatus_${g.id}" class="${statusClass(g.status)}">
          <option value="在职" ${g.status==='在职'?'selected':''}>在职</option>
          <option value="离职" ${g.status==='离职'?'selected':''}>离职</option>
        </select>
      </td>
      <td>
        <button onclick="updateGuide(${g.id})">保存</button>
        <button class="danger" onclick="deleteGuide(${g.id})">删除</button>
      </td>
    </tr>
  `).join('');
}

async function addGuide() {
  const nickname = document.getElementById('guideNickname').value.trim();
  const phone = document.getElementById('guidePhone').value.trim();
  const language = document.getElementById('guideLanguage').value;
  const status = document.getElementById('guideStatus').value;
  if (!nickname || !phone) { alert('请填写完整信息'); return; }
  await api('/api/guides', { method: 'POST', body: JSON.stringify({ nickname, phone, language, status }) });
  document.getElementById('guideNickname').value = '';
  document.getElementById('guidePhone').value = '';
  loadGuides();
}

async function updateGuide(id) {
  const nickname = document.getElementById('guideNickname_' + id).value.trim();
  const phone = document.getElementById('guidePhone_' + id).value.trim();
  const language = document.getElementById('guideLanguage_' + id).value;
  const status = document.getElementById('guideStatus_' + id).value;
  await api('/api/guides', { method: 'PUT', body: JSON.stringify({ id, nickname, phone, language, status }) });
  loadGuides();
}

async function deleteGuide(id) {
  if (!confirm('确定删除该导览员？')) return;
  await api('/api/guides?id=' + id, { method: 'DELETE' });
  loadGuides();
}

async function loadHallOptions() {
  const halls = await api('/api/halls');
  const openHalls = halls.filter(h => h.status === '开放');
  const sel = document.getElementById('resHall');
  sel.innerHTML = '<option value="">选择展厅</option>' +
    openHalls.map(h => `<option value="${h.id}">${h.name} (${h.theme}, 最大${h.max_capacity}人)</option>`).join('');
}

async function loadReservations() {
  const list = await api('/api/reservations');
  const tbody = document.getElementById('reservationList');
  tbody.innerHTML = list.map(r => `
    <tr>
      <td>${r.id}</td>
      <td>${r.guest_name}</td>
      <td>${r.guest_phone}</td>
      <td>${r.hall_name || '-'}</td>
      <td>${r.time_slot}</td>
      <td>${r.people_count}</td>
      <td>${r.guide_name || '<span style="color:#999">未分配</span>'}</td>
      <td>
        <button onclick="openAssignModal(${r.id})">分配导览</button>
        <button class="danger" onclick="deleteReservation(${r.id})">删除</button>
      </td>
    </tr>
  `).join('');
}

async function addReservation() {
  const guest_name = document.getElementById('guestName').value.trim();
  const guest_phone = document.getElementById('guestPhone').value.trim();
  const hall_id = parseInt(document.getElementById('resHall').value);
  let time_slot = document.getElementById('resTime').value;
  const people_count = parseInt(document.getElementById('resPeople').value);
  if (!guest_name || !guest_phone || !hall_id || !time_slot || !people_count) {
    alert('请填写完整信息'); return;
  }
  time_slot = time_slot.replace('T', ' ') + ':00';
  await api('/api/reservations', { method: 'POST', body: JSON.stringify({ guest_name, guest_phone, hall_id, time_slot, people_count }) });
  document.getElementById('guestName').value = '';
  document.getElementById('guestPhone').value = '';
  document.getElementById('resTime').value = '';
  document.getElementById('resPeople').value = '';
  loadReservations();
}

async function deleteReservation(id) {
  if (!confirm('确定删除该预约？')) return;
  await api('/api/reservations?id=' + id, { method: 'DELETE' });
  loadReservations();
}

async function openAssignModal(reservationId) {
  currentAssignReservationId = reservationId;
  const guides = await api('/api/guides');
  const activeGuides = guides.filter(g => g.status === '在职');
  const sel = document.getElementById('assignGuideSelect');
  sel.innerHTML = activeGuides.length ? activeGuides.map(g => `<option value="${g.id}">${g.nickname} (${g.language})</option>`).join('') : '<option value="">无在职导览员</option>';
  document.getElementById('assignModal').style.display = 'flex';
}

function closeAssignModal() {
  document.getElementById('assignModal').style.display = 'none';
  currentAssignReservationId = null;
}

async function confirmAssign() {
  const guide_id = parseInt(document.getElementById('assignGuideSelect').value);
  if (!guide_id) { alert('请选择导览员'); return; }
  await api('/api/assign', { method: 'POST', body: JSON.stringify({ reservation_id: currentAssignReservationId, guide_id }) });
  closeAssignModal();
  loadReservations();
}

async function loadToday() {
  const now = new Date();
  document.getElementById('todayDate').textContent = now.toLocaleDateString('zh-CN');
  const list = await api('/api/reservations/today');
  const tbody = document.getElementById('todayList');
  if (list.length === 0) {
    tbody.innerHTML = '<tr><td colspan="7" style="text-align:center;color:#999;padding:40px">今日暂无预约</td></tr>';
    return;
  }
  tbody.innerHTML = list.map(r => `
    <tr>
      <td>${r.id}</td>
      <td>${r.guest_name}</td>
      <td>${r.guest_phone}</td>
      <td>${r.hall_name || '-'}</td>
      <td>${r.time_slot}</td>
      <td>${r.people_count}</td>
      <td>${r.guide_name || '<span style="color:#999">未分配</span>'}</td>
    </tr>
  `).join('');
}

async function loadStats() {
  const stats = await api('/api/stats/monthly-theme');
  const tbody = document.getElementById('statsList');
  if (stats.length === 0) {
    tbody.innerHTML = '<tr><td colspan="2" style="text-align:center;color:#999;padding:40px">暂无数据</td></tr>';
    return;
  }
  const total = stats.reduce((s, x) => s + x.count, 0);
  tbody.innerHTML = stats.map(s => `
    <tr>
      <td>${s.theme}</td>
      <td>${s.count} 次</td>
    </tr>
  `).join('') + `<tr style="font-weight:bold;background:#fafafa"><td>合计</td><td>${total} 次</td></tr>`;
}

loadHalls();
