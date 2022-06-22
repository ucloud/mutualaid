env: {{ .ENV }} 
server:
  http:
    addr: 0.0.0.0:28000
    timeout: 5s
  grpc:
    addr: 0.0.0.0:29000
    timeout: 5s
data:
  database:
    driver: mysql
    source: {{ .DB_USER }}:{{ .DB_PASSWORD }}@tcp({{ .DB_HOST }}:{{ .DB_PORT }})/{{ .DB_NAME }}?charset=utf8mb4&parseTime=True&loc=Local
    slow_threshold: 1s
  redis:
    addrs:
      - {{ .REDIS_NODE1 }}:6379
      - {{ .REDIS_NODE2 }}:6379
      - {{ .REDIS_NODE3 }}:6379
      - {{ .REDIS_NODE4 }}:6379
    password:
    pool_size: 60
    min_idle_conn: 10
    read_timeout: 0.2s
    write_timeout: 0.2s

proxy:
  wxkey:
    WPAppID: {{ .WECHAT_WP_APP_ID }}
    WPSecret: {{ .WECHAT_WP_SECRET }}
    AppID: {{ .WECHAT_APP_ID }}
    Secret: {{ .WECHAT_SECRET }}
    BaseURL: "https://api.weixin.qq.com"
    endpoint: "endpoint"
    msgTpl:
      帮助被采纳提醒:
        id: dKoQRX2JzYmJdQ1mn3sJpULGNc9OJfhTst6lciVOsGg
        args:
          求助编号: character_string1
          求助内容: thing2
          回复内容: thing3
          采纳时间: time4
      收到帮助提醒:
        id: _rHvH3tzUO-PKfYbvMzCWtvWtqPJakEgOGyz4RQysN4
        args:
          求助编号: character_string1
          求助内容: thing2
          回复人: thing3
          回复内容: thing4
          回复时间: time5
bizConfig:
