var ticket;
//这个是获取用户的信息
$.ajax({
    async: false,
    type: 'post',
    data: {'url': location.href},
    url: '/share/get_user_token',
    success: function (data) {
        alert(data)
    }
});
//这个是上传图片到本地服务器
$.ajax({
    async: false,
    type: 'post',
    data: {'': location.href},
    url: '/share/get_user_token',
    success: function (data) {
        alert(data)
    }
});
//这个是获取调用微信的ticket
$.ajax({
    async: false,
    type: 'get',
    url: '/share/get_ticker',
    success: function (data) {
        ticket = data;
    }
});
//获取config中的signature
var s1 = 'jsapi_ticket=' + ticket + '&noncestr=qwertyuiop&timestamp=1414587457&url=' + location.href.split('#')[0]
var signature = hex_sha1(s1)
console.log(signature)
console.log(location.href.split('#')[0])
wx.config({
    debug: true, // 开启调试模式,调用的所有api的返回值会在客户端alert出来，若要查看传入的参数，可以在pc端打开，参数信息会通过log打出，仅在pc端时才会打印。
    appId: 'wx53d52d70ccd6439f', // 必填，公众号的唯一标识
    timestamp: 1414587457, // 必填，生成签名的时间戳
    nonceStr: 'qwertyuiop', // 必填，生成签名的随机串
    signature: signature,// 必填，签名，见附录1
    jsApiList: ['onMenuShareTimeline',
        'chooseImage',
        'startRecord',
        'stopRecord',
        'playVoice',
        'onMenuShareAppMessage',
        'getLocation',
        'onMenuShareQQ',
        'chooseWXPay',
        'uploadImage'
    ] // 必填，需要使用的JS接口列表，所有JS接口列表见附录2
});
wx.ready(function () {
//检查各种接口
    wx.checkJsApi({
        jsApiList: ['onMenuShareTimeline', 'chooseImage', 'getLocalImgData'], // 需要检测的JS接口列表，所有JS接口列表见附录2,
        success: function (res) {
            // 以键值对的形式返回，可用的api值true，不可用为false
            // 如：{"checkResult":{"chooseImage":true},"errMsg":"checkJsApi:ok"}

        }
    });
//只有触发分享到朋友圈就会调用这个方法
    wx.onMenuShareTimeline({
        title: '分享朋友圈测试，敢点进来我就弄死你！！！', // 分享标题
        link: 'http://service.genyuanlian.com/seven_nightredirecturl', // 分享链接，该链接域名或路径必须与当前页面对应的公众号JS安全域名一致
        imgUrl: 'https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1514438536365&di=09f25c134d4b0589e136a5fddc22a22c&imgtype=0&src=http%3A%2F%2Fwww.haha365.com%2Fuploadfile%2F2014%2F0404%2F20140404063425613.jpg', // 分享图标
        success: function () {
            alert("chengegfadfa");

        },
        cancel: function () {
            // 用户取消分享后执行的回调函数
            /*alert("quxiaofenxiang");*/

        }
    });
//分享给朋友调用的方法
    wx.onMenuShareAppMessage({
        title: '这他么是个测试',
        desc: '还有谁？？？还有谁？？？我就问还有谁',
        link: 'http://service.genyuanlian.com/seven_nightredirecturl',
        imgUrl: 'https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1514438682725&di=7624aa8b34c3d92a2d2b95254aff2450&imgtype=0&src=http%3A%2F%2Fvpic.video.qq.com%2F3170316%2Fb03322rccwa_ori_3.jpg',
        trigger: function (res) {
            alert('分享前调用的方法');
        },
        success: function (res) {
            alert('分享成功');
        },
        cancel: function (res) {
            alert('取消分享');
        },
        fail: function (res) {
            alert(JSON.stringify(res));
        }
    });
//分享到QQ里面
    wx.onMenuShareQQ({
        title: '我是分享到qq里面的', // 分享标题
        desc: '这个是测试', // 分享描述
        link: 'http://service.genyuanlian.com/seven_nightmain/index', // 分享链接
        imgUrl: 'http://demo.open.weixin.qq.com/jssdk/images/p2166127561.jpg', // 分享图标
        trigger: function (res) {
            alert('分享前调用的方法');
        },
        complete: function (res) {
            alert(JSON.stringify(res));
        },
        success: function (res) {
            alert('分享成功');
        },
        cancel: function (res) {
            alert('取消分享');
        },
        fail: function (res) {
            alert(JSON.stringify(res));
        }
    });


});
wx.error(function (res) {
    // config信息验证失败会执行error函数，如签名过期导致验证失败，具体错误信息可以打开config的debug模式查看，也可以在返回的res参数中查看，对于SPA可以在这里更新签名。
    alert("chuowu");
});

var imagelocalIds
$('#chooseImage').click(function () {
//选图片事件
    wx.chooseImage({
        count: 1, // 默认9
        sizeType: ['original', 'compressed'], // 可以指定是原图还是压缩图，默认二者都有
        sourceType: ['album', 'camera'], // 可以指定来源是相册还是相机，默认二者都有
        success: function (res) {
            var localIds = res.localIds; // 返回选定照片的本地ID列表，localId可以作为img标签的src属性显示图片
            imagelocalIds = localIds
//alert(`<img src=url(${localIds}) />`)
            var img = document.createElement('img');
            img.src = localIds;
            $('#img').append(img);
            $("#updateimg").attr("src", localIds)
            //选完就上传
            wx.uploadImage({
                localId: imagelocalIds.toString(), // 需要上传的图片的本地ID，由chooseImage接口获得
                isShowProgressTips: 1, // 默认为1，显示进度提示
                success: function (res) {
                    var serverId = res.serverId; // 返回图片的服务器端ID
                    var media_id = {"media_id": res.serverId};
                    $("#medId").attr("value", serverId)
                    alert($("medId").val())
                }
            });
        }
    });
});


//上传图片
$('#upwixinimage').click(function () {
    alert(imagelocalIds)
    wx.uploadImage({
        localId: imagelocalIds.toString(), // 需要上传的图片的本地ID，由chooseImage接口获得
        isShowProgressTips: 1, // 默认为1，显示进度提示
        success: function (res) {
            var serverId = res.serverId; // 返回图片的服务器端ID
            var media_id = {"media_id": res.serverId};
            $("#medId").attr("value", serverId)
            alert($("medId").val())
        }
    });
})


//这个是事先定义好的变量。会自动把声音放进去
var voice = {
    localId: '',
    serverId: ''
};

//开始录音
$('#startRecord').click(function () {

    wx.startRecord({
        cancel: function () {
            alert('开始录音');
        }
    });
});

//停止录音
$('#stopRecord').click(function () {


    wx.stopRecord({
        success: function (res) {
            voice.localId = res.localId;
        },
        fail: function (res) {
            alert(JSON.stringify(res));
        }
    });
});

// 这个是必须需要的录音和播放录音
// wx.onVoiceRecordEnd({
//  complete: function (res) {
//   voice.localId = res.localId;
//   alert('声音赋值');
// }
// });

//播放刚刚录音的声音
$('#playVoice').click(function () {
    if (voice.localId == '') {
        alert('播放刚刚录音的声音');
        return;
    }
    wx.playVoice({
        localId: voice.localId
    });
});

//分享给朋友”按钮点击状态及自定义分享内容接口,可是内置的接口好像用按钮不管用，同理分享的朋友圈什么
$('#onMenuShareAppMessage').click(function () {
    wx.onMenuShareAppMessage({
        title: '这个是测试',
        desc: '哎，描述啥啊',
        link: 'http://neimeng.datahunter.cn/',
        imgUrl: 'http://demo.open.weixin.qq.com/jssdk/images/p2166127561.jpg',
        trigger: function (res) {

            alert('什么鬼？');
        },
        success: function (res) {
            alert('分享成功');
        },
        cancel: function (res) {
            alert('取消分享');
        },
        fail: function (res) {
            alert(JSON.stringify(res));
        }
    });
});


//获取地理位置
$('#getLocation').click(function () {
    wx.getLocation({
        success: function (res) {
            alert(JSON.stringify(res));
        },
        cancel: function (res) {
            alert('获取你的地理位置是否同意');
        }
    });
});
$('#chooseWXPay').click(function () {
    //微信支付
    wx.chooseWXPay({
        timestamp: 0, // 支付签名时间戳，注意微信jssdk中的所有使用timestamp字段均为小写。但最新版的支付后台生成签名使用的timeStamp字段名需大写其中的S字符
        nonceStr: '', // 支付签名随机串，不长于 32 位
        package: '', // 统一支付接口返回的prepay_id参数值，提交格式如：prepay_id=\*\*\*）
        signType: '', // 签名方式，默认为'SHA1'，使用新版支付需传入'MD5'
        paySign: '', // 支付签名
        success: function (res) {
// 支付成功后的回调函数
        },
        cancel: function (res) {
            alert('取消付款');
        },
        fail: function (res) {
            alert(JSON.stringify(res));
        }
    });
});




