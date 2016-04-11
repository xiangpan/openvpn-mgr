#!/bin/bash

#生成证书请求
function create_csr() {
    if [ ! -f "${1}.key" ];then
        create_key $1
    fi
    sleep 1
    SSL_Common_Name=${SSL_Common_Name}`date +"%s"`
expect <<EOF
        spawn openssl req -new -key $1.key -out $1.csr
        expect "Country Name"
        send "${SSL_CN}\n"
        expect "State or Province Name"
        send "${SSL_Province_Name}\n"
        expect "Locality Name"
        send "${SSL_Locality_Name}\n"
        expect "Organization Name"
        send "${SSL_Org_Name}\n"
        expect "Organizational Unit Name"
        send "${SSL_Org_UN}\n"
        expect "Common Name"
        send "${SSL_Common_Name}\n"
        expect "Email Address"
        send "${SSL_Email}\n"
        expect "A challenge password"
        send "\n"
        expect "An optional company name"
        send "\n"
        expect eof
EOF
}

#生成自签名证书
function create_ca() {
    if [ ! -f "ca.key" ];then
        create_key ca
    fi

expect <<EOF
        spawn openssl req -new -x509 -days 3650 -key ca.key -out ca.crt
        expect "Country Name"
        send "${SSL_CN}\n"
        expect "State or Province Name"
        send "${SSL_Province_Name}\n"
        expect "Locality Name"
        send "${SSL_Locality_Name}\n"
        expect "Organization Name"
        send "${SSL_Org_Name}\n"
        expect "Organizational Unit Name"
        send "${SSL_Org_UN}\n"
        expect "Common Name"
        send "${SSL_Common_Name}\n"
        expect "Email Address"
        send "${SSL_Email}\n"

        expect eof
EOF
}

#生成证书
function create_crt() {
    if [ ! -f "${1}.key" ];then
        create_key $1
    fi
    if [ ! -f "${1}.csr" ];then
        create_csr $1
    fi
    openssl x509 -req -days 3650 -in ${1}.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out ${1}.crt
    openssl verify -CAfile ca.crt ${1}.crt
}


#生成key
function create_key() {
    openssl genrsa -out ${1}.key 2048
}

which expect 2>/dev/null
if [ $? -ne 0 ];then
    echo "此脚本需要 expect 支持"
    exit 1
fi
export SSL_CN="CN"
export SSL_Province_Name="SH"
export SSL_Locality_Name="ShangHai"
export SSL_Org_Name="www.smnode.com"
export SSL_Org_UN="smnode"
export SSL_Common_Name="smnode"
export SSL_Email="support@smnode.com"

create_ca
create_crt server
create_crt client
