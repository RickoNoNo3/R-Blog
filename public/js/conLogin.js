/**
 * 显示登录框
 */
const _loginFailedWithOkNote = [
  '都说了要点取消啦!',
  '怎么不听话, 请点击取消!',
  '请点击取消.',
  '希望你能点击取消.',
  '求求你点击取消.',
  '相信我, 点击取消可以发现新世界.',
  '请马上点击取消!',
  '为什么不点击取消?',
  '快点, 点击取消!',
  '点击取消!',
];
const _loginFailedWithOk = () => {
  BlogPage.PopWindow.openAsNote(
    `loginFailed${Math.floor(Math.random() * 100000000)}`,
    '请点击取消',
    `${_loginFailedWithOkNote[Math.floor(Math.random() * 10)]}`,
    () => {
      _loginFailedWithOk();
    },
  );
};

function login() {
  BlogPage.PopWindow.openAsInput(
    'login',
    '请点击取消',
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
