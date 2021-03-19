let Editor;
let Output;
let WaitAni;
let InFullscreen = false;

// -----------------------------------------
// CodeMirror 初始化
// -----------------------------------------

(function initCodeMirror() {
// 必须保证逆字典序排列, 因为匹配是顺序执行的
  let CodeSepSrc = [
    ['yml', 'text/x-yaml'],
    ['yaml', 'text/x-yaml'],
    ['xml', 'text/xml'],
    ['vue', 'text/x-vue'],
    ['typescript', 'text/javascript'],
    ['tsx', 'text/typescript-jsx'],
    ['ts', 'text/javascript'],
    ['text', 'text/plain'],
    ['tex', 'text/x-stex'],
    ['shell', 'text/x-sh'],
    ['sh', 'text/x-sh'],
    ['swift', 'text/x-swift'],
    ['sqlite', 'text/x-sql'],
    ['sql', 'text/x-sql'],
    ['scala', 'text/x-scala'],
    ['rust', 'text/x-rust'],
    ['r', 'text/x-rsrc'],
    ['python', 'text/x-python'],
    ['py', 'text/x-python'],
    ['powershell', 'application/x-powershell'],
    ['php', 'text/x-php'],
    ['nginx', 'text/x-nginx-conf'],
    ['mysql', 'text/x-sql'],
    ['mssql', 'text/x-sql'],
    ['md', 'text/x-markdown'],
    ['mathjax', 'text/x-stex'],
    ['markdown', 'text/x-markdown'],
    ['lua', 'text/x-lua'],
    ['latex', 'text/x-stex'],
    ['kotlin', 'text/x-kotlin'],
    ['jsx', 'text/jsx'],
    ['json', 'text/x-yaml'],
    ['js', 'text/javascript'],
    ['javascript', 'text/javascript'],
    ['java', 'text/x-java'],
    ['html', 'text/html'],
    ['golang', 'text/x-go'],
    ['go', 'text/x-go'],
    ['dockerfile', 'text/x-dockerfile'],
    ['docker', 'text/x-dockerfile'],
    ['dart', 'application/dart'],
    ['cython', 'text/x-cython'],
    ['css', 'text/css'],
    ['csharp', 'text/x-csharp'],
    ['cs', 'text/x-csharp'],
    ['cpp', 'text/x-c++hdr'],
    ['cmake', 'text/x-cmake'],
    ['c#', 'text/x-csharp'],
    ['c++', 'text/x-c++hdr'],
    ['c', 'text/x-c'],
  ];

  CodeMirror.defineMode('md', function (config) {
    return CodeMirror.multiplexingMode.apply(this, [
      CodeMirror.getMode(config, 'text/x-gfm'),
      {
        open: '$$',
        close: '$$',
        mode: CodeMirror.getMode(config, 'text/x-stex'),
        delimStyle: 'math',
      },
      {
        open: '$',
        close: '$',
        mode: CodeMirror.getMode(config, 'text/x-stex'),
        delimStyle: 'math',
      },
      ...CodeSepSrc.map(([key, mime]) => ({
        open: '```' + key,
        close: '```',
        mode: CodeMirror.getMode(config, mime),
        delimStyle: 'code',
      })),
      {
        open: '```',
        close: '```',
        mode: CodeMirror.getMode(config, 'text/javascript'),
        delimStyle: 'code',
      },
      {
        open: '`',
        close: '`',
        mode: CodeMirror.getMode(config, 'text/javascript'),
        delimStyle: 'code',
      },
    ]);
  });
})();

// -----------------------------------------
// 功能函数
// -----------------------------------------

function save() {
  if (IS_NEW) {
    BlogPage.Ajax.call(
      '/api/admin/new',
      {
        data: editor.getValue().trim(),
        type: 1,
        dirId: PARENT_ID,
      },
      () => {
        location.href = `/admin/edit?type=0&id=${PARENT_ID}`;
      },
    );
  } else {
    BlogPage.Ajax.call(
      '/api/admin/edit',
      {
        data: editor.getValue().trim(),
        type: 1,
        id: ARTICLE_ID,
      },
      () => {
        BlogPage.PopWindow.openAsNote('saveEnd', '保存成功', '是否返回上级目录?', () => {
          location.href = `/admin/edit?type=0&id=${PARENT_ID}`;
        });
      },
    );
  }
}

function toggleFullscreen() {
  let allContent = document.getElementById('allContent');
  let fullScreenIco = document.querySelector('#toolBar>li[option=fullscreen] i.iconfont');
  if (InFullscreen) {
    allContent.classList.remove('fullscreen');
    fullScreenIco.innerHTML = '&#xe7e4';
    // noinspection JSUnresolvedVariable
    let method = document.cancelFullScreen || document.webkitCancelFullScreen || document.mozCancelFullScreen || document.exitFullScreen || window.myCancelFullScreen || (() => {});
    if (method) method.call(document);
  } else {
    allContent.classList.add('fullscreen');
    fullScreenIco.innerHTML = '&#xe7e3';
    // noinspection JSUnresolvedVariable
    let method = document.body.requestFullScreen || document.body.webkitRequestFullScreen || document.body.mozRequestFullScreen || document.body.msRequestFullScreen || window.myRequestFullScreen || (() => {});
    if (method) method.call(document.body);
  }
  InFullscreen = !InFullscreen;
}


// -----------------------------------------
// 实时预览相关
// -----------------------------------------

// 把iframe传来的渲染好的html移到output里
function refreshDone(newValue, styles) {
  $('head>style[category=mathjax]').remove();
  $('head').append(styles);
  $(Output).html(newValue);
  $(WaitAni).css('opacity', 0);
}

// 建立临时iframe, 渲染好再触发refreshDone
function refresh() {
  setTimeout(() => {
    $(WaitAni).css('opacity', 1);
    // 创建新iframe
    $('body>iframe#preview').remove();
    let iframe = document.createElement('iframe');
    iframe.id = 'preview';
    iframe.style.display = 'none';
    iframe.onload = () => {
      // 添加jquery依赖
      let f = window.frames[0];
      let fDoc = f.document;
      let script = fDoc.createElement('script');
      script.src = `${CDN}js/jquery-3.4.1.min.js`;
      script.addEventListener('load', () => {
        // 添加blogPage依赖
        let script = fDoc.createElement('script');
        script.src = `${CDN}js/blogPage.${USE_MIN_STR}js`;
        script.addEventListener('load', () => {
          // 填充文档, 打造一个可以渲染markdown, 渲染结束后触发信号的页面
          $(fDoc.body).html(`
            <div class="mycontent"></div>
            <script>
            let CDN = '${CDN}';
            let USE_MIN_STR = '${USE_MIN_STR}';
            $(window).on('PageCompleted', ()=>{
              $('head>style').attr('category', 'mathjax');
              window.top.refreshDone($('.mycontent').html(), $('head>style[category=mathjax]'));
            });
            BlogPage.Ext.loadShowdown($('.mycontent'), window.top.Editor.getValue().trim(), () => {
              BlogPage.Ext.loadGroup([
                {
                  func: BlogPage.Ext.loadGraph,
                  args: [$('.mycontent')],
                },
                {func: BlogPage.Ext.loadTex},
                {func: BlogPage.Ext.loadHighlight},
              ], () => {
                PageComplete();
              });
            });
            </script>
          `);
        }, false);
        fDoc.head.appendChild(script);
      }, false);
      fDoc.head.appendChild(script);
    };
    document.body.appendChild(iframe);
  });
}

// 整理文档

$(document).ready(() => {
  // init code mirror
  let textarea = document.querySelector('textarea#editor');
  window.Editor = Editor = CodeMirror.fromTextArea(textarea, {
    mode: 'md',
    theme: 'midnight',
    lineNumbers: true,
    lineWrapping: true,
    // scrollbarStyle: 'overlay',
    // styleActiveLine: true,
    // extraKeys: {
    //   'F2': function () {
    //     toggleOptionList();
    //   },
    //   'F11': function () {
    //     toggleFullscreen();
    //   },
    //   // 'F11': function (cm) {
    //   //    cm.setOption('fullScreen', !cm.getOption('fullScreen'));
    //   // },
    // },
    // vim:
    // keyMap: 'vim',
    // matchBrackets: true,
    // showCursorWhenSelecting: true,
    // tex:
    inMathMode: true,
    // markdown:
    singleCursorHeightPerLine: true,
    highlightFormatting: true,
    fencedCodeBlockHighlighting: true,
    // javascript:
    typescript: true,
    json: true,
  });
  // CodeMirror.commands.save = function () { alert('Saving'); };
  Editor.on('change', () => {
    refresh();
  });

  // init page
  Output = $('#output>div');
  WaitAni = $('#outputRefresh');
  setTimeout(refresh);
  ShowPage();
});
