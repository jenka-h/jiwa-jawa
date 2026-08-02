'use strict';

const $ = id => document.getElementById(id);
const state = { profiles: [], game: null, selected: null, poll: null, toast: null };
const boardWidth = Number(document.body.dataset.boardWidth);
const boardHeight = Number(document.body.dataset.boardHeight);
const svgNS = 'http://www.w3.org/2000/svg';

async function api(path, options) {
  try {
    const response = await fetch(path, options);
    if (!response.ok) throw new Error((await response.text()).trim() || `Request failed (${response.status})`);
    return await response.json();
  } catch (error) {
    notify(error.message || 'Could not reach the game server.', 'error');
    throw error;
  }
}

async function loadProfiles() {
  state.profiles = await api('/api/profiles');
  $('profileSelect').replaceChildren(...state.profiles.map(profile => {
    const option = document.createElement('option');
    option.value = profile.id;
    option.textContent = `${profile.name} · ${Math.round(profile.rating)}`;
    return option;
  }));
}

function switchTab(name) {
  const selecting = name === 'select';
  $('selectTab').classList.toggle('active', selecting);
  $('createTab').classList.toggle('active', !selecting);
  $('selectPanel').classList.toggle('active', selecting);
  $('createPanel').classList.toggle('active', !selecting);
}

async function createProfile() {
  const input = $('profileName');
  const name = input.value.trim();
  if (!name) return notify('Enter a player name.', 'error');
  const profile = await api('/api/profiles', jsonOptions('POST', { name }));
  input.value = '';
  await loadProfiles();
  $('profileSelect').value = profile.id;
  switchTab('select');
  notify(`Profile created: ${profile.name}`, 'success');
}

async function startSession(role, button) {
  const profileID = $('profileSelect').value;
  if (!profileID) return notify('Create or select a profile first.', 'error');
  button.disabled = true;
  try {
    const game = await api('/api/session', jsonOptions('POST', {
      profile_id: profileID,
      role,
      peer: $('peer').value.trim()
    }));
    openGame(game);
    notify(role === 'host' ? 'Waiting for an opponent…' : 'Joining game…', 'success');
  } finally {
    button.disabled = false;
  }
}

function openGame(game) {
  $('menu').classList.remove('active');
  $('game').classList.add('active');
  render(game);
  stopPolling();
  state.poll = window.setInterval(refresh, 800);
}

function backToMenu() {
  stopPolling();
  $('game').classList.remove('active');
  $('menu').classList.add('active');
}

function stopPolling() {
  if (state.poll) window.clearInterval(state.poll);
  state.poll = null;
}

async function refresh() {
  if (document.hidden) return;
  try {
    const game = await api('/api/state');
    if (game.active) render(game);
  } catch (_) {
    // api() already displays a useful error; polling will retry.
  }
}

function render(game) {
  state.game = game;
  if (state.selected && !isLocalPiece(game, state.selected)) state.selected = null;
  const myTurn = game.turn === game.local_side;
  const finished = game.finished;

  text('profileNameA', game.local_name || 'Local player');
  text('profileMetaA', `${game.role || 'player'} · ${sideName(game.local_side)}`);
  text('profileEloA', formatRating(game.local_rating));
  text('profileInitialA', initials(game.local_name));
  text('profileNameB', game.remote_name || 'Waiting…');
  text('profileMetaB', `opponent · ${sideName(game.remote_side)}`);
  text('profileEloB', game.remote_name ? formatRating(game.remote_rating) : '—');
  text('profileInitialB', initials(game.remote_name));
  const myPenalty = game.penalty_active && game.penalty_side === game.local_side;
  const myContinuation = game.can_end_turn && myTurn;
  text('turnTitle', finished ? 'Game finished' : myPenalty ? 'Dam Ora Mangan' : myContinuation ? 'Capture again?' : myTurn ? 'Your turn' : game.status === 'ready' ? "Opponent's turn" : 'Waiting for opponent');
  text('turnHint', finished ? `Winner: ${sideName(game.winner) || 'draw'}` : myPenalty ? `Remove ${game.penalties_remaining} opponent piece${game.penalties_remaining === 1 ? '' : 's'}` : game.penalty_active ? 'Opponent is choosing your penalty pieces' : myContinuation ? 'Choose any highlighted piece to capture again, or end your turn' : myTurn ? 'Choose a piece, then a target' : game.status === 'ready' ? 'Waiting for their move' : 'Share your address with the other player');
  text('moveNumber', game.move_number || 0);
  text('darkCaptured', Math.max(0, 16 - countPieces(game, 'player_one')));
  text('lightCaptured', Math.max(0, 16 - countPieces(game, 'player_two')));
  text('timer', formatTime(game.duration_seconds || 0));
  $('surrender').disabled = finished || game.status !== 'ready';
  $('endTurn').hidden = !myContinuation || finished;
  renderBoard(game);
  renderLogs(game.logs || []);
}

function renderBoard(game) {
  document.querySelectorAll('.live-piece').forEach(piece => piece.remove());
  const board = document.querySelector('.board');
  for (const point of game.board.points || []) {
    if (!point.owner) continue;
    const position = project(point.x, point.y);
    const group = svg('g', {
      class: `live-piece piece-group${state.selected === point.id ? ' selected' : ''}${game.penalty_active && game.penalty_side === game.local_side && point.owner === game.remote_side ? ' penalty-target' : ''}${game.can_end_turn && (game.capture_sources || []).includes(point.id) ? ' continuation-piece' : ''}`,
      'data-point': point.id,
      role: 'button',
      tabindex: '0',
      'aria-label': `${sideName(point.owner)} piece at point ${point.id}`
    });
    group.append(
      svg('circle', { class: `piece-ring ${point.owner === 'player_one' ? 'piece-dark' : 'piece-light'}`, cx: position.x, cy: position.y, r: 27 }),
      svg('image', { class: 'piece-image', href: point.owner === 'player_one' ? '/assets/p1.png' : '/assets/p2.png', x: position.x - 23, y: position.y - 23, width: 46, height: 46, preserveAspectRatio: 'xMidYMid slice' })
    );
    board.appendChild(group);
  }
}

function renderLogs(logs) {
  $('gameLogs').replaceChildren(...logs.slice(-12).reverse().map(item => {
    const row = document.createElement('div');
    row.className = 'log-row';
    const number = document.createElement('b'); number.textContent = `${item.no}.`;
    const player = document.createElement('span'); player.textContent = `${item.player} · ${item.action}`;
    const move = document.createElement('b'); move.textContent = item.move;
    row.append(number, player, move);
    return row;
  }));
}

async function choosePoint(id) {
  const game = state.game;
  if (!game || game.finished) return;
  if (game.status !== 'ready') return notify('The opponent has not connected yet.', 'error');
  if (game.turn !== game.local_side) return notify('Wait for your turn.', 'error');
  const point = (game.board.points || []).find(item => item.id === id);

  if (game.penalty_active) {
    if (game.penalty_side !== game.local_side) return notify('Opponent is resolving Dam Ora Mangan.', 'error');
    if (!point || point.owner !== game.remote_side) return notify('Choose one of the opponent’s pieces to remove.', 'error');
    const updated = await api('/api/penalty', jsonOptions('POST', { target: id }));
    state.selected = null;
    render(updated);
    return;
  }

  if (game.can_end_turn && point && point.owner === game.local_side && !(game.capture_sources || []).includes(point.id)) {
    return notify('Choose any highlighted capturing piece, or end your turn.', 'error');
  }

  if (!state.selected) {
    if (!point || point.owner !== game.local_side) return notify('Choose one of your pieces.', 'error');
    state.selected = id;
    renderBoard(game);
    return notify(`Point ${id} selected. Choose a target.`);
  }
  if (state.selected === id) {
    state.selected = null;
    renderBoard(game);
    return;
  }
  if (point && point.owner === game.local_side) {
    state.selected = id;
    renderBoard(game);
    return;
  }

  const source = state.selected;
  const updated = await api('/api/move', jsonOptions('POST', { source, target: id }));
  state.selected = null;
  render(updated);
}

async function endTurn() {
  const updated = await api('/api/end-turn', { method: 'POST' });
  state.selected = null;
  render(updated);
}

async function surrender() {
  if (!window.confirm('Surrender this game?')) return;
  render(await api('/api/finish', { method: 'POST' }));
}

function handleBoardAction(event) {
  const target = event.target.closest('[data-point]');
  if (!target) return;
  if (event.type === 'keydown' && event.key !== 'Enter' && event.key !== ' ') return;
  event.preventDefault();
  choosePoint(Number(target.dataset.point)).catch(() => {});
}

function jsonOptions(method, body) {
  return { method, headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) };
}
function text(id, value) { $(id).textContent = value ?? ''; }
function countPieces(game, owner) { return (game.board.points || []).filter(point => point.owner === owner).length; }
function isLocalPiece(game, id) { return (game.board.points || []).some(point => point.id === id && point.owner === game.local_side); }
function sideName(side) { return side === 'player_one' ? 'Dark' : side === 'player_two' ? 'Light' : ''; }
function initials(name = '') { return name.trim().split(/\s+/).slice(0, 2).map(part => part[0]?.toUpperCase()).join('') || 'P'; }
function formatTime(seconds) { const h = Math.floor(seconds / 3600); const m = Math.floor(seconds % 3600 / 60); const s = seconds % 60; return h ? [h, m, s].map(value => String(value).padStart(2, '0')).join(':') : [m, s].map(value => String(value).padStart(2, '0')).join(':'); }
function formatRating(value) { return Math.round(Number(value) || 1200); }
function project(x, y) { return { x: 70 + x / 8 * (boardWidth - 140), y: 60 + y / 4 * (boardHeight - 120) }; }
function svg(tag, attributes) { const node = document.createElementNS(svgNS, tag); for (const [key, value] of Object.entries(attributes)) node.setAttribute(key, value); return node; }
function notify(message, type = '') { const toast = $('toast'); toast.textContent = message; toast.className = `toast show ${type}`; window.clearTimeout(state.toast); state.toast = window.setTimeout(() => toast.className = 'toast', 2800); }

for (const tab of document.querySelectorAll('[data-tab]')) tab.addEventListener('click', () => switchTab(tab.dataset.tab));
for (const button of document.querySelectorAll('[data-role]')) button.addEventListener('click', () => startSession(button.dataset.role, button).catch(() => {}));
$('createProfile').addEventListener('click', () => createProfile().catch(() => {}));
$('profileName').addEventListener('keydown', event => { if (event.key === 'Enter') createProfile().catch(() => {}); });
$('backMenu').addEventListener('click', backToMenu);
$('endTurn').addEventListener('click', () => endTurn().catch(() => {}));
$('surrender').addEventListener('click', () => surrender().catch(() => {}));
document.querySelector('.board').addEventListener('click', handleBoardAction);
document.querySelector('.board').addEventListener('keydown', handleBoardAction);
document.addEventListener('visibilitychange', () => { if (!document.hidden && state.game) refresh(); });
loadProfiles().catch(() => {});
