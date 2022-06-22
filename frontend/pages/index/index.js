import request from "../../utils/request";
var app = getApp();

Page({
    data: {
        pageNumber: 0,
        pageSize: 10,
        isAll: false,
        latitude: "",
        longitude: "",
        listArr: [],
        pageIsLogin: false, // 当前页面保存用户登录态
    },

    fetchList(latitude, longitude) {
        const { listArr, pageNumber, pageSize, isAll } = this.data;

        // 已加载全部数据
        if (isAll) {
            this.setData({ isAll: true });
            return;
        }

        request({
            url: "/api/discovery",
            method: "GET",
            data: {
                pageNumber,
                pageSize,
                latitude: latitude || app.globalData.latitude,
                longitude: longitude || app.globalData.longitude,
            },
        }).then(
            (res) => {
                console.log("首页-返回值：", res);
                const { list } = res.data;
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
        wx.setStorageSync('needReloadIndex', false);
        this.setData({ pageNumber: 0, listArr: [], isAll: false });
    },

    // 重置加载状态
    resetStatus() {
        wx.stopPullDownRefresh();
    },

    /**
     *  只在初始化触发执行一次
     */
    onLoad: function () {},

    /**
     * 每次打开页面都触发
     */
    onShow: function () {
        const { pageIsLogin } = this.data;
        const { isLogin } = app.globalData;

        // 用户的登录状态发生变化时，清空当前列表数据，重新拉取数据！
        if (pageIsLogin !== isLogin || (!pageIsLogin && !isLogin) || wx.getStorageSync('needReloadIndex')) {
            // 初始化数据集
            this.resetData();

            // 授权获取用户地理位置信息
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

            // 更新当前页的登录态
            this.setData({ pageIsLogin: app.globalData.isLogin });
        }
    },

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

    // 跳转详情页
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
            path: "/pages/index/index",
            // imageUrl: "/images/logo_ucloud.png",
        };
    },
});
