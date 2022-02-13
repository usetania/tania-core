import { http } from '../../services'

export default {
  ApiLogin: function (payload) {
    return http.login('authorize', payload).then(function(response){
      return response
    }).catch(function() {
      throw new Error()
    })
  },
  ApiChangePassword: (payload, cbSuccess, cbError) => {
    http.post('user/change_password', payload, cbSuccess, cbError)
  },
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
    http.get('farms/inventories/materials/available_plant_type', cbSuccess, cbError)
  },
  ApiCreateReservoir: (farmid, payload, cbSuccess, cbError) => {
    http.post('farms/' + farmid + '/reservoirs', payload, cbSuccess, cbError)
  },
  ApiUpdateReservoir: (reservoirid, payload, cbSuccess, cbError) => {
    http.put('farms/reservoirs/' + reservoirid, payload, cbSuccess, cbError)
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
  ApiUpdateArea: (areaid, payload, cbSuccess, cbError) => {
    http.put('farms/areas/' + areaid, payload, cbSuccess, cbError, {
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
  ApiUpdateCrop: (cropid, payload, cbSuccess, cbError) => {
    http.put('farms/crops/' + cropid, payload, cbSuccess, cbError)
  },
  ApiFetchAreaCrop: (areaid, cbSuccess, cbError) => {
    http.get('farms/areas/' + areaid + '/crops' , cbSuccess, cbError)
  },
  ApiFetchCrop: (farmid, pageid, status, cbSuccess, cbError) => {
    http.get('farms/' + farmid + '/crops?page=' + pageid + '&status=' + status, cbSuccess, cbError)
  },
  ApiFetchArchivedCrop: (farmid, pageid, cbSuccess, cbError) => {
    http.get('farms/' + farmid + '/crops/archives?page=' + pageid, cbSuccess, cbError)
  },
  ApiFindCropByUid: (farmid, cropid, cbSuccess, cbError) => {
    http.get('farms/crops/' + cropid, cbSuccess, cbError)
  },
  ApiMoveCrop: (cropid, payload, cbSuccess, cbError) => {
    http.post('farms/crops/'+cropid+'/move', payload, cbSuccess, cbError)
  },
  ApiHarvestCrop: (cropid, payload, cbSuccess, cbError) => {
    http.post('farms/crops/'+cropid+'/harvest', payload, cbSuccess, cbError)
  },
  ApiDumpCrop: (cropid, payload, cbSuccess, cbError) => {
    http.post('farms/crops/'+cropid+'/dump', payload, cbSuccess, cbError)
  },
  ApiPhotoCrop: (cropid, payload, cbSuccess, cbError) => {
    http.post('farms/crops/'+cropid+'/photos', payload, cbSuccess, cbError, {
      'Content-Type': 'multipart/form-data'
    })
  },
  ApiWaterCrop: (cropid, payload, cbSuccess, cbError) => {
    http.post('farms/crops/' + cropid + '/water' , payload, cbSuccess, cbError)
  },
  ApiCreateCropNotes: (cropid, payload, cbSuccess, cbError) => {
    http.post('farms/crops/' + cropid + '/notes' , payload, cbSuccess, cbError)
  },
  ApiDeleteCropNotes: (cropid, noteid, payload, cbSuccess, cbError) => {
    http.delete('farms/crops/' + cropid + '/notes/' + noteid, payload, cbSuccess, cbError)
  },
  ApiFetchMaterial: (pageid, cbSuccess, cbError) => {
    http.get('farms/inventories/materials?page=' + pageid, cbSuccess, cbError)
  },
  ApiFetchAgrochemicalMaterial: (type, cbSuccess, cbError) => {
    http.get('farms/inventories/materials/simple?type=AGROCHEMICAL&type_detail=' + type, cbSuccess, cbError)
  },
  ApiCreateMaterial: (payload, cbSuccess, cbError) => {
    http.post('farms/inventories/materials/' + payload.type, payload, cbSuccess, cbError)
  },
  ApiUpdateMaterial: (materialid, payload, cbSuccess, cbError) => {
    http.put('farms/inventories/materials/' + payload.type + '/' + materialid, payload, cbSuccess, cbError)
  },
  ApiCreateTask: (payload, cbSuccess, cbError) => {
    http.post('tasks', payload, cbSuccess, cbError)
  },
  ApiUpdateTask: (taskid, payload, cbSuccess, cbError) => {
    http.put('tasks/' + taskid, payload, cbSuccess, cbError)
  },
  ApiFetchTask: (pageid, cbSuccess, cbError) => {
    http.get('tasks?page=' + pageid, cbSuccess, cbError)
  },
  ApiFetchActivity: (cropid, cbSuccess, cbError) => {
    http.get('farms/crops/'+ cropid +'/activities', cbSuccess, cbError)
  },
  ApiFetchCropInformation: (farmid, cbSuccess, cbError) => {
    http.get('farms/'+ farmid +'/crops/information', cbSuccess, cbError)
  },
  ApiFindTasksByDomainAndAssetId: (pageid, domain, assetid, cbSuccess, cbError) => {
    http.get('tasks/search?page=' + pageid + '&domain='+ domain +'&asset_id=' + assetid, cbSuccess, cbError)
  },
  ApiFindTasksByCategoryAndPriorityAndStatus: (pageid, category, priority, status, cbSuccess, cbError) => {
    http.get('tasks/search?page=' + pageid + '&category='+ category +'&priority=' + priority + status, cbSuccess, cbError)
  },
  ApiSetTaskDue: (taskid, cbSuccess, cbError) => {
    http.put('tasks/' + taskid + '/due', {}, cbSuccess, cbError)
  },
  ApiSetTaskCompleted: (taskid, cbSuccess, cbError) => {
    http.put('tasks/' + taskid + '/complete', {}, cbSuccess, cbError)
  },
}
