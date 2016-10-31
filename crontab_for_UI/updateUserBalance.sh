#!/bin/bash

if [ -f /tmp/updateUserBalance.cron ]; then
	echo "Cron already running.....!"
else
	touch /tmp/updateUserBalance.cron
	mysql -u root -p'dexiUkrain123!' -h 162.243.124.60 adsgo_for_native -e 'CALL updateUserBalance();'
	rm -rf /tmp/updateUserBalance.cron
fi
