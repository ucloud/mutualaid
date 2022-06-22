const moment = require("../../lib/moment/we-moment");

const transform = (val) => {
  const value = parseInt(val);
  if (!isNaN(value)) {
    if (moment(value * 1000).isBefore(moment().subtract(7, 'days'))) {
      return moment(value * 1000).format("YYYY-MM-DD HH:mm:ss");
    }
    return moment(value * 1000).fromNow();
  }
  return "";
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
    }
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
})