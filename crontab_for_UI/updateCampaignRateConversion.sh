#!/bin/bash

if [ -f /tmp/updateCampaignRateConversion.cron ]; then
	echo "Cron already running.....!"
else
	touch /tmp/updateCampaignRateConversion.cron
	mysql -u root -p'dexiUkrain123!' -h 162.243.124.60 adsgo_for_native -e 'CALL updateCampaignRateConversion();'
	rm -rf /tmp/updateCampaignRateConversion.cron
fi
