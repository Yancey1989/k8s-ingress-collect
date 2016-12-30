# Collect Nginx vts module stat data
## Write to influxdb with the [rest api](https://docs.influxdata.com/influxdb/v1.1/guides/writing_data/)
## Step 1
`git clone https://github.com/Yancey1989/k8s-ingress-collect.git && cd k8s-ingress-collect`

## Step 2
Configurate enviroment variable.

Key | Value
---|---
NGINX_HOST | Nginx host
NGINX_PORT | Nginx stats port
INFLUX_DB_HOST | InfluxDB host
INFLUX_DB_PORT | InfluxDB api port
INFLUX_DB_NAME | Which influxdb database
INTERVAL | Interval for collecting
