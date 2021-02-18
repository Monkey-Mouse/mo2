module.exports = {
  chainWebpack: config => {
    config.module.rules.delete('eslint');
  },

  transpileDependencies: [
    'vuetify'
  ],
  lintOnSave: true,
  //outputDir: path.resolve(__dirname, '../../back/LoveCraft.Kshub/wwwroot'),
  // 放置生成的静态资源 (js、css、img、fonts) 的 (相对于 outputDir 的) 目录。
  // assetsDir: 'static',
  // 指定生成的 index.html 的输出路径 (相对于 outputDir)。也可以是一个绝对路径。
  // indexPath: 'index.html',
  // 默认情况下，生成的静态资源在它们的文件名中包含了 hash 以便更好的控制缓存。然而，这也要求 index 的 HTML 是被 Vue CLI 自动生成的。如果你无法使用 Vue CLI 生成的 index HTML，你可以通过将这个选项设为 false 来关闭文件名哈希。
  filenameHashing: true,
  // 多页面
  pages: undefined,
  // 是否使用包含运行时编译器的 Vue 构建版本。设置为 true 后你就可以在 Vue 组件中使用 template 选项了，但是这会让你的应用额外增加 10kb 左右。
  runtimeCompiler: false,
  // 默认情况下 babel-loader 会忽略所有 node_modules 中的文件。如果你想要通过 Babel 显式转译一个依赖，可以在这个选项中列出来。
  transpileDependencies: [],
  // 如果你不需要生产环境的 source map，可以将其设置为 false 以加速生产环境构建。
  productionSourceMap: false,
  // 设置生成的 HTML 中 <link rel="stylesheet"> 和 <script> 标签的 crossorigin 属性。需要注意的是该选项仅影响由 html-webpack-plugin 在构建时注入的标签 - 直接写在模版 (public/index.html) 中的标签不受影响。
  crossorigin: undefined,
  // 在生成的 HTML 中的 <link rel="stylesheet"> 和 <script> 标签上启用 Subresource Integrity (SRI)。如果你构建后的文件是部署在 CDN 上的，启用该选项可以提供额外的安全性。需要注意的是该选项仅影响由 html-webpack-plugin 在构建时注入的标签 - 直接写在模版 (public/index.html) 中的标签不受影响。另外，当启用 SRI 时，preload resource hints 会被禁用，因为 Chrome 的一个 bug 会导致文件被下载两次。
  integrity: false,
  // 反向代理
  devServer: {
    // http2:true,
    https: false,
    proxy: {
      '/api': {
        // 要访问的跨域的域名
        target: 'http://localhost:5000',
        // target: 'https://limfx.pro',
        ws: true, // 是否启用websockets
        pathRewrite: {
          '^/frontend': ''
        },
        secure: true, // 使用的是http协议则设置为false，https协议则设置为true
        // 开启代理：在本地会创建一个虚拟服务端，然后发送请求的数据，并同时接收请求的数据，这样客户端端和服务端进行数据的交互就不会有跨域问题
        changOrigin: true,
        cookieDomainRewrite: "localhost"
      },
      // '/img': {
      //     // 要访问的跨域的域名
      //     // target: 'http://localhost:8010',
      //     target: 'https://limfx.pro',
      //     ws: true, // 是否启用websockets
      //     secure: true, // 使用的是http协议则设置为false，https协议则设置为true
      //     pathRewrite: { '^/frontend': '' },
      //     // 开启代理：在本地会创建一个虚拟服务端，然后发送请求的数据，并同时接收请求的数据，这样客户端端和服务端进行数据的交互就不会有跨域问题
      //     changOrigin: true
      // }
    }
  },
  configureWebpack: {
    devtool: 'source-map'
  }
}