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
      }
    ]
  },
};
