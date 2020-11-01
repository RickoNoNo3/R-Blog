var LoadedModules = [];

function LoadScript(src, callback) {
	setTimeout(function () {
		if (!callback)
			callback = function () { };
		if (LoadedModules.indexOf(src) == -1) {
			var script = document.createElement('script'),
				head = document.getElementsByTagName('head')[0];
			script.type = 'text/javascript';
			script.charset = 'UTF-8';
			script.src = src;
			if (script.addEventListener) {
				script.addEventListener('load', function () {
					callback();
					LoadedModules.push(src);
				}, false);
			} else if (script.attachEvent) {
				script.attachEvent('onreadystatechange', function () {
					var target = window.event.srcElement;
					if (target.readyState == 'loaded') {
						callback();
						LoadedModules.push(src);
					}
				});
			}
			head.appendChild(script);
		} else {
			callback();
		}
	}, 100);
}
