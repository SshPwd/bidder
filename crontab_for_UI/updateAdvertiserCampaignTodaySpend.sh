#!/bin/bash

if [ -f /tmp/updateAdvertiserCampaignTodaySpend.cron ]; then
	echo "Cron already running.....!"
else
	touch /tmp/updateAdvertiserCampaignTodaySpend.cron
	mysql -u root -p'dexiUkrain123!' -h 162.243.124.60 adsgo_for_native -e 'CALL updateAdvertiserCampaignTodaySpend();'
	rm -rf /tmp/updateAdvertiserCampaignTodaySpend.cron
fi
