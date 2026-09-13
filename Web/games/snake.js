import { BaseGame } from './BaseGame.js';

export class SnakeGame extends BaseGame {
  constructor(canvas, ctx) {
    super(canvas, ctx);
    // 画布中央 600x600 区域
    this.area = { x: 125, y: 25, w: 600, h: 600 };
    this.cell = 20;
    this.cols = this.area.w / this.cell;
    this.rows = this.area.h / this.cell;
  }

  reset() {
    this.snake = [{ x: Math.floor(this.cols / 2), y: Math.floor(this.rows / 2) }];
    this.dir = { x: 1, y: 0 };
    this.nextDir = { x: 1, y: 0 };
    this.food = this.spawnFood();
    this.score = 0;
    this.gameOver = false;
    this.moveTimer = 0;
    this.moveDelay = 150;
  }

  spawnFood() {
    while (true) {
      const p = {
        x: Math.floor(Math.random() * this.cols),
        y: Math.floor(Math.random() * this.rows),
      };
      if (!this.snake.some(s => s.x === p.x && s.y === p.y)) return p;
    }
  }

  onKeyDown(e) {
    const k = e.key;
    if (k === 'ArrowUp'    && this.dir.y !== 1)  this.nextDir = { x: 0, y: -1 };
    if (k === 'ArrowDown'  && this.dir.y !== -1) this.nextDir = { x: 0, y: 1 };
    if (k === 'ArrowLeft'  && this.dir.x !== 1)  this.nextDir = { x: -1, y: 0 };
    if (k === 'ArrowRight' && this.dir.x !== -1) this.nextDir = { x: 1, y: 0 };
    if (k.toLowerCase() === 'r') this.reset();
  }

  update(dt) {
    if (this.gameOver) return;
    this.moveTimer += dt;
    if (this.moveTimer < this.moveDelay) return;
    this.moveTimer = 0;

    this.dir = { ...this.nextDir };
    const head = this.snake[0];
    const nh = { x: head.x + this.dir.x, y: head.y + this.dir.y };

    if (nh.x < 0 || nh.x >= this.cols || nh.y < 0 || nh.y >= this.rows) {
      this.gameOver = true; return;
    }
    if (this.snake.some(s => s.x === nh.x && s.y === nh.y)) {
      this.gameOver = true; return;
    }

    this.snake.unshift(nh);
    if (nh.x === this.food.x && nh.y === this.food.y) {
      this.score += 10;
      this.food = this.spawnFood();
    } else {
      this.snake.pop();
    }
  }

  draw() {
    this.clear();
    const { x: ax, y: ay } = this.area;
    const c = this.ctx;

    // 背景
    c.fillStyle = '#14141e';
    c.fillRect(ax, ay, this.area.w, this.area.h);

    // 蛇
    this.snake.forEach((s, i) => {
      c.fillStyle = i === 0 ? '#3cdc3c' : '#28b428';
      this.roundRect(ax + s.x * this.cell + 1, ay + s.y * this.cell + 1,
                     this.cell - 2, this.cell - 2, 4);
    });

    // 食物
    c.fillStyle = '#ff3c3c';
    this.roundRect(ax + this.food.x * this.cell + 1, ay + this.food.y * this.cell + 1,
                   this.cell - 2, this.cell - 2, 6);

    // 分数
    this.drawText(`Score: ${this.score}`, ax + 10, ay + 10, { size: 24 });

    if (this.gameOver) {
      c.fillStyle = 'rgba(0,0,0,0.75)';
      c.fillRect(ax, ay, this.area.w, this.area.h);
      this.drawText('GAME OVER - Press R', ax + this.area.w / 2, ay + this.area.h / 2,
                    { size: 36, align: 'center', baseline: 'middle' });
    }
  }
}