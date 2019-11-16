# mirror-proxy
Jenkins Update Center mirror proxy

```bash
openssl genrsa -out demo.key 1024
openssl req -new -x509 -days 1095 -key demo.key \
    -out demo.crt \
    -subj "/C=CN/ST=GD/L=SZ/O=vihoo/OU=dev/CN=demo.com/emailAddress=demo@demo.com"
```