export class BaseGame {
  constructor(canvas, ctx) {
    this.canvas = canvas;
    this.ctx = ctx;
    this.width = canvas.width;
    this.height = canvas.height;
    this.keys = {};
    this.running = false;
    this._raf = null;
    this._lastTime = 0;
    this.onExitToMenu = null;

    this._onKeyDown = (e) => {
      if (e.key === 'Escape') {
        if (this.onExitToMenu) this.onExitToMenu();
        return;
      }
      this.keys[e.key.toLowerCase()] = true;
      this.onKeyDown?.(e);
    };
    this._onKeyUp = (e) => {
      this.keys[e.key.toLowerCase()] = false;
      this.onKeyUp?.(e);
    };
    this._onClick = (e) => this.onClick?.(e);
    this._onContextMenu = (e) => {
      e.preventDefault();
      this.onRightClick?.(e);
    };
    this._onMouseMove = (e) => this.onMouseMove?.(e);

    window.addEventListener('keydown', this._onKeyDown);
    window.addEventListener('keyup', this._onKeyUp);
    canvas.addEventListener('click', this._onClick);
    canvas.addEventListener('contextmenu', this._onContextMenu);
    canvas.addEventListener('mousemove', this._onMouseMove);
  }

  reset() {}

  start() {
    this.reset();
    this.running = true;
    this._lastTime = performance.now();
    const loop = (t) => {
      if (!this.running) return;
      const dt = t - this._lastTime;
      this._lastTime = t;
      this.update(dt);
      this.draw();
      this._raf = requestAnimationFrame(loop);
    };
    this._raf = requestAnimationFrame(loop);
  }

  destroy() {
    this.running = false;
    if (this._raf) cancelAnimationFrame(this._raf);
    window.removeEventListener('keydown', this._onKeyDown);
    window.removeEventListener('keyup', this._onKeyUp);
    this.canvas.removeEventListener('click', this._onClick);
    this.canvas.removeEventListener('contextmenu', this._onContextMenu);
    this.canvas.removeEventListener('mousemove', this._onMouseMove);
  }

  update(_dt) {}
  draw() {}

  drawText(text, x, y, { size = 24, color = '#fff', align = 'left', baseline = 'top' } = {}) {
    const c = this.ctx;
    c.font = `${size}px 'Segoe UI', Tahoma, sans-serif`;
    c.fillStyle = color;
    c.textAlign = align;
    c.textBaseline = baseline;
    c.fillText(text, x, y);
  }

  roundRect(x, y, w, h, r, fill, stroke) {
    const c = this.ctx;
    if (w < 2 * r) r = w / 2;
    if (h < 2 * r) r = h / 2;
    c.beginPath();
    c.moveTo(x + r, y);
    c.arcTo(x + w, y, x + w, y + h, r);
    c.arcTo(x + w, y + h, x, y + h, r);
    c.arcTo(x, y + h, x, y, r);
    c.arcTo(x, y, x + w, y, r);
    c.closePath();

    if (fill !== undefined && fill !== null) c.fillStyle = fill;
    c.fill();

    if (stroke) {
      c.strokeStyle = stroke;
      c.lineWidth = 2;
      c.stroke();
    }
  }

  clear(bg = '#191923') {
    this.ctx.fillStyle = bg;
    this.ctx.fillRect(0, 0, this.width, this.height);
  }
}