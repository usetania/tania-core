const isProduction = process.env.NODE_ENV === 'production';

module.exports = async ({ config, mode }) => {
  config.module.rules.push(
    {
      test: /\.pug$/,
      oneOf: [
        // this applies to `<template lang="pug">` in Vue components
        {
          resourceQuery: /^\?vue/,
          use: ['pug-plain-loader'],
        },
        // this applies to pug imports inside JavaScript
        {
          use: ['raw-loader', 'pug-plain-loader'],
        },
      ],
    },
    {
      test: /\.(sa|sc)ss$/,
      use: [{
        loader: 'vue-style-loader',
      }, {
        loader: isProduction ? MiniCssExtractPlugin.loader : 'style-loader',
      }, {
        loader: 'css-loader',
      }, {
        loader: 'sass-loader'
      }],
    },
  );

  return config;
};
