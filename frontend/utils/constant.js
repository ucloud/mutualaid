// 求助类型
const AidType = [
    { label: "食品生活物资", value: 10 },
    { label: "就医", value: 20 },
    { label: "求药", value: 30 },
    { label: "防疫物资", value: 40 },
    { label: "隔离求助", value: 50 },
    { label: "心理援助", value: 60 },
    { label: "其他", value: 99 },
];

// 求助人群
const AidGroup = [
    { label: "重症患者", value: 1 },
    { label: "儿童婴儿", value: 2 },
    { label: "孕妇", value: 3 },
    { label: "老人", value: 40 },
    { label: "残障", value: 50 },
    { label: "外来务工人员", value: 60 },
    { label: "滞留人员", value: 70 },
    { label: "新冠阳性", value: 80 },
    { label: "医护工作者", value: 90 },
    { label: "街道社区", value: 100 },
    { label: "外籍人士", value: 110 },
    { label: "其他", value: 120 },
];

// 紧急程度
const EmergencyLevel = [
    { label: "威胁生命", value: 1 },
    { label: "威胁健康", value: 2 },
    { label: "处境困难", value: 3 },
    { label: "暂无危险", value: 4 },
];

// 任务状态
const Status = [
    { label: "已创建", value: 10 },
    { label: "已关闭", value: 15 },
    { label: "已解决", value: 20 },
];

const getTagText = (type, value) => {
    if (!value || !type) return "";
    const typeMap = {
        type: AidType,
        group: AidGroup,
        emergency: EmergencyLevel,
        status: Status
    };
    const sourceData = typeMap[type];
    const filterResult = sourceData.filter((item) => item["value"] === Number(value));
    return (filterResult.length && filterResult[0].label) || "";
};

export { AidType, AidGroup, EmergencyLevel, Status, getTagText };
