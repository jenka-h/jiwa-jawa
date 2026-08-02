package gui

const pageHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
  <meta name="theme-color" content="#f5efe6">
  <title>HoloDam</title>
  <link rel="icon" href="/assets/logo.png">
  <link rel="apple-touch-icon" href="/assets/logo.png">
  <link rel="stylesheet" href="/assets/app.css?v=6">
</head>
<body data-board-width="{{.Board.Width}}" data-board-height="{{.Board.Height}}">
  <main>
    <section id="menu" class="screen menu active" aria-labelledby="game-title">
      <img class="logo" src="/assets/logo.png" alt="">
      <h1 id="game-title">HoloDam</h1>
      <p class="subtitle">Javanese strategy game</p>

      <div class="tabs" role="tablist" aria-label="Profile options">
        <button id="selectTab" class="tab active" data-tab="select" role="tab">Select profile</button>
        <button id="createTab" class="tab" data-tab="create" role="tab">New profile</button>
      </div>
      <div id="selectPanel" class="tab-panel active">
        <label for="profileSelect">Player</label>
        <select id="profileSelect"></select>
      </div>
      <div id="createPanel" class="tab-panel">
        <label for="profileName">Player name</label>
        <div class="inline-form">
          <input id="profileName" maxlength="40" autocomplete="nickname" placeholder="Your name">
          <button id="createProfile" class="button secondary">Create</button>
        </div>
      </div>

      <div class="connection">
        <label for="peer">Opponent address</label>
        <input id="peer" inputmode="url" autocapitalize="off" spellcheck="false" value="127.0.0.1:9001" placeholder="192.168.1.20:9001">
        <small>Only needed when joining. Use the host device's LAN IP and UDP port.</small>
      </div>
      <div class="menu-actions">
        <button class="button" data-role="host">Host game</button>
        <button class="button secondary" data-role="join">Join game</button>
      </div>
    </section>

    <section id="game" class="screen app" aria-label="Active game">
      <header class="game-header">
        <button id="backMenu" class="icon-button" aria-label="Back to menu">☰</button>
        <div id="turnStatus" class="turn-status">
          <strong id="turnTitle">Waiting</strong>
          <span id="turnHint">Connecting to opponent</span>
        </div>
        <time id="timer">00:00</time>
      </header>

      <section class="players" aria-label="Players">
        <article class="player-card">
          <span id="profileInitialA" class="avatar dark">P1</span>
          <div><strong id="profileNameA">Player 1</strong><small id="profileMetaA">Host · Dark</small><span class="elo-badge">Elo <b id="profileEloA">1200</b></span></div>
        </article>
        <article class="player-card opponent">
          <span id="profileInitialB" class="avatar light">P2</span>
          <div><strong id="profileNameB">Waiting…</strong><small id="profileMetaB">Opponent · Light</small><span class="elo-badge">Elo <b id="profileEloB">—</b></span></div>
        </article>
      </section>

      <div class="game-layout">
        <section class="board-card" aria-label="Game board">
          <svg class="board" viewBox="0 0 {{.Board.Width}} {{.Board.Height}}" role="group" aria-label="HoloDam board">
            {{range .Board.Edges}}<line class="board-line" x1="{{.X1}}" y1="{{.Y1}}" x2="{{.X2}}" y2="{{.Y2}}"></line>{{end}}
            {{range .Board.Points}}<g class="board-point" data-point="{{.ID}}" role="button" tabindex="0" aria-label="Point {{.ID}}"><circle class="touch-target" cx="{{.X}}" cy="{{.Y}}" r="24"></circle><circle class="node" cx="{{.X}}" cy="{{.Y}}" r="8"></circle><text class="point-label" x="{{.X}}" y="{{.Y}}">{{.ID}}</text></g>{{end}}
          </svg>
        </section>

        <aside class="details">
          <section class="stats">
            <div><span>Move</span><strong id="moveNumber">0</strong></div>
            <div><span>Dark taken</span><strong id="darkCaptured">0</strong></div>
            <div><span>Light taken</span><strong id="lightCaptured">0</strong></div>
          </section>
          <div class="game-actions">
            <button id="endTurn" class="button secondary" hidden>End turn</button>
            <button id="surrender" class="button danger">Surrender</button>
          </div>
          <details open>
            <summary>Game log</summary>
            <div id="gameLogs" class="log-list"></div>
          </details>
        </aside>
      </div>
    </section>
  </main>
  <div id="toast" class="toast" role="status" aria-live="polite"></div>
  <script src="/assets/app.js?v=6" defer></script>
</body>
</html>`
