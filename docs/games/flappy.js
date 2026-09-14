import { BaseGame } from './BaseGame.js';

export class FlappyGame extends BaseGame {
  constructor(canvas, ctx) {
    super(canvas, ctx);
    this.areaW = 400;
    this.areaH = 600;
    this.offsetX = Math.floor((this.width - this.areaW) / 2);
    this.offsetY = Math.floor((this.height - this.areaH) / 2);
  }

  reset() {
    this.birdY = this.areaH / 2;
    this.birdVy = 0;
    this.birdR = 15;
    this.gravity = 0.6;
    this.jump = -9.0;
    this.pipes = [];
    this.pipeW = 60;
    this.pipeGap = 170;
    this.pipeSpeed = 3.5;
    this.score = 0;
    this.gameOver = false;
    this.pipeTimer = 0;
    this.pipeDelay = 100;
    this.birdX = 80;

    this.started = false;
  }

  addPipe() {
    const gapY = 80 + Math.random() * (this.areaH - 160 - this.pipeGap);
    this.pipes.push({ x: this.areaW, gapY, passed: false });
  }

  onKeyDown(e) {
    if (e.key.toLowerCase() === 'r') { this.reset(); return; }
    if (e.key === ' ' || e.key === 'ArrowUp') {
      if (this.gameOver) return;
      if (!this.started) {
        this.started = true;
      }
      this.birdVy = this.jump;
    }
  }

  onClick() {
    if (this.gameOver) return;
    if (!this.started) this.started = true;
    this.birdVy = this.jump;
  }

  update(dt) {
    if (this.gameOver || !this.started) return;

    const f = dt / 16.67;
    this.birdVy += this.gravity * f;
    this.birdY += this.birdVy * f;

    if (this.birdY - this.birdR < 0) { this.birdY = this.birdR; this.birdVy = 0; }
    if (this.birdY + this.birdR > this.areaH) { this.gameOver = true; return; }

    this.pipeTimer += dt;
    if (this.pipeTimer >= this.pipeDelay * 16.67) {
      this.pipeTimer = 0;
      this.addPipe();
    }

    for (let i = this.pipes.length - 1; i >= 0; i--) {
      const p = this.pipes[i];
      p.x -= this.pipeSpeed * f;

      if (!p.passed && p.x + this.pipeW < this.birdX) {
        p.passed = true;
        this.score++;
      }

      if (p.x < this.birdX + this.birdR && p.x + this.pipeW > this.birdX - this.birdR) {
        if (this.birdY - this.birdR < p.gapY ||
            this.birdY + this.birdR > p.gapY + this.pipeGap) {
          this.gameOver = true;
          return;
        }
      }

      if (p.x + this.pipeW < 0) this.pipes.splice(i, 1);
    }
  }

  draw() {
    this.clear();
    const c = this.ctx;
    const { offsetX: ox, offsetY: oy, areaW, areaH } = this;

    c.fillStyle = '#87ceeb';
    c.fillRect(ox, oy, areaW, areaH);

    for (const p of this.pipes) {
      c.fillStyle = '#28c828';
      this.roundRect(ox + p.x, oy, this.pipeW, p.gapY, 6);
      this.roundRect(ox + p.x - 6, oy + p.gapY - 22, this.pipeW + 12, 22, 6);
      const by = p.gapY + this.pipeGap;
      this.roundRect(ox + p.x, oy + by, this.pipeW, areaH - by, 6);
      this.roundRect(ox + p.x - 6, oy + by, this.pipeW + 12, 22, 6);
    }

    c.fillStyle = '#ffff32';
    c.beginPath();
    c.arc(ox + this.birdX, oy + this.birdY, this.birdR, 0, Math.PI * 2);
    c.fill();
    c.fillStyle = '#000';
    c.beginPath();
    c.arc(ox + this.birdX + 8, oy + this.birdY - 5, 4, 0, Math.PI * 2);
    c.fill();
    c.fillStyle = '#fff';
    c.beginPath();
    c.arc(ox + this.birdX + 10, oy + this.birdY - 7, 2, 0, Math.PI * 2);
    c.fill();

    this.drawText(String(this.score), ox + areaW / 2, oy + 50,
                  { size: 40, align: 'center', baseline: 'middle' });

    if (!this.started && !this.gameOver) {
      c.fillStyle = 'rgba(0,0,0,0.4)';
      c.fillRect(ox, oy, areaW, areaH);
      this.drawText('Press SPACE / Click to start',
                    ox + areaW / 2, oy + areaH / 2,
                    { size: 22, align: 'center', baseline: 'middle' });
    }

    if (this.gameOver) {
      c.fillStyle = 'rgba(0,0,0,0.6)';
      c.fillRect(ox, oy, areaW, areaH);
      this.drawText('GAME OVER', ox + areaW / 2, oy + areaH / 2 - 30,
                    { size: 40, align: 'center', baseline: 'middle' });
      this.drawText('Press R to restart', ox + areaW / 2, oy + areaH / 2 + 20,
                    { size: 22, align: 'center', baseline: 'middle' });
    }
  }
}
