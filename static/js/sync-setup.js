const modifierKeys = [
    225, //AltGr
    40, //DownArrow
    39, //RightArrow
    38, //UpArrow
    37, //LeftArrow
    36, //Pos1
    35, //End
    34, //PageDown
    33, //PageUp
    27, //Esc
    18, //Alt
    17, //Ctrl
    16, //Shift
];

function isModifierKey(keyCode) {
    for (let i = 0; i < modifierKeys.length; i++) {
        if (modifierKeys[i] === keyCode) {
            return true;
        }
    }
    return false;
}

const ws = new WebSocket('ws://' + location.host + '/editorsync');

let EditorSync = {
    sub: []
};

ws.onopen = function () {
    console.log('Connected to WebSocket')
};

ws.onerror = function (error) {
    console.log('WebSocket Error ' + error);
};

ws.onmessage = function (evt) {
    EditorSync.sub.forEach(function (func) {
        func(evt.data);
    });
};

function sendCursorPosition(pos, editorId) {
    let position = {
        editorId: editorId,
        type: 'click',
        row: pos.row,
        column: pos.column
    };
    ws.send(JSON.stringify(position));
}

function sendCmdKeyEvent(cmdContent, editorId) {
    let event = {
        cmdContent: cmdContent,
        editorId: editorId,
        type: 'cmdUpdate'
    };
    ws.send(JSON.stringify(event));
}

function sendKeyEvent(editorVal, evt, editorId) {
    let event = {
        editorContent: editorVal,
        editorId: editorId,
        altKey: evt.altKey,
        bubbles: evt.bubbles,
        cancelBubble: evt.cancelBubble,
        cancelable: evt.cancelable,
        charCode: evt.charCode,
        code: evt.code,
        composed: evt.composed,
        ctrlKey: evt.ctrlKey,
        defaultPrevented: evt.defaultPrevented,
        detail: evt.detail,
        eventPhase: evt.eventPhase,
        isComposing: evt.isComposing,
        isTrusted: evt.isTrusted,
        key: evt.key,
        keyCode: evt.keyCode,
        location: evt.location,
        metaKey: evt.metaKey,
        repeat: evt.repeat,
        returnValue: evt.returnValue,
        shiftKey: evt.shiftKey,
        type: evt.type,
        which: evt.which
    };
    ws.send(JSON.stringify(event));
}

function sendLogUpdate(editorId, clear, stdout, stderr){
    let event = {
        type: 'logupdate',
        editorId: editorId,
        stdout: stdout,
        stderr: stderr,
        clear: clear
    };
    ws.send(JSON.stringify(event));
}

function sendLogScrollUpdate(editorId, scrollTop){
    let event = {
        type: 'logscrollupdate',
        editorId: editorId,
        scrollTop: scrollTop,
    };
    ws.send(JSON.stringify(event));
}