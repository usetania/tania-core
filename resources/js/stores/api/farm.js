import { http } from '@/services'

export default {
  ApiCreateFarm: (payload, cbSuccess, cbError) => {
    http.post('farms', payload, cbSuccess, cbError)
  },
  ApiGetFarmTypes: (cbSuccess, cbError) => {
    http.get('farms/types', cbSuccess, cbError)
  },
  ApiCreateReservoir: (farmid, payload, cbSuccess, cbError) => {
    http.post('farms/' + farmid + '/reservoirs', payload, cbSuccess, cbError)
  }
}
