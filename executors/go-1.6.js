function getFileName(str) {
    return str.split('\\').pop().split('/').pop();
}

var filename = getFileName(codePath);

var base = codePath.substring(0, codePath.length - filename.length - 1);

exec("docker", "run", "--rm", "-v", base + ":/code", "-w", "/code", "golang:1.6", "go", "run", filename);