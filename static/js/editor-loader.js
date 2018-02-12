var aceInit = false;

slideshow.on("showSlide", function () {
    if (aceInit) {
        return;
    }
    var elems = document.querySelectorAll(".editor");
    var i = 1;
    elems.forEach(function (elem) {
        var theme = "solarized_dark";
        var mode = "php";
        var filename;
        if (elem.dataset.filename) {
            filename = elem.dataset.filename;
        }
        if (!filename) {
            console.error("No filename defined! One editor div is missing the data-filename attribute");
            return;
        }
        var executor;
        if (elem.dataset.executor) {
            executor = elem.dataset.executor;
        }
        if (!executor) {
            console.error("No executor defined! One editor div is missing the data-executor attribute");
            return;
        }
        if (elem.dataset.theme) {
            theme = elem.dataset.theme;
        }
        if (elem.dataset.mode) {
            mode = elem.dataset.mode;
        }
        var editor = ace.edit(elem);
        editor.setTheme("ace/theme/" + theme);
        editor.session.setMode("ace/mode/" + mode);
        editor.on("focus", function () {
            slideshow.pause()
        });
        editor.on("blur", function () {
            slideshow.resume()
        });

        var execButton = document.createElement('button');
        var outputLog = document.createElement("div");
        outputLog.classList.add("outputlog");
        var outputpre = document.createElement("pre");
        outputLog.appendChild(outputpre);
        execButton.innerHTML = "Execute";
        execButton.onclick = function () {
            var last_index = 0;
            var postdata = "executor=" + executor + "&filename=" + encodeURIComponent(filename) + "&payload=" + encodeURIComponent(editor.getValue());
            var xhr = new XMLHttpRequest();
            xhr.open("POST", "/exec", true);
            xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
            xhr.onprogress = function () {
                var curr_index = xhr.responseText.length;
                if (last_index === curr_index) return;
                var s = xhr.responseText.substring(last_index, curr_index);
                last_index = curr_index;
                var resp = JSON.parse(s);
                outputpre.innerHTML += "<span>" + resp.stdout + "</span>";
                outputpre.scrollTop = outputpre.scrollHeight;
            };
            xhr.send(postdata);
        };
        elem.insertAdjacentElement("afterend", execButton);
        execButton.insertAdjacentElement("afterend", outputLog);

        i++;
    });
    aceInit = true;
});