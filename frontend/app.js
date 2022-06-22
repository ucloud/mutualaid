//app.js
import request from "./utils/request";
import {
    serialize
} from './utils/util';

const useLocationDesc = "地理位置用于就近互助，不授权则使用上海中心位置代替。";
const SH_latitude = 31.2381;
const SH_longitude = 121.4692;
const initGlobalData = {
    isSH_latitude_longitude: false,
    onLaunchPathAndQuery: "",
    isLogin: false,
    jwtToken: "",
    code: "",
    userProfile: null,
    userInfo: {
        phone: ""
    },
    latitude: "",
    longitude: ""
};

App({
    onLaunch: function (e) {
        const {
            path,
            query
        } = e;
        this.globalData.onLaunchPathAndQuery = `${path}?${serialize(query)}`;

        // this.getLocation();
        // this.login();

        const updateManager = wx.getUpdateManager();
        updateManager.onUpdateReady(function () {
            wx.showModal({
                title: '更新提示',
                content: '新版本已经准备好，是否重启应用？',
                success: function (res) {
                    if (res.confirm) {
                        // 新的版本已经下载好，调用 applyUpdate 应用新版本并重启
                        updateManager.applyUpdate()
                    }
                }
            })
        })
    },

    getLocation({
        success = () => {},
        fail = () => {}
    } = {}) {
        const {
            latitude,
            longitude,
            isSH_latitude_longitude,
        } = this.globalData;

        if (latitude && longitude) {
            success(latitude, longitude, isSH_latitude_longitude);
            return;
        }

        // 获取地理位置权限
        wx.getSetting({
            success: (res) => {
                if (!res.authSetting['scope.userLocation']) {
                    // 去弹窗开启授权地理位置
                    wx.authorize({
                        scope: 'scope.userLocation',
                        success: () => {
                            this.wxGetLocation(success);
                        },
                        fail: () => {
                            this.globalData.latitude = SH_latitude;
                            this.globalData.longitude = SH_longitude;
                            this.globalData.isSH_latitude_longitude = true;
                            fail(SH_latitude, SH_longitude, true);

                            // this.showLocationSettingModal(() => {
                            //     this.wxGetLocation(success);
                            // }, () => {
                            //     this.globalData.latitude = SH_latitude;
                            //     this.globalData.longitude = SH_longitude;
                            //     this.globalData.isSH_latitude_longitude = true;
                            //     fail(SH_latitude, SH_longitude, true);
                            // });
                        }
                    })
                } else {
                    this.wxGetLocation(success);
                }
            }
        });
    },

    logout: function () {
        this.globalData = {
            ...initGlobalData
        };
        wx.reLaunch({
            url: '/pages/index/index'
        });
    },

    login: function (options = {}) {
        const {
            success = () => {}, fail = () => {}
        } = options;
        if (this.globalData.isLogin) {
            success();
            return;
        }

        wx.showLoading({
            title: '加载中'
        });
        // 微信登录
        wx.login({
            success: (res) => {
                if (res.code) {
                    this.setGlobalData({
                        code: res.code
                    })
                    this.setUserProfile(success, fail);
                } else {
                    wx.showToast({
                        title: res.errMsg,
                        icon: 'none',
                        duration: 2000
                    })
                    fail(res.errMsg);
                }
            },
            fail: (res) => {
                wx.hideLoading();
                fail(res && res.errMsg);
            }
        })
    },

    setUserProfile: function (success, fail) {
        wx.getStorage({
            key: "userProfile",
            success: (res) => {
                console.log('getStorage userProfile success', res.data);
                this.setGlobalData({
                    userProfile: res.data
                })
            },
            fail: (err) => {
                console.log('getStorage userProfile fail', err);
            },
            complete: () => {
                wx.hideLoading();
                if (!this.globalData.userProfile) {
                    // 获取昵称、头像授权
                    wx.showModal({
                        title: '请求微信授权',
                        content: 'UCloud优刻得云计算携手真爱梦想基金，为更高效的地帮助您，需要获取微信授权信息。（昵称、头像）',
                        confirmText: '立即授权',
                        cancelText: '暂不授权',
                        success: (res) => {
                            if (res.confirm) {
                                wx.getUserProfile({
                                    desc: '用于用户的昵称、头像展示。',
                                    success: (res) => {
                                        wx.setStorage({
                                            key: 'userProfile',
                                            data: res.userInfo
                                        });
                                        this.setGlobalData({
                                            userProfile: res.userInfo
                                        })
                                        this.activeUser(() => {
                                            success();
                                        });
                                    },
                                    fail: (res) => {
                                        fail(res && res.errMsg);
                                    }
                                })
                            } else if (res.cancel) {
                                // wx.exitMiniProgram();
                                fail();
                            }
                        }
                    })
                } else {
                    this.activeUser(() => {
                        success();
                    });
                }
            }
        })
    },

    showLocationSettingModal: function (success, fail) {
        wx.showModal({
            title: '提示',
            content: useLocationDesc,
            cancelText: '取消',
            confirmText: '去设置',
            success: (res) => {
                if(res.confirm){
                    wx.openSetting({
                        success: (res) => {
                            if (res.authSetting['scope.userLocation']) {
                                success();
                            }
                        },
                        fail: () => {
                            this.showLocationSettingModal(success, fail);
                        }
                    })
                } else if(res.cancel){
                    fail();
                }
            }
        })
    },

    getUser: function () {
        request({
            url: "/api/user",
            method: "GET",
        }).then(({
            data
        }) => {
            this.setGlobalData({
                userInfo: data.user
            });
        }, (err) => {
            wx.showToast({
                title: err.toString() || '获取用户信息失败',
                icon: 'none',
                duration: 1500
            })
        })
    },

    activeUser: function (success = () => {}) {
        const {
            code,
            userProfile
        } = this.globalData;
        request({
            url: "/api/activeuser",
            method: "POST",
            data: {
                "loginCode": code,
                "name": userProfile.nickName,
                "icon": userProfile.avatarUrl
            }
        }).then(({
            header
        }) => {
            this.setGlobalData({
                isLogin: true,
                jwtToken: header['Jwt-Token']
            });

            this.getUser();
            success();
        }, (err) => {
            // 登录失败
            wx.showModal({
                title: '登录失败',
                content: err,
                cancelText: '退出',
                confirmText: '重新登录',
                success: (res) => {
                    if (res.confirm) {
                        this.activeUser(success);
                        return;
                    }

                    wx.exitMiniProgram();
                }
            })
        })
    },

    setGlobalData: function (data) {
        this.globalData = {
            ...this.globalData,
            ...data
        }
    },

    getGlobalData: function () {
        return {
            ...this.globalData
        };
    },

    globalData: {
        ...initGlobalData
    },

    wxGetLocation: function (success) {
        wx.getLocation({
            success: (result) => {
                const {
                    latitude,
                    longitude
                } = result;
                this.globalData.latitude = latitude;
                this.globalData.longitude = longitude;
                this.globalData.isSH_latitude_longitude = false;
            },
            fail: () => {
                this.globalData.latitude = SH_latitude;
                this.globalData.longitude = SH_longitude;
                this.globalData.isSH_latitude_longitude = true;
            },
            complete: () => {
                const {
                    latitude,
                    longitude,
                    isSH_latitude_longitude
                } = this.globalData;
                success(latitude, longitude, isSH_latitude_longitude);
            },
        })
    },

    wxSubscribeGotHelp: function() {
        wx.requestSubscribeMessage({
            tmplIds: [
                'tmplIds' // 替换为小程序订阅消息中创建的模板ID
            ],
            success:  (res) =>{
                console.log('wxSubscribeGotHelp success', res);
            },
            fail:  (res) =>{
                console.log('wxSubscribeGotHelp fail', res);
            },
        })
    },

    wxSubscribeAcceptHelp: function() {
        wx.requestSubscribeMessage({
            tmplIds: [
                'tmplIds' // 替换为小程序订阅消息中创建的模板ID
            ],
            success:  (res) =>{
                console.log('wxSubscribeAcceptHelp success', res);
            },
            fail:  (res) =>{
                console.log('wxSubscribeAcceptHelp fail', res);
            },
        })
    }
});