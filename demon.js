function initMap(){
var searchService;
var suCb=null;
var keyword='';
var latlngBounds = new qq.maps.LatLngBounds();
//调用Poi检索类
searchService = new qq.maps.SearchService({
complete : function(results){
// console.log(keyword,results);
var pois = results.detail.pois;
if(!pois){
suCb&&suCb();
return;
}
for(var i = 0,l = pois.length;i < l; i++){
var _isDtz = pois[i].category.indexOf('地铁站')!=-1;
var _isJt = pois[i].category.indexOf('交通设施')!=-1;
var _isKey=pois[i].name.indexOf('[地铁站]')!=-1;
if(_isDtz&&_isJt&&_isKey){
suCb&&suCb(pois[i]);
return;
}
// latlngBounds.extend(poi.latLng); 
if(i==pois.length-1){
if(results.type!='POI_LIST'){
console.log(keyword+'无地铁站匹配');
}
suCb&&suCb(pois[0]);
}
}
},
error:function(){
keyword=keyword.replace('地铁站','');
searchService.search(keyword);
}
});

function searchKeyword(_city,_key,cb) {
keyword = _key;
var region = _city;
suCb=cb;
// console.log(_key);
searchService.setLocation(region);
searchService.search(keyword);
}
window.searchKeyword=searchKeyword;
}
