import { BaseGame } from './BaseGame.js';

export class BreakoutGame extends BaseGame {
  constructor(canvas, ctx) {
    super(canvas, ctx);
    this.areaSize = 600;
    this.offsetX = Math.floor((this.width - this.areaSize) / 2);
    this.offsetY = Math.floor((this.height - this.areaSize) / 2);
  }

  reset() {
    this.paddle = { x: this.areaSize / 2 - 60, y: this.areaSize - 40, w: 120, h: 16 };
    this.paddleSpeed = 7;
    this.ball = { x: this.areaSize / 2 - 10, y: this.areaSize - 70, w: 20, h: 20 };
    this.ballSpeed = 8;
    this.ballDx = 1;
    this.ballDy = -1;
    this.normalizeSpeed();
    this.bricks = [];
    this.score = 0;
    this.lives = 3;
    this.gameOver = false;
    this.waiting = true;

    const rows = 5, cols = 8;
    const bw = Math.floor((this.areaSize - 20) / cols) - 4;
    const bh = 22;
    const colors = ['#e63232', '#e69632', '#e6e632', '#32e632', '#3296e6'];
    for (let r = 0; r < rows; r++) {
      for (let c = 0; c < cols; c++) {
        this.bricks.push({
          x: 10 + c * (bw + 4),
          y: 40 + r * (bh + 4),
          w: bw, h: bh,
          color: colors[r % colors.length],
          alive: true,
        });
      }
    }
  }

  onKeyDown(e) {
    if (e.key.toLowerCase() === 'r') { this.reset(); return; }
    if (e.key === ' ' && this.waiting && !this.gameOver) this.waiting = false;
  }

  normalizeSpeed() {
    const target = this.ballSpeed;
    const current = Math.hypot(this.ballDx, this.ballDy);
    if (current === 0) return;
    const scale = target / current;
    this.ballDx *= scale;
    this.ballDy *= scale;
  }

  resetBall() {
    this.ball.x = this.areaSize / 2 - 10;
    this.ball.y = this.areaSize - 70;
    this.ballDx = Math.random() < 0.5 ? 1 : -1;
    this.ballDy = -1;
    this.normalizeSpeed();
  }

  update(dt) {
    if (this.gameOver || this.waiting) return;
    const f = dt / 16.67;

    if (this.keys['arrowleft']  && this.paddle.x > 0)
      this.paddle.x -= this.paddleSpeed * f;
    if (this.keys['arrowright'] && this.paddle.x + this.paddle.w < this.areaSize)
      this.paddle.x += this.paddleSpeed * f;

    this.ball.x += this.ballDx * f;
    this.ball.y += this.ballDy * f;

    if (this.ball.x <= 0 || this.ball.x + this.ball.w >= this.areaSize) this.ballDx *= -1;
    // 顶部反弹
    if (this.ball.y <= 0) this.ballDy *= -1;

    if (this.ball.y + this.ball.h >= this.areaSize) {
      this.lives--;
      if (this.lives <= 0) {
        this.gameOver = true;
      } else {
        this.waiting = true;
        this.resetBall();
      }
      return;
    }

    if (this.rectCollide(this.ball, this.paddle)) {
      this.ballDy = -Math.abs(this.ballDy);
      const hit = (this.ball.x + this.ball.w / 2 - (this.paddle.x + this.paddle.w / 2)) / (this.paddle.w / 2);
      this.ballDx = hit * this.ballSpeed * 0.9;
      if (Math.abs(this.ballDx) < 1.2) this.ballDx = this.ballDx >= 0 ? 1.8 : -1.8;
      this.ball.y = this.paddle.y - this.ball.h;
      this.normalizeSpeed();
    }

    for (const b of this.bricks) {
      if (!b.alive) continue;
      if (this.rectCollide(this.ball, b)) {
        b.alive = false;
        this.score += 10;

        const overlapTop = this.ball.y + this.ball.h - b.y;
        const overlapBottom = b.y + b.h - this.ball.y;
        const overlapLeft = this.ball.x + this.ball.w - b.x;
        const overlapRight = b.x + b.w - this.ball.x;
        const minOv = Math.min(overlapTop, overlapBottom, overlapLeft, overlapRight);

        if (minOv === overlapTop) {
          this.ball.y = b.y - this.ball.h;
          this.ballDy = -Math.abs(this.ballDy);
        } else if (minOv === overlapBottom) {
          this.ball.y = b.y + b.h;
          this.ballDy = Math.abs(this.ballDy);
        } else if (minOv === overlapLeft) {
          this.ball.x = b.x - this.ball.w;
          this.ballDx = -Math.abs(this.ballDx);
        } else {
          this.ball.x = b.x + b.w;
          this.ballDx = Math.abs(this.ballDx);
        }

        this.normalizeSpeed();
        break;
      }
    }

    if (this.bricks.every(b => !b.alive)) this.gameOver = true;
  }

  rectCollide(a, b) {
    return a.x < b.x + b.w && a.x + a.w > b.x && a.y < b.y + b.h && a.y + a.h > b.y;
  }

  draw() {
    this.clear();
    const c = this.ctx;
    const { offsetX: ox, offsetY: oy, areaSize } = this;

    c.fillStyle = '#14141e';
    c.fillRect(ox, oy, areaSize, areaSize);

    for (const b of this.bricks) {
      if (!b.alive) continue;
      this.roundRect(ox + b.x, oy + b.y, b.w, b.h, 8, b.color, '#ffffff');
    }

    this.roundRect(ox + this.paddle.x, oy + this.paddle.y,
                   this.paddle.w, this.paddle.h, 10, '#ffffff', '#c8c8dc');

    c.fillStyle = '#ffff78';
    c.beginPath();
    c.arc(ox + this.ball.x + this.ball.w / 2,
          oy + this.ball.y + this.ball.h / 2,
          this.ball.w / 2, 0, Math.PI * 2);
    c.fill();

    // HUD
    this.drawText(`Score: ${this.score}`, ox + 10, oy + 10, { size: 22 });
    this.drawText(`Lives: ${this.lives}`, ox + areaSize - 120, oy + 10, { size: 22 });

    if (this.waiting && !this.gameOver) {
      this.drawText('Press SPACE to start',
                    ox + areaSize / 2, oy + areaSize / 2 + 30,
                    { size: 28, align: 'center', baseline: 'middle' });
    } else if (this.gameOver) {
      c.fillStyle = 'rgba(0,0,0,0.6)';
      c.fillRect(ox, oy, areaSize, areaSize);
      const msg = this.lives <= 0 ? 'GAME OVER' : 'YOU WIN!';
      const color = this.lives <= 0 ? '#fff' : '#00ff00';
      this.drawText(msg, ox + areaSize / 2, oy + areaSize / 2 - 20,
                    { size: 40, color, align: 'center', baseline: 'middle' });
      this.drawText('Press R to restart',
                    ox + areaSize / 2, oy + areaSize / 2 + 30,
                    { size: 22, align: 'center', baseline: 'middle' });
    }
  }
}
