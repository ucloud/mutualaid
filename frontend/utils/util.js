function serialize(obj) {
    var str = [];
    for (var p in obj)
      if (obj.hasOwnProperty(p)) {
        str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
      }
    return str.join("&");
}

function formatTime(date) {
    var year = date.getFullYear();
    var month = date.getMonth() + 1;
    var day = date.getDate();

    var hour = date.getHours();
    var minute = date.getMinutes();
    var second = date.getSeconds();

    return [year, month, day].map(formatNumber).join("/") + " " + [hour, minute, second].map(formatNumber).join(":");
}

function formatNumber(n) {
    n = n.toString();
    return n[1] ? n : "0" + n;
}

function formatMobile(n) {
    if (n) {
        var mphone = n.substr(0, 3) + "****" + n.substr(7);
        return mphone;
    } else {
        return "";
    }
}

/**
 * 设置微信自带tabber的角标
 * 1、在非tabbar的页面上使用这部分代码会报错，可以在每个tabbar页面的onshow中设定，达到实时显示角标数据。
 * 2、部分机型真机上角标被遮挡问题（暂未找到解决办法，可能需要重写tabbar）：https://developers.weixin.qq.com/community/develop/doc/0006ae0b848f208d04cda10945bc00?highLine=wx.setTabBarBadge
 * @param {*} mySendHelpNum 我帮助的
 * @param {*} myAskHelpNum 我的求助
 */
function setWxTabBarBadge(mySendHelpNum, myAskHelpNum) {
    if (mySendHelpNum) {
        wx.setTabBarBadge({ index: 1, text: mySendHelpNum });
    } else {
        wx.removeTabBarBadge({ index: 1 });
    }

    if (myAskHelpNum) {
        wx.setTabBarBadge({ index: 3, text: myAskHelpNum });
    } else {
        wx.removeTabBarBadge({ index: 3 });
    }
}

function checkIsCorrectNumber(value){
    if (
        !/^0?(13|14|15|17|18|19)[0-9]{9}$/.test(value) &&
        !/^[0-9-()（）]{7,18}$/.test(value)
    ) {
        return false;
    }

    return true;
}

module.exports = {
    formatTime,
    formatNumber,
    formatMobile,
    setWxTabBarBadge,
    serialize,
    checkIsCorrectNumber
};
