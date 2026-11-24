import{a as be,b as xe,f as we,i as At,j as $t,k as Lt,l as Ft,m as Ot,n as _e,o as De,p as Se,q as Ce,r as Me,s as Ee,t as Ie,u as Wt,v as Vt,y as Ye,z as Ae}from"./chunk-MKCXTBEZ.js";import{p as $e}from"./chunk-7TIOW63O.js";import{x as si}from"./chunk-MNVP5QMN.js";import{N as fe,O as he,T as me,U as ke,V as ye,W as ge,X as pe,Y as ve,Z as Te,_ as ot}from"./chunk-7JYBA3TO.js";import{a as de,b as nt,d as vt}from"./chunk-K6ALI33Z.js";import{a as l}from"./chunk-PAESYWDV.js";import{b as wt,c as at}from"./chunk-5GVNIGXK.js";var Ve=wt((Rt,Nt)=>{(function(t,e){typeof Rt=="object"&&typeof Nt<"u"?Nt.exports=e():typeof define=="function"&&define.amd?define(e):(t=typeof globalThis<"u"?globalThis:t||self).dayjs_plugin_isoWeek=e()})(Rt,(function(){"use strict";var t="day";return function(e,s,r){var n=function(w){return w.add(4-w.isoWeekday(),t)},d=s.prototype;d.isoWeekYear=function(){return n(this).year()},d.isoWeek=function(w){if(!this.$utils().u(w))return this.add(7*(w-this.isoWeek()),t);var x,V,I,H,M=n(this),D=(x=this.isoWeekYear(),V=this.$u,I=(V?r.utc:r)().year(x).startOf("year"),H=4-I.isoWeekday(),I.isoWeekday()>4&&(H+=7),I.add(H,t));return M.diff(D,"week")+1},d.isoWeekday=function(w){return this.$utils().u(w)?this.day()||7:this.day(this.day()%7?w:w-7)};var y=d.startOf;d.startOf=function(w,x){var V=this.$utils(),I=!!V.u(x)||x;return V.p(w)==="isoweek"?I?this.date(this.date()-(this.isoWeekday()-1)).startOf("day"):this.date(this.date()-1-(this.isoWeekday()-1)+7).endOf("day"):y.bind(this)(w,x)}}}))});var Pe=wt((Ht,Bt)=>{(function(t,e){typeof Ht=="object"&&typeof Bt<"u"?Bt.exports=e():typeof define=="function"&&define.amd?define(e):(t=typeof globalThis<"u"?globalThis:t||self).dayjs_plugin_customParseFormat=e()})(Ht,(function(){"use strict";var t={LTS:"h:mm:ss A",LT:"h:mm A",L:"MM/DD/YYYY",LL:"MMMM D, YYYY",LLL:"MMMM D, YYYY h:mm A",LLLL:"dddd, MMMM D, YYYY h:mm A"},e=/(\[[^[]*\])|([-_:/.,()\s]+)|(A|a|Q|YYYY|YY?|ww?|MM?M?M?|Do|DD?|hh?|HH?|mm?|ss?|S{1,3}|z|ZZ?)/g,s=/\d/,r=/\d\d/,n=/\d\d?/,d=/\d*[^-_:/,()\s\d]+/,y={},w=function(_){return(_=+_)+(_>68?1900:2e3)},x=function(_){return function(S){this[_]=+S}},V=[/[+-]\d\d:?(\d\d)?|Z/,function(_){(this.zone||(this.zone={})).offset=(function(S){if(!S||S==="Z")return 0;var b=S.match(/([+-]|\d\d)/g),O=60*b[1]+(+b[2]||0);return O===0?0:b[0]==="+"?-O:O})(_)}],I=function(_){var S=y[_];return S&&(S.indexOf?S:S.s.concat(S.f))},H=function(_,S){var b,O=y.meridiem;if(O){for(var B=1;B<=24;B+=1)if(_.indexOf(O(B,0,S))>-1){b=B>12;break}}else b=_===(S?"pm":"PM");return b},M={A:[d,function(_){this.afternoon=H(_,!1)}],a:[d,function(_){this.afternoon=H(_,!0)}],Q:[s,function(_){this.month=3*(_-1)+1}],S:[s,function(_){this.milliseconds=100*+_}],SS:[r,function(_){this.milliseconds=10*+_}],SSS:[/\d{3}/,function(_){this.milliseconds=+_}],s:[n,x("seconds")],ss:[n,x("seconds")],m:[n,x("minutes")],mm:[n,x("minutes")],H:[n,x("hours")],h:[n,x("hours")],HH:[n,x("hours")],hh:[n,x("hours")],D:[n,x("day")],DD:[r,x("day")],Do:[d,function(_){var S=y.ordinal,b=_.match(/\d+/);if(this.day=b[0],S)for(var O=1;O<=31;O+=1)S(O).replace(/\[|\]/g,"")===_&&(this.day=O)}],w:[n,x("week")],ww:[r,x("week")],M:[n,x("month")],MM:[r,x("month")],MMM:[d,function(_){var S=I("months"),b=(I("monthsShort")||S.map((function(O){return O.slice(0,3)}))).indexOf(_)+1;if(b<1)throw new Error;this.month=b%12||b}],MMMM:[d,function(_){var S=I("months").indexOf(_)+1;if(S<1)throw new Error;this.month=S%12||S}],Y:[/[+-]?\d+/,x("year")],YY:[r,function(_){this.year=w(_)}],YYYY:[/\d{4}/,x("year")],Z:V,ZZ:V};function D(_){var S,b;S=_,b=y&&y.formats;for(var O=(_=S.replace(/(\[[^\]]+])|(LTS?|l{1,4}|L{1,4})/g,(function(W,F,m){var v=m&&m.toUpperCase();return F||b[m]||t[m]||b[v].replace(/(\[[^\]]+])|(MMMM|MM|DD|dddd)/g,(function(g,T,a){return T||a.slice(1)}))}))).match(e),B=O.length,j=0;j<B;j+=1){var L=O[j],p=M[L],f=p&&p[0],C=p&&p[1];O[j]=C?{regex:f,parser:C}:L.replace(/^\[|\]$/g,"")}return function(W){for(var F={},m=0,v=0;m<B;m+=1){var g=O[m];if(typeof g=="string")v+=g.length;else{var T=g.regex,a=g.parser,k=W.slice(v),h=T.exec(k)[0];a.call(F,h),W=W.replace(h,"")}}return(function(u){var Y=u.afternoon;if(Y!==void 0){var i=u.hours;Y?i<12&&(u.hours+=12):i===12&&(u.hours=0),delete u.afternoon}})(F),F}}return function(_,S,b){b.p.customParseFormat=!0,_&&_.parseTwoDigitYear&&(w=_.parseTwoDigitYear);var O=S.prototype,B=O.parse;O.parse=function(j){var L=j.date,p=j.utc,f=j.args;this.$u=p;var C=f[1];if(typeof C=="string"){var W=f[2]===!0,F=f[3]===!0,m=W||F,v=f[2];F&&(v=f[2]),y=this.$locale(),!W&&v&&(y=b.Ls[v]),this.$d=(function(k,h,u,Y){try{if(["x","X"].indexOf(h)>-1)return new Date((h==="X"?1e3:1)*k);var i=D(h)(k),A=i.year,c=i.month,G=i.day,o=i.hours,E=i.minutes,$=i.seconds,R=i.milliseconds,N=i.zone,P=i.week,z=new Date,tt=G||(A||c?1:z.getDate()),et=A||z.getFullYear(),lt=0;A&&!c||(lt=c>0?c-1:z.getMonth());var ut,dt=o||0,X=E||0,st=$||0,K=R||0;return N?new Date(Date.UTC(et,lt,tt,dt,X,st,K+60*N.offset*1e3)):u?new Date(Date.UTC(et,lt,tt,dt,X,st,K)):(ut=new Date(et,lt,tt,dt,X,st,K),P&&(ut=Y(ut).week(P).toDate()),ut)}catch{return new Date("")}})(L,C,p,b),this.init(),v&&v!==!0&&(this.$L=this.locale(v).$L),m&&L!=this.format(C)&&(this.$d=new Date("")),y={}}else if(C instanceof Array)for(var g=C.length,T=1;T<=g;T+=1){f[1]=C[T-1];var a=b.apply(this,f);if(a.isValid()){this.$d=a.$d,this.$L=a.$L,this.init();break}T===g&&(this.$d=new Date(""))}else B.call(this,j)}}}))});var ze=wt((jt,Gt)=>{(function(t,e){typeof jt=="object"&&typeof Gt<"u"?Gt.exports=e():typeof define=="function"&&define.amd?define(e):(t=typeof globalThis<"u"?globalThis:t||self).dayjs_plugin_advancedFormat=e()})(jt,(function(){"use strict";return function(t,e){var s=e.prototype,r=s.format;s.format=function(n){var d=this,y=this.$locale();if(!this.isValid())return r.bind(this)(n);var w=this.$utils(),x=(n||"YYYY-MM-DDTHH:mm:ssZ").replace(/\[([^\]]+)]|Q|wo|ww|w|WW|W|zzz|z|gggg|GGGG|Do|X|x|k{1,2}|S/g,(function(V){switch(V){case"Q":return Math.ceil((d.$M+1)/3);case"Do":return y.ordinal(d.$D);case"gggg":return d.weekYear();case"GGGG":return d.isoWeekYear();case"wo":return y.ordinal(d.week(),"W");case"w":case"ww":return w.s(d.week(),V==="w"?1:2,"0");case"W":case"WW":return w.s(d.isoWeek(),V==="W"?1:2,"0");case"k":case"kk":return w.s(String(d.$H===0?24:d.$H),V==="k"?1:2,"0");case"X":return Math.floor(d.$d.getTime()/1e3);case"x":return d.$d.getTime();case"z":return"["+d.offsetName()+"]";case"zzz":return"["+d.offsetName("long")+"]";default:return V}}));return r.bind(this)(x)}}}))});var Re=wt((Xt,Ut)=>{(function(t,e){typeof Xt=="object"&&typeof Ut<"u"?Ut.exports=e():typeof define=="function"&&define.amd?define(e):(t=typeof globalThis<"u"?globalThis:t||self).dayjs_plugin_duration=e()})(Xt,(function(){"use strict";var t,e,s=1e3,r=6e4,n=36e5,d=864e5,y=31536e6,w=2628e6,x=/^(-|\+)?P(?:([-+]?[0-9,.]*)Y)?(?:([-+]?[0-9,.]*)M)?(?:([-+]?[0-9,.]*)W)?(?:([-+]?[0-9,.]*)D)?(?:T(?:([-+]?[0-9,.]*)H)?(?:([-+]?[0-9,.]*)M)?(?:([-+]?[0-9,.]*)S)?)?$/,V=/\[([^\]]+)]|YYYY|YY|Y|M{1,2}|D{1,2}|H{1,2}|m{1,2}|s{1,2}|SSS/g,I={years:y,months:w,days:d,hours:n,minutes:r,seconds:s,milliseconds:1,weeks:6048e5},H=function(L){return L instanceof B},M=function(L,p,f){return new B(L,f,p.$l)},D=function(L){return e.p(L)+"s"},_=function(L){return L<0},S=function(L){return _(L)?Math.ceil(L):Math.floor(L)},b=function(L){return Math.abs(L)},O=function(L,p){return L?_(L)?{negative:!0,format:""+b(L)+p}:{negative:!1,format:""+L+p}:{negative:!1,format:""}},B=(function(){function L(f,C,W){var F=this;if(this.$d={},this.$l=W,f===void 0&&(this.$ms=0,this.parseFromMilliseconds()),C)return M(f*I[D(C)],this);if(typeof f=="number")return this.$ms=f,this.parseFromMilliseconds(),this;if(typeof f=="object")return Object.keys(f).forEach((function(g){F.$d[D(g)]=f[g]})),this.calMilliseconds(),this;if(typeof f=="string"){var m=f.match(x);if(m){var v=m.slice(2).map((function(g){return g!=null?Number(g):0}));return this.$d.years=v[0],this.$d.months=v[1],this.$d.weeks=v[2],this.$d.days=v[3],this.$d.hours=v[4],this.$d.minutes=v[5],this.$d.seconds=v[6],this.calMilliseconds(),this}}return this}var p=L.prototype;return p.calMilliseconds=function(){var f=this;this.$ms=Object.keys(this.$d).reduce((function(C,W){return C+(f.$d[W]||0)*I[W]}),0)},p.parseFromMilliseconds=function(){var f=this.$ms;this.$d.years=S(f/y),f%=y,this.$d.months=S(f/w),f%=w,this.$d.days=S(f/d),f%=d,this.$d.hours=S(f/n),f%=n,this.$d.minutes=S(f/r),f%=r,this.$d.seconds=S(f/s),f%=s,this.$d.milliseconds=f},p.toISOString=function(){var f=O(this.$d.years,"Y"),C=O(this.$d.months,"M"),W=+this.$d.days||0;this.$d.weeks&&(W+=7*this.$d.weeks);var F=O(W,"D"),m=O(this.$d.hours,"H"),v=O(this.$d.minutes,"M"),g=this.$d.seconds||0;this.$d.milliseconds&&(g+=this.$d.milliseconds/1e3,g=Math.round(1e3*g)/1e3);var T=O(g,"S"),a=f.negative||C.negative||F.negative||m.negative||v.negative||T.negative,k=m.format||v.format||T.format?"T":"",h=(a?"-":"")+"P"+f.format+C.format+F.format+k+m.format+v.format+T.format;return h==="P"||h==="-P"?"P0D":h},p.toJSON=function(){return this.toISOString()},p.format=function(f){var C=f||"YYYY-MM-DDTHH:mm:ss",W={Y:this.$d.years,YY:e.s(this.$d.years,2,"0"),YYYY:e.s(this.$d.years,4,"0"),M:this.$d.months,MM:e.s(this.$d.months,2,"0"),D:this.$d.days,DD:e.s(this.$d.days,2,"0"),H:this.$d.hours,HH:e.s(this.$d.hours,2,"0"),m:this.$d.minutes,mm:e.s(this.$d.minutes,2,"0"),s:this.$d.seconds,ss:e.s(this.$d.seconds,2,"0"),SSS:e.s(this.$d.milliseconds,3,"0")};return C.replace(V,(function(F,m){return m||String(W[F])}))},p.as=function(f){return this.$ms/I[D(f)]},p.get=function(f){var C=this.$ms,W=D(f);return W==="milliseconds"?C%=1e3:C=W==="weeks"?S(C/I[W]):this.$d[W],C||0},p.add=function(f,C,W){var F;return F=C?f*I[D(C)]:H(f)?f.$ms:M(f,this).$ms,M(this.$ms+F*(W?-1:1),this)},p.subtract=function(f,C){return this.add(f,C,!0)},p.locale=function(f){var C=this.clone();return C.$l=f,C},p.clone=function(){return M(this.$ms,this)},p.humanize=function(f){return t().add(this.$ms,"ms").locale(this.$l).fromNow(!f)},p.valueOf=function(){return this.asMilliseconds()},p.milliseconds=function(){return this.get("milliseconds")},p.asMilliseconds=function(){return this.as("milliseconds")},p.seconds=function(){return this.get("seconds")},p.asSeconds=function(){return this.as("seconds")},p.minutes=function(){return this.get("minutes")},p.asMinutes=function(){return this.as("minutes")},p.hours=function(){return this.get("hours")},p.asHours=function(){return this.as("hours")},p.days=function(){return this.get("days")},p.asDays=function(){return this.as("days")},p.weeks=function(){return this.get("weeks")},p.asWeeks=function(){return this.as("weeks")},p.months=function(){return this.get("months")},p.asMonths=function(){return this.as("months")},p.years=function(){return this.get("years")},p.asYears=function(){return this.as("years")},L})(),j=function(L,p,f){return L.add(p.years()*f,"y").add(p.months()*f,"M").add(p.days()*f,"d").add(p.hours()*f,"h").add(p.minutes()*f,"m").add(p.seconds()*f,"s").add(p.milliseconds()*f,"ms")};return function(L,p,f){t=f,e=f().$utils(),f.duration=function(F,m){var v=f.locale();return M(F,{$l:v},m)},f.isDuration=H;var C=p.prototype.add,W=p.prototype.subtract;p.prototype.add=function(F,m){return H(F)?j(this,F,1):C.bind(this)(F,m)},p.prototype.subtract=function(F,m){return H(F)?j(this,F,-1):W.bind(this)(F,m)}}}))});function ai(t){return t}var Dt=1,Pt=2,zt=3,_t=4,Le=1e-6;function oi(t){return"translate("+t+",0)"}function ci(t){return"translate(0,"+t+")"}function li(t){return e=>+t(e)}function ui(t,e){return e=Math.max(0,t.bandwidth()-e*2)/2,t.round()&&(e=Math.round(e)),s=>+t(s)+e}function di(){return!this.__axis}function Fe(t,e){var s=[],r=null,n=null,d=6,y=6,w=3,x=typeof window<"u"&&window.devicePixelRatio>1?0:.5,V=t===Dt||t===_t?-1:1,I=t===_t||t===Pt?"x":"y",H=t===Dt||t===zt?oi:ci;function M(D){var _=r??(e.ticks?e.ticks.apply(e,s):e.domain()),S=n??(e.tickFormat?e.tickFormat.apply(e,s):ai),b=Math.max(d,0)+w,O=e.range(),B=+O[0]+x,j=+O[O.length-1]+x,L=(e.bandwidth?ui:li)(e.copy(),x),p=D.selection?D.selection():D,f=p.selectAll(".domain").data([null]),C=p.selectAll(".tick").data(_,e).order(),W=C.exit(),F=C.enter().append("g").attr("class","tick"),m=C.select("line"),v=C.select("text");f=f.merge(f.enter().insert("path",".tick").attr("class","domain").attr("stroke","currentColor")),C=C.merge(F),m=m.merge(F.append("line").attr("stroke","currentColor").attr(I+"2",V*d)),v=v.merge(F.append("text").attr("fill","currentColor").attr(I,V*b).attr("dy",t===Dt?"0em":t===zt?"0.71em":"0.32em")),D!==p&&(f=f.transition(D),C=C.transition(D),m=m.transition(D),v=v.transition(D),W=W.transition(D).attr("opacity",Le).attr("transform",function(g){return isFinite(g=L(g))?H(g+x):this.getAttribute("transform")}),F.attr("opacity",Le).attr("transform",function(g){var T=this.parentNode.__axis;return H((T&&isFinite(T=T(g))?T:L(g))+x)})),W.remove(),f.attr("d",t===_t||t===Pt?y?"M"+V*y+","+B+"H"+x+"V"+j+"H"+V*y:"M"+x+","+B+"V"+j:y?"M"+B+","+V*y+"V"+x+"H"+j+"V"+V*y:"M"+B+","+x+"H"+j),C.attr("opacity",1).attr("transform",function(g){return H(L(g)+x)}),m.attr(I+"2",V*d),v.attr(I,V*b).text(S),p.filter(di).attr("fill","none").attr("font-size",10).attr("font-family","sans-serif").attr("text-anchor",t===Pt?"start":t===_t?"end":"middle"),p.each(function(){this.__axis=L})}return M.scale=function(D){return arguments.length?(e=D,M):e},M.ticks=function(){return s=Array.from(arguments),M},M.tickArguments=function(D){return arguments.length?(s=D==null?[]:Array.from(D),M):s.slice()},M.tickValues=function(D){return arguments.length?(r=D==null?null:Array.from(D),M):r&&r.slice()},M.tickFormat=function(D){return arguments.length?(n=D,M):n},M.tickSize=function(D){return arguments.length?(d=y=+D,M):d},M.tickSizeInner=function(D){return arguments.length?(d=+D,M):d},M.tickSizeOuter=function(D){return arguments.length?(y=+D,M):y},M.tickPadding=function(D){return arguments.length?(w=+D,M):w},M.offset=function(D){return arguments.length?(x=+D,M):x},M}function Oe(t){return Fe(Dt,t)}function We(t){return Fe(zt,t)}var je=at(si(),1),Q=at(de(),1),Ge=at(Ve(),1),Xe=at(Pe(),1),Ue=at(ze(),1),mt=at(de(),1),ni=at(Re(),1);var Zt=(function(){var t=l(function(T,a,k,h){for(k=k||{},h=T.length;h--;k[T[h]]=a);return k},"o"),e=[6,8,10,12,13,14,15,16,17,18,20,21,22,23,24,25,26,27,28,29,30,31,33,35,36,38,40],s=[1,26],r=[1,27],n=[1,28],d=[1,29],y=[1,30],w=[1,31],x=[1,32],V=[1,33],I=[1,34],H=[1,9],M=[1,10],D=[1,11],_=[1,12],S=[1,13],b=[1,14],O=[1,15],B=[1,16],j=[1,19],L=[1,20],p=[1,21],f=[1,22],C=[1,23],W=[1,25],F=[1,35],m={trace:l(function(){},"trace"),yy:{},symbols_:{error:2,start:3,gantt:4,document:5,EOF:6,line:7,SPACE:8,statement:9,NL:10,weekday:11,weekday_monday:12,weekday_tuesday:13,weekday_wednesday:14,weekday_thursday:15,weekday_friday:16,weekday_saturday:17,weekday_sunday:18,weekend:19,weekend_friday:20,weekend_saturday:21,dateFormat:22,inclusiveEndDates:23,topAxis:24,axisFormat:25,tickInterval:26,excludes:27,includes:28,todayMarker:29,title:30,acc_title:31,acc_title_value:32,acc_descr:33,acc_descr_value:34,acc_descr_multiline_value:35,section:36,clickStatement:37,taskTxt:38,taskData:39,click:40,callbackname:41,callbackargs:42,href:43,clickStatementDebug:44,$accept:0,$end:1},terminals_:{2:"error",4:"gantt",6:"EOF",8:"SPACE",10:"NL",12:"weekday_monday",13:"weekday_tuesday",14:"weekday_wednesday",15:"weekday_thursday",16:"weekday_friday",17:"weekday_saturday",18:"weekday_sunday",20:"weekend_friday",21:"weekend_saturday",22:"dateFormat",23:"inclusiveEndDates",24:"topAxis",25:"axisFormat",26:"tickInterval",27:"excludes",28:"includes",29:"todayMarker",30:"title",31:"acc_title",32:"acc_title_value",33:"acc_descr",34:"acc_descr_value",35:"acc_descr_multiline_value",36:"section",38:"taskTxt",39:"taskData",40:"click",41:"callbackname",42:"callbackargs",43:"href"},productions_:[0,[3,3],[5,0],[5,2],[7,2],[7,1],[7,1],[7,1],[11,1],[11,1],[11,1],[11,1],[11,1],[11,1],[11,1],[19,1],[19,1],[9,1],[9,1],[9,1],[9,1],[9,1],[9,1],[9,1],[9,1],[9,1],[9,1],[9,1],[9,2],[9,2],[9,1],[9,1],[9,1],[9,2],[37,2],[37,3],[37,3],[37,4],[37,3],[37,4],[37,2],[44,2],[44,3],[44,3],[44,4],[44,3],[44,4],[44,2]],performAction:l(function(a,k,h,u,Y,i,A){var c=i.length-1;switch(Y){case 1:return i[c-1];case 2:this.$=[];break;case 3:i[c-1].push(i[c]),this.$=i[c-1];break;case 4:case 5:this.$=i[c];break;case 6:case 7:this.$=[];break;case 8:u.setWeekday("monday");break;case 9:u.setWeekday("tuesday");break;case 10:u.setWeekday("wednesday");break;case 11:u.setWeekday("thursday");break;case 12:u.setWeekday("friday");break;case 13:u.setWeekday("saturday");break;case 14:u.setWeekday("sunday");break;case 15:u.setWeekend("friday");break;case 16:u.setWeekend("saturday");break;case 17:u.setDateFormat(i[c].substr(11)),this.$=i[c].substr(11);break;case 18:u.enableInclusiveEndDates(),this.$=i[c].substr(18);break;case 19:u.TopAxis(),this.$=i[c].substr(8);break;case 20:u.setAxisFormat(i[c].substr(11)),this.$=i[c].substr(11);break;case 21:u.setTickInterval(i[c].substr(13)),this.$=i[c].substr(13);break;case 22:u.setExcludes(i[c].substr(9)),this.$=i[c].substr(9);break;case 23:u.setIncludes(i[c].substr(9)),this.$=i[c].substr(9);break;case 24:u.setTodayMarker(i[c].substr(12)),this.$=i[c].substr(12);break;case 27:u.setDiagramTitle(i[c].substr(6)),this.$=i[c].substr(6);break;case 28:this.$=i[c].trim(),u.setAccTitle(this.$);break;case 29:case 30:this.$=i[c].trim(),u.setAccDescription(this.$);break;case 31:u.addSection(i[c].substr(8)),this.$=i[c].substr(8);break;case 33:u.addTask(i[c-1],i[c]),this.$="task";break;case 34:this.$=i[c-1],u.setClickEvent(i[c-1],i[c],null);break;case 35:this.$=i[c-2],u.setClickEvent(i[c-2],i[c-1],i[c]);break;case 36:this.$=i[c-2],u.setClickEvent(i[c-2],i[c-1],null),u.setLink(i[c-2],i[c]);break;case 37:this.$=i[c-3],u.setClickEvent(i[c-3],i[c-2],i[c-1]),u.setLink(i[c-3],i[c]);break;case 38:this.$=i[c-2],u.setClickEvent(i[c-2],i[c],null),u.setLink(i[c-2],i[c-1]);break;case 39:this.$=i[c-3],u.setClickEvent(i[c-3],i[c-1],i[c]),u.setLink(i[c-3],i[c-2]);break;case 40:this.$=i[c-1],u.setLink(i[c-1],i[c]);break;case 41:case 47:this.$=i[c-1]+" "+i[c];break;case 42:case 43:case 45:this.$=i[c-2]+" "+i[c-1]+" "+i[c];break;case 44:case 46:this.$=i[c-3]+" "+i[c-2]+" "+i[c-1]+" "+i[c];break}},"anonymous"),table:[{3:1,4:[1,2]},{1:[3]},t(e,[2,2],{5:3}),{6:[1,4],7:5,8:[1,6],9:7,10:[1,8],11:17,12:s,13:r,14:n,15:d,16:y,17:w,18:x,19:18,20:V,21:I,22:H,23:M,24:D,25:_,26:S,27:b,28:O,29:B,30:j,31:L,33:p,35:f,36:C,37:24,38:W,40:F},t(e,[2,7],{1:[2,1]}),t(e,[2,3]),{9:36,11:17,12:s,13:r,14:n,15:d,16:y,17:w,18:x,19:18,20:V,21:I,22:H,23:M,24:D,25:_,26:S,27:b,28:O,29:B,30:j,31:L,33:p,35:f,36:C,37:24,38:W,40:F},t(e,[2,5]),t(e,[2,6]),t(e,[2,17]),t(e,[2,18]),t(e,[2,19]),t(e,[2,20]),t(e,[2,21]),t(e,[2,22]),t(e,[2,23]),t(e,[2,24]),t(e,[2,25]),t(e,[2,26]),t(e,[2,27]),{32:[1,37]},{34:[1,38]},t(e,[2,30]),t(e,[2,31]),t(e,[2,32]),{39:[1,39]},t(e,[2,8]),t(e,[2,9]),t(e,[2,10]),t(e,[2,11]),t(e,[2,12]),t(e,[2,13]),t(e,[2,14]),t(e,[2,15]),t(e,[2,16]),{41:[1,40],43:[1,41]},t(e,[2,4]),t(e,[2,28]),t(e,[2,29]),t(e,[2,33]),t(e,[2,34],{42:[1,42],43:[1,43]}),t(e,[2,40],{41:[1,44]}),t(e,[2,35],{43:[1,45]}),t(e,[2,36]),t(e,[2,38],{42:[1,46]}),t(e,[2,37]),t(e,[2,39])],defaultActions:{},parseError:l(function(a,k){if(k.recoverable)this.trace(a);else{var h=new Error(a);throw h.hash=k,h}},"parseError"),parse:l(function(a){var k=this,h=[0],u=[],Y=[null],i=[],A=this.table,c="",G=0,o=0,E=0,$=2,R=1,N=i.slice.call(arguments,1),P=Object.create(this.lexer),z={yy:{}};for(var tt in this.yy)Object.prototype.hasOwnProperty.call(this.yy,tt)&&(z.yy[tt]=this.yy[tt]);P.setInput(a,z.yy),z.yy.lexer=P,z.yy.parser=this,typeof P.yylloc>"u"&&(P.yylloc={});var et=P.yylloc;i.push(et);var lt=P.options&&P.options.ranges;typeof z.yy.parseError=="function"?this.parseError=z.yy.parseError:this.parseError=Object.getPrototypeOf(this).parseError;function ut(Z){h.length=h.length-2*Z,Y.length=Y.length-Z,i.length=i.length-Z}l(ut,"popStack");function dt(){var Z;return Z=u.pop()||P.lex()||R,typeof Z!="number"&&(Z instanceof Array&&(u=Z,Z=u.pop()),Z=k.symbols_[Z]||Z),Z}l(dt,"lex");for(var X,st,K,q,ir,It,ft={},bt,it,ue,xt;;){if(K=h[h.length-1],this.defaultActions[K]?q=this.defaultActions[K]:((X===null||typeof X>"u")&&(X=dt()),q=A[K]&&A[K][X]),typeof q>"u"||!q.length||!q[0]){var Yt="";xt=[];for(bt in A[K])this.terminals_[bt]&&bt>$&&xt.push("'"+this.terminals_[bt]+"'");P.showPosition?Yt="Parse error on line "+(G+1)+`:
`+P.showPosition()+`
Expecting `+xt.join(", ")+", got '"+(this.terminals_[X]||X)+"'":Yt="Parse error on line "+(G+1)+": Unexpected "+(X==R?"end of input":"'"+(this.terminals_[X]||X)+"'"),this.parseError(Yt,{text:P.match,token:this.terminals_[X]||X,line:P.yylineno,loc:et,expected:xt})}if(q[0]instanceof Array&&q.length>1)throw new Error("Parse Error: multiple actions possible at state: "+K+", token: "+X);switch(q[0]){case 1:h.push(X),Y.push(P.yytext),i.push(P.yylloc),h.push(q[1]),X=null,st?(X=st,st=null):(o=P.yyleng,c=P.yytext,G=P.yylineno,et=P.yylloc,E>0&&E--);break;case 2:if(it=this.productions_[q[1]][1],ft.$=Y[Y.length-it],ft._$={first_line:i[i.length-(it||1)].first_line,last_line:i[i.length-1].last_line,first_column:i[i.length-(it||1)].first_column,last_column:i[i.length-1].last_column},lt&&(ft._$.range=[i[i.length-(it||1)].range[0],i[i.length-1].range[1]]),It=this.performAction.apply(ft,[c,o,G,z.yy,q[1],Y,i].concat(N)),typeof It<"u")return It;it&&(h=h.slice(0,-1*it*2),Y=Y.slice(0,-1*it),i=i.slice(0,-1*it)),h.push(this.productions_[q[1]][0]),Y.push(ft.$),i.push(ft._$),ue=A[h[h.length-2]][h[h.length-1]],h.push(ue);break;case 3:return!0}}return!0},"parse")},v=(function(){var T={EOF:1,parseError:l(function(k,h){if(this.yy.parser)this.yy.parser.parseError(k,h);else throw new Error(k)},"parseError"),setInput:l(function(a,k){return this.yy=k||this.yy||{},this._input=a,this._more=this._backtrack=this.done=!1,this.yylineno=this.yyleng=0,this.yytext=this.matched=this.match="",this.conditionStack=["INITIAL"],this.yylloc={first_line:1,first_column:0,last_line:1,last_column:0},this.options.ranges&&(this.yylloc.range=[0,0]),this.offset=0,this},"setInput"),input:l(function(){var a=this._input[0];this.yytext+=a,this.yyleng++,this.offset++,this.match+=a,this.matched+=a;var k=a.match(/(?:\r\n?|\n).*/g);return k?(this.yylineno++,this.yylloc.last_line++):this.yylloc.last_column++,this.options.ranges&&this.yylloc.range[1]++,this._input=this._input.slice(1),a},"input"),unput:l(function(a){var k=a.length,h=a.split(/(?:\r\n?|\n)/g);this._input=a+this._input,this.yytext=this.yytext.substr(0,this.yytext.length-k),this.offset-=k;var u=this.match.split(/(?:\r\n?|\n)/g);this.match=this.match.substr(0,this.match.length-1),this.matched=this.matched.substr(0,this.matched.length-1),h.length-1&&(this.yylineno-=h.length-1);var Y=this.yylloc.range;return this.yylloc={first_line:this.yylloc.first_line,last_line:this.yylineno+1,first_column:this.yylloc.first_column,last_column:h?(h.length===u.length?this.yylloc.first_column:0)+u[u.length-h.length].length-h[0].length:this.yylloc.first_column-k},this.options.ranges&&(this.yylloc.range=[Y[0],Y[0]+this.yyleng-k]),this.yyleng=this.yytext.length,this},"unput"),more:l(function(){return this._more=!0,this},"more"),reject:l(function(){if(this.options.backtrack_lexer)this._backtrack=!0;else return this.parseError("Lexical error on line "+(this.yylineno+1)+`. You can only invoke reject() in the lexer when the lexer is of the backtracking persuasion (options.backtrack_lexer = true).
`+this.showPosition(),{text:"",token:null,line:this.yylineno});return this},"reject"),less:l(function(a){this.unput(this.match.slice(a))},"less"),pastInput:l(function(){var a=this.matched.substr(0,this.matched.length-this.match.length);return(a.length>20?"...":"")+a.substr(-20).replace(/\n/g,"")},"pastInput"),upcomingInput:l(function(){var a=this.match;return a.length<20&&(a+=this._input.substr(0,20-a.length)),(a.substr(0,20)+(a.length>20?"...":"")).replace(/\n/g,"")},"upcomingInput"),showPosition:l(function(){var a=this.pastInput(),k=new Array(a.length+1).join("-");return a+this.upcomingInput()+`
`+k+"^"},"showPosition"),test_match:l(function(a,k){var h,u,Y;if(this.options.backtrack_lexer&&(Y={yylineno:this.yylineno,yylloc:{first_line:this.yylloc.first_line,last_line:this.last_line,first_column:this.yylloc.first_column,last_column:this.yylloc.last_column},yytext:this.yytext,match:this.match,matches:this.matches,matched:this.matched,yyleng:this.yyleng,offset:this.offset,_more:this._more,_input:this._input,yy:this.yy,conditionStack:this.conditionStack.slice(0),done:this.done},this.options.ranges&&(Y.yylloc.range=this.yylloc.range.slice(0))),u=a[0].match(/(?:\r\n?|\n).*/g),u&&(this.yylineno+=u.length),this.yylloc={first_line:this.yylloc.last_line,last_line:this.yylineno+1,first_column:this.yylloc.last_column,last_column:u?u[u.length-1].length-u[u.length-1].match(/\r?\n?/)[0].length:this.yylloc.last_column+a[0].length},this.yytext+=a[0],this.match+=a[0],this.matches=a,this.yyleng=this.yytext.length,this.options.ranges&&(this.yylloc.range=[this.offset,this.offset+=this.yyleng]),this._more=!1,this._backtrack=!1,this._input=this._input.slice(a[0].length),this.matched+=a[0],h=this.performAction.call(this,this.yy,this,k,this.conditionStack[this.conditionStack.length-1]),this.done&&this._input&&(this.done=!1),h)return h;if(this._backtrack){for(var i in Y)this[i]=Y[i];return!1}return!1},"test_match"),next:l(function(){if(this.done)return this.EOF;this._input||(this.done=!0);var a,k,h,u;this._more||(this.yytext="",this.match="");for(var Y=this._currentRules(),i=0;i<Y.length;i++)if(h=this._input.match(this.rules[Y[i]]),h&&(!k||h[0].length>k[0].length)){if(k=h,u=i,this.options.backtrack_lexer){if(a=this.test_match(h,Y[i]),a!==!1)return a;if(this._backtrack){k=!1;continue}else return!1}else if(!this.options.flex)break}return k?(a=this.test_match(k,Y[u]),a!==!1?a:!1):this._input===""?this.EOF:this.parseError("Lexical error on line "+(this.yylineno+1)+`. Unrecognized text.
`+this.showPosition(),{text:"",token:null,line:this.yylineno})},"next"),lex:l(function(){var k=this.next();return k||this.lex()},"lex"),begin:l(function(k){this.conditionStack.push(k)},"begin"),popState:l(function(){var k=this.conditionStack.length-1;return k>0?this.conditionStack.pop():this.conditionStack[0]},"popState"),_currentRules:l(function(){return this.conditionStack.length&&this.conditionStack[this.conditionStack.length-1]?this.conditions[this.conditionStack[this.conditionStack.length-1]].rules:this.conditions.INITIAL.rules},"_currentRules"),topState:l(function(k){return k=this.conditionStack.length-1-Math.abs(k||0),k>=0?this.conditionStack[k]:"INITIAL"},"topState"),pushState:l(function(k){this.begin(k)},"pushState"),stateStackSize:l(function(){return this.conditionStack.length},"stateStackSize"),options:{"case-insensitive":!0},performAction:l(function(k,h,u,Y){var i=Y;switch(u){case 0:return this.begin("open_directive"),"open_directive";break;case 1:return this.begin("acc_title"),31;break;case 2:return this.popState(),"acc_title_value";break;case 3:return this.begin("acc_descr"),33;break;case 4:return this.popState(),"acc_descr_value";break;case 5:this.begin("acc_descr_multiline");break;case 6:this.popState();break;case 7:return"acc_descr_multiline_value";case 8:break;case 9:break;case 10:break;case 11:return 10;case 12:break;case 13:break;case 14:this.begin("href");break;case 15:this.popState();break;case 16:return 43;case 17:this.begin("callbackname");break;case 18:this.popState();break;case 19:this.popState(),this.begin("callbackargs");break;case 20:return 41;case 21:this.popState();break;case 22:return 42;case 23:this.begin("click");break;case 24:this.popState();break;case 25:return 40;case 26:return 4;case 27:return 22;case 28:return 23;case 29:return 24;case 30:return 25;case 31:return 26;case 32:return 28;case 33:return 27;case 34:return 29;case 35:return 12;case 36:return 13;case 37:return 14;case 38:return 15;case 39:return 16;case 40:return 17;case 41:return 18;case 42:return 20;case 43:return 21;case 44:return"date";case 45:return 30;case 46:return"accDescription";case 47:return 36;case 48:return 38;case 49:return 39;case 50:return":";case 51:return 6;case 52:return"INVALID"}},"anonymous"),rules:[/^(?:%%\{)/i,/^(?:accTitle\s*:\s*)/i,/^(?:(?!\n||)*[^\n]*)/i,/^(?:accDescr\s*:\s*)/i,/^(?:(?!\n||)*[^\n]*)/i,/^(?:accDescr\s*\{\s*)/i,/^(?:[\}])/i,/^(?:[^\}]*)/i,/^(?:%%(?!\{)*[^\n]*)/i,/^(?:[^\}]%%*[^\n]*)/i,/^(?:%%*[^\n]*[\n]*)/i,/^(?:[\n]+)/i,/^(?:\s+)/i,/^(?:%[^\n]*)/i,/^(?:href[\s]+["])/i,/^(?:["])/i,/^(?:[^"]*)/i,/^(?:call[\s]+)/i,/^(?:\([\s]*\))/i,/^(?:\()/i,/^(?:[^(]*)/i,/^(?:\))/i,/^(?:[^)]*)/i,/^(?:click[\s]+)/i,/^(?:[\s\n])/i,/^(?:[^\s\n]*)/i,/^(?:gantt\b)/i,/^(?:dateFormat\s[^#\n;]+)/i,/^(?:inclusiveEndDates\b)/i,/^(?:topAxis\b)/i,/^(?:axisFormat\s[^#\n;]+)/i,/^(?:tickInterval\s[^#\n;]+)/i,/^(?:includes\s[^#\n;]+)/i,/^(?:excludes\s[^#\n;]+)/i,/^(?:todayMarker\s[^\n;]+)/i,/^(?:weekday\s+monday\b)/i,/^(?:weekday\s+tuesday\b)/i,/^(?:weekday\s+wednesday\b)/i,/^(?:weekday\s+thursday\b)/i,/^(?:weekday\s+friday\b)/i,/^(?:weekday\s+saturday\b)/i,/^(?:weekday\s+sunday\b)/i,/^(?:weekend\s+friday\b)/i,/^(?:weekend\s+saturday\b)/i,/^(?:\d\d\d\d-\d\d-\d\d\b)/i,/^(?:title\s[^\n]+)/i,/^(?:accDescription\s[^#\n;]+)/i,/^(?:section\s[^\n]+)/i,/^(?:[^:\n]+)/i,/^(?::[^#\n;]+)/i,/^(?::)/i,/^(?:$)/i,/^(?:.)/i],conditions:{acc_descr_multiline:{rules:[6,7],inclusive:!1},acc_descr:{rules:[4],inclusive:!1},acc_title:{rules:[2],inclusive:!1},callbackargs:{rules:[21,22],inclusive:!1},callbackname:{rules:[18,19,20],inclusive:!1},href:{rules:[15,16],inclusive:!1},click:{rules:[24,25],inclusive:!1},INITIAL:{rules:[0,1,3,5,8,9,10,11,12,13,14,17,23,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52],inclusive:!0}}};return T})();m.lexer=v;function g(){this.yy={}}return l(g,"Parser"),g.prototype=m,m.Parser=g,new g})();Zt.parser=Zt;var fi=Zt;Q.default.extend(Ge.default);Q.default.extend(Xe.default);Q.default.extend(Ue.default);var Ne={friday:5,saturday:6},J="",te="",ee=void 0,ie="",yt=[],gt=[],re=new Map,ne=[],Mt=[],pt="",se="",qe=["active","done","crit","milestone","vert"],ae=[],ht="",Tt=!1,oe=!1,ce="sunday",Et="saturday",Qt=0,hi=l(function(){ne=[],Mt=[],pt="",ae=[],St=0,Jt=void 0,Ct=void 0,U=[],J="",te="",se="",ee=void 0,ie="",yt=[],gt=[],Tt=!1,oe=!1,Qt=0,re=new Map,ht="",me(),ce="sunday",Et="saturday"},"clear"),mi=l(function(t){ht=t},"setDiagramId"),ki=l(function(t){te=t},"setAxisFormat"),yi=l(function(){return te},"getAxisFormat"),gi=l(function(t){ee=t},"setTickInterval"),pi=l(function(){return ee},"getTickInterval"),vi=l(function(t){ie=t},"setTodayMarker"),Ti=l(function(){return ie},"getTodayMarker"),bi=l(function(t){J=t},"setDateFormat"),xi=l(function(){Tt=!0},"enableInclusiveEndDates"),wi=l(function(){return Tt},"endDatesAreInclusive"),_i=l(function(){oe=!0},"enableTopAxis"),Di=l(function(){return oe},"topAxisEnabled"),Si=l(function(t){se=t},"setDisplayMode"),Ci=l(function(){return se},"getDisplayMode"),Mi=l(function(){return J},"getDateFormat"),Ze=l((t,e)=>{let s=e.toLowerCase().split(/[\s,]+/).filter(r=>r!=="");return[...new Set([...t,...s])]},"mergeTokens"),Ei=l(function(t){yt=Ze(yt,t)},"setIncludes"),Ii=l(function(){return yt},"getIncludes"),Yi=l(function(t){gt=Ze(gt,t)},"setExcludes"),Ai=l(function(){return gt},"getExcludes"),$i=l(function(){return re},"getLinks"),Li=l(function(t){pt=t,ne.push(t)},"addSection"),Fi=l(function(){return ne},"getSections"),Oi=l(function(){let t=He(),e=10,s=0;for(;!t&&s<e;)t=He(),s++;return Mt=U,Mt},"getTasks"),Qe=l(function(t,e,s,r){let n=t.format(e.trim()),d=t.format("YYYY-MM-DD");return r.includes(n)||r.includes(d)?!1:s.includes("weekends")&&(t.isoWeekday()===Ne[Et]||t.isoWeekday()===Ne[Et]+1)||s.includes(t.format("dddd").toLowerCase())?!0:s.includes(n)||s.includes(d)},"isInvalidDate"),Wi=l(function(t){ce=t},"setWeekday"),Vi=l(function(){return ce},"getWeekday"),Pi=l(function(t){Et=t},"setWeekend"),Ke=l(function(t,e,s,r){if(!s.length||t.manualEndTime)return;let n;t.startTime instanceof Date?n=(0,Q.default)(t.startTime):n=(0,Q.default)(t.startTime,e,!0),n=n.add(1,"d");let d;t.endTime instanceof Date?d=(0,Q.default)(t.endTime):d=(0,Q.default)(t.endTime,e,!0);let[y,w]=zi(n,d,e,s,r);t.endTime=y.toDate(),t.renderEndTime=w},"checkTaskDates"),zi=l(function(t,e,s,r,n){let d=!1,y=null,w=e.add(1e4,"d");for(;t<=e;){if(d||(y=e.toDate()),d=Qe(t,s,r,n),d&&(e=e.add(1,"d"),e>w))throw new Error("Failed to find a valid date that was not excluded by `excludes` after 10,000 iterations.");t=t.add(1,"d")}return[e,y]},"fixTaskDates"),Kt=l(function(t,e,s){if(s=s.trim(),l(w=>{let x=w.trim();return x==="x"||x==="X"},"isTimestampFormat")(e)&&/^\d+$/.test(s))return new Date(Number(s));let d=/^after\s+(?<ids>[\d\w- ]+)/.exec(s);if(d!==null){let w=null;for(let V of d.groups.ids.split(" ")){let I=ct(V);I!==void 0&&(!w||I.endTime>w.endTime)&&(w=I)}if(w)return w.endTime;let x=new Date;return x.setHours(0,0,0,0),x}let y=(0,Q.default)(s,e.trim(),!0);if(y.isValid())return y.toDate();{nt.debug("Invalid date:"+s),nt.debug("With date format:"+e.trim());let w=new Date(s);if(w===void 0||isNaN(w.getTime())||w.getFullYear()<-1e4||w.getFullYear()>1e4)throw new Error("Invalid date:"+s);return w}},"getStartDate"),Je=l(function(t){let e=/^(\d+(?:\.\d+)?)([Mdhmswy]|ms)$/.exec(t.trim());return e!==null?[Number.parseFloat(e[1]),e[2]]:[NaN,"ms"]},"parseDuration"),ti=l(function(t,e,s,r=!1){s=s.trim();let d=/^until\s+(?<ids>[\d\w- ]+)/.exec(s);if(d!==null){let I=null;for(let M of d.groups.ids.split(" ")){let D=ct(M);D!==void 0&&(!I||D.startTime<I.startTime)&&(I=D)}if(I)return I.startTime;let H=new Date;return H.setHours(0,0,0,0),H}let y=(0,Q.default)(s,e.trim(),!0);if(y.isValid())return r&&(y=y.add(1,"d")),y.toDate();let w=(0,Q.default)(t),[x,V]=Je(s);if(!Number.isNaN(x)){let I=w.add(x,V);I.isValid()&&(w=I)}return w.toDate()},"getEndDate"),St=0,kt=l(function(t){return t===void 0?(St=St+1,"task"+St):t},"parseId"),Ri=l(function(t,e){let s;e.substr(0,1)===":"?s=e.substr(1,e.length):s=e;let r=s.split(","),n={};le(r,n,qe);for(let y=0;y<r.length;y++)r[y]=r[y].trim();let d="";switch(r.length){case 1:n.id=kt(),n.startTime=t.endTime,d=r[0];break;case 2:n.id=kt(),n.startTime=Kt(void 0,J,r[0]),d=r[1];break;case 3:n.id=kt(r[0]),n.startTime=Kt(void 0,J,r[1]),d=r[2];break;default:}return d&&(n.endTime=ti(n.startTime,J,d,Tt),n.manualEndTime=(0,Q.default)(d,"YYYY-MM-DD",!0).isValid(),Ke(n,J,gt,yt)),n},"compileData"),Ni=l(function(t,e){let s;e.substr(0,1)===":"?s=e.substr(1,e.length):s=e;let r=s.split(","),n={};le(r,n,qe);for(let d=0;d<r.length;d++)r[d]=r[d].trim();switch(r.length){case 1:n.id=kt(),n.startTime={type:"prevTaskEnd",id:t},n.endTime={data:r[0]};break;case 2:n.id=kt(),n.startTime={type:"getStartDate",startData:r[0]},n.endTime={data:r[1]};break;case 3:n.id=kt(r[0]),n.startTime={type:"getStartDate",startData:r[1]},n.endTime={data:r[2]};break;default:}return n},"parseData"),Jt,Ct,U=[],ei={},Hi=l(function(t,e){let s={section:pt,type:pt,processed:!1,manualEndTime:!1,renderEndTime:null,raw:{data:e},task:t,classes:[]},r=Ni(Ct,e);s.raw.startTime=r.startTime,s.raw.endTime=r.endTime,s.id=r.id,s.prevTaskId=Ct,s.active=r.active,s.done=r.done,s.crit=r.crit,s.milestone=r.milestone,s.vert=r.vert,s.vert?s.order=-1:(s.order=Qt,Qt++);let n=U.push(s);Ct=s.id,ei[s.id]=n-1},"addTask"),ct=l(function(t){let e=ei[t];return U[e]},"findTaskById"),Bi=l(function(t,e){let s={section:pt,type:pt,description:t,task:t,classes:[]},r=Ri(Jt,e);s.startTime=r.startTime,s.endTime=r.endTime,s.id=r.id,s.active=r.active,s.done=r.done,s.crit=r.crit,s.milestone=r.milestone,s.vert=r.vert,Jt=s,Mt.push(s)},"addTaskOrg"),He=l(function(){let t=l(function(s){let r=U[s],n="";switch(U[s].raw.startTime.type){case"prevTaskEnd":{let d=ct(r.prevTaskId);r.startTime=d.endTime;break}case"getStartDate":n=Kt(void 0,J,U[s].raw.startTime.startData),n&&(U[s].startTime=n);break}return U[s].startTime&&(U[s].endTime=ti(U[s].startTime,J,U[s].raw.endTime.data,Tt),U[s].endTime&&(U[s].processed=!0,U[s].manualEndTime=(0,Q.default)(U[s].raw.endTime.data,"YYYY-MM-DD",!0).isValid(),Ke(U[s],J,gt,yt))),U[s].processed},"compileTask"),e=!0;for(let[s,r]of U.entries())t(s),e=e&&r.processed;return e},"compileTasks"),ji=l(function(t,e){let s=e;ot().securityLevel!=="loose"&&(s=(0,je.sanitizeUrl)(e)),t.split(",").forEach(function(r){ct(r)!==void 0&&(ri(r,()=>{window.open(s,"_self")}),re.set(r,s))}),ii(t,"clickable")},"setLink"),ii=l(function(t,e){t.split(",").forEach(function(s){let r=ct(s);r!==void 0&&r.classes.push(e)})},"setClass"),Gi=l(function(t,e,s){if(ot().securityLevel!=="loose"||e===void 0)return;let r=[];if(typeof s=="string"){r=s.split(/,(?=(?:(?:[^"]*"){2})*[^"]*$)/);for(let d=0;d<r.length;d++){let y=r[d].trim();y.startsWith('"')&&y.endsWith('"')&&(y=y.substr(1,y.length-2)),r[d]=y}}r.length===0&&r.push(t),ct(t)!==void 0&&ri(t,()=>{$e.runFunc(e,...r)})},"setClickFun"),ri=l(function(t,e){ae.push(function(){let s=ht?`${ht}-${t}`:t,r=document.querySelector(`[id="${s}"]`);r!==null&&r.addEventListener("click",function(){e()})},function(){let s=ht?`${ht}-${t}`:t,r=document.querySelector(`[id="${s}-text"]`);r!==null&&r.addEventListener("click",function(){e()})})},"pushFun"),Xi=l(function(t,e,s){t.split(",").forEach(function(r){Gi(r,e,s)}),ii(t,"clickable")},"setClickEvent"),Ui=l(function(t){ae.forEach(function(e){e(t)})},"bindFunctions"),qi={getConfig:l(()=>ot().gantt,"getConfig"),clear:hi,setDateFormat:bi,getDateFormat:Mi,enableInclusiveEndDates:xi,endDatesAreInclusive:wi,enableTopAxis:_i,topAxisEnabled:Di,setAxisFormat:ki,getAxisFormat:yi,setTickInterval:gi,getTickInterval:pi,setTodayMarker:vi,getTodayMarker:Ti,setAccTitle:ke,getAccTitle:ye,setDiagramTitle:ve,getDiagramTitle:Te,setDiagramId:mi,setDisplayMode:Si,getDisplayMode:Ci,setAccDescription:ge,getAccDescription:pe,addSection:Li,getSections:Fi,getTasks:Oi,addTask:Hi,findTaskById:ct,addTaskOrg:Bi,setIncludes:Ei,getIncludes:Ii,setExcludes:Yi,getExcludes:Ai,setClickEvent:Xi,setLink:ji,getLinks:$i,bindFunctions:Ui,parseDuration:Je,isInvalidDate:Qe,setWeekday:Wi,getWeekday:Vi,setWeekend:Pi};function le(t,e,s){let r=!0;for(;r;)r=!1,s.forEach(function(n){let d="^\\s*"+n+"\\s*$",y=new RegExp(d);t[0].match(y)&&(e[n]=!0,t.shift(1),r=!0)})}l(le,"getTaskTags");mt.default.extend(ni.default);var Zi=l(function(){nt.debug("Something is calling, setConf, remove the call")},"setConf"),Be={monday:De,tuesday:Se,wednesday:Ce,thursday:Me,friday:Ee,saturday:Ie,sunday:_e},Qi=l((t,e)=>{let s=[...t].map(()=>-1/0),r=[...t].sort((d,y)=>d.startTime-y.startTime||d.order-y.order),n=0;for(let d of r)for(let y=0;y<s.length;y++)if(d.startTime>=s[y]){s[y]=d.endTime,d.order=y+e,y>n&&(n=y);break}return n},"getMaxIntersections"),rt,qt=1e4,Ki=l(function(t,e,s,r){let n=ot().gantt;r.db.setDiagramId(e);let d=ot().securityLevel,y;d==="sandbox"&&(y=vt("#i"+e));let w=d==="sandbox"?vt(y.nodes()[0].contentDocument.body):vt("body"),x=d==="sandbox"?y.nodes()[0].contentDocument:document,V=x.getElementById(e);rt=V.parentElement.offsetWidth,rt===void 0&&(rt=1200),n.useWidth!==void 0&&(rt=n.useWidth);let I=r.db.getTasks(),H=I.filter(m=>!m.vert),M=[];for(let m of H)M.push(m.type);M=F(M);let D={},_=2*n.topPadding;if(r.db.getDisplayMode()==="compact"||n.displayMode==="compact"){let m={};for(let g of H)m[g.section]===void 0?m[g.section]=[g]:m[g.section].push(g);let v=0;for(let g of Object.keys(m)){let T=Qi(m[g],v)+1;v+=T,_+=T*(n.barHeight+n.barGap),D[g]=T}}else{_+=H.length*(n.barHeight+n.barGap);for(let m of M)D[m]=H.filter(v=>v.type===m).length}V.setAttribute("viewBox","0 0 "+rt+" "+_);let S=w.select(`[id="${e}"]`),b=Ae().domain([xe(I,function(m){return m.startTime}),be(I,function(m){return m.endTime})]).rangeRound([0,rt-n.leftPadding-n.rightPadding]);function O(m,v){let g=m.startTime,T=v.startTime,a=0;return g>T?a=1:g<T&&(a=-1),a}l(O,"taskCompare"),I.sort(O),B(I,rt,_),he(S,_,rt,n.useMaxWidth),S.append("text").text(r.db.getDiagramTitle()).attr("x",rt/2).attr("y",n.titleTopMargin).attr("class","titleText");function B(m,v,g){let T=n.barHeight,a=T+n.barGap,k=n.topPadding,h=n.leftPadding,u=Ye().domain([0,M.length]).range(["#00B9FA","#F95002"]).interpolate(we);L(a,k,h,v,g,m,r.db.getExcludes(),r.db.getIncludes()),f(h,k,v,g),j(m,a,k,h,T,u,v,g),C(a,k,h,T,u),W(h,k,v,g)}l(B,"makeGantt");function j(m,v,g,T,a,k,h){m.sort((o,E)=>o.vert===E.vert?0:o.vert?1:-1);let u=m.filter(o=>!o.vert),i=[...new Set(u.map(o=>o.order))].map(o=>u.find(E=>E.order===o));S.append("g").selectAll("rect").data(i).enter().append("rect").attr("x",0).attr("y",function(o,E){return E=o.order,E*v+g-2}).attr("width",function(){return h-n.rightPadding/2}).attr("height",v).attr("class",function(o){for(let[E,$]of M.entries())if(o.type===$)return"section section"+E%n.numberSectionStyles;return"section section0"}).enter();let A=S.append("g").selectAll("rect").data(m).enter(),c=r.db.getLinks();if(A.append("rect").attr("id",function(o){return e+"-"+o.id}).attr("rx",3).attr("ry",3).attr("x",function(o){return o.milestone?b(o.startTime)+T+.5*(b(o.endTime)-b(o.startTime))-.5*a:b(o.startTime)+T}).attr("y",function(o,E){return E=o.order,o.vert?n.gridLineStartPadding:E*v+g}).attr("width",function(o){return o.milestone?a:o.vert?.08*a:b(o.renderEndTime||o.endTime)-b(o.startTime)}).attr("height",function(o){return o.vert?u.length*(n.barHeight+n.barGap)+n.barHeight*2:a}).attr("transform-origin",function(o,E){return E=o.order,(b(o.startTime)+T+.5*(b(o.endTime)-b(o.startTime))).toString()+"px "+(E*v+g+.5*a).toString()+"px"}).attr("class",function(o){let E="task",$="";o.classes.length>0&&($=o.classes.join(" "));let R=0;for(let[P,z]of M.entries())o.type===z&&(R=P%n.numberSectionStyles);let N="";return o.active?o.crit?N+=" activeCrit":N=" active":o.done?o.crit?N=" doneCrit":N=" done":o.crit&&(N+=" crit"),N.length===0&&(N=" task"),o.milestone&&(N=" milestone "+N),o.vert&&(N=" vert "+N),N+=R,N+=" "+$,E+N}),A.append("text").attr("id",function(o){return e+"-"+o.id+"-text"}).text(function(o){return o.task}).attr("font-size",n.fontSize).attr("x",function(o){let E=b(o.startTime),$=b(o.renderEndTime||o.endTime);if(o.milestone&&(E+=.5*(b(o.endTime)-b(o.startTime))-.5*a,$=E+a),o.vert)return b(o.startTime)+T;let R=this.getBBox().width;return R>$-E?$+R+1.5*n.leftPadding>h?E+T-5:$+T+5:($-E)/2+E+T}).attr("y",function(o,E){return o.vert?n.gridLineStartPadding+u.length*(n.barHeight+n.barGap)+60:(E=o.order,E*v+n.barHeight/2+(n.fontSize/2-2)+g)}).attr("text-height",a).attr("class",function(o){let E=b(o.startTime),$=b(o.endTime);o.milestone&&($=E+a);let R=this.getBBox().width,N="";o.classes.length>0&&(N=o.classes.join(" "));let P=0;for(let[tt,et]of M.entries())o.type===et&&(P=tt%n.numberSectionStyles);let z="";return o.active&&(o.crit?z="activeCritText"+P:z="activeText"+P),o.done?o.crit?z=z+" doneCritText"+P:z=z+" doneText"+P:o.crit&&(z=z+" critText"+P),o.milestone&&(z+=" milestoneText"),o.vert&&(z+=" vertText"),R>$-E?$+R+1.5*n.leftPadding>h?N+" taskTextOutsideLeft taskTextOutside"+P+" "+z:N+" taskTextOutsideRight taskTextOutside"+P+" "+z+" width-"+R:N+" taskText taskText"+P+" "+z+" width-"+R}),ot().securityLevel==="sandbox"){let o;o=vt("#i"+e);let E=o.nodes()[0].contentDocument;A.filter(function($){return c.has($.id)}).each(function($){var R=E.querySelector("#"+CSS.escape(e+"-"+$.id)),N=E.querySelector("#"+CSS.escape(e+"-"+$.id+"-text"));let P=R.parentNode;var z=E.createElement("a");z.setAttribute("xlink:href",c.get($.id)),z.setAttribute("target","_top"),P.appendChild(z),z.appendChild(R),z.appendChild(N)})}}l(j,"drawRects");function L(m,v,g,T,a,k,h,u){if(h.length===0&&u.length===0)return;let Y,i;for(let{startTime:$,endTime:R}of k)(Y===void 0||$<Y)&&(Y=$),(i===void 0||R>i)&&(i=R);if(!Y||!i)return;if((0,mt.default)(i).diff((0,mt.default)(Y),"year")>5){nt.warn("The difference between the min and max time is more than 5 years. This will cause performance issues. Skipping drawing exclude days.");return}let A=r.db.getDateFormat(),c=[],G=null,o=(0,mt.default)(Y);for(;o.valueOf()<=i;)r.db.isInvalidDate(o,A,h,u)?G?G.end=o:G={start:o,end:o}:G&&(c.push(G),G=null),o=o.add(1,"d");S.append("g").selectAll("rect").data(c).enter().append("rect").attr("id",$=>e+"-exclude-"+$.start.format("YYYY-MM-DD")).attr("x",$=>b($.start.startOf("day"))+g).attr("y",n.gridLineStartPadding).attr("width",$=>b($.end.endOf("day"))-b($.start.startOf("day"))).attr("height",a-v-n.gridLineStartPadding).attr("transform-origin",function($,R){return(b($.start)+g+.5*(b($.end)-b($.start))).toString()+"px "+(R*m+.5*a).toString()+"px"}).attr("class","exclude-range")}l(L,"drawExcludeDays");function p(m,v,g,T){if(g<=0||m>v)return 1/0;let a=v-m,k=mt.default.duration({[T??"day"]:g}).asMilliseconds();return k<=0?1/0:Math.ceil(a/k)}l(p,"getEstimatedTickCount");function f(m,v,g,T){let a=r.db.getDateFormat(),k=r.db.getAxisFormat(),h;k?h=k:a==="D"?h="%d":h=n.axisFormat??"%Y-%m-%d";let u=We(b).tickSize(-T+v+n.gridLineStartPadding).tickFormat(Vt(h)),i=/^([1-9]\d*)(millisecond|second|minute|hour|day|week|month)$/.exec(r.db.getTickInterval()||n.tickInterval);if(i!==null){let A=parseInt(i[1],10);if(isNaN(A)||A<=0)nt.warn(`Invalid tick interval value: "${i[1]}". Skipping custom tick interval.`);else{let c=i[2],G=r.db.getWeekday()||n.weekday,o=b.domain(),E=o[0],$=o[1],R=p(E,$,A,c);if(R>qt)nt.warn(`The tick interval "${A}${c}" would generate ${R} ticks, which exceeds the maximum allowed (${qt}). This may indicate an invalid date or time range. Skipping custom tick interval.`);else switch(c){case"millisecond":u.ticks(At.every(A));break;case"second":u.ticks($t.every(A));break;case"minute":u.ticks(Lt.every(A));break;case"hour":u.ticks(Ft.every(A));break;case"day":u.ticks(Ot.every(A));break;case"week":u.ticks(Be[G].every(A));break;case"month":u.ticks(Wt.every(A));break}}}if(S.append("g").attr("class","grid").attr("transform","translate("+m+", "+(T-50)+")").call(u).selectAll("text").style("text-anchor","middle").attr("fill","#000").attr("stroke","none").attr("font-size",10).attr("dy","1em"),r.db.topAxisEnabled()||n.topAxis){let A=Oe(b).tickSize(-T+v+n.gridLineStartPadding).tickFormat(Vt(h));if(i!==null){let c=parseInt(i[1],10);if(isNaN(c)||c<=0)nt.warn(`Invalid tick interval value: "${i[1]}". Skipping custom tick interval.`);else{let G=i[2],o=r.db.getWeekday()||n.weekday,E=b.domain(),$=E[0],R=E[1];if(p($,R,c,G)<=qt)switch(G){case"millisecond":A.ticks(At.every(c));break;case"second":A.ticks($t.every(c));break;case"minute":A.ticks(Lt.every(c));break;case"hour":A.ticks(Ft.every(c));break;case"day":A.ticks(Ot.every(c));break;case"week":A.ticks(Be[o].every(c));break;case"month":A.ticks(Wt.every(c));break}}}S.append("g").attr("class","grid").attr("transform","translate("+m+", "+v+")").call(A).selectAll("text").style("text-anchor","middle").attr("fill","#000").attr("stroke","none").attr("font-size",10)}}l(f,"makeGrid");function C(m,v){let g=0,T=Object.keys(D).map(a=>[a,D[a]]);S.append("g").selectAll("text").data(T).enter().append(function(a){let k=a[0].split(fe.lineBreakRegex),h=-(k.length-1)/2,u=x.createElementNS("http://www.w3.org/2000/svg","text");u.setAttribute("dy",h+"em");for(let[Y,i]of k.entries()){let A=x.createElementNS("http://www.w3.org/2000/svg","tspan");A.setAttribute("alignment-baseline","central"),A.setAttribute("x","10"),Y>0&&A.setAttribute("dy","1em"),A.textContent=i,u.appendChild(A)}return u}).attr("x",10).attr("y",function(a,k){if(k>0)for(let h=0;h<k;h++)return g+=T[k-1][1],a[1]*m/2+g*m+v;else return a[1]*m/2+v}).attr("font-size",n.sectionFontSize).attr("class",function(a){for(let[k,h]of M.entries())if(a[0]===h)return"sectionTitle sectionTitle"+k%n.numberSectionStyles;return"sectionTitle"})}l(C,"vertLabels");function W(m,v,g,T){let a=r.db.getTodayMarker();if(a==="off")return;let k=S.append("g").attr("class","today"),h=new Date,u=k.append("line");u.attr("x1",b(h)+m).attr("x2",b(h)+m).attr("y1",n.titleTopMargin).attr("y2",T-n.titleTopMargin).attr("class","today"),a!==""&&u.attr("style",a.replace(/,/g,";"))}l(W,"drawToday");function F(m){let v={},g=[];for(let T=0,a=m.length;T<a;++T)Object.prototype.hasOwnProperty.call(v,m[T])||(v[m[T]]=!0,g.push(m[T]));return g}l(F,"checkUnique")},"draw"),Ji={setConf:Zi,draw:Ki},tr=l(t=>`
  .mermaid-main-font {
        font-family: ${t.fontFamily};
  }

  .exclude-range {
    fill: ${t.excludeBkgColor};
  }

  .section {
    stroke: none;
    opacity: 0.2;
  }

  .section0 {
    fill: ${t.sectionBkgColor};
  }

  .section2 {
    fill: ${t.sectionBkgColor2};
  }

  .section1,
  .section3 {
    fill: ${t.altSectionBkgColor};
    opacity: 0.2;
  }

  .sectionTitle0 {
    fill: ${t.titleColor};
  }

  .sectionTitle1 {
    fill: ${t.titleColor};
  }

  .sectionTitle2 {
    fill: ${t.titleColor};
  }

  .sectionTitle3 {
    fill: ${t.titleColor};
  }

  .sectionTitle {
    text-anchor: start;
    font-family: ${t.fontFamily};
  }


  /* Grid and axis */

  .grid .tick {
    stroke: ${t.gridColor};
    opacity: 0.8;
    shape-rendering: crispEdges;
  }

  .grid .tick text {
    font-family: ${t.fontFamily};
    fill: ${t.textColor};
  }

  .grid path {
    stroke-width: 0;
  }


  /* Today line */

  .today {
    fill: none;
    stroke: ${t.todayLineColor};
    stroke-width: 2px;
  }


  /* Task styling */

  /* Default task */

  .task {
    stroke-width: 2;
  }

  .taskText {
    text-anchor: middle;
    font-family: ${t.fontFamily};
  }

  .taskTextOutsideRight {
    fill: ${t.taskTextDarkColor};
    text-anchor: start;
    font-family: ${t.fontFamily};
  }

  .taskTextOutsideLeft {
    fill: ${t.taskTextDarkColor};
    text-anchor: end;
  }


  /* Special case clickable */

  .task.clickable {
    cursor: pointer;
  }

  .taskText.clickable {
    cursor: pointer;
    fill: ${t.taskTextClickableColor} !important;
    font-weight: bold;
  }

  .taskTextOutsideLeft.clickable {
    cursor: pointer;
    fill: ${t.taskTextClickableColor} !important;
    font-weight: bold;
  }

  .taskTextOutsideRight.clickable {
    cursor: pointer;
    fill: ${t.taskTextClickableColor} !important;
    font-weight: bold;
  }


  /* Specific task settings for the sections*/

  .taskText0,
  .taskText1,
  .taskText2,
  .taskText3 {
    fill: ${t.taskTextColor};
  }

  .task0,
  .task1,
  .task2,
  .task3 {
    fill: ${t.taskBkgColor};
    stroke: ${t.taskBorderColor};
  }

  .taskTextOutside0,
  .taskTextOutside2
  {
    fill: ${t.taskTextOutsideColor};
  }

  .taskTextOutside1,
  .taskTextOutside3 {
    fill: ${t.taskTextOutsideColor};
  }


  /* Active task */

  .active0,
  .active1,
  .active2,
  .active3 {
    fill: ${t.activeTaskBkgColor};
    stroke: ${t.activeTaskBorderColor};
  }

  .activeText0,
  .activeText1,
  .activeText2,
  .activeText3 {
    fill: ${t.taskTextDarkColor} !important;
  }


  /* Completed task */

  .done0,
  .done1,
  .done2,
  .done3 {
    stroke: ${t.doneTaskBorderColor};
    fill: ${t.doneTaskBkgColor};
    stroke-width: 2;
  }

  .doneText0,
  .doneText1,
  .doneText2,
  .doneText3 {
    fill: ${t.taskTextDarkColor} !important;
  }

  /* Done task text displayed outside the bar sits against the diagram background,
     not against the done-task bar, so it must use the outside/contrast color. */
  .doneText0.taskTextOutsideLeft,
  .doneText0.taskTextOutsideRight,
  .doneText1.taskTextOutsideLeft,
  .doneText1.taskTextOutsideRight,
  .doneText2.taskTextOutsideLeft,
  .doneText2.taskTextOutsideRight,
  .doneText3.taskTextOutsideLeft,
  .doneText3.taskTextOutsideRight {
    fill: ${t.taskTextOutsideColor} !important;
  }


  /* Tasks on the critical line */

  .crit0,
  .crit1,
  .crit2,
  .crit3 {
    stroke: ${t.critBorderColor};
    fill: ${t.critBkgColor};
    stroke-width: 2;
  }

  .activeCrit0,
  .activeCrit1,
  .activeCrit2,
  .activeCrit3 {
    stroke: ${t.critBorderColor};
    fill: ${t.activeTaskBkgColor};
    stroke-width: 2;
  }

  .doneCrit0,
  .doneCrit1,
  .doneCrit2,
  .doneCrit3 {
    stroke: ${t.critBorderColor};
    fill: ${t.doneTaskBkgColor};
    stroke-width: 2;
    cursor: pointer;
    shape-rendering: crispEdges;
  }

  .milestone {
    transform: rotate(45deg) scale(0.8,0.8);
  }

  .milestoneText {
    font-style: italic;
  }
  .doneCritText0,
  .doneCritText1,
  .doneCritText2,
  .doneCritText3 {
    fill: ${t.taskTextDarkColor} !important;
  }

  /* Done-crit task text outside the bar \u2014 same reasoning as doneText above. */
  .doneCritText0.taskTextOutsideLeft,
  .doneCritText0.taskTextOutsideRight,
  .doneCritText1.taskTextOutsideLeft,
  .doneCritText1.taskTextOutsideRight,
  .doneCritText2.taskTextOutsideLeft,
  .doneCritText2.taskTextOutsideRight,
  .doneCritText3.taskTextOutsideLeft,
  .doneCritText3.taskTextOutsideRight {
    fill: ${t.taskTextOutsideColor} !important;
  }

  .vert {
    stroke: ${t.vertLineColor};
  }

  .vertText {
    font-size: 15px;
    text-anchor: middle;
    fill: ${t.vertLineColor} !important;
  }

  .activeCritText0,
  .activeCritText1,
  .activeCritText2,
  .activeCritText3 {
    fill: ${t.taskTextDarkColor} !important;
  }

  .titleText {
    text-anchor: middle;
    font-size: 18px;
    fill: ${t.titleColor||t.textColor};
    font-family: ${t.fontFamily};
  }
`,"getStyles"),er=tr,lr={parser:fi,db:qi,renderer:Ji,styles:er};export{lr as diagram};
