#!/bin/bash

if [ -f /tmp/updateLowIASBan.cron ]; then
	echo "Cron already running.....!"
else
	touch /tmp/updateLowIASBan.cron
	mysql -u root -p'dexiUkrain123!' -h 162.243.124.60 adsgo_for_native -e 'CALL updateLowIASBan();'
	rm -rf /tmp/updateLowIASBan.cron
fi
