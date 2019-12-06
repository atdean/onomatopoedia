const { src, dest, series } = require('gulp');
const sass = require('gulp-sass');
const del = require('del');

function styles(cb) {
    return src('public/scss/**/*.scss')
        .pipe(sass().on('error', sass.logError))
        .pipe(dest('public/css/'));
}

function clean(cb) {
    return del([
        'public/css/*.css'
    ]);
}

exports.build = series(clean, styles);
exports.clean = clean;