var iflex_spa_callbacks = iflex_spa_callbacks || {};

jQuery(function() {

    var href = window.location.href,
        has_protocol = new RegExp('^(http|https):','i'),
        is_local = new RegExp('^(http:|https:|)//' + window.location.host,'i'),
        is_download = new RegExp('\.(iso|torrent|sig|zip)$'),
        wp_admin = new RegExp('(/wp-admin/|/wp-login.php)'),

    //  jQuery_monkeypatch - add error trapping to the globalEval call
    jQuery_monkeypatch = function() {
        jQuery.globalEval = function(b) {
            var a = window;
            b && jQuery.trim(b) && (a.execScript || function(b) {
                try {
                    a.eval.call(a, b)
                } catch(error) {
                    console.log(error, b);
                }
            })(b)
        }
    },

    //  scroll_top - reposition the page on a given #tag
    scroll_top = function(tag) {
        if(jQuery('#'+tag)) {
            node = jQuery('#'+tag).offset();
            jQuery('html, body').animate({scrollTop: node.top + 'px'}, 'fast');
        }
    },


    //  set_position - reposition the page based on the current #tag and link
    set_position = function(tag,htag) {
        if(htag&&(jQuery('#'+htag).length)) return scroll_top(htag);
        if(tag=='href') jQuery('html,body').scrollTop(0);
    },

    //  move_scripts - move a bunch of scripts to a new target location
    move_scripts = function(scripts, target) {
        jQuery.each(scripts, function() { 
            var script = document.createElement('script'), node = target.lastChild;
            if(this.src) {
                script.type = this.type;
                script.src = this.src;
            } else {
                script.innerHTML = this.innerHTML;
            }
            node.parentNode.insertBefore(script, node.nextSibling);
        });
    },

    //  reload_head - add new entries to page head and remove obsoletes
    reload_head = function(head) {
        var live = jQuery('head').find('*');
        jQuery.each(jQuery(head).find('*'), function() {
            var raw = jQuery(this).get()[0], i=live.length;
            while(--i >= 0) {
                if(raw.isEqualNode(jQuery(live[i]).get()[0])) {
                    live.splice(i,1);
                    break;
                }
            }
            if(i<0) jQuery('head').append(this);
        });
        jQuery.each(jQuery('head').find('*'), function() {
            var raw = jQuery(this).get()[0], old = this;
            jQuery.each(live, function(k,v) {
                if(raw.isEqualNode(jQuery(v).get()[0])) jQuery(old).remove();
            });
        });
    },

    //  load_page - bulk of the work, replace the current page with a new one
    load_page = function(tag, href, data, navigate) {
        var target, base, url, htag,
        replace_page = function(data,status,xhr) {
            var mtype,stype,content_type = xhr.getResponseHeader('Content-Type').split(";")[0];
            [mtype,stype] = content_type.split("/");
            switch(mtype) {
                case 'text':
                    var url, htag, head, body, hscripts, bscripts,
                    node = document.createElement("html");
                    window.onload = null;   // make sure this is clear for things like SMF
                    node.innerHTML = data;
                    head = node.getElementsByTagName("head")[0];
                    document.title = jQuery(head).find('title').text(); 
                    body = node.getElementsByTagName("body")[0];
                    hscripts = jQuery(head).find('script').remove().get();
                    bscripts = jQuery(body).find('script').remove().get();
                    reload_head(head);
                    move_scripts(hscripts, body);
                    move_scripts(bscripts, body);
                    jQuery_monkeypatch();
                    jQuery("body").html(body.innerHTML);
                    [url,htag] = href.split('#');
                    set_position(tag,htag);
                    setTimeout(function(){jQuery(window).trigger('load');},100);
                    break;
                case 'application':
                case 'image':
                    var win = window.open(href, '_blank');
                    win.focus();
                    return false;
                default:
                    alert('Not configured handle file type: ', content_type);
                    return false;
            }
            if(navigate) history.pushState({href: href}, null, href);
            jQuery.each(iflex_spa_callbacks, function(name,routine){
                routine();
            });
        };
        switch(tag) {
            case 'action':
                target = href;
                href = window.location.href;
                jQuery.ajax({
                    url: target,
                    data: data,
                    type: 'POST',
                    success: replace_page,
                    processData: false,
                    contentType: false,
                    cache: false
                });
                break;
            case 'href':
                base = window.location.href.split('#')[0];
                [url,htag] = href.split('#');
                if(htag&&(!url||url==base)) return scroll_top(htag);
                if(!htag||(url&&(url!=base))) jQuery.get(href, replace_page).fail(function(){
                    alert('Unable to locate url '+href);
                });
        }
    },

    //  outside - work out if a given URL is going to cause us to leave site
    outside = function(href) {
        if(has_protocol.test(href) && !is_local.test(href)) return true;
        if(is_download.test(href)||wp_admin.test(href)) return true;
        return false;
    },
   
    //  wrap - assign a new handler to an element (if appropriate)
    wrap = function(self, tag, handler) {

        //  if we have no url, or the url is outside the site ...
        var href = jQuery(self).attr(tag);
        if(!href||(href=='#')||outside(href)) return false;

        //  add a handler to override the default process
        jQuery(self).on(handler, function(event){
            var not_supported = false,
                data = null,
                button = null;

            if(event.isDefaultPrevented()) return false;
            switch(tag) {
                case 'action':

                    //  for POST requests we're going to use the results from FormData, and add the name/value of the submit button
                    data = new FormData(this);
                    button = document.activeElement;
                    if(button.name) data.append(button.name, button.value)
                    break;
                case 'href':
                    break;
                default:
                    not_supported = true;
                    break;
            }
            if(not_supported) return;
            event.preventDefault();
            load_page(tag,href,data,true);
        });
    },

    //  back_button - retask the back button to do a soft back if within site
    back_button = function() {
        //  If we're run out of history, exit the site
        if(!history.state) return window.history.back();
        //  Otherwise, implement our own version of bac
        load_page('href', history.state.href, null, false);
    };

    // wrap all A links for FORMs
    jQuery('a').each(function(){wrap(this,'href','click')});
    jQuery('form').each(function(){wrap(this,'action','submit')});

    // wrap the browser's 'BACK' button
    jQuery(window).off('popstate').on('popstate', back_button);
    if(history.state) return;

    // record first page of the session
    history.pushState({href:href},null,href);
});