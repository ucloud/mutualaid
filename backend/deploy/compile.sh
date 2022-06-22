git config --global url."https://gitlab-ci-token:${CI_JOB_TOKEN}@git.ucloudadmin.com".insteadOf "https://git.ucloudadmin.com"
export GOPROXY=https://goproxy.cn,direct
export GOPRIVATE=git.ucloudadmin.com
go mod download -x
make build

sed \
  -e "s^{{ .ENV }}^prod^" \
  -e "s^{{ .DB_USER }}^$DB_USER_PROD^" \
  -e "s^{{ .DB_PASSWORD }}^$DB_PASSWORD_PROD^" \
  -e "s^{{ .DB_HOST }}^$DB_HOST_PROD^" \
  -e "s^{{ .DB_PORT }}^$DB_PORT_PROD^" \
  -e "s^{{ .DB_NAME }}^$DB_NAME_PROD^" \
  -e "s^{{ .REDIS_NODE1 }}^$REDIS_NODE1_PROD^" \
  -e "s^{{ .REDIS_NODE2 }}^$REDIS_NODE2_PROD^" \
  -e "s^{{ .REDIS_NODE3 }}^$REDIS_NODE3_PROD^" \
  -e "s^{{ .REDIS_NODE4 }}^$REDIS_NODE4_PROD^" \
  -e "s^{{ .REDIS_PASSWORD }}^$REDIS_PASSWORD_PROD^" \
  -e "s^{{ .WECHAT_APP_ID }}^$WECHAT_APP_ID_PROD^" \
  -e "s^{{ .WECHAT_SECRET }}^$WECHAT_SECRET_PROD^" \
  -e "s^{{ .WECHAT_WP_APP_ID }}^$WECHAT_WP_APP_ID_PROD^" \
  -e "s^{{ .WECHAT_WP_SECRET }}^$WECHAT_WP_SECRET_PROD^" \
  ./configs/config.yaml.tpl > ./bin/config.yaml
