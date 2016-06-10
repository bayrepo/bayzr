<!DOCTYPE html>
<html>
<title>{{.ReportName}}</title>
<meta name="viewport" content="width=device-width, initial-scale=1">

<style>
/* W3.CSS 2.52 by Jan Egil and Borge Refsnes. Do not remove this line */
*{-webkit-box-sizing:border-box;-moz-box-sizing:border-box;box-sizing:border-box}
/* Extract from normalize.css by Nicolas Gallagher and Jonathan Neal git.io/normalize */
html{-ms-text-size-adjust:100%;-webkit-text-size-adjust:100%}body{margin:0}
article,aside,details,figcaption,figure,footer,header,hgroup,main,menu,nav,section,summary{display:block}
audio,canvas,video{display:inline-block;vertical-align:baseline}
audio:not([controls]){display:none;height:0}
[hidden],template{display:none}
a{background-color:transparent}
a:active,a:hover{outline:0}abbr[title]{border-bottom:1px dotted}
b,strong,optgroup{font-weight:bold}dfn{font-style:italic}
mark{background:#ff0;color:#000}small{font-size:80%}
sub,sup{font-size:75%;line-height:0;position:relative;vertical-align:baseline}
sup{top:-0.5em}sub{bottom:-0.25em}
img{border:0}svg:not(:root){overflow:hidden}figure{margin:1em 40px}
hr{-moz-box-sizing:content-box;box-sizing:content-box}
code,kbd,pre,samp{font-family:monospace,monospace;font-size:1em}
button,input,optgroup,select,textarea{color:inherit;font:inherit;margin:0}
button{overflow:visible}button,select{text-transform:none}
button,html input[type=button],input[type=reset],input[type=submit]{-webkit-appearance:button;cursor:pointer}
button[disabled],html input[disabled]{cursor:default}
button::-moz-focus-inner,input::-moz-focus-inner{border:0;padding:0}
input[type=checkbox],input[type=radio]{padding:0}
input[type=number]::-webkit-inner-spin-button,input[type=number]::-webkit-outer-spin-button{height:auto}
input[type=search]{box-sizing:content-box;-webkit-appearance:textfield;-moz-box-sizing:content-box;-webkit-box-sizing:content-box}
input[type=search]::-webkit-search-cancel-button,input[type=search]::-webkit-search-decoration{-webkit-appearance:none}
fieldset{border:1px solid #c0c0c0;margin:0 2px;padding:.35em .625em .75em}
legend{border:0;padding:0}pre,textarea{overflow:auto}
/*End extract from normalize.css*/
html,body{font-family:Verdana,sans-serif;font-size:15px;line-height:1.5}
html{overflow-x:hidden}
h1,h2,h3,h4,h5,h6,.w3-slim,.w3-wide{font-family:"Segoe UI",Arial,sans-serif}
h1{font-size:36px}h2{font-size:30px}h3{font-size:24px}h4{font-size:20px}h5{font-size:18px}h6{font-size:16px}
.w3-serif{font-family:"Times New Roman",Times,serif}
h1,h2,h3,h4,h5,h6{font-weight:400;margin:10px 0}
.w3-wide{letter-spacing:4px}
h1 a,h2 a,h3 a,h4 a,h5 a,h6 a{font-weight:inherit}
hr{height:0;border:0;border-top:1px solid #eee;margin:20px 0}
img{margin-bottom:-5px}
a{color:inherit}
table{border-collapse:collapse;border-spacing:0;width:100%;display:table}
table,th,td{border:none}
.w3-table-all{border:1px solid #ccc}
.w3-bordered tr,.w3-table-all tr{border-bottom:1px solid #ddd}
.w3-striped tbody tr:nth-child(even){background-color:#f1f1f1}
.w3-table-all tr:nth-child(odd){background-color:#fff}
.w3-table-all tr:nth-child(even){background-color:#f1f1f1}
.w3-hoverable tbody tr:hover,.w3-ul.w3-hoverable li:hover{background-color:#ccc}
.w3-centered tr th,.w3-centered tr td{text-align:center}
.w3-table td,.w3-table th,.w3-table-all td,.w3-table-all th{padding:6px 8px;display:table-cell;text-align:left;vertical-align:top}
.w3-table th:first-child,.w3-table td:first-child,.w3-table-all th:first-child,.w3-table-all td:first-child{padding-left:16px}
.w3-btn,.w3-btn-block{border:none;display:inline-block;outline:0;padding:6px 16px;vertical-align:middle;overflow:hidden;text-decoration:none!important;color:#fff;background-color:#000;text-align:center;cursor:pointer;white-space:nowrap}
.w3-btn.w3-disabled,.w3-btn-block.w3-disabled,.w3-btn-floating.w3-disabled,.w3-btn:disabled,.w3-btn-floating:disabled,.w3-btn-floating-large.w3-disabled,.w3-btn-floating-large:disabled{cursor:not-allowed;opacity:0.3}
.w3-btn.w3-disabled *,.w3-btn-block.w3-disabled,.w3-btn-floating.w3-disabled *,.w3-btn:disabled *,.w3-btn-floating:disabled *{pointer-events:none}
.w3-btn.w3-disabled:hover,.w3-btn-block.w3-disabled:hover,.w3-btn:disabled:hover,.w3-btn-floating.w3-disabled:hover,.w3-btn-floating:disabled:hover,
.w3-btn-floating-large.w3-disabled:hover,.w3-btn-floating-large:disabled:hover{box-shadow:none}
.w3-btn:hover,.w3-btn-block:hover,.w3-btn-floating:hover,.w3-btn-floating-large:hover{box-shadow:0 8px 16px 0 rgba(0,0,0,0.2),0 6px 20px 0 rgba(0,0,0,0.19)}
.w3-btn-block{width:100%}
.w3-btn,.w3-btn-floating,.w3-btn-floating-large,.w3-closenav,.w3-opennav{-webkit-touch-callout:none;-webkit-user-select:none;-khtml-user-select:none;-moz-user-select:none;-ms-user-select:none;user-select:none}   
.w3-btn-floating,.w3-btn-floating-large{display:inline-block;text-align:center;color:#fff;background-color:#000;position:relative;overflow:hidden;z-index:1;padding:0;border-radius:50%;cursor:pointer;font-size:24px}
.w3-btn-floating{width:40px;height:40px;line-height:40px}
.w3-btn-floating-large{width:56px;height:56px;line-height:56px}
.w3-btn-group .w3-btn{float:left}
.w3-btn.w3-ripple{position:relative}
.w3-ripple:after{content:"";background:#90EE90;display:block;position:absolute;padding-top:300%;padding-left:350%;margin-left:-20px!important;margin-top:-120%;opacity:0;transition:.8s}
.w3-ripple:active:after{padding:0;margin:0;opacity:1;transition:0s}
.w3-badge,.w3-tag,.w3-sign{background-color:#000;color:#fff;display:inline-block;padding-left:8px;padding-right:8px;font-weight:bold;text-align:center}
.w3-badge{border-radius:50%}
ul.w3-ul{list-style-type:none;padding:0;margin:0}
ul.w3-ul li{padding:6px 2px 6px 16px;border-bottom:1px solid #ddd}
ul.w3-ul li:last-child{border-bottom:none}
.w3-image,.w3-tooltip,.w3-display-container{position:relative}
img.w3-image,.w3-image img{max-width:100%;height:auto}
.w3-image .w3-title{position:absolute;bottom:8px;left:16px;color:#fff;font-size:20px}
.w3-fluid{max-width:100%;height:auto}
.w3-tooltip .w3-text{display:none}
.w3-tooltip:hover .w3-text{display:inline-block}
.w3-navbar{list-style-type:none;margin:0;padding:0;overflow:hidden}
.w3-navbar li{float:left}.w3-navbar li a{display:block;padding:8px 16px}.w3-navbar li a:hover{color:#000;background-color:#ccc}
.w3-navbar .w3-dropdown-hover,.w3-navbar .w3-dropdown-click{position:static}
.w3-navbar .w3-dropdown-hover:hover,.w3-navbar .w3-dropdown-hover:first-child,.w3-navbar .w3-dropdown-click:hover{background-color:#ccc;color:#000}
.w3-navbar a,.w3-topnav a,.w3-sidenav a,.w3-dropnav a,.w3-dropdown-content a,.w3-accordion-content a{text-decoration:none!important}
.w3-navbar .w3-opennav.w3-right{float:right!important}
.w3-topnav{padding:8px 8px}
.w3-topnav a{padding:0 8px;border-bottom:3px solid transparent;-webkit-transition:border-bottom .3s;transition:border-bottom .3s}
.w3-topnav a:hover{border-bottom:3px solid #fff}
.w3-topnav .w3-dropdown-hover a{border-bottom:0}
.w3-opennav,.w3-closenav{color:inherit}
.w3-opennav:hover,.w3-closenav:hover{cursor:pointer;opacity:0.8}
.w3-btn,.w3-btn-floating,.w3-btn-floating-large,.w3-btn-block,.w3-hover-shadow,.w3-hover-opacity,
.w3-navbar a,.w3-sidenav a,.w3-dropnav a,.w3-pagination li a,.w3-hoverable tbody tr,.w3-hoverable li,.w3-accordion-content a,.w3-dropdown-content a,.w3-dropdown-click:hover,.w3-dropdown-hover:hover,.w3-opennav,.w3-closenav,.w3-closebtn,
.w3-hover-amber,.w3-hover-aqua,.w3-hover-blue,.w3-hover-light-blue,.w3-hover-brown,.w3-hover-cyan,.w3-hover-blue-grey,.w3-hover-green,.w3-hover-light-green,.w3-hover-indigo,.w3-hover-khaki,.w3-hover-lime,.w3-hover-orange,.w3-hover-deep-orange,.w3-hover-pink,
.w3-hover-purple,.w3-hover-deep-purple,.w3-hover-red,.w3-hover-sand,.w3-hover-teal,.w3-hover-yellow,.w3-hover-white,.w3-hover-black,.w3-hover-grey,.w3-hover-light-grey,.w3-hover-dark-grey,.w3-hover-text-amber,.w3-hover-text-aqua,.w3-hover-text-blue,.w3-hover-text-light-blue,
.w3-hover-text-brown,.w3-hover-text-cyan,.w3-hover-text-blue-grey,.w3-hover-text-green,.w3-hover-text-light-green,.w3-hover-text-indigo,.w3-hover-text-khaki,.w3-hover-text-lime,.w3-hover-text-orange,.w3-hover-text-deep-orange,.w3-hover-text-pink,.w3-hover-text-purple,
.w3-hover-text-deep-purple,.w3-hover-text-red,.w3-hover-text-sand,.w3-hover-text-teal,.w3-hover-text-yellow,.w3-hover-text-white,.w3-hover-text-black,.w3-hover-text-grey,.w3-hover-text-light-grey,.w3-hover-text-dark-grey
{-webkit-transition:background-color .3s,color .15s,box-shadow .3s,opacity 0.3s;transition:background-color .3s,color .15s,box-shadow .3s,opacity 0.3s}
.w3-sidenav{height:100%;width:200px;background-color:#fff;position:fixed!important;z-index:1;overflow:auto}
.w3-sidenav a{padding:4px 2px 4px 16px}
.w3-sidenav a:hover{background-color:#ccc}
.w3-sidenav a,.w3-dropnav a{display:block}
.w3-sidenav .w3-dropdown-hover:hover,.w3-sidenav .w3-dropdown-hover:first-child,.w3-sidenav .w3-dropdown-click:hover{background-color:#ccc;color:#000}
.w3-sidenav .w3-dropdown-hover,.w3-sidenav .w3-dropdown-click {width:100%}.w3-sidenav .w3-dropdown-hover .w3-dropdown-content,.w3-sidenav .w3-dropdown-click .w3-dropdown-content{min-width:100%}
.w3-main,#main{transition:margin-left .4s}
.w3-dropnav{background-color:#fff}
.w3-dropnav a:hover{text-decoration:underline!important}
.w3-modal{z-index:3;display:none;padding-top:100px;position:fixed;left:0;top:0;width:100%;height:100%;overflow:auto;background-color:rgb(0,0,0);background-color:rgba(0,0,0,0.4)}
.w3-modal-content{margin:auto;background-color:#fff;position:relative;padding:0;outline:0;width:600px}.w3-closebtn{text-decoration:none;float:right;font-size:24px;font-weight:bold;color:inherit}
.w3-closebtn:hover,.w3-closebtn:focus{color:#000;text-decoration:none;cursor:pointer}
.w3-pagination{display:inline-block;padding:0;margin:0}
.w3-pagination li{display:inline}
.w3-pagination li a{text-decoration:none;color:#000;float:left;padding:8px 16px}
.w3-pagination li a:hover,.w3-pagination li a:focus{background-color:#ccc}
.w3-input-group,.w3-group{margin-top:24px;margin-bottom:24px}
.w3-input{padding:8px;display:block;border:none;border-bottom:1px solid #808080;width:100%}
.w3-label{color:#009688}
.w3-input:not(:valid)~.w3-validate{color:#f44336}
.w3-select{padding:4px 0;width:100%;color:#000;border:1px solid transparent;border-bottom:1px solid #009688}
.w3-select select:focus{color:#000;border:1px solid #009688}.w3-select option[disabled]{color:#009688}
.w3-dropdown-click,.w3-dropdown-hover{position:relative;display:inline-block;cursor:pointer}
.w3-dropdown-hover:hover .w3-dropdown-content{display:block;z-index:1}
.w3-dropdown-content{cursor:auto;color:#000;background-color:#fff;display:none;position:absolute;min-width:160px;margin:0;padding:0}
.w3-dropdown-content a{padding:6px 16px;display:block}
.w3-dropdown-content a:hover{background-color:#ccc}
.w3-accordion {width:100%;cursor:pointer}
.w3-accordion-content{cursor:auto;display:none;position:relative;width:100%;margin:0;padding:0}
.w3-accordion-content a{padding:6px 16px;display:block}
.w3-accordion-content a:hover{background-color:#ccc}
.w3-progress-container{width:100%;height:1.5em;position:relative;background-color:#f1f1f1}
.w3-progressbar{background-color:#757575;height:100%;position:absolute;line-height:inherit}
input[type=checkbox].w3-check,input[type=radio].w3-radio{width:24px;height:24px;position:relative;top:6px}
input[type=checkbox].w3-check:checked+.w3-validate,input[type=radio].w3-radio:checked+.w3-validate{color:#009688} 
input[type=checkbox].w3-check:disabled+.w3-validate,input[type=radio].w3-radio:disabled+.w3-validate{color:#aaa}
.w3-responsive{overflow-x:auto}
.w3-container:after,.w3-row:after,.w3-row-padding:after,.w3-topnav:after,.w3-clear:after,.w3-btn-group:before,.w3-btn-group:after
{content:"";display:table;clear:both}
.w3-col,.w3-half,.w3-third,.w3-twothird,.w3-threequarter,.w3-quarter{float:left;width:100%}
.w3-col.s1{width:8.33333%}
.w3-col.s2{width:16.66666%}
.w3-col.s3{width:24.99999%}
.w3-col.s4{width:33.33333%}
.w3-col.s5{width:41.66666%}
.w3-col.s6{width:49.99999%}
.w3-col.s7{width:58.33333%}
.w3-col.s8{width:66.66666%}
.w3-col.s9{width:74.99999%}
.w3-col.s10{width:83.33333%}
.w3-col.s11{width:91.66666%}
.w3-col.s12,.w3-half,.w3-third,.w3-twothird,.w3-threequarter,.w3-quarter{width:99.99999%}
@media only screen and (min-width:601px){
.w3-col.m1{width:8.33333%}
.w3-col.m2{width:16.66666%}
.w3-col.m3,.w3-quarter{width:24.99999%}
.w3-col.m4,.w3-third{width:33.33333%}
.w3-col.m5{width:41.66666%}
.w3-col.m6,.w3-half{width:49.99999%}
.w3-col.m7{width:58.33333%}
.w3-col.m8,.w3-twothird{width:66.66666%}
.w3-col.m9,.w3-threequarter{width:74.99999%}
.w3-col.m10{width:83.33333%}
.w3-col.m11{width:91.66666%}
.w3-col.m12{width:99.99999%}}
@media only screen and (min-width:993px){
.w3-col.l1{width:8.33333%}
.w3-col.l2{width:16.66666%}
.w3-col.l3,.w3-quarter{width:24.99999%}
.w3-col.l4,.w3-third{width:33.33333%}
.w3-col.l5{width:41.66666%}
.w3-col.l6,.w3-half{width:49.99999%}
.w3-col.l7{width:58.33333%}
.w3-col.l8,.w3-twothird{width:66.66666%}
.w3-col.l9,.w3-threequarter{width:74.99999%}
.w3-col.l10{width:83.33333%}
.w3-col.l11{width:91.66666%}
.w3-col.l12{width:99.99999%}}
.w3-content{max-width:980px;margin:auto}
.w3-rest{overflow:hidden}
.w3-hide{display:none!important}.w3-show-block,.w3-show{display:block!important}.w3-show-inline-block{display:inline-block!important}
@media (max-width:600px){.w3-modal-content{margin:0 10px;width:auto!important}.w3-modal{padding-top:30px}}
@media (max-width:768px){.w3-modal-content{width:500px}.w3-modal{padding-top:50px}}
@media (min-width:993px){.w3-modal-content{width:900px}}
@media screen and (max-width:600px){.w3-topnav a{display:block}.w3-navbar li:not(.w3-opennav){float:none;width:100%!important}.w3-navbar li.w3-right{float:none!important}}	
@media screen and (max-width:600px){.w3-topnav .w3-dropdown-hover .w3-dropdown-content,.w3-navbar .w3-dropdown-click .w3-dropdown-content,.w3-navbar .w3-dropdown-hover .w3-dropdown-content{position:relative}}	
@media screen and (max-width:600px){.w3-topnav,.w3-navbar{text-align:center}}
@media (max-width:600px){.w3-hide-small{display:none!important}}
@media (max-width:992px) and (min-width:601px){.w3-hide-medium{display:none!important}}
@media (min-width:993px){.w3-hide-large{display:none!important}}
@media screen and (max-width:992px){.w3-sidenav.w3-collapse{display:none}.w3-main{margin-left:0!important}}
@media screen and (min-width:992px){.w3-sidenav.w3-collapse{display:block!important}}
.w3-top,.w3-bottom{position:fixed;width:100%;z-index:1}.w3-top{top:0}.w3-bottom{bottom:0}
.w3-overlay{position:fixed;display:none;width:100%;height:100%;top:0;left:0;right:0;bottom:0;background-color:rgba(0,0,0,0.5);z-index:2}
.w3-left{float:left!important}.w3-right{float:right!important}
.w3-tiny{font-size:10px!important}.w3-small{font-size:12px!important}
.w3-medium{font-size:15px!important}
.w3-large{font-size:18px!important}
.w3-xlarge{font-size:24px!important}
.w3-xxlarge{font-size:36px!important}
.w3-xxxlarge{font-size:48px!important}
.w3-jumbo{font-size:64px!important}
.w3-vertical{word-break:break-all;line-height:1;text-align:center;width:0.6em}
.w3-left-align{text-align:left!important}.w3-right-align{text-align:right!important}
.w3-justify{text-align:justify!important}
.w3-center{text-align:center!important}
.w3-display-topleft{position:absolute;left:0;top:0}.w3-display-topright{position:absolute;right:0;top:0}
.w3-display-bottomleft{position:absolute;left:0;bottom:0}.w3-display-bottomright{position:absolute;right:0;bottom:0}
.w3-display-middle{position:absolute;left:0;bottom:50%;width:100%;text-align:center}
.w3-display-topmiddle{position:absolute;left:0;top:0;width:100%;text-align:center}.w3-display-bottommiddle{position:absolute;left:0;bottom:0;width:100%;text-align:center}
.w3-circle{border-radius:50%!important}
.w3-round-small{border-radius:2px!important}.w3-round,.w3-round-medium{border-radius:4px!important}
.w3-round-large{border-radius:8px!important}.w3-round-xlarge{border-radius:16px!important}
.w3-round-xxlarge{border-radius:32px!important}.w3-round-jumbo{border-radius:64px!important}
.w3-border-0{border:0!important}
.w3-border{border:1px solid #ccc!important}
.w3-border-top{border-top:1px solid #ccc!important}.w3-border-bottom{border-bottom:1px solid #ccc!important}
.w3-border-left{border-left:1px solid #ccc!important}.w3-border-right{border-right:1px solid #ccc!important}
.w3-margin{margin:16px!important}.w3-margin-0{margin:0!important}
.w3-margin-top{margin-top:16px!important}.w3-margin-bottom{margin-bottom:16px!important}
.w3-margin-left{margin-left:16px!important}.w3-margin-right{margin-right:16px!important}
.w3-section{margin-top:16px!important;margin-bottom:16px!important}
/* Might be removed in a later version */
.w3-margin-4{margin:4px!important}.w3-margin-8{margin:8px!important}
.w3-margin-12{margin:12px!important}.w3-margin-16{margin:16px!important}.w3-margin-24{margin:24px!important}
.w3-margin-32{margin:32px!important}.w3-margin-48{margin:48px!important}.w3-margin-64{margin:64px!important}
/* End remove */
.w3-padding-tiny{padding:2px 4px!important}
.w3-padding-small{padding:4px 8px!important}
.w3-padding-medium,.w3-padding,.w3-form{padding:8px 16px!important}
.w3-padding-large{padding:12px 24px!important}
.w3-padding-xlarge{padding:16px 32px!important}
.w3-padding-xxlarge{padding:24px 48px!important}
.w3-padding-jumbo{padding:32px 64px!important}
.w3-padding-4,.w3-padding-hor-4{padding-top:4px!important;padding-bottom:4px!important}
.w3-padding-8,.w3-padding-hor-8{padding-top:8px!important;padding-bottom:8px!important}
.w3-padding-12,.w3-padding-hor-12{padding-top:12px!important;padding-bottom:12px!important}
.w3-padding-16,.w3-padding-hor-16{padding-top:16px!important;padding-bottom:16px!important}
.w3-padding-24,.w3-padding-hor-24{padding-top:24px!important;padding-bottom:24px!important}
.w3-padding-32,.w3-padding-hor-32{padding-top:32px!important;padding-bottom:32px!important}
.w3-padding-48,.w3-padding-hor-48{padding-top:48px!important;padding-bottom:48px!important}
.w3-padding-64,.w3-padding-hor-64{padding-top:64px!important;padding-bottom:64px!important}
.w3-padding-128,.w3-padding-hor-128{padding-top:128px!important;padding-bottom:128px!important}
.w3-padding-0{padding:0!important}
.w3-padding-ver-4{padding-left:4px!important;padding-right:4px!important}
.w3-padding-ver-8{padding-left:8px!important;padding-right:8px!important}
.w3-padding-ver-12{padding-left:12px!important;padding-right:12px!important}
.w3-padding-ver-16{padding-left:16px!important;padding-right:16px!important}
.w3-padding-ver-24{padding-left:24px!important;padding-right:24px!important}
.w3-padding-ver-32{padding-left:32px!important;padding-right:32px!important}
.w3-padding-ver-48{padding-left:48px!important;padding-right:48px!important}
.w3-padding-ver-64{padding-left:64px!important;padding-right:64px!important}
.w3-padding-top{padding-top:8px!important}.w3-padding-bottom{padding-bottom:8px!important}
.w3-padding-left{padding-left:16px!important}.w3-padding-right{padding-right:16px!important}
.w3-topbar{border-top:6px solid #ccc!important}.w3-bottombar{border-bottom:6px solid #ccc!important}
.w3-leftbar{border-left:6px solid #ccc!important}.w3-rightbar{border-right:6px solid #ccc!important}
.w3-row-padding,.w3-row-padding>.w3-half,.w3-row-padding>.w3-third,.w3-row-padding>.w3-twothird,.w3-row-padding>.w3-threequarter,.w3-row-padding>.w3-quarter,.w3-row-padding>.w3-col{padding:0 8px}
.w3-spin{animation:w3-spin 2s infinite linear;-webkit-animation:w3-spin 2s infinite linear}
@-webkit-keyframes w3-spin{
0%{-webkit-transform:rotate(0deg);transform:rotate(0deg)}
100%{-webkit-transform:rotate(359deg);transform:rotate(359deg)}}
@keyframes w3-spin{
0%{-webkit-transform:rotate(0deg);transform: rotate(0deg)}
100%{-webkit-transform:rotate(359deg);transform:rotate(359deg)}}
.w3-container{padding:0.01em 16px}
.w3-example{background-color:#f1f1f1;padding:0.01em 16px}
.w3-code{font-family:Consolas,"courier new";font-size:16px;line-height:1.4;width:auto;background-color:#fff;padding:8px 12px;border-left:4px solid #009688;word-wrap:break-word}
.w3-example,.w3-code,.w3-reference{margin:20px 0}
.w3-card{border:1px solid #ccc}
.w3-card-2,.w3-example{box-shadow:0 2px 4px 0 rgba(0,0,0,0.16),0 2px 10px 0 rgba(0,0,0,0.12)!important}
.w3-card-4,.w3-hover-shadow:hover{box-shadow:0 4px 8px 0 rgba(0,0,0,0.2),0 6px 20px 0 rgba(0,0,0,0.19)!important}
.w3-card-8{box-shadow:0 8px 16px 0 rgba(0,0,0,0.2),0 6px 20px 0 rgba(0,0,0,0.19)!important}
.w3-card-12{box-shadow:0 12px 16px 0 rgba(0,0,0,0.24),0 17px 50px 0 rgba(0,0,0,0.19)!important}
.w3-card-16{box-shadow:0 16px 24px 0 rgba(0,0,0,0.22),0 25px 55px 0 rgba(0,0,0,0.21)!important}
.w3-card-24{box-shadow:0 24px 24px 0 rgba(0,0,0,0.2),0 40px 77px 0 rgba(0,0,0,0.22)!important}
.w3-animate-fading{-webkit-animation:fading 10s infinite;animation:fading 10s infinite}
@-webkit-keyframes fading{0%{opacity:0}50%{opacity:1}100%{opacity:0}}
@keyframes fading{0%{opacity:0}50%{opacity:1}100%{opacity:0}}
.w3-animate-opacity{-webkit-animation:opac 1.5s;animation:opac 1.5s}
@-webkit-keyframes opac{from{opacity:0} to{opacity:1}}
@keyframes opac{from{opacity:0} to{opacity:1}}
.w3-animate-top{position:relative;-webkit-animation:animatetop 0.4s;animation:animatetop 0.4s}
@-webkit-keyframes animatetop{from{top:-300px;opacity:0} to{top:0;opacity:1}}
@keyframes animatetop{from{top:-300px;opacity:0} to{top:0;opacity:1}}
.w3-animate-left{position:relative;-webkit-animation:animateleft 0.4s;animation:animateleft 0.4s}
@-webkit-keyframes animateleft{from{left:-300px;opacity:0} to{left:0;opacity:1}}
@keyframes animateleft{from{left:-300px;opacity:0} to{left:0;opacity:1}}
.w3-animate-right{position:relative;-webkit-animation:animateright 0.4s;animation:animateright 0.4s}
@-webkit-keyframes animateright{from{right:-300px;opacity:0} to{right:0;opacity:1}}
@keyframes animateright{from{right:-300px;opacity:0} to{right:0;opacity:1}}
.w3-animate-bottom{position:relative;-webkit-animation:animatebottom 0.4s;animation:animatebottom 0.4s}
@-webkit-keyframes animatebottom{from{bottom:-300px;opacity:0} to{bottom:0px;opacity:1}}
@keyframes animatebottom{from{bottom:-300px;opacity:0} to{bottom:0;opacity:1}}
.w3-animate-zoom {-webkit-animation:animatezoom 0.6s;animation:animatezoom 0.6s}
@-webkit-keyframes animatezoom{from{-webkit-transform:scale(0)} to{-webkit-transform:scale(1)}}
@keyframes animatezoom{from{transform:scale(0)} to{transform:scale(1)}}
.w3-animate-input{-webkit-transition:width 0.4s ease-in-out;transition:width 0.4s ease-in-out}.w3-animate-input:focus{width:100%!important}
.w3-transparent{background-color:transparent!important}
.w3-hover-none:hover{box-shadow:none!important;background-color:transparent!important}
.w3-amber,.w3-hover-amber:hover{color:#000!important;background-color:#ffc107!important}
.w3-aqua,.w3-hover-aqua:hover{color:#000!important;background-color:#00ffff!important}
.w3-blue,.w3-hover-blue:hover{color:#fff!important;background-color:#2196F3!important}
.w3-light-blue,.w3-hover-light-blue:hover{color:#000!important;background-color:#87CEEB!important}
.w3-brown,.w3-hover-brown:hover{color:#fff!important;background-color:#795548!important}
.w3-cyan,.w3-hover-cyan:hover{color:#000!important;background-color:#00bcd4!important}
.w3-blue-grey,.w3-hover-blue-grey:hover{color:#fff!important;background-color:#607d8b!important}
.w3-green,.w3-hover-green:hover{color:#fff!important;background-color:#4CAF50!important}
.w3-light-green,.w3-hover-light-green:hover{color:#000!important;background-color:#8bc34a!important}
.w3-indigo,.w3-hover-indigo:hover{color:#fff!important;background-color:#3f51b5!important}
.w3-khaki,.w3-hover-khaki:hover{color:#000!important;background-color:#f0e68c!important}
.w3-lime,.w3-hover-lime:hover{color:#000!important;background-color:#cddc39!important}
.w3-orange,.w3-hover-orange:hover{color:#000!important;background-color:#ff9800!important}
.w3-deep-orange,.w3-hover-deep-orange:hover{color:#fff!important;background-color:#ff5722!important}
.w3-pink,.w3-hover-pink:hover{color:#fff!important;background-color:#e91e63!important}
.w3-purple,.w3-hover-purple:hover{color:#fff!important;background-color:#9c27b0!important}
.w3-deep-purple,.w3-hover-deep-purple:hover{color:#fff!important;background-color:#673ab7!important}
.w3-red,.w3-hover-red:hover{color:#fff!important;background-color:#f44336!important}
.w3-sand,.w3-hover-sand:hover{color:#000!important;background-color:#fdf5e6!important}
.w3-teal,.w3-hover-teal:hover{color:#fff!important;background-color:#009688!important}
.w3-yellow,.w3-hover-yellow:hover{color:#000!important;background-color:#ffeb3b!important}
.w3-white,.w3-hover-white:hover{color:#000!important;background-color:#fff!important}
.w3-black,.w3-hover-black:hover{color:#fff!important;background-color:#000!important}
.w3-grey,.w3-hover-grey:hover{color:#000!important;background-color:#9e9e9e!important}
.w3-light-grey,.w3-hover-light-grey:hover{color:#000!important;background-color:#f1f1f1!important}
.w3-dark-grey,.w3-hover-dark-grey:hover{color:#fff!important;background-color:#616161!important}
.w3-pale-red,.w3-hover-pale-red:hover{color:#000!important;background-color:#ffe7e7!important}.w3-pale-green,.w3-hover-pale-green:hover{color:#000!important;background-color:#e7ffe7!important}
.w3-pale-yellow,.w3-hover-pale-yellow:hover{color:#000!important;background-color:#ffffd7!important}.w3-pale-blue,.w3-hover-pale-blue:hover{color:#000!important;background-color:#e7ffff!important}
.w3-text-amber,.w3-hover-text-amber:hover{color:#ffc107!important}
.w3-text-aqua,.w3-hover-text-aqua:hover{color:#00ffff!important}
.w3-text-blue,.w3-hover-text-blue:hover{color:#2196F3!important}
.w3-text-light-blue,.w3-hover-text-light-blue:hover{color:#87CEEB!important}
.w3-text-brown,.w3-hover-text-brown:hover{color:#795548!important}
.w3-text-cyan,.w3-hover-text-cyan:hover{color:#00bcd4!important}
.w3-text-blue-grey,.w3-hover-text-blue-grey:hover{color:#607d8b!important}
.w3-text-green,.w3-hover-text-green:hover{color:#4CAF50!important}
.w3-text-light-green,.w3-hover-text-light-green:hover{color:#8bc34a!important}
.w3-text-indigo,.w3-hover-text-indigo:hover{color:#3f51b5!important}
.w3-text-khaki,.w3-hover-text-khaki:hover{color:#b4aa50!important}
.w3-text-lime,.w3-hover-text-lime:hover{color:#cddc39!important}
.w3-text-orange,.w3-hover-text-orange:hover{color:#ff9800!important}
.w3-text-deep-orange,.w3-hover-text-deep-orange:hover{color:#ff5722!important}
.w3-text-pink,.w3-hover-text-pink:hover{color:#e91e63!important}
.w3-text-purple,.w3-hover-text-purple:hover{color:#9c27b0!important}
.w3-text-deep-purple,.w3-hover-text-deep-purple:hover{color:#673ab7!important}
.w3-text-red,.w3-hover-text-red:hover{color:#f44336!important}
.w3-text-sand,.w3-hover-text-sand:hover{color:#fdf5e6!important}
.w3-text-teal,.w3-hover-text-teal:hover{color:#009688!important}
.w3-text-yellow,.w3-hover-text-yellow:hover{color:#d2be0e!important}
.w3-text-white,.w3-hover-text-white:hover{color:#fff!important}
.w3-text-black,.w3-hover-text-black:hover{color:#000!important}
.w3-text-grey,.w3-hover-text-grey:hover{color:#757575!important}
.w3-text-light-grey,.w3-hover-text-light-grey:hover{color:#f1f1f1!important}
.w3-text-dark-grey,.w3-hover-text-dark-grey:hover{color:#3a3a3a!important}
.w3-border-amber,.w3-hover-border-amber:hover{border-color:#ffc107!important}
.w3-border-aqua,.w3-hover-border-aqua:hover{border-color:#00ffff!important}
.w3-border-blue,.w3-hover-border-blue:hover{border-color:#2196F3!important}
.w3-border-light-blue,.w3-hover-border-light-blue:hover{border-color:#87CEEB!important}
.w3-border-brown,.w3-hover-border-brown:hover{border-color:#795548!important}
.w3-border-cyan,.w3-hover-border-cyan:hover{border-color:#00bcd4!important}
.w3-border-blue-grey,.w3-hover-blue-grey:hover{border-color:#607d8b!important}
.w3-border-green,.w3-hover-border-green:hover{border-color:#4CAF50!important}
.w3-border-light-green,.w3-hover-border-light-green:hover{border-color:#8bc34a!important}
.w3-border-indigo,.w3-hover-border-indigo:hover{border-color:#3f51b5!important}
.w3-border-khaki,.w3-hover-border-khaki:hover{border-color:#f0e68c!important}
.w3-border-lime,.w3-hover-border-lime:hover{border-color:#cddc39!important}
.w3-border-orange,.w3-hover-border-orange:hover{border-color:#ff9800!important}
.w3-border-deep-orange,.w3-hover-border-deep-orange:hover{border-color:#ff5722!important}
.w3-border-pink,.w3-hover-border-pink:hover{border-color:#e91e63!important}
.w3-border-purple,.w3-hover-border-purple:hover{border-color:#9c27b0!important}
.w3-border-deep-purple,.w3-hover-border-deep-purple:hover{border-color:#673ab7!important}
.w3-border-red,.w3-hover-border-red:hover{border-color:#f44336!important}
.w3-border-sand,.w3-hover-border-sand:hover{border-color:#fdf5e6!important}
.w3-border-teal,.w3-hover-border-teal:hover{border-color:#009688!important}
.w3-border-yellow,.w3-hover-border-yellow:hover{border-color:#ffeb3b!important}
.w3-border-white,.w3-hover-border-white:hover{border-color:#fff!important}
.w3-border-black,.w3-hover-border-black:hover{border-color:#000!important}
.w3-border-grey,.w3-hover-border-grey:hover{border-color:#9e9e9e!important}
.w3-border-light-grey,.w3-hover-border-light-grey:hover{border-color:#f1f1f1!important}
.w3-border-dark-grey,.w3-hover-border-dark-grey:hover{border-color:#616161!important}
.w3-border-pale-red,.w3-hover-border-pale-red:hover{border-color:#ffe7e7!important}.w3-border-pale-green,.w3-hover-border-pale-green:hover{border-color:#e7ffe7!important}
.w3-border-pale-yellow,.w3-hover-border-pale-yellow:hover{border-color:#ffffd7!important}.w3-border-pale-blue,.w3-hover-border-pale-blue:hover{border-color:#e7ffff!important}
.w3-opacity,.w3-hover-opacity:hover{opacity:0.60}
.w3-text-shadow{text-shadow:1px 1px 0 #444}.w3-text-shadow-white{text-shadow:1px 1px 0 #ddd}

.w3-theme-l5 {color:#000 !important; background-color:#f8fbf4 !important}
.w3-theme-l4 {color:#000 !important; background-color:#e8f3db !important}
.w3-theme-l3 {color:#000 !important; background-color:#d1e7b7 !important}
.w3-theme-l2 {color:#000 !important; background-color:#b9db93 !important}
.w3-theme-l1 {color:#000 !important; background-color:#a2cf6f !important}
.w3-theme-d1 {color:#fff !important; background-color:#7eb63d !important}
.w3-theme-d2 {color:#fff !important; background-color:#70a236 !important}
.w3-theme-d3 {color:#fff !important; background-color:#628e2f !important}
.w3-theme-d4 {color:#fff !important; background-color:#547a29 !important}
.w3-theme-d5 {color:#fff !important; background-color:#466522 !important}

.w3-theme-light {color:#000 !important; background-color:#f8fbf4 !important}
.w3-theme-dark {color:#fff !important; background-color:#466522 !important}
.w3-theme-action {color:#fff !important; background-color:#466522 !important}

.w3-theme {color:#fff !important; background-color:#8bc34a !important}
.w3-text-theme {color:#8bc34a !important}
.w3-theme-border {border-color:#8bc34a !important}
.w3-hover-theme:hover {color:#fff !important; background-color:#8bc34a !important}

/*

Visual Studio-like style based on original C# coloring by Jason Diamond <jason@diamond.name>

*/
.hljs {
  display: block;
  overflow-x: auto;
  padding: 0.5em;
  background: white;
  color: black;
}

.hljs-comment,
.hljs-quote,
.hljs-variable {
  color: #008000;
}

.hljs-keyword,
.hljs-selector-tag,
.hljs-built_in,
.hljs-name,
.hljs-tag {
  color: #00f;
}

.hljs-string,
.hljs-title,
.hljs-section,
.hljs-attribute,
.hljs-literal,
.hljs-template-tag,
.hljs-template-variable,
.hljs-type,
.hljs-addition {
  color: #a31515;
}

.hljs-deletion,
.hljs-selector-attr,
.hljs-selector-pseudo,
.hljs-meta {
  color: #2b91af;
}

.hljs-doctag {
  color: #808080;
}

.hljs-attr {
  color: #f00;
}

.hljs-symbol,
.hljs-bullet,
.hljs-link {
  color: #00b0e8;
}


.hljs-emphasis {
  font-style: italic;
}

.hljs-strong {
  font-weight: bold;
}
</style>
<script>
/*! highlight.js v9.3.0 | BSD3 License | git.io/hljslicense */
!function(e){var n="object"==typeof window&&window||"object"==typeof self&&self;"undefined"!=typeof exports?e(exports):n&&(n.hljs=e({}),"function"==typeof define&&define.amd&&define([],function(){return n.hljs}))}(function(e){function n(e){return e.replace(/&/gm,"&amp;").replace(/</gm,"&lt;").replace(/>/gm,"&gt;")}function t(e){return e.nodeName.toLowerCase()}function r(e,n){var t=e&&e.exec(n);return t&&0==t.index}function a(e){return/^(no-?highlight|plain|text)$/i.test(e)}function i(e){var n,t,r,i=e.className+" ";if(i+=e.parentNode?e.parentNode.className:"",t=/\blang(?:uage)?-([\w-]+)\b/i.exec(i))return w(t[1])?t[1]:"no-highlight";for(i=i.split(/\s+/),n=0,r=i.length;r>n;n++)if(w(i[n])||a(i[n]))return i[n]}function o(e,n){var t,r={};for(t in e)r[t]=e[t];if(n)for(t in n)r[t]=n[t];return r}function u(e){var n=[];return function r(e,a){for(var i=e.firstChild;i;i=i.nextSibling)3==i.nodeType?a+=i.nodeValue.length:1==i.nodeType&&(n.push({event:"start",offset:a,node:i}),a=r(i,a),t(i).match(/br|hr|img|input/)||n.push({event:"stop",offset:a,node:i}));return a}(e,0),n}function c(e,r,a){function i(){return e.length&&r.length?e[0].offset!=r[0].offset?e[0].offset<r[0].offset?e:r:"start"==r[0].event?e:r:e.length?e:r}function o(e){function r(e){return" "+e.nodeName+'="'+n(e.value)+'"'}f+="<"+t(e)+Array.prototype.map.call(e.attributes,r).join("")+">"}function u(e){f+="</"+t(e)+">"}function c(e){("start"==e.event?o:u)(e.node)}for(var s=0,f="",l=[];e.length||r.length;){var g=i();if(f+=n(a.substr(s,g[0].offset-s)),s=g[0].offset,g==e){l.reverse().forEach(u);do c(g.splice(0,1)[0]),g=i();while(g==e&&g.length&&g[0].offset==s);l.reverse().forEach(o)}else"start"==g[0].event?l.push(g[0].node):l.pop(),c(g.splice(0,1)[0])}return f+n(a.substr(s))}function s(e){function n(e){return e&&e.source||e}function t(t,r){return new RegExp(n(t),"m"+(e.cI?"i":"")+(r?"g":""))}function r(a,i){if(!a.compiled){if(a.compiled=!0,a.k=a.k||a.bK,a.k){var u={},c=function(n,t){e.cI&&(t=t.toLowerCase()),t.split(" ").forEach(function(e){var t=e.split("|");u[t[0]]=[n,t[1]?Number(t[1]):1]})};"string"==typeof a.k?c("keyword",a.k):Object.keys(a.k).forEach(function(e){c(e,a.k[e])}),a.k=u}a.lR=t(a.l||/\w+/,!0),i&&(a.bK&&(a.b="\\b("+a.bK.split(" ").join("|")+")\\b"),a.b||(a.b=/\B|\b/),a.bR=t(a.b),a.e||a.eW||(a.e=/\B|\b/),a.e&&(a.eR=t(a.e)),a.tE=n(a.e)||"",a.eW&&i.tE&&(a.tE+=(a.e?"|":"")+i.tE)),a.i&&(a.iR=t(a.i)),void 0===a.r&&(a.r=1),a.c||(a.c=[]);var s=[];a.c.forEach(function(e){e.v?e.v.forEach(function(n){s.push(o(e,n))}):s.push("self"==e?a:e)}),a.c=s,a.c.forEach(function(e){r(e,a)}),a.starts&&r(a.starts,i);var f=a.c.map(function(e){return e.bK?"\\.?("+e.b+")\\.?":e.b}).concat([a.tE,a.i]).map(n).filter(Boolean);a.t=f.length?t(f.join("|"),!0):{exec:function(){return null}}}}r(e)}function f(e,t,a,i){function o(e,n){for(var t=0;t<n.c.length;t++)if(r(n.c[t].bR,e))return n.c[t]}function u(e,n){if(r(e.eR,n)){for(;e.endsParent&&e.parent;)e=e.parent;return e}return e.eW?u(e.parent,n):void 0}function c(e,n){return!a&&r(n.iR,e)}function g(e,n){var t=N.cI?n[0].toLowerCase():n[0];return e.k.hasOwnProperty(t)&&e.k[t]}function p(e,n,t,r){var a=r?"":E.classPrefix,i='<span class="'+a,o=t?"":"</span>";return i+=e+'">',i+n+o}function h(){if(!k.k)return n(M);var e="",t=0;k.lR.lastIndex=0;for(var r=k.lR.exec(M);r;){e+=n(M.substr(t,r.index-t));var a=g(k,r);a?(B+=a[1],e+=p(a[0],n(r[0]))):e+=n(r[0]),t=k.lR.lastIndex,r=k.lR.exec(M)}return e+n(M.substr(t))}function d(){var e="string"==typeof k.sL;if(e&&!R[k.sL])return n(M);var t=e?f(k.sL,M,!0,y[k.sL]):l(M,k.sL.length?k.sL:void 0);return k.r>0&&(B+=t.r),e&&(y[k.sL]=t.top),p(t.language,t.value,!1,!0)}function b(){L+=void 0!==k.sL?d():h(),M=""}function v(e,n){L+=e.cN?p(e.cN,"",!0):"",k=Object.create(e,{parent:{value:k}})}function m(e,n){if(M+=e,void 0===n)return b(),0;var t=o(n,k);if(t)return t.skip?M+=n:(t.eB&&(M+=n),b(),t.rB||t.eB||(M=n)),v(t,n),t.rB?0:n.length;var r=u(k,n);if(r){var a=k;a.skip?M+=n:(a.rE||a.eE||(M+=n),b(),a.eE&&(M=n));do k.cN&&(L+="</span>"),k.skip||(B+=k.r),k=k.parent;while(k!=r.parent);return r.starts&&v(r.starts,""),a.rE?0:n.length}if(c(n,k))throw new Error('Illegal lexeme "'+n+'" for mode "'+(k.cN||"<unnamed>")+'"');return M+=n,n.length||1}var N=w(e);if(!N)throw new Error('Unknown language: "'+e+'"');s(N);var x,k=i||N,y={},L="";for(x=k;x!=N;x=x.parent)x.cN&&(L=p(x.cN,"",!0)+L);var M="",B=0;try{for(var C,j,I=0;;){if(k.t.lastIndex=I,C=k.t.exec(t),!C)break;j=m(t.substr(I,C.index-I),C[0]),I=C.index+j}for(m(t.substr(I)),x=k;x.parent;x=x.parent)x.cN&&(L+="</span>");return{r:B,value:L,language:e,top:k}}catch(O){if(-1!=O.message.indexOf("Illegal"))return{r:0,value:n(t)};throw O}}function l(e,t){t=t||E.languages||Object.keys(R);var r={r:0,value:n(e)},a=r;return t.filter(w).forEach(function(n){var t=f(n,e,!1);t.language=n,t.r>a.r&&(a=t),t.r>r.r&&(a=r,r=t)}),a.language&&(r.second_best=a),r}function g(e){return E.tabReplace&&(e=e.replace(/^((<[^>]+>|\t)+)/gm,function(e,n){return n.replace(/\t/g,E.tabReplace)})),E.useBR&&(e=e.replace(/\n/g,"<br>")),e}function p(e,n,t){var r=n?x[n]:t,a=[e.trim()];return e.match(/\bhljs\b/)||a.push("hljs"),-1===e.indexOf(r)&&a.push(r),a.join(" ").trim()}function h(e){var n=i(e);if(!a(n)){var t;E.useBR?(t=document.createElementNS("http://www.w3.org/1999/xhtml","div"),t.innerHTML=e.innerHTML.replace(/\n/g,"").replace(/<br[ \/]*>/g,"\n")):t=e;var r=t.textContent,o=n?f(n,r,!0):l(r),s=u(t);if(s.length){var h=document.createElementNS("http://www.w3.org/1999/xhtml","div");h.innerHTML=o.value,o.value=c(s,u(h),r)}o.value=g(o.value),e.innerHTML=o.value,e.className=p(e.className,n,o.language),e.result={language:o.language,re:o.r},o.second_best&&(e.second_best={language:o.second_best.language,re:o.second_best.r})}}function d(e){E=o(E,e)}function b(){if(!b.called){b.called=!0;var e=document.querySelectorAll("pre code");Array.prototype.forEach.call(e,h)}}function v(){addEventListener("DOMContentLoaded",b,!1),addEventListener("load",b,!1)}function m(n,t){var r=R[n]=t(e);r.aliases&&r.aliases.forEach(function(e){x[e]=n})}function N(){return Object.keys(R)}function w(e){return e=(e||"").toLowerCase(),R[e]||R[x[e]]}var E={classPrefix:"hljs-",tabReplace:null,useBR:!1,languages:void 0},R={},x={};return e.highlight=f,e.highlightAuto=l,e.fixMarkup=g,e.highlightBlock=h,e.configure=d,e.initHighlighting=b,e.initHighlightingOnLoad=v,e.registerLanguage=m,e.listLanguages=N,e.getLanguage=w,e.inherit=o,e.IR="[a-zA-Z]\\w*",e.UIR="[a-zA-Z_]\\w*",e.NR="\\b\\d+(\\.\\d+)?",e.CNR="(-?)(\\b0[xX][a-fA-F0-9]+|(\\b\\d+(\\.\\d*)?|\\.\\d+)([eE][-+]?\\d+)?)",e.BNR="\\b(0b[01]+)",e.RSR="!|!=|!==|%|%=|&|&&|&=|\\*|\\*=|\\+|\\+=|,|-|-=|/=|/|:|;|<<|<<=|<=|<|===|==|=|>>>=|>>=|>=|>>>|>>|>|\\?|\\[|\\{|\\(|\\^|\\^=|\\||\\|=|\\|\\||~",e.BE={b:"\\\\[\\s\\S]",r:0},e.ASM={cN:"string",b:"'",e:"'",i:"\\n",c:[e.BE]},e.QSM={cN:"string",b:'"',e:'"',i:"\\n",c:[e.BE]},e.PWM={b:/\b(a|an|the|are|I'm|isn't|don't|doesn't|won't|but|just|should|pretty|simply|enough|gonna|going|wtf|so|such|will|you|your|like)\b/},e.C=function(n,t,r){var a=e.inherit({cN:"comment",b:n,e:t,c:[]},r||{});return a.c.push(e.PWM),a.c.push({cN:"doctag",b:"(?:TODO|FIXME|NOTE|BUG|XXX):",r:0}),a},e.CLCM=e.C("//","$"),e.CBCM=e.C("/\\*","\\*/"),e.HCM=e.C("#","$"),e.NM={cN:"number",b:e.NR,r:0},e.CNM={cN:"number",b:e.CNR,r:0},e.BNM={cN:"number",b:e.BNR,r:0},e.CSSNM={cN:"number",b:e.NR+"(%|em|ex|ch|rem|vw|vh|vmin|vmax|cm|mm|in|pt|pc|px|deg|grad|rad|turn|s|ms|Hz|kHz|dpi|dpcm|dppx)?",r:0},e.RM={cN:"regexp",b:/\//,e:/\/[gimuy]*/,i:/\n/,c:[e.BE,{b:/\[/,e:/\]/,r:0,c:[e.BE]}]},e.TM={cN:"title",b:e.IR,r:0},e.UTM={cN:"title",b:e.UIR,r:0},e.METHOD_GUARD={b:"\\.\\s*"+e.UIR,r:0},e});hljs.registerLanguage("cpp",function(t){var e={cN:"keyword",b:"\\b[a-z\\d_]*_t\\b"},r={cN:"string",v:[t.inherit(t.QSM,{b:'((u8?|U)|L)?"'}),{b:'(u8?|U)?R"',e:'"',c:[t.BE]},{b:"'\\\\?.",e:"'",i:"."}]},i={cN:"number",v:[{b:"\\b(\\d+(\\.\\d*)?|\\.\\d+)(u|U|l|L|ul|UL|f|F)"},{b:t.CNR}],r:0},s={cN:"meta",b:"#",e:"$",k:{"meta-keyword":"if else elif endif define undef warning error line pragma ifdef ifndef"},c:[{b:/\\\n/,r:0},{bK:"include",e:"$",k:{"meta-keyword":"include"},c:[t.inherit(r,{cN:"meta-string"}),{cN:"meta-string",b:"<",e:">",i:"\\n"}]},r,t.CLCM,t.CBCM]},a=t.IR+"\\s*\\(",c={keyword:"int float while private char catch export virtual operator sizeof dynamic_cast|10 typedef const_cast|10 const struct for static_cast|10 union namespace unsigned long volatile static protected bool template mutable if public friend do goto auto void enum else break extern using class asm case typeid short reinterpret_cast|10 default double register explicit signed typename try this switch continue inline delete alignof constexpr decltype noexcept static_assert thread_local restrict _Bool complex _Complex _Imaginary atomic_bool atomic_char atomic_schar atomic_uchar atomic_short atomic_ushort atomic_int atomic_uint atomic_long atomic_ulong atomic_llong atomic_ullong",built_in:"std string cin cout cerr clog stdin stdout stderr stringstream istringstream ostringstream auto_ptr deque list queue stack vector map set bitset multiset multimap unordered_set unordered_map unordered_multiset unordered_multimap array shared_ptr abort abs acos asin atan2 atan calloc ceil cosh cos exit exp fabs floor fmod fprintf fputs free frexp fscanf isalnum isalpha iscntrl isdigit isgraph islower isprint ispunct isspace isupper isxdigit tolower toupper labs ldexp log10 log malloc realloc memchr memcmp memcpy memset modf pow printf putchar puts scanf sinh sin snprintf sprintf sqrt sscanf strcat strchr strcmp strcpy strcspn strlen strncat strncmp strncpy strpbrk strrchr strspn strstr tanh tan vfprintf vprintf vsprintf endl initializer_list unique_ptr",literal:"true false nullptr NULL"},n=[e,t.CLCM,t.CBCM,i,r];return{aliases:["c","cc","h","c++","h++","hpp"],k:c,i:"</",c:n.concat([s,{b:"\\b(deque|list|queue|stack|vector|map|set|bitset|multiset|multimap|unordered_map|unordered_set|unordered_multiset|unordered_multimap|array)\\s*<",e:">",k:c,c:["self",e]},{b:t.IR+"::",k:c},{v:[{b:/=/,e:/;/},{b:/\(/,e:/\)/},{bK:"new throw return else",e:/;/}],k:c,c:n.concat([{b:/\(/,e:/\)/,c:n.concat(["self"]),r:0}]),r:0},{cN:"function",b:"("+t.IR+"[\\*&\\s]+)+"+a,rB:!0,e:/[{;=]/,eE:!0,k:c,i:/[^\w\s\*&]/,c:[{b:a,rB:!0,c:[t.TM],r:0},{cN:"params",b:/\(/,e:/\)/,k:c,r:0,c:[t.CLCM,t.CBCM,r,i]},t.CLCM,t.CBCM,s]}])}});
</script>
<script>hljs.initHighlightingOnLoad();</script>

<body>
<div class="w3-card">
<div class="w3-container w3-theme w3-card-2">
  <h1>{{.ReportName}}</h1>
</div>
</div>
<ul class="w3-ul w3-border">
{{range $index, $value := .ListOfCheckers}}
  <li>{{$index}} <span class="w3-badge w3-green w3-right">{{$value}}</span></li>
{{end}}
</ul>

<h2>List of files</h2>


<div class="w3-accordion">
{{$parent:=.}}
{{$LineVar:=index .Consts 3}}
{{$NormalVar:=index .Consts 0}}
{{$WarningVar:=index .Consts 1}}
{{$DangerVar:=index .Consts 2}}
{{$ErrVar:=index .Consts 4}}
{{$ErrContVar:=index .Consts 5}}
{{range $index, $value := .ListOfFiles}}
{{$cont := index $parent.ListOfLong $value}}
  <button onclick="myFunction('repo{{$index}}')" class="w3-btn-block w3-left-align w3-khaki w3-border-black w3-border-bottom">{{$value}} <span class="w3-badge w3-green">{{$cont.Normal}}</span><span class="w3-badge w3-red">{{$cont.Critical}}</span><span class="w3-badge w3-yellow">{{$cont.Warning}}</span></button>
  <div id="repo{{$index}}" class="w3-accordion-content w3-container w3-small">  
    <div>
      <table class="w3-table tbl_code">
{{range $iind, $value_list := $cont.List}}
	{{if eq $value_list.Value $LineVar}}
       <tr>
        <td style="width:80px" class="w3-padding-0 w3-padding-left">{{$value_list.Number}}</td>
        <td class="w3-border-green w3-border-left w3-padding-0  w3-padding-left"><pre class="w3-padding-0 w3-margin-0 code">{{$value_list.Line}}</pre></td>
       </tr>
    {{else if eq $value_list.Value $NormalVar}}
	   <tr>
        <td colspan="2" class="w3-padding-0">
		 <p class="w3-container w3-light-blue w3-margin-0 w3-padding-0  w3-padding-left">{{$value_list.Plugin}}: {{$value_list.Line}}</p>
	{{if $value_list.Link}}
	     <p class="w3-container w3-light-blue w3-margin-0 w3-padding-0  w3-padding-left">
	      <div class="w3-dropdown-hover">
		    <button class="w3-btn w3-red w3-tiny w3-padding-0 w3-margin-0">Link</button>
		    <div class="w3-dropdown-content w3-border w3-light-grey">
		{{range $iiind, $value_list_inner := $value_list.Link}}
			{{if eq $value_list_inner.FndStr false}}
			  <pre class="w3-padding-0 w3-margin-0 w3-light-grey">{{$value_list_inner.Line}}</pre>
			{{else}}
			  <pre class="w3-padding-0 w3-margin-0">{{$value_list_inner.Line}}</pre>
			{{end}}
		{{end}}	
		    </div>
		  </div>
		 </p>
	{{end}}	 
		</td>
       </tr>
    {{else if eq $value_list.Value $WarningVar}}
	   <tr>
        <td colspan="2" class="w3-padding-0">
		 <p class="w3-container w3-khaki w3-margin-0 w3-padding-0  w3-padding-left">{{$value_list.Plugin}}: {{$value_list.Line}}</p>
	{{if $value_list.Link}}
	     <p class="w3-container w3-khaki w3-margin-0 w3-padding-0  w3-padding-left">
	      <div class="w3-dropdown-hover">
		    <button class="w3-btn w3-red w3-tiny w3-padding-0 w3-margin-0">Link</button>
		    <div class="w3-dropdown-content w3-border w3-light-grey">
		{{range $iiind, $value_list_inner := $value_list.Link}}
			{{if eq $value_list_inner.FndStr false}}
			  <pre class="w3-padding-0 w3-margin-0 w3-light-grey">{{$value_list_inner.Line}}</pre>
			{{else}}
			  <pre class="w3-padding-0 w3-margin-0">{{$value_list_inner.Line}}</pre>
			{{end}}
		{{end}}	
		    </div>
		  </div>
		 </p>
	{{end}}	 
		</td>
       </tr>
    {{else if eq $value_list.Value $DangerVar}}
	   <tr>
        <td colspan="2" class="w3-padding-0">
		 <p class="w3-container w3-red w3-margin-0 w3-padding-0  w3-padding-left">{{$value_list.Plugin}}: {{$value_list.Line}}</p>
	{{if $value_list.Link}}
	     <p class="w3-container w3-red w3-margin-0 w3-padding-0  w3-padding-left">
	      <div class="w3-dropdown-hover">
		    <button class="w3-btn w3-red w3-tiny w3-padding-0 w3-margin-0">Link</button>
		    <div class="w3-dropdown-content w3-border w3-light-grey">
		{{range $iiind, $value_list_inner := $value_list.Link}}
			{{if eq $value_list_inner.FndStr false}}
			  <pre class="w3-padding-0 w3-margin-0 w3-light-grey">{{$value_list_inner.Line}}</pre>
			{{else}}
			  <pre class="w3-padding-0 w3-margin-0">{{$value_list_inner.Line}}</pre>
			{{end}}
		{{end}}	
		    </div>
		  </div>
		 </p>
	{{end}}	
		</td>
       </tr>
    {{else if eq $value_list.Value $ErrVar}}
       <tr>
        <td style="width:80px" class="w3-padding-0 w3-padding-left">{{$value_list.Number}}</td>
        <td class="w3-border-green w3-border-left w3-padding-0  w3-padding-left w3-light-grey"><pre class="w3-padding-0 w3-margin-0 w3-light-grey">{{$value_list.Line}}</pre></td>
       </tr>
    {{else if eq $value_list.Value $ErrContVar}}
       <tr>
        <td style="width:80px" class="w3-padding-0 w3-padding-left">{{$value_list.Number}}</td>
        <td class="w3-border-green w3-border-left w3-padding-0  w3-padding-left w3-light-grey"><pre class="w3-padding-0 w3-margin-0 w3-light-grey">{{$value_list.Line}}</pre></td>
       </tr>
    {{else}}
	   <tr>
        <td colspan="2">
		 <p>...</p>
		</td>
       </tr>   
    {{end}}
{{end}}
      </table>
	</div>
  </div>
{{end}}
</div>

<script>
function myFunction(id) {
    var x = document.getElementById(id);
    if (x.className.indexOf("w3-show") == -1) {
        x.className += " w3-show";
		x.previousElementSibling.className = 
        x.previousElementSibling.className.replace(" w3-khaki", " w3-blue");
    } else { 
        x.className = x.className.replace(" w3-show", "");
        x.previousElementSibling.className = 
        x.previousElementSibling.className.replace(" w3-blue", " w3-khaki");
    }
}

window.onload = function() {
    var aCodes = document.getElementsByTagName('pre');
    for (var i=0; i < aCodes.length; i++) {
        hljs.highlightBlock(aCodes[i]);
    }
};
</script>

</body>
</html>








