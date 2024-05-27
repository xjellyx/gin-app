<script lang="ts" setup>
import { defineComponent, onMounted, ref } from 'vue';
import 'ol/ol.css';
import { Map, View } from 'ol';
import TileLayer from 'ol/layer/Tile';
import OSM from 'ol/source/OSM';
import { fromLonLat, get as getProjection, toLonLat,transformExtent } from 'ol/proj';
import WMTS from 'ol/source/WMTS';
import WMTSTileGrid from 'ol/tilegrid/WMTS';
import { getTopLeft, getWidth } from 'ol/extent';
import ImageLayer from 'ol/layer/Image';
import Static from 'ol/source/ImageStatic';


const projection = getProjection('EPSG:3857');
const projectionExtent = projection.getExtent();
const size = getWidth(projectionExtent) / 256;
const resolutions = Array.from({ length: 19 });
const matrixIds = Array.from({ length: 19 });
for (let z = 0; z < 19; ++z) {
  // generate resolutions and matrixIds arrays for this WMTS
  resolutions[z] = size / 2 ** z;
  matrixIds[z] = z;
}
const coords = ref<{ lat: number; lon: number } | null>(null);
const mapDiv = ref<HTMLDivElement | null>(null);
mapDiv.value = 'map';

const copyCoords = () => {
  if (coords.value) {
    const text = `${coords.value.lon},${coords.value.lat}`;
    navigator.clipboard.writeText(text).then(() => {});
  }
};

function usgsWmts() {
  return new WMTS({
    attributions: 'Tiles © <a href="https://mrdata.usgs.gov/geology/state/"' + ' target="_blank">USGS</a>',
    url: 'https://mrdata.usgs.gov/mapcache/wmts',
    layer: 'sgmc2',
    matrixSet: 'GoogleMapsCompatible',
    format: 'image/png',
    projection,
    tileGrid: new WMTSTileGrid({
      origin: getTopLeft(projectionExtent),
      resolutions,
      matrixIds
    }),
    style: 'default',
    wrapX: true
  });
}
function addStatic(src: string) {
    return new ImageLayer({
      opacity: 0.7,
      source: new Static({
        url: src,
        projection,
        imageExtent: transformExtent([107.54044733888277,24.59310726828751,107.54140765074139,24.594345681935053],"EPSG:4326","EPSG:3857"),
      })
    })
}
function tiandiWmts() {
  return new WMTS({
    url: 'https://t2.tianditu.gov.cn/img_w/wmts?tk=ce0125d2666eba13871c6f22bc456d3b&TILEMATRIXSET=w',
    service: 'WMTS',
    request: 'GetTile',
    version: '1.0.0',
    layer: 'img',
    style: 'default',
    format: 'tiles',
    projection,
    tileGrid: new WMTSTileGrid({
      origin: getTopLeft(projectionExtent),
      resolutions,
      matrixIds
    }),
    wrapX: true
  });
}

onMounted(() => {
  if (mapDiv.value) {
    const map = new Map({
      target: mapDiv.value,
      layers: [
        new TileLayer({
          source: new OSM()
        }),
        new TileLayer({
          opacity: 0.7,
          source: tiandiWmts()
        }),
        addStatic("http://192.168.3.4:8888/api/v1/resource/5e329bec2a3b25b542cab5af6e4087f.jpg"),
      ],
      view: new View({
        center: fromLonLat([107.54044733888277,24.59310726828751]),
        zoom: 15
      })
    });
    map.on('click', event => {
      const [lon, lat] = toLonLat(event.coordinate);
      coords.value = { lat, lon };
      mapDiv.value = null;
    });
  }
});
</script>

<template>
  <div id="map" class="map">
    <div v-if="coords" class="coordinates">
      <button @click="copyCoords">{{ coords.lon }},{{ coords.lat }}</button>
    </div>
  </div>
</template>

<style scoped>
#map {
  width: 100%;
  height: 100vh;
}
.coordinates {
  position: absolute;
  left: 220px;
  bottom: 60px;
  padding: 10px;
  color: #712d2d;
  z-index: 9999;
}
button {
  opacity: 0.8;
  margin-left: 10px;
}
</style>
