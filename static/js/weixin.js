var ticket;
var imagelocalIds;
var shareUrl
var addr
var medId
//这个是获取用户的信息
$.ajax({
    async: false,
    type: 'post',
    data: {'url': location.href},
    url: '/share/get_user_token',
    success: function (data) {
    }
});
//这个是上传图片到本地服务器
$.ajax({
    async: false,
    type: 'post',
    data: {'': location.href},
    url: '/share/get_user_token',
    success: function (data) {
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
wx.config({
    debug: false, // 开启调试模式,调用的所有api的返回值会在客户端alert出来，若要查看传入的参数，可以在pc端打开，参数信息会通过log打出，仅在pc端时才会打印。
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
        'uploadImage',
        'chooseWXPay'
    ] // 必填，需要使用的JS接口列表，所有JS接口列表见附录2
});
//分享朋友圈和分享给朋友调用的方法
var share=function () {
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
            link: shareUrl, // 分享链接，该链接域名或路径必须与当前页面对应的公众号JS安全域名一致
            imgUrl: 'https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1514438536365&di=09f25c134d4b0589e136a5fddc22a22c&imgtype=0&src=http%3A%2F%2Fwww.haha365.com%2Fuploadfile%2F2014%2F0404%2F20140404063425613.jpg', // 分享图标
            success: function () {
                $.get("/sharesucess?openid="+openid, function (res, status) {

                    })
            },
            cancel: function () {
            }
        });
//分享给朋友调用的方法
        wx.onMenuShareAppMessage({
            title: '链爱的表白',
            desc: '分享一下 不变的爱',
            link: shareUrl,
            imgUrl: 'https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1514438682725&di=7624aa8b34c3d92a2d2b95254aff2450&imgtype=0&src=http%3A%2F%2Fvpic.video.qq.com%2F3170316%2Fb03322rccwa_ori_3.jpg',
            trigger: function (res) {
                },
            success: function (res) {
                $.get("/sharesucess?openid="+openid, function (res, status) {

                })
            },
            cancel: function (res) {
            },
            fail: function (res) {
            }
        });

    });
}



wx.error(function (res) {
    // config信息验证失败会执行error函数，如签名过期导致验证失败，具体错误信息可以打开config的debug模式查看，也可以在返回的res参数中查看，对于SPA可以在这里更新签名。
    alert("chuowu");
});

var chooseimge=function(){
    //选图片事件
    wx.chooseImage({
        count: 1, // 默认9
        sizeType: ['original', 'compressed'], // 可以指定是原图还是压缩图，默认二者都有
        sourceType: ['album', 'camera'], // 可以指定来源是相册还是相机，默认二者都有
        success: function (res) {
            var localIds = res.localIds; // 返回选定照片的本地ID列表，localId可以作为img标签的src属性显示图片
            imagelocalIds = localIds
            $("#updateimg").attr("src", localIds)
            $("#updateimg1").attr("src", localIds)
            //选完就上传
            wx.uploadImage({
                localId: imagelocalIds.toString(), // 需要上传的图片的本地ID，由chooseImage接口获得
                isShowProgressTips: 1, // 默认为1，显示进度提示
                success: function (res) {
                    var serverId = res.serverId; // 返回图片的服务器端ID
                    medId=serverId
                }
            });
        }
    });
}
//上传图片
$('#upwixinimage').click(function () {

    wx.uploadImage({
        localId: imagelocalIds.toString(), // 需要上传的图片的本地ID，由chooseImage接口获得
        isShowProgressTips: 1, // 默认为1，显示进度提示
        success: function (res) {
            var serverId = res.serverId; // 返回图片的服务器端ID
            var media_id = {"media_id": res.serverId};
            $("#medId").attr("value", serverId)
        }
    });
})

//微信支付
var pay=function () {
    debugger
    a=""
    $.ajaxSettings.async = false;
    $.get("/getpayid",function(data,status){
        da=$.parseXML(data)
        $xml=$(da)
        $title=$xml.find("prepay_id")
        a=$title.text()
    });
    debugger
    //微信支付
    var sinNosha1="appid=wx53d52d70ccd6439f&nonceStr=werewr2r3r2r2&package="+a+"&signType=SHA1&timeStamp=1414587457";
    sign=hex_sha1(sinNosha1)
    wx.chooseWXPay({
        timestamp: 1414587457, // 支付签名时间戳，注意微信jssdk中的所有使用timestamp字段均为小写。但最新版的支付后台生成签名使用的timeStamp字段名需大写其中的S字符
        nonceStr: 'werewr2r3r2r2', // 支付签名随机串，不长于 32 位
        package: a, // 统一支付接口返回的prepay_id参数值，提交格式如：prepay_id=\*\*\*）
        signType: 'SHA1', // 签名方式，默认为'SHA1'，使用新版支付需传入'MD5'
        paySign: sign, // 支付签名
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

}







