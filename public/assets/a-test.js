(function(w,d){
	if(!w.irtb)
	{
		w.irtb = {
			endpoint: '/bid?seat_id={seat_id}&secret_key={secret_key}',
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
            curGuid: false,
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
            ready: (function(){    

                var readyList,
                    DOMContentLoaded,
                    class2type = {};
                    class2type["[object Boolean]"] = "boolean";
                    class2type["[object Number]"] = "number";
                    class2type["[object String]"] = "string";
                    class2type["[object Function]"] = "function";
                    class2type["[object Array]"] = "array";
                    class2type["[object Date]"] = "date";
                    class2type["[object RegExp]"] = "regexp";
                    class2type["[object Object]"] = "object";

                var ReadyObj = {
                    // Is the DOM ready to be used? Set to true once it occurs.
                    isReady: false,
                    // A counter to track how many items to wait for before
                    // the ready event fires. See #6781
                    readyWait: 1,
                    // Hold (or release) the ready event
                    holdReady: function( hold ) {
                        if ( hold ) {
                            ReadyObj.readyWait++;
                        } else {
                            ReadyObj.ready( true );
                        }
                    },
                    // Handle when the DOM is ready
                    ready: function( wait ) {
                        // Either a released hold or an DOMready/load event and not yet ready
                        if ( (wait === true && !--ReadyObj.readyWait) || (wait !== true && !ReadyObj.isReady) ) {
                            // Make sure body exists, at least, in case IE gets a little overzealous (ticket #5443).
                            if ( !document.body ) {
                                return setTimeout( ReadyObj.ready, 1 );
                            }

                            // Remember that the DOM is ready
                            ReadyObj.isReady = true;
                            // If a normal DOM Ready event fired, decrement, and wait if need be
                            if ( wait !== true && --ReadyObj.readyWait > 0 ) {
                                return;
                            }
                            // If there are functions bound, to execute
                            readyList.resolveWith( document, [ ReadyObj ] );

                            // Trigger any bound ready events
                            //if ( ReadyObj.fn.trigger ) {
                            //  ReadyObj( document ).trigger( "ready" ).unbind( "ready" );
                            //}
                        }
                    },
                    bindReady: function() {
                        if ( readyList ) {
                            return;
                        }
                        readyList = ReadyObj._Deferred();

                        // Catch cases where $(document).ready() is called after the
                        // browser event has already occurred.
                        if ( document.readyState === "complete" ) {
                            // Handle it asynchronously to allow scripts the opportunity to delay ready
                            return setTimeout( ReadyObj.ready, 1 );
                        }

                        // Mozilla, Opera and webkit nightlies currently support this event
                        if ( document.addEventListener ) {
                            // Use the handy event callback
                            document.addEventListener( "DOMContentLoaded", DOMContentLoaded, false );
                            // A fallback to window.onload, that will always work
                            window.addEventListener( "load", ReadyObj.ready, false );

                        // If IE event model is used
                        } else if ( document.attachEvent ) {
                            // ensure firing before onload,
                            // maybe late but safe also for iframes
                            document.attachEvent( "onreadystatechange", DOMContentLoaded );

                            // A fallback to window.onload, that will always work
                            window.attachEvent( "onload", ReadyObj.ready );

                            // If IE and not a frame
                            // continually check to see if the document is ready
                            var toplevel = false;

                            try {
                                toplevel = window.frameElement == null;
                            } catch(e) {}

                            if ( document.documentElement.doScroll && toplevel ) {
                                doScrollCheck();
                            }
                        }
                    },
                    _Deferred: function() {
                        var // callbacks list
                            callbacks = [],
                            // stored [ context , args ]
                            fired,
                            // to avoid firing when already doing so
                            firing,
                            // flag to know if the deferred has been cancelled
                            cancelled,
                            // the deferred itself
                            deferred  = {

                                // done( f1, f2, ...)
                                done: function() {
                                    if ( !cancelled ) {
                                        var args = arguments,
                                            i,
                                            length,
                                            elem,
                                            type,
                                            _fired;
                                        if ( fired ) {
                                            _fired = fired;
                                            fired = 0;
                                        }
                                        for ( i = 0, length = args.length; i < length; i++ ) {
                                            elem = args[ i ];
                                            type = ReadyObj.type( elem );
                                            if ( type === "array" ) {
                                                deferred.done.apply( deferred, elem );
                                            } else if ( type === "function" ) {
                                                callbacks.push( elem );
                                            }
                                        }
                                        if ( _fired ) {
                                            deferred.resolveWith( _fired[ 0 ], _fired[ 1 ] );
                                        }
                                    }
                                    return this;
                                },

                                // resolve with given context and args
                                resolveWith: function( context, args ) {
                                    if ( !cancelled && !fired && !firing ) {
                                        // make sure args are available (#8421)
                                        args = args || [];
                                        firing = 1;
                                        try {
                                            while( callbacks[ 0 ] ) {
                                                callbacks.shift().apply( context, args );//shifts a callback, and applies it to document
                                            }
                                        }
                                        finally {
                                            fired = [ context, args ];
                                            firing = 0;
                                        }
                                    }
                                    return this;
                                },

                                // resolve with this as context and given arguments
                                resolve: function() {
                                    deferred.resolveWith( this, arguments );
                                    return this;
                                },

                                // Has this deferred been resolved?
                                isResolved: function() {
                                    return !!( firing || fired );
                                },

                                // Cancel
                                cancel: function() {
                                    cancelled = 1;
                                    callbacks = [];
                                    return this;
                                }
                            };

                        return deferred;
                    },
                    type: function( obj ) {
                        return obj == null ?
                            String( obj ) :
                            class2type[ Object.prototype.toString.call(obj) ] || "object";
                    }
                }
                // The DOM ready check for Internet Explorer
                function doScrollCheck() {
                    if ( ReadyObj.isReady ) {
                        return;
                    }

                    try {
                        // If IE is used, use the trick by Diego Perini
                        // http://javascript.nwbox.com/IEContentLoaded/
                        document.documentElement.doScroll("left");
                    } catch(e) {
                        setTimeout( doScrollCheck, 1 );
                        return;
                    }

                    // and execute any waiting functions
                    ReadyObj.ready();
                }
                // Cleanup functions for the document ready method
                if ( document.addEventListener ) {
                    DOMContentLoaded = function() {
                        document.removeEventListener( "DOMContentLoaded", DOMContentLoaded, false );
                        ReadyObj.ready();
                    };

                } else if ( document.attachEvent ) {
                    DOMContentLoaded = function() {
                        // Make sure body exists, at least, in case IE gets a little overzealous (ticket #5443).
                        if ( document.readyState === "complete" ) {
                            document.detachEvent( "onreadystatechange", DOMContentLoaded );
                            ReadyObj.ready();
                        }
                    };
                }
                function ready( fn ) {
                    // Attach the listeners
                    ReadyObj.bindReady();

                    var type = ReadyObj.type( fn );

                    // Add the callback
                    readyList.done( fn );//readyList is result of _Deferred()
                }
                return ready;
            })(),
			loadAds: function(){
				var request = {
					id: w.irtb.curGuid,
					imp: [],
					site: {
						page: (function(){
                            if(document.location.hash)
                                return 'http://' + document.location.hash.replace('#','');
                            else
                                return w.irtb.page;
                        })()
					},
					device: {
						ua: null,
						ip: null
					}
				};
				
                var readyToLoad = false;
                
                for(var uniqid in w.irtb.ads)
                {
                    if(w.irtb.ads[uniqid].loaded)
                        continue;
                    else
                        w.irtb.ads[uniqid].loaded = true;
                    
                    readyToLoad = true;
                    
                    var ad = w.irtb.ads[uniqid];
                    ad.x = parseInt(ad.x);
                    ad.y = parseInt(ad.y);
                    
                    request.imp.push({
						id: uniqid.toString(),
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
                                                h: 250
                                            }
                                        },	
                                        {
                                            id: 3,
                                            required: 1,
                                            data: {
                                                type: 1
                                            }
                                        },
                                        {
                                            id: 4,
                                            required: 1,
                                            img: {
                                                type: 2,
                                                h: 35
                                            }
                                        },
                                        {
                                            id: 5,
                                            required: 1,
                                            data: {
                                                type: 2
                                            }
                                        },	
                                    ]
								});
							})()
						}
					});
                }
                
                if(readyToLoad == false)
                    return;

				return w.irtb.ajax.post(w.irtb.endpoint.replace('{seat_id}',ad.seat).replace('{secret_key}',ad.key),JSON.stringify(request),function(res){
					var bid = null;
					try {
						bid = JSON.parse(res);
					} catch(e) { return; }
					

                    var ads = bid.seatbid;
   
                    for(var a in ads)
                    {
                        var wrapper = document.getElementById('irtb_' + ads[a].bid[0].impid);
                        var ad = w.irtb.ads[ads[a].bid[0].impid];
                        var placements = [];
                        
                        for(var t in ads[a]['bid'])
                        {
                            var adHolder = ads[a]['bid'][t];
                            var native = JSON.parse(adHolder.adm);
                            native = native.native;
                            
                            placements.push({
                                impression: adHolder.nurl,
                                url: native.link.url,
                                sponsor: native.assets[2]['data']['value'],
                                title: native.assets[0]['title']['text'] + ' ('+adHolder.price+')',
                                desc: native.assets[4]['data']['value'],
                                image: native.assets[1]['img']['url'],
                                brand_logo: native.assets[3]['img']['url']
                            });
                        }
                        
                        
                        placements = w.irtb.chunk(placements,ad.x);
                        var html = [];
                        html.push('<div class="irtb-wrapper">');
                        html.push('<h1 class="irtb-yml">You May Like</h1>');
                        for(var o in placements)
                        {
                            var placementGroup = placements[o];
                            var width = Math.floor(100/ad.x); //Math.floor(100/placementGroup.length);
                            html.push('<div class="irtb-group">');
                            for(var i in placementGroup)
                            {
                                var placement = placementGroup[i];
                                html.push('<div class="irtb-item irtb-col-'+placementGroup.length+'" style="width:'+width+'%;" >');
                                html.push('<a href="'+placement.url+'" target="_blank">');
                                html.push('<div class="irtb-img-wrapper">');
                                html.push('<img src="'+placement.image+'" />');
                                html.push('<img class="irtb-brand-logo" src="'+placement.brand_logo+'"></img>');
                                html.push('<div class="irtb-sponsor">Sponsored By '+placement.sponsor+'</div>');
                                html.push('</div>');
                                html.push('<div class="irtb-title">'+placement.title+'</div>');
                                html.push('<div class="irtb-caption">'+placement.desc+'</div>');
                                html.push(w.irtb.createPixel(placement.impression));
                                html.push('</a>');
                                html.push('</div>');
                            }
                            html.push('</div>');
                        }
                        html.push('<div class="irtb-pb">Promoted by <a target="_blank" href="http://adsgo.com">adsGO</a></div>');
                        html.push('</div>');
                        
                        wrapper.innerHTML = html.join('');
                        
                        
                    }
                

				},true);
			}
		};
        w.irtb.curGuid = w.irtb.guid();
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

			placement.uniqid = Object.keys(w.irtb.ads).length;
			scripts[s].outerHTML+='<div id="'+ 'irtb_' + Object.keys(w.irtb.ads).length +'"></div>';

			w.irtb.ads[placement.uniqid] = placement;
            w.irtb.ready(function(){
                w.irtb.loadAds();
            });
			//w.irtb.loadAd(placement.uniqid);
		}
	}
})(window,document);
