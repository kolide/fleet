require('es6-promise').polyfill();

const path = require('path');
const webpack = require('webpack');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const bourbon = require('node-bourbon').includePaths;
const HtmlWebpackPlugin = require('html-webpack-plugin');
const WebpackNotifierPlugin = require('webpack-notifier');
const repo = __dirname;

let plugins = [
  new webpack.NoEmitOnErrorsPlugin(),
  new HtmlWebpackPlugin({
    filename: '../frontend/templates/react.tmpl',
    inject: false,
    template: 'frontend/templates/react.ejs',
  }),
  new WebpackNotifierPlugin({
    title: 'Kolide',
    contentImage: path.join(
      repo,
      'assets',
      'images',
      'kolide-logo-icon@3x.png'
    ),
    excludeWarnings: true,
  }),
];

if (process.env.NODE_ENV === 'production') {
  plugins = plugins.concat([
    new webpack.optimize.UglifyJsPlugin({
      compress: { warnings: false },
      output: { comments: false },
    }),
    new webpack.DefinePlugin({
      'process.env': { NODE_ENV: JSON.stringify('production') },
    }),
    new ExtractTextPlugin({
      filename: 'bundle-[contenthash].css',
      allChunks: false,
    }),
  ]);
} else {
  plugins = plugins.concat([
    new ExtractTextPlugin({ filename: 'bundle.css', allChunks: false }),
  ]);
}

const config = {
  entry: {
    bundle: path.join(repo, 'frontend/index.jsx'),
  },
  output: {
    path: path.join(repo, 'assets/'),
    publicPath: '/assets/',
    filename: '[name].js',
  },
  plugins,
  module: {
    // The following noParse suppresses the warning about sqlite-parser being a
    // pre-compiled JS file. See https://goo.gl/N4s6bB.
    noParse: /node_modules\/sqlite-parser\/dist\/sqlite-parser-min.js/,
    rules: [
      {
        test: /\.(png|gif)$/,
        use: { loader: 'url-loader?name=[name]@[hash].[ext]&limit=6000' },
      },
      {
        test: /\.(pdf|ico|jpg|svg|eot|otf|woff|ttf|mp4|webm)$/,
        use: { loader: 'file-loader?name=[name]@[hash].[ext]' },
      },
      { test: /\.json$/, use: { loader: 'raw-loader' } },
      {
        test: /\.tsx?$/,
        exclude: /node_modules/,
        use: { loader: 'ts-loader' },
      },
      {
        test: /\.scss$/,
        exclude: /node_modules/,
        use: ExtractTextPlugin.extract({
          fallback: 'style-loader',
          use: [
            { loader: 'css-loader' },
            {
              loader: 'sass-loader',
              options: {
                sourceMap: true,
                includePaths: [bourbon],
              },
            },
            { loader: 'import-glob-loader' },
          ],
        }),
      },
      {
        test: /\.css$/,
        use: ExtractTextPlugin.extract({
          fallback: 'style-loader',
          use: 'css-loader',
        }),
      },
      {
        test: /\.(js|jsx)$/,
        include: path.join(repo, 'frontend'),
        use: {
          loader: 'babel-loader?cacheDirectory',
          options: {
            presets: [
              [
                '@babel/preset-env',
                {
                  modules: false,
                  targets: {
                    browsers: '> 1%',
                  },
                  useBuiltIns: 'entry',
                },
              ],
              '@babel/preset-react',
            ],
            plugins: [
              '@babel/plugin-syntax-dynamic-import',
              '@babel/plugin-transform-spread',
              '@babel/plugin-proposal-object-rest-spread',
              '@babel/plugin-proposal-export-default-from',
              [
                '@babel/plugin-proposal-class-properties',
                {
                  spec: true,
                },
              ],
            ],
            env: {
              test: {
                plugins: ['@babel/plugin-transform-modules-commonjs'],
              },
            },
          },
        },
      },
    ],
  },
  resolve: {
    extensions: ['.js', '.jsx', '.json'],
    modules: [path.resolve(path.join(repo, './frontend')), 'node_modules'],
  },
};

if (process.env.NODE_ENV === 'production') {
  config.output.filename = '[name]-[hash].js';
}

module.exports = config;
