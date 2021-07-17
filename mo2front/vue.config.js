const path = require('path');
module.exports = {
  // chainWebpack: config => {
  //   config.module.rules.delete('eslint');
  // },
  pwa: {
    iconPaths: {
      favicon32: 'https://cdn.mo2.leezeeyee.com/dist/img/icons/favicon-32x32.png',
      favicon16: 'https://cdn.mo2.leezeeyee.com/dist/img/icons/favicon-16x16.png',
      appleTouchIcon: 'https://cdn.mo2.leezeeyee.com/dist/img/icons/tile150x150.png',
      maskIcon: 'https://cdn.mo2.leezeeyee.com/dist/img/icons/tile150x150.png',
      msTileImage: 'https://cdn.mo2.leezeeyee.com/dist/img/icons/tile150x150.png',
    },
    manifestOptions: {
      name: "Mo2",
      short_name: "Mo2",
      start_url: "https://www.motwo.cn",
      display: "standalone",
      description: "超好用又好看的博客软件。速度快，应用小。安装它只需要不到一秒钟的时间。",
      features: [
        "Cross Platform",
        "Fast",
        "Small",
        "Mordern Rich Text Editor"
      ],
      "icons": [{
        "src": "https://cdn.mo2.leezeeyee.com/dist/img/icons/android-chrome-192x192.png",
        "sizes": "192x192",
        "type": "image/png"
      }, {
        "src": "https://cdn.mo2.leezeeyee.com/dist/img/icons/android-chrome-512x512.png",
        "sizes": "512x512",
        "type": "image/png"
      }, {
        "src": "https://cdn.mo2.leezeeyee.com/dist/img/icons/android-chrome-maskable-192x192.png",
        "sizes": "192x192",
        "type": "image/png",
        "purpose": "maskable"
      }, {
        "src": "https://cdn.mo2.leezeeyee.com/dist/img/icons/android-chrome-maskable-512x512.png",
        "sizes": "512x512",
        "type": "image/png",
        "purpose": "maskable"
      }],
      screenshots: [{
          src: "https://cdn.mo2.leezeeyee.com/mobile.png~parallax",
          sizes: "300x200",
          type: "image/png"
        },
        {
          src: "https://cdn.mo2.leezeeyee.com/home.png~parallax",
          sizes: "300x200",
          type: "image/png"
        },
        {
          src: "https://cdn.mo2.leezeeyee.com/home-dark.png~parallax",
          sizes: "300x200",
          type: "image/png"
        },
        {
          src: "https://cdn.mo2.leezeeyee.com/customize.png~parallax",
          sizes: "300x200",
          type: "image/png"
        },
        {
          src: "https://cdn.mo2.leezeeyee.com/editor.png~parallax",
          sizes: "300x200",
          type: "image/png"
        },
        {
          src: "https://cdn.mo2.leezeeyee.com/user-home.png~parallax",
          sizes: "300x200",
          type: "image/png"
        },
      ],
    },
    workboxOptions: {
      skipWaiting: true,
      navigateFallback: 'https://cdn.mo2.leezeeyee.com/dist/index.html',
      runtimeCaching: [{
        urlPattern: new RegExp('^https://www.motwo.cn/api/'),
        handler: 'NetworkFirst',
        options: {
          networkTimeoutSeconds: 1,
          cacheName: 'api-cache',
        },
      }, {
        urlPattern: new RegExp(".*\\.(jpg|png|tif|ico|txt|css|webp)$"),
        handler: 'StaleWhileRevalidate',
        options: {
          cacheName: 'static-cache',
        },
      }],

    }
  },
  transpileDependencies: [
    'vuetify'
  ],
  lintOnSave: true,
  outputDir: path.resolve(__dirname, '../dist'),
  publicPath: process.env.NODE_ENV === 'production' ?
    '//cdn.mo2.leezeeyee.com/dist/' : '/',
  //publicPath: '/static',
  // 放置生成的静态资源 (js、css、img、fonts) 的 (相对于 outputDir 的) 目录。
  //assetsDir: 'static',
  // 指定生成的 index.html 的输出路径 (相对于 outputDir)。也可以是一个绝对路径。
  // indexPath: 'index.html',
  // 默认情况下，生成的静态资源在它们的文件名中包含了 hash 以便更好的控制缓存。然而，这也要求 index 的 HTML 是被 Vue CLI 自动生成的。如果你无法使用 Vue CLI 生成的 index HTML，你可以通过将这个选项设为 false 来关闭文件名哈希。
  filenameHashing: true,
  // 多页面
  pages: undefined,
  // 是否使用包含运行时编译器的 Vue 构建版本。设置为 true 后你就可以在 Vue 组件中使用 template 选项了，但是这会让你的应用额外增加 10kb 左右。
  runtimeCompiler: true,
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
        target: 'http://localhost:5001',
        // target: 'https://limfx.pro',
        ws: true, // 是否启用websockets
        pathRewrite: {
          '^/frontend': ''
        },
        secure: false, // 使用的是http协议则设置为false，https协议则设置为true
        // 开启代理：在本地会创建一个虚拟服务端，然后发送请求的数据，并同时接收请求的数据，这样客户端端和服务端进行数据的交互就不会有跨域问题
        // changOrigin: true,
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
    devtool: 'source-map',
    optimization: {
      splitChunks: {
        minSize: 10000,
        maxSize: 249856,
      }
    }
  }
}