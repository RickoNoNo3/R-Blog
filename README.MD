# R-Blog Core  (R-Blog 核心)

轻量级博客系统 R-Blog 的功能核心.

Go + Sqlite 3. **同步支持 客户端-服务器模式 和 浏览器-服务器模式 和 本地模式**! 

实现**纯绿色便携博客系统**, 多平台打开即用, 妈妈再也不用担心我换环境了!

-----

需要注意, 启用高级渲染(图论/数学公式/代码高亮)需要安装Chrome浏览器. 如果只阅读, 不添加和修改, 或者不使用高级渲染则不需要.

Linux环境下, 使用Chromium替换Chrome时, 需要手动将Chromium链接到`PATH`下的`google-chrome`文件上, 构造一个假Chrome供程序识别.

## 使用方法

请完整阅读本章.

原则上, 下载本核心程序的Release版本, 执行相应的环境检查, 进行程序配置, 即有一个完整的博客系统可供使用.

然而, 为了提高扩展性, 本程序被设计成以命令行参数来执行博客功能, 因此直接供人类使用并不方便. **建议直接使用以下UI封装版本:**

| UI封装                                                  | 平台依赖       | 描述                                                         | 支持本地博客核心 | 支持连接到网络博客 | 支持全部博客功能 | 支持HTML图形化展示 | 当前核心版本 |
| ------------------------------------------------------- | -------------- | ------------------------------------------------------------ | ---------------- | ------------------ | ---------------- | ------------------ | ------------ |
| [Shell](https://github.com/RickoNoNo3/R-Blog-shell)     | Go 和 SQLite 3 | 基础的命令行实用工具. 提供核心程序调用的替代性写法(基本语法无变化). 同时尽可能格式化核心程序的输出. 适合作为核心程序的后台伴侣. | 是               | -                  | 是               | -                  | 0.4.0        |
| [Website](https://github.com/RickoNoNo3/R-Blog-website) | Go 和 SQLite 3 | Web服务器(带视图层). 提供HTML图形界面和Web API双功能. 可以直接作为博客网站使用, 同时也是网络间多端共用博客的必备环节. | 是               | -                  | 是               | 是                 | 未实现       |
| Desktop                                                 | WPF (.Net 4.5) | 适用于桌面版PC. 以桌面应用程序(GUI)的形式展现所有博客功能. 可选使用本地博客或连接到网络博客. | 是               | 是                 | 是               | 是                 | 未实现       |
| Mobile                                                  | Flutter Mobile | 适用于移动设备. 仅支持连接到网络博客.                        | -                | 是                 | 是               | 是                 | 未实现       |
| ...                                                     |                | *本项目是开源的! 期待你自己的优秀实现!*                      |                  |                    |                  |                    |              |

- 如需搭建Web服务器, 建议直接使用以下架构方式:
	- **在服务器上下载含核心的`Shell`, 不含核心的`Website`**. 使用`Website`建立网络服务器, 并将本地核心路径指向`Shell`内的Blog核心, `Shell`作为后台管理工具备用. **此时已经可以通过网页访问博客**;
	- (可选)**在需要的设备上下载不含核心的`Desktop`, `Mobile`等**, 选择连接到网络博客, 连接上一步中建立的网络服务器. **此时即可通过桌面应用、移动应用来访问博客**.

- 如仅需在桌面端使用, 可以选择直接下载含核心的`Desktop`, 选择使用本地博客.
- 如仅需在移动端使用, 目前不支持. 请换用Web服务器方式.
- 如需为一个不含核心的UI封装启用本地博客, 只要下载本仓库的核心程序, 在UI封装中进行适当配置即可.
- 如需升级版本, 目前需要自行对比新老版本差异进行升级. (TODO: 处理版本升级逻辑).

## 依赖项

### 核心依赖

[go-sqlite3](https://github.com/mattn/go-sqlite3), [sqlx](https://github.com/jmoiron/sqlx)

### Markdown静态渲染

[blackfriday.v2](https://github.com/russross/blackfriday)

### Markdown动态渲染(基于浏览器的JS高级渲染)

[chromedp](https://github.com/chromedp/chromedp), [highlight.js](https://github.com/highlightjs/highlight.js), [mathjax.js](https://github.com/mathjax/MathJax), [katex.js](https://github.com/KaTeX/KaTeX), [viz.js](https://github.com/mdaines/viz.js)