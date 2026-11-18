import{a as ot}from"./chunk-QRO5TZGF.js";import{a as at}from"./chunk-BHZAHZCC.js";import{a as nt}from"./chunk-UGBNWPB5.js";import"./chunk-T2647YD2.js";import"./chunk-V2IKIZ2E.js";import"./chunk-DGL6ZPDT.js";import{w as tt}from"./chunk-MKCXTBEZ.js";import"./chunk-Q3RTYOM3.js";import"./chunk-H7ARKS5I.js";import"./chunk-DOMQEUEU.js";import"./chunk-DHIDQXOH.js";import"./chunk-5YUNXXPO.js";import"./chunk-YUPO6TJ5.js";import"./chunk-I7Y7QKED.js";import"./chunk-BEAVMZXW.js";import"./chunk-5HCT2Y7B.js";import"./chunk-2S364EUL.js";import"./chunk-47X3M3DT.js";import"./chunk-YGNX3VI3.js";import"./chunk-W3Q3274G.js";import{n as rt,o as it}from"./chunk-7TIOW63O.js";import{a as R,d as et}from"./chunk-MNVP5QMN.js";import{O as V,T as X,U as Z,V as j,W as q,X as J,Y as K,Z as Q,_ as Y,k as U}from"./chunk-7JYBA3TO.js";import{b as T}from"./chunk-K6ALI33Z.js";import{a as l}from"./chunk-PAESYWDV.js";import"./chunk-5GVNIGXK.js";var lt=U.pie,L={sections:new Map,showData:!1,config:lt},b=L.sections,O=L.showData,St=structuredClone(lt),xt=l(()=>structuredClone(St),"getConfig"),wt=l(()=>{b=new Map,O=L.showData,X()},"clear"),Ct=l(({label:t,value:a})=>{if(a<0)throw new Error(`"${t}" has invalid value: ${a}. Negative values are not allowed in pie charts. All slice values must be >= 0.`);b.has(t)||(b.set(t,a),T.debug(`added new section: ${t}, with value: ${a}`))},"addSection"),$t=l(()=>b,"getSections"),Dt=l(t=>{O=t},"setShowData"),yt=l(()=>O,"getShowData"),st={getConfig:xt,clear:wt,setDiagramTitle:K,getDiagramTitle:Q,setAccTitle:Z,getAccTitle:j,setAccDescription:q,getAccDescription:J,addSection:Ct,getSections:$t,setShowData:Dt,getShowData:yt},Tt=l((t,a)=>{ot(t,a),a.setShowData(t.showData),t.sections.map(a.addSection)},"populateDb"),bt={parse:l(async t=>{let a=await nt("pie",t);T.debug(a),Tt(a,st)},"parse")},At=l(t=>`
  .pieCircle{
    stroke: ${t.pieStrokeColor};
    stroke-width : ${t.pieStrokeWidth};
    opacity : ${t.pieOpacity};
  }
  .pieCircle.highlighted{
    scale: 1.05;
    opacity: 1;
  }
  .pieCircle.highlightedOnHover:hover{
    transition-duration: 250ms;
    scale: 1.05;
    opacity: 1;
  }
  .pieOuterCircle{
    stroke: ${t.pieOuterStrokeColor};
    stroke-width: ${t.pieOuterStrokeWidth};
    fill: none;
  }
  .pieTitleText {
    text-anchor: middle;
    font-size: ${t.pieTitleTextSize};
    fill: ${t.pieTitleTextColor};
    font-family: ${t.fontFamily};
  }
  .slice {
    font-family: ${t.fontFamily};
    fill: ${t.pieSectionTextColor};
    font-size:${t.pieSectionTextSize};
    // fill: white;
  }
  .legend text {
    fill: ${t.pieLegendTextColor};
    font-family: ${t.fontFamily};
    font-size: ${t.pieLegendTextSize};
  }
`,"getStyles"),kt=At,_t=l(t=>{let a=[...t.values()].reduce((n,m)=>n+m,0),W=[...t.entries()].map(([n,m])=>({label:n,value:m})).filter(n=>n.value/a*100>=1);return et().value(n=>n.value).sort(null)(W)},"createPieArcs"),zt=l((t,a,W,F)=>{T.debug(`rendering pie chart
`+t);let n=F.db,m=Y(),h=it(n.getConfig(),m.pie),H=40,i=18,c=4,S=450,x=S,A=at(a),$=A.append("g");$.attr("transform","translate("+x/2+","+S/2+")");let{themeVariables:o}=m,[M]=rt(o.pieOuterStrokeWidth);M??=2;let ct=h.legendPosition,P=h.textPosition,dt=h.donutHole>0&&h.donutHole<=.9?h.donutHole:0,f=Math.min(x,S)/2-H,gt=R().innerRadius(dt*f).outerRadius(f),pt=R().innerRadius(f*P).outerRadius(f*P),w=$.append("g");w.append("circle").attr("cx",0).attr("cy",0).attr("r",f+M/2).attr("class","pieOuterCircle");let D=n.getSections(),ht=_t(D),ft=[o.pie1,o.pie2,o.pie3,o.pie4,o.pie5,o.pie6,o.pie7,o.pie8,o.pie9,o.pie10,o.pie11,o.pie12],k=0;D.forEach(e=>{k+=e});let G=ht.filter(e=>(e.data.value/k*100).toFixed(0)!=="0"),_=tt(ft).domain([...D.keys()]);w.selectAll("mySlices").data(G).enter().append("path").attr("d",gt).attr("fill",e=>_(e.data.label)).attr("class",e=>{let r="pieCircle";return h.highlightSlice==="hover"?r+=" highlightedOnHover":h.highlightSlice===e.data.label&&(r+=" highlighted"),r}),w.selectAll("mySlices").data(G).enter().append("text").text(e=>(e.data.value/k*100).toFixed(0)+"%").attr("transform",e=>"translate("+pt.centroid(e)+")").style("text-anchor","middle").attr("class","slice");let ut=$.append("text").text(n.getDiagramTitle()).attr("x",0).attr("y",-(S-50)/2).attr("class","pieTitleText"),C=[...D.entries()].map(([e,r])=>({label:e,value:r})),u=$.selectAll(".legend").data(C).enter().append("g").attr("class","legend");u.append("rect").attr("width",i).attr("height",i).style("fill",e=>_(e.label)).style("stroke",e=>_(e.label)),u.append("text").attr("x",i+c).attr("y",i-c).text(e=>n.getShowData()?`${e.label} [${e.value}]`:e.label);let v=Math.max(...u.selectAll("text").nodes().map(e=>e?.getBoundingClientRect().width??0)),y=S,z=x+H,s=i+c,E=C.length*s;switch(ct){case"center":u.attr("transform",(e,r)=>{let d=s*C.length/2,g=-v/2-(i+c),p=r*s-d;return"translate("+g+","+p+")"});break;case"top":y+=E,u.attr("transform",(e,r)=>{let d=f,g=-v/2-(i+c),p=r*s-d;return`translate(${g}, ${p})`}),w.attr("transform",()=>`translate(0, ${E+s})`);break;case"bottom":y+=E,u.attr("transform",(e,r)=>{let d=-f-s,g=-v/2-(i+c),p=r*s-d;return"translate("+g+","+p+")"});break;case"left":z+=i+c+v,u.attr("transform",(e,r)=>{let d=s*C.length/2,g=-f-(i+c),p=r*s-d;return"translate("+g+","+p+")"}),w.attr("transform",()=>`translate(${v+i+c}, 0)`);break;case"right":default:z+=i+c+v,u.attr("transform",(e,r)=>{let d=s*C.length/2,g=12*i,p=r*s-d;return"translate("+g+","+p+")"});break}let B=ut.node()?.getBoundingClientRect().width??0,mt=x/2-B/2,vt=x/2+B/2,N=Math.min(0,mt),I=Math.max(z,vt)-N;A.attr("viewBox",`${N} 0 ${I} ${y}`),V(A,y,I,h.useMaxWidth)},"draw"),Et={draw:zt},Bt={parser:bt,db:st,renderer:Et,styles:kt};export{Bt as diagram};
