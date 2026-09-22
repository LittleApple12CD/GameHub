package games {
    import BaseGame;
    import flash.events.KeyboardEvent;
    import flash.events.MouseEvent;
    import flash.ui.Keyboard;

    public class MinesweeperGame extends BaseGame {
        private var rows:int = 16, cols:int = 16;
        private var mineCount:int = 40;
        private var cell:Number = 30;
        private var boardW:Number = 480;
        private var boardH:Number = 480;
        private var offsetX:Number = 185;
        private var offsetY:Number = 105;

        private var board:Array;
        private var revealed:Array;
        private var flagged:Array;
        private var gameOver:Boolean;
        private var won:Boolean;
        private var firstClick:Boolean;

        public function MinesweeperGame() {}

        override public function reset():void {
            board = []; revealed = []; flagged = [];
            for (var y:int = 0; y < rows; y++) {
                var r1:Array = [], r2:Array = [], r3:Array = [];
                for (var x:int = 0; x < cols; x++) { r1.push(0); r2.push(false); r3.push(false); }
                board.push(r1); revealed.push(r2); flagged.push(r3);
            }
            gameOver = false;
            won = false;
            firstClick = true;
        }

        private function placeMines(sx:int, sy:int):void {
            var placed:int = 0;
            while (placed < mineCount) {
                var x:int = int(Math.random() * cols);
                var y:int = int(Math.random() * rows);
                if (board[y][x] == -1) continue;
                if (Math.abs(x - sx) <= 1 && Math.abs(y - sy) <= 1) continue;
                board[y][x] = -1;
                placed++;
            }
            for (var yy:int = 0; yy < rows; yy++) {
                for (var xx:int = 0; xx < cols; xx++) {
                    if (board[yy][xx] == -1) continue;
                    var cnt:int = 0;
                    for (var dy:int = -1; dy <= 1; dy++)
                        for (var dx:int = -1; dx <= 1; dx++) {
                            var nx:int = xx + dx, ny:int = yy + dy;
                            if (nx >= 0 && nx < cols && ny >= 0 && ny < rows && board[ny][nx] == -1)
                                cnt++;
                        }
                    board[yy][xx] = cnt;
                }
            }
        }

        private function reveal(x:int, y:int):void {
            if (x < 0 || x >= cols || y < 0 || y >= rows) return;
            if (revealed[y][x] || flagged[y][x]) return;
            revealed[y][x] = true;
            if (board[y][x] == -1) { gameOver = true; return; }
            if (board[y][x] == 0) {
                for (var dy:int = -1; dy <= 1; dy++)
                    for (var dx:int = -1; dx <= 1; dx++)
                        reveal(x + dx, y + dy);
            }
            checkWin();
        }

        private function checkWin():void {
            var cnt:int = 0;
            for (var y:int = 0; y < rows; y++)
                for (var x:int = 0; x < cols; x++)
                    if (revealed[y][x]) cnt++;
            if (cnt == rows * cols - mineCount) won = true;
        }

        override public function onKeyDown(e:KeyboardEvent):void {
            if (e.keyCode == Keyboard.R) reset();
        }

        private function getCell(e:MouseEvent):Object {
            var mx:Number = e.stageX - offsetX;
            var my:Number = e.stageY - offsetY;
            var x:int = int(mx / cell);
            var y:int = int(my / cell);
            if (x < 0 || x >= cols || y < 0 || y >= rows) return null;
            return { x: x, y: y };
        }

        override public function onClick(e:MouseEvent):void {
            if (gameOver || won) return;
            var c:Object = getCell(e);
            if (c == null) return;
            if (firstClick) { placeMines(c.x, c.y); firstClick = false; }
            if (!flagged[c.y][c.x]) reveal(c.x, c.y);
        }

        override public function onRightClick(e:MouseEvent):void {
            if (gameOver || won) return;
            var c:Object = getCell(e);
            if (c == null) return;
            if (!revealed[c.y][c.x]) flagged[c.y][c.x] = !flagged[c.y][c.x];
        }

        override public function draw():void {
            beginTextFrame();
            clear();

            for (var y:int = 0; y < rows; y++) {
                for (var x:int = 0; x < cols; x++) {
                    var px:Number = offsetX + x * cell + 1;
                    var py:Number = offsetY + y * cell + 1;
                    var size:Number = cell - 2;

                    if (revealed[y][x]) {
                        if (board[y][x] == -1) {
                            // 炸弹
                            roundRect(px, py, size, size, 4, 0xc83232);
                            graphics.beginFill(0x000000);
                            graphics.drawCircle(px + size / 2, py + size / 2, 8);
                            graphics.endFill();
                        } else {
                            roundRect(px, py, size, size, 4, 0xbebec8);
                            var n:int = board[y][x];
                            if (n > 0) {
                                var colors:Array = [0, 0x0000ff, 0x009600, 0xff0000,
                                                    0x0000b4, 0x960000, 0x009696, 0x000000];
                                drawText(String(n), px + size / 2, py + size / 2,
                                         21, colors[n], "center");
                            }
                        }
                    } else {
                        var color:uint = flagged[y][x] ? 0xd2d23c : 0x505069;
                        roundRect(px, py, size, size, 4, color);
                        if (flagged[y][x]) {
                            // 真实三角旗
                            var px2:Number = px + size / 2;
                            var py2:Number = py + size / 2;
                            graphics.lineStyle(2, 0x000000);
                            graphics.moveTo(px2, py2 - 8);
                            graphics.lineTo(px2, py2 + 8);
                            graphics.lineStyle();
                            graphics.beginFill(0xff3232);
                            graphics.moveTo(px2, py2 - 8);
                            graphics.lineTo(px2 + 8, py2 - 3);
                            graphics.lineTo(px2, py2 + 2);
                            graphics.endFill();
                        }
                    }
                }
            }

            if (gameOver || won) {
                graphics.beginFill(0x000000, 0.6);
                graphics.drawRect(offsetX, offsetY, boardW, boardH);
                graphics.endFill();
                var msg:String = won ? "YOU WIN!" : "GAME OVER";
                var col:uint = won ? 0x00ff00 : 0xffffff;
                drawText(msg, offsetX + boardW / 2, offsetY + boardH / 2 - 20,
                         46, col, "center");
                drawText("Press R to restart", offsetX + boardW / 2, offsetY + boardH / 2 + 30,
                         26, 0xffffff, "center");
            }
        }
    }
}