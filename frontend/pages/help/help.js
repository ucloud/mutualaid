// pages/help.js
import request from "../../utils/request";
import {
    AidType,
    AidGroup,
    EmergencyLevel
} from "../../utils/constant";

import { checkIsCorrectNumber } from "../../utils/util"

Page({
    /**
     * 页面的初始数据
     */
    data: {
        aidType: AidType,
        aidGroup: AidGroup,
        emergencyLevel: EmergencyLevel,
        aidTypeValue: -1,
        aidGroupValue: -1,
        emergencyLevelValue: -1,
        content: "",
        btnDisabled: true,
        latitude: "",
        longitude: "",
        aidAddress: getApp().globalData.isSH_latitude_longitude ? '上海市' : '',
        aidNumber: getApp().globalData.userInfo.phone,
    },
    checkSubmitBtn: function () {
        this.setData({
            btnDisabled: !(this.data.aidTypeValue >= 0 && this.data.aidGroupValue >= 0 && this.data.emergencyLevelValue >= 0 && !!this.data.content && !!this.data.aidAddress && !!this.data.aidNumber)
        })
    },
    bindInputChange: function (event) {
        this.setData({
            content: event.detail.value
        }, () => {
            this.checkSubmitBtn()
        })
    },
    bindAidTypeChange: function (event) {
        this.setData({
            aidTypeValue: event.detail.value
        }, () => {
            this.checkSubmitBtn();
        })
    },
    bindAidGroupChange: function (event) {
        this.setData({
            aidGroupValue: event.detail.value
        }, () => {
            this.checkSubmitBtn();
        })
    },
    bindEmergencyLevelChange: function (event) {
        this.setData({
            emergencyLevelValue: event.detail.value
        }, () => {
            this.checkSubmitBtn();
        })
    },
    submitForm() {
        const {
            latitude,
            longitude,
            content,
            aidType,
            aidTypeValue,
            aidGroup,
            aidGroupValue,
            emergencyLevel,
            emergencyLevelValue,
            aidAddress,
            aidNumber
        } = this.data;

        if (!checkIsCorrectNumber(aidNumber)) {
            wx.showToast({
                title: '联系电话请填写正确的手机号或座机号',
                icon: 'none',
                duration: 2000
            })

            return;
        }

        request({
            url: "/api/user/aid",
            method: "POST",
            data: {
                latitude,
                longitude,
                "type": aidType[aidTypeValue].value,
                "group": aidGroup[aidGroupValue].value,
                "emergency": emergencyLevel[emergencyLevelValue].value,
                content,
                phone: aidNumber,
                addr: aidAddress
            }
        }).then(() => {
            getApp().wxSubscribeGotHelp();

            wx.showToast({
                title: '提交成功',
                mask: true,
                success: () => {
                    setTimeout(() => {
                        wx.reLaunch({
                            url: '/pages/askHelp/askHelp'
                        });
                    }, 1500)
                }
            })
        }, (err) => {
            wx.showToast({
                title: err,
                icon: 'none',
                mask: true
            })
        })
    },

    onAddressFocus() {
        if(getApp().globalData.isSH_latitude_longitude){
            return;
        }

        wx.chooseLocation({
            success: (res) => {
                const { address = "", name = "" } = res;
                this.setData({
                    aidAddress: `${address}${name}`,
                    latitude: res.latitude,
                    longitude: res.longitude,
                }, () => {
                    this.checkSubmitBtn();
                })
            }
        });
    },

    onNumberChange(event) {
        this.setData({
            aidNumber: event.detail.value
        }, () => {
            this.checkSubmitBtn();
        })
    },

    getPhoneNumber(e) {
        const phoneCode = e.detail.code;
        if(!phoneCode){
            return;
        }

        request({
            url: `/api/wx_phone_number?phoneCode=${phoneCode}`,
            method: "GET"
        }).then(({ data }) => {
            this.setData({
                aidNumber: data.phone
            }, () => {
                this.checkSubmitBtn();
            })
        }, (err) => {
            wx.showToast({
                title: err,
                icon: 'none',
                mask: true
            })
        })
    },

    /**
     * 生命周期函数--监听页面加载
     */
    onLoad: function (options) {
    },

    /**
     * 生命周期函数--监听页面初次渲染完成
     */
    onReady: function () {
        wx.setNavigationBarTitle({
            title: '发布求助信息'
        });

        const app = getApp();

        app.getLocation({
            success: (latitude, longitude, isSH_latitude_longitude)=>{
                if(isSH_latitude_longitude){
                    this.data.latitude = latitude;
                    this.data.longitude = longitude;
                }
            },
            fail: (latitude, longitude, isSH_latitude_longitude)=>{
                if(isSH_latitude_longitude){
                    this.data.latitude = latitude;
                    this.data.longitude = longitude;
                }
            }
        });
        app.login({
            fail: ()=>{
                wx.reLaunch({
                    url: '/pages/index/index'
                });
            },
        });
    },

    /**
     * 生命周期函数--监听页面显示
     */
    onShow: function () {

    },

    /**
     * 生命周期函数--监听页面隐藏
     */
    onHide: function () {

    },

    /**
     * 生命周期函数--监听页面卸载
     */
    onUnload: function () {

    },

    /**
     * 页面相关事件处理函数--监听用户下拉动作
     */
    onPullDownRefresh: function () {

    },

    /**
     * 页面上拉触底事件的处理函数
     */
    onReachBottom: function () {

    },

    /**
     * 用户点击右上角分享
     */
    onShareAppMessage: function () {

    }
})