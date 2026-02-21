export function initialGraphState() {
  return { dataset:{nodes:[],edges:[],truncated:false,cursor:''}, viewport:{zoom:1,panX:0,panY:0,focusedId:'',selectedId:''}, filters:{depth:1,direction:'BOTH',relationTypes:[],extensions:[]}, loading:{neighborhood:false,error:''}, layout:{running:false,iterations:0,durationMs:0}, render:{fps:0,frameMs:[]} };
}
export function graphReducer(state, action) {
  switch(action.type){
    case 'loadStart': return { ...state, loading:{neighborhood:true,error:''} };
    case 'loadSuccess': return { ...state, loading:{neighborhood:false,error:''}, dataset:{nodes:action.nodes,edges:action.edges,truncated:action.truncated,cursor:action.cursor} };
    case 'setFocus': return { ...state, viewport:{...state.viewport,focusedId:action.nodeId} };
    case 'setFilters': return { ...state, filters: action.filters };
    default: return state;
  }
}
export function selectVisibleGraph(state){
  if(!state.filters.relationTypes.length && !state.filters.extensions.length) return {nodes:state.dataset.nodes, edges:state.dataset.edges};
  const set = new Set(state.dataset.nodes.filter(n=>!state.filters.extensions.length || state.filters.extensions.includes(n.ext)).map(n=>n.id));
  return {nodes:state.dataset.nodes.filter(n=>set.has(n.id)), edges:state.dataset.edges.filter(e=>set.has(e.from_id)&&set.has(e.to_id)&&(!state.filters.relationTypes.length || state.filters.relationTypes.includes(e.relation_type)))};
}
