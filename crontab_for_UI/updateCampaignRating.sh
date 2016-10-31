#!/bin/bash

if [ -f /tmp/updateCampaignRating.cron ]; then
	echo "Cron already running.....!"
else
	touch /tmp/updateCampaignRating.cron
	mysql -u root -p'dexiUkrain123!' -h 162.243.124.60 adsgo_for_native -e 'CALL updateCampaignRating();'
	rm -rf /tmp/updateCampaignRating.cron
fi
