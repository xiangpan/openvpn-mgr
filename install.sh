#!/bin/bash

cd `dirname $0` || exit 1

if [ -f '/etc/redhat-release' ];then
    release='Centos'
fi

if [ -f '/etc/issue' ];then
    grep -q 'Ubuntu' /etc/issue
    if [ $? -eq 0 ];then
        release='Ubuntu'
    fi
fi

if [ -z "$release" ];then
    echo '此脚本仅支持Centos、Ubuntu'
    exit 1
fi

function get_ip() {
    curl ip.cn 2>/dev/null | sed -n 's/[^0-9]*\([0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\).*/\1/p'
}

if [ "$release" == "Centos" ];then
    release_version=`sed -n 's/[^0-9]*\([0-9]\).*/\1/p' /etc/redhat-release`
    if [ -z "${release_version}" ];then
        echo "具体发行版本未知，默认6"
        release_version="6"
    fi
    wget -O /etc/yum.repos.d/epel.repo http://mirrors.aliyun.com/repo/epel-${release_version}.repo
    yum -y install curl expect  openssl-devel openvpn
else
    apt-get -y   install curl expect libssl-dev openvpn
fi

if [ ! -d "/etc/openvpn" ];then
    echo "找不到openvpn安装目录"
    exit 1
else 
    mv /etc/openvpn /etc/openvpn-`date +"%s"`
    mkdir /etc/openvpn
fi


openssl dhparam -out dh2048.pem 2048
openvpn --genkey --secret ta.key


./create_crt.sh
if [ $? -ne 0 ];then
    exit 1
fi

install -m 600  *key *crt *pem  /etc/openvpn/
install -d  /etc/openvpn/instance
install -d  /etc/openvpn/template
install -m 755 instance/*  /etc/openvpn/instance
install -m 755 template/*  /etc/openvpn/template


mgrfile="openvpn-mgr"
if [ "$release" == "Centos" ];then
    release_version=`sed -n 's/[^0-9]*\([0-9]\).*/\1/p' /etc/redhat-release`
    if [ -z "${release_version}" ];then
        echo "具体发行版本未知，默认6"
        release_version="6"
    fi
    wget -O  /etc/openvpn/$mgrfile http://www.smnode.com/smstatic/${mgrfile}/${mgrfile}_${release_version}
else 
    wget -O  /etc/openvpn/$mgrfile http://www.smnode.com/smstatic/${mgrfile}/${mgrfile}_7
fi

if [ -f "/etc/openvpn/$mgrfile" ];then
    chmod 755 /etc/openvpn/$mgrfile
    /etc/openvpn/$mgrfile -c initcfg
else 
    echo "安装失败"
    exit 1
fi

