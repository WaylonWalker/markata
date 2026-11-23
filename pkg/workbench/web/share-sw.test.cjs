'use strict';
const test=require('node:test');
const assert=require('node:assert/strict');
const fs=require('node:fs');
const path=require('node:path');
const vm=require('node:vm');

const source=fs.readFileSync(path.join(__dirname,'share-sw.js'),'utf8');

function workerHarness(offline) {
  const handlers=new Map();
  const cached=[];
  let skipWaitingCalls=0;
  const offlinePage=new Response('PUBLIC OFFLINE CAPTURE', {headers:{'Content-Type':'text/html'}});
  const caches={
    open:async name=>{
      assert.equal(name,'markata-offline-capture-v1');
      return {add:async pathname=>{
        cached.push(pathname);
        assert.equal(pathname,'/editor/offline-capture.html');
      }};
    },
    keys:async()=>['markata-offline-capture-v1'],
    delete:async()=>true,
    match:async pathname=>pathname==='/editor/offline-capture.html'?offlinePage.clone():undefined
  };
  const self={
    registration:{scope:'https://markata.example/'},
    location:{origin:'https://markata.example'},
    clients:{claim:async()=>{}},
    skipWaiting:async()=>{skipWaitingCalls++;},
    addEventListener:(name,fn)=>handlers.set(name,fn)
  };
  vm.runInNewContext(source,{self,URL,Response,caches,console,
    importScripts:()=>{},
    fetch:async()=>{
      if(offline)throw new Error('network unavailable');
      return new Response('PRIVATE EDITOR RESPONSE');
    }
  });
  async function dispatch(name,event={}) {
    let pending,responded;
    handlers.get(name)({...event,
      waitUntil(p){pending=p;},
      respondWith(p){responded=p;}
    });
    if(pending)await pending;
    return responded;
  }
  return {dispatch,cached,skips:()=>skipWaitingCalls};
}

test('offline navigation serves only a public capture page, never a cached editor',async()=>{
  const worker=workerHarness(true);
  await worker.dispatch('install');
  assert.deepEqual(worker.cached,['/editor/offline-capture.html']);
  assert.equal(worker.skips(),0,'upgrades must not take over an open editor');
  const response=await worker.dispatch('fetch',{
    request:{url:'https://markata.example/',method:'GET',mode:'navigate'}
  });
  assert.equal(await response.text(),'PUBLIC OFFLINE CAPTURE');
  assert.deepEqual(worker.cached,['/editor/offline-capture.html']);
});

test('online editor navigation stays network first and private APIs are not intercepted',async()=>{
  const worker=workerHarness(false);
  await worker.dispatch('install');
  const response=await worker.dispatch('fetch',{
    request:{url:'https://markata.example/',method:'GET',mode:'navigate'}
  });
  assert.equal(await response.text(),'PRIVATE EDITOR RESPONSE');
  const intercepted=await worker.dispatch('fetch',{
    request:{url:'https://markata.example/workspace/file',method:'GET',mode:'navigate'}
  });
  assert.equal(intercepted,undefined);
  assert.deepEqual(worker.cached,['/editor/offline-capture.html']);
});

test('service worker activates only on an explicit same-origin update request',async()=>{
  const worker=workerHarness(false);
  await worker.dispatch('install');
  await worker.dispatch('message',{data:{type:'MARKATA_ACTIVATE_UPDATE'},
    source:{type:'window',url:'https://other.example/'}
  });
  assert.equal(worker.skips(),0);
  await worker.dispatch('message',{data:{type:'MARKATA_ACTIVATE_UPDATE'},
    source:{type:'window',url:'https://markata.example/'}
  });
  assert.equal(worker.skips(),1);
});
