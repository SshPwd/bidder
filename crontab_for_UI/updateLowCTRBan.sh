#!/bin/bash

if [ -f /tmp/updateLowCTRBan.cron ]; then
	echo "Cron already running.....!"
else
	touch /tmp/updateLowCTRBan.cron
	mysql -u root -p'dexiUkrain123!' -h 162.243.124.60 adsgo_for_native -e 'CALL updateLowCTRBan();'
	rm -rf /tmp/updateLowCTRBan.cron
fi
