import { BaseGame } from './BaseGame.js';

export class TicTacToeGame extends BaseGame {
  constructor(canvas, ctx) {
    super(canvas, ctx);
    this.areaSize = 500;
    this.offsetX = Math.floor((this.width - this.areaSize) / 2);
    this.offsetY = Math.floor((this.height - this.areaSize) / 2) - 20;
    this.cellSize = this.areaSize / 3;
  }

  reset() {
    this.board = Array.from({ length: 3 }, () => ['', '', '']);
    this.currentPlayer = 'X';
    this.winner = null;
    this.gameOver = false;
    this.moveCount = 0;
  }

  checkWinner() {
    const b = this.board;
    for (let i = 0; i < 3; i++) {
      if (b[i][0] && b[i][0] === b[i][1] && b[i][1] === b[i][2]) return b[i][0];
      if (b[0][i] && b[0][i] === b[1][i] && b[1][i] === b[2][i]) return b[0][i];
    }
    if (b[0][0] && b[0][0] === b[1][1] && b[1][1] === b[2][2]) return b[0][0];
    if (b[0][2] && b[0][2] === b[1][1] && b[1][1] === b[2][0]) return b[0][2];
    return null;
  }

  onKeyDown(e) {
    if (e.key.toLowerCase() === 'r') this.reset();
  }

  onClick(e) {
    if (this.gameOver) return;
    const rect = this.canvas.getBoundingClientRect();
    const scaleX = this.canvas.width / rect.width;
    const scaleY = this.canvas.height / rect.height;
    const mx = (e.clientX - rect.left) * scaleX - this.offsetX;
    const my = (e.clientY - rect.top) * scaleY - this.offsetY;
    const x = Math.floor(mx / this.cellSize);
    const y = Math.floor(my / this.cellSize);
    if (x < 0 || x > 2 || y < 0 || y > 2) return;
    if (this.board[y][x] !== '') return;

    this.board[y][x] = this.currentPlayer;
    this.moveCount++;
    this.winner = this.checkWinner();
    if (this.winner || this.moveCount === 9) this.gameOver = true;
    else this.currentPlayer = this.currentPlayer === 'X' ? 'O' : 'X';
  }

  draw() {
    this.clear();
    const c = this.ctx;
    const { offsetX: ox, offsetY: oy, cellSize, areaSize } = this;

    c.fillStyle = '#1e1e2a';
    c.fillRect(ox, oy, areaSize, areaSize);

    // 格子
    for (let r = 0; r < 3; r++) {
      for (let col = 0; col < 3; col++) {
        this.roundRect(ox + col * cellSize + 2, oy + r * cellSize + 2,
                       cellSize - 4, cellSize - 4, 8, '#323241');
        const v = this.board[r][col];
        if (v) {
          this.drawText(v, ox + col * cellSize + cellSize / 2,
                        oy + r * cellSize + cellSize / 2,
                        { size: 80, color: v === 'X' ? '#3cc8ff' : '#ffc83c',
                          align: 'center', baseline: 'middle' });
        }
      }
    }

    // 分隔线
    c.strokeStyle = '#505064';
    c.lineWidth = 3;
    for (let i = 1; i < 3; i++) {
      c.beginPath();
      c.moveTo(ox + i * cellSize, oy + 5);
      c.lineTo(ox + i * cellSize, oy + areaSize - 5);
      c.stroke();
      c.beginPath();
      c.moveTo(ox + 5, oy + i * cellSize);
      c.lineTo(ox + areaSize - 5, oy + i * cellSize);
      c.stroke();
    }

    // 状态
    const statusY = oy + areaSize + 30;
    if (this.gameOver) {
      const msg = this.winner ? `Player ${this.winner} Wins!` : 'Draw!';
      const color = this.winner ? '#00ff00' : '#ffff64';
      this.drawText(msg, this.width / 2, statusY,
                    { size: 40, color, align: 'center', baseline: 'middle' });
      this.drawText('Press R to restart', this.width / 2, statusY + 40,
                    { size: 22, color: '#c8c8c8', align: 'center', baseline: 'middle' });
    } else {
      this.drawText(`Player ${this.currentPlayer}'s Turn`, this.width / 2, statusY,
                    { size: 36, align: 'center', baseline: 'middle' });
    }
  }
}