getent passwd contra > /dev/null 2>&1 || useradd -r -d /opt/contra -s /sbin/nologin contra
cp /etc/contra.conf /etc/contra.conf.rpmsave || true