package games {
    import BaseGame;
    import flash.events.KeyboardEvent;
    import flash.events.MouseEvent;
    import flash.ui.Keyboard;

    public class FlappyGame extends BaseGame {
        private var areaW:Number = 400, areaH:Number = 600;
        private var offsetX:Number = 225, offsetY:Number = 48;

        private var birdY:Number, birdVy:Number;
        private var birdR:Number = 15;
        private var gravity:Number = 0.6;
        private var jump:Number = -10.0;
        private var birdX:Number = 80;

        private var pipes:Array;
        private var pipeW:Number = 60;
        private var pipeGap:Number = 170;
        private var pipeSpeed:Number = 3.5;
        private var score:int;
        private var gameOver:Boolean;
        private var started:Boolean;
        private var pipeTimer:Number = 0;
        private var pipeDelay:Number = 1500;

        public function FlappyGame() {}

        override public function reset():void {
            birdY = areaH / 2;
            birdVy = 0;
            pipes = [];
            score = 0;
            gameOver = false;
            started = false;
            pipeTimer = 0;
        }

        private function addPipe():void {
            var gapY:Number = 80 + Math.random() * (areaH - 160 - pipeGap);
            pipes.push({ x: areaW, gapY: gapY, passed: false });
        }

        private function flap():void {
            if (gameOver) return;
            if (!started) started = true;
            birdVy = jump;
        }

        override public function onKeyDown(e:KeyboardEvent):void {
            if (e.keyCode == Keyboard.R) { reset(); return; }
            if (e.keyCode == Keyboard.SPACE || e.keyCode == Keyboard.UP) flap();
        }

        override public function onClick(e:MouseEvent):void { flap(); }

        override public function update(dt:Number):void {
            if (gameOver || !started) return;
            var f:Number = dt / 16.67;
            birdVy += gravity * f;
            birdY += birdVy * f;

            if (birdY - birdR < 0) { birdY = birdR; birdVy = 0; }
            if (birdY + birdR > areaH) { gameOver = true; return; }

            pipeTimer += dt;
            if (pipeTimer >= pipeDelay) {
                pipeTimer = 0;
                addPipe();
            }

            for (var i:int = pipes.length - 1; i >= 0; i--) {
                var p:Object = pipes[i];
                p.x -= pipeSpeed * f;

                if (!p.passed && p.x + pipeW < birdX) {
                    p.passed = true;
                    score++;
                }

                if (p.x < birdX + birdR && p.x + pipeW > birdX - birdR) {
                    if (birdY - birdR < p.gapY || birdY + birdR > p.gapY + pipeGap) {
                        gameOver = true;
                        return;
                    }
                }

                if (p.x + pipeW < 0) pipes.splice(i, 1);
            }
        }

        override public function draw():void {
            beginTextFrame();
            clear();

            // 天空
            graphics.beginFill(0x87ceeb);
            graphics.drawRect(offsetX, offsetY, areaW, areaH);
            graphics.endFill();

            // 管道（每根独立 beginFill）
            for (var i:int = 0; i < pipes.length; i++) {
                var p:Object = pipes[i];
                var px:Number = offsetX + p.x;
                var py:Number = offsetY;

                // 上管道主体
                graphics.beginFill(0x28c828);
                graphics.drawRect(px, py, pipeW, p.gapY);
                graphics.endFill();
                // 上管帽（独立）
                graphics.beginFill(0x28c828);
                graphics.drawRoundRect(px - 6, py + p.gapY - 22,
                                       pipeW + 12, 22, 8, 8);
                graphics.endFill();

                // 下管道主体
                var by:Number = p.gapY + pipeGap;
                graphics.beginFill(0x28c828);
                graphics.drawRect(px, py + by, pipeW, areaH - by);
                graphics.endFill();
                // 下管帽（独立）
                graphics.beginFill(0x28c828);
                graphics.drawRoundRect(px - 6, py + by, pipeW + 12, 22, 8, 8);
                graphics.endFill();
            }

            // 鸟
            graphics.beginFill(0xffff32);
            graphics.drawCircle(offsetX + birdX, offsetY + birdY, birdR);
            graphics.endFill();
            graphics.beginFill(0x000000);
            graphics.drawCircle(offsetX + birdX + 8, offsetY + birdY - 5, 4);
            graphics.endFill();
            graphics.beginFill(0xffffff);
            graphics.drawCircle(offsetX + birdX + 10, offsetY + birdY - 7, 2);
            graphics.endFill();

            // 分数
            drawText(String(score), offsetX + areaW / 2, offsetY + 50,
                     46, 0xffffff, "center");

            if (!started && !gameOver) {
                graphics.beginFill(0x000000, 0.4);
                graphics.drawRect(offsetX, offsetY, areaW, areaH);
                graphics.endFill();
                drawText("Press SPACE / Click to start",
                         offsetX + areaW / 2, offsetY + areaH / 2,
                         24, 0xffffff, "center");
            }

            if (gameOver) {
                graphics.beginFill(0x000000, 0.6);
                graphics.drawRect(offsetX, offsetY, areaW, areaH);
                graphics.endFill();
                drawText("GAME OVER", offsetX + areaW / 2, offsetY + areaH / 2 - 30,
                         44, 0xffffff, "center");
                drawText("Press R to restart", offsetX + areaW / 2, offsetY + areaH / 2 + 20,
                         24, 0xffffff, "center");
            }
        }
    }
}