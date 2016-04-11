# openvpn-mgr

openvpn-mgr 是开源的openvpn管理工具


##版权声明
    本软件及文档遵从AGPL许可协议, 版权属于上海守茂网络科技有限公司。

###用途：
    openvpn 自动化安装、用户管理     	

    需要expect软件支持

###使用说明：
    默认所有文件生成后，会复制到 /etc/openvpn/目录
    ./openvpn-mgr -h 显示帮助信息
    Usage of ./openvpn-mgr:
      -c string
            指定执行命令(auth list search add del modify initcfg (default "auth")
      -p string
            密码
      -u string
            用户

      -c string
            指定执行命令(auth list search add del modify initcfg (default "auth")

        auth: 根据openvpn via-env规则，认证用户.
        list: 显示所有用户
        del:  删除用户
        modify: 修改用户
        search：查询用户 
        initcfg: openvpn安装后，生成openvpn配置文件 
    



###配置文件详细说明：
        



###一些说明：

    systemd 请使用下面命令设置openvpn自启动
        systemctl -f enable openvpn@server.service
        systemctl start openvpn@server.service

##安装

    服务端安装
    1. curl http://www.smnode.com/smstatic/install_openvpn_mgr.sh | bash
    2. 运行: /etc/openvpn/openvpn-mgr  -c initcfg 生成server.conf文件 
    3. 为了使vpn用户能访问其它服务器，vpn server端需要做snat 参考脚本  /etc/openvpn/instance/snat-rules.iptables

    客户端安装
    1. 下载openvpn客户端软件 
        https://openvpn.net/index.php/open-source/downloads.html

    2. 运行/etc/openvpn/build_client.sh
        将client目录中的文件，复制到客户端的config目录中

    
##交流QQ群
    323865245


## 守茂网络科技 Copyright
