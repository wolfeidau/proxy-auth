var webpack = require('webpack');
var plugins = [];

plugins.push(
	new webpack.optimize.UglifyJsPlugin({
		compressor: {
			warnings: false
		}
	}),
	new webpack.optimize.DedupePlugin()
);

module.exports = {
  context: __dirname + "/app",
  entry: {
    javascript: "./app.js",
    html: "./index.html",
  },
  plugins: plugins,
  output: {
    filename: "bundle.js",
    path: __dirname + "/dist",
  },
  module: {
    loaders: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loaders: ["babel-loader"],
      },
      {
        test: /\.html$/,
        loader: "file?name=[name].[ext]",
      },
    ],
  },
}
