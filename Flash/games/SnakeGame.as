package games {
    import BaseGame;
    import flash.events.KeyboardEvent;
    import flash.ui.Keyboard;

    public class SnakeGame extends BaseGame {
        private var areaX:Number = 125, areaY:Number = 48;
        private var areaW:Number = 600, areaH:Number = 600;
        private var cell:Number = 20;
        private var cols:int, rows:int;

        private var snake:Array;
        private var dirX:int, dirY:int;
        private var nextDirX:int, nextDirY:int;
        private var foodX:int, foodY:int;
        private var score:int;
        private var gameOver:Boolean;
        private var moveTimer:Number;
        private var moveDelay:Number = 150;

        public function SnakeGame() {
            cols = int(areaW / cell);
            rows = int(areaH / cell);
        }

        override public function reset():void {
            snake = [{ x: int(cols/2), y: int(rows/2) }];
            dirX = 1; dirY = 0;
            nextDirX = 1; nextDirY = 0;
            spawnFood();
            score = 0;
            gameOver = false;
            moveTimer = 0;
        }

        private function spawnFood():void {
            while (true) {
                foodX = int(Math.random() * cols);
                foodY = int(Math.random() * rows);
                var hit:Boolean = false;
                for (var i:int = 0; i < snake.length; i++) {
                    if (snake[i].x == foodX && snake[i].y == foodY) { hit = true; break; }
                }
                if (!hit) return;
            }
        }

        override public function onKeyDown(e:KeyboardEvent):void {
            switch (e.keyCode) {
                case Keyboard.UP:    if (dirY != 1)  { nextDirX = 0;  nextDirY = -1; } break;
                case Keyboard.DOWN:  if (dirY != -1) { nextDirX = 0;  nextDirY = 1;  } break;
                case Keyboard.LEFT:  if (dirX != 1)  { nextDirX = -1; nextDirY = 0;  } break;
                case Keyboard.RIGHT: if (dirX != -1) { nextDirX = 1;  nextDirY = 0;  } break;
                case Keyboard.R:     reset(); break;
            }
        }

        override public function update(dt:Number):void {
            if (gameOver) return;
            moveTimer += dt;
            if (moveTimer < moveDelay) return;
            moveTimer = 0;

            dirX = nextDirX; dirY = nextDirY;
            var head:Object = snake[0];
            var nhx:int = head.x + dirX;
            var nhy:int = head.y + dirY;

            if (nhx < 0 || nhx >= cols || nhy < 0 || nhy >= rows) { gameOver = true; return; }
            for (var i:int = 0; i < snake.length; i++) {
                if (snake[i].x == nhx && snake[i].y == nhy) { gameOver = true; return; }
            }

            snake.unshift({ x: nhx, y: nhy });
            if (nhx == foodX && nhy == foodY) {
                score += 10;
                spawnFood();
            } else {
                snake.pop();
            }
        }

        override public function draw():void {
            beginTextFrame();
            clear();

            // 区域背景
            graphics.beginFill(0x14141e);
            graphics.drawRect(areaX, areaY, areaW, areaH);
            graphics.endFill();

            // 蛇
            for (var i:int = 0; i < snake.length; i++) {
                var s:Object = snake[i];
                roundRect(areaX + s.x * cell + 1, areaY + s.y * cell + 1,
                          cell - 2, cell - 2, 4, i == 0 ? 0x3cdc3c : 0x28b428);
            }

            // 食物
            roundRect(areaX + foodX * cell + 1, areaY + foodY * cell + 1,
                      cell - 2, cell - 2, 6, 0xff3c3c);

            // 分数
            drawText("Score: " + score, areaX + 10, areaY + 10, 28);

            if (gameOver) {
                graphics.beginFill(0x000000, 0.75);
                graphics.drawRect(areaX, areaY, areaW, areaH);
                graphics.endFill();
                drawText("GAME OVER - Press R",
                         areaX + areaW / 2, areaY + areaH / 2,
                         42, 0xffffff, "center");
            }
        }
    }
}