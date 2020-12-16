#!/bin/bash
umask 0002
chmod g+w /data
chown -R qumine:qumine /data
exec qumine-server server