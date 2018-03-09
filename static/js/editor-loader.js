let aceInit = false;

slideshow.on('showSlide', function () {
    if (aceInit) {
        return;
    }
    const elems = document.querySelectorAll('.editor');
    let i = 1;
    elems.forEach(function (elem) {
        let theme = 'solarized_dark',
            mode = 'php',
            editorId = i,
            executor,
            executors,
            filename,
            execButton,
            clearButton,
            outputLog,
            outputPre,
            editor,
            cmdContainer,
            cmdLine;
        if (elem.dataset.filename) {
            filename = elem.dataset.filename;
        }
        if (!filename) {
            console.error('No filename defined! One editor div is missing the data-filename attribute');
            return;
        }
        if (elem.dataset.executor) {
            executor = elem.dataset.executor;
        } else if (elem.dataset.executors) {
            let execs = elem.dataset.executors.split(';');
            executor = execs[0];
            if (execs.length > 1) {
                executors = execs;
            }
        }
        if (!executor) {
            console.error('No executor defined! One editor div is missing the data-executor and data-executors attribute');
            return;
        }
        if (elem.dataset.theme) {
            theme = elem.dataset.theme;
        }
        if (elem.dataset.mode) {
            mode = elem.dataset.mode;
        }
        editor = ace.edit(elem);
        editor.setTheme('ace/theme/' + theme);
        editor.session.setMode('ace/mode/' + mode);
        editor.on('focus', function () {
            slideshow.pause()
        });
        editor.on('blur', function () {
            slideshow.resume()
        });

        if (elem.dataset.showcmd === "true") {
            cmdContainer = document.createElement('div');
            cmdContainer.setAttribute('class', 'cmdLineContainer');
            cmdLine = document.createElement('input');
            cmdLine.setAttribute('type', 'text');
            cmdLine.setAttribute('class', 'cmdLine');
            cmdLine.addEventListener('focus', function () {
                slideshow.pause();
            });
            cmdLine.addEventListener('blur', function () {
                slideshow.resume();
            });
            cmdContainer.appendChild(cmdLine);
        }

        execButton = document.createElement('button');
        clearButton = document.createElement('button');
        outputLog = document.createElement('div');
        outputLog.classList.add('outputlog');
        outputPre = document.createElement('pre');
        outputLog.appendChild(outputPre);
        execButton.innerHTML = 'Execute code';
        execButton.classList.add('editorbutton');
        execButton.setAttribute('accesskey', 'x');
        clearButton.innerHTML = 'Clear log';
        clearButton.classList.add('editorbutton');
        clearButton.setAttribute('accesskey', 'l');

        let textarea = editor.textInput.getElement();

        /******** Dynamic Editor/Log view width&height ********/

        if (elem.dataset.editorheight) {
            elem.style.height = elem.dataset.editorheight;
        }

        if (elem.dataset.editorwidth) {
            elem.style.width = elem.dataset.editorwidth;
            cmdContainer.width = elem.dataset.editorwidth;
            cmdLine.width = elem.dataset.editorwidth;
        }

        if (elem.dataset.logheight) {
            outputLog.style.height = elem.dataset.logheight;
            outputPre.style.height = elem.dataset.logheight;
        }

        if (elem.dataset.logwidth) {
            outputLog.style.width = elem.dataset.logwidth;
            outputPre.style.width = elem.dataset.logwidth;
        }

        /******** Event Setup ************/



        editor.on('click', function (evt) {
            sendCursorPosition(editor.getCursorPosition(), editorId);
        });

        textarea.addEventListener('keydown', function (evt) {
            if (document.activeElement !== textarea) { // prevent event listener loop
                return
            }
            sendKeyEvent(editor.getValue(), evt, editorId);
            if (!isModifierKey(evt.keyCode)) { // synchronize cursor position after content update
                sendCursorPosition(editor.getCursorPosition(), editorId);
            }
        });

        textarea.addEventListener('keyup', function (evt) {
            if (document.activeElement !== textarea) { // prevent event listener loop
                return
            }
            sendKeyEvent(editor.getValue(), evt, editorId);
            if (!isModifierKey(evt.keyCode)) { // synchronize cursor position after content update
                sendCursorPosition(editor.getCursorPosition(), editorId);
            }
        });

        if (cmdLine) {
            cmdLine.addEventListener('keyup', function (evt) {
                if (document.activeElement !== cmdLine) { // prevent event listener loop
                    return
                }
                sendCmdKeyEvent(cmdLine.value, editorId);
            });
        }

        EditorSync.sub.push(function (evt) {
            let event = JSON.parse(evt);
            if (editorId === event.editorId && event.type === 'cmdUpdate' && cmdLine) {
                cmdLine.value = event.cmdContent;
            }
            else if (editorId === event.editorId && event.type === 'logupdate') {
                if (event.clear === true) {
                    outputPre.innerText = "";
                    return;
                }
                if (event.stdout !== undefined && event.stdout !== "") {
                    outputPre.innerHTML += "<span>" + event.stdout + "</span>";
                }
                if (event.stderr !== undefined && event.stderr !== "") {
                    outputPre.innerHTML += "<span style='color: red;'>" + event.stderr + "</span>";
                }
                outputPre.scrollTop = outputPre.scrollHeight;
            }
            else if (editorId === event.editorId && document.activeElement !== textarea) {
                if (event.type === 'click') {
                    editor.selection.moveTo(event.row, event.column);
                    return;
                }

                // if there is a content update, replace editor content and return.
                if (!isModifierKey(event.keyCode)) {
                    editor.setValue(event.editorContent, 1);
                    return;
                }

                // position/modifier keys can be sent directly to the textarea. "content" keys don't work.
                let e = document.createEvent("Event");
                e.initEvent(event.type, true, true);
                e.altKey = event.altKey;
                e.charCode = event.charCode;
                e.code = event.code;
                e.ctrlKey = event.ctrlKey;
                e.key = event.key;
                e.metaKey = event.metaKey;
                e.repeat = event.repeat;
                e.shiftKey = event.shiftKey;
                e.keyCodeVal = event.keyCode;
                e.whichVal = event.which;
                Object.defineProperty(e, 'keyCode', {
                    get: function () {
                        return this.keyCodeVal;
                    }
                });
                Object.defineProperty(e, 'which', {
                    get: function () {
                        return this.whichVal;
                    }
                });
                textarea.dispatchEvent(e);
            }
        });

        clearButton.onclick = function () {
            sendLogUpdate(editorId, true)
        };

        execButton.onclick = function () {
            let cmdArgs = "";
            if (cmdLine) {
                cmdArgs = cmdLine.value;
            }
            let postData = 'editorId=' + editorId + '&executor=' + executor + '&filename=' + encodeURIComponent(filename) + '&payload=' + encodeURIComponent(editor.getValue()) + '&cmdargs=' + encodeURIComponent(cmdArgs),
                xhr = new XMLHttpRequest();
            xhr.open('POST', '/exec', true);
            xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
            xhr.send(postData);
        };

        /******** /Event Setup ************/
        elem.insertAdjacentElement('afterend', clearButton);
        elem.insertAdjacentElement('afterend', execButton);
        if (executors) {
            let select = document.createElement('select');
            executors.forEach(function (exec) {
                select.options.add(new Option(exec, exec));
            });
            select.onchange = function () {
                executor = select.value;
            };
            execButton.insertAdjacentElement('afterend', select);
        }
        clearButton.insertAdjacentElement('afterend', outputLog);
        if (cmdContainer) {
            elem.insertAdjacentElement('afterend', cmdContainer)
        }
        i++;
    });
    aceInit = true;
});