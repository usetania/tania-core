import NProgress from 'nprogress';

import * as types from '../../mutation-types';
import FarmApi from '../../api/farm';
import { StubArea } from '../../stubs';

const state = {
  area: Object.assign({}, StubArea),
  areas: [],
  areanotes: [],
};

const getters = {
  getCurrentArea: stateObj => stateObj.area,
  getAllAreas: stateObj => stateObj.areas,
};

const actions = {

  submitArea({ commit, state, getters }, payload) {
    const farm = getters.getCurrentFarm;

    NProgress.start();
    return new Promise((resolve, reject) => {
      const formData = new FormData();
      formData.append('name', payload.name);
      formData.append('size', payload.size);
      formData.append('size_unit', payload.size_unit);
      formData.append('type', payload.type);
      formData.append('location', payload.location);
      formData.append('reservoir_id', payload.reservoir_id);
      formData.append('photo', payload.photo);
      formData.append('farm_id', farm.uid);
      if (payload.uid !== '') {
        FarmApi.ApiUpdateArea(payload.uid, formData, ({ data }) => {
          commit(types.UPDATE_AREA, data.data);
          resolve(payload);
        }, error => reject(error.response));
      } else {
        FarmApi.ApiCreateArea(farm.uid, formData, ({ data }) => {
          commit(types.CREATE_AREA, data.data);
          resolve(payload);
        }, error => reject(error.response));
      }
    });
  },
  fetchAreas({ commit, state, getters }, payload) {
    const farm = getters.getCurrentFarm;

    NProgress.start();
    return new Promise((resolve, reject) => {
      FarmApi.ApiFetchArea(farm.uid, ({ data }) => {
        commit(types.FETCH_AREA, data.data);
        resolve(data);
      }, error => reject(error.response));
    });
  },
  fetchAreaCrops({ commit, state, getters }, areaId) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi.ApiFetchAreaCrop(areaId, ({ data }) => {
        resolve(data);
      }, error => reject(error.response));
    });
  },
  getAreaByUid({ commit, state, getters }, areaId) {
    const farm = getters.getCurrentFarm;

    NProgress.start();
    return new Promise((resolve, reject) => {
      FarmApi.ApiFindAreaByUid(farm.uid, areaId, ({ data }) => {
        resolve(data);
      }, error => reject(error.response));
    });
  },
  createAreaNotes({ commit, state }, payload) {
    NProgress.start();
    return new Promise((resolve, reject) => {
      let areaId = payload.obj_uid;
      FarmApi
        .ApiCreateAreaNotes(areaId, payload, ({ data }) => {
          payload = data.data;
          commit(types.CREATE_AREA_NOTES, payload);
          resolve(payload);
        }, error => reject(error.response));
    });
  },
  deleteAreaNote({ commit, state }, payload) {
    NProgress.start();
    return new Promise((resolve, reject) => {
      let areaId = payload.obj_uid;
      let noteId = payload.uid;
      FarmApi
        .ApiDeleteAreaNotes(areaId, noteId, payload, ({ data }) => {
          payload = data.data;
          commit(types.DELETE_AREA_NOTES, payload);
          resolve(payload);
        }, error => reject(error.response));
    });
  },
};

const mutations = {
  [types.CREATE_AREA](state, payload) {
    state.areas.push(payload);
  },
  [types.UPDATE_AREA](state, payload) {
    const areas = state.areas;
    state.areas = areas.map(area => (area.uid === payload.uid) ? payload : area);
  },
  [types.SET_AREA](state, payload) {
    state.area = payload;
  },
  [types.FETCH_AREA](state, payload) {
    state.areas = payload;
  },
  [types.CREATE_AREA_NOTES](state, payload) {
    state.areanotes.push(payload);
  },
  [types.DELETE_AREA_NOTES](state, payload) {
    state.areanotes.push(payload);
  },
};

export default {
  state, getters, actions, mutations,
};
