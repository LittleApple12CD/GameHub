package games {
    import BaseGame;
    import flash.events.KeyboardEvent;
    import flash.ui.Keyboard;

    public class TetrisGame extends BaseGame {
        private var SHAPES:Object = {
            I: [[0,0,0,0],[1,1,1,1],[0,0,0,0],[0,0,0,0]],
            O: [[1,1],[1,1]],
            T: [[0,1,0],[1,1,1],[0,0,0]],
            S: [[0,1,1],[1,1,0],[0,0,0]],
            Z: [[1,1,0],[0,1,1],[0,0,0]],
            L: [[1,0,0],[1,1,1],[0,0,0]],
            J: [[0,0,1],[1,1,1],[0,0,0]]
        };
        private var COLORS:Object = {
            I: 0x00f0f0, O: 0xf0f000, T: 0xb400f0,
            S: 0x00f000, Z: 0xf00000, L: 0xf0a000, J: 0x0000f0
        };
        private var shapeKeys:Array = ["I","O","T","S","Z","L","J"];

        private var cols:int = 10;
        private var rows:int = 20;
        private var cell:Number = 30;
        private var boardW:Number = 300;
        private var boardH:Number = 600;
        private var offsetX:Number = 275;
        private var offsetY:Number = 48;

        private var board:Array;
        private var score:int;
        private var gameOver:Boolean;
        private var fallTimer:Number;
        private var fallDelay:Number = 500;

        private var pieceKey:String;
        private var piece:Array;
        private var px:int, py:int;

        public function TetrisGame() {}

        override public function reset():void {
            board = [];
            for (var r:int = 0; r < rows; r++) {
                var row:Array = [];
                for (var c:int = 0; c < cols; c++) row.push(0);
                board.push(row);
            }
            score = 0;
            gameOver = false;
            fallTimer = 0;
            spawn();
        }

        private function spawn():void {
            pieceKey = shapeKeys[int(Math.random() * 7)];
            var src:Array = SHAPES[pieceKey];
            piece = [];
            for (var r:int = 0; r < src.length; r++) {
                var row:Array = [];
                for (var c:int = 0; c < src[r].length; c++) row.push(src[r][c]);
                piece.push(row);
            }
            px = int(cols / 2 - piece[0].length / 2);
            py = 0;
            if (collide(piece, px, py)) gameOver = true;
        }

        private function collide(shape:Array, x:int, y:int):Boolean {
            for (var r:int = 0; r < shape.length; r++) {
                for (var c:int = 0; c < shape[r].length; c++) {
                    if (!shape[r][c]) continue;
                    var bx:int = x + c, by:int = y + r;
                    if (bx < 0 || bx >= cols || by >= rows) return true;
                    if (by >= 0 && board[by][bx]) return true;
                }
            }
            return false;
        }

        private function rotatePiece():void {
            var newPiece:Array = [];
            var rowsN:int = piece[0].length;
            for (var i:int = 0; i < rowsN; i++) {
                var newRow:Array = [];
                for (var j:int = piece.length - 1; j >= 0; j--) {
                    newRow.push(piece[j][i]);
                }
                newPiece.push(newRow);
            }
            if (!collide(newPiece, px, py)) piece = newPiece;
        }

        private function lock():void {
            for (var r:int = 0; r < piece.length; r++) {
                for (var c:int = 0; c < piece[r].length; c++) {
                    if (piece[r][c]) {
                        var by:int = py + r, bx:int = px + c;
                        if (by >= 0) board[by][bx] = pieceKey;
                    }
                }
            }
            clearLines();
            spawn();
        }

        private function clearLines():void {
            var cleared:int = 0;
            for (var r:int = rows - 1; r >= 0; r--) {
                var full:Boolean = true;
                for (var c:int = 0; c < cols; c++) if (!board[r][c]) { full = false; break; }
                if (full) {
                    board.splice(r, 1);
                    var empty:Array = [];
                    for (var k:int = 0; k < cols; k++) empty.push(0);
                    board.unshift(empty);
                    cleared++;
                    r++;
                }
            }
            var pts:Array = [0, 100, 300, 500, 800];
            score += pts[cleared];
        }

        override public function onKeyDown(e:KeyboardEvent):void {
            if (e.keyCode == Keyboard.R) { reset(); return; }
            if (gameOver) return;
            switch (e.keyCode) {
                case Keyboard.LEFT:
                    if (!collide(piece, px - 1, py)) px--;
                    break;
                case Keyboard.RIGHT:
                    if (!collide(piece, px + 1, py)) px++;
                    break;
                case Keyboard.DOWN:
                    if (!collide(piece, px, py + 1)) py++;
                    break;
                case Keyboard.UP:
                case Keyboard.SPACE:
                    rotatePiece();
                    break;
            }
        }

        override public function update(dt:Number):void {
            if (gameOver) return;
            fallTimer += dt;
            if (fallTimer >= fallDelay) {
                fallTimer = 0;
                if (!collide(piece, px, py + 1)) py++;
                else lock();
            }
        }

        override public function draw():void {
            beginTextFrame();
            clear();

            // 棋盘背景
            graphics.beginFill(0x14141e);
            graphics.drawRect(offsetX, offsetY, boardW, boardH);
            graphics.endFill();

            // 已落方块
            for (var r:int = 0; r < rows; r++) {
                for (var c:int = 0; c < cols; c++) {
                    var v:String = board[r][c];
                    if (v) {
                        roundRect(offsetX + c * cell + 1, offsetY + r * cell + 1,
                                  cell - 2, cell - 2, 4, COLORS[v]);
                    }
                }
            }

            // 当前方块
            if (!gameOver) {
                for (var pr:int = 0; pr < piece.length; pr++) {
                    for (var pc:int = 0; pc < piece[pr].length; pc++) {
                        if (piece[pr][pc]) {
                            roundRect(offsetX + (px + pc) * cell + 1,
                                      offsetY + (py + pr) * cell + 1,
                                      cell - 2, cell - 2, 4, COLORS[pieceKey]);
                        }
                    }
                }
            }

            drawText("Score: " + score, offsetX + boardW + 20, offsetY + 20, 32);

            if (gameOver) {
                graphics.beginFill(0x000000, 0.7);
                graphics.drawRect(offsetX, offsetY, boardW, boardH);
                graphics.endFill();
                drawText("GAME OVER",
                         offsetX + boardW / 2, offsetY + boardH / 2 - 20,
                         46, 0xffffff, "center");
                drawText("Press R to restart",
                         offsetX + boardW / 2, offsetY + boardH / 2 + 30,
                         26, 0xffffff, "center");
            }
        }
    }
}