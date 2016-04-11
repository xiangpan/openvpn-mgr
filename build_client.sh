#!/bin/bash
cd  `dirname $0`/client || exit 
real_ip=`curl  ip.cn 2>/dev/null | sed -n 's/[^0-9]*\([0-9.]*\).*/\1/p'`

cp ../template/vpn-server.ovpn  .
cp ../{ta.key,client.key,client.crt,ca.crt}  .
sed -i "s/SED_IP/$real_ip/" vpn-server.ovpn
