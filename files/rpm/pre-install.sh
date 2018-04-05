getent passwd contra > /dev/null 2>&1 || useradd -r -d /opt/contra -s /sbin/nologin contra
mv /etc/contra.conf /etc/contra.conf.rpmsave || true