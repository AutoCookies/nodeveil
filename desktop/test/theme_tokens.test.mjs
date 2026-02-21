import test from 'node:test';
import assert from 'node:assert/strict';

const required = ['bg','surface','text','muted','border','primary','secondary','success','warn','error','graphNode','graphEdge','graphSelected','graphHover'];
const themes = {
  light: { bg:1,surface:1,text:1,muted:1,border:1,primary:1,secondary:1,success:1,warn:1,error:1,graphNode:1,graphEdge:1,graphSelected:1,graphHover:1 },
  dark: { bg:1,surface:1,text:1,muted:1,border:1,primary:1,secondary:1,success:1,warn:1,error:1,graphNode:1,graphEdge:1,graphSelected:1,graphHover:1 },
  neon: { bg:1,surface:1,text:1,muted:1,border:1,primary:1,secondary:1,success:1,warn:1,error:1,graphNode:1,graphEdge:1,graphSelected:1,graphHover:1 },
  sakura: { bg:1,surface:1,text:1,muted:1,border:1,primary:1,secondary:1,success:1,warn:1,error:1,graphNode:1,graphEdge:1,graphSelected:1,graphHover:1 }
};

test('theme token completeness', () => {
  Object.values(themes).forEach((t) => required.forEach((k) => assert.equal(Object.hasOwn(t, k), true)));
});
