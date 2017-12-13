const mix = require('laravel-mix');
const webpack = require('webpack');
const path = require('path');
const glob = require('glob');

const HtmlWebpackPlugin = require('html-webpack-plugin');
const PurifyCSSPlugin = require('purifycss-webpack');

mix.webpackConfig({
  output: {
    chunkFilename: mix.inProduction() ? 'js/build/[name].[chunkhash].js' : 'js/build/[name].js',
    publicPath: '/'
  },
  plugins: [
    new webpack.optimize.CommonsChunkPlugin({
      name: '/js/vendor',
      minChunks: (module) => {
        return module.context && module.context.indexOf('node_modules') !== -1
      }
    }),
    new PurifyCSSPlugin({
      paths: glob.sync(path.join(__dirname, 'resources/assets/js/**/*.vue')),
    })
  ],
  resolve: {
    extensions: [
      ".vue"
    ]
  }
})

mix.js('resources/js/app.js', 'public/js')
   .sass('resources/sass/app.scss', 'public/css');
