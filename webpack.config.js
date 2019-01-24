const path = require('path')
const { VueLoaderPlugin } = require('vue-loader')

module.exports = {
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'resources/js')
    }
  },
  module: {
    rules: process.env.NODE_ENV === 'test'
      ? [
        {
          test: /\.vue$/,
          use: 'vue-loader'
        },
        {
          test: /\.hbs$/,
          use: {
            loader: 'handlebars-loader'
          }
        },
        {
          test: /\.pug$/,
          oneOf: [
            // this applies to `<template lang="pug">` in Vue components
            {
              resourceQuery: /^\?vue/,
              use: ['pug-plain-loader']
            },
            // this applies to pug imports inside JavaScript
            {
              use: ['raw-loader', 'pug-plain-loader']
            }
          ]
        },
        {
          test: /\.js$/,
          exclude: /node_modules/,
          use: {
            loader: 'babel-loader',
            options: {
              presets: ['@babel/env']
            }
          }
        },
        {
          test: /\.(png|jpg|gif|svg)$/,
          loader: 'file-loader',
          options: {
            name: '[name].[ext]?[hash]'
          }
        },
        {
          test: /\.(css|sass|scss)$/,
          loader: ['raw-loader', 'sass-loader']
        }
      ]
      : []
  },
  plugins: [
    new VueLoaderPlugin()
  ]
}

// test specific setups
if (process.env.NODE_ENV === 'test') {
  module.exports.externals = [require('webpack-node-externals')()]
  module.exports.devtool = 'inline-cheap-module-source-map'
}
