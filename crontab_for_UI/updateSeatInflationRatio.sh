#!/bin/bash

if [ -f /tmp/updateSeatInflationRatio.cron ]; then
	echo "Cron already running.....!"
else
	touch /tmp/updateSeatInflationRatio.cron
	mysql -u root -p'dexiUkrain123!' -h 162.243.124.60 adsgo_for_native -e 'CALL updateSeatInflationRatio();'
	rm -rf /tmp/updateSeatInflationRatio.cron
fi
