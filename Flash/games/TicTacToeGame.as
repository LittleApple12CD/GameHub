package games {
    import BaseGame;
    import flash.events.KeyboardEvent;
    import flash.events.MouseEvent;
    import flash.ui.Keyboard;

    public class TicTacToeGame extends BaseGame {
        private var areaSize:Number = 500;
        private var offsetX:Number = 175, offsetY:Number = 65;
        private var cellSize:Number = 500 / 3;

        private var board:Array;
        private var currentPlayer:String = "X";
        private var winner:String = null;
        private var gameOver:Boolean;
        private var moveCount:int;

        public function TicTacToeGame() {}

        override public function reset():void {
            board = [["","",""],["","",""],["","",""]];
            currentPlayer = "X";
            winner = null;
            gameOver = false;
            moveCount = 0;
        }

        private function checkWinner():String {
            var b:Array = board;
            for (var i:int = 0; i < 3; i++) {
                if (b[i][0] != "" && b[i][0] == b[i][1] && b[i][1] == b[i][2]) return b[i][0];
                if (b[0][i] != "" && b[0][i] == b[1][i] && b[1][i] == b[2][i]) return b[0][i];
            }
            if (b[0][0] != "" && b[0][0] == b[1][1] && b[1][1] == b[2][2]) return b[0][0];
            if (b[0][2] != "" && b[0][2] == b[1][1] && b[1][1] == b[2][0]) return b[0][2];
            return null;
        }

        override public function onKeyDown(e:KeyboardEvent):void {
            if (e.keyCode == Keyboard.R) reset();
        }

        override public function onClick(e:MouseEvent):void {
            if (gameOver) return;
            var mx:Number = e.stageX - offsetX;
            var my:Number = e.stageY - offsetY;
            var x:int = int(mx / cellSize);
            var y:int = int(my / cellSize);
            if (x < 0 || x > 2 || y < 0 || y > 2) return;
            if (board[y][x] != "") return;

            board[y][x] = currentPlayer;
            moveCount++;
            winner = checkWinner();
            if (winner != null || moveCount == 9) gameOver = true;
            else currentPlayer = (currentPlayer == "X") ? "O" : "X";
        }

        override public function draw():void {
            beginTextFrame();
            clear();

            // 背景
            graphics.beginFill(0x1e1e2a);
            graphics.drawRect(offsetX, offsetY, areaSize, areaSize);
            graphics.endFill();

            // 格子 + 分隔线
            for (var r:int = 0; r < 3; r++) {
                for (var c:int = 0; c < 3; c++) {
                    roundRect(offsetX + c * cellSize + 2, offsetY + r * cellSize + 2,
                              cellSize - 4, cellSize - 4, 8, 0x323241);
                    var v:String = board[r][c];
                    if (v != "") {
                        drawText(v, offsetX + c * cellSize + cellSize / 2,
                                 offsetY + r * cellSize + cellSize / 2,
                                 92, v == "X" ? 0x3cc8ff : 0xffc83c, "center");
                    }
                }
            }

            // 状态文字
            var statusY:Number = offsetY + areaSize + 30;
            if (gameOver) {
                var msg:String = (winner != null) ? ("Player " + winner + " Wins!") : "Draw!";
                var col:uint = (winner != null) ? 0x00ff00 : 0xffff64;
                drawText(msg, 425, statusY, 46, col, "center");
                drawText("Press R to restart", 425, statusY + 40, 26, 0xc8c8c8, "center");
            } else {
                drawText("Player " + currentPlayer + "'s Turn", 425, statusY,
                         42, 0xffffff, "center");
            }
        }
    }
}