{
//-------------------
//  弹出菜单控制(应放于菜单项管理脚本之后)
//-------------------
  // 确定列表项长度, 并设置给css变量来控制样式
  $(document.body).css('--menu-cnt', $('.mynav>.mynavline>.mymenubar>.mymenuli').length);
  // 初始化
  const mynav = $('.mynav');
  mynav.css('opacity', '1');
  // 替换点击事件
  const links = mynav.find('a');
  links.each((index, link) => {
    $(link).data('href', $(link).attr('href'));
    $(link).attr('href', 'javascript:void(0)');
    $(link).click(function () {
      if (mynav[0].offsetTop === 0 && !mynav.hasClass('preview')) {
        window.location.href = $(this).data('href');
      }
    });
  });

  $(window).on('touchstart', e=>e.preventDefault());
  $(window).on('NavMenuPreview', ()=>{
    setTimeout(() => {
      mynav.addClass('preview');
      setTimeout(()=>{
        mynav.removeClass('preview');
      }, 1000);
    }, 1000);
  });

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
