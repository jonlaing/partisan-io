var dest = './front_dist';
var src = './front_src';
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
    src: src + '/styles/main.scss',
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
    src: 'font-src/*.html',
    dest: dest
  },
  fonts: {
    src: 'front-src/fonts/*.{ttf,woff,woff2,eot,svg}',
    dest: dest + '/fonts'
  },
  icons: {
    src: './node_modules/font-awesome/fonts/**.*',
    dest: dest + '/fonts'
  },
  images: {
    src: src + '/images/**.*',
    dest: dest + '/images'
  },
  skeleton: {
    src: src + '/styles/skeleton/*.css',
    dest: dest + '/styles/skeleton'
  },
  watch: {
    src: 'front-src/**/*.*',
    tasks: ['build']
  }
};
