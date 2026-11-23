const test=require('node:test');
const assert=require('node:assert/strict');
const share=require('./share-target-utils.js');
test('accepts explicit and extracted HTTP links',()=>{
  assert.equal(share.normalizeSharedText({text:'Look at https://example.com/post).'}).url,'https://example.com/post');
  assert.equal(share.normalizeSharedText({url:'https://example.com/a',text:'https://example.com/b'}).url,'https://example.com/a');
  assert.equal(share.normalizeSharedText({text:'Just words'}).url,'');
});
test('rejects unsafe URLs',()=>{
  for(const url of ['javascript:alert(1)','file:///secret','https://user:pass@host.com','bad url'])
    assert.equal(share.safeURL(url),'');
});
test('rejects unsupported, large or too many images',()=>{
  assert.equal(share.validateImages([{type:'image/jpeg',size:123}]).length,1);
  assert.throws(()=>share.validateImages([{type:'image/heic',size:1}]),/Only PNG/);
  assert.throws(()=>share.validateImages([{type:'image/jpeg',size:15*1024*1024+1}]),/15 MiB/);
  assert.throws(()=>share.validateImages(Array.from({length:7},()=>({type:'image/png',size:1}))),/six/);
});
