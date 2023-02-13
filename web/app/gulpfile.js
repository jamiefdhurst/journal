const gulp = require('gulp');
const plugins = require('gulp-load-plugins')();
const sass = require('gulp-sass')(require('node-sass'));
const webpackStream = require('webpack-stream');

gulp.task('sass', function () {
    return gulp.src('./scss/default.scss')
        .pipe(sass({outputStyle: 'compressed'}).on('error', sass.logError))
        .pipe(plugins.rename('default.min.css'))
        .pipe(gulp.dest('./../static/css'));
});

gulp.task('webpack', function () {
    return gulp.src('./**/*.js')
        .pipe(webpackStream({
            entry: './js/default.js',
            mode: 'production',
            output: {
                path: __dirname + '/../static/js',
                filename: 'default.min.js'
            }
        }))
        .pipe(plugins.uglify())
        .pipe(gulp.dest('./../static/js'));
});

gulp.task('default', function () {
    return gulp.watch(
        ['./scss/**/*.scss', './js/**/*.js'],
        gulp.parallel('sass', 'webpack')
    );
});
