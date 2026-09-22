package {
    import flash.display.Sprite;
    import flash.events.Event;
    import flash.events.KeyboardEvent;
    import flash.events.MouseEvent;
    import flash.ui.Keyboard;
    import flash.utils.getTimer;
    import flash.text.TextField;
    import flash.text.TextFormat;
    import flash.text.TextFieldAutoSize;

    public class BaseGame extends Sprite {
        public var keys:Object = {};
        public var running:Boolean = false;
        public var onExitToMenu:Function = null;

        protected var _lastTime:Number = 0;
        protected var _bgColor:uint = 0x191923;

        // TextField 缓存池
        private var _textPool:Array = [];
        private var _textUsed:Array = [];

        public function BaseGame() {
            addEventListener(Event.ADDED_TO_STAGE, _onAdded);
        }

        private function _onAdded(e:Event):void {
            removeEventListener(Event.ADDED_TO_STAGE, _onAdded);
            stage.addEventListener(KeyboardEvent.KEY_DOWN, _onKeyDown);
            stage.addEventListener(KeyboardEvent.KEY_UP, _onKeyUp);
            addEventListener(MouseEvent.CLICK, _onClick);
            addEventListener(MouseEvent.RIGHT_CLICK, _onRightClick);
            stage.addEventListener(MouseEvent.RIGHT_CLICK, _suppress);
        }

        private function _suppress(e:MouseEvent):void { e.stopPropagation(); }

        private function _onKeyDown(e:KeyboardEvent):void {
            if (e.keyCode == Keyboard.ESCAPE) {
                if (onExitToMenu != null) onExitToMenu();
                return;
            }
            keys[String.fromCharCode(e.keyCode).toLowerCase()] = true;
            keys[e.keyCode] = true;
            onKeyDown(e);
        }
        private function _onKeyUp(e:KeyboardEvent):void {
            keys[String.fromCharCode(e.keyCode).toLowerCase()] = false;
            keys[e.keyCode] = false;
            onKeyUp(e);
        }
        private function _onClick(e:MouseEvent):void { onClick(e); }
        private function _onRightClick(e:MouseEvent):void { onRightClick(e); }

        public function reset():void {}
        public function onKeyDown(e:KeyboardEvent):void {}
        public function onKeyUp(e:KeyboardEvent):void {}
        public function onClick(e:MouseEvent):void {}
        public function onRightClick(e:MouseEvent):void {}

        public function start():void {
            reset();
            running = true;
            _lastTime = getTimer();
            addEventListener(Event.ENTER_FRAME, _loop);
        }

        private function _loop(e:Event):void {
            if (!running) return;
            var t:Number = getTimer();
            var dt:Number = t - _lastTime;
            _lastTime = t;
            update(dt);
            draw();
        }

        public function destroy():void {
            running = false;
            removeEventListener(Event.ENTER_FRAME, _loop);
            if (stage) {
                stage.removeEventListener(KeyboardEvent.KEY_DOWN, _onKeyDown);
                stage.removeEventListener(KeyboardEvent.KEY_UP, _onKeyUp);
                stage.removeEventListener(MouseEvent.RIGHT_CLICK, _suppress);
            }
            removeEventListener(MouseEvent.CLICK, _onClick);
            removeEventListener(MouseEvent.RIGHT_CLICK, _onRightClick);
            for (var i:int = 0; i < _textPool.length; i++) {
                if (_textPool[i] && _textPool[i].parent) {
                    _textPool[i].parent.removeChild(_textPool[i]);
                }
            }
            _textPool = [];
            _textUsed = [];
        }

        public function update(dt:Number):void {}
        public function draw():void {}

        // ============== 绘制工具 ==============

        public function clear(bg:uint = 0x191923):void {
            graphics.clear();
            graphics.beginFill(bg);
            graphics.drawRect(0, 0, 850, 700);
            graphics.endFill();
        }

        public function roundRect(x:Number, y:Number, w:Number, h:Number,
                                  r:Number, fill:uint, stroke:* = null):void {
            if (w < 2 * r) r = w / 2;
            if (h < 2 * r) r = h / 2;

            graphics.beginFill(fill);
            graphics.drawRoundRect(x, y, w, h, r * 2, r * 2);
            graphics.endFill();

            if (stroke != null) {
                graphics.lineStyle(2, uint(stroke));
                graphics.drawRoundRect(x, y, w, h, r * 2, r * 2);
                graphics.lineStyle();
            }
        }

        protected function beginTextFrame():void {
            for (var i:int = 0; i < _textPool.length; i++) {
                _textPool[i].visible = false;
            }
            _textUsed = [];
        }

        public function drawText(text:String, x:Number, y:Number,
                                 size:Number = 24, color:uint = 0xffffff,
                                 align:String = "left"):TextField {
            var tf:TextField;
            if (_textUsed.length < _textPool.length) {
                tf = _textPool[_textUsed.length];
            } else {
                tf = new TextField();
                tf.selectable = false;
                tf.mouseEnabled = false;
                tf.multiline = false;
                tf.wordWrap = false;
                addChild(tf);
                _textPool.push(tf);
            }
            _textUsed.push(tf);

            tf.visible = true;
            tf.text = text;
            tf.width = 850;
            tf.height = size + 12;

            var fmt:TextFormat = new TextFormat();
            fmt.font = "_sans";
            fmt.size = size;
            fmt.color = color;
            fmt.align = "left";
            tf.defaultTextFormat = fmt;
            tf.setTextFormat(fmt);

            tf.autoSize = TextFieldAutoSize.LEFT;
            var textW:Number = tf.textWidth;

            if (align == "center") {
                tf.x = x - textW / 2;
            } else if (align == "right") {
                tf.x = x - textW;
            } else {
                tf.x = x;
            }
            tf.y = y - size * 0.58;
            return tf;
        }
    }
}