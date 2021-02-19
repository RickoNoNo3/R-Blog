/**
 * 显示登录框
 */
function login() {
	BlogPage.PopWindow.openAsInput(
		'login',
		'请点击Cancel',
		'password',
		undefined,
		function (pswd) {
			BlogPage.Ajax.call(
				'/api/admin/login',
				{
					pswd,
				},
				function (data) {
					if (data['text'] === 'ok') {
						location.reload();
					}
				},
			);
		},
	);
}

