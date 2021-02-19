// 判断初始化过程中页面内容是否渲染完成时使用.
// 在最后一个影响文档内容的插件加载完毕后的回调函数中, 调用PageComplete()来给此变量赋值为true.
var PageCompleted = false;

function PageComplete() {
	PageCompleted = true;
	// 给需要等待页面内容渲染完成才能开始的各种任务发送消息
	$(window).trigger('PageCompleted');
	$(window).trigger('BgScrollStart');
}

var BlogPage = {
	Ajax: {
		/**
		 * 封装的Ajax函数
		 * @param url {string}
		 * @param data {any?}
		 * @param success {function(any)?}
		 * @param error {function()?}
		 * @param complete {function()?}
		 * @param dataType {string?}
		 * @param contentType {string?}
		 * @param method {string?}
		 */
		call: (url, data, success, error, complete, dataType = 'json', contentType = 'application/json', method = 'post') => {
			if (typeof data === 'object') {
				data = JSON.stringify(data);
			}
			$.ajax({
				url,
				method,
				contentType,
				data,
				dataType,
				success: function (res) {
					if (typeof success === 'function') {
						success(res);
					} else {
						location.reload();
					}
				},
				error: function () {
					if (typeof error === 'function') {
						error();
					} else {
						BlogPage.PopWindow.openAsNote('callAjaxFailed', '失败', '请检查您的网络连接, 或稍后再试.');
					}
				},
				complete,
			});
		},
	},
	PopWindow: {
		/** 存放于popWindows中的win信息 */
		WinInfo: class {
			/**
			 * 窗口id
			 * @member {string}
			 */
			id;

			/**
			 * 窗口节点
			 * @member {jQuery}
			 */
			node;

			/**
			 * 点击OK按钮时的回调函数
			 * @member {function(jQuery)}
			 */
			callbackOk;

			/**
			 * 点击Cancel按钮时的回调函数
			 * @member {function()}
			 */
			callbackCancel;

			/**
			 * @param id {string}
			 * @param node {jQuery}
			 * @param ok {function(jQuery)}
			 * @param cancel {function()}
			 */
			constructor(id, node, ok, cancel) {
				this.id = id;
				this.node = node;
				this.callbackOk = ok;
				this.callbackCancel = cancel;
			};

			ok = () => {
				if (typeof this.callbackOk === 'function')
					this.callbackOk(this.node.children('.pop-content'));
				BlogPage.PopWindow.close(this.id);
			};

			cancel = () => {
				if (typeof this.callbackCancel === 'function')
					this.callbackCancel();
				BlogPage.PopWindow.close(this.id);
			};
		},
		/**
		 * 已弹出的窗口列表
		 * @type {WinInfo[]}
		 */
		wins: [],
		/**
		 * 弹出一个窗口
		 * @param id {string} 窗口id
		 * @param title {string} 标题
		 * @param html {string|HTMLElement|jQuery} 内嵌的html
		 * @param done {function(jQuery)?} 弹出后的回调函数
		 * @param ok {function(jQuery)?} 点击确定按钮时的回调函数
		 * @param cancel {function()?} 点击取消按钮或背景时的回调函数
		 * @returns {boolean}
		 */
		open: (id, title, html, done, ok, cancel) => {
			if (BlogPage.PopWindow.find(id) !== -1) {
				return false;
			}
			try {
				let win = new BlogPage.PopWindow.WinInfo(id, undefined, ok, cancel);
				let node =
					$('<div class="pop-window"></div>').attr('id', 'PopWindow-' + id).append(
						$('<div class="pop-background"></div>').click(win.cancel),
					).append(
						$('<div class="mybox pop-content"></div>').append(
							$('<p class="pop-title"></p>').text(title),
						).append(
							$('<div class="pop-html"></div>').html(html),
						).append(
							$('<div class="btn-bar"></div>').append(
								$('<a href="javascript:void(0);" class="button colorful">Cancel</a>').click(win.cancel),
							).append(
								$('<a href="javascript:void(0);" class="button colorful">OK</a>').click(win.ok),
							),
						),
					);
				win.node = node;
				$(top.document.body).append(node);
				BlogPage.PopWindow.wins.push(win);
				node.fadeTo(100, 1);
				if (typeof done === 'function')
					done(node.children('.pop-content'));
			} catch (e) {
				console.log(e);
				return false;
			}
			return true;
		},
		/**
		 * 通过窗口id关闭一个弹窗
		 * @param id {string} id
		 * @returns {boolean}
		 */
		close: (id) => {
			let index = BlogPage.PopWindow.find(id);
			if (index !== -1) {
				BlogPage.PopWindow.wins[index].node.fadeTo(100, 0, undefined, function () {
					$(BlogPage.PopWindow.wins[index].node).remove();
					BlogPage.PopWindow.wins.splice(index, 1);
				});
				return true;
			} else {
				return false;
			}
		},
		/**
		 * 通过窗口id寻找一个弹窗
		 * @param id {string} 窗口id
		 * @returns {number} 成功返回其处于BlogPopWindows的下标, 失败返回-1
		 */
		find: (id) => {
			for (let i = 0; i < BlogPage.PopWindow.wins.length; i++) {
				if (BlogPage.PopWindow.wins[i].id === id) {
					return i;
				}
			}
			return -1;
		},
		/**
		 * 建立一个输入框弹窗(open的二次封装)
		 * @param id {string} id
		 * @param title {string} 标题
		 * @param type {string} input的类型
		 * @param text {string?} 如果是文本输入框, 可以填充一个默认值, 可为空
		 * @param ok {function(any)?} 回调函数, 用户确认后触发, 带一个参数, 为用户输入的值
		 * @param cancel {function()?} 回调函数, 用户取消时触发
		 */
		openAsInput: (id, title, type, text, ok, cancel) => {
			BlogPage.PopWindow.open(id, title, (
				$('<input type="' + type + '">')
			), (ele) => {
				let inputEle = ele.find('input');
				inputEle.keypress(function (e) {
					if (e.keyCode === 13) {
						ele.find('.btn-bar>.button[value=OK]').click();
					}
				});
				text = text ?? '';
				inputEle.val(text);
				inputEle.focus();
			}, function (ele) {
				if (typeof ok === 'function')
					ok(ele.find('input').val());
			}, function () {
				if (typeof cancel === 'function')
					cancel();
			});
		},
		/**
		 * 显示一个提示(note)
		 * @param id {string} id
		 * @param title {string} 标题
		 * @param content {string|HTMLElement|jQuery} 提示内容(可以是html)
		 * @param ok {function(jQuery)?} 回调函数, 用户确认后触发, 带有一个JQEle参数
		 * @param cancel {function()?} 回调函数, 用户取消时触发
		 */
		openAsNote: (id, title, content, ok, cancel) => {
			BlogPage.PopWindow.open(id, title, (
				$('<div></div>').html(content)
			), undefined, ok, cancel);
		},
		/**
		 * 显示目录选择器
		 * @param ok {function(string)?} 回调函数, 用户确认后触发, 带有一个dirId参数
		 * @param cancel {function()?} 回调函数, 用户取消时触发
		 */
		openAsDirSelector: (ok, cancel) => {
			BlogPage.Ajax.call(
				'/tool/dirSelector',
				null,
				function (res) {
					BlogPage.PopWindow.open('dirSelector', '目录选择', (
						$(res)
					), null, function (ele) {
						if (typeof ok === 'function') {
							ok(ele.find('.dir-selector .radio[active]').attr('dir-id'));
						}
					}, cancel);
				},
				function () {
					BlogPage.PopWindow.openAsNote('noDirSelector', '请求失败', '无法打开目录选择器.');
				}, undefined, 'html',
			);
		},
	},
	Ext: {
		/**
		 * 已加载的插件列表
		 * @type {HTMLScriptElement[]}
		 */
		exts: [],
		/**
		 * 动态加载一个JS文件(插件)
		 * 并放入BlogExts
		 * @param {string} id
		 * @param {string} src
		 * @param {function()?} callback
		 * @param {boolean?} doCallbackAlready
		 */
		loadJS: (id, src, callback = function () {}, doCallbackAlready = true) => {
			if (BlogPage.Ext.findJS(id) !== -1) {
				if (doCallbackAlready && typeof callback === 'function')
					callback();
				return;
			}
			let script = document.createElement('script');
			script.type = 'text/javascript';
			script.id = id;
			script.src = src;
			script.async = true;
			script.addEventListener('load', () => {
				if (typeof callback === 'function')
					callback();
				BlogPage.Ext.exts.push(script);
			}, false);
			setTimeout(function () {
				document.head.appendChild(script);
			}, 100);
		},
		/**
		 * 通过插件id寻找一个弹窗
		 * @param id {string} 插件id
		 * @returns {number} 成功返回其处于BlogExts的下标, 失败返回-1
		 */
		findJS: (id) => {
			for (let i = 0; i < BlogPage.Ext.exts.length; i++) {
				if (BlogPage.Ext.exts[i].id === id) {
					return i;
				}
			}
			return -1;
		},
		/**
		 * 等待多个异步操作完成, 然后再执行回调函数
		 *
		 * 输入的group必须按照{func, args, callback}[]格式给出, 且不允许省略参数.
		 *
		 * args指除了最后一个callback外的所有参数的数组
		 * @param group {{func:Function, args:any[]?, callback:function()?}[]} 异步函数的集合
		 * @param callback {function()?} 全部完成后的回调函数
		 * @param timeout {number?} 等待加载的延时, 超时直接调用callback, 在加载未结束时跳过此组加载. 0或负数表示不跳过, 即使出现了问题.
		 */
		loadGroup: (group, callback, timeout = 8000) => {
			let all = group.length;
			let completed = 0;
			let called = false;
			if (timeout > 0) {
				setTimeout(() => {
					if (typeof callback === 'function' && !called) {
						// console.log('加载未完成, 触发超时');
						called = true;
						callback();
					}
				}, timeout);
			}
			group.forEach(v => {
				// v.args ??= [];
				// v.callback ??= () => {};
				v.args = v.args ?? [];
				v.callback = v.callback ?? (() => {});
				let tCallback = () => {
					v.callback();
					if (++completed >= all) {
						if (typeof callback === 'function' && !called) {
							// console.log('全部加载完成');
							called = true;
							callback();
						}
					}
				};
				v.args.push(tCallback);
				setTimeout(() => {
					v.func.apply(this, v.args);
				});
			});
		},
		/**
		 * 加载markdown文章, dst应当有undisplay类
		 * @param dst {HTMLElement|jQuery} 渲染结果的放置位置
		 * @param src {string|HTMLElement|jQuery} 文档源, 可以从DOM中获取, 也可直接渲染字符串. 建议用字符串以避免HTML相关转义问题
		 * @param callback {function()?} load完成后的回调函数
		 */
		loadShowdown: (dst, src, callback) => {
			let srcText;
			if (typeof src === 'string') {
				srcText = src;
			} else {
				srcText = $(src).text();
			}
			srcText = srcText.replace(/\\\\/g, `\\\\\\\\`);
			BlogPage.Ext.loadJS('showdown', `${CDN}js/showdown.min.js`, () => {
				let converter = new showdown.Converter();
				converter.setOption('literalMidWordUnderscores', true);
				converter.setOption('tables', true);
				converter.setOption('tasklists', true);
				converter.setOption('smoothLivePreview', true);
				converter.setOption('openLinksInNewWindow', true);
				$(dst).html(converter.makeHtml(srcText));
				Array.from($('h1,h2,h3,h4,h5,h6')).forEach(ele => {
					$(ele).attr('id', $(ele).text());
				});
				$(dst).css('transition', 'opacity 0.5s ease-out');
				$(dst).removeClass('undisplay');
				if (typeof callback === 'function')
					callback();
			});
		},
		/**
		 * 加载代码高亮模块
		 * @param callback {function()?}
		 */
		loadHighlight: (callback) => {
			BlogPage.Ext.loadJS('hljs', `${CDN}js/highlight/highlight.pack.js`, () => {
				hljs.initHighlighting();
				if (typeof callback === 'function') {
					callback();
				}
			});
		},
		/**
		 * 加载图论模块
		 * @param callback {function()?}
		 */
		loadGraph: (callback) => {
			if (typeof callback === 'function') {
				callback();
			}
		},
		/**
		 * 加载数学公式模块
		 * @param callback {function()?}
		 */
		loadTex: (callback) => {
			window.MathJax = {
				tex: {
					inlineMath: [['$', '$']],
					displayMath: [['$$', '$$']],
				},
				startup: {
					pageReady: () => {
						return MathJax.startup.defaultPageReady().then(() => {
							let mjStyleNode = document.getElementById('MJX-CHTML-styles');
							// 替换字体路径为CDN资源
							mjStyleNode.innerHTML = mjStyleNode.innerHTML.replace(/"[^"]*(woff-v2\/[^"]*)"/g, `"${CDN}fonts/$1"`);
							if (typeof callback === 'function') {
								callback();
							}
						});
					},
				},
			};
			BlogPage.Ext.loadJS('mathJax', `${CDN}js/mathjax/tex-chtml.js`, () => {});
		},
		/**
		 * 加载静态js文件, 典型如几个con*.js
		 * @param callback {function()?} load完成时的回调函数
		 */
		loadStaticFiles: (callback) => {
			BlogPage.Ext.loadGroup([
				{
					func: BlogPage.Ext.loadJS,
					args: ['conBg', `${CDN}js/conBg.${USE_MIN_STR}js`],
				},
				{
					func: BlogPage.Ext.loadJS,
					args: ['conLogin', `${CDN}js/conLogin.${USE_MIN_STR}js`],
					callback: () => {
						BlogPage.Ext.loadJS('conNavLoc', `${CDN}js/conNavLoc.${USE_MIN_STR}js`);
					},
				},
			], callback);
		},
	},
};
