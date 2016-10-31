#!/bin/bash

if [ -f /tmp/updateCreativeCTR.cron ]; then
	echo "Cron already running.....!"
else
	touch /tmp/updateCreativeCTR.cron
	mysql -u root -p'dexiUkrain123!' -h 162.243.124.60 adsgo_for_native -e 'CALL updateCreativeCTR();'
	rm -rf /tmp/updateCreativeCTR.cron
fi
