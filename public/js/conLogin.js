/**
 * 显示登录框
 */
const _loginFailedWithOkNote = [
  '都说了要点Cancel啦!',
  '怎么不听话, 请点击Cancel!',
  '请点击Cancel.',
  '希望你能点击Cancel.',
  '求求你点击Cancel.',
  '相信我, 点击Cancel可以发现新世界.',
  '请马上点击Cancel!',
  '为什么不点击Cancel?',
  '快点, 点击Cancel!',
  '点击Cancel!',
];
const _loginFailedWithOk = () => {
  BlogPage.PopWindow.openAsNote(
    `loginFailed${Math.floor(Math.random() * 100000000)}`,
    '请点击Cancel',
    `${_loginFailedWithOkNote[Math.floor(Math.random() * 10)]}`,
    () => {
      _loginFailedWithOk();
    },
  );
};

function login() {
  BlogPage.PopWindow.openAsInput(
    'login',
    '请点击Cancel',
    'password',
    undefined,
    pswd => {
      BlogPage.Ajax.call(
        '/api/login',
        {pswd},
        data => {
          if (data['res'] === 'ok') {
            location.reload();
          } else {
            _loginFailedWithOk();
          }
        },
      );
    },
  );
}

function logout() {
  BlogPage.PopWindow.openAsNote('logout', '退出登录', '确定要退出登陆吗?', () => {
    BlogPage.Ajax.call(
      '/api/logout',
      undefined,
      () => {
        location.reload();
      },
    );
  });
}