let aceInit = false;

slideshow.on("showSlide", function () {
    if (aceInit) {
        return;
    }
    const elems = document.querySelectorAll(".editor");
    let i = 1;
    elems.forEach(function (elem) {
        let theme = "solarized_dark",
            mode = "php",
            executor,
            executors,
            filename,
            execButton,
            clearButton,
            outputLog,
            outputPre,
            editor;
        if (elem.dataset.filename) {
            filename = elem.dataset.filename;
        }
        if (!filename) {
            console.error("No filename defined! One editor div is missing the data-filename attribute");
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
            console.error("No executor defined! One editor div is missing the data-executor and data-executors attribute");
            return;
        }
        if (elem.dataset.theme) {
            theme = elem.dataset.theme;
        }
        if (elem.dataset.mode) {
            mode = elem.dataset.mode;
        }
        editor = ace.edit(elem);
        editor.setTheme("ace/theme/" + theme);
        editor.session.setMode("ace/mode/" + mode);
        editor.on("focus", function () {
            slideshow.pause()
        });
        editor.on("blur", function () {
            slideshow.resume()
        });

        execButton = document.createElement('button');
        clearButton = document.createElement('button');
        outputLog = document.createElement("div");
        outputLog.classList.add("outputlog");
        outputPre = document.createElement("pre");
        outputLog.appendChild(outputPre);
        execButton.innerHTML = "Execute code";
        execButton.classList.add("editorbutton");
        execButton.setAttribute("accesskey", "x");
        clearButton.innerHTML = "Clear log";
        clearButton.classList.add("editorbutton");
        clearButton.setAttribute("accesskey", "l");
        clearButton.onclick = function () {
            outputPre.innerText = "";
        };
        execButton.onclick = function () {
            let last_index = 0,
                postData = "executor=" + executor + "&filename=" + encodeURIComponent(filename) + "&payload=" + encodeURIComponent(editor.getValue()),
                xhr = new XMLHttpRequest();

            xhr.open("POST", "/exec", true);
            xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
            xhr.onprogress = function () {
                let curr_index = xhr.responseText.length,
                    s,
                    resp,
                    stdout,
                    stderr;
                if (last_index === curr_index) return;
                s = xhr.responseText.substring(last_index, curr_index);
                last_index = curr_index;
                resp = JSON.parse(s);
                stdout = resp.stdout;
                stderr = resp.stderr;
                if (stdout !== undefined && stdout !== "") {
                    outputPre.innerHTML += "<span>" + resp.stdout + "</span>";
                }
                if (stderr !== undefined && stderr !== "") {
                    outputPre.innerHTML += "<span style='color: red;'>" + resp.stderr + "</span>";
                }
                outputPre.scrollTop = outputPre.scrollHeight;
            };
            xhr.send(postData);
        };
        elem.insertAdjacentElement("afterend", clearButton);
        elem.insertAdjacentElement("afterend", execButton);
        if (executors) {
            let select = document.createElement("select");
            executors.forEach(function (exec) {
                select.options.add(new Option(exec, exec));
            });
            select.onchange = function () {
                executor = select.value;
            };
            execButton.insertAdjacentElement("afterend", select);
        }
        clearButton.insertAdjacentElement("afterend", outputLog);

        i++;
    });
    aceInit = true;
});