package {
    import flash.display.Sprite;
    import flash.display.StageAlign;
    import flash.display.StageScaleMode;
    import flash.events.Event;
    import flash.events.KeyboardEvent;
    import flash.events.MouseEvent;
    import flash.ui.Keyboard;
    import flash.text.TextField;
    import flash.text.TextFormat;
    import flash.system.fscommand;
    import flash.utils.getDefinitionByName;

    import games.SnakeGame;
    import games.TetrisGame;
    import games.MinesweeperGame;
    import games.FlappyGame;
    import games.TankGame;
    import games.BreakoutGame;
    import games.TicTacToeGame;

    public class Main extends Sprite {
        private var gameList:Array = [
            { name: "Snake",       cls: SnakeGame },
            { name: "Tetris",      cls: TetrisGame },
            { name: "Minesweeper", cls: MinesweeperGame },
            { name: "Flappy Bird", cls: FlappyGame },
            { name: "TankBattle",  cls: TankGame },
            { name: "Breakout",    cls: BreakoutGame },
            { name: "Tic Tac Toe", cls: TicTacToeGame }
        ];

        private var menuLayer:Sprite;
        private var gameLayer:Sprite;
        private var headerLayer:Sprite;
        private var currentGame:BaseGame = null;
        private var selectedIndex:int = 0;
        private var itemBoxes:Array = [];
        private var itemTexts:Array = [];
        private var titleField:TextField;
        private var headerTitle:TextField;

        public function Main() {
            if (stage) init();
            else addEventListener(Event.ADDED_TO_STAGE, onAdded);
        }

        private function onAdded(e:Event):void {
            removeEventListener(Event.ADDED_TO_STAGE, onAdded);
            init();
        }

        private function init():void {
            stage.scaleMode = StageScaleMode.NO_SCALE;
            stage.align = StageAlign.TOP_LEFT;
            buildMenu();
            stage.addEventListener(KeyboardEvent.KEY_DOWN, onMenuKey);
        }

        // ==================== 菜单 ====================

        private function buildMenu():void {
            menuLayer = new Sprite();
            addChild(menuLayer);

            // 背景
            var bg:Sprite = new Sprite();
            bg.graphics.beginFill(0x191923);
            bg.graphics.drawRect(0, 0, 850, 700);
            bg.graphics.endFill();
            menuLayer.addChild(bg);

            // 标题
            titleField = new TextField();
            titleField.text = "GAME HUB";
            titleField.selectable = false;
            var tfmt:TextFormat = new TextFormat("_sans", 52, 0xffffff, true);
            tfmt.letterSpacing = 6;
            titleField.defaultTextFormat = tfmt;
            titleField.setTextFormat(tfmt);
            titleField.width = 850;
            titleField.height = 70;
            titleField.y = 40;
            tfmt.align = "center";
            titleField.setTextFormat(tfmt);
            menuLayer.addChild(titleField);

            // 选项
            itemBoxes = [];
            itemTexts = [];
            var itemW:Number = 280;
            var itemH:Number = 52;
            var gap:Number = 10;
            var startY:Number = 145;
            var startX:Number = (850 - itemW) / 2;

            for (var i:int = 0; i < gameList.length; i++) {
                var box:Sprite = new Sprite();
                box.name = "box" + i;
                box.addEventListener(MouseEvent.CLICK, makeLauncher(i));
                box.addEventListener(MouseEvent.MOUSE_OVER, makeHover(i));
                menuLayer.addChild(box);
                box.x = startX;
                box.y = startY + i * (itemH + gap);
                drawMenuBox(box, i == 0);
                itemBoxes.push(box);

                var tf:TextField = new TextField();
                tf.text = gameList[i].name;
                tf.selectable = false;
                tf.mouseEnabled = false;
                tf.width = itemW;
                tf.height = itemH;
                var fmt:TextFormat = new TextFormat("_sans", 26,
                    i == 0 ? 0x0f0f17 : 0xffffff);
                fmt.align = "center";
                tf.defaultTextFormat = fmt;
                tf.setTextFormat(fmt);
                tf.x = startX;
                tf.y = startY + i * (itemH + gap) + 11;
                menuLayer.addChild(tf);
                itemTexts.push(tf);
            }

            // EXIT 按钮
            var exitBtn:Sprite = new Sprite();
            exitBtn.name = "exitBtn";
            var exW:Number = 200;
            var exH:Number = 48;
            exitBtn.graphics.beginFill(0xc83232);
            exitBtn.graphics.drawRoundRect(0, 0, exW, exH, 20, 20);
            exitBtn.graphics.endFill();
            exitBtn.graphics.lineStyle(2, 0x7878a0);
            exitBtn.graphics.drawRoundRect(0, 0, exW, exH, 20, 20);
            exitBtn.graphics.lineStyle();
            exitBtn.x = (850 - exW) / 2;
            exitBtn.y = startY + gameList.length * (itemH + gap) + 20;
            exitBtn.buttonMode = true;
            exitBtn.addEventListener(MouseEvent.CLICK, function(e:MouseEvent):void {
                onExitRequest();
            });
            menuLayer.addChild(exitBtn);

            var exitTf:TextField = new TextField();
            exitTf.text = "EXIT";
            exitTf.selectable = false;
            exitTf.mouseEnabled = false;
            exitTf.width = exW;
            exitTf.height = exH;
            var efmt:TextFormat = new TextFormat("_sans", 24, 0xffffff);
            efmt.align = "center";
            exitTf.defaultTextFormat = efmt;
            exitTf.setTextFormat(efmt);
            exitTf.x = exitBtn.x;
            exitTf.y = exitBtn.y + 10;
            menuLayer.addChild(exitTf);
        }

        private function drawMenuBox(box:Sprite, selected:Boolean):void {
            box.graphics.clear();
            box.graphics.beginFill(selected ? 0x50c8ff : 0x37374b);
            box.graphics.drawRoundRect(0, 0, 280, 52, 24, 24);
            box.graphics.endFill();
            box.graphics.lineStyle(2, selected ? 0x78d8ff : 0x7878a0);
            box.graphics.drawRoundRect(0, 0, 280, 52, 24, 24);
            box.graphics.lineStyle();
        }

        private function makeLauncher(idx:int):Function {
            return function(e:MouseEvent):void { launchGame(idx); };
        }

        private function makeHover(idx:int):Function {
            return function(e:MouseEvent):void {
                selectedIndex = idx;
                updateSelection();
            };
        }

        private function updateSelection():void {
            for (var i:int = 0; i < itemBoxes.length; i++) {
                var sel:Boolean = (i == selectedIndex);
                drawMenuBox(itemBoxes[i], sel);
                var fmt:TextFormat = new TextFormat("_sans", 26,
                    sel ? 0x0f0f17 : 0xffffff);
                fmt.align = "center";
                itemTexts[i].setTextFormat(fmt);
            }
        }

        private function onMenuKey(e:KeyboardEvent):void {
            if (menuLayer == null || !menuLayer.visible) return;
            if (e.keyCode == Keyboard.UP) {
                selectedIndex = (selectedIndex - 1 + gameList.length) % gameList.length;
                updateSelection();
            } else if (e.keyCode == Keyboard.DOWN) {
                selectedIndex = (selectedIndex + 1) % gameList.length;
                updateSelection();
            } else if (e.keyCode == Keyboard.ENTER) {
                launchGame(selectedIndex);
            }
        }

        private function onExitRequest():void {
            try {
                var NativeAppClass:Class = getDefinitionByName("flash.desktop.NativeApplication") as Class;
                var nativeApp:Object = NativeAppClass["nativeApplication"];
                nativeApp["exit"]();
                return;
            } catch (e:Error) {
            
            }

            try {
                fscommand("quit", "");
                return;
            } catch (e:Error) {

            }

            showExitScreen();
        }

        private function showExitScreen():void {
            if (menuLayer) menuLayer.visible = false;
            if (gameLayer) gameLayer.visible = false;
            if (headerLayer) headerLayer.visible = false;

            var bye:Sprite = new Sprite();
            bye.name = "bye";
            bye.graphics.beginFill(0x0f0f17);
            bye.graphics.drawRect(0, 0, 850, 700);
            bye.graphics.endFill();
            addChild(bye);

            var t1:TextField = new TextField();
            t1.text = "GAME HUB";
            t1.selectable = false;
            t1.width = 850; t1.height = 60;
            var f1:TextFormat = new TextFormat("_sans", 52, 0x50c8ff, true);
            f1.align = "center";
            t1.defaultTextFormat = f1; t1.setTextFormat(f1);
            t1.y = 220;
            addChild(t1);

            var t2:TextField = new TextField();
            t2.text = "Game exited. Thanks for playing!";
            t2.selectable = false;
            t2.width = 850; t2.height = 40;
            var f2:TextFormat = new TextFormat("_sans", 22, 0xaaaaaa);
            f2.align = "center";
            t2.defaultTextFormat = f2; t2.setTextFormat(f2);
            t2.y = 300;
            addChild(t2);

            var t3:TextField = new TextField();
            t3.text = "You can close the window now.";
            t3.selectable = false;
            t3.width = 850; t3.height = 40;
            var f3:TextFormat = new TextFormat("_sans", 22, 0xaaaaaa);
            f3.align = "center";
            t3.defaultTextFormat = f3; t3.setTextFormat(f3);
            t3.y = 340;
            addChild(t3);
        }

        // ==================== 游戏切换 ====================

        private function launchGame(idx:int):void {
            menuLayer.visible = false;
            if (currentGame) {
                currentGame.destroy();
                if (currentGame.parent) currentGame.parent.removeChild(currentGame);
                currentGame = null;
            }
            if (gameLayer) { removeChild(gameLayer); gameLayer = null; }
            if (headerLayer) { removeChild(headerLayer); headerLayer = null; }

            gameLayer = new Sprite();
            addChild(gameLayer);

            var C:Class = gameList[idx].cls as Class;
            currentGame = new C() as BaseGame;
            currentGame.onExitToMenu = backToMenu;
            gameLayer.addChild(currentGame);
            currentGame.start();

            buildHeader(gameList[idx].name);
        }

        private function buildHeader(name:String):void {
            headerLayer = new Sprite();
            headerLayer.y = 0;
            addChild(headerLayer);

            // 顶部 40px 条
            headerLayer.graphics.beginFill(0x23232f);
            headerLayer.graphics.drawRect(0, 0, 850, 40);
            headerLayer.graphics.endFill();
            headerLayer.graphics.lineStyle(1, 0x37374b);
            headerLayer.graphics.moveTo(0, 40);
            headerLayer.graphics.lineTo(850, 40);
            headerLayer.graphics.lineStyle();

            // 返回按钮
            var back:Sprite = new Sprite();
            back.graphics.beginFill(0x23232f);
            back.graphics.drawRoundRect(0, 0, 70, 28, 8, 8);
            back.graphics.endFill();
            back.graphics.lineStyle(1, 0x50c8ff);
            back.graphics.drawRoundRect(0, 0, 70, 28, 8, 8);
            back.graphics.lineStyle();
            back.x = 14; back.y = 6;
            back.buttonMode = true;
            back.addEventListener(MouseEvent.CLICK, function(e:MouseEvent):void { backToMenu(); });
            headerLayer.addChild(back);

            var btf:TextField = new TextField();
            btf.text = "← Back";
            btf.selectable = false;
            btf.mouseEnabled = false;
            btf.width = 70; btf.height = 28;
            var bfmt:TextFormat = new TextFormat("_sans", 13, 0x50c8ff);
            bfmt.align = "center";
            btf.defaultTextFormat = bfmt; btf.setTextFormat(bfmt);
            btf.x = 14; btf.y = 13;
            headerLayer.addChild(btf);

            // 游戏名
            headerTitle = new TextField();
            headerTitle.text = name;
            headerTitle.selectable = false;
            headerTitle.mouseEnabled = false;
            headerTitle.width = 300; headerTitle.height = 40;
            var hfmt:TextFormat = new TextFormat("_sans", 14, 0xaaaaaa);
            hfmt.align = "right";
            headerTitle.defaultTextFormat = hfmt;
            headerTitle.setTextFormat(hfmt);
            headerTitle.x = 850 - 320;
            headerTitle.y = 12;
            headerLayer.addChild(headerTitle);
        }

        private function backToMenu():void {
            if (currentGame) {
                currentGame.destroy();
                if (currentGame.parent) currentGame.parent.removeChild(currentGame);
                currentGame = null;
            }
            if (gameLayer) { removeChild(gameLayer); gameLayer = null; }
            if (headerLayer) { removeChild(headerLayer); headerLayer = null; }
            menuLayer.visible = true;
            selectedIndex = 0;
            updateSelection();
        }
    }
}