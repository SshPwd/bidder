#!/bin/bash

if [ -f /tmp/updateInventoryCTR.cron ]; then
	echo "Cron already running.....!"
else
	touch /tmp/updateInventoryCTR.cron
	mysql -u root -p'dexiUkrain123!' -h 162.243.124.60 adsgo_for_native -e 'CALL updateInventoryCTR();'
	rm -rf /tmp/updateInventoryCTR.cron
fi
