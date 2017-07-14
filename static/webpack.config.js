const webpack = require('webpack');
const BabiliPlugin = require("babili-webpack-plugin");

if (process.env.NODE_ENV && !['production', 'development'].includes(process.env.NODE_ENV)) {
  console.log("錯誤：環境變數 NODE_ENV 僅能爲 production | development");
  process.exit(1);
}
if (process.env.silent && !['true', 'flase'].includes(process.env.silent)) {
  console.log("錯誤：環境變數 silent 僅能爲 true | false");
  process.exit(1);
}

module.exports = {
  entry: './jsx/basepage.jsx',
  output: {
    filename: './dist/app.js'       
  },
  module: {
    loaders: [
      {
        test: /\.js|jsx$/,
        loader: 'babel',
        query:
          {
            presets:['react']
          }
      },
      {
        test: /less$/,
        loader: 'style-loader!css-loader!less-loader'
      },
      {
        test: /jpg$/,
        loader: 'url-loader'
      }
    ]
  },
  plugins: (() => {
    let ret = [
      new webpack.DefinePlugin({
        'process.env.NODE_ENV': `"${process.env.NODE_ENV || 'development'}"`,
        'process.env.silent': `${process.env.silent == 'true' ? true : false}`,
      }),
    ];
    if (process.env.NODE_ENV == 'production') {
      ret.push(new BabiliPlugin());
    }
    return ret;
  })()
};
