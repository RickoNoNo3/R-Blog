{
//-------------------
//  弹出菜单控制(应放于菜单项管理脚本之后)
//-------------------
	// 确定列表项长度, 并设置给css变量来控制样式
	$(document.body).css('--menu-cnt', $('.mynav>.mynavline>.mymenubar>.mymenuli').length);
	// 初始化
	let mynav = document.getElementsByClassName('mynav')[0];
	mynav.style.opacity = 1;
	// 替换点击事件, 由menuClick接管
	let linkHrefs = {};
	let links = mynav.getElementsByTagName('a');
	for (let i = 0; i < links.length; ++i) {
		linkHrefs[i] = links[i].href;
		links[i].href = 'javascript:void(0);';
		links[i].addEventListener('click', function () {
			menuClick(i);
		});
		links[i].setAttribute('onclick', `menuClick(${i});`);
	}

	// 为触屏专门适配过的点击事件处理
	// 仅在完全弹出菜单后点击按钮, 才会跳转
	function menuClick(index) {
		// 判断是否完整弹出
		if (mynav.offsetTop === 0) {
			window.location.href = linkHrefs[index];
		}
	}

	window.addEventListener('touchstart', function () { }, true);

//-------------------
//  地址栏相关
//-------------------
	let locPrefix = '当前位置: ';
	let locDOM = document.getElementsByClassName('myloc')[0];
	let locLinkDOM = document.getElementById('loclink');
	let locLink_full = locLinkDOM.innerHTML.replace(locPrefix, '').trim();

	let locList = locLink_full.split(' &gt; ');

	// 绘制locLink栏, 将会根据显示宽度自动决定内容
	// 可能为如下几种:
	//   1. 当前位置: 博客 > dir1 > dir2 > dir3
	//   2. 博客 > dir1 > dir2 > dir3
	//   3. ... > dir2 > dir3 (类似, 能显示几级显示几级, 最少显示当前文件夹)
	function locRedraw() {
		locLinkDOM.innerHTML = locPrefix + locLink_full;
		if (checkLocLength()) {
			for (let i = 0; i < locList.length; ++i) {
				if (!i) {
					locLinkDOM.innerHTML = '';
				} else {
					locLinkDOM.innerHTML = locList[i - 1].replace(/(<a.*>).*(<\/a>)/g, '$1...$2');
				}
				for (let j = i; j < locList.length; ++j) {
					locLinkDOM.innerHTML += (j ? ' &gt; ' : '') + locList[j];
				}
				if (!checkLocLength())
					break;
			}
		}
	}

	function checkLocLength() {
		return (locDOM.offsetWidth > window.innerWidth);
	}

	window.addEventListener('resize', locRedraw);
	locRedraw();

	// initial anime
	locDOM.style.opacity = '1';
	locDOM.style.left = '0';
}
