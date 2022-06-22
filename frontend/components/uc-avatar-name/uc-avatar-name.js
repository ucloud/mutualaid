const defaultAvatar = "../../images/avatar_default.png";
const defaultName = "匿名";

Component({
  properties: {
    avatar: {
      type: String,
      value: defaultAvatar
    },
    name: {
      type: String,
      value: defaultName
    }
  },
  data: {
    defaultAvatar: defaultAvatar,
    defaultName: defaultName
  },
  methods: {
    // 这里是一个自定义方法
  }
})