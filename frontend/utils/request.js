const baseUrl = 'https://mutualaid.ucloud.cn';

const request = options => {
    const app = getApp();

    return new Promise((resolve, reject) => {
        const {
            url,
        } = options;
        const header = options.header || {};
        if (app.globalData.isLogin) {
            header['Authorization'] = `bearer ${app.globalData.jwtToken}`;
        }

        wx.request({
            ...options,
            url: `${baseUrl}${url}`,
            header: {
                'Content-Type': 'application/json',
                ...header
            },
            success: function (res) {
                const {
                    statusCode
                } = res;

                if (statusCode === 200) {
                    resolve(res);
                    return;
                }

                reject( res && res.data && res.data.message);
            },
            fail: function (res) {
                const { status,  statusCode } = res;

                if (statusCode === 403 || status === 403) {
                    reject("未登录或登录过期");
                    app.logout();
                    return;
                }

                reject(res && res.errMsg);
            }
        })
    })
}

module.exports = request;