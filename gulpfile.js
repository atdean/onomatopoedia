const gulp = require('gulp');
const sass = require('gulp-sass');
const del = require('del');

function styles(cb) {
    return gulp.src('public/scss/**/*.scss')
        .pipe(sass().on('error', sass.logError))
        .pipe(gulp.dest('public/css/'));
}

function clean(cb) {
    return del([
        'public/css/*.css'
    ]);
}

function watch(cb) {
    gulp.watch('public/scss/**/*.scss', gulp.series(clean, styles));
}

exports.build = gulp.series(clean, styles);
exports.clean = clean;
exports.watch = watch;