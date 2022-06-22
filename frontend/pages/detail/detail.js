// pages/detail/detail.js
import request from "../../utils/request";
import {
  checkIsCorrectNumber
} from "../../utils/util";
const {
  getTemplate,
  templateStyle
} = require('../../utils/template');

var app = getApp();

// 求助状态映射
const getStatus = (statusId) => {
  const sMap = {
    10: "Created", // 已创建
    15: "Canceled", // 已关闭
    20: "Finished" // 已完成
  };

  return sMap[statusId] || "Created";
}

Page({
  /**
   * 页面的初始数据
   */
  data: {
    isLogin: app.globalData.isLogin, // 登录态
    status: "", // 求助状态
    isMyAid: false, // 我的求助
    isMyHelp: false, // 我帮助过的
    aidInfo: {}, // 求助信息
    aidMessages: [], // 留言
    formData: {}, // 我要帮助表单数据
    showHelpForm: false, // 显示我要帮助弹窗
    showSelectForm: false, // 显示完成帮助弹窗
    rules: [{
      name: 'phone',
      rules: {
        validator: function (rule, value) {
          if (
            !checkIsCorrectNumber(value)
          ) {
            return "请填写正确的手机号或座机号";
          }
        }
      },
    }, {
      name: 'info',
      rules: [{
        required: true,
        message: '请输入帮助说明'
      }],
    }],
    selectDialogBtns: [{
      type: 'primary',
      text: '提交',
      value: 0
    }],
    selectedHelpKey: "$$by_self",
    error: "", // 全局错误信息
    isPhoneError: false,
    latitude: 0,
    longitude: 0
  },

  onFormControlChange: function (e) {
    const {
      field
    } = e.currentTarget.dataset;

    this.setData({
      [`formData.${field}`]: e.detail.value,
      error: "",
      isPhoneError: false
    })
  },

  // 提交帮助
  submitForm: function () {
    this.selectComponent('#form').validate((valid, errors) => {
      if (!valid) {
        const firstError = Object.keys(errors)
        if (firstError.length) {
          this.setData({
            error: errors[firstError[0]].message,
            isPhoneError: errors[firstError[0]].name === "phone"
          })
        }
        return false;
      }
      request({
        url: '/api/user/aid/message',
        method: "post",
        data: {
          id: this.data.id,
          phone: this.data.formData.phone,
          content: this.data.formData.info
        }
      }).then((res) => {
        getApp().wxSubscribeAcceptHelp();
        console.log('res.data', res.data);
        this.setData({
          showHelpForm: false,
        });
        wx.hideLoading();
        wx.showToast({
          title: '提交成功'
        });
        // 告知列表页重新刷新
        wx.setStorageSync('needReloadIndex', true);
        wx.setStorageSync('needReloadMyHelp', true);
        // 刷新
        this.getUserAid(this.data.id, this.data.latitude, this.data.longitude);
      }, (error) => {
        console.log('error', error);
        this.setData({
          error
        });
        wx.hideLoading();
      });
    })
  },

  handleHelp: function () {
    this.setData({
      showHelpForm: true
    })
  },

  cancelHelp: function () {
    wx.showModal({
      content: '您确定要关闭本条求助吗？',
      success: (res) => {
        if (res.confirm) {
          request({
            url: `/api/user/aid/${this.data.id}`,
            method: "DELETE",
          }).then(() => {
            wx.showToast({
              title: '取消成功',
              success: () => {
                setTimeout(() => {
                  wx.reLaunch({
                    url: '/pages/askHelp/askHelp'
                  });
                }, 1500)
              }
            });
          }, (error) => {
            console.log('error', error);
            wx.showToast({
              title: '取消失败'
            })
          });
        }
      }
    })
  },

  finishlHelp: function () {
    this.setData({
      showSelectForm: true
    });
  },

  selectHelp: function () {
    const {
      id,
      selectedHelpKey: messageId
    } = this.data;
    request({
      url: `/api/user/aid/${this.data.id}`,
      method: "PUT",
      data: {
        "id": id,
        "messageId": messageId === "$$by_self" ? 0 : messageId
      },
    }).then(() => {
      this.setData({
        showSelectForm: false
      });
      wx.showToast({
        title: '设置成功',
        success: () => {
          setTimeout(() => {
            wx.reLaunch({
              url: '/pages/askHelp/askHelp'
            });
          }, 1500)
        }
      });
    }, (error) => {
      console.log('error', error);
      wx.showToast({
        title: '设置失败',
        icon: "error"
      })
    });
  },

  helpRadioChange: function (e) {
    const selected = e.currentTarget.dataset.key;
    this.setData({
      selectedHelpKey: selected
    })
  },

  callPhone: function (e) {
    const phone = e.currentTarget.dataset.phone;
    wx.makePhoneCall({
      phoneNumber: phone
    })
  },

  // 获取手机号
  getPhoneNumber(e) {
    const phoneCode = e.detail.code;
    if (!phoneCode) {
      return;
    }

    request({
      url: `/api/wx_phone_number?phoneCode=${phoneCode}`,
      method: "GET"
    }).then(({
      data
    }) => {
      this.setData({
        'formData.phone': data.phone,
        error: "",
        isPhoneError: false
      })
    }, (err) => {
      wx.showToast({
        title: err,
        icon: 'none',
        mask: true
      })
    })
  },

  // 求助详情
  getUserAid: function (id, latitude, longitude) {
    this.setData({
      latitude,
      longitude
    });
    wx.showLoading({
      title: "加载中"
    });
    request({
      url: `/api/user/aid/${id}?latitude=${latitude}&longitude=${longitude}`,
    }).then((res) => {
      console.log('res.data', res.data);
      const aidMessages = ((res.data.aid && res.data.aid.message) || []).sort((a, b) => {
        return b.createTime - a.createTime
      })
      this.setData({
        aidInfo: res.data.aid,
        aidMessages: aidMessages,
        status: getStatus(res.data.aid && res.data.aid.status),
        isMyAid: res.data.isMyAid,
        isMyHelp: res.data.isMyHelp
      });
      wx.hideLoading();
    }, (error) => {
      console.log('error', error);
      this.setData({
        error
      });
      wx.hideLoading();
    });
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    const {
      id,
    } = options;
    // id = "95073702466551810"

    // 不带id参数非法，重定向首页
    if (!id) {
      return wx.reLaunch({
        url: '/pages/index/index'
      });
    }

    this.setData({
      id: id
    });

    const appInst = getApp();
    // 登录
    appInst.login({
      success: () => {
        this.setData({
          isLogin: app.globalData.isLogin
        });
        // 获取位置信息
        appInst.getLocation({
          success: (latitude, longitude) => {
            this.getUserAid(id, latitude, longitude);
          },
          fail: (latitude, longitude) => {
            // 未授权位置 用默认坐标
            this.getUserAid(id, latitude, longitude);
          }
        });
      },
      fail: () => {
        // 未授权登录 重定向首页
        wx.reLaunch({
          url: '/pages/index/index'
        });
      },
    });
  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady: function () {
    console.log("ready");
  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow: function () {
    console.log("show");
  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide: function () {},

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload: function () {},

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh: function () {},

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom: function () {},

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage: function () {
    const {
      id
    } = this.data;
    const path = `/pages/detail/detail?id=${id}`;

    const promise = new Promise(resolve => {
      this.widget = this.selectComponent('.share-canvas');
      const p1 = this.widget.renderToCanvas({
        wxml: getTemplate(this.data),
        style: templateStyle
      });
      p1.then(() => {
        const p2 = this.widget.canvasToTempFilePath()
        p2.then(res => {
          resolve({
            title: "TA需要您的帮助",
            imageUrl: res.tempFilePath
          })
        })
      }).catch(() => {
        // hack：真机上会走到error里
        const p2 = this.widget.canvasToTempFilePath()
        p2.then(res => {
          resolve({
            title: "TA需要您的帮助",
            imageUrl: res.tempFilePath
          })
        })
      });
    });
    return {
      title: "TA需要您的帮助",
      path: path,
      promise
    };
  },
});