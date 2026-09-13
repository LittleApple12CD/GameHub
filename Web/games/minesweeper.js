import { BaseGame } from './BaseGame.js';

export class MinesweeperGame extends BaseGame {
  constructor(canvas, ctx) {
    super(canvas, ctx);
    this.rows = 16;
    this.cols = 16;
    this.mineCount = 40;
    this.cell = 30;
    this.boardW = this.cols * this.cell;
    this.boardH = this.rows * this.cell;
    this.offsetX = Math.floor((this.width - this.boardW) / 2);
    this.offsetY = Math.floor((this.height - this.boardH) / 2);
  }

  reset() {
    this.board = Array.from({ length: this.rows }, () => Array(this.cols).fill(0));
    this.revealed = Array.from({ length: this.rows }, () => Array(this.cols).fill(false));
    this.flagged = Array.from({ length: this.rows }, () => Array(this.cols).fill(false));
    this.gameOver = false;
    this.won = false;
    this.firstClick = true;
  }

  placeMines(sx, sy) {
    let placed = 0;
    while (placed < this.mineCount) {
      const x = Math.floor(Math.random() * this.cols);
      const y = Math.floor(Math.random() * this.rows);
      if (this.board[y][x] === -1) continue;
      if (Math.abs(x - sx) <= 1 && Math.abs(y - sy) <= 1) continue;
      this.board[y][x] = -1;
      placed++;
    }
    for (let y = 0; y < this.rows; y++) {
      for (let x = 0; x < this.cols; x++) {
        if (this.board[y][x] === -1) continue;
        let cnt = 0;
        for (let dy = -1; dy <= 1; dy++)
          for (let dx = -1; dx <= 1; dx++) {
            const nx = x + dx, ny = y + dy;
            if (nx >= 0 && nx < this.cols && ny >= 0 && ny < this.rows && this.board[ny][nx] === -1)
              cnt++;
          }
        this.board[y][x] = cnt;
      }
    }
  }

  reveal(x, y) {
    if (x < 0 || x >= this.cols || y < 0 || y >= this.rows) return;
    if (this.revealed[y][x] || this.flagged[y][x]) return;
    this.revealed[y][x] = true;
    if (this.board[y][x] === -1) { this.gameOver = true; return; }
    if (this.board[y][x] === 0) {
      for (let dy = -1; dy <= 1; dy++)
        for (let dx = -1; dx <= 1; dx++)
          this.reveal(x + dx, y + dy);
    }
    this.checkWin();
  }

  checkWin() {
    let revealed = 0;
    for (let y = 0; y < this.rows; y++)
      for (let x = 0; x < this.cols; x++)
        if (this.revealed[y][x]) revealed++;
    if (revealed === this.rows * this.cols - this.mineCount) this.won = true;
  }

  onKeyDown(e) {
    if (e.key.toLowerCase() === 'r') this.reset();
  }

  onClick(e) {
    if (this.gameOver || this.won) return;
    const cell = this._getCellFromEvent(e);
    if (!cell) return;
    const { x, y } = cell;
    if (this.firstClick) {
      this.placeMines(x, y);
      this.firstClick = false;
    }
    if (!this.flagged[y][x]) this.reveal(x, y);
  }

  onRightClick(e) {
    if (this.gameOver || this.won) return;
    const cell = this._getCellFromEvent(e);
    if (!cell) return;
    const { x, y } = cell;
    if (!this.revealed[y][x]) {
      this.flagged[y][x] = !this.flagged[y][x];
    }
  }

  _getCellFromEvent(e) {
    const rect = this.canvas.getBoundingClientRect();
    // 处理 canvas 被 CSS 缩放的情况
    const scaleX = this.canvas.width / rect.width;
    const scaleY = this.canvas.height / rect.height;
    const mx = (e.clientX - rect.left) * scaleX - this.offsetX;
    const my = (e.clientY - rect.top) * scaleY - this.offsetY;
    const x = Math.floor(mx / this.cell);
    const y = Math.floor(my / this.cell);
    if (x < 0 || x >= this.cols || y < 0 || y >= this.rows) return null;
    return { x, y };
  }

  draw() {
    this.clear();
    const c = this.ctx;
    const { offsetX: ox, offsetY: oy, cell } = this;

    for (let y = 0; y < this.rows; y++) {
      for (let x = 0; x < this.cols; x++) {
        const px = ox + x * cell + 1;
        const py = oy + y * cell + 1;
        const size = cell - 2;

        if (this.revealed[y][x]) {
          if (this.board[y][x] === -1) {
            this.roundRect(px, py, size, size, 4, '#c83232');
            c.fillStyle = '#000';
            c.beginPath();
            c.arc(px + size / 2, py + size / 2, 8, 0, Math.PI * 2);
            c.fill();
          } else {
            this.roundRect(px, py, size, size, 4, '#bebec8');
            const n = this.board[y][x];
            if (n > 0) {
              const colors = ['', '#0000ff', '#009600', '#ff0000', '#0000b4', '#960000', '#009696', '#000'];
              this.drawText(String(n), px + size / 2, py + size / 2,
                            { size: 18, color: colors[n], align: 'center', baseline: 'middle' });
            }
          }
        } else {
          const color = this.flagged[y][x] ? '#d2d23c' : '#505069';
          this.roundRect(px, py, size, size, 4, color);
          if (this.flagged[y][x]) {
            this.drawText('F', px + size / 2, py + size / 2,
                          { size: 18, color: '#ff3232', align: 'center', baseline: 'middle' });
          }
        }
      }
    }

    if (this.gameOver || this.won) {
      c.fillStyle = 'rgba(0,0,0,0.6)';
      c.fillRect(ox, oy, this.boardW, this.boardH);
      const msg = this.won ? 'YOU WIN!' : 'GAME OVER';
      const color = this.won ? '#00ff00' : '#ffffff';
      this.drawText(msg, ox + this.boardW / 2, oy + this.boardH / 2 - 20,
                    { size: 40, color, align: 'center', baseline: 'middle' });
      this.drawText('Press R to restart', ox + this.boardW / 2, oy + this.boardH / 2 + 30,
                    { size: 22, align: 'center', baseline: 'middle' });
    }
  }
}