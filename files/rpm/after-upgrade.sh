diff /etc/contra.conf.dist /etc/contra.conf.rpmsave > /dev/null 2>&1 || cp /etc/contra.conf.rpmsave /etc/contra.conf
getent passwd contra > /dev/null 2>&1 || useradd -r -d /opt/contra -s /sbin/nologin contra
