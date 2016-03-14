var gulp = require('gulp'),
    plugins = require('gulp-load-plugins')(),
    webpack = require('webpack'),
    webpackStream = require('webpack-stream');

gulp.task('lib', function () {
    gulp.src(['./node_modules/medium-editor/dist/css/medium-editor.css'])
        .pipe(plugins.rename('_medium-editor.scss'))
        .pipe(gulp.dest('./build/scss/plugins/'));
    gulp.src(['./node_modules/normalize-css/normalize.css'])
        .pipe(plugins.rename('_normalize.scss'))
        .pipe(gulp.dest('./build/scss/plugins/'));
});

gulp.task('sass', function () {
    gulp.src('./build/scss/default.scss')
        .pipe(plugins.sass({outputStyle: 'compressed'}).on('error', plugins.sass.logError))
        .pipe(plugins.rename('default.min.css'))
        .pipe(gulp.dest('./public/css'));
});

gulp.task('webpack', function () {
    gulp.src('./build/**/*.js')
        .pipe(webpackStream({
            entry: './build/js/default.js',
            output: {
                path: __dirname + '/public/js',
                filename: 'default.min.js'
            }
        }))
        .pipe(plugins.uglify())
        .pipe(gulp.dest('./public/js'));
});

gulp.task('default', function () {
    gulp.watch('./build/scss/**/*.scss', ['sass']);
    gulp.watch('./build/js/**/*.js', ['webpack']);
});
