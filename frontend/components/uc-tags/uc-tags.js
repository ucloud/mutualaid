const { getTagText } = require("../../utils/constant");
const colorMap = new Map().set(1, "level-one").set(2, "level-two").set(3, "level-three").set(4, "level-four");

Component({
    properties: {
        tags: {
            type: Object,
            value: {},
            observer: function (tags) {
                this.setData({
                    tagList: {
                        ...tags,
                        typeText: getTagText("type", tags.type),
                        groupText: getTagText("group", tags.group),
                        emergencyText: getTagText("emergency", tags.emergency),
                        colorLevel: colorMap.get(tags.emergency),
                    },
                });
            },
        },
        // 禁用（已解决/已完成）
        disabled: {
            type: Boolean,
            value: false,
        },
    },

    data: {
        tagList: {},
    },

    methods: {},

    // 组件所在页面的生命周期函数
    pageLifetimes: {
        show() {},
    },

    lifetimes: {
        attached: function () {
            const { tags } = this.data;
            this.setData({
                tagList: {
                    ...tags,
                    typeText: getTagText("type", tags.type),
                    groupText: getTagText("group", tags.group),
                    emergencyText: getTagText("emergency", tags.emergency),
                    colorLevel: colorMap.get(tags.emergency),
                },
            });
        },
    },
});
