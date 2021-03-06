(function(w,d){
	if(!w.irtb)
	{
		w.irtb = {
			endpoint: '//irtb.io/bid?seat_id={seat_id}&secret_key={secret_key}',
			css: (function(){
				var link = d.createElement("link");
				link.href = '//irtb.io/assets/a.css';
				link.type = "text/css";
				link.rel = "stylesheet";
				d.getElementsByTagName("head")[0].appendChild(link);
			})(),
			createPixel: function(url) {
				return '<iframe src="'+url+'" scrolling="no" frameborder="0" height="0" width="0" style="position:absolute"></iframe>';
			},
			guid: function() { 
				return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) { var r = Math.random()*16|0, v = c == 'x' ? r : (r&0x3|0x8); return v.toString(16); }); 
			},
			page: (function(){
				var ret = document.location.href;
				try
				{
					if (window.top === window.self)
						ret = window.location.href;
					else
						ret = top.window.location.href;
						
					return ret;
				}catch(i){
					if(typeof window.location.ancestorOrigins !== 'undefined')
					{
						var ancestors = window.location.ancestorOrigins;
						if(ancestors.length == 1)
						{
							try
							{
								ret = document.referrer;
							}catch(i){
								ret = ancestors[ancestors.length - 1];
							}
						}
						else
						{
							ret = ancestors[ancestors.length - 1];
						}
					}
					else
					{
						try
						{
							ret = document.referrer;
						}catch(i){
							return ret;
						}
					}
				}
				return ret;
			})(),
			hasAttribute: function (element, attr) {
				try
				{
					if('hasAttribute' in element)
					{
						return element.hasAttribute(attr);
					}
				} catch(e) {}
				try
				{
					var x = element.getAttribute(attr);
					if(x)
						return true;
					else
						return false;
				} catch(e) {
					return false;
				}	
			},
			ajax: (function(){
				var ajax = {};
				ajax.x = function() {
					if (typeof XMLHttpRequest !== 'undefined') {
						return new XMLHttpRequest();  
					}
					var versions = [
						"MSXML2.XmlHttp.6.0",
						"MSXML2.XmlHttp.5.0",   
						"MSXML2.XmlHttp.4.0",  
						"MSXML2.XmlHttp.3.0",   
						"MSXML2.XmlHttp.2.0",  
						"Microsoft.XmlHttp"
					];
				
					var xhr;
					for(var i = 0; i < versions.length; i++) {  
						try {  
							xhr = new ActiveXObject(versions[i]);  
							break;  
						} catch (e) {
						}  
					}
					return xhr;
				};
				ajax.send = function(url, callback, method, data, sync) {
					var x = ajax.x();
					x.open(method, url, sync);
					x.onreadystatechange = function() {
						if (x.readyState == 4) {
							callback(x.responseText)
						}
					};
					if (method == 'POST') {
						x.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
					}
					x.send(data)
				};
				ajax.get = function(url, data, callback, sync) {
					var query = [];
					for (var key in data) {
						query.push(encodeURIComponent(key) + '=' + encodeURIComponent(data[key]));
					}
					ajax.send(url + (query.length ? '?' + query.join('&') : ''), callback, 'GET', null, sync)
				};
				ajax.post = function(url, data, callback, sync) {
					ajax.send(url, callback, 'POST', data, sync)
				};
				return ajax;
			})(),
			chunk: function(input,size,preserve_keys){
				var x, p = '',
					i = 0,
					c = -1,
					l = input.length || 0,
					n = [];
				
				if (size < 1) {
					return null;
				}
				
				if (Object.prototype.toString.call(input) === '[object Array]') {
					if (preserve_keys) {
					while (i < l) {
						(x = i % size) ? n[c][i] = input[i] : n[++c] = {}, n[c][i] = input[i];
						i++;
					}
					} else {
					while (i < l) {
						(x = i % size) ? n[c][x] = input[i] : n[++c] = [input[i]];
						i++;
					}
					}
				} else {
					if (preserve_keys) {
					for (p in input) {
						if (input.hasOwnProperty(p)) {
						(x = i % size) ? n[c][p] = input[p] : n[++c] = {}, n[c][p] = input[p];
						i++;
						}
					}
					} else {
					for (p in input) {
						if (input.hasOwnProperty(p)) {
						(x = i % size) ? n[c][x] = input[p] : n[++c] = [input[p]];
						i++;
						}
					}
					}
				}
				return n;
			},
			ads: {},
			loadAd: function(uniqid){
				var ad = w.irtb.ads[uniqid];
				ad.x = parseInt(ad.x);
				ad.y = parseInt(ad.y);
				
				var request = JSON.stringify({
					id: uniqid,
					imp: [{
						id: '1',
						bidfloor: ad.floor,
						native: {
							request: (function(){
								return JSON.stringify({
                                    plcmtcnt: ad.x * ad.y,
                                    assets: [
                                        {
                                            id: 1,
                                            required: 1,
                                            title: {
                                                len: 70
                                            }
                                        },
                                        {
                                            id: 2,
                                            required: 1,
                                            img: {
                                                type: 3,
                                                w: 300,
                                                h: 150
                                            }
                                        },	
                                        {
                                            id: 3,
                                            required: 1,
                                            data: {
                                                type: 1
                                            }
                                        }
                                    ]
								});
							})()
						}
					}],
					site: {
						page: w.irtb.page,
					},
					device: {
						ua: false,
						ip: false
					}
				});
				
				return w.irtb.ajax.post(w.irtb.endpoint.replace('{seat_id}',ad.seat).replace('{secret_key}',ad.key),request,function(res){
					var bid = null;
					try {
						bid = JSON.parse(res);
					} catch(e) { return; }
					
					var wrapper = document.getElementById(bid.id);
					var placements = [];
					if(wrapper)
					{
						var ads = bid.seatbid;
						for(var a in ads)
						{
							for(var t in ads[a]['bid'])
							{
								var adHolder = ads[a]['bid'][t];
								var native = JSON.parse(adHolder.adm);
								native = native.native;
								
								placements.push({
									impression: adHolder.nurl,
									url: native.link.url,
									sponsor: native.assets[2]['data']['value'],
									title: native.assets[0]['title']['text'],
									image: native.assets[1]['img']['url']
								});
							}
						}
					}
					placements = w.irtb.chunk(placements,ad.x);
					var html = [];
					html.push('<div class="irtb-wrapper">');
					html.push('<h1 class="irtb-yml">You May Like</h1>');
					for(var o in placements)
					{
						var placementGroup = placements[o];
						var width = Math.floor(100/placementGroup.length);
						html.push('<div class="irtb-group">');
						for(var i in placementGroup)
						{
							var placement = placementGroup[i];
							html.push('<div class="irtb-item irtb-col-'+placementGroup.length+'" style="width:'+width+'%;" >');
							html.push('<a href="'+placement.url+'" target="_blank">');
							html.push('<img src="'+placement.image+'" />');
							html.push('<div>'+placement.title+'</div>');
							html.push(w.irtb.createPixel(placement.impression));
							html.push('</a>');
							html.push('</div>');
						}
						html.push('</div>');
					}
					html.push('<div class="irtb-pb">Promoted by <a target="_blank" href="http://adsgo.com">adsGO</a></div>');
					html.push('</div>');
					
					wrapper.innerHTML = html.join('');
				},true);
			}
		};
	}
	var scripts = document.getElementsByTagName('script');
	for(var s=0;s<scripts.length;s++)
	{
		if(w.irtb.hasAttribute(scripts[s],'irtb') == true && w.irtb.hasAttribute(scripts[s],'loaded') == false)
		{
			scripts[s].setAttribute('loaded','1');
			
			var placement = {
				seat: null,
				key: null,
				x: null,
				y: null
			};
			
			for(var p in placement)
			{
				if(w.irtb.hasAttribute(scripts[s],'data-'+p) == true)
					placement[p] =  scripts[s].getAttribute('data-'+p);
				else 
					break;
			}

            if(w.irtb.hasAttribute(scripts[s],'data-floor') == true)
                placement.floor =  scripts[s].getAttribute('data-floor');
            else 
                placement.floor =  0;

			placement.uniqid = 'irtb_' + w.irtb.guid();
			scripts[s].outerHTML+='<div id="'+placement.uniqid+'"></div>';

			w.irtb.ads[placement.uniqid] = placement;
			w.irtb.loadAd(placement.uniqid);

		}
	}
})(window,document);
