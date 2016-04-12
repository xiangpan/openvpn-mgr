#!/bin/bash
mkdir -p  `dirname $0`/client
cd  `dirname $0`/client || exit 
real_ip=`curl  ip.cn 2>/dev/null | sed -n 's/[^0-9]*\([0-9.]*\).*/\1/p'`

cp ../template/vpn-server.ovpn  .
cp ../ta.key ta_smnode.key
cp ../client.key client_smnode.key
cp ../client.crt client_smnode.crt
cp ../ca.crt ca_smnode.crt
sed -i "s/SED_IP/$real_ip/" vpn-server.ovpn
