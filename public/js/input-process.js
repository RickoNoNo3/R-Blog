/**
 * 在Pick等操作前, 必须初始化一个FileInputDOM
 * @param {Function?} callback
 * @param {string?} acceptType 默认为"image/*"
 * @param {HTMLElement?} parentDOM
 * @author DHC
 */
function InitFileInputDOM(callback, acceptType = 'image/*', parentDOM = document.body) {
  let DOM = document.getElementById('FILEINPUT');
  if (DOM === null) {
    $(parentDOM).append($(`<input id="FILEINPUT" style="display:none;" type="file" accept="${acceptType}">`));
  } else {
    $(DOM).attr('accept', acceptType);
  }
  if (typeof callback === 'function') {
    callback();
  }
}

/**
 * 通过inputDOM(默认为"#FILEINPUT")拾取图像,
 * 将图像的base64字符串传入callback.
 *
 * @param {Function?} callback(base64)
 * @param {Function?} callbackError
 * @param {boolean?} withoutZip 禁用压缩, 默认为false
 * @param {HTMLInputElement?} DOM 默认为"#FILEINPUT"
 * @author DHC
 * @modified 2020/4/27
 */
function PickImg(callback, callbackError, withoutZip = true, DOM = document.getElementById('FILEINPUT')) {
  DOM.onchange = function () {
    try {
      if (DOM.files[0]) {
        FileToDataURL(DOM.files[0], function (res) {
          if (typeof callback === 'function') {
            callback(res);
          }
        }, withoutZip);
      } else {
        if (typeof callbackError === 'function') {
          callbackError();
        }
      }
    } catch (e) {
      console.log(e);
    }
  };
  DOM.click();
}

/**
 * 通过inputDOM(默认为"#FILEINPUT")拾取文件,
 * 将文件的名称和二进制数据传入callback.
 *
 * @param {function(name:string, length:number, data:Blob|ArrayBuffer)} callback
 * @param {function?} callbackError
 * @param {boolean} readToBuf 为true则file转为ArrayBuffer, 否则直接返回Blob(File)
 * @param {HTMLInputElement?} DOM 默认为"#FILEINPUT"
 * @author DHC
 * @modified 2021/2/9
 */
function PickFile(callback, callbackError, readToBuf = false, DOM = document.getElementById('FILEINPUT')) {
  DOM.onchange = function () {
    try {
      if (DOM.files[0]) {
        if (readToBuf) {
          FileToBuffer(DOM.files[0], callback, callbackError);
        } else {
          callback(DOM.files[0].name, DOM.files[0].size, DOM.files[0]);
        }
      } else {
        if (typeof callbackError === 'function') {
          callbackError();
        }
      }
    } catch (e) {
      if (typeof callbackError === 'function') {
        callbackError();
      }
    }
  };
  DOM.click();
}

/**
 * 将input的file转换为ArrayBuffer
 */
function FileToBuffer(file, callback, callbackError) {
  var reader = new FileReader();
  reader.onloadend = function (e) {
    try {
      let res = e.target.result;
      if (e.target.error !== null)
        throw 'read failed: ' + e.target.error.message;
      if (typeof callback === 'function')
        callback(file.name, res.byteLength, res);
    } catch (ex) {
      console.log(ex);
      if (typeof callbackError === 'function')
        callbackError();
    }
  };
  reader.readAsArrayBuffer(file);
}

/**
 * 将input的file转换为base64编码(限图片)
 */
function FileToDataURL(file, callback, withoutZip) {
  var reader = new FileReader();
  reader.onloadend = function (e) {
    try {
      let res = e.target.result;
      if (withoutZip) {
        if (typeof callback === 'function')
          callback(res);
        return;
      }
      ZipImgFile(res, {
        zoom: 0.75,
        width: 800,
      }, function (res1) {
        res = res1;
        if (typeof callback === 'function')
          callback(res);
      });
    } catch (ex) {
      console.log(ex);
    }
  };
  reader.readAsDataURL(file);
}

/**
 * @param {String} base64
 * @param {Object?} config
 * @param {Function?} callback
 * 压缩base64格式的图片文件, 通过callback(res)返回压缩后的base64字符串res
 * config中包含width, height, quality(0~1), zoom 4个属性, 可以任意省略.
 * zoom 与width和height有冲突关系, zoom结果大于width时, 会取width
 * 在图片本身宽或高小于50时, 完全屏蔽压缩
 */
function ZipImgFile(base64, config, callback) {
  if (!base64) return base64;
  var img = new Image();
  img.src = base64;
  img.onload = function () {
    var res;
    var ext = base64.substring(base64.indexOf('/') + 1, base64.indexOf(';')).toLowerCase();
    if (this.width > 50 && this.height > 50 && ext !== 'gif') {
      // 默认按比例压缩, 除非同时指定宽和高
      var scale = this.width / this.height;
      var w, h;
      if (config.zoom) {
        w = this.width * config.zoom;
        h = this.height * config.zoom;
      }
      if (!config.zoom || (config.width && w > config.width) || (config.height && h > config.height)) {
        w = config.width || this.width;
        h = config.height || (w / scale);
        if (config.height && !config.width) { // 只指定height时的弥补操作
          w = h * scale;
        }
      }
      //生成canvas
      var canvas = document.createElement('canvas');
      var ctx = canvas.getContext('2d');
      // 创建属性节点
      var anw = document.createAttribute('width');
      var anh = document.createAttribute('height');
      anw.nodeValue = w;
      anh.nodeValue = h;
      canvas.setAttributeNode(anw);
      canvas.setAttributeNode(anh);
      ctx.drawImage(this, 0, 0, w, h);
      res = canvas.toDataURL('image/' + ext, config.quality || 0.8);
      if (res.length >= base64.length) {
        // 反向压缩, 神仙操作
        // 你别说还真有这种, 所以要还原回来
        res = base64;
      }
    } else {
      res = base64;
    }
    console.log('压缩前大小: ' + base64.length + '  |  压缩后大小: ' + res.length);
    if (typeof callback === 'function')
      callback(res);
  };
}

/* 和summernote对接, 压缩富文本文档中的base64图片 */
function ZipImgHTML(htmlDOM, callback) {
  let imgs = $(' img[src*="base64"]:not([class~="drawn"])', htmlDOM);
  if (imgs.length === 0) {
    if (typeof callback === 'function')
      callback(htmlDOM.html());
    return;
  }
  let zipCnt = 0;
  for (let i of imgs) {
    ZipImgFile(imgs[i].src, {
      zoom: 0.75,
      width: 500,
    }, function (base64) {
      imgs[i].src = base64;
      imgs[i].classList.add('drawn');
      if (++zipCnt === imgs.length) {
        if (typeof callback === 'function')
          callback(htmlDOM.html());
      }
    });
  }
}

