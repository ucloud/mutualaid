const transform = (value) => {
  const _val = parseFloat(value);
  if (isNaN(_val)) {
    return "";
  }
  const result = Math.round(_val / 1000);
  return result > 1 ? `${result}km` : `≤ 1km`;
}

Component({
  properties: {
    value: {
      type: String,
      value: "",
      observer: function (val) {
        this.setData({
          valueForShow: transform(val)
        })
      }
    },
    address: {
      type: String,
      value: "",
    },
    disabled: {
      type: Boolean,
      value: false,
    },
  },

  data: {
    valueForShow: ""
  },

  lifetimes: {
    attached: function () {
      this.setData({
        valueForShow: transform(this.data.value)
      })
    },
  },

  methods: {}
})