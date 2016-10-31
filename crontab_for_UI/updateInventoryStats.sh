#!/bin/bash

if [ -f /tmp/updateInventoryStats.cron ]; then
	echo "Cron already running.....!"
else
	touch /tmp/updateInventoryStats.cron
	mysql -u root -p'dexiUkrain123!' -h 162.243.124.60 adsgo_for_native -e 'CALL updateInventoryStats();'
	rm -rf /tmp/updateInventoryStats.cron
fi
