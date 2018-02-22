// This file's only purpose is to have convenient code completion in IDEs
// These variables and functions are exposed from go

/***
 * The current OS.
 * https://github.com/golang/go/blob/master/src/go/build/syslist.go
 * @type {string}
 */
var os ="";

/***
 * The current architecture
 * https://github.com/golang/go/blob/master/src/go/build/syslist.go
 * @type {string}
 */
var arch = "";

/***
 * The file containing the code, with full path. E.g. /tmp/foo/bar.php
 * @type {string}
 */
var codePath = "";

/***
 * Check if the given program is installed.
 * Supports Windows and any OS that knows the "which" command
 *
 * @param bin string
 * @return bool
 */
function isInstalled(bin){}

/***
 * A binary and some optional arguments. This will be executed on the host system, be careful!
 * Command output will be sent to the browser
 * @param bin string
 * @param args ...string
 * @return bool
 */
function exec(bin, ...args) {}

/***
 * A binary and some optional arguments. This will be executed on the host system, be careful!
 * Command output will not be sent to the browser
 * @param bin string
 * @param args ...string
 * @return bool
 */
function execQuiet(bin, ...args) {}

/***
 * Send message to browser, using stdout json field
 * @param message string
 * @return bool
 */
function sendStdOut(message){}

/***
 * Send message to browser, using stderr json field
 * @param message string
 * @return bool
 */
function sendStdErr(message){}

/***
 * Checks if a given image, e.g. php or php:7.2 is installed on the system
 * @param image string
 * @return bool
 */
function isDockerImageInstalled(image){}

/***
 * Pulls the given docker image
 * To pull from e.g. a local registry, see https://docs.docker.com/engine/reference/commandline/pull/#pull-from-a-different-registry
 * Will be executed immediatly with no browser output!
 * @param image
 * @return bool
 */
function pullDockerImage(image){}