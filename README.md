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

`/update-center.json?version=`
