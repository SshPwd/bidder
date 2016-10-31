#!/bin/bash

if [ -f /tmp/updateInventoryCTROptimizationTables.cron ]; then
	echo "Cron already running.....!"
else
	touch /tmp/updateInventoryCTROptimizationTables.cron
	mysql -u root -p'dexiUkrain123!' -h 162.243.124.60 adsgo_for_native -e 'CALL updateInventoryCTROptimizationTables();'
	rm -rf /tmp/updateInventoryCTROptimizationTables.cron
fi
