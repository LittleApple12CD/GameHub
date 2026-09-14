import { BaseGame } from './BaseGame.js';

export class TankGame extends BaseGame {
  constructor(canvas, ctx) {
    super(canvas, ctx);
    this.areaSize = 600;
    this.offsetX = Math.floor((this.width - this.areaSize) / 2);
    this.offsetY = Math.floor((this.height - this.areaSize) / 2);
    this.cell = 30;
    this.gridSize = this.areaSize / this.cell; // 20
  }

  reset() {
    this.walls = this.generateWalls();
    this.player = { x: 1, y: 1, dir: { x: 1, y: 0 }, alive: true };
    this.enemies = [{
      x: this.gridSize - 2, y: this.gridSize - 2,
      dir: { x: -1, y: 0 }, alive: true, timer: 0,
    }];
    this.bullets = [];
    this.score = 0;
    this.gameOver = false;
    this.enemySpawnTimer = 0;
    this.enemySpawnDelay = 120;
    this.playerMoveTimer = 0;
    this.enemyMoveTimer = 0;
    this.moveDelay = 200;
  }

  generateWalls() {
    const walls = new Set();
    const key = (x, y) => `${x},${y}`;
    for (let i = 0; i < this.gridSize; i++) {
      walls.add(key(i, 0));
      walls.add(key(i, this.gridSize - 1));
      walls.add(key(0, i));
      walls.add(key(this.gridSize - 1, i));
    }
    for (let k = 0; k < 12; k++) {
      const x = 2 + Math.floor(Math.random() * (this.gridSize - 4));
      const y = 2 + Math.floor(Math.random() * (this.gridSize - 4));
      if ((x === 1 && y === 1) || (x === this.gridSize - 2 && y === this.gridSize - 2)) continue;
      walls.add(key(x, y));
    }
    return walls;
  }

  canMove(x, y, isPlayer, self) {
    if (x < 0 || x >= this.gridSize || y < 0 || y >= this.gridSize) return false;
    if (this.walls.has(`${x},${y}`)) return false;
    if (isPlayer) {
      if (this.enemies.some(e => e.alive && e.x === x && e.y === y)) return false;
    } else {
      if (this.player.alive && this.player.x === x && this.player.y === y) return false;
      if (this.enemies.some(e => e !== self && e.alive && e.x === x && e.y === y)) return false;
    }
    return true;
  }

  shoot(x, y, dir, isPlayer) {
    this.bullets.push({
      x: x + dir.x, y: y + dir.y, dir, isPlayer, alive: true,
    });
  }

  onKeyDown(e) {
    if (e.key.toLowerCase() === 'r') { this.reset(); return; }
    if (e.key === ' ' && !this.gameOver && this.player.alive) {
      this.shoot(this.player.x, this.player.y, this.player.dir, true);
    }
  }

  update(dt) {
    if (this.gameOver) return;

    // 玩家移动
    if (this.player.alive) {
      let dx = 0, dy = 0;
      if (this.keys['w']) dy = -1;
      else if (this.keys['s']) dy = 1;
      else if (this.keys['a']) dx = -1;
      else if (this.keys['d']) dx = 1;

      if (dx || dy) {
        this.player.dir = { x: dx, y: dy };
        this.playerMoveTimer += dt;
        if (this.playerMoveTimer >= this.moveDelay) {
          this.playerMoveTimer = 0;
          const nx = this.player.x + dx, ny = this.player.y + dy;
          if (this.canMove(nx, ny, true, this.player)) {
            this.player.x = nx;
            this.player.y = ny;
          }
        }
      } else {
        this.playerMoveTimer = 0;
      }
    }

    // 子弹移动
    for (const b of this.bullets) {
      if (!b.alive) continue;
      let bx = b.x, by = b.y;
      const steps = 5;
      for (let s = 0; s < steps; s++) {
        bx += b.dir.x / steps;
        by += b.dir.y / steps;
        const ix = Math.round(bx), iy = Math.round(by);
        if (ix < 0 || ix >= this.gridSize || iy < 0 || iy >= this.gridSize) { b.alive = false; break; }
        if (this.walls.has(`${ix},${iy}`)) { b.alive = false; break; }
        if (b.isPlayer) {
          const enemy = this.enemies.find(e => e.alive && e.x === ix && e.y === iy);
          if (enemy) { enemy.alive = false; b.alive = false; this.score += 10; break; }
        } else {
          if (this.player.alive && this.player.x === ix && this.player.y === iy) {
            this.player.alive = false;
            b.alive = false;
            this.gameOver = true;
            break;
          }
        }
      }
      b.x = bx; b.y = by;
    }
    this.bullets = this.bullets.filter(b => b.alive);

    // 敌人
    for (const e of this.enemies) {
      if (!e.alive) continue;
      e.timer++;
      if (e.timer >= 20) {
        e.timer = 0;
        if (Math.random() < 0.25) {
          const dirs = [{x:1,y:0},{x:-1,y:0},{x:0,y:1},{x:0,y:-1}];
          e.dir = dirs[Math.floor(Math.random() * 4)];
        }
        e.moveTimer = (e.moveTimer || 0) + 20;
        if (e.moveTimer >= 250) {
          e.moveTimer = 0;
          const nx = e.x + e.dir.x, ny = e.y + e.dir.y;
          if (this.canMove(nx, ny, false, e)) { e.x = nx; e.y = ny; }
        }
      }
      if (Math.random() < 0.015) this.shoot(e.x, e.y, e.dir, false);
    }

    // 生成敌人
    this.enemySpawnTimer++;
    if (this.enemySpawnTimer >= this.enemySpawnDelay) {
      this.enemySpawnTimer = 0;
      if (this.enemies.filter(e => e.alive).length < 3) {
        this.enemies.push({
          x: this.gridSize - 2, y: this.gridSize - 2,
          dir: { x: -1, y: 0 }, alive: true, timer: 0, moveTimer: 0,
        });
      }
    }
  }

  drawTank(x, y, body, gun, dir) {
    const c = this.ctx;
    const { offsetX: ox, offsetY: oy, cell } = this;
    const cx = ox + x * cell + cell / 2;
    const cy = oy + y * cell + cell / 2;
    c.fillStyle = body;
    this.roundRect(ox + x * cell + 2, oy + y * cell + 2, cell - 4, cell - 4, 8);
    c.fillStyle = gun;
    const t = 4;
    if (dir.x === 1)      c.fillRect(cx + t, cy - 4, 14, 8);
    else if (dir.x === -1) c.fillRect(cx - 18, cy - 4, 14, 8);
    else if (dir.y === -1) c.fillRect(cx - 4, cy - 18, 8, 14);
    else                   c.fillRect(cx - 4, cy + t, 8, 14);
    c.fillStyle = body;
    c.beginPath();
    c.arc(cx, cy, 6, 0, Math.PI * 2);
    c.fill();
  }

  draw() {
    this.clear();
    const c = this.ctx;
    const { offsetX: ox, offsetY: oy, cell } = this;

    c.fillStyle = '#1e1e2a';
    c.fillRect(ox, oy, this.areaSize, this.areaSize);

    // 墙
    c.fillStyle = '#646482';
    for (const w of this.walls) {
      const [x, y] = w.split(',').map(Number);
      this.roundRect(ox + x * cell + 1, oy + y * cell + 1, cell - 2, cell - 2, 8, '#646482');
    }

    // 坦克
    if (this.player.alive) this.drawTank(this.player.x, this.player.y, '#32e632', '#1ec81e', this.player.dir);
    for (const e of this.enemies) if (e.alive) this.drawTank(e.x, e.y, '#e63232', '#c81e1e', e.dir);

    // 子弹
    for (const b of this.bullets) {
      c.fillStyle = b.isPlayer ? '#ffff50' : '#ff9632';
      c.beginPath();
      c.arc(ox + b.x * cell + cell / 2, oy + b.y * cell + cell / 2, 5, 0, Math.PI * 2);
      c.fill();
    }

    // 分数
    this.drawText(`Score: ${this.score}`, ox + 10, oy + 10, { size: 22 });

    if (this.gameOver) {
      c.fillStyle = 'rgba(0,0,0,0.6)';
      c.fillRect(ox, oy, this.areaSize, this.areaSize);
      this.drawText('GAME OVER', ox + this.areaSize / 2, oy + this.areaSize / 2 - 20,
                    { size: 40, align: 'center', baseline: 'middle' });
      this.drawText('Press R to restart', ox + this.areaSize / 2, oy + this.areaSize / 2 + 30,
                    { size: 22, align: 'center', baseline: 'middle' });
    }
  }
}