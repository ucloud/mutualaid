import request from "../../utils/request";
import { getTagText } from "../../utils/constant";
var app = getApp();

Page({
    data: {
        pageNumber: 0,
        pageSize: 10,
        isAll: false,
        latitude: "",
        longitude: "",
        listArr: [],
    },

    fetchList(latitude, longitude) {
        const { listArr, pageNumber, pageSize, isAll } = this.data;

        // 已加载全部数据
        if (isAll) {
            this.setData({ isAll: true });
            return;
        }

        request({
            url: "/api/user/aid/offered",
            method: "GET",
            data: {
                pageNumber,
                pageSize,
                latitude: latitude || app.globalData.latitude,
                longitude: longitude || app.globalData.longitude,
            },
        }).then(
            (res) => {
                console.log("我的帮助-返回值：", res);
                const { list } = res.data;
                list.forEach((item) => {
                    item["isDisableStatus"] = [15, 20].includes(item.status);
                    item["statusText"] = getTagText("status", item.status);
                });
                this.setData({
                    listArr: [...listArr, ...list],
                    isAll: list.length < pageSize,
                });
                this.resetStatus();
            },
            (err) => {
                this.resetStatus();
                wx.showToast({ title: err, icon: "none", mask: true });
            }
        );
    },

    // 重置数据
    resetData() {
        this.setData({ pageNumber: 0, listArr: [], isAll: false });
    },

    // 重置加载状态
    resetStatus() {
        wx.stopPullDownRefresh();
    },

    /**
     * 只在初始化触发执行一次
     */
    onLoad: function () {
        this.loadList();
    },

    /**
     * 每次打开页面都触发
     */
    onShow: function () {
        if (wx.getStorageSync("needReloadMyHelp")) {
            this.loadList();
            wx.setStorageSync("needReloadMyHelp", false);
        }
    },

    loadList() {
        // 初始化数据集
        this.resetData();

        // 1、拦截授权登录
        app.login({
            success: () => {
                // 2、授权获取用户地理位置信息
                app.getLocation({
                    success: (latitude, longitude) => {
                        this.setData({ latitude, longitude });
                        this.fetchList(latitude, longitude);
                    },
                    fail: (SH_latitude, SH_longitude) => {
                        this.setData({
                            latitude: SH_latitude,
                            longitude: SH_longitude,
                        });
                        this.fetchList(SH_latitude, SH_longitude);
                    },
                });
            },
            fail: () => {
                wx.reLaunch({ url: "/pages/index/index" });
            },
        });
    },

    /**
     * 当前tab被点击时触发
     * 该方法可解决tabbar再次点击切换不触发生命周期钩子
     */
    onTabItemTap(item) {},

    /**
     * 触底加载更多
     */
    onReachBottom: function () {
        this.setData({ pageNumber: this.data.pageNumber + 1 });
        const { latitude, longitude } = this.data;
        this.fetchList(latitude, longitude);
    },

    /**
     * 下拉刷新
     */
    onPullDownRefresh: function () {
        const { latitude, longitude } = this.data;
        this.resetData();
        this.fetchList(latitude, longitude);
    },

    toDetail: function (e) {
        const id = e.currentTarget.dataset.id;
        wx.navigateTo({
            url: `/pages/detail/detail?id=${id}`,
        });
    },

    /**
     * 用户点击右上角分享
     */
    onShareAppMessage: function () {
        return {
            title: "TA们需要您的帮助",
            path: "/pages/sendHelp/sendHelp",
            // imageUrl: "/images/logo_ucloud.png",
        };
    },
});
