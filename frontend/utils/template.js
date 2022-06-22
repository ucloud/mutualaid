const moment = require("../lib/moment/we-moment");
const {
  getTagText
} = require("./constant");
const colorMap = new Map().set(1, "level-one").set(2, "level-two").set(3, "level-three").set(4, "level-four");

const transform = (val) => {
  const value = parseInt(val);
  if (!isNaN(value)) {
    return moment(parseInt(val) * 1000).format("YYYY-MM-DD HH:mm:ss")
  }
  return "";
}

const getTags = (tags) => {
  const tagList = {
    ...tags,
    typeText: getTagText("type", tags.type),
    groupText: getTagText("group", tags.group),
    emergencyText: getTagText("emergency", tags.emergency),
    colorLevel: colorMap.get(tags.emergency),
  };

  return `
  <view class="item-tag-box">
      <view class="tag-item tag-type"><text>${tagList.typeText}</text></view>
      <view class="tag-item tag-person">
        <text>${tagList.groupText}</text>
      </view>
      <view class="tag-item tag-level ${tagList.colorLevel}">
        <text>${tagList.emergencyText}</text>
      </view>
  </view>
  `;
};

const transformDistance = (value) => {
  const _val = parseFloat(value);
  if (isNaN(_val)) {
    return "";
  }
  const result = Math.round(_val / 1000);
  return result > 1 ? `${result}km` : `≤ 1km`;
}

const getTemplate = (data) => `
<view class='detail-main'>
  <!-- 求助信息 -->
  <view class="detail-cells detail-cells-info">
    <view class="detail-cells-title">
      <view class="uc-avatar-name">
        <image class="uc-avatar uc-avatar-sm uc-avatar-name-avatar" src="${data.aidInfo.user.icon || '../../images/avatar_default.png'}" />
        <text class="uc-avatar-name-name">${data.aidInfo.user.name || "匿名"}</text>
      </view>
    </view>
    <view class="detail-line"></view>
    <view class="detail-cell detail-cell-content">
      <view class="uc-create-time">
        <text>求助时间：${transform(data.aidInfo.createTime)}</text>
      </view>
      <view class="weui-cell-desc weui-cell-desc-dark"><text>${data.aidInfo.content}</text></view>
      ${getTags(data.aidInfo)}
      <view class="uc-contact">
        <image class="uc-contact-icon" src="../../images/lock.png"></image>
        <view class="uc-contact-detail">
          <view class="uc-contact-detail-main">
            <text>联系方式回复后可见</text>
          </view>
          <view class="uc-contact-detail-desc">
            <text>点击下方“我要帮助”</text>
          </view>
        </view>
      </view>

      <view class="uc-distance">
        <image class="uc-distance-icon" src="../../images/icon_position.png"></image>
        <text class="uc-distance-value">${transformDistance(data.aidInfo.distance)}</text>
        <text class="uc-distance-address">${data.aidInfo.address}</text>
      </view>

      <view style="margin-top: 8rpx;">
        <uc-distance value="{{aidInfo.distance}}" address="{{aidInfo.address}}"></uc-distance>
      </view>
    </view>
  </view>
</view>
`;

const templateStyle = {
  detailMain: {
    width: 400,
    height: 300,
    overflow: 'hidden',
    flexDirection: 'column'
  },
  detailCells: {
    width: 400,
    height: 300,
    marginTop: 10,
    marginBottom: 10,
    backgroundColor: 'white'
  },
  detailCellsInfo: {
    display: 'flex',
    flexShrink: 0,
    flexGrow: 0,
    overflow: 'hidden',
    flexDirection: 'column'
  },
  detailCellsTitle: {
    paddingTop: 16,
    paddingRight: 16,
    paddingBottom: 8,
    marginLeft: 16,
    fontSize: 14.5,
    fontWeight: 400,
  },
  detailCell: {
    paddingTop: 16,
    paddingRight: 16,
    paddingBottom: 16,
    paddingLeft: 16,
  },
  detailLine: {
    width: 400,
    height: 1,
    marginLeft: 16,
    backgroundColor: "#e0e0e0"
  },
  detailCellContent: {
    overflow: 'auto',
  },
  weuiCellDesc: {
    opacity: '0.5',
    fontSize: 14,
    color: '#000000',
    fontWeight: 400,
    marginBottom: 32,
    wordBreak: 'break-all',
    flexWrap: "wrap"
  },
  weuiCellDescDark: {
    opacity: 1,
    color: '#333333',
    textAlign: 'justify',
  },
  ucAvatarName: {
    fontSize: 14,
    color: '#7D7D8A',
    fontWeight: 400,
    flexDirection: 'row',
  },
  ucAvatar: {
    width: 24,
    height: 24,
    borderRadius: 4,
  },
  ucAvatarSm: {
    width: 20,
    height: 20
  },
  ucAvatarNameAvatar: {
    marginRight: 4,
    verticalAlign: 'middle',
    lineHeight: 1
  },
  ucAvatarNameName: {
    width: 400,
    verticalAlign: 'middle'
  },
  ucCreateTime: {
    width: 400,
    marginBottom: 32,
    fontSize: 14,
    color: '#666666',
    fontWeight: 200,
  },
  itemTagBox: {
    width: 400,
    height: 24,
    flexDirection: 'row',
  },
  tagItem: {
    width: 90,
    flexShrink: 0,
    flexGrow: 0,
    paddingTop: 2,
    paddingRight: 8,
    paddingBottom: 2,
    paddingLeft: 8,
    borderRadius: 7,
    fontSize: 12,
    marginRight: 9,
    textAlign: 'center'
  },
  tagType: {
    backgroundColor: '#e1e6f6',
    color: '#333'
  },
  tagPerson: {
    backgroundColor: '#e1e6f6',
    color: '#333'
  },
  tagLevel: {
    color: '#fff'
  },
  levelOne: {
    backgroundColor: '#ff605c'
  },
  levelTwo: {
    backgroundColor: '#f78532'
  },
  levelThree: {
    backgroundColor: '#3c97cc'
  },
  levelFour: {
    backgroundColor: '#07c160'
  },
  ucContact: {
    paddingTop: 5,
    paddingRight: 5,
    paddingBottom: 5,
    paddingLeft: 5,
    marginTop: 8,
    backgroundColor: '#F5F5F9',
    alignItems: 'flex-start',
    flexDirection: 'row'
  },
  ucContactIcon: {
    backgroundColor: '#FFFFFF',
    width: 48,
    height: 48,
    marginRight: 9
  },
  ucContactDetail: {
    width: 300,
    flexShrink: 1,
    flexGrow: 1,
    flexDirection: 'column'
  },
  ucContactDetailMain: {
    marginTop: 6,
    marginBottom: 20,
    fontSize: 13,
    color: '#000000',
    fontWeight: 400
  },
  ucContactDetailDesc: {
    fontSize: 12,
    color: '#525266',
    fontWeight: 400
  },
  ucDistance: {
    width: 300,
    flexDirection: 'row',
    alignItems: 'flex-start',
    marginTop: 8
  },
  ucDistanceIcon: {
    width: 14,
    height: 14,
    marginRight: 8,
    marginTop: 2,
    verticalAlign: 'middle'
  },
  ucDistanceValue: {
    width: 50,
    fontSize: 14,
    color: '#333333',
    fontWeight: 200,
    marginRight: 12,
    verticalAlign: 'middle'
  },
  ucDistanceAddress: {
    width: 200,
    fontSize: 14,
    color: '#666666',
    fontWeight: 200,
    verticalAlign: 'middle'
  }
};

module.exports = {
  getTemplate,
  templateStyle
}