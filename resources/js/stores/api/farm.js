import { http } from '@/services'

export default {
  ApiFetchFarm: (cbSuccess, cbError) => {
    http.get('farms', cbSuccess, cbError)
  },
  ApiCreateFarm: (payload, cbSuccess, cbError) => {
    http.post('farms', payload, cbSuccess, cbError)
  },
  ApiGetFarmTypes: (cbSuccess, cbError) => {
    http.get('farms/types', cbSuccess, cbError)
  },
  ApiFetchFarmInventories: (cbSuccess, cbError) => {
    http.get('farms/inventories/materials/available_seed', cbSuccess, cbError)
  },
  ApiCreateReservoir: (farmid, payload, cbSuccess, cbError) => {
    http.post('farms/' + farmid + '/reservoirs', payload, cbSuccess, cbError)
  },
  ApiFetchReservoir: (farmid, cbSuccess, cbError) => {
    http.get('farms/' + farmid + '/reservoirs', cbSuccess, cbError)
  },
  ApiFindReservoirByUid: (farmid, reservoirid, cbSuccess, cbError) => {
    http.get('farms/' + farmid + '/reservoirs/' + reservoirid, cbSuccess, cbError)
  },
  ApiCreateReservoirNotes: (reservoirid, payload, cbSuccess, cbError) => {
    http.post('farms/reservoirs/' + reservoirid + '/notes' , payload, cbSuccess, cbError)
  },
  ApiDeleteReservoirNotes: (reservoirid, noteid, payload, cbSuccess, cbError) => {
    http.delete('farms/reservoirs/' + reservoirid + '/notes/' + noteid, payload, cbSuccess, cbError)
  },
  ApiCreateArea: (farmid, payload, cbSuccess, cbError) => {
    http.post('farms/' + farmid + '/areas', payload, cbSuccess, cbError, {
      'Content-Type': 'multipart/form-data'
    })
  },
  ApiFetchArea: (farmid, cbSuccess, cbError) => {
    http.get('farms/' + farmid + '/areas', cbSuccess, cbError)
  },
  ApiFindAreaByUid: (farmid, areaid, cbSuccess, cbError) => {
    http.get('farms/' + farmid + '/areas/' + areaid, cbSuccess, cbError)
  },
  ApiCreateAreaNotes: (areaid, payload, cbSuccess, cbError) => {
    http.post('farms/areas/' + areaid + '/notes' , payload, cbSuccess, cbError)
  },
  ApiDeleteAreaNotes: (areaid, noteid, payload, cbSuccess, cbError) => {
    http.delete('farms/areas/' + areaid + '/notes/' + noteid, payload, cbSuccess, cbError)
  },
  ApiCreateCrop: (areaid, payload, cbSuccess, cbError) => {
    http.post('farms/areas/' + areaid + '/crops' , payload, cbSuccess, cbError)
  },
  ApiFetchCrop: (farmid, cbSuccess, cbError) => {
    http.get('farms/' + farmid + '/crops', cbSuccess, cbError)
  },
  ApiFindCropByUid: (farmid, cropid, cbSuccess, cbError) => {
    http.get('farms/crops/' + cropid, cbSuccess, cbError)
  },
  ApiCreateCropNotes: (cropid, payload, cbSuccess, cbError) => {
    http.post('farms/crops/' + cropid + '/notes' , payload, cbSuccess, cbError)
  },
  ApiDeleteCropNotes: (cropid, noteid, payload, cbSuccess, cbError) => {
    http.delete('farms/crops/' + cropid + '/notes/' + noteid, payload, cbSuccess, cbError)
  },
  ApiCreateMaterial: (payload, cbSuccess, cbError) => {
    http.post('farms/inventories/materials/' + payload.type, payload, cbSuccess, cbError)
  },
}
