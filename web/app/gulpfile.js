var gulp = require('gulp'),
    plugins = require('gulp-load-plugins')(),
    webpack = require('webpack'),
    webpackStream = require('webpack-stream');

gulp.task('lib', function () {
    gulp.src(['./node_modules/medium-editor/dist/css/medium-editor.css'])
        .pipe(plugins.rename('_medium-editor.scss'))
        .pipe(gulp.dest('./scss/plugins/'));
    gulp.src(['./node_modules/normalize-css/normalize.css'])
        .pipe(plugins.rename('_normalize.scss'))
        .pipe(gulp.dest('./scss/plugins/'));
});

gulp.task('sass', function () {
    gulp.src('./scss/default.scss')
        .pipe(plugins.sass({outputStyle: 'compressed'}).on('error', plugins.sass.logError))
        .pipe(plugins.rename('default.min.css'))
        .pipe(gulp.dest('./../static/css'));
});

gulp.task('webpack', function () {
    gulp.src('./**/*.js')
        .pipe(webpackStream({
            entry: './js/default.js',
            output: {
                path: __dirname + '/../static/js',
                filename: 'default.min.js'
            }
        }))
        .pipe(plugins.uglify())
        .pipe(gulp.dest('./../static/js'));
});

gulp.task('default', function () {
    gulp.watch('./scss/**/*.scss', ['sass']);
    gulp.watch('./js/**/*.js', ['webpack']);
});
