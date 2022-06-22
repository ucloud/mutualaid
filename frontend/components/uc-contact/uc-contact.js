Component({
  properties: {
    phone: {
      type: String,
      value: ""
    },
    address: {
      type: String,
      value: ""
    },
    hide: {
      type: Boolean,
      value: false
    },
  },
  data: {},
  methods: {
    callPhone: function (e) {
      const phone = e.currentTarget.dataset.phone;
      wx.makePhoneCall({
        phoneNumber: phone
      })
    }
  }
})