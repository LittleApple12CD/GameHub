import { SnakeGame } from './games/snake.js';
import { TetrisGame } from './games/tetris.js';
import { MinesweeperGame } from './games/minesweeper.js';
import { FlappyGame } from './games/flappy.js';
import { TankGame } from './games/tank.js';
import { BreakoutGame } from './games/breakout.js';
import { TicTacToeGame } from './games/tictactoe.js';

const GAMES = [
  { name: 'Snake',         cls: SnakeGame },
  { name: 'Tetris',        cls: TetrisGame },
  { name: 'Minesweeper',   cls: MinesweeperGame },
  { name: 'Flappy Bird',   cls: FlappyGame },
  { name: 'Tank',          cls: TankGame },
  { name: 'Breakout',      cls: BreakoutGame },
  { name: 'Tic Tac Toe',   cls: TicTacToeGame },
];

const menuEl   = document.getElementById('menu');
const listEl   = document.getElementById('game-list');
const exitBtn  = document.getElementById('exit-btn');
const gameScr  = document.getElementById('game-screen');
const canvas   = document.getElementById('game-canvas');
const ctx      = canvas.getContext('2d');
const titleEl  = document.getElementById('game-title');
const backBtn  = document.getElementById('back-btn');

let currentGame = null;
let selectedIndex = 0;

// 构建菜单
GAMES.forEach((g, i) => {
  const li = document.createElement('li');
  li.textContent = g.name;
  li.dataset.index = i;
  li.addEventListener('mouseenter', () => {
    selectedIndex = i;
    updateSelection();
  });
  li.addEventListener('click', () => launchGame(i));
  listEl.appendChild(li);
});

function updateSelection() {
  [...listEl.children].forEach((li, i) =>
    li.classList.toggle('selected', i === selectedIndex)
  );
}
updateSelection();

function exitApp() {
  window.close();

  setTimeout(() => {
    document.body.innerHTML = `
      <div style="
        display:flex; flex-direction:column; align-items:center; justify-content:center;
        min-height:100vh; background:#0f0f17; color:#fff;
        font-family:'Segoe UI',sans-serif; gap:16px;">
        <h1 style="font-size:42px; color:#50c8ff; letter-spacing:4px;">GAME HUB</h1>
        <p style="font-size:20px; color:#aaa;">已退出，感谢游玩！</p>
        <p style="font-size:14px; color:#666;">你可以直接关闭这个标签页。</p>
        <button onclick="location.reload()" style="
          margin-top:12px; padding:10px 30px; font-size:16px;
          background:#50c8ff; color:#0f0f17; border:none;
          border-radius:10px; cursor:pointer;">重新开始</button>
      </div>
    `;
  }, 100);
}

function confirmExit() {
  return confirm('确定要退出 Game Hub 吗？');
}

function backToMenu() {
  if (currentGame) {
    currentGame.destroy();
    currentGame = null;
  }
  gameScr.classList.remove('active');
  menuEl.classList.add('active');
  ctx.clearRect(0, 0, canvas.width, canvas.height);
}

document.addEventListener('keydown', (e) => {
  if (menuEl.classList.contains('active')) {
    if (e.key === 'ArrowUp') {
      selectedIndex = (selectedIndex - 1 + GAMES.length) % GAMES.length;
      updateSelection();
    } else if (e.key === 'ArrowDown') {
      selectedIndex = (selectedIndex + 1) % GAMES.length;
      updateSelection();
    } else if (e.key === 'Enter') {
      launchGame(selectedIndex);
    } else if (e.key === 'Escape') {
      if (confirmExit()) exitApp();
    }
  }
});

exitBtn.addEventListener('click', () => {
  if (confirmExit()) exitApp();
});

backBtn.addEventListener('click', backToMenu);

function launchGame(index) {
  const { name, cls } = GAMES[index];
  menuEl.classList.remove('active');
  gameScr.classList.add('active');
  titleEl.textContent = name;

  if (currentGame) currentGame.destroy();
  currentGame = new cls(canvas, ctx);
  currentGame.onExitToMenu = backToMenu;
  currentGame.start();
}

window.__gameHubBackToMenu = backToMenu;