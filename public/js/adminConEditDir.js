$(document).ready(() => {
  ShowPage();
});

// -----------------------------------------------
// Select
// -----------------------------------------------

function toggleSelect(type, id) {
  var ele = $(`.list > .list-r[entity-type=${type}][entity-id=${id}]`);
  var icon = ele.find('.checkbox-d>i');
  // 动画, 等待延迟
  icon.fadeOut(50, function () {
    // 改checkbox
    if (typeof ele.attr('checked') !== 'undefined') {
      icon.html('&#xe600;');
      ele.removeAttr('checked');
    } else {
      icon.html('&#xe601;');
      ele.attr('checked', '1');
    }
    // 动画, 等待延迟
    icon.fadeIn(200);
    // 改toolBar
    if (getSelectList().length === 0) {
      $('#toolBar>li[option=move], #toolBar>li[option=remove]').removeClass('show');
    } else {
      $('#toolBar>li[option=move], #toolBar>li[option=remove]').addClass('show');
    }
  });
}

function getSelectList() {
  var list = $('.list>.list-r[checked]');
  var data = [];
  if (!list || list.length === 0) return data;
  list.each(function () {
    var ele = $(this);
    data.push({
      type: parseInt(ele.attr('entity-type')),
      id: parseInt(ele.attr('entity-id')),
    });
  });
  return data;
}

// -----------------------------------------------
// Jump
// -----------------------------------------------

function jump(type, id) {
  switch (type) {
    case 0:
      location.href = `/admin/edit?type=0&id=${id}`;
      break;
    case 1:
      window.open(`/blog/article/${id}`, '_blank');
      break;
    case 2:
      window.open(`/blog/file/${id}`, '_blank');
      break;
  }
}

// -----------------------------------------------
// Edit/New
// -----------------------------------------------

function popNameInput(popId, popTitle, oldText = '') {
  return new Promise(resolve => {
    BlogPage.PopWindow.openAsInput(popId, popTitle, 'text',
      oldText,
      (data) => {
        data = data.trim();
        if (data === '') {
          BlogPage.PopWindow.openAsNote(`${popId}Fail`, '失败', '输入无效!');
        } else {
          resolve(data);
        }
      },
    );
  });
}

function newDir() {
  popNameInput('dirNameInput', '输入目录名').then(data => {
    BlogPage.Ajax.call(
      '/api/admin/new',
      {
        data,
        type: 0,
        dirId: DIR_ID,
      },
    );
  });
}

function editDir(id) {
  var ele = $(`.list>.list-r[entity-type=0][entity-id=${id}] .name-d i+span`);
  popNameInput('dirNameInput', '输入目录名', ele.text()).then(data => {
    BlogPage.Ajax.call(
      '/api/admin/edit',
      {
        data,
        type: 0,
        id,
      },
    );
  });
}

function newArticle() {
  location.href = `/admin/edit?type=1&id=-1&parentId=${DIR_ID}`;
}

function editArticle(id) {
  location.href = `/admin/edit?type=1&id=${id}`;
}

function newFile() {
  BlogPage.PopWindow.openAsFileUpload(() => {
    location.reload();
  }, undefined, DIR_ID, false, true);
}

function editFile(id) {
  var ele = $(`.list>.list-r[entity-type=2][entity-id=${id}] .name-d i+span`);
  popNameInput('fileNameInput', '输入文件名', ele.text()).then(data => {
    BlogPage.Ajax.call(
      '/api/admin/edit',
      {
        data,
        type: 2,
        id,
      },
    );
  });
}

function edit(type, id) {
  switch (type) {
    case 0:
      return editDir(id);
    case 1:
      return editArticle(id);
    case 2:
      return editFile(id);
  }
}

// -----------------------------------------------
// Remove
// -----------------------------------------------

function remove_inner(list) {
  BlogPage.PopWindow.openAsNote('remove', '确认删除', '删除操作是永久的<br/>删除目录会连同所有内容一并删除!', () => {
    BlogPage.Ajax.call('/api/admin/remove', {list});
  });
}

function removeOne(type, id) {
  remove_inner([{type, id}]);
}

function removeList() {
  remove_inner(getSelectList());
}

// -----------------------------------------------
// Move
// -----------------------------------------------

function move_inner(list) {
  BlogPage.PopWindow.openAsDirSelector(dirId => {
    BlogPage.Ajax.call('/api/admin/move', {list, dirId});
  });
}

function moveOne(type, id) {
  move_inner([{type, id}]);
}

function moveList() {
  move_inner(getSelectList());
}
