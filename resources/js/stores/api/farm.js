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
    http.get('farms/inventories', cbSuccess, cbError)
  },
  ApiCreateReservoir: (farmid, payload, cbSuccess, cbError) => {
    http.post('farms/' + farmid + '/reservoirs', payload, cbSuccess, cbError)
  },
  ApiFetchReservoir: (farmid, cbSuccess, cbError) => {
    http.get('farms/' + farmid + '/reservoirs', cbSuccess, cbError)
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
  ApiCreateAreaNotes: (areaId, payload, cbSuccess, cbError) => {
    http.post('farms/areas/' + areaId + '/notes' , payload, cbSuccess, cbError)
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
  ApiCreateCropNotes: (cropId, payload, cbSuccess, cbError) => {
    http.post('farms/crops/' + cropId + '/notes' , payload, cbSuccess, cbError)
  },
}
