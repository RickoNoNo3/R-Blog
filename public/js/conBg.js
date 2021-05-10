//-------------------
//  背景图相关
//-------------------
(function () {
	var BG = document.getElementById('BG');
	var imgHeight = 1080;
	var imgWidth = 1386;
	var imgRatio = imgWidth / imgHeight;

	function getSeenRatio() {
		var seenHeight = window.innerHeight;
		var seenWidth = window.innerWidth;
		return seenWidth / seenHeight;
	}

	function bgScroll() {
		// 页面没加载完时, 背景不要滚动
		if (!PageCompleted) {
			return;
		}
		// 注意, top的设置值和偏移量老是反的, 因为是向上偏移
		// ----------------------
		// 当前滚动位置
		var scrollTop = Math.max(document.body.scrollTop, document.documentElement.scrollTop);
		// 视窗宽度
		var seenWidth = window.innerWidth;
		// 视窗高度
		var seenHeight = window.innerHeight;
		// 文档内容总高度
		var allHeight = document.body.scrollHeight;
		// 图片最大顶部偏移量 (有实际大小就是图片自高-视窗高; 否则使用图片预高(带缩放)-视窗高)
		var maxTop = (BG.offsetHeight !== 0 ? BG.offsetHeight : (seenWidth / imgWidth * imgHeight)) - seenHeight;
		// 错误强制修复(当前偏移量超过最大偏移量底部留黑 || 当前偏移量是负数顶部留黑)
		(function () {
			var nowTop = -parseInt(BG.style.top.replace(/px$/g, ''));
			// console.log(nowTop);
			// console.log(maxTop);
			if (nowTop > maxTop || nowTop < 0) {
				BG.style.top = '0px';
			}
		})();
		// 只在横屏时滚动, 竖屏保持在左上角
		if (getSeenRatio() > 1) {
			// 内容超过一屏, 按当前阅读比例设定偏移量
			// 内容不超过一瓶, 通过动画效果滚动到底部
			if (Math.abs(allHeight - seenHeight) >= 10) { // allHeight和seenHeight可能有轻微偏差
				var newTop = ((scrollTop / (allHeight - seenHeight)) * 1.1 * maxTop); // 1.1倍可加快露出猫猫, 不过需要和下一句溢出判断来配合
				BG.style.top = '-' + (newTop < maxTop ? newTop : maxTop) + 'px'; // 最多滚动到底部
			} else {
				BG.style.top = -maxTop + 'px';
			}
		} else {
			BG.style.top = '0px';
		}
	}

	function bgResize() {
		var seenHeight = window.innerHeight;
		var seenWidth = window.innerWidth;
		var seenRatio = seenWidth / seenHeight;

		// 可视宽高比 <= 图像宽高比, 高度优先
		// 可视宽高比 >  图像宽高比, 宽度优先
		if (seenRatio <= imgRatio) {
			BG.style.height = '100vh';
			BG.style.width = 'unset';
		} else {
			BG.style.height = 'unset';
			BG.style.width = '100vw';
		}
		bgScroll();
	}

	$(window).on('scroll', bgScroll);
	$(window).on('resize', bgResize);
	$(window).on('BgScrollStart', bgScroll);

	bgResize();
	setInterval(function () {
		bgResize();
	}, 60000);

	function showBg() {
		BG.onload(new Event('must', {}));
		BG.style.display = 'block';
		BG.style.opacity = '0.7';
		bgScroll();
	}

	if (BG.getAttribute('complete') === '1') {
		showBg();
	} else {
		setTimeout(() => {
			showBg();
		}, 3000);
	}
})();
