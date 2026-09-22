package games {
    import BaseGame;
    import flash.events.KeyboardEvent;
    import flash.ui.Keyboard;

    public class TankGame extends BaseGame {
        private var areaSize:Number = 600;
        private var offsetX:Number = 125, offsetY:Number = 48;
        private var cell:Number = 30;
        private var gridSize:int = 20;

        private var walls:Object;
        private var player:Object;
        private var enemies:Array;
        private var bullets:Array;
        private var score:int;
        private var gameOver:Boolean;
        private var enemySpawnTimer:int = 0;
        private var enemySpawnDelay:int = 120;
        private var playerMoveTimer:Number = 0;
        private var moveDelay:Number = 200;

        public function TankGame() {}

        override public function reset():void {
            walls = generateWalls();
            player = { x: 1, y: 1, dirX: 1, dirY: 0, alive: true };
            enemies = [{ x: gridSize - 2, y: gridSize - 2, dirX: -1, dirY: 0,
                         alive: true, timer: 0, moveTimer: 0 }];
            bullets = [];
            score = 0;
            gameOver = false;
            enemySpawnTimer = 0;
            playerMoveTimer = 0;
        }

        private function wallKey(x:int, y:int):String { return x + "," + y; }

        private function generateWalls():Object {
            var w:Object = {};
            for (var i:int = 0; i < gridSize; i++) {
                w[wallKey(i, 0)] = true;
                w[wallKey(i, gridSize - 1)] = true;
                w[wallKey(0, i)] = true;
                w[wallKey(gridSize - 1, i)] = true;
            }
            var k:int = 0;
            while (k < 12) {
                var x:int = 2 + int(Math.random() * (gridSize - 4));
                var y:int = 2 + int(Math.random() * (gridSize - 4));
                if ((x == 1 && y == 1) || (x == gridSize - 2 && y == gridSize - 2)) continue;
                w[wallKey(x, y)] = true;
                k++;
            }
            return w;
        }

        private function canMove(x:int, y:int, isPlayer:Boolean, self:Object):Boolean {
            if (x < 0 || x >= gridSize || y < 0 || y >= gridSize) return false;
            if (walls[wallKey(x, y)]) return false;
            if (isPlayer) {
                for (var i:int = 0; i < enemies.length; i++) {
                    var e:Object = enemies[i];
                    if (e.alive && e.x == x && e.y == y) return false;
                }
            } else {
                if (player.alive && player.x == x && player.y == y) return false;
                for (var j:int = 0; j < enemies.length; j++) {
                    var e2:Object = enemies[j];
                    if (e2 != self && e2.alive && e2.x == x && e2.y == y) return false;
                }
            }
            return true;
        }

        private function shoot(x:int, y:int, dx:int, dy:int, isPlayer:Boolean):void {
            bullets.push({ x: x + dx, y: y + dy, dx: dx, dy: dy,
                           isPlayer: isPlayer, alive: true });
        }

        override public function onKeyDown(e:KeyboardEvent):void {
            if (e.keyCode == Keyboard.R) { reset(); return; }
            if (e.keyCode == Keyboard.SPACE && !gameOver && player.alive) {
                shoot(player.x, player.y, player.dirX, player.dirY, true);
            }
        }

        override public function update(dt:Number):void {
            if (gameOver) return;

            // 玩家移动
            if (player.alive) {
                var dx:int = 0, dy:int = 0;
                if (keys[87]) dy = -1;        // W
                else if (keys[83]) dy = 1;    // S
                else if (keys[65]) dx = -1;   // A
                else if (keys[68]) dx = 1;    // D

                if (dx != 0 || dy != 0) {
                    player.dirX = dx; player.dirY = dy;
                    playerMoveTimer += dt;
                    if (playerMoveTimer >= moveDelay) {
                        playerMoveTimer = 0;
                        var nx:int = player.x + dx;
                        var ny:int = player.y + dy;
                        if (canMove(nx, ny, true, player)) {
                            player.x = nx; player.y = ny;
                        }
                    }
                } else {
                    playerMoveTimer = 0;
                }
            }

            // 子弹
            for (var i:int = 0; i < bullets.length; i++) {
                var b:Object = bullets[i];
                if (!b.alive) continue;
                var bx:Number = b.x, by:Number = b.y;
                var steps:int = 5;
                for (var s:int = 0; s < steps; s++) {
                    bx += b.dx / steps;
                    by += b.dy / steps;
                    var ix:int = Math.round(bx);
                    var iy:int = Math.round(by);
                    if (ix < 0 || ix >= gridSize || iy < 0 || iy >= gridSize) { b.alive = false; break; }
                    if (walls[wallKey(ix, iy)]) { b.alive = false; break; }
                    if (b.isPlayer) {
                        var hit:Boolean = false;
                        for (var ei:int = 0; ei < enemies.length; ei++) {
                            var en:Object = enemies[ei];
                            if (en.alive && en.x == ix && en.y == iy) {
                                en.alive = false;
                                b.alive = false;
                                score += 10;
                                hit = true;
                                break;
                            }
                        }
                        if (hit) break;
                    } else {
                        if (player.alive && player.x == ix && player.y == iy) {
                            player.alive = false;
                            b.alive = false;
                            gameOver = true;
                            break;
                        }
                    }
                }
                b.x = bx; b.y = by;
            }

            // 清掉死子弹
            var live:Array = [];
            for (var k:int = 0; k < bullets.length; k++) if (bullets[k].alive) live.push(bullets[k]);
            bullets = live;

            // 敌人
            for (var e2i:int = 0; e2i < enemies.length; e2i++) {
                var e:Object = enemies[e2i];
                if (!e.alive) continue;
                e.timer++;
                if (e.timer >= 20) {
                    e.timer = 0;
                    if (Math.random() < 0.25) {
                        var dirs:Array = [[1,0],[-1,0],[0,1],[0,-1]];
                        var d:Array = dirs[int(Math.random() * 4)];
                        e.dirX = d[0]; e.dirY = d[1];
                    }
                    e.moveTimer += 20;
                    if (e.moveTimer >= 250) {
                        e.moveTimer = 0;
                        var ex:int = e.x + e.dirX;
                        var ey:int = e.y + e.dirY;
                        if (canMove(ex, ey, false, e)) { e.x = ex; e.y = ey; }
                    }
                }
                if (Math.random() < 0.015) {
                    shoot(e.x, e.y, e.dirX, e.dirY, false);
                }
            }

            // 生成敌人
            enemySpawnTimer++;
            if (enemySpawnTimer >= enemySpawnDelay) {
                enemySpawnTimer = 0;
                var aliveCount:int = 0;
                for (var ai:int = 0; ai < enemies.length; ai++) if (enemies[ai].alive) aliveCount++;
                if (aliveCount < 3) {
                    enemies.push({ x: gridSize - 2, y: gridSize - 2,
                                   dirX: -1, dirY: 0, alive: true,
                                   timer: 0, moveTimer: 0 });
                }
            }
        }

        private function drawTank(x:int, y:int, body:uint, gun:uint, dx:int, dy:int):void {
            var cx:Number = offsetX + x * cell + cell / 2;
            var cy:Number = offsetY + y * cell + cell / 2;

            // 车身
            roundRect(offsetX + x * cell + 2, offsetY + y * cell + 2,
                      cell - 4, cell - 4, 8, body);

            // 炮管
            graphics.beginFill(gun);
            var t:Number = 4;
            if (dx == 1)       graphics.drawRect(cx + t, cy - 4, 14, 8);
            else if (dx == -1) graphics.drawRect(cx - 18, cy - 4, 14, 8);
            else if (dy == -1) graphics.drawRect(cx - 4, cy - 18, 8, 14);
            else               graphics.drawRect(cx - 4, cy + t, 8, 14);
            graphics.endFill();

            // 圆顶
            graphics.beginFill(body);
            graphics.drawCircle(cx, cy, 6);
            graphics.endFill();
        }

        override public function draw():void {
            beginTextFrame();
            clear();

            // 区域背景
            graphics.beginFill(0x1e1e2a);
            graphics.drawRect(offsetX, offsetY, areaSize, areaSize);
            graphics.endFill();

            // 墙
            for (var key:String in walls) {
                var parts:Array = key.split(",");
                var wx:int = int(parts[0]);
                var wy:int = int(parts[1]);
                roundRect(offsetX + wx * cell + 1, offsetY + wy * cell + 1,
                          cell - 2, cell - 2, 8, 0x646482);
            }

            // 坦克
            if (player.alive) drawTank(player.x, player.y, 0x32e632, 0x1ec81e, player.dirX, player.dirY);
            for (var i:int = 0; i < enemies.length; i++) {
                var e:Object = enemies[i];
                if (e.alive) drawTank(e.x, e.y, 0xe63232, 0xc81e1e, e.dirX, e.dirY);
            }

            // 子弹
            for (var b:int = 0; b < bullets.length; b++) {
                var bl:Object = bullets[b];
                graphics.beginFill(bl.isPlayer ? 0xffff50 : 0xff9632);
                graphics.drawCircle(offsetX + bl.x * cell + cell / 2,
                                    offsetY + bl.y * cell + cell / 2, 5);
                graphics.endFill();
            }

            // 分数
            drawText("Score: " + score, offsetX + 10, offsetY + 10, 25);

            if (gameOver) {
                graphics.beginFill(0x000000, 0.6);
                graphics.drawRect(offsetX, offsetY, areaSize, areaSize);
                graphics.endFill();
                drawText("GAME OVER",
                         offsetX + areaSize / 2, offsetY + areaSize / 2 - 20,
                         46, 0xffffff, "center");
                drawText("Press R to restart",
                         offsetX + areaSize / 2, offsetY + areaSize / 2 + 30,
                         26, 0xffffff, "center");
            }
        }
    }
}