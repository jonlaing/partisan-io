var gulp = require('gulp');
var config = require('../config').skeleton;

gulp.task('skeleton', function() {
  return gulp.src(config.src)
    .pipe(gulp.dest(config.dest));
});
