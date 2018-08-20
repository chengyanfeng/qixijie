var payfortime=1
var ifprize=""
var g1=1000
var g2=0
var g3=0
var g4=0
var g5=0
var g6=0
var lunpan =function(){
    //幸运大转盘抽奖
    //获得指针元素
    var zhizhen=document.getElementById("zhizhen");
    //存放间隔动画id，用来清除运动
    var dbox=null;
    //间隔动画所用时间，表示转动快慢
    var dtime=1;
    //角度，和css设置对应,初始为0
    var deg=null;
    //初始角度
    var sdeg=0;
    //由js设置默认角度
    zhizhen.style.transform="rotate(" + sdeg + "deg)";
    //变化增量
    var cc=5;
    //旋转基本圈数
    var quan=3;
    //多余角度
    var odeg=null;
    //停止时的角度
    var stopdeg=null;
    //区间奖项
    var jiang=[
        [1,51,"9.9BSTK"], //未中奖
        [52,102,"99BSTK"],//6等奖
        [103,153,"蜂蜜(38元)"],//5等奖
        [154,203,"小红薯一筐(48元)"],//4等奖
        [204,251,"算力包(50元)"],//3等奖
        [252,307,"石斛粉(119元)"],//2等奖
        [307,360,"MAC口红套装(1588元)"]//1等奖
    ];
    //可用次数
    cishu=payfortime;
    //奖项判定函数
    function is(deg){
        var res="未中奖";
        for(var i=0;i<jiang.length;i++){
            if(deg>=jiang[i][0] && deg<=jiang[i][1]){
                res=jiang[i][2];
            };
        };
        return res;
    };
    //是否在动画中
    var able=false;
    //概率
    var gailv=[[g1,"9.9BSTK"],[g2,"99BSTK"],[g3,"蜂蜜(38元)"],[g4,"小红薯一筐(48元)"],[g5,"算力包(50元)"],[g6,"石斛粉(119元)"],[0,"MAC口红套装(1588元)"]];
    //开始到结束总时间
    var xq=0;
    //通过奖项设置额外角度的表现
    function set(real){
        var mindeg,maxdeg;
        for(var i=0;i<jiang.length;i++){
            if(real==jiang[i][2]){
                mindeg=jiang[i][0];
                maxdeg=jiang[i][1];
            };
        };
        return mindeg+Math.floor(Math.random()*(maxdeg-mindeg+1));//生成min-max的随机数

    };
    //监听点击事件
    zhizhen.onclick=function(){
        if(!able){
            if(cishu==0){//可用次数处理
                alert("次数耗光，等待下次机会！");


            }else{
                cishu-=1;//次数减少
                deg=sdeg;
                cc=5;
                xq=0;
                var allarr=[];//长度1000，存放0-6 表示奖项
                for(var i=0;i<gailv.length;i++){
                    for(var j=0;j<gailv[i][0];j++){
                        allarr.push(gailv[i][1]);
                    };
                };
                var real=allarr[Math.floor(Math.random()*1000)];
                odeg=set(real);
                stopdeg=quan*360+odeg;
                alltime=stopdeg/cc*dtime;
                dbox=setInterval(dong,dtime);
            };
        };

    };
    //大转盘转动函数
    function dong(){
        able=true;
        deg+=cc;
        if(deg>stopdeg){
            clearInterval(dbox);
            setTimeout(function(){
                able=false;
                if (is(odeg)=="未中奖"){
                    debugger
                    ifprize=is(odeg)
                    //显示未中奖div
                    $("#regret").css("display","block")
                }else {
			ifprize=is(odeg)
			var dialog =weui.dialog({
    title: '中奖了',
    content: '恭喜您获得了'+ifprize,
	className: 'custom-classname',
    buttons: [{
        label: '确定',
        type: 'primary',
        onClick: function () { $("#numberAndName").css("display","block") }
    }]
});

                    //显示未中奖div
                }


            },500)
        }else{
            if(deg%360==0){//判断第几圈
		console.log(deg)
                xq+=1;
		console.log(xq)
                if(xq==quan-1){//到最后一圈
                    cc=1;//增量变为1 变慢旋转
                };
                zhizhen.style.transform="rotate(" + deg + "deg)";
            }else{
                zhizhen.style.transform="rotate(" + deg + "deg)";
            };
        };
    }

};

