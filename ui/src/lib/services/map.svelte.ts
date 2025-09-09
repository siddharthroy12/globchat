
import type { Map } from "maplibre-gl";

let map: null | Map = $state(null);

export function getMap() {
	return map;
}

export function setMap(m: Map) {
	map = m;
}

