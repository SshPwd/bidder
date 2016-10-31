#!/bin/bash

if [ -f /tmp/updateAdvertiserReportingDay.cron ]; then
	echo "Cron already running.....!"
else
	touch /tmp/updateAdvertiserReportingDay.cron
	mysql -u root -p'dexiUkrain123!' -h 162.243.124.60 adsgo_for_native -e 'CALL updateAdvertiserReportingDay();'
	rm -rf /tmp/updateAdvertiserReportingDay.cron
fi
