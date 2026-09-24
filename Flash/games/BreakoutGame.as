package games {
    import BaseGame;
    import flash.events.KeyboardEvent;
    import flash.ui.Keyboard;

    public class BreakoutGame extends BaseGame {
        private var areaSize:Number = 600;
        private var offsetX:Number = 125, offsetY:Number = 48;

        private var paddle:Object;
        private var paddleSpeed:Number = 7;
        private var ball:Object;
        private var ballSpeed:Number = 8;
        private var ballDx:Number = 1;
        private var ballDy:Number = -1;
        private var bricks:Array;
        private var score:int;
        private var lives:int;
        private var gameOver:Boolean;
        private var waiting:Boolean;

        public function BreakoutGame() {}

        override public function reset():void {
            paddle = { x: areaSize / 2 - 60, y: areaSize - 40, w: 120, h: 16 };
            ball = { x: areaSize / 2 - 10, y: areaSize - 70, w: 20, h: 20 };
            ballDx = 1; ballDy = -1;
            normalizeSpeed();
            bricks = [];
            score = 0;
            lives = 3;
            gameOver = false;
            waiting = true;

            var rows:int = 5, cols:int = 8;
            var bw:Number = Math.floor((areaSize - 20) / cols) - 4;
            var bh:Number = 22;
            var colors:Array = [0xe63232, 0xe69632, 0xe6e632, 0x32e632, 0x3296e6];
            for (var r:int = 0; r < rows; r++) {
                for (var c:int = 0; c < cols; c++) {
                    bricks.push({
                        x: 10 + c * (bw + 4),
                        y: 40 + r * (bh + 4),
                        w: bw, h: bh,
                        color: colors[r % colors.length],
                        alive: true
                    });
                }
            }
        }

        private function normalizeSpeed():void {
            var current:Number = Math.sqrt(ballDx * ballDx + ballDy * ballDy);
            if (current == 0) return;
            var scale:Number = ballSpeed / current;
            ballDx *= scale;
            ballDy *= scale;
        }

        private function resetBall():void {
            ball.x = areaSize / 2 - 10;
            ball.y = areaSize - 70;
            ballDx = Math.random() < 0.5 ? 1 : -1;
            ballDy = -1;
            normalizeSpeed();
        }

        override public function onKeyDown(e:KeyboardEvent):void {
            if (e.keyCode == Keyboard.R) { reset(); return; }
            if (e.keyCode == Keyboard.SPACE && waiting && !gameOver) waiting = false;
        }

        private function rectCollide(a:Object, b:Object):Boolean {
            return a.x < b.x + b.w && a.x + a.w > b.x &&
                   a.y < b.y + b.h && a.y + a.h > b.y;
        }

        override public function update(dt:Number):void {
            if (gameOver || waiting) return;
            var f:Number = dt / 16.67;

            if (keys[Keyboard.LEFT] && paddle.x > 0)
                paddle.x -= paddleSpeed * f;
            if (keys[Keyboard.RIGHT] && paddle.x + paddle.w < areaSize)
                paddle.x += paddleSpeed * f;

            ball.x += ballDx * f;
            ball.y += ballDy * f;

            if (ball.x <= 0 || ball.x + ball.w >= areaSize) ballDx *= -1;
            if (ball.y <= 0) ballDy *= -1;

            if (ball.y + ball.h >= areaSize) {
                lives--;
                if (lives <= 0) {
                    gameOver = true;
                } else {
                    waiting = true;
                    resetBall();
                }
                return;
            }

            if (rectCollide(ball, paddle)) {
                ballDy = -Math.abs(ballDy);
                var hit:Number = (ball.x + ball.w / 2 - (paddle.x + paddle.w / 2)) / (paddle.w / 2);
                ballDx = hit * ballSpeed * 0.9;
                if (Math.abs(ballDx) < 1.2) ballDx = ballDx >= 0 ? 1.8 : -1.8;
                normalizeSpeed();
            }

            for (var i:int = 0; i < bricks.length; i++) {
                var b:Object = bricks[i];
                if (!b.alive) continue;
                if (rectCollide(ball, b)) {
                    b.alive = false;
                    score += 10;

                    var oTop:Number = ball.y + ball.h - b.y;
                    var oBot:Number = b.y + b.h - ball.y;
                    var oLef:Number = ball.x + ball.w - b.x;
                    var oRig:Number = b.x + b.w - ball.x;
                    var minOv:Number = Math.min(Math.min(oTop, oBot), Math.min(oLef, oRig));

                    if (minOv == oTop) {
                        ball.y = b.y - ball.h; ballDy = -Math.abs(ballDy);
                    } else if (minOv == oBot) {
                        ball.y = b.y + b.h; ballDy = Math.abs(ballDy);
                    } else if (minOv == oLef) {
                        ball.x = b.x - ball.w; ballDx = -Math.abs(ballDx);
                    } else {
                        ball.x = b.x + b.w; ballDx = Math.abs(ballDx);
                    }
                    normalizeSpeed();
                    break;
                }
            }

            var allDead:Boolean = true;
            for (var k:int = 0; k < bricks.length; k++) {
                if (bricks[k].alive) { allDead = false; break; }
            }
            if (allDead) gameOver = true;
        }

        override public function draw():void {
            beginTextFrame();
            clear();

            // 区域背景
            graphics.beginFill(0x14141e);
            graphics.drawRect(offsetX, offsetY, areaSize, areaSize);
            graphics.endFill();

            // 砖块
            for (var i:int = 0; i < bricks.length; i++) {
                var b:Object = bricks[i];
                if (!b.alive) continue;
                roundRect(offsetX + b.x, offsetY + b.y, b.w, b.h, 8, b.color, 0xffffff);
            }

            // 挡板
            roundRect(offsetX + paddle.x, offsetY + paddle.y,
                      paddle.w, paddle.h, 10, 0xffffff, 0xc8c8dc);

            // 球
            graphics.beginFill(0xffff78);
            graphics.drawCircle(offsetX + ball.x + ball.w / 2,
                                offsetY + ball.y + ball.h / 2,
                                ball.w / 2);
            graphics.endFill();

            // HUD
            drawText("Score: " + score, offsetX + 10, offsetY + 10, 25);
            drawText("Lives: " + lives, offsetX + areaSize - 130, offsetY + 10, 25);

            if (waiting && !gameOver) {
                drawText("Press SPACE to start",
                         offsetX + areaSize / 2, offsetY + areaSize / 2 + 30,
                         32, 0xffffff, "center");
            } else if (gameOver) {
                graphics.beginFill(0x000000, 0.6);
                graphics.drawRect(offsetX, offsetY, areaSize, areaSize);
                graphics.endFill();
                var msg:String = lives <= 0 ? "GAME OVER" : "YOU WIN!";
                var col:uint = lives <= 0 ? 0xffffff : 0x00ff00;
                drawText(msg, offsetX + areaSize / 2, offsetY + areaSize / 2 - 20,
                         46, col, "center");
                drawText("Press R to restart",
                         offsetX + areaSize / 2, offsetY + areaSize / 2 + 30,
                         26, 0xffffff, "center");
            }
        }
    }
}
