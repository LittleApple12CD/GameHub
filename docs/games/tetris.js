import { BaseGame } from './BaseGame.js';

const SHAPES = {
  I: [[0,0,0,0],[1,1,1,1],[0,0,0,0],[0,0,0,0]],
  O: [[1,1],[1,1]],
  T: [[0,1,0],[1,1,1],[0,0,0]],
  S: [[0,1,1],[1,1,0],[0,0,0]],
  Z: [[1,1,0],[0,1,1],[0,0,0]],
  L: [[1,0,0],[1,1,1],[0,0,0]],
  J: [[0,0,1],[1,1,1],[0,0,0]],
};
const COLORS = {
  I: '#00f0f0', O: '#f0f000', T: '#b400f0',
  S: '#00f000', Z: '#f00000', L: '#f0a000', J: '#0000f0',
};

export class TetrisGame extends BaseGame {
  constructor(canvas, ctx) {
    super(canvas, ctx);
    this.cols = 10;
    this.rows = 20;
    this.cell = 30;
    this.boardW = this.cols * this.cell;   // 300
    this.boardH = this.rows * this.cell;   // 600
    this.offsetX = Math.floor((this.width - this.boardW - 140) / 2);
    this.offsetY = Math.floor((this.height - this.boardH) / 2);
  }

  reset() {
    this.board = Array.from({ length: this.rows }, () => Array(this.cols).fill(0));
    this.score = 0;
    this.gameOver = false;
    this.fallTimer = 0;
    this.fallDelay = 500;
    this.spawn();
  }

  spawn() {
    this.pieceKey = Object.keys(SHAPES)[Math.floor(Math.random() * 7)];
    this.piece = SHAPES[this.pieceKey].map(r => [...r]);
    this.px = Math.floor(this.cols / 2 - this.piece[0].length / 2);
    this.py = 0;
    if (this.collide(this.piece, this.px, this.py)) this.gameOver = true;
  }

  collide(shape, x, y) {
    for (let r = 0; r < shape.length; r++) {
      for (let c = 0; c < shape[r].length; c++) {
        if (!shape[r][c]) continue;
        const bx = x + c, by = y + r;
        if (bx < 0 || bx >= this.cols || by >= this.rows) return true;
        if (by >= 0 && this.board[by][bx]) return true;
      }
    }
    return false;
  }

  rotate() {
    const s = this.piece;
    const rotated = s[0].map((_, i) => s.map(row => row[i]).reverse());
    if (!this.collide(rotated, this.px, this.py)) this.piece = rotated;
  }

  lock() {
    this.piece.forEach((row, r) => row.forEach((v, c) => {
      if (v) {
        const by = this.py + r, bx = this.px + c;
        if (by >= 0) this.board[by][bx] = this.pieceKey;
      }
    }));
    this.clearLines();
    this.spawn();
  }

  clearLines() {
    let cleared = 0;
    for (let r = this.rows - 1; r >= 0; r--) {
      if (this.board[r].every(v => v)) {
        this.board.splice(r, 1);
        this.board.unshift(Array(this.cols).fill(0));
        cleared++;
      }
    }
    if (cleared) this.score += [0, 100, 300, 500, 800][cleared];
  }

  onKeyDown(e) {
    if (e.key.toLowerCase() === 'r') { this.reset(); return; }
    if (this.gameOver) return;
    switch (e.key) {
      case 'ArrowLeft':
        if (!this.collide(this.piece, this.px - 1, this.py)) this.px--;
        break;
      case 'ArrowRight':
        if (!this.collide(this.piece, this.px + 1, this.py)) this.px++;
        break;
      case 'ArrowDown':
        if (!this.collide(this.piece, this.px, this.py + 1)) this.py++;
        break;
      case 'ArrowUp':
      case ' ':
        this.rotate();
        break;
    }
  }

  update(dt) {
    if (this.gameOver) return;
    this.fallTimer += dt;
    if (this.fallTimer >= this.fallDelay) {
      this.fallTimer = 0;
      if (!this.collide(this.piece, this.px, this.py + 1)) this.py++;
      else this.lock();
    }
  }

  draw() {
    this.clear();
    const c = this.ctx;
    const { offsetX: ox, offsetY: oy, cell } = this;

    // 棋盘背景
    c.fillStyle = '#14141e';
    c.fillRect(ox, oy, this.boardW, this.boardH);

    // 已落方块
    for (let r = 0; r < this.rows; r++) {
      for (let col = 0; col < this.cols; col++) {
        const v = this.board[r][col];
        if (v) {
          c.fillStyle = COLORS[v];
          this.roundRect(ox + col * cell + 1, oy + r * cell + 1, cell - 2, cell - 2, 4);
        }
      }
    }

    // 当前方块
    if (!this.gameOver) {
      c.fillStyle = COLORS[this.pieceKey];
      this.piece.forEach((row, r) => row.forEach((v, col) => {
        if (v) {
          this.roundRect(ox + (this.px + col) * cell + 1,
                         oy + (this.py + r) * cell + 1,
                         cell - 2, cell - 2, 4);
        }
      }));
    }

    // 分数
    this.drawText(`Score: ${this.score}`, ox + this.boardW + 20, oy + 20, { size: 28 });

    if (this.gameOver) {
      c.fillStyle = 'rgba(0,0,0,0.7)';
      c.fillRect(ox, oy, this.boardW, this.boardH);
      this.drawText('GAME OVER', ox + this.boardW / 2, oy + this.boardH / 2 - 20,
                    { size: 40, align: 'center', baseline: 'middle' });
      this.drawText('Press R to restart', ox + this.boardW / 2, oy + this.boardH / 2 + 30,
                    { size: 22, align: 'center', baseline: 'middle' });
    }
  }
}