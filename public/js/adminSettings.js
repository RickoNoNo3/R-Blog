$(document).ready(() => {
  ShowPage();
  // 初始化IsInDebug按钮
  $('input:checkbox[name=isInDebug]~a.button').click(function () {
    $('input:checkbox[name=isInDebug]').prop('checked', ($(this).attr('option') === 'true') ? 'checked' : '');
    $(this).siblings().removeAttr('active');
    $(this).attr('active', '1');
  });
  $(`input:checkbox[name=isInDebug] ~ a.button:nth-of-type(${DEBUG_BUTTON_ACTIVE})`).click();
});

function save() {
  const data = {
    m2obj: {
      Version: [
        parseInt($('input[name=version-0]').val().trim()),
        parseInt($('input[name=version-1]').val().trim()),
        parseInt($('input[name=version-2]').val().trim()),
      ].join('.'),
      Blog: {
        CDN: $('input[name=cdn]').val().trim(),
        BGImg: $('input[name=bgImg]').val().trim(),
        IconImg: $('input[name=iconImg]').val().trim(),
        Favicon: $('input[name=favicon]').val().trim(),
      },
      Info: {
        Email: $('input[name=email]').val().trim(),
        QQ: $('input[name=qq]').val().trim(),
        Telegram: $('input[name=telegram]').val().trim(),
      },
      AdminPSWD: $('input[name=adminPSWD]').val().trim(),
      ServerPort: parseInt($('input[name=serverPort]').val().trim()),
      IsInDebug: $('input[name=isInDebug]').prop('checked'),
      LogFile: {
        ConsoleLog: $('input[name=logFile-con]').val().trim(),
        WebLog: $('input[name=logFile-web]').val().trim(),
      },
    },
  };
  BlogPage.Ajax.call(
    '/api/admin/settings/save',
    data,
    res => {
      if (res['res'] === 'ok') {
        location.reload();
      } else {
        BlogPage.PopWindow.openAsNote('restartFailed', '保存失败', '发生未知错误');
      }
    },
  );
}

function reset() {
  BlogPage.PopWindow.openAsNote('resetConfirm', '确定重置', '确定要重置所有配置项吗?<br>如果这和当前环境不匹配, 可能导致系统无法使用!', () => {
    BlogPage.Ajax.call(
      '/api/admin/settings/reset',
      {},
      res => {
        if (res['res'] === 'ok') {
          BlogPage.PopWindow.openAsNote('restartSuccess', '重置成功', '请根据需要重新启动服务器', () => {
            location.reload();
          }, () => {
            location.reload();
          });
        } else {
          BlogPage.PopWindow.openAsNote('restartFailed', '重置失败', '发生未知错误');
        }
      },
    );
  });
}

function restart() {
  BlogPage.Ajax.call(
    '/api/admin/restart',
    {},
    res => {
      if (res['res'] === 'ok') {
        BlogPage.PopWindow.openAsNote('restartSuccess', '服务器重启中', '请稍等片刻，然后刷新页面');
      } else {
        BlogPage.PopWindow.openAsNote('restartFailed', '重启失败', '发生未知错误');
      }
    },
  );
}

function submit(data) {
  console.log(data);
}
