var dest = './dist';
var src = './src';
var gutil = require('gulp-util');

module.exports = {
  server: {
    settings: {
      root: dest,
      host: 'localhost',
      port: 8080,
      livereload: {
        port: 35929
      }
    }
  },
  sass: {
    src: src + '/styles/**/*.{sass,scss,css}',
    dest: dest + '/styles',
    settings: {
      indentedSyntax: false, // Enable .sass syntax?
      imagePath: '/images' // Used by the image-url helper
    }
  },
  browserify: {
    settings: {
      transform: ['babelify', 'reactify']
    },
    src: src + '/js/main.jsx',
    dest: dest + '/js/',
    outputName: 'main.js',
    debug: gutil.env.type === 'dev'
  },
  html: {
    src: 'src/*.html',
    dest: dest
  },
  fonts: {
    src: 'src/fonts/*.{ttf,woff,eot,svg}',
    dest: dest + '/styles'
  },
  watch: {
    src: 'src/**/*.*',
    tasks: ['build']
  }
};
