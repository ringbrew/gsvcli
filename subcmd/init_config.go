package subcmd

const ConfigTmpl = `environment: test
debug: true

port:
  rpc: 3000
  proxy: 3001
  api: 3002
  admin: 3003

mysql:
  user_name: root
  password: 123456
  host: localhost:3306
  database: test

redis:
  host: 127.0.0.1:6379
  password: redispass
  db: 0

trace:
  type: ''
  endpoint: ''
`
