import {test} from 'node:test';
import assert from 'node:assert/strict';
import {ToolsClient} from '../dist/index.js';
test('owner paths preserve CAS, false values, cancellation and caller isolation',async()=>{
 const calls=[];const client=new ToolsClient({request:async(...args)=>{calls.push(args);return {items:[]}}});
 const signal=new AbortController().signal;
 await client.listToolSettings(signal);
 const input={enabled:false,expected_revision:0,tool_version:'1'};
 await client.updateToolSetting('calculate',input,signal);
 assert.deepEqual(calls,[['/tools/preferences',{signal}],['/tools/preferences/calculate',{method:'PUT',body:input,signal}]]);
 for(const key of ['../secrets','..','/','a?x=1','a#x',''])assert.throws(()=>client.updateToolSetting(key,input));
 assert.equal(calls.length,2);
});
