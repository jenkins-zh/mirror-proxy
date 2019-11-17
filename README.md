## Jenkins Update Center Mirror Proxy

The proxy is a HTTP server which serve for several different Jenkins Update Center sites.

## Service

The Mirror Proxy server has deployed to here.

[https://updates.jenkins-zh.cn/update-center.json](https://updates.jenkins-zh.cn/update-center.json)

## Get started

Run it as demo on MacOS, please follow this:

`make darwin cert run`

On Linux, please follow this:

`make linux cert run-linux`

Run it as a Docker container:

`docker run -v rootCA:/rootCA docker.pkg.github.com/jenkins-zh/mirror-proxy/mirror-proxy:0.0.1 --cert /rootCA/demo.crt --key /rootCA/demo.key`

## API

The only API path is:

|API|Description|
|---|---|
|`GET /update-center.json?version=2.190.2`|Get the update-center.json which allows you give different query conditions|
|`GET /json-servers`|Get all JSON servers|
|`GET /providers`|Get all mirror storage providers|
|`GET /providers/default`|Get the default mirror storage provider|

### Update Center

Below are the query ways for the update center of the mirror:

|Key|Description|
|---|---|
|`version`|The version of current Jenkins|
|`mirror-experimental`|Indicate if you want to use the experimental of plugins|
|`mirror-jsonServer`|Specific the JSON server|
|`mirror-provider`|Specific the mirror storage provider|

**All keys come from query and header. Header value will override the query ones.**
