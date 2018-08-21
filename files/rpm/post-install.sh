chown -R contra /opt/contra
stat /etc/contra.conf || cp /etc/contra.conf.dist /etc/contra.conf
chown contra /etc/contra.conf