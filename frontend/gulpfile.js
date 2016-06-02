const gulp = require('gulp');
const rename = require('gulp-rename');
const psi = require('psi');
const imagemin = require('gulp-imagemin');
const imageResize = require('gulp-image-resize');
const autoprefixer = require('gulp-autoprefixer');
const sass = require('gulp-sass');
const cleanCSS = require('gulp-clean-css');


const config = {
    siteUrl: "http://c713296e.ngrok.io",
    imageSrc: 'src/images/*',
    imageDest: 'src/images',
    thumbnails: 'src/images/thumbnails/*',
    cssInput: 'src/css/*.scss',
    cssOutput: 'src'
};

const sassOptions = {
  errLogToConsole: true,
  outputStyle: 'expanded'
};


gulp.task('images', function() {
   gulp.src(config.imageSrc)

   .pipe(imagemin({
       progressive: true,
       optimizationLevel: 7
   }))

   .pipe(gulp.dest(config.imageDest));
});

gulp.task('thumbnail', function(){
    gulp.src(config.thumbnails)

    .pipe(imageResize({ width: 115}))

    .pipe(gulp.dest(config.viewImageDest));
})


gulp.task('mobile', function() {
    psi.output(config.siteUrl, {
        nokey: 'true',
        strategy: 'mobile',
    }).then(function (data) {
        console.log('Done');
    });
});

gulp.task('desktop', function() {
    psi.output(config.siteUrl, {
        nokey: 'true',
        strategy: 'desktop',
    }).then(function () {
        console.log('Done');
    });
});


gulp.task('build-css', function() {
  gulp.src(config.cssInput)
  .pipe(sass(sassOptions).on('error', sass.logError))
  .pipe(autoprefixer())
  .pipe(cleanCSS())
  .pipe(rename({suffix: '.min'}))
  .pipe(gulp.dest(config.cssOutput));
})

gulp.task('watch', function() {
  gulp.watch(config.cssInput, ['build-css']);
})



gulp.task('imageOpt', ['images', 'viewImages', 'thumbnail'])
gulp.task('build', ['build-css'])
gulp.task('default', ['build-css'])
