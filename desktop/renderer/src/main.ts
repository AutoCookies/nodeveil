import { createInitialState } from './state.js';

const state = createInitialState();
const statusNode = document.getElementById('engine-status');

if (statusNode) {
  statusNode.textContent = state.statusText;
}
