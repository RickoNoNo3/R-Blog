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

// TODO

function save() {

}

function reset() {

}

function restart() {
  BlogPage.Ajax.call(
    '/api/admin/restart',
    {
      data: Editor.getValue().trim(),
      type: 1,
      dirId: PARENT_ID,
    },
    res => {
      if (res['res'] === 'ok') {
        location.href = `/admin/edit?type=0&id=${PARENT_ID}`;
      } else {
        BlogPage.PopWindow.openAsNote('saveFailed', '保存失败', '发生未知错误');
      }
    },
  );
}

function submit(data) {
  console.log(data);
}
