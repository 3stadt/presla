remark.macros.scaleImg = function (scaling) {
    let url = this;
    return '<img src="' + url + '" style="max-width: ' + scaling + '" />';
};

remark.macros.styleImg = function (sty) {
    let url = this;
    return '<img src="' + url + '" style="' + sty + '" />';
};

remark.macros.mp4 = function (width, height) {
    let video = this;
    return '<video class="vid" width="' + width + '" height="' + height + '" controls>' +
        '<source src="' + video + '" type="video/mp4">' +
        'Your browser does not support the video tag.' +
        '</video>'
};

let slideshow = remark.create({
    ratio: '16:9',
    highlightStyle: 'rainbow',
    sourceUrl: 'md/' + pres + '.md',
    countIncrementalSlides: false
});
