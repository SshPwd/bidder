#!/bin/bash

if [ -f /tmp/updateFeatureBreakdown.cron ]; then
	echo "Cron already running.....!"
else
	touch /tmp/updateFeatureBreakdown.cron
	mysql -u root -p'dexiUkrain123!' -h 162.243.124.60 adsgo_for_native -e 'CALL updateFeatureBreakdown();'
	rm -rf /tmp/updateFeatureBreakdown.cron
fi
